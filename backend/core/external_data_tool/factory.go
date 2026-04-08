package externaldatatool

import (
	"fmt"
	"reflect"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/extension"
	"mlib.com/gofy/server/core/external_data_tool/base"
	"mlib.com/gofy/server/global"
	"mlib.com/mlog"
)

type ExternalDataToolFactory struct {
	__extension_instance base.ExternalDataToolor
}

func NewExternalDataToolFactory(
	name, tenantID, appID, variable string,
	config map[string]interface{},
) (*ExternalDataToolFactory, error) {
	extension_class := global.CodeBasedExtension.ExtensionClass(extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)
	if extension_class == nil {
		mlog.Errorf("extension module(%s) name(%s) not exist", extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)
		return nil, fmt.Errorf("extension module(%s) name(%s) not exist", extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)
	}
	if _, ok := any(extension_class).(base.ExternalDataToolor); !ok {
		mlog.Errorf("extension module(%s) name(%s) is not a ExternalDataToolor", extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)
		return nil, fmt.Errorf("extension module(%s) name(%s) is not a ExternalDataToolor", extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)
	}

	typ := reflect.TypeOf(extension_class)
	ptrValue := reflect.New(typ)

	return &ExternalDataToolFactory{
		__extension_instance: ptrValue.Interface().(base.ExternalDataToolor),
	}, nil
}

func (factory *ExternalDataToolFactory) ValidateConfig(name string, tenant_id string, config map[string]any) {
	/*
	   Validate the incoming form config data.

	   :param name: the name of external data tool
	   :param tenant_id: the id of workspace
	   :param config: the form config data
	   :return:
	*/
	global.CodeBasedExtension.ValidateFormSchema(extension.ExtensionModule_EXTERNAL_DATA_TOOL, name, config)
	extension_class := global.CodeBasedExtension.ExtensionClass(extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)
	if _, ok := any(extension_class).(base.ExternalDataToolor); !ok {
		mlog.Errorf("extension module(%s) name(%s) is not a ExternalDataToolor", extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)
		panic(exceptions.NewValueError(fmt.Sprintf("extension module(%s) name(%s) is not a ExternalDataToolor", extension.ExtensionModule_EXTERNAL_DATA_TOOL, name)))
	}
	// FIXME mypy issue here, figure out how to fix it
	any(extension_class).(base.ExternalDataToolor).ValidateConfig(tenant_id, config) // type: ignore
}

func (factory *ExternalDataToolFactory) Query(inputs map[string]any, query string) string {
	/*
	   Query the external data tool.

	   :param inputs: user inputs
	   :param query: the query of chat app
	   :return: the tool query result
	*/
	return factory.__extension_instance.Query(inputs, query)
}
