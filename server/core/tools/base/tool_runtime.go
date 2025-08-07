package base

import (
	"encoding/json"
	"maps"

	"mlib.com/gofy/server/core/exceptions"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	"mlib.com/mlog"
)

type ToolRuntime struct {
	TenantID          string                            `json:"tenant_id"`
	ToolID            string                            `json:"tool_id"`
	InvokeFrom        appenumtypes.InvokeFrom           `json:"invoke_from"`
	ToolInvokeFrom    toolsenumtypes.ToolInvokeFromType `json:"tool_invoke_from"`
	Credentials       map[string]any                    `json:"credentials"`
	CredentialType    toolsenumtypes.CredentialType     `json:"credential_type"` //default=CredentialType.API_KEY)
	RuntimeParameters map[string]any                    `json:"runtime_parameters"`
}

func NewToolRuntime(data any) *ToolRuntime {
	tool_runtime := new(ToolRuntime)
	if data == nil {
		return tool_runtime
	}
	if _, ok := data.(map[string]any); ok {
		bindata, _ := json.Marshal(data.(map[string]any))
		err := json.Unmarshal(bindata, tool_runtime)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("new ToolRuntime failed, because of provaded invalid initial data"))
		}
		return tool_runtime
	} else if _, ok := data.(*ToolRuntime); ok {
		tool_runtime.TenantID = data.(*ToolRuntime).TenantID
		tool_runtime.ToolID = data.(*ToolRuntime).ToolID
		tool_runtime.InvokeFrom = data.(*ToolRuntime).InvokeFrom
		tool_runtime.ToolInvokeFrom = data.(*ToolRuntime).ToolInvokeFrom
		tool_runtime.Credentials = maps.Clone(data.(*ToolRuntime).Credentials)
		tool_runtime.RuntimeParameters = maps.Clone(data.(*ToolRuntime).RuntimeParameters)
		return tool_runtime
	} else {
		mlog.Warningf("provaded initial data=%#v is unsupported", data)
		tool_runtime.CredentialType = toolsenumtypes.Credential_API_KEY
		return tool_runtime
	}
}
