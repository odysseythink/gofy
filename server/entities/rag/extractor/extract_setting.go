package extractor

import (
	"encoding/json"

	extractorenumtypes "mlib.com/gofy/server/enum_types/rag/extractor"
	"mlib.com/gofy/server/models"
)

// NotionInfo represents the import info from Notion
type NotionInfo struct {
	NotionWorkspaceID string           `json:"notion_workspace_id"`
	NotionObjID       string           `json:"notion_obj_id"`
	NotionPageType    string           `json:"notion_page_type"`
	Document          *models.Document `json:"document,omitempty"`
	TenantID          string           `json:"tenant_id"`
}

// Helper function to create NotionInfo
func NewNotionInfo(data map[string]any) *NotionInfo {
	info := &NotionInfo{}
	err := json.Unmarshal(dataToJSON(data), info)
	if err != nil {
		panic(err)
	}
	return info
}

// WebsiteInfo represents the import info from a website
type WebsiteInfo struct {
	Provider        string `json:"provider"`
	JobID           string `json:"job_id"`
	URL             string `json:"url"`
	Mode            string `json:"mode"`
	TenantID        string `json:"tenant_id"`
	OnlyMainContent bool   `json:"only_main_content"`
}

// Helper function to create WebsiteInfo
func NewWebsiteInfo(data map[string]any) (*WebsiteInfo, error) {
	info := &WebsiteInfo{}
	err := json.Unmarshal(dataToJSON(data), info)
	return info, err
}

// ExtractSetting represents the settings for extraction
type ExtractSetting struct {
	DatasourceType extractorenumtypes.DatasourceType `json:"datasource_type"`
	UploadFile     *models.UploadFile                `json:"upload_file,omitempty"`
	NotionInfo     *NotionInfo                       `json:"notion_info,omitempty"`
	WebsiteInfo    *WebsiteInfo                      `json:"website_info,omitempty"`
	DocumentModel  *string                           `json:"document_model,omitempty"`
}

// Helper function to create ExtractSetting
func NewExtractSetting(data map[string]any) (*ExtractSetting, error) {
	setting := &ExtractSetting{}
	err := json.Unmarshal(dataToJSON(data), setting)
	return setting, err
}

// Helper function to convert data to JSON
func dataToJSON(data map[string]any) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}
