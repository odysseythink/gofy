package builtintool

import (
	"mlib.com/gofy/server/core/tools/base"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
)

type BuiltinToolProviderController struct {
	*base.BaseToolProviderController
	tools []*BuiltinTool
}

func NewBuiltinToolProviderController() {

}

func (controller *BuiltinToolProviderController) ProviderType() toolsenumtypes.ToolProviderType {
	return toolsenumtypes.ToolProvider_BUILT_IN
}
