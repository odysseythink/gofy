package hostingconfiguration

import (
	"strings"
	"sync"

	"github.com/odysseythink/confy"
	providerentities "github.com/odysseythink/gofy/backend/entities/provider"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	providerenumtypes "github.com/odysseythink/gofy/backend/enum_types/provider"
)

type HostingQuotaer interface {
	Type() providerenumtypes.ProviderQuotaType
	GetRestrictModels() []*providerentities.RestrictModel
	SetRestrictModels([]*providerentities.RestrictModel)
}

// HostingQuota 定义 Hosting 配额
type HostingQuota struct {
	RestrictModels []*providerentities.RestrictModel
}

func (q *HostingQuota) GetRestrictModels() []*providerentities.RestrictModel {
	return q.RestrictModels
}
func (q *HostingQuota) SetRestrictModels(val []*providerentities.RestrictModel) {
	q.RestrictModels = val
}

// TrialHostingQuota 定义试用 Hosting 配额
type TrialHostingQuota struct {
	*HostingQuota
	QuotaLimit int
}

func (q *TrialHostingQuota) Type() providerenumtypes.ProviderQuotaType {
	return providerenumtypes.ProviderQuota_TRIAL
}

// PaidHostingQuota 定义付费 Hosting 配额
type PaidHostingQuota struct {
	*HostingQuota
}

func (q *PaidHostingQuota) Type() providerenumtypes.ProviderQuotaType {
	return providerenumtypes.ProviderQuota_PAID
}

// FreeHostingQuota 定义免费 Hosting 配额
type FreeHostingQuota struct {
	*HostingQuota
}

func (q *FreeHostingQuota) Type() providerenumtypes.ProviderQuotaType {
	return providerenumtypes.ProviderQuota_FREE
}

// HostingProvider 定义 Hosting 提供商
type HostingProvider struct {
	Enabled     bool
	Credentials map[string]any                  // 使用 any来表示任意类型
	QuotaUnit   providerenumtypes.QuotaUnitType // 使用指针来表示可选字段
	Quotas      []HostingQuotaer
}

// HostedModerationConfig 定义托管内容审核配置
type HostedModerationConfig struct {
	Enabled   bool
	Providers []string
}

// HostingConfiguration 定义 Hosting 配置
type HostingConfiguration struct {
	ProviderMap      map[string]*HostingProvider
	moderationConfig *HostedModerationConfig
}

// func NewHostingConfiguration(gofyConfig GofyConfig) *HostingConfiguration {
// 	return &HostingConfiguration{
// 		providerMap: make(map[string]*HostingProvider),
// 		gofyConfig:  gofyConfig,
// 	}
// }

var (
	mInstance *HostingConfiguration
	mOnce     sync.Once
)

func Instance() *HostingConfiguration {
	mOnce.Do(func() {
		mInstance = &HostingConfiguration{
			ProviderMap: make(map[string]*HostingProvider),
		}
		mInstance.InitApp()
	})
	return mInstance
}

func (hc *HostingConfiguration) InitApp() {
	hc.ProviderMap["azure_openai"] = hc.initAzureOpenAI()
	hc.ProviderMap["openai"] = hc.initOpenAI()
	hc.ProviderMap["anthropic"] = hc.initAnthropic()
	hc.ProviderMap["minimax"] = hc.initMinimax()
	hc.ProviderMap["spark"] = hc.initSpark()
	hc.ProviderMap["zhipuai"] = hc.initZhipuai()

	hc.moderationConfig = hc.initModerationConfig()
}

func (hc *HostingConfiguration) initAzureOpenAI() *HostingProvider {
	quotaUnit := providerenumtypes.QuotaUnit_TIMES
	if confy.GetWithDefault[bool]("hosted_azure_openai_config.enable", false) {
		credentials := map[string]any{
			"openai_api_key":  confy.GetWithDefault[string]("hosted_azure_openai_config.api_key", ""),
			"openai_api_base": confy.GetWithDefault[string]("hosted_azure_openai_config.api_base", ""),
			"base_model_name": "gpt-35-turbo",
		}

		hostedQuotaLimit := confy.GetWithDefault[int]("hosted_azure_openai_config.quota_limit", 200)
		trialQuota := &TrialHostingQuota{
			HostingQuota: &HostingQuota{
				RestrictModels: []*providerentities.RestrictModel{
					&providerentities.RestrictModel{Model: "gpt-4", BaseModelName: "gpt-4", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-4o", BaseModelName: "gpt-4o", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-4o-mini", BaseModelName: "gpt-4o-mini", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-4-32k", BaseModelName: "gpt-4-32k", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-4-1106-preview", BaseModelName: "gpt-4-1106-preview", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-4-vision-preview", BaseModelName: "gpt-4-vision-preview", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-35-turbo", BaseModelName: "gpt-35-turbo", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-35-turbo-1106", BaseModelName: "gpt-35-turbo-1106", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-35-turbo-instruct", BaseModelName: "gpt-35-turbo-instruct", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "gpt-35-turbo-16k", BaseModelName: "gpt-35-turbo-16k", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "text-davinci-003", BaseModelName: "text-davinci-003", ModelType: modelruntimeenumtypes.Model_LLM},
					&providerentities.RestrictModel{Model: "text-embedding-ada-002", BaseModelName: "text-embedding-ada-002", ModelType: modelruntimeenumtypes.Model_TEXT_EMBEDDING},
					&providerentities.RestrictModel{Model: "text-embedding-3-small", BaseModelName: "text-embedding-3-small", ModelType: modelruntimeenumtypes.Model_TEXT_EMBEDDING},
					&providerentities.RestrictModel{Model: "text-embedding-3-large", BaseModelName: "text-embedding-3-large", ModelType: modelruntimeenumtypes.Model_TEXT_EMBEDDING},
				},
			},
			QuotaLimit: hostedQuotaLimit,
		}

		return &HostingProvider{
			Enabled:     true,
			Credentials: credentials,
			QuotaUnit:   quotaUnit,
			Quotas:      []HostingQuotaer{trialQuota},
		}
	}

	return &HostingProvider{
		Enabled:   false,
		QuotaUnit: quotaUnit,
	}
}

func (hc *HostingConfiguration) initOpenAI() *HostingProvider {
	quotaUnit := providerenumtypes.QuotaUnit_CREDITS
	quotas := make([]HostingQuotaer, 0)

	if confy.GetWithDefault[bool]("hosted_openai_config.trial_enable", false) {
		hostedQuotaLimit := confy.GetWithDefault[int]("hosted_openai_config.quota_limit", 200)
		trialModels := hc.parseRestrictModelsFromEnv("HOSTED_OPENAI_TRIAL_MODELS")
		trialQuota := &TrialHostingQuota{
			HostingQuota: &HostingQuota{
				RestrictModels: trialModels,
			},
			QuotaLimit: hostedQuotaLimit,
		}
		quotas = append(quotas, trialQuota)
	}

	if confy.GetWithDefault[bool]("hosted_openai_config.paid_enable", false) {
		paidModels := hc.parseRestrictModelsFromEnv("HOSTED_OPENAI_PAID_MODELS")
		paidQuota := &PaidHostingQuota{
			HostingQuota: &HostingQuota{
				RestrictModels: paidModels,
			},
		}
		quotas = append(quotas, paidQuota)
	}

	if len(quotas) > 0 {
		credentials := map[string]any{
			"openai_api_key":      confy.GetWithDefault[string]("hosted_openai_config.api_key", ""),
			"openai_api_base":     confy.GetWithDefault[string]("hosted_openai_config.api_base", ""),
			"openai_organization": confy.GetWithDefault[string]("hosted_openai_config.api_organization", ""),
		}

		return &HostingProvider{
			Enabled:     true,
			Credentials: credentials,
			QuotaUnit:   quotaUnit,
			Quotas:      quotas,
		}
	}

	return &HostingProvider{
		Enabled:   false,
		QuotaUnit: quotaUnit,
	}
}

func (hc *HostingConfiguration) initAnthropic() *HostingProvider {
	quotaUnit := providerenumtypes.QuotaUnit_TOKENS
	quotas := make([]HostingQuotaer, 0)

	if confy.GetWithDefault[bool]("hosted_anthropic_config.trial_enable", false) {
		hostedQuotaLimit := confy.GetWithDefault[int]("hosted_anthropic_config.quota_limit", 200)
		trialQuota := &TrialHostingQuota{
			HostingQuota: &HostingQuota{},
			QuotaLimit:   hostedQuotaLimit,
		}
		quotas = append(quotas, trialQuota)
	}

	if confy.GetWithDefault[bool]("hosted_anthropic_config.pid_enable", false) {
		paidQuota := &PaidHostingQuota{
			HostingQuota: &HostingQuota{},
		}
		quotas = append(quotas, paidQuota)
	}

	if len(quotas) > 0 {
		credentials := map[string]any{
			"anthropic_api_key":  confy.GetWithDefault[string]("hosted_anthropic_config.api_key", ""),
			"anthropic_api_base": confy.GetWithDefault[string]("hosted_anthropic_config.api_base", ""),
		}

		return &HostingProvider{
			Enabled:     true,
			Credentials: credentials,
			QuotaUnit:   quotaUnit,
			Quotas:      quotas,
		}
	}

	return &HostingProvider{
		Enabled:   false,
		QuotaUnit: quotaUnit,
	}
}

func (hc *HostingConfiguration) initMinimax() *HostingProvider {
	quotaUnit := providerenumtypes.QuotaUnit_TOKENS
	if confy.GetWithDefault[bool]("hosted_minmax_config.enable", false) {
		quotas := []HostingQuotaer{&FreeHostingQuota{HostingQuota: &HostingQuota{}}}

		return &HostingProvider{
			Enabled:     true,
			Credentials: nil,
			QuotaUnit:   quotaUnit,
			Quotas:      quotas,
		}
	}

	return &HostingProvider{
		Enabled:   false,
		QuotaUnit: quotaUnit,
	}
}

func (hc *HostingConfiguration) initSpark() *HostingProvider {
	quotaUnit := providerenumtypes.QuotaUnit_TOKENS
	if confy.GetWithDefault[bool]("hosted_spark_config.enable", false) {
		quotas := []HostingQuotaer{&FreeHostingQuota{HostingQuota: &HostingQuota{}}}

		return &HostingProvider{
			Enabled:     true,
			Credentials: nil,
			QuotaUnit:   quotaUnit,
			Quotas:      quotas,
		}
	}

	return &HostingProvider{
		Enabled:   false,
		QuotaUnit: quotaUnit,
	}
}

func (hc *HostingConfiguration) initZhipuai() *HostingProvider {
	quotaUnit := providerenumtypes.QuotaUnit_TOKENS
	if confy.GetWithDefault[bool]("hosted_zhipu_ai_config.enable", false) {
		quotas := []HostingQuotaer{&FreeHostingQuota{HostingQuota: &HostingQuota{}}}

		return &HostingProvider{
			Enabled:     true,
			Credentials: nil,
			QuotaUnit:   quotaUnit,
			Quotas:      quotas,
		}
	}

	return &HostingProvider{
		Enabled:   false,
		QuotaUnit: quotaUnit,
	}
}

func (hc *HostingConfiguration) initModerationConfig() *HostedModerationConfig {
	if confy.GetWithDefault[bool]("hosted_moderation_config.enable", false) && confy.GetWithDefault[string]("hosted_moderation_config.providers", "") != "" {
		providers := strings.Split(confy.GetWithDefault[string]("hosted_moderation_config.providers", ""), ",")
		return &HostedModerationConfig{
			Enabled:   true,
			Providers: providers,
		}
	}

	return &HostedModerationConfig{
		Enabled: false,
	}
}

func (hc *HostingConfiguration) parseRestrictModelsFromEnv(envVar string) []*providerentities.RestrictModel {
	return nil
}
