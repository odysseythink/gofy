package python3

import codeexecutorenumtypes "github.com/odysseythink/gofy/backend/enum_types/code_executor"

type Python3CodeProvider struct{}

func (provider *Python3CodeProvider) GetLanguage() codeexecutorenumtypes.CodeLanguage {
	return codeexecutorenumtypes.CodeLanguage_PYTHON3
}
func (provider *Python3CodeProvider) IsAcceptLanguage(language codeexecutorenumtypes.CodeLanguage) bool {
	return language == provider.GetLanguage()
}
func (provider *Python3CodeProvider) GetDefaultCode() string {
	return `
def main(arg1: str, arg2: str) -> dict:
	return {
		"result": arg1 + arg2,
	}	
	`
}
func (provider *Python3CodeProvider) GetDefaultConfig() map[string]any {
	return map[string]any{
		"type": "code",
		"config": map[string]any{
			"variables": []any{
				map[string]any{"variable": "arg1", "value_selector": []any{}},
				map[string]any{"variable": "arg2", "value_selector": []any{}},
			},
			"code_language": provider.GetLanguage(),
			"code":          provider.GetDefaultCode(),
			"outputs":       map[string]any{"result": map[string]any{"type": "string", "children": nil}},
		},
	}
}
