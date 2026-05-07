package base

import "github.com/odysseythink/gofy/backend/core/extension"

type ExternalDataToolor interface {
	extension.Extensiblor
	ValidateConfig(tenant_id string, config map[string]any)
	Query(inputs map[string]any, query string) string
	GetAppID() string
	SetAppID(string)
	GetVariable() string
	SetVariable(string)
}
type ExternalDataTool struct {
	*extension.Extensible
	AppID    string `json:"app_id"`   //the id of app
	Variable string `json:"variable"` // the tool variable name of app tool
}

func (extdt *ExternalDataTool) Module() extension.ExtensionModuleType {
	return extension.ExtensionModule_EXTERNAL_DATA_TOOL
}
func (extdt *ExternalDataTool) GetAppID() string {
	return extdt.AppID
}
func (extdt *ExternalDataTool) SetAppID(val string) {
	extdt.AppID = val
}
func (extdt *ExternalDataTool) GetVariable() string {
	return extdt.Variable
}
func (extdt *ExternalDataTool) SetVariable(val string) {
	extdt.Variable = val
}
func NewExternalDataTool(tenant_id string, app_id string, variable string, config map[string]any) *ExternalDataTool {
	return &ExternalDataTool{
		Extensible: extension.NewExtensible(tenant_id, config),
		AppID:      app_id,
		Variable:   variable,
	}
}
