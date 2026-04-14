package datamanageres

import (
	"encoding/json"

	"github.com/odysseythink/mlog"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type ApiKeyAuthManager struct {
}

func (mgr ApiKeyAuthManager) GetAuthCredentials(tenant_id string, category string, provider string) map[string]any {
	data_source_api_key_bindings := new(models.DataSourceApiKeyAuthBinding)
	err := dbengine.Instance().DB.Model(&models.DataSourceApiKeyAuthBinding{}).Where("tenant_id = ?", tenant_id).Where("category = ?", category).Where("provider = ?", provider).Where("disabled = ?", false).First(data_source_api_key_bindings).Error
	if err != nil {
		mlog.Errorf("get DataSourceApiKeyAuthBinding failed:%v", err)
		data_source_api_key_bindings = nil
	}
	if data_source_api_key_bindings == nil {
		return nil
	}
	var credentials map[string]any
	err = json.Unmarshal([]byte(data_source_api_key_bindings.Credentials), &credentials)
	if err != nil {
		mlog.Errorf("json.Unmarshal failed:%v", err)
	}
	return credentials
}
