package suggestedquestionsafteranswer

import (
	"errors"
)

// SuggestedQuestionsAfterAnswerConfigManager 管理器结构体
type SuggestedQuestionsAfterAnswerConfigManager struct{}

// Convert 将配置转换为模型配置
func (m *SuggestedQuestionsAfterAnswerConfigManager) Convert(config map[string]any) bool {
	suggestedQuestionsAfterAnswer := false
	suggestedQuestionsAfterAnswerDict, ok := config["suggested_questions_after_answer"].(map[string]any)
	if ok {
		if enabled, ok := suggestedQuestionsAfterAnswerDict["enabled"].(bool); ok {
			suggestedQuestionsAfterAnswer = enabled
		}
	}

	return suggestedQuestionsAfterAnswer
}

// ValidateAndSetDefaults 验证并设置建议问题功能的默认值
func (m *SuggestedQuestionsAfterAnswerConfigManager) ValidateAndSetDefaults(config map[string]any) (map[string]any, []string, error) {
	if config["suggested_questions_after_answer"] == nil {
		config["suggested_questions_after_answer"] = map[string]any{"enabled": false}
	}

	suggestedQuestionsAfterAnswerDict, ok := config["suggested_questions_after_answer"].(map[string]any)
	if !ok {
		return nil, nil, errors.New("suggested_questions_after_answer must be of dict type")
	}

	if suggestedQuestionsAfterAnswerDict["enabled"] == nil {
		suggestedQuestionsAfterAnswerDict["enabled"] = false
	}

	if _, ok := suggestedQuestionsAfterAnswerDict["enabled"].(bool); !ok {
		return nil, nil, errors.New("enabled in suggested_questions_after_answer must be of boolean type")
	}

	return config, []string{"suggested_questions_after_answer"}, nil
}

// func main() {
// 	// 示例使用
// 	config := map[string]any{
// 		"suggested_questions_after_answer": map[string]any{"enabled": true},
// 	}

// 	manager := SuggestedQuestionsAfterAnswerConfigManager{}
// 	enabled := manager.Convert(config)
// 	config, fields, err := manager.ValidateAndSetDefaults(config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	println("Suggested Questions After Answer Enabled:", enabled)
// 	println("Config:", config)
// 	println("Fields:", fields)
// }
