package llmgenerator

import (
	"fmt"
	"reflect"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/utils"
	"mlib.com/gofy/server/utils/validate"
)

type RuleConfigGeneratorOutputParser struct{}

func (parser *RuleConfigGeneratorOutputParser) GetFormatInstructions() (string, string, string) {
	return RULE_CONFIG_PROMPT_GENERATE_TEMPLATE,
		RULE_CONFIG_PARAMETER_GENERATE_TEMPLATE,
		RULE_CONFIG_STATEMENT_GENERATE_TEMPLATE
}

var (
	expected_keys = []string{"prompt", "variables", "opening_statement"}
)

func (parser *RuleConfigGeneratorOutputParser) Parse(text string) map[string]any {
	return func() map[string]any {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					panic(exceptions.NewOutputParserError(fmt.Sprintf("Parsing text\n{%s}\n of rule config generator raised following error:\n{%v}", text, err)))
				} else {
					panic(r)
				}
			}
		}()
		parsed := utils.ParseAndCheckJsonMarkdown(text, expected_keys)
		err := validate.StringMapTypeVerify(parsed, validate.Rules{
			"prompt": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("Expected 'prompt' to be a string.:%v", err)
			panic(exceptions.NewValueError("Expected 'prompt' to be a string.:" + err.Error()))
		}
		err = validate.StringMapTypeVerify(parsed, validate.Rules{
			"variables": {validate.RuleTypeOfField(reflect.Slice), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("Expected 'variables' to be a list.:%v", err)
			panic(exceptions.NewValueError("Expected 'variables' to be a list.:" + err.Error()))
		}
		err = validate.StringMapTypeVerify(parsed, validate.Rules{
			"opening_statement": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("Expected 'opening_statement' to be a str.:%v", err)
			panic(exceptions.NewValueError("Expected 'opening_statement' to be a str.:" + err.Error()))
		}
		return parsed
	}()

}
