package code

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"mlib.com/confy"
	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	codenodesexceptions "mlib.com/gofy/server/core/exceptions/nodes/code"
	codeexecutor "mlib.com/gofy/server/core/helper/code_executor"
	python3codeexecutor "mlib.com/gofy/server/core/helper/code_executor/template_transformer/python3"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	codenodesentities "mlib.com/gofy/server/entities/nodes/code"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	codeexecutorenumtypes "mlib.com/gofy/server/enum_types/code_executor"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)

type CodeNode struct {
	*base.BaseNode[*codenodesentities.CodeNodeData]
}

func (n *CodeNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_CODE
}

func (n *CodeNode) Run() (run_result *workflowentities.NodeRunResult, stream_run_result iter.Seq[any]) {
	// Get code language
	code_language := n.NodeData.CodeLanguage
	code := n.NodeData.Code

	// Get variables
	variables := map[string]any{}
	for _, variable_selector := range n.NodeData.Variables {
		variable_name := variable_selector.Variable
		variable := n.GetGraphRuntimeState().VariablePool.Get(variable_selector.ValueSelector)
		if variable != nil {
			variables[variable_name] = variable.ToObject()
		} else {
			variables[variable_name] = nil
		}
	}
	// Run code
	mlog.Debugf("---code_language=%v", code_language)
	defer func() {
		if r := recover(); r != nil {
			if exp, ok := r.(*codenodesexceptions.CodeExecutionError); ok {
				run_result = &workflowentities.NodeRunResult{
					Status:    models.WorkflowNodeExecutionStatus_FAILED,
					Inputs:    variables,
					Error:     exp.Error(),
					ErrorType: "CodeExecutionError",
				}
			} else if exp, ok := r.(*codenodesexceptions.CodeNodeError); ok {
				run_result = &workflowentities.NodeRunResult{
					Status:    models.WorkflowNodeExecutionStatus_FAILED,
					Inputs:    variables,
					Error:     exp.Error(),
					ErrorType: "CodeNodeError",
				}
			} else {
				panic(r)
			}
		}
	}()
	result := codeexecutor.ExecuteWorkflowCodeTemplate(
		codeexecutorenumtypes.CodeLanguage(code_language),
		code,
		variables,
	)

	// Transform result
	result = transformResult(result, n.NodeData.Outputs, "", 1)

	return &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:  variables,
		Outputs: result,
	}, nil
}

func (n *CodeNode) GetDefaultConfig(filters map[string]any) map[string]any {
	/*
	   Get default config of node.
	   :param filters: filter by node config parameters.
	   :return:
	*/
	code_language := codeexecutorenumtypes.CodeLanguage_PYTHON3
	if len(filters) > 0 {
		if _, ok := filters["code_language"]; ok {
			if _, ok := filters["code_language"].(string); ok {
				code_language = codeexecutorenumtypes.CodeLanguage(filters["code_language"].(string))
			}
		}
	}
	providers := []codeexecutor.CodeNodeProvider{&python3codeexecutor.Python3CodeProvider{}}
	var code_provider codeexecutor.CodeNodeProvider
	for _, v := range providers {
		if v.IsAcceptLanguage(code_language) {
			code_provider = v
			break
		}
	}
	if code_provider != nil {
		return code_provider.GetDefaultConfig()
	} else {
		return nil
	}
}
func (n *CodeNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	/*
		Extract variable selector to variable mapping
		:param graph_config: graph config
		:param node_id: node id
		:param node_data: node data
		:return:
	*/
	typed_node_data, err := mapstruct.MapToStruct1[*codenodesentities.CodeNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}

	res := map[string][]string{}
	for _, variable_selector := range typed_node_data.Variables {
		res[node_id+"."+variable_selector.Variable] = variable_selector.ValueSelector
	}
	return res
}

func checkString(value string, variable string) string {
	/*
	   Check string
	   :param value: value
	   :param variable: variable
	   :return:
	*/

	if len(value) > confy.GetWithDefault[int]("code_execution_config.max_string_length", 80000) {
		panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("The length of output variable `%#v` must be less than %d characters", variable, confy.GetWithDefault[int]("code_execution_config.max_string_length", 80000))))
	}
	return strings.ReplaceAll(value, "\x00", "")
}
func checkNumber[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](value T, variable string) T {
	/*
	   Check number
	   :param value: value
	   :param variable: variable
	   :return:
	*/

	if int(value) > confy.GetWithDefault[int]("code_execution_config.max_number", 9223372036854775807) || int(value) < confy.GetWithDefault[int]("code_execution_config.min_number", -9223372036854775808) {
		panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output variable `%s` is out of range, it must be between %d and %d.", variable, confy.GetWithDefault[int]("code_execution_config.min_number", -9223372036854775808), confy.GetWithDefault[int]("code_execution_config.max_number", 9223372036854775807))))
	}
	switch realdata := any(value).(type) {
	case float32:
		if len(strings.Split(fmt.Sprintf("%v", realdata), ".")[1]) > confy.GetWithDefault[int]("code_execution_config.max_precision", 20) {
			panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output variable `%s` has too high precision, it must be less than %d digits.", variable, confy.GetWithDefault[int]("code_execution_config.max_precision", 20))))
		}
	case float64:
		if len(strings.Split(fmt.Sprintf("%v", realdata), ".")[1]) > confy.GetWithDefault[int]("code_execution_config.max_precision", 20) {
			panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output variable `%s` has too high precision, it must be less than %d digits.", variable, confy.GetWithDefault[int]("code_execution_config.max_precision", 20))))
		}
	}

	return value
}

func transformResult(
	result map[string]any,
	outputSchema map[string]*codenodesentities.Output,
	prefix string,
	depth int, /* = 1*/
) map[string]any {
	if depth > confy.GetWithDefault[int]("code_execution_config.max_depth", 5) {
		panic(codenodesexceptions.NewDepthLimitError(fmt.Sprintf("Depth limit %d reached, object too deep.", confy.GetWithDefault[int]("code_execution_config.max_depth", 5))))
	}

	transformedResult := make(map[string]any)
	parametersValidated := make(map[string]bool)

	if outputSchema == nil {
		for outputName, outputValue := range result {
			if outputValue == nil {
				continue
			}
			realprefix := outputName
			if prefix != "" {
				realprefix = fmt.Sprintf("%s.%s", prefix, outputName)
			}

			switch v := outputValue.(type) {
			case map[string]any:
				transformResult(v, nil, realprefix, depth+1)
			case int:
				checkNumber(v, realprefix)
			case float64:
				checkNumber(v, realprefix)
			case string:
				checkString(v, realprefix)
			case []any:
				if len(v) == 0 {
					break
				}

				firstElement := v[0]
				if firstElement != nil {
					if _, ok := firstElement.(int); ok && !slices.ContainsFunc(v, func(sv any) bool {
						if _, ok := sv.(int); !ok {
							return true
						}
						return false
					}) {
						for i, sv := range v {
							checkNumber(sv.(int), fmt.Sprintf("%s[%d]", realprefix, i))
						}
					} else if _, ok := firstElement.(float64); ok && !slices.ContainsFunc(v, func(sv any) bool {
						if _, ok := sv.(float64); !ok {
							return true
						}
						return false
					}) {
						for i, sv := range v {
							checkNumber(sv.(float64), fmt.Sprintf("%s[%d]", realprefix, i))
						}
					} else if _, ok := firstElement.(string); ok && !slices.ContainsFunc(v, func(sv any) bool {
						if _, ok := sv.(string); !ok {
							return true
						}
						return false
					}) {
						for i, sv := range v {
							checkString(sv.(string), fmt.Sprintf("%s[%d]", realprefix, i))
						}
					} else if _, ok := firstElement.(map[string]any); ok && !slices.ContainsFunc(v, func(sv any) bool {
						if _, ok := sv.(map[string]any); !ok {
							return true
						}
						return false
					}) {
						for i, sv := range v {
							transformResult(sv.(map[string]any), nil, fmt.Sprintf("%s[%d]", realprefix, i), depth+1)
						}
					} else {
						panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not a valid array. make sure all elements are of the same type.", realprefix)))
					}
				}
			default:
				panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not a valid type.", realprefix)))
			}
		}
		return result
	}

	for outputName, outputConfig := range outputSchema {
		// dot := ""
		realprefix := outputName
		if prefix != "" {
			// dot = "."
			realprefix = fmt.Sprintf("%s.%s", prefix, outputName)
		}

		if _, ok := result[outputName]; !ok {
			panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is missing.", realprefix)))
		}
		if outputConfig == nil {
			panic(codenodesexceptions.NewOutputValidationError("Output can't be nil."))
		}

		switch outputConfig.Type {
		case "object":
			if value, ok := result[outputName].(map[string]any); ok {
				transformed := transformResult(value, outputConfig.Children, realprefix, depth+1)

				transformedResult[outputName] = transformed
			} else if result[outputName] == nil {
				transformedResult[outputName] = nil
			} else {
				panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not an object, got %T instead.", realprefix, result[outputName])))
			}
		case "number":
			if _, ok := result[outputName].(int); ok {
				validated := checkNumber(result[outputName].(int), realprefix)

				transformedResult[outputName] = validated
			} else if _, ok := result[outputName].(float64); ok {
				validated := checkNumber(result[outputName].(float64), realprefix)

				transformedResult[outputName] = validated
			} else {
				panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not an number, got %T instead.", realprefix, result[outputName])))
			}
		case "string":
			if value, ok := result[outputName].(string); ok {
				validated := checkString(value, realprefix)

				transformedResult[outputName] = validated
			} else {
				panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not a string, got %T instead.", realprefix, result[outputName])))
			}
		case "array[number]":
			if value, ok := result[outputName].([]any); ok {
				if len(value) > confy.GetWithDefault[int]("code_execution_config.max_number_array_length", 1000) {
					panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("The length of output variable `%s` must be less than %d elements.", realprefix, confy.GetWithDefault[int]("code_execution_config.max_number_array_length", 1000))))
				}

				validated := make([]any, 0)
				for i, sv := range value {
					if _, ok := sv.(int); ok {
						valid := checkNumber(sv.(int), fmt.Sprintf("%s[%d]", realprefix, i))
						validated = append(validated, valid)
					} else if _, ok := sv.(float64); ok {
						valid := checkNumber(sv.(float64), fmt.Sprintf("%s[%d]", realprefix, i))
						validated = append(validated, valid)
					} else {
						panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not an number array, got %#v instead.", realprefix, result[outputName])))
					}
				}
				transformedResult[outputName] = validated
			} else if result[outputName] == nil {
				transformedResult[outputName] = nil
			} else {
				panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not an array, got %T instead.", realprefix, result[outputName])))
			}
		case "array[string]":
			if value, ok := result[outputName].([]any); ok {
				if len(value) > confy.GetWithDefault[int]("code_execution_config.max_string_array_length", 30) {
					panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("The length of output variable `%s` must be less than %d elements.", realprefix, confy.GetWithDefault[int]("code_execution_config.max_string_array_length", 30))))
				}

				validated := make([]any, len(value))
				for i, val := range value {
					if strVal, ok := val.(string); ok {
						valid := checkString(strVal, fmt.Sprintf("%s[%d]", realprefix, i))
						validated[i] = valid
					} else {
						panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s[%d] is not a string, got %T instead.", realprefix, i, val)))
					}
				}
				transformedResult[outputName] = validated
			} else if result[outputName] == nil {
				transformedResult[outputName] = nil
			} else {
				panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not an array, got %T instead.", realprefix, result[outputName])))
			}
		case "array[object]":
			if value, ok := result[outputName].([]any); ok {
				if len(value) > confy.GetWithDefault[int]("code_execution_config.max_object_array_length", 30) {
					panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("The length of output variable `%s` must be less than %d elements.", realprefix, confy.GetWithDefault[int]("code_execution_config.max_object_array_length", 30))))
				}

				validated := make([]any, len(value))
				for i, val := range value {
					if val == nil {
						validated[i] = nil
					} else {
						if obj, ok := val.(map[string]any); ok {
							transformed := transformResult(obj, outputConfig.Children, fmt.Sprintf("%s[%d]", realprefix, i), depth+1)
							validated[i] = transformed
						} else {
							panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s[%d] is not an object, got %T instead at index %d.", realprefix, i, val, i)))
						}
					}
				}
				transformedResult[outputName] = validated
			} else if result[outputName] == nil {
				transformedResult[outputName] = nil
			} else {
				panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output %s is not an array, got %T instead.", realprefix, result[outputName])))
			}
		default:
			panic(codenodesexceptions.NewOutputValidationError(fmt.Sprintf("Output type %s is not supported.", outputConfig.Type)))
		}

		parametersValidated[outputName] = true
	}

	if len(parametersValidated) != len(result) {
		panic(codenodesexceptions.NewCodeNodeError("Not all output parameters are validated."))
	}

	return transformedResult
}
func New() *CodeNode {
	return &CodeNode{
		BaseNode: &base.BaseNode[*codenodesentities.CodeNodeData]{},
	}
}
func init() {
	nodesconstants.Regist(New())
}
