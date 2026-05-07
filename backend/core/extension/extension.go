package extension

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/mlog"
)

type Extension struct {
	ModuleExtensions map[ExtensionModuleType]map[string]*ModuleExtension
}

func (ext *Extension) ModuleExtension(module ExtensionModuleType, extension_name string) *ModuleExtension {
	var module_extensions map[string]*ModuleExtension
	if _, ok := ext.ModuleExtensions[module]; ok {
		module_extensions = ext.ModuleExtensions[module]
	}
	if module_extensions == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Extension Module {%v} not found", module)))
	}
	var module_extension *ModuleExtension
	if _, ok := module_extensions[extension_name]; ok {
		module_extension = module_extensions[extension_name]
	}
	if module_extension == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Extension {%s} not found", extension_name)))
	}
	return module_extension
}
func (ext *Extension) ExtensionClass(module ExtensionModuleType, extension_name string) Extensiblor {
	module_extension := ext.ModuleExtension(module, extension_name)
	return module_extension.ExtensionClass
}
func (ext *Extension) ValidateFormSchema(module ExtensionModuleType, extension_name string, config map[string]any) {
	module_extension := ext.ModuleExtension(module, extension_name)
	form_schema := module_extension.FormSchema

	// TODO validate form_schema
	mlog.Infof("------form_schema=%#v", form_schema)
}
