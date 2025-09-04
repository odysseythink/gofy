package indextool

import (
	"slices"

	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/core/app/runner/base"
	agenttoolcbhandler "mlib.com/gofy/server/core/callback_handler/agent_tool"
	"mlib.com/gofy/server/core/exceptions"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	modelmgr "mlib.com/gofy/server/core/manageres/model_manager"
	"mlib.com/gofy/server/core/memory"
	dbengine "mlib.com/gofy/server/db_engine"
	agententities "mlib.com/gofy/server/entities/agent"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	ragentities "mlib.com/gofy/server/entities/rag"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/main/sandbox/runner"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)
type DatasetIndexToolCallbackHandler struct{
_queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage]
_app_id string
_message_id string
_user_id string
_invoke_from  appenumtypes.InvokeFrom
}

func New(
        queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage], app_id  string, message_id  string, user_id  string, invoke_from  appenumtypes.InvokeFrom,
    ) *DatasetIndexToolCallbackHandler{
		return &DatasetIndexToolCallbackHandler{
        _queue_manager: queue_manager,
        _app_id: app_id,
        _message_id: message_id,
        _user_id: user_id,
        _invoke_from: invoke_from,
		}
}
func (handler *DatasetIndexToolCallbackHandler) OnQuery(query  string, dataset_id  string) {
        dataset_query := &models.DatasetQuery{
			ID: uuid.NewV4().String(),
            DatasetID:dataset_id,
            Content:query,
            Source:"app",
            SourceAppID:handler._app_id,
             CreatedBy:handler._user_id,
        }
		if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_EXPLORE, appenumtypes.InvokeFrom_DEBUGGER}, handler._invoke_from){
dataset_query.CreatedByRole = "account"
		} else {
dataset_query.CreatedByRole = "end_user"
		}
		dbengine.Instance().DB.Create(dataset_query)
}
func (handler *DatasetIndexToolCallbackHandler) OnToolEnd(documents []*ragentities.Document) {
        for _, document := range documents{
            if len(document.Metadata) > 0{
                document_id := mapstruct.Get(document.metadata,"document_id", "")
                dataset_document := new(models.Document)
				err := dbengine.Instance().DB.Where(&models.Document{}).Where("id =?", document_id).First(dataset_document).Error
				if err != nil {
					mlog.Errorf("get Document failed:%v", err)
					dataset_document = nil
				}

                if dataset_document == nil{
                    mlog.Warningf(
                        "Expected DatasetDocument record to exist, but none was found, document_id=%s",
                        document_id,
                    )
                    continue
				}
                if dataset_document.doc_form == IndexType.PARENT_CHILD_INDEX{
                    child_chunk = (
                        db.session.query(ChildChunk)
                        .where(
                            ChildChunk.index_node_id == document.metadata["doc_id"],
                            ChildChunk.dataset_id == dataset_document.dataset_id,
                            ChildChunk.document_id == dataset_document.id,
                        )
                        .first()
                    )
                    if child_chunk{
                        segment = (
                            db.session.query(DocumentSegment)
                            .where(DocumentSegment.id == child_chunk.segment_id)
                            .update(
                                {DocumentSegment.hit_count: DocumentSegment.hit_count + 1}, synchronize_session=False
                            )
                        )
					}
                }else{
                    query = db.session.query(DocumentSegment).where(
                        DocumentSegment.index_node_id == document.metadata["doc_id"]
                    )

                    if "dataset_id" in document.metadata{
                        query = query.where(DocumentSegment.dataset_id == document.metadata["dataset_id"])
					}
                    // add hit count to document segment
                    query.update({DocumentSegment.hit_count: DocumentSegment.hit_count + 1}, synchronize_session=False)

                db.session.commit()
			}
		}
}
    // TODO(-LAN-): Improve type check
func (handler *DatasetIndexToolCallbackHandler) return_retriever_resource_info(resource: Sequence[RetrievalSourceMetadata]){
        handler._queue_manager.publish(
            QueueRetrieverResourcesEvent(retriever_resources=resource), PublishFrom.APPLICATION_MANAGER
        )
		}