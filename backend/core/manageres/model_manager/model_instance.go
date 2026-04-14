package modelmanager

import (
	"errors"
	"fmt"
	"iter"
	"reflect"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	coreentities "mlib.com/gofy/server/entities/core"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	providerentities "mlib.com/gofy/server/entities/provider"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
)

type ModelInstance struct {
	ProviderModelBundle  *coreentities.ProviderModelBundle
	Model                string
	Provider             string
	Credentials          map[string]any
	ModelTypeInstance    modelruntimeentities.AIModeler
	LoadBalancingManager *LBModelManager
}

func NewModelInstance(provider_model_bundle *coreentities.ProviderModelBundle, model string) *ModelInstance {
	mi := &ModelInstance{}
	mi.ProviderModelBundle = provider_model_bundle
	mi.Model = model
	mi.Provider = provider_model_bundle.Configuration.Provider.Provider
	mi.Credentials = mi._fetch_credentials_from_bundle(provider_model_bundle, model)
	mi.ModelTypeInstance = mi.ProviderModelBundle.ModelTypeInstance
	mi.LoadBalancingManager = mi._get_load_balancing_manager(
		provider_model_bundle.Configuration,
		provider_model_bundle.ModelTypeInstance.ModelType(),
		model,
		mi.Credentials,
	)
	return mi
}

func (mi *ModelInstance) _fetch_credentials_from_bundle(provider_model_bundle *coreentities.ProviderModelBundle, model string) map[string]any {
	/*
		Fetch credentials from provider model bundle
		:param provider_model_bundle: provider model bundle
		:param model: model name
		:return:
	*/
	configuration := provider_model_bundle.Configuration
	model_type := provider_model_bundle.ModelTypeInstance.ModelType()
	credentials := configuration.GetCurrentCredentials(model_type, model)
	if credentials == nil {
		mlog.Errorf("GetCurrentCredentials failed")
		panic(exceptions.NewProviderTokenNotInitError(fmt.Sprintf("Model %s credentials is not initialized.", model)))
	}

	return credentials

}
func (mi *ModelInstance) _get_load_balancing_manager(
	configuration *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string, credentials map[string]any,
) *LBModelManager {
	/*
		Get load balancing model credentials
		:param configuration: provider configuration
		:param model_type: model type
		:param model: model name
		:param credentials: model credentials
		:return:
	*/
	if len(configuration.ModelSettings) > 0 && configuration.UsingProviderType == providerenumtypes.Provider_CUSTOM {
		var current_model_setting *providerentities.ModelSetting
		// check if model is disabled by admin
		for _, model_setting := range configuration.ModelSettings {
			if model_setting.ModelType == model_type && model_setting.Model == model {
				current_model_setting = model_setting
				break
			}
		}
		// check if load balancing is enabled
		if current_model_setting != nil && len(current_model_setting.LoadBalancingConfigs) > 0 {
			// use load balancing proxy to choose credentials
			var managed_credentials map[string]any
			if configuration.CustomConfiguration.Provider != nil {
				managed_credentials = credentials
			}
			lb_model_manager := NewLBModelManager(
				configuration.TenantID,
				configuration.Provider.Provider,
				model_type,
				model,
				current_model_setting.LoadBalancingConfigs,
				managed_credentials,
			)

			return lb_model_manager
		}
	}
	return nil
}

func (mi *ModelInstance) GetLLMNumTokens(
	prompt_messages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool,
) int {
	/*
		Get number of tokens for llm

		:param prompt_messages: prompt messages
		:param tools: tools for tool calling
		:return:
	*/
	if !reflect.TypeOf(mi.ModelTypeInstance).Implements(reflect.TypeOf((*modelruntimeentities.LargeLanguageModeler)(nil)).Elem()) {
		panic(errors.New("model type instance is not LargeLanguageModel"))
	}

	new_model_type_instance := any(mi.ModelTypeInstance).(modelruntimeentities.LargeLanguageModeler)
	val := mi._round_robin_invoke(new_model_type_instance.GetNumTokens, mi.Model, mi.Credentials, prompt_messages, tools)
	return val.(int)
}

func (mi *ModelInstance) _round_robin_invoke_function(function any, args ...any) any {
	switch realfunc := function.(type) {

	case func(string, map[string]any, []modelruntimeentities.PromptMessager, []*modelruntimeentities.PromptMessageTool) int:
		if len(args) != 4 {
			mlog.Errorf("need 4 args, but only provided %d args", len(args))
			panic(exceptions.NewValueError(fmt.Sprintf("need 4 args, but only provided %d args", len(args))))
		}
		if _, ok := args[0].(string); !ok {
			mlog.Errorf("args[0]=%#v must be string", args[0])
			panic(exceptions.NewValueError(fmt.Sprintf("args[0]=%#v must be string", args[0])))
		}
		if args[1] != nil {
			if _, ok := args[1].(map[string]any); !ok {
				mlog.Errorf("args[1]=%#v must be map[string]any", args[1])
				panic(exceptions.NewValueError(fmt.Sprintf("args[1]=%#v must be map[string]any", args[1])))
			}
		}
		if args[2] != nil {
			if _, ok := args[2].([]modelruntimeentities.PromptMessager); !ok {
				mlog.Errorf("args[2]=%#v must be []modelruntimeentities.PromptMessager", args[2])
				panic(exceptions.NewValueError(fmt.Sprintf("args[2]=%#v must be []modelruntimeentities.PromptMessager", args[2])))
			}
		}
		if args[3] != nil {
			if _, ok := args[3].([]*modelruntimeentities.PromptMessageTool); !ok {
				mlog.Errorf("args[3]=%#v must be []*modelruntimeentities.PromptMessageTool", args[3])
				panic(exceptions.NewValueError(fmt.Sprintf("args[3]=%#v must be []*modelruntimeentities.PromptMessageTool", args[3])))
			}
		}
		return realfunc(args[0].(string), args[1].(map[string]any), args[2].([]modelruntimeentities.PromptMessager), args[3].([]*modelruntimeentities.PromptMessageTool))
	case func(string, map[string]any, []modelruntimeentities.PromptMessager, map[string]any, []*modelruntimeentities.PromptMessageTool, []string, string) *modelruntimeentities.LLMResult:
		if len(args) != 7 {
			mlog.Errorf("need 7 args, but only provided %d args", len(args))
			panic(exceptions.NewValueError(fmt.Sprintf("need 8 args, but only provided %d args", len(args))))
		}
		if _, ok := args[0].(string); !ok {
			mlog.Errorf("args[0]=%#v must be string", args[0])
			panic(exceptions.NewValueError(fmt.Sprintf("args[0]=%#v must be string", args[0])))
		}
		if args[1] != nil {
			if _, ok := args[1].(map[string]any); !ok {
				mlog.Errorf("args[1]=%#v must be map[string]any", args[1])
				panic(exceptions.NewValueError(fmt.Sprintf("args[1]=%#v must be map[string]any", args[1])))
			}
		}
		if args[2] != nil {
			if _, ok := args[2].([]modelruntimeentities.PromptMessager); !ok {
				mlog.Errorf("args[2]=%#v must be []modelruntimeentities.PromptMessager", args[2])
				panic(exceptions.NewValueError(fmt.Sprintf("args[2]=%#v must be []modelruntimeentities.PromptMessager", args[2])))
			}
		}
		if args[3] != nil {
			if _, ok := args[3].(map[string]any); !ok {
				mlog.Errorf("args[3]=%#v must be map[string]any", args[3])
				panic(exceptions.NewValueError(fmt.Sprintf("args[3]=%#v must be map[string]any", args[3])))
			}
		}
		if args[4] != nil {
			if _, ok := args[4].([]*modelruntimeentities.PromptMessageTool); !ok {
				mlog.Errorf("args[4]=%#v must be []*modelruntimeentities.PromptMessageTool", args[4])
				panic(exceptions.NewValueError(fmt.Sprintf("args[4]=%#v must be []*modelruntimeentities.PromptMessageTool", args[4])))
			}
		}
		if args[5] != nil {
			if _, ok := args[5].([]string); !ok {
				mlog.Errorf("args[5]=%#v must be []string", args[5])
				panic(exceptions.NewValueError(fmt.Sprintf("args[5]=%#v must be []string", args[5])))
			}
		}
		if _, ok := args[6].(string); !ok {
			mlog.Errorf("args[6]=%#v must be string", args[6])
			panic(exceptions.NewValueError(fmt.Sprintf("args[6]=%#v must be string", args[6])))
		}
		return realfunc(args[0].(string), args[1].(map[string]any), args[2].([]modelruntimeentities.PromptMessager), args[3].(map[string]any), args[4].([]*modelruntimeentities.PromptMessageTool), args[5].([]string), args[6].(string))
	case func(string, map[string]any, []modelruntimeentities.PromptMessager, map[string]any, []*modelruntimeentities.PromptMessageTool, []string, string) iter.Seq[*modelruntimeentities.LLMResultChunk]:
		if len(args) != 7 {
			mlog.Errorf("need 7 args, but only provided %d args", len(args))
			panic(exceptions.NewValueError(fmt.Sprintf("need 8 args, but only provided %d args", len(args))))
		}
		if _, ok := args[0].(string); !ok {
			mlog.Errorf("args[0]=%#v must be string", args[0])
			panic(exceptions.NewValueError(fmt.Sprintf("args[0]=%#v must be string", args[0])))
		}
		if args[1] != nil {
			if _, ok := args[1].(map[string]any); !ok {
				mlog.Errorf("args[1]=%#v must be map[string]any", args[1])
				panic(exceptions.NewValueError(fmt.Sprintf("args[1]=%#v must be map[string]any", args[1])))
			}
		}
		if args[2] != nil {
			if _, ok := args[2].([]modelruntimeentities.PromptMessager); !ok {
				mlog.Errorf("args[2]=%#v must be []modelruntimeentities.PromptMessager", args[2])
				panic(exceptions.NewValueError(fmt.Sprintf("args[2]=%#v must be []modelruntimeentities.PromptMessager", args[2])))
			}
		}
		if args[3] != nil {
			if _, ok := args[3].(map[string]any); !ok {
				mlog.Errorf("args[3]=%#v must be map[string]any", args[3])
				panic(exceptions.NewValueError(fmt.Sprintf("args[3]=%#v must be map[string]any", args[3])))
			}
		}
		if args[4] != nil {
			if _, ok := args[4].([]*modelruntimeentities.PromptMessageTool); !ok {
				mlog.Errorf("args[4]=%#v must be []*modelruntimeentities.PromptMessageTool", args[4])
				panic(exceptions.NewValueError(fmt.Sprintf("args[4]=%#v must be []*modelruntimeentities.PromptMessageTool", args[4])))
			}
		}
		if args[5] != nil {
			if _, ok := args[5].([]string); !ok {
				mlog.Errorf("args[5]=%#v must be []string", args[5])
				panic(exceptions.NewValueError(fmt.Sprintf("args[5]=%#v must be []string", args[5])))
			}
		}
		if _, ok := args[6].(string); !ok {
			mlog.Errorf("args[6]=%#v must be string", args[6])
			panic(exceptions.NewValueError(fmt.Sprintf("args[6]=%#v must be string", args[6])))
		}
		return realfunc(args[0].(string), args[1].(map[string]any), args[2].([]modelruntimeentities.PromptMessager), args[3].(map[string]any), args[4].([]*modelruntimeentities.PromptMessageTool), args[5].([]string), args[6].(string))
	}
	panic(exceptions.NewValueError(fmt.Sprintf("unsurported function(%#v)", function)))
}

func (mi *ModelInstance) _round_robin_invoke(function any, args ...any) any {
	/*
	   Round-robin invoke
	   :param function: function to invoke
	   :param args: function args
	   :param kwargs: function kwargs
	   :return:
	*/
	if mi.LoadBalancingManager == nil {
		return mi._round_robin_invoke_function(function, args...)
	}

	var last_exception error
	// last_exception: Union[InvokeRateLimitError, InvokeAuthorizationError, InvokeConnectionError, None] = None
	for {
		lb_config := mi.LoadBalancingManager.fetch_next()
		if lb_config == nil {
			if last_exception == nil {
				panic(exceptions.NewProviderTokenNotInitError("Model credentials is not initialized."))
			} else {
				panic(last_exception)
			}
		}
		is_continue := false
		res := func() any {
			defer func() {
				if r := recover(); r != nil {
					if exp, ok := r.(*modelruntimeexceptions.InvokeRateLimitError); ok {
						// expire in 60 seconds
						mi.LoadBalancingManager.cooldown(lb_config, 60)
						last_exception = exp
						is_continue = true
					} else if exp, ok := r.(*modelruntimeexceptions.InvokeAuthorizationError); ok {
						// expire in 10 seconds
						mi.LoadBalancingManager.cooldown(lb_config, 10)
						last_exception = exp
						is_continue = true
					} else if exp, ok := r.(*modelruntimeexceptions.InvokeConnectionError); ok {
						// expire in 10 seconds
						mi.LoadBalancingManager.cooldown(lb_config, 10)
						last_exception = exp
						is_continue = true
					} else {
						panic(r)
					}
				}
			}()
			return mi._round_robin_invoke_function(function, args...)
		}()

		if is_continue {
			continue
		}
		return res
	}
	// return nil, nil
}

func (mi *ModelInstance) InvokeLLM(
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
	callbacks []modelruntimeentities.Callbacker,
) *modelruntimeentities.LLMResult {
	/*
		Invoke large language model

		:param prompt_messages: prompt messages
		:param model_parameters: model parameters
		:param tools: tools for tool calling
		:param stop: stop words
		:param stream: is stream response
		:param user: unique user id
		:param callbacks: callbacks
		:return: full response or stream response chunk generator result
	*/
	if !reflect.TypeOf(mi.ModelTypeInstance).Implements(reflect.TypeOf((*modelruntimeentities.LargeLanguageModeler)(nil)).Elem()) {
		panic(errors.New("model type instance is not LargeLanguageModel"))
	}
	new_model_type_instance := any(mi.ModelTypeInstance).(modelruntimeentities.LargeLanguageModeler)
	val := mi._round_robin_invoke(
		new_model_type_instance.Invoke,
		mi.Model,
		mi.Credentials,
		prompt_messages,
		model_parameters,
		tools,
		stop,
		user,
	)

	return val.(*modelruntimeentities.LLMResult)
}

func (mi *ModelInstance) InvokeLLMStream(
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
	callbacks []modelruntimeentities.Callbacker,
) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	/*
		Invoke large language model

		:param prompt_messages: prompt messages
		:param model_parameters: model parameters
		:param tools: tools for tool calling
		:param stop: stop words
		:param stream: is stream response
		:param user: unique user id
		:param callbacks: callbacks
		:return: full response or stream response chunk generator result
	*/
	if !reflect.TypeOf(mi.ModelTypeInstance).Implements(reflect.TypeOf((*modelruntimeentities.LargeLanguageModeler)(nil)).Elem()) {
		panic(errors.New("model type instance is not LargeLanguageModel"))
	}
	new_model_type_instance := any(mi.ModelTypeInstance).(modelruntimeentities.LargeLanguageModeler)
	val := mi._round_robin_invoke(
		new_model_type_instance.InvokeStream,
		mi.Model,
		mi.Credentials,
		prompt_messages,
		model_parameters,
		tools,
		stop,
		user,
	)
	return val.(iter.Seq[*modelruntimeentities.LLMResultChunk])
}
