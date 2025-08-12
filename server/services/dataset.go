package services

import (
	"encoding/json"
	"fmt"
	"strings"

	dbengine "mlib.com/gofy/server/db_engine"
	enumtypes "mlib.com/gofy/server/enum_types"
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
	// if document{
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
	// if not external_knowledge_binding{
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
func (s *DatasetService) GetDatasetsByIDs(ids []string, tenant_id string) ([]*models.Dataset, int64) {
	var datas []*models.Dataset
	err := dbengine.Instance().DB.Where("id in ? and tenant_id = ?", ids, tenant_id).Find(&datas).Error
	if err != nil {
		mlog.Errorf("get Dataset failed:%v", err)
		return nil, 0
	}

	return datas, int64(len(datas))
}
func (s *DatasetService) GetDatasets(page int32, per_page int32, tenant_id string, user *models.Account, search string, tag_ids []string, include_all bool) ([]*models.Dataset, int64) {
	db := dbengine.Instance().DB.Model(&models.Dataset{}).Where("tenant_id = ?", tenant_id).Order("created_at DESC")

	if user != nil {
		// get permitted dataset ids
		var dataset_permissions []*models.DatasetPermission
		err := dbengine.Instance().DB.Where("account_id = ? and tenant_id = ?", user.ID, tenant_id).Find(&dataset_permissions).Error
		if err != nil {
			mlog.Errorf("get DatasetPermission failed:%v", err)
		}
		permitted_dataset_ids := []string{}
		if len(dataset_permissions) > 0 {
			for _, v := range dataset_permissions {
				permitted_dataset_ids = append(permitted_dataset_ids, v.DatasetID)
			}
		}

		if enumtypes.TenantAccountRole(user.CurrentRole()) == enumtypes.TenantAccountRole_DATASET_OPERATOR {
			// only show datasets that the user has permission to access
			if len(permitted_dataset_ids) > 0 {
				db = db.Where("id in ?", permitted_dataset_ids)
			} else {
				return []*models.Dataset{}, 0
			}
		} else {
			if enumtypes.TenantAccountRole(user.CurrentRole()) == enumtypes.TenantAccountRole_OWNER || !include_all {
				// show all datasets that the user has permission to access
				if len(permitted_dataset_ids) > 0 {
					db = db.Where("permission = ? or (permission = ? and created_by = ?) or (permission = ? and id in ?)", enumtypes.DatasetPermission_ALL_TEAM, enumtypes.DatasetPermission_ONLY_ME, user.ID, enumtypes.DatasetPermission_PARTIAL_TEAM, permitted_dataset_ids)
				} else {
					db = db.Where("permission = ? or (permission = ? and created_by = ?)", enumtypes.DatasetPermission_ALL_TEAM, enumtypes.DatasetPermission_ONLY_ME, user.ID)
				}
			}
		}
	} else {
		// if no user, only show datasets that are shared with all team members
		db = db.Where("permission = ?", enumtypes.DatasetPermission_ALL_TEAM)
	}
	if search != "" {
		db = db.Where("name LIKE ?", "%"+search+"%")
	}
	if len(tag_ids) > 0 {
		target_ids := ServiceGroupApp.Tag.GetTargetIDsByIDs("knowledge", tenant_id, tag_ids)
		if len(target_ids) > 0 {
			db = db.Where("id in ?", target_ids)
		} else {
			return []*models.Dataset{}, 0
		}
	}
	var total int64

	err := db.Count(&total).Error
	if err != nil {
		mlog.Errorf("count dataset failed:%v", err)
		return []*models.Dataset{}, 0
	}
	offset := int(per_page * (page - 1))
	var datas []*models.Dataset
	err = db.Limit(int(per_page)).Offset(offset).Scan(&datas).Error
	if err != nil {
		mlog.Errorf("get dataset failed:%v", err)
		return []*models.Dataset{}, 0
	}

	return datas, total
}

func (s *DatasetService) GetDatasetPartialMemberList(dataset_id string) []string {
	var datas []string
	err := dbengine.Instance().DB.Model(&models.DatasetPermission{}).Select("account_id").Where("dataset_id = ?", dataset_id).Find(&datas).Error
	if err != nil {
		mlog.Errorf("get DatasetPermission failed:%v", err)
		return nil
	}

	return datas
}
