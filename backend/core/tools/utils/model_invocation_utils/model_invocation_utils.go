package modelinvocationutils

import (
	"encoding/json"
	"time"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	toolsexceptions "mlib.com/gofy/server/core/exceptions/tools"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	dbengine "mlib.com/gofy/server/db_engine"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
)

func GetMaxLLMContextTokens(
	tenant_id string,
) int {
	model_instance := (&modelmanager.ModelManager{}).GetDefaultModelInstance(
		tenant_id,
		modelruntimeenumtypes.Model_LLM,
	)

	if model_instance == nil {
		panic(toolsexceptions.NewInvokeModelError("Model not found"))
	}

	llm_model := any(model_instance.ModelTypeInstance).(modelruntimeentities.LargeLanguageModeler)
	schema := llm_model.GetModelSchema(model_instance.ModelTypeInstance, model_instance.Model, model_instance.Credentials)

	if schema == nil {
		panic(toolsexceptions.NewInvokeModelError("No model schema found"))
	}

	max_tokens := 0
	if _, ok := schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE]; ok {
		if _, ok := schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE].(int); ok {
			max_tokens = schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE].(int)
		}
	}
	if max_tokens <= 0 {
		return 2048
	}
	return max_tokens
}
func CalculateTokens(tenant_id string, prompt_messages []modelruntimeentities.PromptMessager) int {
	// get model instance
	model_instance := (&modelmanager.ModelManager{}).GetDefaultModelInstance(
		tenant_id,
		modelruntimeenumtypes.Model_LLM,
	)

	if model_instance == nil {
		panic(toolsexceptions.NewInvokeModelError("Model not found"))
	}

	// get tokens
	return model_instance.GetLLMNumTokens(prompt_messages, nil)
}

func Invoke(
	user_id string, tenant_id string, tool_type string, tool_name string, prompt_messages []modelruntimeentities.PromptMessager,
) *modelruntimeentities.LLMResult {
	// get model manager
	model_instance := (&modelmanager.ModelManager{}).GetDefaultModelInstance(
		tenant_id,
		modelruntimeenumtypes.Model_LLM,
	)
	if model_instance == nil {
		panic(toolsexceptions.NewInvokeModelError("Model not found"))
	}
	// get prompt tokens
	prompt_tokens := model_instance.GetLLMNumTokens(prompt_messages, nil)

	model_parameters := map[string]any{
		"temperature": 0.8,
		"top_p":       0.8,
	}

	// create tool model invoke
	now := time.Now()
	tool_model_invoke := &models.ToolModelInvoke{
		ID:       uuid.NewV4().String(),
		UserID:   user_id,
		TenantID: tenant_id,
		Provider: model_instance.Provider,
		ToolType: tool_type,
		ToolName: tool_name,
		// ModelParameters        :,
		// PromptMessages         :,
		ModelResponse: "",
		PromptTokens:  prompt_tokens,
		Currency:      "USD",
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	bindata, _ := json.Marshal(model_parameters)
	tool_model_invoke.ModelParameters = string(bindata)
	bindata, _ = json.Marshal(prompt_messages)
	tool_model_invoke.PromptMessages = string(bindata)
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
		dbengine.Instance().DB.Create(tool_model_invoke)
	}()
	response := func() *modelruntimeentities.LLMResult {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if e, ok := r.(*modelruntimeexceptions.InvokeRateLimitError); ok {
					panic(toolsexceptions.NewInvokeModelError("Invoke rate limit error: " + e.Error()))
				} else if e, ok := r.(*modelruntimeexceptions.InvokeBadRequestError); ok {
					panic(toolsexceptions.NewInvokeModelError("Invoke bad request error: " + e.Error()))
				} else if e, ok := r.(*modelruntimeexceptions.InvokeConnectionError); ok {
					panic(toolsexceptions.NewInvokeModelError("Invoke connection error: " + e.Error()))
				} else if _, ok := r.(*modelruntimeexceptions.InvokeAuthorizationError); ok {
					panic(toolsexceptions.NewInvokeModelError("Invoke authorization error"))
				} else if e, ok := r.(*modelruntimeexceptions.InvokeServerUnavailableError); ok {
					panic(toolsexceptions.NewInvokeModelError("Invoke server unavailable error: " + e.Error()))
				} else if e, ok := r.(error); ok {
					panic(toolsexceptions.NewInvokeModelError("Invoke error: " + e.Error()))
				} else {
					panic(r)
				}
			}
		}()
		return model_instance.InvokeLLM(
			prompt_messages,
			model_parameters,
			nil,
			nil,
			user_id,
			nil,
		)
	}()
	// update tool model invoke
	tool_model_invoke.ModelResponse = response.Message.Content
	if response.Usage != nil {
		tool_model_invoke.AnswerTokens = response.Usage.CompletionTokens
		tool_model_invoke.AnswerUnitPrice = response.Usage.CompletionUnitPrice
		tool_model_invoke.AnswerPriceUnit = response.Usage.CompletionPriceUnit
		tool_model_invoke.ProviderResponseLatency = response.Usage.Latency
		tool_model_invoke.TotalPrice = response.Usage.TotalPrice
		tool_model_invoke.Currency = response.Usage.Currency
	}

	return response
}
