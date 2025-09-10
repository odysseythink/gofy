package toolmanager

import (
	"encoding/json"
	"fmt"
	"sync"

	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/tools/base"
	builtintool "mlib.com/gofy/server/core/tools/builtin_tool"
	dbengine "mlib.com/gofy/server/db_engine"
	agententities "mlib.com/gofy/server/entities/agent"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	"mlib.com/gofy/server/models"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/mlog"
)

type ToolManager struct {
	builtin_provider_lock    sync.Mutex
	builtin_providers        map[string]*builtintool.BuiltinToolProviderController
	builtin_providers_loaded bool
	builtin_tools_labels     map[string]*commontypes.I18nObject
}

func NewToolManager() *ToolManager {
	return &ToolManager{
		builtin_providers:    make(map[string]*builtintool.BuiltinToolProviderController),
		builtin_tools_labels: make(map[string]*commontypes.I18nObject),
	}
}

func (tm *ToolManager) GetBuiltinProvider(provider string) (*builtintool.BuiltinToolProviderController, error) {
	/*
		get the builtin provider

		:param provider: the name of the provider
		:return: the provider
	*/
	// if len(tm.builtin_providers) == 0 {
	// 	// init the builtin providers
	// 	tm.load_builtin_providers_cache()
	// }
	if _, ok := tm.builtin_providers[provider]; !ok {
		return nil, exceptions.NewToolProviderNotFoundError(fmt.Sprintf("builtin provider %s not found", provider))
	}

	return tm.builtin_providers[provider], nil
}

// func(tm *ToolManager) GetBuiltinTool(provider, tool_name string) -> Union[BuiltinTool, Tool]{
// 	/*
// 	get the builtin tool

// 	:param provider: the name of the provider
// 	:param tool_name: the name of the tool

// 	:return: the provider, the tool
// 	*/
// 	provider_controller := tm.GetBuiltinProvider(provider)
// 	tool = provider_controller.get_tool(tool_name)
// 	if tool is None:
// 		raise ToolNotFoundError(f"tool {tool_name} not found")

//		return tool
//	}
func (tm *ToolManager) GetToolRuntime(
	provider_type string,
	provider_id string,
	tool_name string,
	tenant_id string,
	invoke_from appenumtypes.InvokeFrom, /* = InvokeFrom.DEBUGGER*/
	tool_invoke_from toolsenumtypes.ToolInvokeFromType, /* = ToolInvokeFrom.AGENT*/
) /*-> Union[BuiltinTool, ApiTool, Tool]*/ (any, error) {
	/*
		get the tool runtime

		:param provider_type: the type of the provider
		:param provider_name: the name of the provider
		:param tool_name: the name of the tool

		:return: the tool
	*/
	return nil, nil
	// controller: Union[BuiltinToolProviderController, ApiToolProviderController, WorkflowToolProviderController]
	// if provider_type == "builtin"{
	// 	builtin_tool = tm.get_builtin_tool(provider_id, tool_name)

	// 	// check if the builtin tool need credentials
	// 	provider_controller = tm.get_builtin_provider(provider_id)
	// 	if not provider_controller.need_credentials{
	// 		return builtin_tool.fork_tool_runtime(
	// 			runtime={
	// 				"tenant_id": tenant_id,
	// 				"credentials": {},
	// 				"invoke_from": invoke_from,
	// 				"tool_invoke_from": tool_invoke_from,
	// 			}
	// 		)
	// 	}
	// 	// get credentials
	// 	builtin_provider: Optional[BuiltinToolProvider] = (
	// 		db.session.query(BuiltinToolProvider)
	// 		.filter(
	// 			BuiltinToolProvider.tenant_id == tenant_id,
	// 			BuiltinToolProvider.provider == provider_id,
	// 		)
	// 		.first()
	// 	)

	// 	if builtin_provider is None{
	// 		raise ToolProviderNotFoundError(f"builtin provider {provider_id} not found")
	// 	}
	// 	// decrypt the credentials
	// 	credentials = builtin_provider.credentials
	// 	controller = tm.get_builtin_provider(provider_id)
	// 	tool_configuration = ToolConfigurationManager(tenant_id=tenant_id, provider_controller=controller)

	// 	decrypted_credentials = tool_configuration.decrypt_tool_credentials(credentials)

	// 	return builtin_tool.fork_tool_runtime(
	// 		runtime={
	// 			"tenant_id": tenant_id,
	// 			"credentials": decrypted_credentials,
	// 			"runtime_parameters": {},
	// 			"invoke_from": invoke_from,
	// 			"tool_invoke_from": tool_invoke_from,
	// 		}
	// 	)

	// }else if provider_type == "api"{
	// 	if tenant_id is None{
	// 		raise ValueError("tenant id is required for api provider")
	// 	}
	// 	api_provider, credentials = tm.get_api_provider_controller(tenant_id, provider_id)

	// 	// decrypt the credentials
	// 	tool_configuration = ToolConfigurationManager(tenant_id=tenant_id, provider_controller=api_provider)
	// 	decrypted_credentials = tool_configuration.decrypt_tool_credentials(credentials)

	// 	return api_provider.get_tool(tool_name).fork_tool_runtime(
	// 		runtime={
	// 			"tenant_id": tenant_id,
	// 			"credentials": decrypted_credentials,
	// 			"invoke_from": invoke_from,
	// 			"tool_invoke_from": tool_invoke_from,
	// 		}
	// 	)
	// }else if provider_type == "workflow"{
	// 	workflow_provider: Optional[WorkflowToolProvider] = (
	// 		db.session.query(WorkflowToolProvider)
	// 		.filter(WorkflowToolProvider.tenant_id == tenant_id, WorkflowToolProvider.id == provider_id)
	// 		.first()
	// 	)

	// 	if workflow_provider is None{
	// 		raise ToolProviderNotFoundError(f"workflow provider {provider_id} not found")
	// 	}
	// 	controller = ToolTransformService.workflow_provider_to_controller(db_provider=workflow_provider)
	// 	controller_tools: Optional[list[Tool]] = controller.get_tools(
	// 		user_id="", tenant_id=workflow_provider.tenant_id
	// 	)
	// 	if controller_tools is None or len(controller_tools) == 0{
	// 		raise ToolProviderNotFoundError(f"workflow provider {provider_id} not found")
	// 	}
	// 	return controller_tools[0].fork_tool_runtime(
	// 		runtime={
	// 			"tenant_id": tenant_id,
	// 			"credentials": {},
	// 			"invoke_from": invoke_from,
	// 			"tool_invoke_from": tool_invoke_from,
	// 		}
	// 	)
	// }else if provider_type == "app"{
	// 	raise NotImplementedError("app provider not implemented")
	// }else{
	// 	raise ToolProviderNotFoundError(f"provider type {provider_type} not found")
	// }
}

func (tm *ToolManager) GetAgentToolRuntime(tenant_id, app_id string, agent_tool agententities.AgentToolEntity, invoke_from appenumtypes.InvokeFrom /* = InvokeFrom.DEBUGGER*/) *base.Tool {
	/*
		get the agent tool runtime
	*/
	// tool_entity = tm.get_tool_runtime(
	// 	provider_type=agent_tool.provider_type,
	// 	provider_id=agent_tool.provider_id,
	// 	tool_name=agent_tool.tool_name,
	// 	tenant_id=tenant_id,
	// 	invoke_from=invoke_from,
	// 	tool_invoke_from=ToolInvokeFrom.AGENT,
	// )
	// runtime_parameters = {}
	// parameters = tool_entity.get_all_runtime_parameters()
	// for parameter in parameters{
	// 	// check file types
	// 	if (
	// 		parameter.type
	// 		in {
	// 			ToolParameter.ToolParameterType.SYSTEM_FILES,
	// 			ToolParameter.ToolParameterType.FILE,
	// 			ToolParameter.ToolParameterType.FILES,
	// 		}
	// 		and parameter.required
	// 	){
	// 		raise ValueError(f"file type parameter {parameter.name} not supported in agent")
	// 	}
	// 	if parameter.form == ToolParameter.ToolParameterForm.FORM{
	// 		// save tool parameter to tool entity memory
	// 		value = tm._init_runtime_parameter(parameter, agent_tool.tool_parameters)
	// 		runtime_parameters[parameter.name] = value
	// 	}
	// }
	// // decrypt runtime parameters
	// encryption_manager = ToolParameterConfigurationManager(
	// 	tenant_id=tenant_id,
	// 	tool_runtime=tool_entity,
	// 	provider_name=agent_tool.provider_id,
	// 	provider_type=agent_tool.provider_type,
	// 	identity_id=f"AGENT.{app_id}",
	// )
	// runtime_parameters = encryption_manager.decrypt_tool_parameters(runtime_parameters)
	// if tool_entity.runtime is None or tool_entity.runtime.runtime_parameters is None{
	// 	raise ValueError("runtime not found or runtime parameters not found")
	// }
	// tool_entity.runtime.runtime_parameters.update(runtime_parameters)
	// return tool_entity
	return nil
}

func (tm *ToolManager) GetToolIcon(tenant_id string, provider_type string, provider_id string) /*-> Union[str, dict]:*/ any {
	/*
	   get the tool icon

	   :param tenant_id: the id of the tenant
	   :param provider_type: the type of the provider
	   :param provider_id: the id of the provider
	   :return:
	*/

	// provider: Optional[Union[BuiltinToolProvider, ApiToolProvider, WorkflowToolProvider]] = None
	if provider_type == "builtin" {
		return confy.GetWithDefault[string]("console_api_url", "http://127.0.0.1:5001") + "/console/api/workspaces/current/tool-provider/builtin/" + provider_id + "/icon"
	} else if provider_type == "api" {
		// try:
		provider := new(models.ApiToolProvider)
		err := dbengine.Instance().DB.Model(&models.ApiToolProvider{}).Where("tenant_id = ? and id = ?", tenant_id, provider_id).First(provider).Error
		if err != nil {
			mlog.Errorf("get ApiToolProvider failed:%v", err)
			provider = nil
		}

		if provider == nil {
			mlog.Errorf("api provider %s not found", provider_id)
			return map[string]any{"background": "#252525", "content": `\ud83d\ude01`}
		}
		var icon map[string]any
		err = json.Unmarshal([]byte(provider.Icon), &icon)
		if err != nil {
			mlog.Errorf("json unmashal(%s) failed:%v", provider.Icon, err)
			return provider.Icon
		}

		return icon
	} else if provider_type == "workflow" {
		provider := new(models.WorkflowToolProvider)
		err := dbengine.Instance().DB.Model(&models.WorkflowToolProvider{}).Where("tenant_id = ? and id = ?", tenant_id, provider_id).First(provider).Error
		if err != nil {
			mlog.Errorf("get WorkflowToolProvider failed:%v", err)
			provider = nil
		}
		if provider == nil {
			mlog.Errorf("workflow provider %s not found", provider_id)
			return map[string]any{"background": "#252525", "content": `\ud83d\ude01`}
		}

		var icon map[string]any
		err = json.Unmarshal([]byte(provider.Icon), &icon)
		if err != nil {
			mlog.Errorf("json unmashal(%s) failed:%v", provider.Icon, err)
			return provider.Icon
		}

		return icon
	} else {
		mlog.Errorf("provider type %s not found", provider_type)
		return nil
	}
}
