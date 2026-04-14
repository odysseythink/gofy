package ifelse

import (
	"iter"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	"mlib.com/gofy/server/core/workflow/utils/condition"
	ifelsenodesentities "mlib.com/gofy/server/entities/nodes/if_else"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	conditionentities "mlib.com/gofy/server/entities/workflow/condition"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
)

type IfElseNode struct {
	*base.BaseNode[*ifelsenodesentities.IfElseNodeData]
}

func (n *IfElseNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_IF_ELSE
}

func (n *IfElseNode) Run() (run_result *workflowentities.NodeRunResult, run_stream_result iter.Seq[any]) {

	/*
	   Run node
	   :return:
	*/
	node_inputs := map[string]any{"conditions": []any{}}

	process_data := map[string]any{"condition_results": []any{}}

	input_conditions := []map[string]any{}
	final_result := false
	selected_case_id := ""
	condition_processor := &condition.ConditionProcessor{}
	group_result := []bool{}
	mlog.Debugf("*********NodeData=%#v", n.BaseNode.NodeData)
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				run_result = &workflowentities.NodeRunResult{
					Status:      models.WorkflowNodeExecutionStatus_FAILED,
					Inputs:      node_inputs,
					ProcessData: process_data,
					Error:       err.Error(),
				}
			} else {
				panic(r)
			}
		}
	}()
	// Check if the new cases structure is used
	if len(n.BaseNode.NodeData.Cases) > 0 {
		for _, item := range n.BaseNode.NodeData.Cases {
			mlog.Debugf("-------item=%#v", item)
			for _, v := range item.Conditions {
				mlog.Debugf("         condition=%#v", v)
			}
			input_conditions, group_result, final_result = condition_processor.ProcessConditions(
				n.GetGraphRuntimeState().VariablePool,
				item.Conditions,
				item.LogicalOperator,
			)
			if _, ok := process_data["condition_results"]; !ok {
				process_data["condition_results"] = make([]map[string]any, 0)
			}
			process_data["condition_results"] = append(process_data["condition_results"].([]any), map[string]any{
				"group":        item.ModelDump(),
				"results":      group_result,
				"final_result": final_result,
			})

			// Break if a case passes (logical short-circuit)
			if final_result {
				selected_case_id = item.CaseID // Capture the ID of the passing case
				break
			}
		}
	} else {
		// TODO: Update database then remove this
		// Fallback to old structure if cases are not defined
		logical_operator := "and"
		if n.BaseNode.NodeData.LogicalOperator != "" {
			logical_operator = n.BaseNode.NodeData.LogicalOperator
		}
		input_conditions, group_result, final_result = _should_not_use_old_function(
			condition_processor,
			n.GetGraphRuntimeState().VariablePool,
			n.BaseNode.NodeData.Conditions,
			logical_operator,
		)
		if final_result {
			selected_case_id = "true"
		} else {
			selected_case_id = "false"
		}
		if _, ok := process_data["condition_results"]; !ok {
			process_data["condition_results"] = make([]map[string]any, 0)
		}
		process_data["condition_results"] = append(process_data["condition_results"].([]any), map[string]any{"group": "default", "results": group_result, "final_result": final_result})
	}
	node_inputs["conditions"] = input_conditions
	outputs := map[string]any{"result": final_result, "selected_case_id": selected_case_id}

	run_result = &workflowentities.NodeRunResult{
		Status:      models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:      node_inputs,
		ProcessData: process_data,
		// edge_source_handle=selected_case_id or "false",  // Use case ID or 'default'
		Outputs: outputs,
	}
	if selected_case_id != "" {
		run_result.EdgeSourceHandle = selected_case_id
	} else {
		run_result.EdgeSourceHandle = "false"
	}

	return run_result, nil
}

func (n *IfElseNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *ifelsenodesentities.IfElseNodeData) map[string][]string {

	return map[string][]string{}
}
func New() *IfElseNode {
	return &IfElseNode{
		BaseNode: &base.BaseNode[*ifelsenodesentities.IfElseNodeData]{},
	}
}

func _should_not_use_old_function(
	condition_processor *condition.ConditionProcessor,
	variable_pool *workflowentities.VariablePool,
	conditions []*conditionentities.Condition,
	operator string, /* Literal["and", "or"]*/
) ([]map[string]any, []bool, bool) {
	return condition_processor.ProcessConditions(
		variable_pool,
		conditions,
		operator,
	)
}
