package services

import (
	"encoding/json"
	"fmt"
	"strings"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type DatasetService struct {
}

func (s *DatasetService) DatasetKeywordTable(ds *models.Dataset) *models.DatasetKeywordTable {
	dskt := new(models.DatasetKeywordTable)
	err := dbengine.Instance().DB.Model(&models.DatasetKeywordTable{}).Where("dataset_id=?", ds.ID).First(dskt).Error
	if err != nil {
		mlog.Errorf("get DatasetKeywordTable by dataset_id = %s failed:%v", ds.ID, err)
		return nil
	}
	return dskt
}

func (s *DatasetService) LatestProcessRule(ds *models.Dataset) *models.DatasetProcessRule {
	// return (
	//     DatasetProcessRule.query.filter(DatasetProcessRule.dataset_id == self.id)
	//     .order_by(DatasetProcessRule.created_at.desc())
	//     .first()
	// )
	dpr := new(models.DatasetProcessRule)
	err := dbengine.Instance().DB.Model(&models.DatasetProcessRule{}).Where("dataset_id=?", ds.ID).Order("created_at DESC").First(dpr).Error
	if err != nil {
		mlog.Errorf("get DatasetProcessRule by dataset_id = %s failed:%v", ds.ID, err)
		return nil
	}
	return dpr
}
func (s *DatasetService) AppCount(ds *models.Dataset) int64 {
	var count int64
	err := dbengine.Instance().DB.Model(&models.AppDatasetJoin{}).Where("dataset_id=?", ds.ID).Count(&count).Error
	if err != nil {
		mlog.Errorf("get AppDatasetJoin by dataset_id = %s failed:%v", ds.ID, err)
		return 0
	}
	return count
}
func (s *DatasetService) DocumentCount(ds *models.Dataset) int64 {
	var count int64
	err := dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id=?", ds.ID).Count(&count).Error
	if err != nil {
		mlog.Errorf("get AppDatasetJoin by dataset_id = %s failed:%v", ds.ID, err)
		return 0
	}
	return count
}
func (s *DatasetService) AvailableDocumentCount(ds *models.Dataset) int64 {
	var count int64
	err := dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id=? and indexing_status = ? and enabled = ? and archived = ?", ds.ID, "indexing_status", true, false).Count(&count).Error
	if err != nil {
		mlog.Errorf("get AppDatasetJoin by dataset_id = %s failed:%v", ds.ID, err)
		return 0
	}
	return count
}
func (s *DatasetService) AvailableSegmentCount(ds *models.Dataset) int64 {
	var count int64
	err := dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id=? and status = ? and enabled = ?", ds.ID, "completed", true).Count(&count).Error
	if err != nil {
		mlog.Errorf("get AppDatasetJoin by dataset_id = %s failed:%v", ds.ID, err)
		return 0
	}
	return count
}
func (s *DatasetService) WordCount(ds *models.Dataset) int {
	var count int
	datas := []*models.Document{}
	err := dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id=?", ds.ID).Find(datas).Error
	if err != nil {
		mlog.Errorf("get AppDatasetJoin by dataset_id = %s failed:%v", ds.ID, err)
		return 0
	}
	for _, v := range datas {
		count += v.WordCount
	}
	return count
}
func (s *DatasetService) DocForm(ds *models.Dataset) string {
	doc := new(models.Document)
	err := dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id=?", ds.ID).First(doc).Error
	if err != nil {
		mlog.Errorf("get AppDatasetJoin by dataset_id = %s failed:%v", ds.ID, err)
		return ""
	}
	return doc.DocForm

	// document = db.session.query(Document).filter(Document.dataset_id == self.id).first()
	// if document:
	//     return document.doc_form
	// return None

}

func (s *DatasetService) Tags(ds *models.Dataset) []*models.Tag {
	tbs := []*models.TagBinding{}
	err := dbengine.Instance().DB.Model(&models.TagBinding{}).Where("target_id=? and tenant_id=? and type=?", ds.ID, ds.TenantID, "knowledge").Preload("Tag").Find(tbs).Error
	if err != nil {
		mlog.Errorf("get TagBinding by dataset_id = %s failed:%v", ds.ID, err)
		return nil
	}
	var tags []*models.Tag
	for _, tb := range tbs {
		if tags == nil {
			tags = []*models.Tag{}
		}
		tags = append(tags, tb.Tag)
	}
	return tags
	// tags = (
	//     db.session.query(Tag)
	//     .join(TagBinding, Tag.id == TagBinding.tag_id)
	//     .filter(
	//         TagBinding.target_id == self.id,
	//         TagBinding.tenant_id == self.tenant_id,
	//         Tag.tenant_id == self.tenant_id,
	//         Tag.type == "knowledge",
	//     )
	//     .all()
	// )

	// return tags or []
}
func (s *DatasetService) ExternalKnowledgeInfo(ds *models.Dataset) map[string]any {
	if ds.Provider != "external" {
		return nil
	}
	external_knowledge_binding := new(models.ExternalKnowledgeBinding)
	err := dbengine.Instance().DB.Model(&models.ExternalKnowledgeBinding{}).Where("dataset_id=?", ds.ID).Preload("ExternalKnowledgeApi").First(external_knowledge_binding).Error
	if err != nil {
		mlog.Errorf("get ExternalKnowledgeBinding by dataset_id = %s failed:%v", ds.ID, err)
		return nil
	}
	// external_knowledge_binding = (
	//     db.session.query(ExternalKnowledgeBindings).filter(ExternalKnowledgeBindings.dataset_id == self.id).first()
	// )
	// if not external_knowledge_binding:
	//     return None

	// external_knowledge_api = (
	//     db.session.query(ExternalKnowledgeApis)
	//     .filter(ExternalKnowledgeApis.id == external_knowledge_binding.external_knowledge_api_id)
	//     .first()
	// )
	if external_knowledge_binding.ExternalKnowledgeApi == nil {
		return nil
	}
	endpoint := ""
	settings := map[string]any{}
	err = json.Unmarshal([]byte(external_knowledge_binding.ExternalKnowledgeApi.Settings), &settings)
	if err != nil {
		mlog.Warningf("json unmarshal Settings(%s) failed:%v", external_knowledge_binding.ExternalKnowledgeApi.Settings, err)
	} else {
		if _, ok := settings["endpoint"]; ok {
			if _, ok := settings["endpoint"].(string); ok {
				endpoint = settings["endpoint"].(string)
			}
		}
	}
	return map[string]any{
		"external_knowledge_id":           external_knowledge_binding.ExternalKnowledgeID,
		"external_knowledge_api_id":       external_knowledge_binding.ExternalKnowledgeApi.ID,
		"external_knowledge_api_name":     external_knowledge_binding.ExternalKnowledgeApi.Name,
		"external_knowledge_api_endpoint": endpoint,
	}
}
func (s *DatasetService) GenCollectionNameByID(dataset_id string) string {
	normalized_dataset_id := strings.ReplaceAll(dataset_id, "-", "_")
	return fmt.Sprintf("Vector_index_%s_Node", normalized_dataset_id)
}

func (s *DatasetService) GetByIDAndTenantID(id, tenant_id string) (*models.Dataset, error) {
	ds := new(models.Dataset)
	err := dbengine.Instance().DB.Where("id = ? and tenant_id = ?", id, tenant_id).First(ds).Error
	if err != nil {
		mlog.Errorf("get Dataset failed:%v", err)
		return nil, err
	}
	return ds, nil
}
