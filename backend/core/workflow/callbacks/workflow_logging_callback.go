package callbacks

import (
	"fmt"

	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
)

var (
	TEXT_COLOR_MAPPING = map[string]string{
		"blue":   "36;1",
		"yellow": "33;1",
		"pink":   "38;5;200",
		"green":  "32;1",
		"red":    "31;1",
	}
)

type WorkflowLoggingCallback struct {
	currentNodeID string
}

func (cb *WorkflowLoggingCallback) OnEvent(event graphengineentities.GraphEngineEvent) {
	if _, ok := any(event).(*graphengineentities.GraphRunStartedEvent); ok {
		cb.PrintText("\n[GraphRunStartedEvent]", "pink", "")
	} else if _, ok := any(event).(*graphengineentities.GraphRunSucceededEvent); ok {
		cb.PrintText("\n[GraphRunSucceededEvent]", "green", "")
	} else if _, ok := any(event).(*graphengineentities.GraphRunPartialSucceededEvent); ok {
		cb.PrintText("\n[GraphRunPartialSucceededEvent]", "pink", "")
	} else if ev, ok := any(event).(*graphengineentities.GraphRunFailedEvent); ok {
		cb.PrintText(fmt.Sprintf("\n[GraphRunFailedEvent] reason: %s", ev.Error), "red", "")
	} else if ev, ok := any(event).(*graphengineentities.NodeRunStartedEvent); ok {
		cb.onWorkflowNodeExecuteStarted(ev)
	} else if ev, ok := any(event).(*graphengineentities.NodeRunSucceededEvent); ok {
		cb.onWorkflowNodeExecuteSucceeded(ev)
	} else if ev, ok := any(event).(*graphengineentities.NodeRunFailedEvent); ok {
		cb.onWorkflowNodeExecuteFailed(ev)
	} else if ev, ok := any(event).(*graphengineentities.NodeRunStreamChunkEvent); ok {
		cb.onNodeTextChunk(ev)
	} else if ev, ok := any(event).(*graphengineentities.ParallelBranchRunStartedEvent); ok {
		cb.onWorkflowParallelStarted(ev)
	} else if ev, ok := any(event).(*graphengineentities.ParallelBranchRunSucceededEvent); ok {
		cb.onWorkflowParallelCompleted(ev)
	} else if ev, ok := any(event).(*graphengineentities.ParallelBranchRunFailedEvent); ok {
		cb.onWorkflowParallelCompleted(ev)
	} else if ev, ok := any(event).(*graphengineentities.IterationRunStartedEvent); ok {
		cb.onWorkflowIterationStarted(ev)
	} else if ev, ok := any(event).(*graphengineentities.IterationRunNextEvent); ok {
		cb.onWorkflowIterationNext(ev)
	} else if ev, ok := any(event).(*graphengineentities.IterationRunSucceededEvent); ok {
		cb.onWorkflowIterationCompleted(ev)
	} else if ev, ok := any(event).(*graphengineentities.IterationRunFailedEvent); ok {
		cb.onWorkflowIterationCompleted(ev)
	} else {
		cb.PrintText(fmt.Sprintf("\n[%#v]", event), "blue", "")
	}
}

func (cb *WorkflowLoggingCallback) onWorkflowNodeExecuteStarted(event *graphengineentities.NodeRunStartedEvent) {
	/*
		Workflow node execute started
	*/
	cb.PrintText("\n[NodeRunStartedEvent]", "yellow", "")
	cb.PrintText(fmt.Sprintf("Node ID: %s", event.NodeID), "yellow", "")
	cb.PrintText(fmt.Sprintf("Node Title: %s", event.NodeData.Title), "yellow", "")
	cb.PrintText(fmt.Sprintf("Type: %s", event.NodeType), "yellow", "")

}
func (cb *WorkflowLoggingCallback) onWorkflowNodeExecuteSucceeded(event *graphengineentities.NodeRunSucceededEvent) {
	/*
		Workflow node execute succeeded
	*/
	route_node_state := event.RouteNodeState

	cb.PrintText("\n[NodeRunSucceededEvent]", "green", "")
	cb.PrintText(fmt.Sprintf("Node ID: %s", event.NodeID), "green", "")
	cb.PrintText(fmt.Sprintf("Node Title: %s", event.NodeData.Title), "green", "")
	cb.PrintText(fmt.Sprintf("Type: %s", event.NodeType), "green", "")

	if route_node_state.NodeRunResult != nil {
		node_run_result := route_node_state.NodeRunResult
		cb.PrintText(
			fmt.Sprintf("Inputs: %#v", node_run_result.Inputs),
			"green",
			"",
		)

		cb.PrintText(
			fmt.Sprintf("Process Data: %#v", node_run_result.ProcessData),
			"green",
			"",
		)

		cb.PrintText(
			fmt.Sprintf("Outputs: %#v", node_run_result.Outputs),
			"green",
			"",
		)
		cb.PrintText(
			fmt.Sprintf("Metadata: %#v", node_run_result.Metadata),
			"green",
			"",
		)
	}
}
func (cb *WorkflowLoggingCallback) onWorkflowNodeExecuteFailed(event *graphengineentities.NodeRunFailedEvent) {
	/*
		Workflow node execute failed
	*/
	route_node_state := event.RouteNodeState

	cb.PrintText("\n[NodeRunFailedEvent]", "red", "")
	cb.PrintText(fmt.Sprintf("Node ID: %s", event.NodeID), "red", "")
	cb.PrintText(fmt.Sprintf("Node Title: %s", event.NodeData.Title), "red", "")
	cb.PrintText(fmt.Sprintf("Type: %s", event.NodeType), "red", "")

	if route_node_state.NodeRunResult != nil {
		node_run_result := route_node_state.NodeRunResult
		cb.PrintText(fmt.Sprintf("Error: %v", node_run_result.Error), "red", "")
		cb.PrintText(
			fmt.Sprintf("Inputs: %#v", node_run_result.Inputs),
			"red",
			"",
		)
		cb.PrintText(
			fmt.Sprintf("Process Data: %#v", node_run_result.ProcessData),
			"red",
			"",
		)
		cb.PrintText(
			fmt.Sprintf("Outputs: %#v", node_run_result.Outputs),
			"red",
			"",
		)
	}
}
func (cb *WorkflowLoggingCallback) onNodeTextChunk(event *graphengineentities.NodeRunStreamChunkEvent) {
	/*
		Publish text chunk
	*/
	route_node_state := event.RouteNodeState
	if cb.currentNodeID == "" || cb.currentNodeID != route_node_state.NodeID {
		cb.currentNodeID = route_node_state.NodeID
		cb.PrintText("\n[NodeRunStreamChunkEvent]", "", "")
		cb.PrintText(fmt.Sprintf("Node ID: %s", route_node_state.NodeID), "", "")

		node_run_result := route_node_state.NodeRunResult
		if node_run_result != nil {
			cb.PrintText(
				fmt.Sprintf("Metadata: %#v", node_run_result.Metadata),
				"",
				"",
			)
		}
	}
	cb.PrintText(event.ChunkContent, "pink", "")

}
func (cb *WorkflowLoggingCallback) onWorkflowParallelStarted(event *graphengineentities.ParallelBranchRunStartedEvent) {
	/*
		Publish parallel started
	*/
	cb.PrintText("\n[ParallelBranchRunStartedEvent]", "blue", "")
	cb.PrintText(fmt.Sprintf("Parallel ID: %s", event.ParallelID), "blue", "")
	cb.PrintText(fmt.Sprintf("Branch ID: %s", event.ParallelStartNodeID), "blue", "")
	if event.InIterationID != "" {
		cb.PrintText(fmt.Sprintf("Iteration ID: %s", event.InIterationID), "blue", "")
	}
}
func (cb *WorkflowLoggingCallback) onWorkflowParallelCompleted(
	event graphengineentities.GraphEngineEvent, /*(*graphengineentities.ParallelBranchRunSucceededEvent | ParallelBranchRunFailedEvent*/
) {
	/*
		Publish parallel completed
	*/
	color := ""
	if ev, ok := any(event).(*graphengineentities.ParallelBranchRunSucceededEvent); ok {
		color = "blue"
		cb.PrintText("\n[ParallelBranchRunSucceededEvent]", color, "")

		cb.PrintText(fmt.Sprintf("Parallel ID: %s", ev.ParallelID), color, "")
		cb.PrintText(fmt.Sprintf("Branch ID: %s", ev.ParallelStartNodeID), color, "")
		if ev.InIterationID != "" {
			cb.PrintText(fmt.Sprintf("Iteration ID: %s", ev.InIterationID), "blue", "")
		}
	} else if ev, ok := any(event).(*graphengineentities.ParallelBranchRunFailedEvent); ok {
		color = "red"
		cb.PrintText("\n[ParallelBranchRunFailedEvent]", color, "")

		cb.PrintText(fmt.Sprintf("Parallel ID: %s", ev.ParallelID), color, "")
		cb.PrintText(fmt.Sprintf("Branch ID: %s", ev.ParallelStartNodeID), color, "")
		if ev.InIterationID != "" {
			cb.PrintText(fmt.Sprintf("Iteration ID: %s", ev.InIterationID), "blue", "")
		}
		cb.PrintText(fmt.Sprintf("Error: %s", ev.Error), color, "")
	}
}
func (cb *WorkflowLoggingCallback) onWorkflowIterationStarted(event *graphengineentities.IterationRunStartedEvent) {
	/*
		Publish iteration started
	*/
	cb.PrintText("\n[IterationRunStartedEvent]", "blue", "")
	cb.PrintText(fmt.Sprintf("Iteration Node ID: %s", event.IterationID), "blue", "")

}
func (cb *WorkflowLoggingCallback) onWorkflowIterationNext(event *graphengineentities.IterationRunNextEvent) {
	/*
		Publish iteration next
	*/
	cb.PrintText("\n[IterationRunNextEvent]", "blue", "")
	cb.PrintText(fmt.Sprintf("Iteration Node ID: %s", event.IterationID), "blue", "")
	cb.PrintText(fmt.Sprintf("Iteration Index: %v", event.Index), "blue", "")

}
func (cb *WorkflowLoggingCallback) onWorkflowIterationCompleted(event graphengineentities.GraphEngineEvent /* *graphengineentities.IterationRunSucceededEvent | IterationRunFailedEvent*/) {
	/*
		Publish iteration completed
	*/
	if ev, ok := any(event).(*graphengineentities.IterationRunSucceededEvent); ok {
		cb.PrintText("\n[IterationRunSucceededEvent]", "blue", "")
		cb.PrintText(fmt.Sprintf("Node ID: %s", ev.IterationID), "blue", "")
	} else if ev, ok := any(event).(*graphengineentities.IterationRunFailedEvent); ok {
		cb.PrintText("\n[IterationRunFailedEvent]", "blue", "")
		cb.PrintText(fmt.Sprintf("Node ID: %s", ev.IterationID), "blue", "")
	}
}

func (cb *WorkflowLoggingCallback) PrintText(text string, color string, end string /*= "\n"*/) {
	/*Print text with highlighting and no end characters.*/
	if end == "" {
		end = "\n"
	}
	text_to_print := ""
	if color != "" {
		text_to_print = cb.getColoredText(text, color)
	} else {
		text_to_print = text
	}
	fmt.Println(text_to_print, end)
}

func (cb *WorkflowLoggingCallback) getColoredText(text string, color string) string {
	/*Get colored text.*/
	color_str := TEXT_COLOR_MAPPING[color]
	return fmt.Sprintf("\u001b[%sm\033[1;3m%s\u001b[0m", color_str, text)
}
