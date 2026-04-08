package codeexecutor

import codeexecutorenumtypes "mlib.com/gofy/server/enum_types/code_executor"

type CodeNodeProvider interface {
	GetLanguage() codeexecutorenumtypes.CodeLanguage
	IsAcceptLanguage(codeexecutorenumtypes.CodeLanguage) bool
	GetDefaultCode() string
	GetDefaultConfig() map[string]any
}

// @classmethod
// def is_accept_language(cls, language: str) -> bool:
//     return language == cls.get_language()

// @classmethod
// def get_default_config(cls) -> dict:
//     return {
//         "type": "code",
//         "config": {
//             "variables": [{"variable": "arg1", "value_selector": []}, {"variable": "arg2", "value_selector": []}],
//             "code_language": cls.get_language(),
//             "code": cls.get_default_code(),
//             "outputs": {"result": {"type": "string", "children": None}},
//         },
//     }
