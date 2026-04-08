package app

import "mlib.com/gofy/server/core/exceptions"

type InvokeFrom string

const (
	InvokeFrom_SERVICE_API InvokeFrom = "service-api"
	InvokeFrom_WEB_APP     InvokeFrom = "web-app"
	InvokeFrom_EXPLORE     InvokeFrom = "explore"
	InvokeFrom_DEBUGGER    InvokeFrom = "debugger"
)

func (ivk InvokeFrom) Validate() {
	if ivk != InvokeFrom_SERVICE_API &&
		ivk != InvokeFrom_WEB_APP &&
		ivk != InvokeFrom_EXPLORE &&
		ivk != InvokeFrom_DEBUGGER {
		panic(exceptions.NewValueError("invalid InvokeFrom"))
	}
}
func (ivk InvokeFrom) ToSource() string {
	/*
		Get source of invoke from.

		:return: source
	*/
	if ivk == InvokeFrom_WEB_APP {
		return "web_app"
	} else if ivk == InvokeFrom_DEBUGGER {
		return "dev"
	} else if ivk == InvokeFrom_EXPLORE {
		return "explore_app"
	} else if ivk == InvokeFrom_SERVICE_API {
		return "api"
	}
	return "dev"
}
