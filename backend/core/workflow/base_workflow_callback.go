package workflow

import (
	appqueueentities "github.com/odysseythink/gofy/backend/entities/app/queue"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
)

type WorkflowCallback interface {
	OnWorkflowRunStarted()
	OnWorkflowRunSucceeded()
	OnWorkflowRunFailed(err string)
	OnWorkflowNodeExecuteStarted(node_id string, node_type nodesenumtypes.NodeType, title string, desc string, node_run_index int, predecessor_node_id string)
	OnWorkflowNodeExecuteSucceeded(node_id string, node_type nodesenumtypes.NodeType, title string, desc string, inputs map[string]any,
		process_data map[string]any,
		outputs map[string]any,
		execution_metadata map[string]any)
	OnWorkflowNodeExecuteFailed(node_id string, node_type nodesenumtypes.NodeType, title string, desc string,
		err string,
		inputs map[string]any,
		outputs map[string]any,
		process_data map[string]any)
	OnNodeTextChunk(node_id string, text string, metadata map[string]any)
	OnWorkflowIterationStarted(node_id string, node_type nodesenumtypes.NodeType,
		node_run_index int,
		title string, desc string,
		inputs map[string]any,
		predecessor_node_id string,
		metadata map[string]any)
	OnWorkflowIterationNext(node_id string, node_type nodesenumtypes.NodeType,
		index int,
		node_run_index int,
		output any)
	OnWorkflowIterationCompleted(node_id string, node_type nodesenumtypes.NodeType,
		node_run_index int,
		outputs map[string]any)
	OnEvent(event appqueueentities.AppQueueEventer)
}
