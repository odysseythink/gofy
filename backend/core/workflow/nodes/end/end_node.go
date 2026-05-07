package end

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	endnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/end"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
)

type EndNode struct {
	*base.BaseNode[*endnodesentities.EndNodeData]
}

func (n *EndNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_END
}

func (n *EndNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	output_variables := n.NodeData.Outputs

	outputs := map[string]any{}
	for _, variable_selector := range output_variables {
		variable := n.GetGraphRuntimeState().VariablePool.Get(variable_selector.ValueSelector)
		var value any
		if variable != nil {
			value = variable.ToObject()
		}
		outputs[variable_selector.Variable] = value
	}
	return &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:  outputs,
		Outputs: outputs,
	}, nil
}

func (n *EndNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *endnodesentities.EndNodeData) map[string][]string {
	return map[string][]string{}
}

func New() *EndNode {
	return &EndNode{
		BaseNode: &base.BaseNode[*endnodesentities.EndNodeData]{},
	}
}
