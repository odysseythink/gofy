package services

import (
	"mlib.com/confy"
	servicesenumtypes "mlib.com/gofy/server/enum_types/services"
	pbentities "mlib.com/gofy/server/proto/entities"
)

func NewSubscriptionModel(plan, interval string) *pbentities.SubscriptionModel {
	if plan == "" {
		plan = "sandbox"
	}
	return &pbentities.SubscriptionModel{
		Plan:     plan,
		Interval: interval,
	}
}

func NewBillingModel(enabled bool, subscription *pbentities.SubscriptionModel) *pbentities.BillingModel {
	if subscription == nil {
		subscription = NewSubscriptionModel("", "")
	}
	return &pbentities.BillingModel{
		Enabled:      enabled,
		Subscription: subscription,
	}
}

func NewLimitationModel(sz, limit int64) *pbentities.LimitationModel {
	return &pbentities.LimitationModel{
		Size:  sz,
		Limit: limit,
	}
}

func NewFeatureModel() *pbentities.FeatureModel {
	return &pbentities.FeatureModel{
		Billing:              NewBillingModel(false, nil),
		Members:              NewLimitationModel(0, 1),
		Apps:                 NewLimitationModel(0, 10),
		VectorSpace:          NewLimitationModel(0, 5),
		AnnotationQuotaLimit: NewLimitationModel(0, 10),
		DocumentsUploadQuota: NewLimitationModel(0, 50),
		DocsProcessing:       "standard",
	}
}

type SystemFeature struct {
	SsoEnforcedForSignin         bool          `json:"sso_enforced_for_signin"`
	SsoEnforcedForSigninProtocol string        `json:"sso_enforced_for_signin_protocol"`
	SsoEnforcedForWeb            bool          `json:"sso_enforced_for_web"`
	SsoEnforcedForWebProtocol    string        `json:"sso_enforced_for_web_protocol"`
	EnableWebSsoSwitchComponent  bool          `json:"enable_web_sso_switch_component"`
	EnableEmailCodeLogin         bool          `json:"enable_email_code_login"`
	EnableEmailPasswordLogin     bool          `json:"enable_email_password_login"`
	EnableSocialOauthLogin       bool          `json:"enable_social_oauth_login"`
	IsAllowRegister              bool          `json:"is_allow_register"`
	IsAllowCreateWorkspace       bool          `json:"is_allow_create_workspace"`
	IsEmailSetup                 bool          `json:"is_email_setup"`
	License                      *LicenseModel `json:"license"`
}

type LicenseModel struct {
	Status    servicesenumtypes.LicenseStatus `json:"status"`
	ExpiredAt string                          `json:"expired_at"`
}

func NewLicenseModel() *pbentities.LicenseModel {
	return &pbentities.LicenseModel{
		Status: string(servicesenumtypes.LicenseStatus_NONE),
	}
}

type FeatureService struct {
}

func (s *FeatureService) GetSystemFeatures() *pbentities.SystemFeature {
	system_features := &pbentities.SystemFeature{
		License: NewLicenseModel(),
	}

	// if confy.Get[bool]("ENTERPRISE_ENABLED") {
	// 	system_features.EnableWebSsoSwitchComponent = true
	// 	s.fulfillParamsFromEnterprise(system_features)
	// }
	// return system_features
	s.fulfillSystemParamsFromEnv(system_features)

	if confy.Get[bool]("ENTERPRISE_ENABLED") {
		system_features.EnableWebSsoSwitchComponent = true

		s.fulfillParamsFromEnterprise(system_features)
	}
	return system_features
}

func (s *FeatureService) fulfillSystemParamsFromEnv(features *pbentities.SystemFeature) {
	features.EnableEmailCodeLogin = confy.Get[bool]("enable_email_code_login")
	features.EnableEmailPasswordLogin = confy.Get[bool]("enable_email_password_login")
	features.EnableSocialOauthLogin = confy.Get[bool]("enable_social_oauth_login")
	features.IsAllowRegister = confy.Get[bool]("allow_register")
	features.IsAllowCreateWorkspace = confy.Get[bool]("allow_create_workspace")
	if confy.Get[string]("mail_type") != "" {
		features.IsEmailSetup = true
	}

}

func (s *FeatureService) fulfillParamsFromEnterprise(features *pbentities.SystemFeature) {
	enterprise_info := (&EnterpriseService{}).GetInfo()

	if enterprise_info != nil {
		if _, ok := enterprise_info["sso_enforced_for_signin"]; ok {
			if _, ok := enterprise_info["sso_enforced_for_signin"].(bool); ok {
				features.SsoEnforcedForSignin = enterprise_info["sso_enforced_for_signin"].(bool)
			}
		}

		if _, ok := enterprise_info["sso_enforced_for_signin_protocol"]; ok {
			if _, ok := enterprise_info["sso_enforced_for_signin_protocol"].(string); ok {
				features.SsoEnforcedForSigninProtocol = enterprise_info["sso_enforced_for_signin_protocol"].(string)
			}
		}

		if _, ok := enterprise_info["sso_enforced_for_web"]; ok {
			if _, ok := enterprise_info["sso_enforced_for_web"].(bool); ok {
				features.SsoEnforcedForWeb = enterprise_info["sso_enforced_for_web"].(bool)
			}
		}

		if _, ok := enterprise_info["sso_enforced_for_web_protocol"]; ok {
			if _, ok := enterprise_info["sso_enforced_for_web_protocol"].(string); ok {
				features.SsoEnforcedForWebProtocol = enterprise_info["sso_enforced_for_web_protocol"].(string)
			}
		}
	}
}

func (s *FeatureService) fulfillParamsFromEnv(f *pbentities.FeatureModel) {
	f.CanReplaceLogo = confy.Get[bool]("CAN_REPLACE_LOGO")
	f.ModelLoadBalancingEnabled = confy.Get[bool]("MODEL_LB_ENABLED")
	f.DatasetOperatorEnabled = confy.Get[bool]("dataset_operator_enabled")
}

func (s *FeatureService) fulfillParamsFromBillingApi(f *pbentities.FeatureModel, tenant_id string) {
	billing_info := (&BillingService{}).GetInfo(tenant_id)

	if _, ok := billing_info["enabled"]; ok {
		if _, ok := billing_info["enabled"].(bool); ok {
			f.Billing.Enabled = billing_info["enabled"].(bool)
		}
	}
	if _, ok := billing_info["subscription"]; ok {
		if val, ok := billing_info["subscription"].(*pbentities.SubscriptionModel); ok && val != nil {
			f.Billing.Subscription.Plan = val.Plan
			f.Billing.Subscription.Interval = val.Interval
		}
	}

	if _, ok := billing_info["members"]; ok {
		if val, ok := billing_info["members"].(*pbentities.LimitationModel); ok && val != nil {
			f.Members.Size = val.Size
			f.Members.Limit = val.Limit
		}
	}
	if _, ok := billing_info["apps"]; ok {
		if val, ok := billing_info["apps"].(*pbentities.LimitationModel); ok && val != nil {
			f.Apps.Size = val.Size
			f.Apps.Limit = val.Limit
		}
	}
	if _, ok := billing_info["vector_space"]; ok {
		if val, ok := billing_info["vector_space"].(*pbentities.LimitationModel); ok && val != nil {
			f.VectorSpace.Size = val.Size
			f.VectorSpace.Limit = val.Limit
		}
	}
	if _, ok := billing_info["documents_upload_quota"]; ok {
		if val, ok := billing_info["documents_upload_quota"].(*pbentities.LimitationModel); ok && val != nil {
			f.DocumentsUploadQuota.Size = val.Size
			f.DocumentsUploadQuota.Limit = val.Limit
		}
	}
	if _, ok := billing_info["annotation_quota_limit"]; ok {
		if val, ok := billing_info["annotation_quota_limit"].(*pbentities.LimitationModel); ok && val != nil {
			f.AnnotationQuotaLimit.Size = val.Size
			f.AnnotationQuotaLimit.Limit = val.Limit
		}
	}
	if _, ok := billing_info["docs_processing"]; ok {
		if val, ok := billing_info["docs_processing"].(string); ok {
			f.DocsProcessing = val
		}
	}

	if _, ok := billing_info["can_replace_logo"]; ok {
		if val, ok := billing_info["can_replace_logo"].(bool); ok {
			f.CanReplaceLogo = val
		}
	}

	if _, ok := billing_info["model_load_balancing_enabled"]; ok {
		if val, ok := billing_info["model_load_balancing_enabled"].(bool); ok {
			f.ModelLoadBalancingEnabled = val
		}
	}

}

func (s *FeatureService) GetFeatures(tenant_id string) *pbentities.FeatureModel {
	f := NewFeatureModel()

	s.fulfillParamsFromEnv(f)

	if confy.Get[bool]("billing_enabled") {
		s.fulfillParamsFromBillingApi(f, tenant_id)
	}
	// bindata, _ := json.Marshal(f)
	// feature_map := map[string]any{}
	// json.Unmarshal(bindata, &feature_map)
	// return system_features
	return f
}
