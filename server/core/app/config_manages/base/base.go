package base

import (
	retrievalresource "mlib.com/gofy/server/core/app/config_manages/features/retrieval_resource"
	suggestedquestionsafteranswer "mlib.com/gofy/server/core/app/config_manages/features/suggested_questions_after_answer"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	"mlib.com/gofy/server/models"
)

type IManager interface {
	Convert(config map[string]any)
}

// 定义FileExtraConfig结构体
type FileExtraConfig struct {
	ImageConfig map[string]any `json:"image_config,omitempty"`
}

// 定义TracingConfigEntity结构体
type TracingConfigEntity struct {
	Enabled         bool   `json:"enabled"`
	TracingProvider string `json:"tracing_provider"`
}

// 定义BaseAppConfigManager结构体
type BaseAppConfigManage struct {
}

// ConvertFeatures 方法
func (m *BaseAppConfigManage) ConvertFeatures(configDict map[string]any, appMode models.AppMode) *appconfigentities.AppAdditionalFeatures {
	additionalFeatures := &appconfigentities.AppAdditionalFeatures{}

	showRetrieveSource := (&retrievalresource.RetrievalResourceConfigManager{}).Convert(configDict)
	additionalFeatures.ShowRetrieveSource = showRetrieveSource

	// isVision := appMode == enumtypes.AppModeCHAT || appMode == enumtypes.AppModeCOMPLETION || appMode == enumtypes.AppModeAGENT_CHAT
	// fileUpload := FileUploadConfigManager{}.Convert(configDict, isVision)
	// additionalFeatures.FileUpload = fileUpload

	additionalFeatures.SuggestedQuestionsAfterAnswer = (&suggestedquestionsafteranswer.SuggestedQuestionsAfterAnswerConfigManager{}).Convert(configDict)

	return additionalFeatures
}

// func main() {
// 	// 示例使用
// 	config := map[string]any{
// 		"file_upload": map[string]any{
// 			"image": map[string]any{
// 				"enabled":          true,
// 				"number_limits":    5,
// 				"transfer_methods": []any{"remote_url", "local_file"},
// 				"detail":           "high",
// 			},
// 		},
// 		"retriever_resource": map[string]any{
// 			"enabled": true,
// 		},
// 	}
// 	appMode := AppModeCHAT

// 	manager := BaseAppConfigManage{}
// 	features := manager.ConvertFeatures(config, appMode)
// 	fmt.Printf("%+v\n", features)
// }
