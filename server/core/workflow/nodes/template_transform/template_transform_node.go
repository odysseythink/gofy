package templatetransform

import (
	"fmt"
	"iter"

	"mlib.com/confy"
	codeexecutor "mlib.com/gofy/server/core/helper/code_executor"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	templatetransformnodesentities "mlib.com/gofy/server/entities/nodes/template_transform"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	codeexecutorenumtypes "mlib.com/gofy/server/enum_types/code_executor"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type TemplateTransformNode struct {
	*base.BaseNode[*templatetransformnodesentities.TemplateTransformNodeData]
}

func (n *TemplateTransformNode) GetDefaultConfig(filters map[string]any) map[string]any {
	return map[string]any{
		"type": "template-transform",
		"config": map[string]any{
			"variables": []map[string]any{
				{"variable": "arg1", "value_selector": []any{}},
			},
			"template": "{{ arg1 }}",
		},
	}
}

func (n *TemplateTransformNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_TEMPLATE_TRANSFORM
}

func (n *TemplateTransformNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	variables := map[string]any{}
	for _, variable_selector := range n.NodeData.Variables {
		variable_name := variable_selector.Variable
		value := n.GetGraphRuntimeState().VariablePool.Get(variable_selector.ValueSelector)
		if value != nil {
			variables[variable_name] = value.ToObject()
		}
	}
	// Run code
	result := func() (res *workflowentities.NodeRunResult) {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(error); ok {
					res = &workflowentities.NodeRunResult{
						Inputs: variables,
						Status: models.WorkflowNodeExecutionStatus_FAILED,
						Error:  exp.Error(),
					}
				} else {
					panic(r)
				}
			}
		}()
		result_map := codeexecutor.ExecuteWorkflowCodeTemplate(codeexecutorenumtypes.CodeLanguage_JINJA2, n.NodeData.Template, variables)
		result_str := ""
		if _, ok := result_map["result"]; ok {
			if _, ok := result_map["result"].(string); ok {
				result_str = result_map["result"].(string)
			} else {
				mlog.Errorf("execute workflow code template return result field(%#v) must be string", result_map["result"])
				res = &workflowentities.NodeRunResult{
					Inputs: variables,
					Status: models.WorkflowNodeExecutionStatus_FAILED,
					Error:  "execute workflow code template return result field must be string",
				}
				return
			}
		} else {
			mlog.Errorf("execute workflow code template doesn't return=%#v result field", result_map)
			res = &workflowentities.NodeRunResult{
				Inputs: variables,
				Status: models.WorkflowNodeExecutionStatus_FAILED,
				Error:  "execute workflow code template doesn't return result field",
			}
			return
		}
		if len(result_str) > confy.GetWithDefault[int]("template_transform_max_length", 80000) {
			mlog.Errorf("Output length exceeds %d characters", confy.GetWithDefault[int]("template_transform_max_length", 80000))
			res = &workflowentities.NodeRunResult{
				Inputs: variables,
				Status: models.WorkflowNodeExecutionStatus_FAILED,
				Error:  fmt.Sprintf("Output length exceeds %d characters", confy.GetWithDefault[int]("template_transform_max_length", 80000)),
			}
			return
		}
		res = &workflowentities.NodeRunResult{
			Inputs:  variables,
			Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
			Outputs: map[string]any{"output": result_str},
		}
		return
	}()

	return result, nil
}

func (n *TemplateTransformNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *templatetransformnodesentities.TemplateTransformNodeData) map[string][]string {
	ret := map[string][]string{}
	for _, variable_selector := range node_data.Variables {
		ret[node_id+"."+variable_selector.Variable] = variable_selector.ValueSelector
	}
	return ret
}
func New() *TemplateTransformNode {
	return &TemplateTransformNode{
		BaseNode: &base.BaseNode[*templatetransformnodesentities.TemplateTransformNodeData]{},
	}
}
