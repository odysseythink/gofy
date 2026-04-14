package llmgenerator

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/odysseythink/mlog"
)

type SuggestedQuestionsAfterAnswerOutputParser struct {
}

func (parser *SuggestedQuestionsAfterAnswerOutputParser) GetFormatInstructions() string {
	return SUGGESTED_QUESTIONS_AFTER_ANSWER_INSTRUCTION_PROMPT
}
func (parser *SuggestedQuestionsAfterAnswerOutputParser) Parse(text string) []string {
	text = strings.TrimPrefix(text, " ")
	text = strings.TrimSuffix(text, " ")
	re := regexp.MustCompile(`\[.*?\]`)
	match := re.FindString(text)
	if match != "" {
		// 解析 JSON 字符串
		jsonObj := []string{}
		if err := json.Unmarshal([]byte(match), &jsonObj); err != nil {
			mlog.Error("Error parsing JSON:", err)
			return []string{}
		}
		return jsonObj
	}
	return []string{}
}
