package services

import (
	"encoding/json"
	"fmt"
	"strings"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	enumtypes "github.com/odysseythink/gofy/backend/enum_types"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
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
func (s *DatasetService) GetExternalKnowledgeAPIs(
	page int32, per_page int32, tenant_id string, search string,
) ([]*models.ExternalKnowledgeApi, int64) {
	db := dbengine.Instance().DB.Model(&models.ExternalKnowledgeApi{}).Where("tenant_id = ?", tenant_id).Order("created_at DESC")
	if search != "" {
		db = db.Where("name like ?", "%"+search+"%")
	}
	var total int64

	err := db.Count(&total).Error
	if err != nil {
		mlog.Errorf("count ExternalKnowledgeApi failed:%v", err)
		return []*models.ExternalKnowledgeApi{}, 0
	}
	offset := int(per_page * (page - 1))
	var datas []*models.ExternalKnowledgeApi
	err = db.Limit(int(per_page)).Offset(offset).Scan(&datas).Error
	if err != nil {
		mlog.Errorf("get ExternalKnowledgeApi failed:%v", err)
		return []*models.ExternalKnowledgeApi{}, 0
	}

	return datas, total
}

func (s *DatasetService) CreateEmptyDataset(tenantID, name string, description string, indexingTechnique string, account *models.Account, permission string, provider string, embeddingModelProvider string, embeddingModelName string) (*models.Dataset, error) {
	// Check for duplicate name
	var count int64
	dbengine.Instance().DB.Model(&models.Dataset{}).Where("tenant_id = ? AND name = ?", tenantID, name).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("dataset with name '%s' already exists", name)
	}

	if permission == "" {
		permission = "only_me"
	}
	if provider == "" {
		provider = "vendor"
	}

	dataset := &models.Dataset{
		ID:                     uuid.NewV4().String(),
		TenantID:               tenantID,
		Name:                   name,
		Description:            description,
		Provider:               provider,
		Permission:             permission,
		IndexingTechnique:      indexingTechnique,
		CreatedBy:              account.ID,
		EmbeddingModelProvider: embeddingModelProvider,
		EmbeddingModel:         embeddingModelName,
	}

	if err := dbengine.Instance().DB.Create(dataset).Error; err != nil {
		return nil, err
	}
	return dataset, nil
}

func (s *DatasetService) DeleteDataset(datasetID string, user *models.Account) error {
	dataset, err := s.GetByIDAndTenantID(datasetID, user.CurrentTenantID())
	if err != nil {
		return err
	}
	if dataset == nil {
		return fmt.Errorf("dataset not found")
	}

	if err := s.CheckDatasetPermission(dataset, user); err != nil {
		return err
	}

	return dbengine.Instance().DB.Delete(dataset).Error
}

func (s *DatasetService) CheckDatasetPermission(dataset *models.Dataset, user *models.Account) error {
	if dataset.TenantID != user.CurrentTenantID() {
		return fmt.Errorf("no permission to access this dataset")
	}

	// Check user role - get tenant account join
	var join models.TenantAccountJoin
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND account_id = ?", dataset.TenantID, user.ID).First(&join).Error; err != nil {
		return fmt.Errorf("no permission to access this dataset")
	}

	if join.Role == "owner" || join.Role == "admin" {
		return nil
	}

	if dataset.Permission == "only_me" && dataset.CreatedBy != user.ID {
		return fmt.Errorf("no permission to access this dataset")
	}

	if dataset.Permission == "partial_members" && dataset.CreatedBy != user.ID {
		var permCount int64
		dbengine.Instance().DB.Model(&models.DatasetPermission{}).Where("dataset_id = ? AND account_id = ?", dataset.ID, user.ID).Count(&permCount)
		if permCount == 0 {
			return fmt.Errorf("no permission to access this dataset")
		}
	}

	return nil
}

func (s *DatasetService) GetProcessRules(datasetID string) map[string]any {
	var rule models.DatasetProcessRule
	err := dbengine.Instance().DB.Where("dataset_id = ?", datasetID).Order("created_at DESC").First(&rule).Error
	if err != nil {
		// Return default rules
		return map[string]any{
			"mode": "automatic",
			"rules": map[string]any{
				"pre_processing_rules": []map[string]any{
					{"id": "remove_extra_spaces", "enabled": true},
					{"id": "remove_urls_emails", "enabled": false},
				},
				"segmentation": map[string]any{
					"separator":  "###",
					"max_tokens": 500,
				},
			},
		}
	}
	return map[string]any{
		"mode":  rule.Mode,
		"rules": rule.Rules,
	}
}

func (s *DatasetService) GetRelatedApps(datasetID string) []*models.AppDatasetJoin {
	var joins []*models.AppDatasetJoin
	dbengine.Instance().DB.Where("dataset_id = ?", datasetID).Order("created_at DESC").Find(&joins)
	return joins
}

func (s *DatasetService) GetDatasetQueries(datasetID string, page, perPage int) ([]*models.DatasetQuery, int64) {
	if perPage > 100 {
		perPage = 100
	}
	var queries []*models.DatasetQuery
	var total int64
	db := dbengine.Instance().DB.Where("dataset_id = ?", datasetID)
	db.Model(&models.DatasetQuery{}).Count(&total)
	db.Order("created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&queries)
	return queries, total
}

func (s *DatasetService) DatasetUseCheck(datasetID string) bool {
	var count int64
	dbengine.Instance().DB.Model(&models.AppDatasetJoin{}).Where("dataset_id = ?", datasetID).Count(&count)
	return count > 0
}

func (s *DatasetService) UpdateDatasetApiStatus(datasetID string, status bool) error {
	return dbengine.Instance().DB.Model(&models.Dataset{}).Where("id = ?", datasetID).Update("is_api_enabled", status).Error
}

func (s *DatasetService) GetDatasetAutoDisableLogs(datasetID string) []*models.DatasetAutoDisableLog {
	var logs []*models.DatasetAutoDisableLog
	dbengine.Instance().DB.Where("dataset_id = ?", datasetID).Order("created_at DESC").Find(&logs)
	return logs
}

// === Dataset Update ===

func (s *DatasetService) UpdateDataset(datasetID string, user *models.Account, data map[string]any) (*models.Dataset, error) {
	dataset, err := s.GetByIDAndTenantID(datasetID, user.CurrentTenantID())
	if err != nil || dataset == nil {
		return nil, fmt.Errorf("dataset not found")
	}
	if err := s.CheckDatasetPermission(dataset, user); err != nil {
		return nil, err
	}
	// Check duplicate name
	if name, ok := data["name"].(string); ok && name != dataset.Name {
		var count int64
		dbengine.Instance().DB.Model(&models.Dataset{}).Where("tenant_id = ? AND name = ? AND id != ?", dataset.TenantID, name, dataset.ID).Count(&count)
		if count > 0 {
			return nil, fmt.Errorf("dataset with name '%s' already exists", name)
		}
	}
	updates := map[string]any{}
	allowedFields := []string{"name", "description", "permission", "indexing_technique", "retrieval_model", "embedding_model_provider", "embedding_model"}
	for _, field := range allowedFields {
		if v, ok := data[field]; ok {
			updates[field] = v
		}
	}
	if len(updates) > 0 {
		updates["updated_by"] = user.ID
		if err := dbengine.Instance().DB.Model(dataset).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return dataset, nil
}

func (s *DatasetService) CheckDatasetModelSetting(dataset *models.Dataset) error {
	if dataset.IndexingTechnique == "high_quality" {
		if dataset.EmbeddingModelProvider == "" || dataset.EmbeddingModel == "" {
			return fmt.Errorf("embedding model not configured for high quality dataset")
		}
	}
	return nil
}

func (s *DatasetService) GetDatasetCollectionBinding(providerName, modelName, collectionType string) *models.DatasetCollectionBinding {
	var binding models.DatasetCollectionBinding
	if err := dbengine.Instance().DB.Where("provider_name = ? AND model_name = ? AND type = ?", providerName, modelName, collectionType).First(&binding).Error; err != nil {
		return nil
	}
	return &binding
}

func (s *DatasetService) UpdatePartialMemberList(tenantID, datasetID string, userIDs []string) error {
	// Clear existing permissions
	dbengine.Instance().DB.Where("dataset_id = ?", datasetID).Delete(&models.DatasetPermission{})
	// Add new permissions
	for _, userID := range userIDs {
		perm := &models.DatasetPermission{
			ID:        uuid.NewV4().String(),
			DatasetID: datasetID,
			AccountID: userID,
			TenantID:  tenantID,
		}
		dbengine.Instance().DB.Create(perm)
	}
	return nil
}

func (s *DatasetService) ClearPartialMemberList(datasetID string) error {
	return dbengine.Instance().DB.Where("dataset_id = ?", datasetID).Delete(&models.DatasetPermission{}).Error
}

func (s *DatasetService) CheckDatasetOperatorPermission(user *models.Account, dataset *models.Dataset) error {
	if dataset == nil || user == nil {
		return fmt.Errorf("invalid parameters")
	}
	role := ""
	var join models.TenantAccountJoin
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND account_id = ?", dataset.TenantID, user.ID).First(&join).Error; err == nil {
		role = join.Role
	}
	if role == "owner" || role == "admin" || role == "editor" || role == "dataset_operator" {
		return nil
	}
	return fmt.Errorf("no permission to operate on this dataset")
}

// CreateProcessRule creates a new dataset process rule.
func (s *DatasetService) CreateProcessRule(datasetID, mode, rules string) (*models.DatasetProcessRule, error) {
	rule := &models.DatasetProcessRule{
		ID:        uuid.NewV4().String(),
		DatasetID: datasetID,
		Mode:      mode,
		Rules:     rules,
	}
	if err := dbengine.Instance().DB.Create(rule).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

// CheckDocumentsUploadQuota validates that the tenant hasn't exceeded upload limits.
func (s *DatasetService) CheckDocumentsUploadQuota(tenantID string, count int) error {
	// Default limit: 50 documents per dataset
	// TODO: Check against billing/subscription plan
	if count > 500 {
		return fmt.Errorf("batch upload limit exceeded (max 500)")
	}
	return nil
}

// DataSourceArgsValidate validates document data source arguments.
func (s *DatasetService) DataSourceArgsValidate(dataSourceType string, dataSourceInfoList []map[string]any) error {
	if len(dataSourceInfoList) == 0 {
		return fmt.Errorf("data_source_info_list is required")
	}
	switch dataSourceType {
	case "upload_file":
		for _, info := range dataSourceInfoList {
			if _, ok := info["upload_file_id"]; !ok {
				return fmt.Errorf("upload_file_id is required for upload_file source")
			}
		}
	case "notion_import":
		for _, info := range dataSourceInfoList {
			if _, ok := info["page_id"]; !ok {
				return fmt.Errorf("page_id is required for notion_import source")
			}
		}
	case "website_crawl":
		for _, info := range dataSourceInfoList {
			if _, ok := info["url"]; !ok {
				return fmt.Errorf("url is required for website_crawl source")
			}
		}
	default:
		return fmt.Errorf("unsupported data source type: %s", dataSourceType)
	}
	return nil
}

// ProcessRuleArgsValidate validates process rule arguments.
func (s *DatasetService) ProcessRuleArgsValidate(mode string, rules string) error {
	if mode == "" {
		return fmt.Errorf("process rule mode is required")
	}
	validModes := []string{"automatic", "custom", "hierarchical"}
	valid := false
	for _, m := range validModes {
		if mode == m {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid process rule mode: %s", mode)
	}
	if mode == "custom" && rules == "" {
		return fmt.Errorf("rules are required for custom mode")
	}
	return nil
}

// === Dataset Update Internals ===

func (s *DatasetService) updateExternalDataset(dataset *models.Dataset, data map[string]any) {
	if rm, ok := data["external_retrieval_model"]; ok {
		rmJSON, _ := json.Marshal(rm)
		dbengine.Instance().DB.Model(dataset).Update("retrieval_model", string(rmJSON))
	}
}

func (s *DatasetService) updateInternalDataset(dataset *models.Dataset, data map[string]any, user *models.Account) error {
	// Remove external-only fields
	delete(data, "external_knowledge_api_id")
	delete(data, "external_knowledge_id")
	delete(data, "external_retrieval_model")

	updates := map[string]any{}
	allowedFields := map[string]string{
		"name": "name", "description": "description", "permission": "permission",
		"indexing_technique": "indexing_technique", "retrieval_model": "retrieval_model",
	}
	for key, col := range allowedFields {
		if v, ok := data[key]; ok && v != nil {
			if key == "retrieval_model" {
				vJSON, _ := json.Marshal(v)
				updates[col] = string(vJSON)
			} else {
				updates[col] = v
			}
		}
	}

	// Handle indexing technique change
	action := s.handleIndexingTechniqueChange(dataset, data, updates)

	// Handle embedding model update
	if action == "" {
		s.handleEmbeddingModelUpdate(dataset, data, updates)
	}

	if len(updates) > 0 {
		updates["updated_by"] = user.ID
		return dbengine.Instance().DB.Model(dataset).Updates(updates).Error
	}
	return nil
}

func (s *DatasetService) handleIndexingTechniqueChange(dataset *models.Dataset, data map[string]any, updates map[string]any) string {
	newTechnique, ok := data["indexing_technique"].(string)
	if !ok || newTechnique == dataset.IndexingTechnique {
		return ""
	}

	if newTechnique == "economy" {
		// Switching to economy: clear embedding settings
		updates["embedding_model_provider"] = ""
		updates["embedding_model"] = ""
		updates["collection_binding_id"] = ""
		return "remove"
	}

	// Switching to high_quality
	if provider, ok := data["embedding_model_provider"].(string); ok {
		updates["embedding_model_provider"] = provider
	}
	if model, ok := data["embedding_model"].(string); ok {
		updates["embedding_model"] = model
	}
	return "add"
}

func (s *DatasetService) handleEmbeddingModelUpdate(dataset *models.Dataset, data map[string]any, updates map[string]any) {
	if dataset.IndexingTechnique != "high_quality" {
		return
	}
	newProvider, hasProvider := data["embedding_model_provider"].(string)
	newModel, hasModel := data["embedding_model"].(string)

	if hasProvider && newProvider != dataset.EmbeddingModelProvider {
		updates["embedding_model_provider"] = newProvider
	}
	if hasModel && newModel != dataset.EmbeddingModel {
		updates["embedding_model"] = newModel
	}
}

// CheckDocForm validates that document form is compatible with dataset.
func (s *DatasetService) CheckDocForm(dataset *models.Dataset, docForm string) error {
	currentDocForm := dataset.DocForm()
	if currentDocForm != "" && currentDocForm != docForm {
		return fmt.Errorf("document form '%s' is not compatible with dataset form '%s'", docForm, currentDocForm)
	}
	return nil
}

// CheckEmbeddingModelSetting validates embedding model is configured for high-quality datasets.
func (s *DatasetService) CheckEmbeddingModelSetting(tenantID, embeddingModelProvider, embeddingModel string) error {
	if embeddingModelProvider == "" {
		return fmt.Errorf("embedding model provider is required")
	}
	if embeddingModel == "" {
		return fmt.Errorf("embedding model name is required")
	}
	// TODO: Validate model exists via model_runtime provider
	return nil
}

// CheckRerankingModelSetting validates reranking model is configured.
func (s *DatasetService) CheckRerankingModelSetting(tenantID, rerankingModelProvider, rerankingModel string) error {
	if rerankingModelProvider == "" || rerankingModel == "" {
		return fmt.Errorf("reranking model provider and name are required")
	}
	// TODO: Validate model exists via model_runtime provider
	return nil
}

// CheckIsMultimodalModel checks if a model supports multimodal input.
func (s *DatasetService) CheckIsMultimodalModel(tenantID, modelProvider, model string) bool {
	// TODO: Query model_runtime for model capabilities
	return false
}

// DocumentCreateArgsValidate validates document creation arguments.
func (s *DatasetService) DocumentCreateArgsValidate(dataSourceType string, processRuleMode string) error {
	if dataSourceType == "" && processRuleMode == "" {
		return fmt.Errorf("data_source or process_rule is required")
	}
	return nil
}

// EstimateArgsValidate validates estimation request arguments.
func (s *DatasetService) EstimateArgsValidate(infoList []map[string]any, processRule map[string]any) error {
	if len(infoList) == 0 {
		return fmt.Errorf("info_list is required")
	}
	if processRule != nil {
		mode, _ := processRule["mode"].(string)
		if mode == "" {
			return fmt.Errorf("process_rule.mode is required")
		}
		validModes := map[string]bool{"automatic": true, "custom": true, "hierarchical": true}
		if !validModes[mode] {
			return fmt.Errorf("invalid process_rule mode: %s", mode)
		}
	}
	return nil
}
