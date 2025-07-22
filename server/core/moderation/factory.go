package moderation

type ModerationFactory struct {
}

func (f *ModerationFactory) ValidateConfig(name string, tenant_id string, config map[string]any) {
	/*
	   Validate the incoming form config data.

	   :param name: the name of extension
	   :param tenant_id: the id of workspace
	   :param config: the form config data
	   :return:
	*/
	// code_based_extension.validate_form_schema(ExtensionModule.MODERATION, name, config)
	// extension_class = code_based_extension.extension_class(ExtensionModule.MODERATION, name)
	// FIXME: mypy error, try to fix it instead of using type: ignore
	// extension_class.validate_config(tenant_id, config) // type: ignore
}
