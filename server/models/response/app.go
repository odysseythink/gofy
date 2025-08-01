package response

import (
	"encoding/json"
	"time"

	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type SiteFields struct {
	AccessToken            string     `json:"access_token"` //": fields.String(attribute="code"),
	Code                   string     `json:"code"`
	Title                  string     `json:"title"`
	IconType               string     `json:"icon_type"`
	Icon                   string     `json:"icon"`
	IconBackground         string     `json:"icon_background"`
	IconURL                string     `json:"icon_url"`
	Description            string     `json:"description"`
	DefaultLanguage        string     `json:"default_language"`
	ChatColorTheme         string     `json:"chat_color_theme"`
	ChatColorThemeInverted bool       `json:"chat_color_theme_inverted"`
	CustomizeDomain        string     `json:"customize_domain"`
	Copyright              string     `json:"copyright"`
	PrivacyPolicy          string     `json:"privacy_policy"`
	CustomDisclaimer       string     `json:"custom_disclaimer"`
	CustomizeTokenStrategy string     `json:"customize_token_strategy"`
	PromptPublic           bool       `json:"prompt_public"`
	AppBaseURL             string     `json:"app_base_url"`
	ShowWorkflowSteps      bool       `json:"show_workflow_steps"`
	UseIconAsAnswerIcon    bool       `json:"use_icon_as_answer_icon"`
	CreatedBy              string     `json:"created_by"`
	CreatedAt              *time.Time `json:"created_at"`
	UpdatedBy              string     `json:"updated_by"`
	UpdatedAt              *time.Time `json:"updated_at"`
}

type ModelConfigFields struct {
	OpeningStatement              string           `json:"opening_statement"`
	SuggestedQuestions            []any            `json:"suggested_questions"`              //": fields.Raw(attribute="suggested_questions_list"),
	SuggestedQuestionsAfterAnswer map[string]any   `json:"suggested_questions_after_answer"` //": fields.Raw(attribute="suggested_questions_after_answer_dict"),
	SpeechToText                  map[string]any   `json:"speech_to_text"`                   //": fields.Raw(attribute="speech_to_text_dict"),
	TextToSpeech                  map[string]any   `json:"text_to_speech"`                   //": fields.Raw(attribute="text_to_speech_dict"),
	RetrieverResource             map[string]any   `json:"retriever_resource"`               //": fields.Raw(attribute="retriever_resource_dict"),
	AnnotationReply               map[string]any   `json:"annotation_reply"`                 //": fields.Raw(attribute="annotation_reply_dict"),
	MoreLikeThis                  map[string]any   `json:"more_like_this"`                   //": fields.Raw(attribute="more_like_this_dict"),
	SensitiveWordAvoidance        map[string]any   `json:"sensitive_word_avoidance"`         //": fields.Raw(attribute="sensitive_word_avoidance_dict"),
	ExternalDataTools             []map[string]any `json:"external_data_tools"`              //": fields.Raw(attribute="external_data_tools_list"),
	Model                         map[string]any   `json:"model"`                            //": fields.Raw(attribute="model_dict"),
	UserInputForm                 []map[string]any `json:"user_input_form"`                  //": fields.Raw(attribute="user_input_form_list"),
	DatasetQueryVariable          string           `json:"dataset_query_variable"`
	PrePrompt                     string           `json:"pre_prompt"`
	AgentMode                     map[string]any   `json:"agent_mode"` //": fields.Raw(attribute="agent_mode_dict"),
	PromptType                    string           `json:"prompt_type"`
	ChatPromptConfig              map[string]any   `json:"chat_prompt_config"`       //": fields.Raw(attribute="chat_prompt_config_dict"),
	CompletionPromptConfig        map[string]any   `json:"completion_prompt_config"` //": fields.Raw(attribute="completion_prompt_config_dict"),
	DatasetConfigs                map[string]any   `json:"dataset_configs"`          //": fields.Raw(attribute="dataset_configs_dict"),
	FileUpload                    map[string]any   `json:"file_upload"`              //": fields.Raw(attribute="file_upload_dict"),
	CreatedBy                     string           `json:"created_by"`
	CreatedAt                     int64            `json:"created_at"`
	UpdatedBy                     string           `json:"updated_by"`
	UpdatedAt                     int64            `json:"updated_at"`
}

type AppDetailWithSiteResponse struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Mode                models.AppMode         `json:"mode"` //": fields.String(attribute="mode_compatible_with_agent"),
	IconType            string                 `json:"icon_type"`
	Icon                string                 `json:"icon"`
	Icon_background     string                 `json:"icon_background"`
	IconURL             string                 `json:"icon_url"` //: AppIconUrlField,
	EnableSite          bool                   `json:"enable_site"`
	EnableApi           bool                   `json:"enable_api"`
	ModelConfig         *ModelConfigFields     `json:"model_config"` //": fields.Nested(model_config_fields, attribute="app_model_config", allow_null=True),
	Workflow            *WorkflowPartialFields `json:"workflow"`     //": fields.Nested(workflow_partial_fields, allow_null=True),
	Site                *SiteFields            `json:"site"`         //": fields.Nested(site_fields),
	ApiBaseURL          string                 `json:"api_base_url"`
	UseIconAsAnswerIcon bool                   `json:"use_icon_as_answer_icon"`
	CreatedBy           string                 `json:"created_by"`
	CreatedAt           int64                  `json:"created_at"`
	UpdatedBy           string                 `json:"updated_by"`
	UpdatedAt           int64                  `json:"updated_at"`
	DeletedTools        []string               `json:"deleted_tools"` //": fields.List(fields.String),
}

type AppDetailResponse struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Mode                models.AppMode         `json:"mode"` //": fields.String(attribute="mode_compatible_with_agent"),
	Icon                string                 `json:"icon"`
	Icon_background     string                 `json:"icon_background"`
	EnableSite          bool                   `json:"enable_site"`
	EnableApi           bool                   `json:"enable_api"`
	ModelConfig         *ModelConfigFields     `json:"model_config"` //": fields.Nested(model_config_fields, attribute="app_model_config", allow_null=True),
	Workflow            *WorkflowPartialFields `json:"workflow"`     //": fields.Nested(workflow_partial_fields, allow_null=True),
	Tracing             string                 `json:"tracing"`
	UseIconAsAnswerIcon bool                   `json:"use_icon_as_answer_icon"`
	CreatedBy           string                 `json:"created_by"`
	CreatedAt           int64                  `json:"created_at"`
	UpdatedBy           string                 `json:"updated_by"`
	UpdatedAt           int64                  `json:"updated_at"`
}

func NewAppDetailResponse(args any) *AppDetailResponse {

	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(AppDetailResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AppDetailResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(AppDetailResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AppDetailResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if app, ok := args.(*models.App); ok && app != nil {
		rsp := &AppDetailResponse{
			ID:                  app.ID,
			Name:                app.Name,
			Description:         app.Description,
			Mode:                app.Mode,
			Icon:                app.Icon,
			Icon_background:     app.IconBackground,
			EnableSite:          app.EnableSite,
			EnableApi:           app.EnableAPI,
			Tracing:             app.Tracing,
			UseIconAsAnswerIcon: app.UseIconAsAnswerIcon,
			CreatedBy:           app.CreatedBy,
			CreatedAt:           app.CreatedAt.Unix(),
			UpdatedBy:           app.UpdatedBy,
			UpdatedAt:           app.UpdatedAt.Unix(),
		}
		app_model_config := app.AppModelConfig()
		if app_model_config != nil {
			rsp.ModelConfig = &ModelConfigFields{
				OpeningStatement:              app_model_config.OpeningStatement,
				SuggestedQuestions:            app_model_config.SuggestedQuestionsList(),
				SuggestedQuestionsAfterAnswer: app_model_config.SuggestedQuestionsAfterAnswerDict(),
				SpeechToText:                  app_model_config.SpeechToTextDict(),
				TextToSpeech:                  app_model_config.TextToSpeechDict(),
				RetrieverResource:             app_model_config.RetrieverResourceDict(),
				AnnotationReply:               app_model_config.AnnotationReplyDict(),
				MoreLikeThis:                  app_model_config.MoreLikeThisDict(),
				SensitiveWordAvoidance:        app_model_config.SensitiveWordAvoidanceDict(),
				ExternalDataTools:             app_model_config.ExternalDataToolsList(),
				Model:                         app_model_config.ModelDict(),
				UserInputForm:                 app_model_config.UserInputFormList(),
				DatasetQueryVariable:          app_model_config.DatasetQueryVariable,
				PrePrompt:                     app_model_config.PrePrompt,
				AgentMode:                     app_model_config.AgentModeDict(),
				PromptType:                    app_model_config.PromptType,
				ChatPromptConfig:              app_model_config.ChatPromptConfigDict(),
				CompletionPromptConfig:        app_model_config.CompletionPromptConfigDict(),
				DatasetConfigs:                app_model_config.DatasetConfigsDict(),
				FileUpload:                    app_model_config.FileUploadDict(),
				CreatedBy:                     app_model_config.CreatedBy,
				CreatedAt:                     app_model_config.CreatedAt.Unix(),
				UpdatedBy:                     app_model_config.UpdatedBy,
				UpdatedAt:                     app_model_config.UpdatedAt.Unix(),
			}
		}
		wf := app.Workflow()
		if wf != nil {
			rsp.Workflow = &WorkflowPartialFields{
				ID:        wf.ID,
				CreatedBy: wf.CreatedBy,
				CreatedAt: wf.CreatedAt.Unix(),
				UpdatedBy: wf.UpdatedBy,
				UpdatedAt: wf.UpdatedAt.Unix(),
			}
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(AppDetailResponse)
}

type TagFields struct {
	ID   string `json:"id  "`
	Name string `json:"name"`
	Type string `json:"type"`
}

type AppPartialResponse struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	MaxActiveRequests   int                    `json:"max_active_requests"`
	Description         string                 `json:"description"`
	Mode                models.AppMode         `json:"mode"` //": fields.String(attribute="mode_compatible_with_agent"),
	IconType            string                 `json:"icon_type"`
	Icon                string                 `json:"icon"`
	IconBackground      string                 `json:"icon_background"`
	IconURL             string                 `json:"icon_url"`
	ModelConfig         *ModelConfigFields     `json:"model_config"` //": fields.Nested(model_config_fields, attribute="app_model_config", allow_null=True),
	Workflow            *WorkflowPartialFields `json:"workflow"`     //": fields.Nested(workflow_partial_fields, allow_null=True),
	Tracing             string                 `json:"tracing"`
	UseIconAsAnswerIcon bool                   `json:"use_icon_as_answer_icon"`
	CreatedBy           string                 `json:"created_by"`
	CreatedAt           int64                  `json:"created_at"`
	UpdatedBy           string                 `json:"updated_by"`
	UpdatedAt           int64                  `json:"updated_at"`
	Tags                []*TagFields           `json:"tags"`
}
type AppPaginationResponse struct {
	Page    int32                 `json:"page"`
	Limit   int32                 `json:"limit"` //"per_page"),
	Total   int64                 `json:"total"`
	HasMore bool                  `json:"has_more"` //"has_next"),
	Data    []*AppPartialResponse `json:"data"`     //"items"),
}

func NewAppPaginationResponse(args any) *AppPaginationResponse {
	rsp := new(AppPaginationResponse)
	if real_args, ok := args.(string); ok && real_args != "" {
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AppPaginationResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AppPaginationResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return rsp
}

func NewAppDetailWithSiteResponse(args any) *AppDetailWithSiteResponse {

	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(AppDetailWithSiteResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AppDetailWithSiteResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(AppDetailWithSiteResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AppDetailWithSiteResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if app_model, ok := args.(*models.App); ok && app_model != nil {
		rsp := &AppDetailWithSiteResponse{
			ID:                  app_model.ID,
			Name:                app_model.Name,
			Description:         app_model.Description,
			Mode:                app_model.Mode,
			IconType:            app_model.IconType,
			Icon:                app_model.Icon,
			Icon_background:     app_model.IconBackground,
			IconURL:             app_model.IconUrl,
			EnableSite:          app_model.EnableSite,
			EnableApi:           app_model.EnableAPI,
			ApiBaseURL:          app_model.ApiBaseUrl(),
			UseIconAsAnswerIcon: app_model.UseIconAsAnswerIcon,
			CreatedBy:           app_model.CreatedBy,
			UpdatedBy:           app_model.UpdatedBy,
			DeletedTools:        app_model.DeletedTools(),
		}
		model_config := app_model.AppModelConfig()
		if model_config != nil {
			rsp.ModelConfig = &ModelConfigFields{
				OpeningStatement:              model_config.OpeningStatement,
				SuggestedQuestions:            model_config.SuggestedQuestionsList(),
				SuggestedQuestionsAfterAnswer: model_config.SuggestedQuestionsAfterAnswerDict(),
				SpeechToText:                  model_config.SpeechToTextDict(),
				TextToSpeech:                  model_config.TextToSpeechDict(),
				RetrieverResource:             model_config.RetrieverResourceDict(),
				AnnotationReply:               model_config.AnnotationReplyDict(),
				MoreLikeThis:                  model_config.MoreLikeThisDict(),
				SensitiveWordAvoidance:        model_config.SensitiveWordAvoidanceDict(),
				ExternalDataTools:             model_config.ExternalDataToolsList(),
				Model:                         model_config.ModelDict(),
				UserInputForm:                 model_config.UserInputFormList(),
				DatasetQueryVariable:          model_config.DatasetQueryVariable,
				PrePrompt:                     model_config.PrePrompt,
				AgentMode:                     model_config.AgentModeDict(),
				PromptType:                    model_config.PromptType,
				ChatPromptConfig:              model_config.ChatPromptConfigDict(),
				CompletionPromptConfig:        model_config.CompletionPromptConfigDict(),
				DatasetConfigs:                model_config.DatasetConfigsDict(),
				FileUpload:                    model_config.FileUploadDict(),
				CreatedBy:                     model_config.CreatedBy,
				UpdatedBy:                     model_config.UpdatedBy,
			}
			if model_config.CreatedAt != nil {
				rsp.ModelConfig.CreatedAt = model_config.CreatedAt.Unix()
			}
			if model_config.UpdatedAt != nil {
				rsp.ModelConfig.UpdatedAt = model_config.UpdatedAt.Unix()
			}
		}
		wf := app_model.Workflow()
		if wf != nil {
			rsp.Workflow = &WorkflowPartialFields{
				ID:        wf.ID,
				CreatedBy: wf.CreatedBy,
				UpdatedBy: wf.UpdatedBy,
			}
			if wf.CreatedAt != nil {
				rsp.Workflow.CreatedAt = wf.CreatedAt.Unix()
			}
			if wf.UpdatedAt != nil {
				rsp.Workflow.UpdatedAt = wf.UpdatedAt.Unix()
			}
		}
		site := app_model.Site()
		if site != nil {
			rsp.Site = &SiteFields{
				AccessToken:            site.Code,
				Code:                   site.Code,
				Title:                  site.Title,
				IconType:               site.IconType,
				Icon:                   site.Icon,
				IconBackground:         site.IconBackground,
				IconURL:                site.Icon,
				Description:            site.Description,
				DefaultLanguage:        site.DefaultLanguage,
				ChatColorTheme:         site.ChatColorTheme,
				ChatColorThemeInverted: site.ChatColorThemeInverted,
				CustomizeDomain:        site.CustomizeDomain,
				Copyright:              site.Copyright,
				PrivacyPolicy:          site.PrivacyPolicy,
				CustomDisclaimer:       site.CustomDisclaimer,
				CustomizeTokenStrategy: site.CustomizeTokenStrategy,
				PromptPublic:           site.PromptPublic,
				AppBaseURL:             site.AppBaseURL(),
				ShowWorkflowSteps:      site.ShowWorkflowSteps,
				UseIconAsAnswerIcon:    site.UseIconAsAnswerIcon,
				CreatedBy:              site.CreatedBy,
				CreatedAt:              site.CreatedAt,
				UpdatedBy:              site.UpdatedBy,
				UpdatedAt:              site.UpdatedAt,
			}
		}
		if app_model.CreatedAt != nil {
			rsp.CreatedAt = app_model.CreatedAt.Unix()
		}
		if app_model.UpdatedAt != nil {
			rsp.UpdatedAt = app_model.UpdatedAt.Unix()
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(AppDetailWithSiteResponse)

}

type AppImportResponse struct {
	ID                 string `json:"id"`
	Status             string `json:"status"`
	AppID              string `json:"app_id"`
	CurrentDSLVersion  string `json:"current_dsl_version"`
	ImportedDSLVersion string `json:"imported_dsl_version"`
	Error              string `json:"error"`
}
