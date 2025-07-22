package extension

type ExtensionModuleType string

const (
	ExtensionModule_MODERATION         ExtensionModuleType = "moderation"
	ExtensionModule_EXTERNAL_DATA_TOOL ExtensionModuleType = "external_data_tool"
)

type ModuleExtension struct {
	ExtensionClass any              `json:"extension_class,omitempty"`
	Name           string           `json:"name"`
	Label          map[string]any   `json:"label,omitempty"`
	FormSchema     []map[string]any `json:"form_schema,omitempty"`
	Builtin        bool             `json:"builtin"`
	Position       int              `json:"position,omitempty"`
}

func NewModuleExtension() *ModuleExtension {
	return &ModuleExtension{
		Builtin: true,
	}
}
