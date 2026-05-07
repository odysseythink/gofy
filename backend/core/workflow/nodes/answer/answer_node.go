package answer

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/file"
	"github.com/odysseythink/gofy/backend/core/variables"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	answergeneraterouter "github.com/odysseythink/gofy/backend/core/workflow/nodes_generate_router/answer"
	variabletemplateparser "github.com/odysseythink/gofy/backend/core/workflow/utils/variable_template_parser"
	answernodesentities "github.com/odysseythink/gofy/backend/entities/nodes/answer"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	answernodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes/answer"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type AnswerNode struct {
	answer string
	*base.BaseNode[*answernodesentities.AnswerNodeData]
}

func (n *AnswerNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_ANSWER
}

func (n *AnswerNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	// generate routes
	generate_routes := (&answergeneraterouter.AnswerStreamGeneratorRouter{}).ExtractGenerateRouteFromNodeData(n.NodeData)

	answer := ""
	files := []*file.File{}
	for _, part := range generate_routes {
		if part.Type() == answernodesenumtypes.GenerateRouteChunk_VAR {
			new_part := any(part).(*answernodesentities.VarGenerateRouteChunk)
			value_selector := new_part.ValueSelector
			variable := n.GetGraphRuntimeState().VariablePool.Get(value_selector)
			if variable != nil {
				mlog.Debugf("------value_selector=%#v, variable=%#v", value_selector, variable)
				if real_segment, ok := any(variable).(*variables.FileVariable); ok {
					files = append(files, real_segment.Value)
				} else if real_segment, ok := any(variable).(*variables.ArrayFileVariable); ok {
					files = append(files, real_segment.Value...)
				}
				answer += variable.Markdown()
			}
		} else {
			new_part := any(part).(*answernodesentities.TextGenerateRouteChunk)
			mlog.Debugf("------Text=%#v", new_part.Text)
			answer += new_part.Text
		}
	}
	return &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Outputs: map[string]any{"answer": answer, "files": files},
	}, nil
}

func (n *AnswerNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *answernodesentities.AnswerNodeData) map[string][]string {
	variable_template_parser := variabletemplateparser.NewVariableTemplateParser(node_data.Answer)
	variable_selectors := variable_template_parser.ExtractVariableSelectors()

	variable_mapping := map[string][]string{}
	for _, variable_selector := range variable_selectors {
		variable_mapping[node_id+"."+variable_selector.Variable] = variable_selector.ValueSelector
	}
	return variable_mapping
}

func New() *AnswerNode {
	return &AnswerNode{
		BaseNode: &base.BaseNode[*answernodesentities.AnswerNodeData]{},
	}
}
