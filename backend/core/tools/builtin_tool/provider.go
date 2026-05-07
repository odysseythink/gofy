package builtintool

import (
	"github.com/odysseythink/gofy/backend/core/tools/base"
	toolsenumtypes "github.com/odysseythink/gofy/backend/enum_types/tools"
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
