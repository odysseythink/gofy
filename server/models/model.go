package models

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

var (
	DEFAULT_APP_TEMPLATES = map[AppMode]map[string]map[string]any{
		// workflow default mode
		AppMode_WORKFLOW: {
			"app": {
				"mode":        AppMode_WORKFLOW,
				"enable_site": true,
				"enable_api":  true,
			},
		},
		// completion default mode
		AppMode_COMPLETION: {
			"app": {
				"mode":        AppMode_COMPLETION,
				"enable_site": true,
				"enable_api":  true,
			},
			"model_config": {
				"model": map[string]any{
					"provider":          "openai",
					"name":              "gpt-4o",
					"mode":              "chat",
					"completion_params": map[string]any{},
				},
				"user_input_form": `[{"paragraph": {"label": "Query","variable": "query","required": true,"default": ""}}]`,
				"pre_prompt":      "{{query}}",
			},
		},
		// chat default mode
		AppMode_CHAT: {
			"app": {
				"mode":        AppMode_CHAT,
				"enable_site": true,
				"enable_api":  true,
			},
			"model_config": {
				"model": map[string]any{
					"provider":          "openai",
					"name":              "gpt-4o",
					"mode":              "chat",
					"completion_params": map[string]any{},
				},
			},
		},
		// advanced-chat default mode
		AppMode_ADVANCED_CHAT: {
			"app": {
				"mode":        AppMode_ADVANCED_CHAT,
				"enable_site": true,
				"enable_api":  true,
			},
		},
		// agent-chat default mode
		AppMode_AGENT_CHAT: {
			"app": {
				"mode":        AppMode_AGENT_CHAT,
				"enable_site": true,
				"enable_api":  true,
			},
			"model_config": {
				"model": map[string]any{
					"provider":          "openai",
					"name":              "gpt-4o",
					"mode":              "chat",
					"completion_params": map[string]any{},
				},
			},
		},
	}
)

// 定义AppMode类型
type AppMode string

type IconType string

const (
	AppMode_COMPLETION    AppMode = "completion"
	AppMode_WORKFLOW      AppMode = "workflow"
	AppMode_CHAT          AppMode = "chat"
	AppMode_ADVANCED_CHAT AppMode = "advanced-chat"
	AppMode_AGENT_CHAT    AppMode = "agent-chat"
	AppMode_CHANNEL       AppMode = "channel"

	Icon_IMAGE IconType = "image"
	Icon_EMOJI IconType = "emoji"
)

// App [...]
type App struct {
	ID                  string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID            string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Name                string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Mode                AppMode    `gorm:"column:mode;type:varchar(255);not null" json:"mode"`
	Icon                string     `gorm:"column:icon;type:varchar(255)" json:"icon"`
	IconBackground      string     `gorm:"column:icon_background;type:varchar(255)" json:"icon_background"`
	AppModelConfigID    string     `gorm:"column:app_model_config_id;type:varchar(36)" json:"app_model_config_id"`
	Status              string     `gorm:"column:status;type:varchar(255);default:normal" json:"status"`
	EnableSite          bool       `gorm:"column:enable_site;type:tinyint(1);not null" json:"enable_site"`
	EnableAPI           bool       `gorm:"column:enable_api;type:tinyint(1);not null" json:"enable_api"`
	APIRpm              int        `gorm:"column:api_rpm;type:int;not null;default:0" json:"api_rpm"`
	APIRph              int        `gorm:"column:api_rph;type:int;not null;default:0" json:"api_rph"`
	IsDemo              bool       `gorm:"column:is_demo;type:tinyint(1);not null;default:0" json:"is_demo"`
	IsPublic            bool       `gorm:"column:is_public;type:tinyint(1);not null;default:0" json:"is_public"`
	CreatedAt           *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsUniversal         bool       `gorm:"column:is_universal;type:tinyint(1);not null;default:0" json:"is_universal"`
	WorkflowID          string     `gorm:"column:workflow_id;type:varchar(36)" json:"workflow_id"`
	Description         string     `gorm:"column:description;type:varchar(255);default:''" json:"description"`
	Tracing             string     `gorm:"column:tracing;type:text" json:"tracing"`
	MaxActiveRequests   int        `gorm:"column:max_active_requests;type:int" json:"max_active_requests"`
	IconType            string     `gorm:"column:icon_type;type:varchar(255)" json:"icon_type"`
	IconUrl             string     `json:"icon_url" form:"-" gorm:"-:all"`
	CreatedBy           string     `gorm:"column:created_by;type:varchar(36)" json:"created_by"`
	UpdatedBy           string     `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UseIconAsAnswerIcon bool       `gorm:"column:use_icon_as_answer_icon;type:tinyint(1);not null;default:0" json:"use_icon_as_answer_icon"`
}

// TableName get sql table name.获取数据库表名
func (App) TableName() string {
	return "apps"
}

func (a *App) DescOrPrompt() string {
	if a.Description != "" {
		return a.Description
	} else {
		app_model_config := a.AppModelConfig()
		if app_model_config != nil {
			return app_model_config.PrePrompt
		} else {
			return ""
		}
	}
}

func (a *App) AppModelConfig() *AppModelConfig {
	if a.AppModelConfigID != "" {
		app_model_config := new(AppModelConfig)
		err := dbengine.Instance().DB.Model(&AppModelConfig{}).Where("id = ?", a.AppModelConfigID).First(app_model_config).Error
		if err != nil {
			mlog.Errorf("get AppModelConfig failed:%v", err)
			return nil
		}
		return app_model_config
	}
	return nil
}
func (a *App) IsAgent() bool {
	app_model_config := a.AppModelConfig()
	if app_model_config == nil {
		return false
	}
	if app_model_config.AgentMode == "" {
		return false
	}
	agent_mode_dict := app_model_config.AgentModeDict()
	enabled := false
	if _, ok := agent_mode_dict["enabled"]; ok {
		if _, ok := agent_mode_dict["enabled"].(bool); ok {
			enabled = agent_mode_dict["enabled"].(bool)
		}
	}
	strategy := ""
	if _, ok := agent_mode_dict["strategy"]; ok {
		if _, ok := agent_mode_dict["strategy"].(string); ok {
			strategy = agent_mode_dict["strategy"].(string)
		}
	}
	if enabled && slices.Contains([]string{"function_call", "react"}, strategy) {
		a.Mode = AppMode_AGENT_CHAT
		dbengine.Instance().DB.Updates(&App{ID: a.ID, Mode: a.Mode})
		return true
	}
	return false
}

func (a *App) ModeCompatibleWithAgent() string {
	if a.Mode == AppMode_CHAT && a.IsAgent() {
		return string(AppMode_AGENT_CHAT)
	}
	return string(a.Mode)
}

func (a *App) DeletedTools() []string {
	// get agent mode tools
	app_model_config := a.AppModelConfig()
	if app_model_config == nil {
		return nil
	}
	if app_model_config.AgentMode == "" {
		return nil
	}
	agent_mode_dict := app_model_config.AgentModeDict()
	tools := make([]map[string]any, 0)
	if _, ok := agent_mode_dict["tools"]; ok {
		if _, ok := agent_mode_dict["tools"].([]map[string]any); ok {
			tools = agent_mode_dict["tools"].([]map[string]any)
		}
	}

	provider_ids := make([]string, 0)

	for _, tool := range tools {
		if len(tool) >= 4 {
			provider_type := ""
			if _, ok := tool["provider_type"]; ok {
				if _, ok := tool["provider_type"].(string); ok {
					provider_type = tool["provider_type"].(string)
				}
			}
			provider_id := ""
			if _, ok := tool["provider_id"]; ok {
				if _, ok := tool["provider_id"].(string); ok {
					provider_id = tool["provider_id"].(string)
				}
			}
			if provider_type == "api" {
				// // check if provider id is a uuid string, if not, skip
				if uuid.FromStringOrNil(provider_type) != uuid.Nil {
					provider_ids = append(provider_ids, provider_id)
				}
			}
		}
	}
	if len(provider_ids) < 1 {
		return nil
	}
	api_providers := make([]string, 0)
	err := dbengine.Instance().DB.Raw("SELECT id FROM tool_api_providers WHERE id IN ?", provider_ids).Scan(&api_providers).Error
	if err != nil {
		return nil
	}

	deleted_tools := make([]string, 0)
	current_api_provider_ids := api_providers

	for _, tool := range tools {
		if len(tool) >= 4 {
			provider_type := ""
			if _, ok := tool["provider_type"]; ok {
				if _, ok := tool["provider_type"].(string); ok {
					provider_type = tool["provider_type"].(string)
				}
			}
			provider_id := ""
			if _, ok := tool["provider_id"]; ok {
				if _, ok := tool["provider_id"].(string); ok {
					provider_id = tool["provider_id"].(string)
				}
			}
			tool_name := ""
			if _, ok := tool["tool_name"]; ok {
				if _, ok := tool["tool_name"].(string); ok {
					tool_name = tool["tool_name"].(string)
				}
			}
			exist := false
			for _, v := range current_api_provider_ids {
				if provider_id == v {
					exist = true
					break
				}
			}
			if provider_type == "api" && !exist && tool_name != "" {
				deleted_tools = append(deleted_tools, tool_name)
			}
		}
	}
	return deleted_tools
}

func (a *App) Tags() []*Tag {
	ret := make([]*Tag, 0)
	err := dbengine.Instance().DB.Raw("SELECT t.* FROM tags t join tag_bindings tb on t.id = tb.tag_id WHERE tb.target_id = ? AND tb.tenant_id = ? AND t.tenant_id = ? AND t.type = 'app';", a.ID, a.TenantID, a.TenantID).Scan(&ret).Error
	if err != nil {
		mlog.Errorf("get tags from mysql failed:%v", err)
		return nil
	}
	return ret
}

func (a *App) Workflow() *Workflow {
	if a.WorkflowID != "" {
		ret := new(Workflow)
		err := dbengine.Instance().DB.Model(&Workflow{}).Where("id = ?", a.WorkflowID).Scan(ret).Error
		if err != nil {
			mlog.Errorf("get Workflow from mysql failed:%v", err)
			return nil
		}
		return ret
	}
	return nil
}

func (a *App) Site() *Site {
	ret := new(Site)
	err := dbengine.Instance().DB.Model(&Site{}).Where("app_id = ?", a.ID).Scan(ret).Error
	if err != nil {
		mlog.Errorf("get Workflsite from mysql failed:%v", err)
		return nil
	}
	return ret
}

func (a *App) Tenant() *Tenant {
	if a.TenantID != "" {
		ret := new(Tenant)
		err := dbengine.Instance().DB.Model(&Tenant{}).Where("id = ?", a.TenantID).Scan(ret).Error
		if err != nil {
			mlog.Errorf("get Tenant from mysql failed:%v", err)
			return nil
		}
		return ret
	}
	return nil
}
func (a *App) ApiBaseUrl() string {
	service_api_url := viper.GetString("service_api_url")
	if service_api_url == "" {
		service_api_url = fmt.Sprintf("http://%s:%d", utils.GetIP(), viper.GetInt("system.addr"))
	}
	return service_api_url + "/v1"
}

func NewApp(args map[string]any) *App {
	app := new(App)
	if len(args) > 0 {
		bindata, _ := json.Marshal(args)
		err := json.Unmarshal(bindata, app)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("create app by map[string]any failed"))
		}
		return app
	}
	return app
}

// DifySetup [...]
type DifySetup struct {
	Version string     `gorm:"primaryKey;column:version;type:varchar(255);not null" json:"version"`
	SetupAt *time.Time `gorm:"column:setup_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"setup_at"`
}

// TableName get sql table name.获取数据库表名
func (DifySetup) TableName() string {
	return "dify_setups"
}

// AppModelConfig [...]
type AppModelConfig struct {
	ID                            string                `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID                         string                `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	Provider                      string                `gorm:"column:provider;type:varchar(255)" json:"provider"`
	ModelID                       string                `gorm:"column:model_id;type:varchar(255)" json:"model_id"`
	Configs                       datatypes.JSON        `gorm:"column:configs;type:json" json:"configs"`
	CreatedAt                     *time.Time            `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                     *time.Time            `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	OpeningStatement              string                `gorm:"column:opening_statement;type:text" json:"opening_statement"`
	SuggestedQuestions            string                `gorm:"column:suggested_questions;type:text" json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer string                `gorm:"column:suggested_questions_after_answer;type:text" json:"suggested_questions_after_answer"`
	MoreLikeThis                  string                `gorm:"column:more_like_this;type:text" json:"more_like_this"`
	Model                         string                `gorm:"column:model;type:text" json:"model"`
	UserInputForm                 string                `gorm:"column:user_input_form;type:text" json:"user_input_form"`
	PrePrompt                     string                `gorm:"column:pre_prompt;type:text" json:"pre_prompt"`
	AgentMode                     string                `gorm:"column:agent_mode;type:text" json:"agent_mode"`
	SpeechToText                  string                `gorm:"column:speech_to_text;type:text" json:"speech_to_text"`
	SensitiveWordAvoidance        string                `gorm:"column:sensitive_word_avoidance;type:text" json:"sensitive_word_avoidance"`
	RetrieverResource             string                `gorm:"column:retriever_resource;type:text" json:"retriever_resource"`
	DatasetQueryVariable          string                `gorm:"column:dataset_query_variable;type:varchar(255)" json:"dataset_query_variable"`
	PromptType                    string                `gorm:"column:prompt_type;type:varchar(255);default:simple" json:"prompt_type"`
	ChatPromptConfig              string                `gorm:"column:chat_prompt_config;type:text" json:"chat_prompt_config"`
	CompletionPromptConfig        string                `gorm:"column:completion_prompt_config;type:text" json:"completion_prompt_config"`
	DatasetConfigs                string                `gorm:"column:dataset_configs;type:text" json:"dataset_configs"`
	ExternalDataTools             string                `gorm:"column:external_data_tools;type:text" json:"external_data_tools"`
	FileUpload                    string                `gorm:"column:file_upload;type:text" json:"file_upload"`
	TextToSpeech                  string                `gorm:"column:text_to_speech;type:text" json:"text_to_speech"`
	CreatedBy                     string                `gorm:"column:created_by;type:varchar(36)" json:"created_by"`
	UpdatedBy                     string                `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	AnnotationSetting             *AppAnnotationSetting `gorm:"ForeignKey:AppID;references:AppID" json:"annotation_setting"`
}

func NewAppModelConfig(args map[string]any) *AppModelConfig {
	app := new(AppModelConfig)
	if len(args) > 0 {
		bindata, _ := json.Marshal(args)
		err := json.Unmarshal(bindata, app)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("create AppModelConfig by map[string]any failed"))
		}
		return app
	}
	return app
}

//     @property
//     def app(self):
//         app = db.session.query(App).filter(App.id == self.app_id).first()
//         return app

func (amc *AppModelConfig) ModelDict() map[string]any {
	if amc.Model != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.Model), &dic)
		if err != nil {
			mlog.Errorf("Model(%s) json unmarshal() to map failed:%v", amc.Model, err)
		} else {
			return dic
		}
	}
	return map[string]any{}
}

//         return json.loads(self.model) if self.model else {}

func (amc *AppModelConfig) SuggestedQuestionsList() []any {
	if amc.SuggestedQuestions != "" {
		qustions := []any{}
		err := json.Unmarshal([]byte(amc.SuggestedQuestions), &qustions)
		if err != nil {
			mlog.Errorf("SuggestedQuestions(%s) json unmarshal() to list failed:%v", amc.SuggestedQuestions, err)
		} else {
			return qustions
		}
	}
	return nil

}

func (amc *AppModelConfig) SuggestedQuestionsAfterAnswerDict() map[string]any {
	if amc.SuggestedQuestionsAfterAnswer != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.SuggestedQuestionsAfterAnswer), &dic)
		if err != nil {
			mlog.Errorf("SuggestedQuestionsAfterAnswer(%s) json unmarshal() to map failed:%v", amc.SuggestedQuestionsAfterAnswer, err)
		} else {
			return dic
		}
	}
	return map[string]any{"enabled": false}
}

func (amc *AppModelConfig) SpeechToTextDict() map[string]any {
	if amc.SpeechToText != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.SpeechToText), &dic)
		if err != nil {
			mlog.Errorf("SpeechToText(%s) json unmarshal() to map failed:%v", amc.SpeechToText, err)
		} else {
			return dic
		}
	}
	return map[string]any{"enabled": false}
}

func (amc *AppModelConfig) TextToSpeechDict() map[string]any {
	if amc.TextToSpeech != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.TextToSpeech), &dic)
		if err != nil {
			mlog.Errorf("TextToSpeech(%s) json unmarshal() to map failed:%v", amc.TextToSpeech, err)
		} else {
			return dic
		}
	}
	return map[string]any{"enabled": false}

}

func (amc *AppModelConfig) RetrieverResourceDict() map[string]any {
	if amc.RetrieverResource != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.RetrieverResource), &dic)
		if err != nil {
			mlog.Errorf("RetrieverResource(%s) json unmarshal() to map failed:%v", amc.RetrieverResource, err)
		} else {
			return dic
		}
	}
	return map[string]any{"enabled": true}

	// return json.loads(self.retriever_resource) if self.retriever_resource else {"enabled": true}
}

func (amc *AppModelConfig) AnnotationReplyDict() map[string]any {
	if amc.AnnotationSetting != nil {
		collection_binding_detail := amc.AnnotationSetting.CollectionBindingDetail
		return map[string]any{
			"id":              amc.AnnotationSetting.ID,
			"enabled":         true,
			"score_threshold": amc.AnnotationSetting.ScoreThreshold,
			"embedding_model": map[string]any{
				"embedding_provider_name": collection_binding_detail.ProviderName,
				"embedding_model_name":    collection_binding_detail.ModelName,
			},
		}

	}
	return map[string]any{"enabled": false}

}

func (amc *AppModelConfig) MoreLikeThisDict() map[string]any {
	if amc.MoreLikeThis != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.MoreLikeThis), &dic)
		if err != nil {
			mlog.Errorf("MoreLikeThis(%s) json unmarshal() to map failed:%v", amc.MoreLikeThis, err)
		} else {
			return dic
		}
	}
	return map[string]any{"enabled": false}
	//         return json.loads(self.more_like_this) if self.more_like_this else {"enabled": false}
}

func (amc *AppModelConfig) SensitiveWordAvoidanceDict() map[string]any {
	if amc.SensitiveWordAvoidance != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.SensitiveWordAvoidance), &dic)
		if err != nil {
			mlog.Errorf("SensitiveWordAvoidance(%s) json unmarshal() to map failed:%v", amc.SensitiveWordAvoidance, err)
		} else {
			return dic
		}
	}
	return map[string]any{"enabled": false, "type": "", "configs": []any{}}
	// }
	//         return (
	//             json.loads(self.sensitive_word_avoidance)
	//             if self.sensitive_word_avoidance
	//             else {"enabled": false, "type": "", "configs": []}
	//         )

}
func (amc *AppModelConfig) ExternalDataToolsList() []map[string]any {
	if amc.ExternalDataTools != "" {
		list := []map[string]any{}
		err := json.Unmarshal([]byte(amc.ExternalDataTools), &list)
		if err != nil {
			mlog.Errorf("ExternalDataTools(%s) json unmarshal() to map failed:%v", amc.ExternalDataTools, err)
		} else {
			return list
		}
	}
	return []map[string]any{}

	// return json.loads(self.external_data_tools) if self.external_data_tools else []

}
func (amc *AppModelConfig) UserInputFormList() []map[string]any {
	if amc.UserInputForm != "" {
		list := []map[string]any{}
		err := json.Unmarshal([]byte(amc.UserInputForm), &list)
		if err != nil {
			mlog.Errorf("UserInputForm(%s) json unmarshal() to map failed:%v", amc.UserInputForm, err)
		} else {
			return list
		}
	}
	return []map[string]any{}
	// return json.loads(self.user_input_form) if self.user_input_form else []

}
func (amc *AppModelConfig) ChatPromptConfigDict() map[string]any {
	if amc.ChatPromptConfig != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.ChatPromptConfig), &dic)
		if err != nil {
			mlog.Errorf("ChatPromptConfig(%s) json unmarshal() to map failed:%v", amc.ChatPromptConfig, err)
		} else {
			return dic
		}
	}
	return map[string]any{}
	// }
	//         return json.loads(self.chat_prompt_config) if self.chat_prompt_config else {}

}
func (amc *AppModelConfig) CompletionPromptConfigDict() map[string]any {
	if amc.CompletionPromptConfig != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.CompletionPromptConfig), &dic)
		if err != nil {
			mlog.Errorf("CompletionPromptConfig(%s) json unmarshal() to map failed:%v", amc.CompletionPromptConfig, err)
		} else {
			return dic
		}
	}
	return map[string]any{}
	// }
	//         return json.loads(self.completion_prompt_config) if self.completion_prompt_config else {}

}
func (amc *AppModelConfig) DatasetConfigsDict() map[string]any {
	if amc.DatasetConfigs != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.DatasetConfigs), &dic)
		if err != nil {
			mlog.Errorf("DatasetConfigs(%s) json unmarshal() to map failed:%v", amc.DatasetConfigs, err)
		} else {
			if _, ok := dic["retrieval_model"]; !ok {
				return map[string]any{"retrieval_model": "single"}
			}
			return dic
		}
	}
	return map[string]any{
		"retrieval_model": "multiple",
	}
	// }
	//         if self.dataset_configs:
	//             dataset_configs: dict = json.loads(self.dataset_configs)
	//             if "retrieval_model" not in dataset_configs:
	//                 return {"retrieval_model": "single"}
	//             else:
	//                 return dataset_configs
	//         return {
	//             "retrieval_model": "multiple",
	//         }

}
func (amc *AppModelConfig) FileUploadDict() map[string]any {
	if amc.FileUpload != "" {
		dic := map[string]any{}
		err := json.Unmarshal([]byte(amc.FileUpload), &dic)
		if err != nil {
			mlog.Errorf("FileUpload(%s) json unmarshal() to map failed:%v", amc.FileUpload, err)
		} else {
			return dic
		}
	}
	return map[string]any{
		"image": map[string]any{
			"enabled":          false,
			"number_limits":    3,
			"detail":           "high",
			"transfer_methods": []any{"remote_url", "local_file"},
		},
	}
}

func (amc *AppModelConfig) FromModelConfigDict(model_config map[string]any) {
	if _, ok := model_config["opening_statement"]; ok {
		if _, ok := model_config["opening_statement"].(string); ok {
			amc.OpeningStatement = model_config["opening_statement"].(string)
		}
	}
	if _, ok := model_config["suggested_questions"]; ok {
		bindata, _ := json.Marshal(model_config["suggested_questions"])
		amc.SuggestedQuestions = string(bindata)
	}
	if _, ok := model_config["suggested_questions_after_answer"]; ok {
		bindata, _ := json.Marshal(model_config["suggested_questions_after_answer"])
		amc.SuggestedQuestionsAfterAnswer = string(bindata)
	}
	if _, ok := model_config["speech_to_text"]; ok {
		bindata, _ := json.Marshal(model_config["speech_to_text"])
		amc.SpeechToText = string(bindata)
	}
	if _, ok := model_config["text_to_speech"]; ok {
		bindata, _ := json.Marshal(model_config["text_to_speech"])
		amc.TextToSpeech = string(bindata)
	}
	if _, ok := model_config["more_like_this"]; ok {
		bindata, _ := json.Marshal(model_config["more_like_this"])
		amc.MoreLikeThis = string(bindata)
	}
	if _, ok := model_config["sensitive_word_avoidance"]; ok {
		bindata, _ := json.Marshal(model_config["sensitive_word_avoidance"])
		amc.SensitiveWordAvoidance = string(bindata)
	}
	if _, ok := model_config["external_data_tools"]; ok {
		bindata, _ := json.Marshal(model_config["external_data_tools"])
		amc.ExternalDataTools = string(bindata)
	}
	if _, ok := model_config["model"]; ok {
		bindata, _ := json.Marshal(model_config["model"])
		amc.Model = string(bindata)
	}
	if _, ok := model_config["user_input_form"]; ok {
		bindata, _ := json.Marshal(model_config["user_input_form"])
		amc.UserInputForm = string(bindata)
	}
	if _, ok := model_config["dataset_query_variable"]; ok {
		if _, ok := model_config["dataset_query_variable"].(string); ok {
			amc.DatasetQueryVariable = model_config["dataset_query_variable"].(string)
		}
	}
	if _, ok := model_config["pre_prompt"]; ok {
		if _, ok := model_config["pre_prompt"].(string); ok {
			amc.PrePrompt = model_config["pre_prompt"].(string)
		}
	}
	if _, ok := model_config["agent_mode"]; ok {
		bindata, _ := json.Marshal(model_config["agent_mode"])
		amc.AgentMode = string(bindata)
	}
	if _, ok := model_config["retriever_resource"]; ok {
		bindata, _ := json.Marshal(model_config["retriever_resource"])
		amc.RetrieverResource = string(bindata)
	}
	amc.PromptType = "simple"
	if _, ok := model_config["prompt_type"]; ok {
		if _, ok := model_config["prompt_type"].(string); ok {
			amc.PromptType = model_config["prompt_type"].(string)
		}
	}
	if _, ok := model_config["chat_prompt_config"]; ok {
		bindata, _ := json.Marshal(model_config["chat_prompt_config"])
		amc.ChatPromptConfig = string(bindata)
	}
	if _, ok := model_config["completion_prompt_config"]; ok {
		bindata, _ := json.Marshal(model_config["completion_prompt_config"])
		amc.CompletionPromptConfig = string(bindata)
	}
	if _, ok := model_config["dataset_configs"]; ok {
		bindata, _ := json.Marshal(model_config["dataset_configs"])
		amc.DatasetConfigs = string(bindata)
	}
	if _, ok := model_config["file_upload"]; ok {
		bindata, _ := json.Marshal(model_config["file_upload"])
		amc.FileUpload = string(bindata)
	}
}

func (amc *AppModelConfig) ToDict() map[string]any {
	return map[string]any{
		"opening_statement":                amc.OpeningStatement,
		"suggested_questions":              amc.SuggestedQuestionsList(),
		"suggested_questions_after_answer": amc.SuggestedQuestionsAfterAnswerDict(),
		"speech_to_text":                   amc.SpeechToTextDict(),
		"text_to_speech":                   amc.TextToSpeechDict(),
		"retriever_resource":               amc.RetrieverResourceDict(),
		"annotation_reply":                 amc.AnnotationReplyDict(),
		"more_like_this":                   amc.MoreLikeThisDict(),
		"sensitive_word_avoidance":         amc.SensitiveWordAvoidanceDict(),
		"external_data_tools":              amc.ExternalDataToolsList(),
		"model":                            amc.ModelDict(),
		"user_input_form":                  amc.UserInputFormList(),
		"dataset_query_variable":           amc.DatasetQueryVariable,
		"pre_prompt":                       amc.PrePrompt,
		"agent_mode":                       amc.AgentModeDict(),
		"prompt_type":                      amc.PromptType,
		"chat_prompt_config":               amc.ChatPromptConfigDict(),
		"completion_prompt_config":         amc.CompletionPromptConfigDict(),
		"dataset_configs":                  amc.DatasetConfigsDict(),
		"file_upload":                      amc.FileUploadDict(),
	}
}

// TableName get sql table name.获取数据库表名
func (AppModelConfig) TableName() string {
	return "app_model_configs"
}

func (ac *AppModelConfig) AgentModeDict() map[string]any {
	ref := map[string]any{"enabled": false, "strategy": "", "tools": []map[string]any{}, "prompt": ""}
	if ac.AgentMode != "" {
		ret := make(map[string]any)
		err := json.Unmarshal([]byte(ac.AgentMode), &ret)
		if err != nil {
			return ref
		} else {
			return ret
		}
	}
	return ref
}

// RecommendedApp [...]
type RecommendedApp struct {
	ID               string         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID            string         `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App              *App           `gorm:"ForeignKey:AppID;AssociationForeignKey:ID" json:"app"`
	Description      datatypes.JSON `gorm:"column:description;type:json;not null" json:"description"`
	Copyright        string         `gorm:"column:copyright;type:varchar(255);not null" json:"copyright"`
	PrivacyPolicy    string         `gorm:"column:privacy_policy;type:varchar(255);not null" json:"privacy_policy"`
	Category         string         `gorm:"column:category;type:varchar(255);not null" json:"category"`
	Position         int            `gorm:"column:position;type:int;not null" json:"position"`
	IsListed         bool           `gorm:"column:is_listed;type:tinyint(1);not null" json:"is_listed"`
	InstallCount     int            `gorm:"column:install_count;type:int;not null" json:"install_count"`
	CreatedAt        *time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Language         string         `gorm:"column:language;type:varchar(255);default:en-US" json:"language"`
	CustomDisclaimer string         `gorm:"column:custom_disclaimer;type:text;not null" json:"custom_disclaimer"`
}

// TableName get sql table name.获取数据库表名
func (RecommendedApp) TableName() string {
	return "recommended_apps"
}

// InstalledApp [...]
type InstalledApp struct {
	ID               string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID         string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant           *Tenant    `gorm:"ForeignKey:TenantID;AssociationForeignKey:ID" json:"tenant"`
	AppID            string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App              *App       `gorm:"ForeignKey:AppID;AssociationForeignKey:ID" json:"app"`
	AppOwnerTenantID string     `gorm:"column:app_owner_tenant_id;type:varchar(36);not null" json:"app_owner_tenant_id"`
	Position         int        `gorm:"column:position;type:int;not null" json:"position"`
	IsPinned         bool       `gorm:"column:is_pinned;type:tinyint(1);not null;default:0" json:"is_pinned"`
	LastUsedAt       *time.Time `gorm:"column:last_used_at;type:timestamp" json:"last_used_at"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (InstalledApp) TableName() string {
	return "installed_apps"
}

// Conversation [...]
type Conversation struct {
	ID                      string          `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID                   string          `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App                     *App            `gorm:"ForeignKey:AppID;AssociationForeignKey:ID" json:"app"`
	AppModelConfigID        string          `gorm:"column:app_model_config_id;type:varchar(36)" json:"app_model_config_id"`
	AppModelConfig          *AppModelConfig `gorm:"ForeignKey:AppModelConfigID;AssociationForeignKey:ID" json:"app_model_config"`
	ModelProvider           string          `gorm:"column:model_provider;type:varchar(255)" json:"model_provider"`
	OverrideModelConfigsStr string          `gorm:"column:override_model_configs;type:text" json:"override_model_configs"`
	ModelID                 string          `gorm:"column:model_id;type:varchar(255)" json:"model_id"`
	Mode                    AppMode         `gorm:"column:mode;type:varchar(255);not null" json:"mode"`
	Name                    string          `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Summary                 string          `gorm:"column:summary;type:text" json:"summary"`
	InputsJson              string          `gorm:"column:inputs;" json:"inputs"`

	Introduction            string     `gorm:"column:introduction;type:text" json:"introduction"`
	SystemInstruction       string     `gorm:"column:system_instruction;type:text" json:"system_instruction"`
	SystemInstructionTokens int        `gorm:"column:system_instruction_tokens;type:int;not null;default:0" json:"system_instruction_tokens"`
	Status                  string     `gorm:"column:status;type:varchar(255);not null" json:"status"`
	FromSource              string     `gorm:"column:from_source;type:varchar(255);not null" json:"from_source"`
	FromEndUserID           string     `gorm:"column:from_end_user_id;type:varchar(36)" json:"from_end_user_id"`
	FromAccountID           string     `gorm:"column:from_account_id;type:varchar(36)" json:"from_account_id"`
	ReadAt                  *time.Time `gorm:"column:read_at;type:timestamp" json:"read_at"`
	ReadAccountID           string     `gorm:"column:read_account_id;type:varchar(36)" json:"read_account_id"`
	CreatedAt               *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsDeleted               bool       `gorm:"column:is_deleted;type:tinyint(1);not null;default:0" json:"is_deleted"`
	InvokeFrom              string     `gorm:"column:invoke_from;type:varchar(255)" json:"invoke_from"`
	DialogueCount           int        `gorm:"column:dialogue_count;type:int;not null;default:0" json:"dialogue_count"`
}

// TableName get sql table name.获取数据库表名
func (Conversation) TableName() string {
	return "conversations"
}

func (c *Conversation) Inputs() map[string]any {
	var inputs map[string]any
	if c.InputsJson != "" {
		json.Unmarshal([]byte(c.InputsJson), &inputs)
	}
	// // Convert file mapping to File object
	// for key, value := range inputs{
	//     // NOTE: It's not the best way to implement this, but it's the only way to avoid circular import for now.
	//     from factories import file_factory

	//     if isinstance(value, dict) and value.get("dify_model_identity") == FILE_MODEL_IDENTITY:
	//         if value["transfer_method"] == FileTransferMethod.TOOL_FILE:
	//             value["tool_file_id"] = value["related_id"]
	//         }else if value["transfer_method"] == FileTransferMethod.LOCAL_FILE:
	//             value["upload_file_id"] = value["related_id"]
	//         inputs[key] = file_factory.build_from_mapping(mapping=value, tenant_id=value["tenant_id"])
	//     }else if isinstance(value, list) and all(
	//         isinstance(item, dict) and item.get("dify_model_identity") == FILE_MODEL_IDENTITY for item in value
	//     ):
	//         inputs[key] = []
	//         for item in value:
	//             if item["transfer_method"] == FileTransferMethod.TOOL_FILE:
	//                 item["tool_file_id"] = item["related_id"]
	//             }else if item["transfer_method"] == FileTransferMethod.LOCAL_FILE:
	//                 item["upload_file_id"] = item["related_id"]
	//             inputs[key].append(file_factory.build_from_mapping(mapping=item, tenant_id=item["tenant_id"]))
	// 		}
	// 	}
	// }
	return inputs
}
func (c *Conversation) SetInputs(val map[string]any) {
	if len(val) > 0 {
		bindata, _ := json.Marshal(val)
		c.InputsJson = string(bindata)
	} else {
		c.InputsJson = "{}"
	}
	// inputs = dict(value)
	// for k, v in inputs.items():
	//     if isinstance(v, File):
	//         inputs[k] = v.model_dump()
	//     }else if isinstance(v, list) and all(isinstance(item, File) for item in v):
	//         inputs[k] = [item.model_dump() for item in v]
	// self._inputs = inputs
}

func (c *Conversation) OverrideModelConfigs() map[string]any {
	var data map[string]any
	if c.OverrideModelConfigsStr != "" {
		json.Unmarshal([]byte(c.OverrideModelConfigsStr), &data)
	}
	return data
}
func (c *Conversation) SetOverrideModelConfigs(val map[string]any) {
	if len(val) > 0 {
		bindata, _ := json.Marshal(val)
		c.OverrideModelConfigsStr = string(bindata)
	} else {
		c.OverrideModelConfigsStr = "{}"
	}
}
func (c *Conversation) ModelConfig() map[string]any {
	model_config := map[string]any{}
	var app_model_config *AppModelConfig

	if c.Mode == AppMode_ADVANCED_CHAT {
		if c.OverrideModelConfigsStr != "" {
			json.Unmarshal([]byte(c.OverrideModelConfigsStr), &model_config)
		}
	} else {
		if c.OverrideModelConfigsStr != "" {
			var override_model_configs map[string]any
			json.Unmarshal([]byte(c.OverrideModelConfigsStr), &override_model_configs)

			if _, ok := override_model_configs["model"]; ok {
				app_model_config = new(AppModelConfig)
				app_model_config.FromModelConfigDict(override_model_configs)
				if app_model_config == nil {
					mlog.Error("app model config not found")
				} else {
					model_config = app_model_config.ToDict()
				}
			} else {
				model_config["configs"] = override_model_configs
			}
		} else {
			app_model_config = new(AppModelConfig)
			err := dbengine.Instance().DB.Model(&AppModelConfig{}).Where("id = ?", c.AppModelConfigID).First(app_model_config).Error
			if err != nil {
				mlog.Errorf("get AppModelConfig failed:%v", err)
				app_model_config = nil
			}
			if app_model_config != nil {
				model_config = app_model_config.ToDict()
			}
		}
	}
	model_config["model_id"] = c.ModelID
	model_config["provider"] = c.ModelProvider

	return model_config
}

// @property
// def summary_or_query(self):
//     if c.summary:
//         return c.summary
//     else:
//         first_message = c.first_message
//         if first_message:
//             return first_message.query
//         else:
//             return ""

func (c *Conversation) Annotated() bool {
	var count int64
	err := dbengine.Instance().DB.Model(&MessageAnnotation{}).Where("conversation_id =?", c.ID).Count(&count).Error
	if err != nil {
		mlog.Errorf("count MessageAnnotation failed:%v", err)
	}

	return count > 0
}

// @property
// def annotation(self):
//     return db.session.query(MessageAnnotation).filter(MessageAnnotation.conversation_id == c.id).first()

func (c *Conversation) MessageCount() int64 {
	var count int64
	err := dbengine.Instance().DB.Model(&Message{}).Where("conversation_id =?", c.ID).Count(&count).Error
	if err != nil {
		mlog.Errorf("count Message failed:%v", err)
	}
	return count
}

func (c *Conversation) UserFeedbackStats() (int64, int64) {
	var like int64
	err := dbengine.Instance().DB.Model(&MessageFeedback{}).Where("conversation_id =? and from_source = ? and rating = ?", c.ID, "user", "like").Count(&like).Error
	if err != nil {
		mlog.Errorf("count MessageFeedback failed:%v", err)
	}
	var dislike int64
	err = dbengine.Instance().DB.Model(&MessageFeedback{}).Where("conversation_id =? and from_source = ? and rating = ?", c.ID, "user", "dislike").Count(&dislike).Error
	if err != nil {
		mlog.Errorf("count MessageFeedback failed:%v", err)
	}

	return like, dislike
}
func (c *Conversation) AdminFeedbackStats() (int64, int64) {
	var like int64
	err := dbengine.Instance().DB.Model(&MessageFeedback{}).Where("conversation_id =? and from_source = ? and rating = ?", c.ID, "admin", "like").Count(&like).Error
	if err != nil {
		mlog.Errorf("count MessageFeedback failed:%v", err)
	}
	var dislike int64
	err = dbengine.Instance().DB.Model(&MessageFeedback{}).Where("conversation_id =? and from_source = ? and rating = ?", c.ID, "admin", "dislike").Count(&dislike).Error
	if err != nil {
		mlog.Errorf("count MessageFeedback failed:%v", err)
	}

	return like, dislike
}

func (c *Conversation) StatusCount() map[string]int {
	var messages []*Message
	err := dbengine.Instance().DB.Model(&Message{}).Where("conversation_id =?", c.ID).Find(&messages).Error
	if err != nil {
		mlog.Errorf("count Message failed:%v", err)
	}
	status_counts := map[WorkflowRunStatus]int{
		WorkflowRunStatus_RUNNING:           0,
		WorkflowRunStatus_SUCCEEDED:         0,
		WorkflowRunStatus_FAILED:            0,
		WorkflowRunStatus_STOPPED:           0,
		WorkflowRunStatus_PARTIAL_SUCCESSED: 0,
	}

	for _, message := range messages {
		workflow_run := message.WorkflowRun()
		if workflow_run != nil {
			WorkflowRunStatus(workflow_run.Status).Validate()
			status_counts[WorkflowRunStatus(workflow_run.Status)] += 1
		}
	}
	if len(messages) > 0 {
		return map[string]int{
			"success":         status_counts[WorkflowRunStatus_SUCCEEDED],
			"failed":          status_counts[WorkflowRunStatus_FAILED],
			"partial_success": status_counts[WorkflowRunStatus_PARTIAL_SUCCESSED],
		}
	} else {
		return nil
	}
}

// @property
// def first_message(self):
//     return db.session.query(Message).filter(Message.conversation_id == c.id).first()

// @property
// def app(self):
//     return db.session.query(App).filter(App.id == c.app_id).first()

func (c *Conversation) FromEndUserSessionID() string {
	if c.FromEndUserID != "" {
		end_user := new(EndUser)
		err := dbengine.Instance().DB.Model(&EndUser{}).Where("id = ?", c.FromEndUserID).First(end_user).Error
		if err != nil {
			mlog.Errorf("get EndUser from mysql failed:%v", err)
			end_user = nil
		}

		if end_user != nil {
			return end_user.SessionID
		}
	}
	return ""
}
func (c *Conversation) FromAccountName() string {
	if c.FromAccountID != "" {
		account := new(Account)
		err := dbengine.Instance().DB.Model(&Account{}).Where("id = ?", c.FromAccountID).First(account).Error
		if err != nil {
			mlog.Errorf("get Account from mysql failed:%v", err)
			account = nil
		}

		if account != nil {
			return account.Name
		}
	}
	return ""
}
func (c *Conversation) InDebugMode() bool {
	return len(c.OverrideModelConfigs()) > 0
}

// Message [...]
type Message struct {
	ID                      string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID                   string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App                     *App       `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	ModelProvider           string     `gorm:"column:model_provider;type:varchar(255)" json:"model_provider"`
	ModelID                 string     `gorm:"column:model_id;type:varchar(255)" json:"model_id"`
	OverrideModelConfigsStr string     `gorm:"column:override_model_configs;type:text" json:"override_model_configs"`
	ConversationID          string     `gorm:"column:conversation_id;type:varchar(36);not null" json:"conversation_id"`
	InputsStr               string     `gorm:"column:inputs;" json:"inputs"`
	Query                   string     `gorm:"column:query;type:text;not null" json:"query"`
	MessageJson             string     `gorm:"column:message" json:"message"`
	MessageTokens           int        `gorm:"column:message_tokens;type:int;not null;default:0" json:"message_tokens"`
	MessageUnitPrice        float64    `gorm:"column:message_unit_price;type:decimal(10,4);not null" json:"message_unit_price"`
	Answer                  string     `gorm:"column:answer;type:text;not null" json:"answer"`
	AnswerTokens            int        `gorm:"column:answer_tokens;type:int;not null;default:0" json:"answer_tokens"`
	AnswerUnitPrice         float64    `gorm:"column:answer_unit_price;type:decimal(10,4);not null" json:"answer_unit_price"`
	ProviderResponseLatency float64    `gorm:"column:provider_response_latency;type:double;not null;default:0" json:"provider_response_latency"`
	TotalPrice              float64    `gorm:"column:total_price;type:decimal(10,7)" json:"total_price"`
	Currency                string     `gorm:"column:currency;type:varchar(255);not null" json:"currency"`
	FromSource              string     `gorm:"column:from_source;type:varchar(255);not null" json:"from_source"`
	FromEndUserID           string     `gorm:"column:from_end_user_id;type:varchar(36)" json:"from_end_user_id"`
	FromAccountID           string     `gorm:"column:from_account_id;type:varchar(36)" json:"from_account_id"`
	CreatedAt               *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	AgentBased              bool       `gorm:"column:agent_based;type:tinyint(1);not null;default:0" json:"agent_based"`
	MessagePriceUnit        float64    `gorm:"column:message_price_unit;type:decimal(10,7);not null;default:0.0010000" json:"message_price_unit"`
	AnswerPriceUnit         float64    `gorm:"column:answer_price_unit;type:decimal(10,7);not null;default:0.0010000" json:"answer_price_unit"`
	WorkflowRunID           string     `gorm:"column:workflow_run_id;type:varchar(36)" json:"workflow_run_id"`
	// WorkflowRun             *WorkflowRun `json:"workflow_run" form:"workflow_run" gorm:"foreignKey:WorkflowRunID;references:ID;"`
	Status             string `gorm:"column:status;type:varchar(255);default:normal" json:"status"`
	Error              string `gorm:"column:error;type:text" json:"error"`
	MessageMetadataStr string `gorm:"column:message_metadata;type:text" json:"message_metadata"`
	InvokeFrom         string `gorm:"column:invoke_from;type:varchar(255)" json:"invoke_from"`
	ParentMessageID    string `gorm:"column:parent_message_id;type:varchar(36)" json:"parent_message_id"`
}

// TableName get sql table name.获取数据库表名
func (Message) TableName() string {
	return "messages"
}

func (msg *Message) Inputs() map[string]any {
	inputs := map[string]any{}
	if msg.InputsStr != "" {
		json.Unmarshal([]byte(msg.InputsStr), &inputs)
	}
	// for key, value in inputs.items():
	//     // NOTE: It's not the best way to implement this, but it's the only way to avoid circular import for now.
	//     from factories import file_factory

	//     if isinstance(value, dict) and value.get("dify_model_identity") == FILE_MODEL_IDENTITY:
	//         if value["transfer_method"] == FileTransferMethod.TOOL_FILE:
	//             value["tool_file_id"] = value["related_id"]
	//         }else if value["transfer_method"] == FileTransferMethod.LOCAL_FILE:
	//             value["upload_file_id"] = value["related_id"]
	//         inputs[key] = file_factory.build_from_mapping(mapping=value, tenant_id=value["tenant_id"])
	//     }else if isinstance(value, list) and all(
	//         isinstance(item, dict) and item.get("dify_model_identity") == FILE_MODEL_IDENTITY for item in value
	//     ):
	//         inputs[key] = []
	//         for item in value:
	//             if item["transfer_method"] == FileTransferMethod.TOOL_FILE:
	//                 item["tool_file_id"] = item["related_id"]
	//             }else if item["transfer_method"] == FileTransferMethod.LOCAL_FILE:
	//                 item["upload_file_id"] = item["related_id"]
	//             inputs[key].append(file_factory.build_from_mapping(mapping=item, tenant_id=item["tenant_id"]))
	return inputs
}
func (msg *Message) SetInputs(val map[string]any) {
	if len(val) > 0 {
		bindata, _ := json.Marshal(val)
		msg.InputsStr = string(bindata)
	} else {
		msg.InputsStr = "{}"
	}
}

func (msg *Message) OverrideModelConfigs() map[string]any {
	data := map[string]any{}
	if msg.InputsStr != "" {
		json.Unmarshal([]byte(msg.OverrideModelConfigsStr), &data)
	}
	return data
}
func (msg *Message) SetOverrideModelConfigs(val map[string]any) {
	if len(val) > 0 {
		bindata, _ := json.Marshal(val)
		msg.OverrideModelConfigsStr = string(bindata)
	} else {
		msg.OverrideModelConfigsStr = "{}"
	}
}

func (msg *Message) MessageMetadata() map[string]any {
	data := map[string]any{}
	if msg.InputsStr != "" {
		json.Unmarshal([]byte(msg.MessageMetadataStr), &data)
	}
	return data
}
func (msg *Message) SetMessageMetadata(val map[string]any) {
	if len(val) > 0 {
		bindata, _ := json.Marshal(val)
		msg.MessageMetadataStr = string(bindata)
	} else {
		msg.MessageMetadataStr = "{}"
	}
}

//     @property
//     def re_sign_file_url_answer(self) -> str:
//         if not self.answer:
//             return self.answer

//         pattern = r"\[!?.*?\]\((((http|https):\/\/.+)?\/files\/(tools\/)?[\w-]+.*?timestamp=.*&nonce=.*&sign=.*)\)"
//         matches = re.findall(pattern, self.answer)

//         if not matches:
//             return self.answer

//         urls = [match[0] for match in matches]

//         // remove duplicate urls
//         urls = list(set(urls))

//         if not urls:
//             return self.answer

//         re_sign_file_url_answer = self.answer
//         for url in urls:
//             if "files/tools" in url:
//                 // get tool file id
//                 tool_file_id_pattern = r"\/files\/tools\/([\.\w-]+)?\?timestamp="
//                 result = re.search(tool_file_id_pattern, url)
//                 if not result:
//                     continue

//                 tool_file_id = result.group(1)

//                 // get extension
//                 if "." in tool_file_id:
//                     split_result = tool_file_id.split(".")
//                     extension = f".{split_result[-1]}"
//                     if len(extension) > 10:
//                         extension = ".bin"
//                     tool_file_id = split_result[0]
//                 else:
//                     extension = ".bin"

//                 if not tool_file_id:
//                     continue

//                 sign_url = ToolFileParser.get_tool_file_manager().sign_file(
//                     tool_file_id=tool_file_id, extension=extension
//                 )
//             }else if "file-preview" in url:
//                 // get upload file id
//                 upload_file_id_pattern = r"\/files\/([\w-]+)\/file-preview?\?timestamp="
//                 result = re.search(upload_file_id_pattern, url)
//                 if not result:
//                     continue

//                 upload_file_id = result.group(1)
//                 if not upload_file_id:
//                     continue
//                 sign_url = file_helpers.get_signed_file_url(upload_file_id)
//             }else if "image-preview" in url:
//                 // image-preview is deprecated, use file-preview instead
//                 upload_file_id_pattern = r"\/files\/([\w-]+)\/image-preview?\?timestamp="
//                 result = re.search(upload_file_id_pattern, url)
//                 if not result:
//                     continue
//                 upload_file_id = result.group(1)
//                 if not upload_file_id:
//                     continue
//                 sign_url = file_helpers.get_signed_file_url(upload_file_id)
//             else:
//                 continue

//             re_sign_file_url_answer = re_sign_file_url_answer.replace(url, sign_url)

//         return re_sign_file_url_answer

//     @property
//     def user_feedback(self):
//         feedback = (
//             db.session.query(MessageFeedback)
//             .filter(MessageFeedback.message_id == self.id, MessageFeedback.from_source == "user")
//             .first()
//         )
//         return feedback

func (msg *Message) AdminFeedback() *MessageFeedback {
	feedback := new(MessageFeedback)
	err := dbengine.Instance().DB.Model(&MessageFeedback{}).Where("message_id=? and from_source = ?", msg.ID, "admin").First(feedback).Error
	if err != nil {
		mlog.Errorf("get MessageFeedback failed:%v", err)
		feedback = nil
	}
	return feedback
}

func (msg *Message) Feedbacks() []*MessageFeedback {
	var feedbacks []*MessageFeedback
	err := dbengine.Instance().DB.Model(&MessageFeedback{}).Where("message_id=?", msg.ID).Find(&feedbacks).Error
	if err != nil {
		mlog.Errorf("get MessageFeedback failed:%v", err)
		feedbacks = nil
	}
	return feedbacks
}
func (msg *Message) Annotation() *MessageAnnotation {
	annotation := new(MessageAnnotation)
	err := dbengine.Instance().DB.Model(&MessageAnnotation{}).Where("message_id=? ", msg.ID).First(annotation).Error
	if err != nil {
		mlog.Errorf("get MessageAnnotation failed:%v", err)
		annotation = nil
	}
	return annotation
}
func (msg *Message) AnnotationHitHistory() *MessageAnnotation {
	annotation_history := new(AppAnnotationHitHistory)
	err := dbengine.Instance().DB.Model(&AppAnnotationHitHistory{}).Where("message_id=? ", msg.ID).First(annotation_history).Error
	if err != nil {
		mlog.Errorf("get AppAnnotationHitHistory failed:%v", err)
		annotation_history = nil
	}

	if annotation_history != nil {
		annotation := new(MessageAnnotation)
		err := dbengine.Instance().DB.Model(&MessageAnnotation{}).Where("id=? ", annotation_history.AnnotationID).First(annotation).Error
		if err != nil {
			mlog.Errorf("get MessageAnnotation failed:%v", err)
			annotation = nil
		}
		return annotation
	}
	return nil

}
func (msg *Message) AppModelConfig() *AppModelConfig {
	conversation := new(Conversation)
	err := dbengine.Instance().DB.Model(&Conversation{}).Where("id=? ", msg.ConversationID).First(conversation).Error
	if err != nil {
		mlog.Errorf("get Conversation failed:%v", err)
		conversation = nil
	}

	if conversation != nil {
		data := new(AppModelConfig)
		err := dbengine.Instance().DB.Model(&AppModelConfig{}).Where("id=? ", conversation.AppModelConfigID).First(data).Error
		if err != nil {
			mlog.Errorf("get AppModelConfig failed:%v", err)
			data = nil
		}

		return data
	}
	return nil

}
func (msg *Message) InDebugMode() bool {
	return len(msg.OverrideModelConfigs()) > 0
}
func (msg *Message) MessageMetadataDict() map[string]any {
	if msg.MessageMetadataStr != "" {
		tmp_dict := map[string]any{}
		err := json.Unmarshal([]byte(msg.MessageMetadataStr), &tmp_dict)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", msg.MessageMetadataStr, err)
		}
		return tmp_dict
	} else {
		return map[string]any{}
	}
}
func (msg *Message) AgentThoughts() []*MessageAgentThought {
	var thoughts []*MessageAgentThought
	err := dbengine.Instance().DB.Model(&MessageAgentThought{}).Where("message_id=?", msg.ID).Order("position ASC").Find(&thoughts).Error
	if err != nil {
		mlog.Errorf("get MessageAgentThought failed:%v", err)
		thoughts = nil
	}
	return thoughts
}
func (msg *Message) RetrieverResources() []*DatasetRetrieverResource {
	var resources []*DatasetRetrieverResource
	err := dbengine.Instance().DB.Model(&DatasetRetrieverResource{}).Where("message_id=?", msg.ID).Order("position ASC").Find(&resources).Error
	if err != nil {
		mlog.Errorf("get DatasetRetrieverResource failed:%v", err)
		resources = nil
	}
	return resources
}

// func (msg *Message) MessageFiles(){
//         // from factories import file_factory

//         var message_files []*MessageFile
// 		err := dbengine.Instance().DB.Model(&MessageFile{}).Where("message_id =? ", msg.ID).Find(&message_files).Error
// 		if err != nil {
// 			mlog.Errorf("get MessageFile failed:%v", err)
// 		}
//         current_app := new(App)
// 		err = dbengine.Instance().DB.Model(&App{}).Where("id =?", msg.AppID).First(current_app).Error
// 				if err != nil {
// 			mlog.Errorf("get App failed:%v", err)
// 			current_app = nil
// 		}
//         if current_app == nil{
//             panic(exceptions.NewValueError(fmt.Sprintf("App {%s} not found", msg.AppID)))
// 		}
//         files = []
//         for message_file in message_files{
//             if message_file.transfer_method == "local_file"{
//                 if message_file.upload_file_id is None{
//                      panic(exceptions.NewValueError(fmt.Sprintf("MessageFile {%s} is a local file but has no upload_file_id", message_file.ID)))
// 				}
//                 file = file_factory.build_from_mapping(
//                     mapping={
//                         "id": message_file.id,
//                         "upload_file_id": message_file.upload_file_id,
//                         "transfer_method": message_file.transfer_method,
//                         "type": message_file.type,
//                     },
//                     tenant_id=current_app.tenant_id,
//                 )
//             }else if message_file.transfer_method == "remote_url"{
//                 if message_file.url is None{
//                      panic(exceptions.NewValueError(fmt.Sprintf("MessageFile {%s} is a remote url but has no url", message_file.ID)))
// 				}
//                 file = file_factory.build_from_mapping(
//                     mapping={
//                         "id": message_file.id,
//                         "type": message_file.type,
//                         "transfer_method": message_file.transfer_method,
//                         "url": message_file.url,
//                     },
//                     tenant_id=current_app.tenant_id,
//                 )
//             }else if message_file.transfer_method == "tool_file"{
//                 if message_file.upload_file_id is None{
//                     assert message_file.url is not None
//                     message_file.upload_file_id = message_file.url.split("/")[-1].split(".")[0]
// 				}
//                 mapping = {
//                     "id": message_file.id,
//                     "type": message_file.type,
//                     "transfer_method": message_file.transfer_method,
//                     "tool_file_id": message_file.upload_file_id,
//                 }
//                 file = filefactory.BuildFromMapping(
//                     mapping=mapping,
//                     tenant_id=current_app.tenant_id,
//                 )
//             else{
//                 panic(exceptions.NewValueError(fmt.Sprintf("MessageFile {%s} has an invalid transfer_method {%s}", message_file.ID, message_file.TransferMethod)))
// 			}
//             files.append(file)
// 		}

//         result = [
//             {"belongs_to": message_file.belongs_to, **file.to_dict()}
//             for (file, message_file) in zip(files, message_files)
//         ]

//	        db.session.commit()
//	        return result
//	}
func (msg *Message) WorkflowRun() *WorkflowRun {
	if msg.WorkflowRunID != "" {
		workflow_run := new(WorkflowRun)
		err := dbengine.Instance().DB.Model(&WorkflowRun{}).Where("id =?", msg.WorkflowRunID).First(workflow_run).Error
		if err != nil {
			mlog.Errorf("count WorkflowRun failed:%v", err)
			workflow_run = nil
		}
		return workflow_run
	}
	return nil
}
func (msg *Message) ToDict() map[string]any {
	return map[string]any{
		"id":               msg.ID,
		"app_id":           msg.AppID,
		"conversation_id":  msg.ConversationID,
		"model_id":         msg.ModelID,
		"inputs":           msg.Inputs,
		"query":            msg.Query,
		"total_price":      msg.TotalPrice,
		"message":          msg.MessageJson,
		"answer":           msg.Answer,
		"status":           msg.Status,
		"error":            msg.Error,
		"message_metadata": msg.MessageMetadataDict(),
		"from_source":      msg.FromSource,
		"from_end_user_id": msg.FromEndUserID,
		"from_account_id":  msg.FromAccountID,
		"created_at":       msg.CreatedAt.Format(time.DateTime),
		"updated_at":       msg.UpdatedAt.Format(time.DateTime),
		"agent_based":      msg.AgentBased,
		"workflow_run_id":  msg.WorkflowRunID,
	}
}

// func NewMessageFromDict(data map[string]any)*Message{

//		return cls(
//		    id=data["id"],
//		    app_id=data["app_id"],
//		    conversation_id=data["conversation_id"],
//		    model_id=data["model_id"],
//		    inputs=data["inputs"],
//		    total_price=data["total_price"],
//		    query=data["query"],
//		    message=data["message"],
//		    answer=data["answer"],
//		    status=data["status"],
//		    error=data["error"],
//		    message_metadata=json.dumps(data["message_metadata"]),
//		    from_source=data["from_source"],
//		    from_end_user_id=data["from_end_user_id"],
//		    from_account_id=data["from_account_id"],
//		    created_at=data["created_at"],
//		    updated_at=data["updated_at"],
//		    agent_based=data["agent_based"],
//		    workflow_run_id=data["workflow_run_id"],
//		)
//	}
//
// MessageFeedback [...]
type MessageFeedback struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App            *App       `gorm:"ForeignKey:AppID;AssociationForeignKey:ID" json:"app"`
	ConversationID string     `gorm:"column:conversation_id;type:varchar(36);not null" json:"conversation_id"`
	MessageID      string     `gorm:"column:message_id;type:varchar(36);not null" json:"message_id"`
	Rating         string     `gorm:"column:rating;type:varchar(255);not null" json:"rating"`
	Content        string     `gorm:"column:content;type:text" json:"content"`
	FromSource     string     `gorm:"column:from_source;type:varchar(255);not null" json:"from_source"`
	FromEndUserID  string     `gorm:"column:from_end_user_id;type:varchar(36)" json:"from_end_user_id"`
	FromAccountID  string     `gorm:"column:from_account_id;type:varchar(36)" json:"from_account_id"`
	Account        *Account   `gorm:"ForeignKey:FromAccountID;AssociationForeignKey:ID" json:"account"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (MessageFeedback) TableName() string {
	return "message_feedbacks"
}
func (mfb *MessageFeedback) FromAccount() *Account {
	account := new(Account)
	err := dbengine.Instance().DB.Model(&Account{}).Where("id =?", mfb.FromAccountID).First(account).Error
	if err != nil {
		mlog.Errorf("get Account failed:%v", err)
		account = nil
	}
	return account
}

// MessageFile [...]
type MessageFile struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	MessageID      string     `gorm:"column:message_id;type:varchar(36);not null" json:"message_id"`
	Type           string     `gorm:"column:type;type:varchar(255);not null" json:"type"`
	TransferMethod string     `gorm:"column:transfer_method;type:varchar(255);not null" json:"transfer_method"`
	URL            string     `gorm:"column:url;type:text" json:"url"`
	UploadFileID   string     `gorm:"column:upload_file_id;type:varchar(36)" json:"upload_file_id"`
	CreatedByRole  string     `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy      string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	BelongsTo      string     `gorm:"column:belongs_to;type:varchar(255)" json:"belongs_to"`
}

// TableName get sql table name.获取数据库表名
func (MessageFile) TableName() string {
	return "message_files"
}

// MessageAnnotation [...]
type MessageAnnotation struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App            *App       `gorm:"ForeignKey:AppID;AssociationForeignKey:ID" json:"app"`
	ConversationID string     `gorm:"column:conversation_id;type:varchar(36)" json:"conversation_id"`
	MessageID      string     `gorm:"column:message_id;type:varchar(36)" json:"message_id"`
	Content        string     `gorm:"column:content;type:text;not null" json:"content"`
	AccountID      string     `gorm:"column:account_id;type:varchar(36);not null" json:"account_id"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Question       string     `gorm:"column:question;type:text" json:"question"`
	HitCount       int        `gorm:"column:hit_count;type:int;not null;default:0" json:"hit_count"`
}

// TableName get sql table name.获取数据库表名
func (MessageAnnotation) TableName() string {
	return "message_annotations"
}
func (ma *MessageAnnotation) Account() *Account {
	account := new(Account)
	err := dbengine.Instance().DB.Model(&Account{}).Where("id =?", ma.AccountID).First(account).Error
	if err != nil {
		mlog.Errorf("get Account failed:%v", err)
		account = nil
	}
	return account
}

func (ma *MessageAnnotation) AnnotationCreateAccount() *Account {
	account := new(Account)
	err := dbengine.Instance().DB.Model(&Account{}).Where("id =?", ma.AccountID).First(account).Error
	if err != nil {
		mlog.Errorf("get Account failed:%v", err)
		account = nil
	}
	return account
}

// AppAnnotationSetting [...]
type AppAnnotationSetting struct {
	ID                      string                    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID                   string                    `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App                     *App                      `gorm:"ForeignKey:AppID;AssociationForeignKey:ID" json:"app"`
	ScoreThreshold          float64                   `gorm:"column:score_threshold;type:double;not null;default:0" json:"score_threshold"`
	CollectionBindingID     string                    `gorm:"column:collection_binding_id;type:varchar(36);not null" json:"collection_binding_id"`
	CreatedUserID           string                    `gorm:"column:created_user_id;type:varchar(36);not null" json:"created_user_id"`
	CreatedAccount          *Account                  `gorm:"ForeignKey:CreatedUserID;AssociationForeignKey:ID" json:"created_account"`
	CreatedAt               *time.Time                `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedUserID           string                    `gorm:"column:updated_user_id;type:varchar(36);not null" json:"updated_user_id"`
	UpdatedAccount          *Account                  `gorm:"ForeignKey:UpdatedUserID;AssociationForeignKey:ID" json:"updated_account"`
	UpdatedAt               *time.Time                `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CollectionBindingDetail *DatasetCollectionBinding `gorm:"ForeignKey:CollectionBindingID;AssociationForeignKey:ID" json:"collection_binding_detail"`
	//     @property
	//     def created_account(self):
	//         account = (
	//             db.session.query(Account)
	//             .join(AppAnnotationSetting, AppAnnotationSetting.created_user_id == Account.id)
	//             .filter(AppAnnotationSetting.id == self.annotation_id)
	//             .first()
	//         )
	//         return account

	//     @property
	//     def updated_account(self):
	//         account = (
	//             db.session.query(Account)
	//             .join(AppAnnotationSetting, AppAnnotationSetting.updated_user_id == Account.id)
	//             .filter(AppAnnotationSetting.id == self.annotation_id)
	//             .first()
	//         )
	//         return account

}

// TableName get sql table name.获取数据库表名
func (AppAnnotationSetting) TableName() string {
	return "app_annotation_settings"
}

// AppAnnotationHitHistory [...]
type AppAnnotationHitHistory struct {
	ID                       string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID                    string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	AnnotationID             string     `gorm:"column:annotation_id;type:varchar(36);not null" json:"annotation_id"`
	Source                   string     `gorm:"column:source;type:text;not null" json:"source"`
	Question                 string     `gorm:"column:question;type:text;not null" json:"question"`
	AccountID                string     `gorm:"column:account_id;type:varchar(36);not null" json:"account_id"`
	AnnotationCreatedAccount *Account   `gorm:"ForeignKey:AccountID;AssociationForeignKey:ID" json:"annotation_create_account"`
	CreatedAt                *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Score                    float64    `gorm:"column:score;type:double;not null;default:0" json:"score"`
	MessageID                string     `gorm:"column:message_id;type:varchar(36);not null" json:"message_id"`
	AnnotationQuestion       string     `gorm:"column:annotation_question;type:text;not null" json:"annotation_question"`
	AnnotationContent        string     `gorm:"column:annotation_content;type:text;not null" json:"annotation_content"`

	//     @property
	//     def account(self):
	//         account = (
	//             db.session.query(Account)
	//             .join(MessageAnnotation, MessageAnnotation.account_id == Account.id)
	//             .filter(MessageAnnotation.id == self.annotation_id)
	//             .first()
	//         )
	//         return account
}

// TableName get sql table name.获取数据库表名
func (AppAnnotationHitHistory) TableName() string {
	return "app_annotation_hit_histories"
}

// OperationLog [...]
type OperationLog struct {
	ID        string         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID  string         `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AccountID string         `gorm:"column:account_id;type:varchar(36);not null" json:"account_id"`
	Action    string         `gorm:"column:action;type:varchar(255);not null" json:"action"`
	Content   datatypes.JSON `gorm:"column:content;type:json" json:"content"`
	CreatedAt *time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedIP string         `gorm:"column:created_ip;type:varchar(255);not null" json:"created_ip"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// EndUser [...]
type EndUser struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID       string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36)" json:"app_id"`
	Type           string     `gorm:"column:type;type:varchar(255);not null" json:"type"`
	ExternalUserID string     `gorm:"column:external_user_id;type:varchar(255)" json:"external_user_id"`
	Name           string     `gorm:"column:name;type:varchar(255)" json:"name"`
	IsAnonymous    bool       `gorm:"column:is_anonymous" json:"is_anonymous"`
	SessionID      string     `gorm:"column:session_id;type:varchar(255);not null" json:"session_id"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (EndUser) TableName() string {
	return "end_users"
}

func CreateOrGetEndUserByUserID(app_model *App, user_id string) *EndUser {
	if user_id == "" {
		user_id = "DEFAULT-USER"
	}
	end_user := new(EndUser)
	err := dbengine.Instance().DB.Model(&EndUser{}).Where("tenant_id = ? and app_id = ? and session_id = ? and type =?", app_model.TenantID, app_model.ID, user_id, "service_api").First(end_user).Error
	if err != nil {
		mlog.Errorf("get end_user failed:%v", err)
		end_user = nil
	}

	if end_user == nil {
		now := time.Now()
		end_user = &EndUser{
			ID:          uuid.NewV4().String(),
			TenantID:    app_model.TenantID,
			AppID:       app_model.ID,
			Type:        "service_api",
			IsAnonymous: user_id == "DEFAULT-USER",
			SessionID:   user_id,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}
		dbengine.Instance().DB.Create(end_user)
	}
	return end_user
}

// Site [...]
type Site struct {
	ID                     string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID                  string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	Title                  string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Icon                   string     `gorm:"column:icon;type:varchar(255)" json:"icon"`
	IconBackground         string     `gorm:"column:icon_background;type:varchar(255)" json:"icon_background"`
	Description            string     `gorm:"column:description;type:text" json:"description"`
	DefaultLanguage        string     `gorm:"column:default_language;type:varchar(255);not null" json:"default_language"`
	Copyright              string     `gorm:"column:copyright;type:varchar(255)" json:"copyright"`
	PrivacyPolicy          string     `gorm:"column:privacy_policy;type:varchar(255)" json:"privacy_policy"`
	CustomizeDomain        string     `gorm:"column:customize_domain;type:varchar(255)" json:"customize_domain"`
	CustomizeTokenStrategy string     `gorm:"column:customize_token_strategy;type:varchar(255);not null" json:"customize_token_strategy"`
	PromptPublic           bool       `gorm:"column:prompt_public;type:tinyint(1);not null;default:0" json:"prompt_public"`
	Status                 string     `gorm:"column:status;type:varchar(255);default:normal" json:"status"`
	CreatedAt              *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Code                   string     `gorm:"column:code;type:varchar(255)" json:"code"`
	CustomDisclaimer       string     `gorm:"column:custom_disclaimer;type:text;not null" json:"custom_disclaimer"`
	ShowWorkflowSteps      bool       `gorm:"column:show_workflow_steps;type:tinyint(1);not null;default:1" json:"show_workflow_steps"`
	ChatColorTheme         string     `gorm:"column:chat_color_theme;type:varchar(255)" json:"chat_color_theme"`
	ChatColorThemeInverted bool       `gorm:"column:chat_color_theme_inverted;type:tinyint(1);not null;default:0" json:"chat_color_theme_inverted"`
	IconType               string     `gorm:"column:icon_type;type:varchar(255)" json:"icon_type"`
	CreatedBy              string     `gorm:"column:created_by;type:varchar(36)" json:"created_by"`
	UpdatedBy              string     `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UseIconAsAnswerIcon    bool       `gorm:"column:use_icon_as_answer_icon;type:tinyint(1);not null;default:0" json:"use_icon_as_answer_icon"`

	//     @property
	//     def custom_disclaimer(self):
	//         return self._custom_disclaimer

	//     @custom_disclaimer.setter
	//     def custom_disclaimer(self, value: str):
	//         if len(value) > 512:
	//             raise ValueError("Custom disclaimer cannot exceed 512 characters.")
	//         self._custom_disclaimer = value

}

func (s *Site) AppBaseURL() string {
	return viper.GetString("APP_WEB_URL")
}

// TableName get sql table name.获取数据库表名
func (Site) TableName() string {
	return "sites"
}
func (Site) GenerateCode(n int) string {
	maxretrytimes := 20
	retrycount := 0
	for {
		result := utils.GenerateString(n)
		var count int64
		if err := dbengine.Instance().Model(&Site{}).Where("code=?", result).Count(&count).Error; err != nil || count > 0 {
			retrycount++
			if retrycount > maxretrytimes {
				panic("gen site code exceed max retry times")
			}
			continue
		}
		return result
	}
}

// ApiToken [...]
type ApiToken struct {
	ID         string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID      string     `gorm:"column:app_id;type:varchar(36)" json:"app_id"`
	Type       string     `gorm:"column:type;type:varchar(16);not null" json:"type"`
	Token      string     `gorm:"column:token;type:varchar(255);not null" json:"token"`
	LastUsedAt *time.Time `gorm:"column:last_used_at;type:timestamp" json:"last_used_at"`
	CreatedAt  *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	TenantID   string     `gorm:"column:tenant_id;type:varchar(36)" json:"tenant_id"`
}

// TableName get sql table name.获取数据库表名
func (ApiToken) TableName() string {
	return "api_tokens"
}

func (ApiToken) GenerateApiKey(prefix string, n int) string {
	maxretrytimes := 20
	retrycount := 0
	for {
		result := prefix + utils.GenerateString(n)
		var count int64
		if err := dbengine.Instance().Model(&ApiToken{}).Where("token=?", result).Count(&count).Error; err != nil || count > 0 {
			retrycount++
			if retrycount > maxretrytimes {
				panic("gen api key exceed max retry times")
			}
			continue
		}
		return result
	}
}

// UploadFile [...]
type UploadFile struct {
	ID            string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID      string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	StorageType   string     `gorm:"column:storage_type;type:varchar(255);not null" json:"storage_type"`
	Key           string     `gorm:"column:key;type:varchar(255);not null" json:"key"`
	Name          string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Size          int        `gorm:"column:size;type:int;not null" json:"size"`
	Extension     string     `gorm:"column:extension;type:varchar(255);not null" json:"extension"`
	MimeType      string     `gorm:"column:mime_type;type:varchar(255)" json:"mime_type"`
	CreatedBy     string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt     *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Used          bool       `gorm:"column:used;type:tinyint(1);not null;default:0" json:"used"`
	UsedBy        string     `gorm:"column:used_by;type:varchar(36)" json:"used_by"`
	UsedAt        *time.Time `gorm:"column:used_at;type:timestamp" json:"used_at"`
	Hash          string     `gorm:"column:hash;type:varchar(255)" json:"hash"`
	CreatedByRole string     `gorm:"column:created_by_role;type:varchar(255);default:account" json:"created_by_role"`
	SourceURL     string     `gorm:"column:source_url;type:varchar(255);default:''" json:"source_url"`
}

// TableName get sql table name.获取数据库表名
func (UploadFile) TableName() string {
	return "upload_files"
}

// ApiRequest [...]
type ApiRequest struct {
	ID         string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID   string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	APITokenID string     `gorm:"column:api_token_id;type:varchar(36);not null" json:"api_token_id"`
	Path       string     `gorm:"column:path;type:varchar(255);not null" json:"path"`
	Request    string     `gorm:"column:request;type:text" json:"request"`
	Response   string     `gorm:"column:response;type:text" json:"response"`
	IP         string     `gorm:"column:ip;type:varchar(255);not null" json:"ip"`
	CreatedAt  *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (ApiRequest) TableName() string {
	return "api_requests"
}

// MessageChain [...]
type MessageChain struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	MessageID string     `gorm:"column:message_id;type:varchar(36);not null" json:"message_id"`
	Type      string     `gorm:"column:type;type:varchar(255);not null" json:"type"`
	Input     string     `gorm:"column:input;type:text" json:"input"`
	Output    string     `gorm:"column:output;type:text" json:"output"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (MessageChain) TableName() string {
	return "message_chains"
}

// MessageAgentThought [...]
type MessageAgentThought struct {
	ID               string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	MessageID        string     `gorm:"column:message_id;type:varchar(36);not null" json:"message_id"`
	MessageChainID   string     `gorm:"column:message_chain_id;type:varchar(36)" json:"message_chain_id"`
	Position         int        `gorm:"column:position;type:int;not null" json:"position"`
	Thought          string     `gorm:"column:thought;type:text" json:"thought"`
	Tool             string     `gorm:"column:tool;type:text" json:"tool"`
	ToolInput        string     `gorm:"column:tool_input;type:text" json:"tool_input"`
	Observation      string     `gorm:"column:observation;type:text" json:"observation"`
	ToolProcessData  string     `gorm:"column:tool_process_data;type:text" json:"tool_process_data"`
	Message          string     `gorm:"column:message;type:text" json:"message"`
	MessageToken     int        `gorm:"column:message_token;type:int" json:"message_token"`
	MessageUnitPrice float64    `gorm:"column:message_unit_price;type:decimal(10,0)" json:"message_unit_price"`
	Answer           string     `gorm:"column:answer;type:text" json:"answer"`
	AnswerToken      int        `gorm:"column:answer_token;type:int" json:"answer_token"`
	AnswerUnitPrice  float64    `gorm:"column:answer_unit_price;type:decimal(10,0)" json:"answer_unit_price"`
	Tokens           int        `gorm:"column:tokens;type:int" json:"tokens"`
	TotalPrice       float64    `gorm:"column:total_price;type:decimal(10,0)" json:"total_price"`
	Currency         string     `gorm:"column:currency;type:varchar(255)" json:"currency"`
	Latency          float64    `gorm:"column:latency;type:double" json:"latency"`
	CreatedByRole    string     `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy        string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	MessagePriceUnit float64    `gorm:"column:message_price_unit;type:decimal(10,7);not null;default:0.0010000" json:"message_price_unit"`
	AnswerPriceUnit  float64    `gorm:"column:answer_price_unit;type:decimal(10,7);not null;default:0.0010000" json:"answer_price_unit"`
	MessageFiles     string     `gorm:"column:message_files;type:text" json:"message_files"`
	ToolLabelsStr    string     `gorm:"column:tool_labels_str;type:varchar(255);default:{}" json:"tool_labels_str"`
	ToolMetaStr      string     `gorm:"column:tool_meta_str;type:varchar(255);default:{}" json:"tool_meta_str"`

	//     @property
	//     def tools(self) -> list[str]:
	//         return self.tool.split(";") if self.tool else []

	//     @property
	//     def tool_meta(self) -> dict:
	//         try:
	//             if self.tool_meta_str:
	//                 return cast(dict, json.loads(self.tool_meta_str))
	//             else:
	//                 return {}
	//         except Exception as e:
	//             return {}

	//     @property
	//     def tool_inputs_dict(self) -> dict:
	//         tools = self.tools
	//         try:
	//             if self.tool_input:
	//                 data = json.loads(self.tool_input)
	//                 result = {}
	//                 for tool in tools:
	//                     if tool in data:
	//                         result[tool] = data[tool]
	//                     else:
	//                         if len(tools) == 1:
	//                             result[tool] = data
	//                         else:
	//                             result[tool] = {}
	//                 return result
	//             else:
	//                 return {tool: {} for tool in tools}
	//         except Exception as e:
	//             return {}

	//     @property
	//     def tool_outputs_dict(self) -> dict:
	//         tools = self.tools
	//         try:
	//             if self.observation:
	//                 data = json.loads(self.observation)
	//                 result = {}
	//                 for tool in tools:
	//                     if tool in data:
	//                         result[tool] = data[tool]
	//                     else:
	//                         if len(tools) == 1:
	//                             result[tool] = data
	//                         else:
	//                             result[tool] = {}
	//                 return result
	//             else:
	//                 return {tool: {} for tool in tools}
	//         except Exception as e:
	//             if self.observation:
	//                 return dict.fromkeys(tools, self.observation)
	//             else:
	//                 return {}
}

func (mat *MessageAgentThought) ToolLabels() map[string]any {
	if mat.ToolLabelsStr != "" {
		var dict map[string]any
		err := json.Unmarshal([]byte(mat.ToolLabelsStr), &dict)
		if err != nil {
			mlog.Errorf("json marshal(%s) failed:%v", mat.ToolLabelsStr, err)
			return map[string]any{}
		}
		return dict
	}
	return map[string]any{}
}

func (mat *MessageAgentThought) Files() []any {
	if mat.MessageFiles != "" {
		var tmplist []any
		err := json.Unmarshal([]byte(mat.MessageFiles), &tmplist)
		if err != nil {
			mlog.Errorf("json marshal(%s) failed:%v", mat.MessageFiles, err)
			return []any{}
		}
		return tmplist
	} else {
		return []any{}
	}
}

// TableName get sql table name.获取数据库表名
func (MessageAgentThought) TableName() string {
	return "message_agent_thoughts"
}

// DatasetRetrieverResource [...]
type DatasetRetrieverResource struct {
	ID              string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	MessageID       string     `gorm:"column:message_id;type:varchar(36);not null" json:"message_id"`
	Position        int        `gorm:"column:position;type:int;not null" json:"position"`
	DatasetID       string     `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	DatasetName     string     `gorm:"column:dataset_name;type:text;not null" json:"dataset_name"`
	DocumentID      string     `gorm:"column:document_id;type:varchar(36)" json:"document_id"`
	DocumentName    string     `gorm:"column:document_name;type:text;not null" json:"document_name"`
	DataSourceType  string     `gorm:"column:data_source_type;type:text" json:"data_source_type"`
	SegmentID       string     `gorm:"column:segment_id;type:varchar(36)" json:"segment_id"`
	Score           float64    `gorm:"column:score;type:double" json:"score"`
	Content         string     `gorm:"column:content;type:text;not null" json:"content"`
	HitCount        int        `gorm:"column:hit_count;type:int" json:"hit_count"`
	WordCount       int        `gorm:"column:word_count;type:int" json:"word_count"`
	SegmentPosition int        `gorm:"column:segment_position;type:int" json:"segment_position"`
	IndexNodeHash   string     `gorm:"column:index_node_hash;type:text" json:"index_node_hash"`
	RetrieverFrom   string     `gorm:"column:retriever_from;type:text;not null" json:"retriever_from"`
	CreatedBy       string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (DatasetRetrieverResource) TableName() string {
	return "dataset_retriever_resources"
}

// Tag [...]
type Tag struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID  string     `gorm:"column:tenant_id;type:varchar(36)" json:"tenant_id"`
	Type      string     `gorm:"column:type;type:varchar(16);not null" json:"type"`
	Name      string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	CreatedBy string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (Tag) TableName() string {
	return "tags"
}

var TAG_TYPE_LIST = []string{"knowledge", "app"}

// TagBinding [...]
type TagBinding struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID  string     `gorm:"column:tenant_id;type:varchar(36)" json:"tenant_id"`
	Tenant    *Tenant    `gorm:"ForeignKey:TenantID;AssociationForeignKey:ID" json:"tenant"`
	TagID     string     `gorm:"column:tag_id;type:varchar(36)" json:"tag_id"`
	Tag       *Tag       `gorm:"ForeignKey:TagID;AssociationForeignKey:ID" json:"tag"`
	TargetID  string     `gorm:"column:target_id;type:varchar(36)" json:"target_id"`
	CreatedBy string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (TagBinding) TableName() string {
	return "tag_bindings"
}

// TraceAppConfig [...]
type TraceAppConfig struct {
	ID              string         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID           string         `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	TracingProvider string         `gorm:"column:tracing_provider;type:varchar(255)" json:"tracing_provider"`
	TracingConfig   datatypes.JSON `gorm:"column:tracing_config;type:json" json:"tracing_config"`
	CreatedAt       *time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsActive        bool           `gorm:"column:is_active;type:tinyint(1);not null;default:1" json:"is_active"`

	//     @property
	//     def tracing_config_dict(self):
	//         return self.tracing_config or {}

	//     @property
	//     def tracing_config_str(self):
	//         return json.dumps(self.tracing_config_dict)

	//     def to_dict(self):
	//         return {
	//             "id": self.id,
	//             "app_id": self.app_id,
	//             "tracing_provider": self.tracing_provider,
	//             "tracing_config": self.tracing_config_dict,
	//             "is_active": self.is_active,
	//             "created_at": str(self.created_at) if self.created_at else None,
	//             "updated_at": str(self.updated_at) if self.updated_at else None,
	//         }
}

// TableName get sql table name.获取数据库表名
func (TraceAppConfig) TableName() string {
	return "trace_app_config"
}
