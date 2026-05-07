package datamanager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/mlog"
)

type ProviderCredentialsCacheType string

const (
	ProviderCredentialsCache_PROVIDER             ProviderCredentialsCacheType = "provider"
	ProviderCredentialsCache_MODEL                ProviderCredentialsCacheType = "provider_model"
	ProviderCredentialsCache_LOAD_BALANCING_MODEL ProviderCredentialsCacheType = "load_balancing_provider_model"
)

type ProviderCredentialsCache struct {
	cache_key string
}

func NewProviderCredentialsCache(tenant_id string, identity_id string, cache_type ProviderCredentialsCacheType) *ProviderCredentialsCache {
	return &ProviderCredentialsCache{
		cache_key: fmt.Sprintf("%v_credentials:tenant_id:%s:id:%s", cache_type, tenant_id, identity_id),
	}
}

func (c *ProviderCredentialsCache) Get() map[string]any {
	/*
	   Get cached model provider credentials.

	   :return:
	*/
	cached_provider_credentials := cache.Instance().GetString(c.cache_key)
	if cached_provider_credentials != "" {
		var credentials map[string]any
		err := json.Unmarshal([]byte(cached_provider_credentials), &credentials)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", cached_provider_credentials, err)
			return nil
		}
		return credentials
	}
	return nil
}
func (c *ProviderCredentialsCache) Set(credentials map[string]any) {
	/*
	   Cache model provider credentials.

	   :param credentials: provider credentials
	   :return:
	*/
	bindata, _ := json.Marshal(credentials)
	cache.Instance().SetEx(c.cache_key, string(bindata), 86400*time.Second)
}
func (c *ProviderCredentialsCache) Delete() {
	/*
	   Delete cached model provider credentials.

	   :return:
	*/
	cache.Instance().DelKey(c.cache_key)
}
