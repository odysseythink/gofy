package indextool

import (
	"slices"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	ragentities "mlib.com/gofy/server/entities/rag"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	indexprocessorenumtypes "mlib.com/gofy/server/enum_types/rag/index_processor"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
)

type DatasetIndexToolCallbackHandler struct {
	_queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage]
	_app_id        string
	_message_id    string
	_user_id       string
	_invoke_from   appenumtypes.InvokeFrom
}

func New(
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage], app_id string, message_id string, user_id string, invoke_from appenumtypes.InvokeFrom,
) *DatasetIndexToolCallbackHandler {
	return &DatasetIndexToolCallbackHandler{
		_queue_manager: queue_manager,
		_app_id:        app_id,
		_message_id:    message_id,
		_user_id:       user_id,
		_invoke_from:   invoke_from,
	}
}
func (handler *DatasetIndexToolCallbackHandler) OnQuery(query string, dataset_id string) {
	dataset_query := &models.DatasetQuery{
		ID:          uuid.NewV4().String(),
		DatasetID:   dataset_id,
		Content:     query,
		Source:      "app",
		SourceAppID: handler._app_id,
		CreatedBy:   handler._user_id,
	}
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_EXPLORE, appenumtypes.InvokeFrom_DEBUGGER}, handler._invoke_from) {
		dataset_query.CreatedByRole = "account"
	} else {
		dataset_query.CreatedByRole = "end_user"
	}
	dbengine.Instance().DB.Create(dataset_query)
}
func (handler *DatasetIndexToolCallbackHandler) OnToolEnd(documents []*ragentities.Document) {
	for _, document := range documents {
		if len(document.Metadata) > 0 {
			mlog.Infof("------document.metadata=%#v", document.Metadata)
			document_id := mapstruct.Get(document.Metadata, "document_id", "")
			dataset_document := new(models.Document)
			err := dbengine.Instance().DB.Where(&models.Document{}).Where("id =?", document_id).First(dataset_document).Error
			if err != nil {
				mlog.Errorf("get Document failed:%v", err)
				dataset_document = nil
			}

			if dataset_document == nil {
				mlog.Warningf(
					"Expected DatasetDocument record to exist, but none was found, document_id=%s",
					document_id,
				)
				continue
			}
			if dataset_document.DocForm == string(indexprocessorenumtypes.Index_PARENT_CHILD_INDEX) {
				child_chunk := new(models.ChildChunk)
				err := dbengine.Instance().DB.Model(&models.ChildChunk{}).Where("index_node_id =? and dataset_id = ? and document_id = ?", mapstruct.Get(document.Metadata, "doc_id", ""), dataset_document.DatasetID, dataset_document.ID).First(child_chunk).Error
				if err != nil {
					mlog.Errorf("get ChildChunk failed:%v", err)
					child_chunk = nil
				}
				if child_chunk != nil {
					dbengine.Instance().DB.Exec("UPDATE document_segments SET hit_count = hit_count + 1 WHERE id = ?", child_chunk.SegmentID)
				}
			} else {
				if _, ok := document.Metadata["dataset_id"]; ok {
					dbengine.Instance().DB.Exec("UPDATE document_segments SET hit_count = hit_count + 1 WHERE index_node_id = ? and dataset_id =?", mapstruct.Get(document.Metadata, "doc_id", ""), mapstruct.Get(document.Metadata, "dataset_id", ""))
				} else {
					dbengine.Instance().DB.Exec("UPDATE document_segments SET hit_count = hit_count + 1 WHERE index_node_id = ?", mapstruct.Get(document.Metadata, "doc_id", ""))
				}
			}
		}
	}
}

// TODO(-LAN-): Improve type check
func (handler *DatasetIndexToolCallbackHandler) ReturnRetrieverResourceInfo(resource []*ragentities.RetrievalSourceMetadata) {
	handler._queue_manager.Publish(&appqueueentities.QueueRetrieverResourcesEvent{RetrieverResources: resource}, appenumtypes.PublishFrom_APPLICATION_MANAGER)
}
