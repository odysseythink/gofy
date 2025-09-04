package extension

type ExtensionModuleType string

const (
	ExtensionModule_MODERATION         ExtensionModuleType = "moderation"
	ExtensionModule_EXTERNAL_DATA_TOOL ExtensionModuleType = "external_data_tool"
)

type ModuleExtension struct {
	ExtensionClass Extensiblor      `json:"extension_class,omitempty"`
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

type Extensiblor interface {
	Module() ExtensionModuleType
	Name() string
	SetTenantID(string)
	GetTenantID() string
	SetConfig(map[string]any)
	GetConfig() map[string]any
}
type Extensible struct {
	// Module ExtensionModuleType `json:"module"`

	// Name     string         `json:"name"`
	TenantID string         `json:"tenant_id"`
	Config   map[string]any `json:"config"`
}

func NewExtensible(tenant_id string, config map[string]any) *Extensible {
	return &Extensible{
		TenantID: tenant_id,
		Config:   config,
	}
}
func (ext *Extensible) SetTenantID(val string) {
	ext.TenantID = val
}
func (ext *Extensible) GetTenantID() string {
	return ext.TenantID
}
func (ext *Extensible) SetConfig(val map[string]any) {
	ext.Config = val
}
func (ext *Extensible) GetConfig() map[string]any {
	return ext.Config
}
