package services

type OpsService struct {
}

func (service *OpsService) GetTraceAppConfig(app_id string, tracing_provider string) map[string]any {
	return nil
	// 	trace_config_data:= new(models.TraceAppConfig)
	// 	err := dbengine.Instance().DB.Model(&models.TraceAppConfig{}).Where("app_id = ? and tracing_provider = ?",
	// 	app_id,
	// 	tracing_provider,
	// 	).First(trace_config_data).Error
	// 	if err != nil {
	// 		mlog.Errorf("get trace config failed:%v", err)

	// 	}

	// 	if trace_config_data == nil{
	// 		return nil
	// }
	// 	// decrypt_token and obfuscated_token
	// 	tenant = db.session.query(App).filter(App.id == app_id).first()
	// 	if not tenant:
	// 		return None
	// 	tenant_id = tenant.tenant_id
	// 	decrypt_tracing_config = OpsTraceManager.decrypt_tracing_config(
	// 		tenant_id, tracing_provider, trace_config_data.tracing_config
	// 	)
	// 	new_decrypt_tracing_config = OpsTraceManager.obfuscated_decrypt_token(tracing_provider, decrypt_tracing_config)

	// 	if tracing_provider == "langfuse" and (
	// 		"project_key" not in decrypt_tracing_config or not decrypt_tracing_config.get("project_key")
	// 	):
	// 		try:
	// 			project_key = OpsTraceManager.get_trace_config_project_key(decrypt_tracing_config, tracing_provider)
	// 			new_decrypt_tracing_config.update(
	// 				{
	// 					"project_url": "{host}/project/{key}".format(
	// 						host=decrypt_tracing_config.get("host"), key=project_key
	// 					)
	// 				}
	// 			)
	// 		except Exception:
	// 			new_decrypt_tracing_config.update(
	// 				{"project_url": "{host}/".format(host=decrypt_tracing_config.get("host"))}
	// 			)

	// 	if tracing_provider == "langsmith" and (
	// 		"project_url" not in decrypt_tracing_config or not decrypt_tracing_config.get("project_url")
	// 	):
	// 		try:
	// 			project_url = OpsTraceManager.get_trace_config_project_url(decrypt_tracing_config, tracing_provider)
	// 			new_decrypt_tracing_config.update({"project_url": project_url})
	// 		except Exception:
	// 			new_decrypt_tracing_config.update({"project_url": "https://smith.langchain.com/"})

	// 	if tracing_provider == "opik" and (
	// 		"project_url" not in decrypt_tracing_config or not decrypt_tracing_config.get("project_url")
	// 	):
	// 		try:
	// 			project_url = OpsTraceManager.get_trace_config_project_url(decrypt_tracing_config, tracing_provider)
	// 			new_decrypt_tracing_config.update({"project_url": project_url})
	// 		except Exception:
	// 			new_decrypt_tracing_config.update({"project_url": "https://www.comet.com/opik/"})

	// trace_config_data.tracing_config = new_decrypt_tracing_config
	// return trace_config_data.to_dict()
}
