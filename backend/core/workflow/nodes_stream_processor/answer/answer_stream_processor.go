package answer

import (
	"fmt"
	"iter"
	"slices"

	"github.com/odysseythink/gofy/backend/core/workflow/graph"
	streamprocessor "github.com/odysseythink/gofy/backend/core/workflow/stream_processor"
	"github.com/odysseythink/gofy/backend/core/workflow/utils/condition"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
	answernodesentities "github.com/odysseythink/gofy/backend/entities/nodes/answer"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	answernodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes/answer"
)

type AnswerStreamProcessor struct {
	*streamprocessor.BaseStreamProcessor
	generateRoutes                      *answernodesentities.AnswerStreamGenerateRoute
	routePosition                       map[string]int
	currentStreamChunkGeneratingNodeIDs map[string][]string
}

// NewEndStreamProcessor creates a new EndStreamProcessor
func NewAnswerStreamProcessor(graph *graph.Graph, variablePool *workflowentities.VariablePool) *AnswerStreamProcessor {
	esp := &AnswerStreamProcessor{
		BaseStreamProcessor: &streamprocessor.BaseStreamProcessor{
			Graph:        graph,
			VariablePool: variablePool,
		},

		generateRoutes:                      graph.AnswerStreamGenerateRoutes,
		routePosition:                       make(map[string]int),
		currentStreamChunkGeneratingNodeIDs: make(map[string][]string),
	}

	// Initialize routePosition
	for answer_node_id := range esp.generateRoutes.AnswerGenerateRoute {
		esp.routePosition[answer_node_id] = 0
	}

	return esp
}

// Reset resets the processor
func (esp *AnswerStreamProcessor) Reset() {
	esp.routePosition = make(map[string]int)
	for answer_node_id := range esp.generateRoutes.AnswerGenerateRoute {
		esp.routePosition[answer_node_id] = 0
	}
	esp.RestNodeIDs = slices.Clone(esp.Graph.NodeIDs)
	esp.currentStreamChunkGeneratingNodeIDs = make(map[string][]string)
}

// GenerateStreamOutputsWhenNodeFinished generates stream outputs when a node finishes
func (esp *AnswerStreamProcessor) GenerateStreamOutputsWhenNodeFinished(event *graphengineentities.NodeRunSucceededEvent) iter.Seq[graphengineentities.GraphEngineEvent] {
	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		for answer_node_id := range esp.routePosition {
			// all depends on end node id not in rest node ids
			if event.RouteNodeState.NodeID != answer_node_id {
				not_all_in := true
				for _, dep_id := range esp.generateRoutes.AnswerDependencies[answer_node_id] {
					if slices.Contains(esp.RestNodeIDs, dep_id) {
						not_all_in = false
						break
					}
				}

				if !slices.Contains(esp.RestNodeIDs, answer_node_id) || not_all_in {
					continue
				}
			}

			route_position := esp.routePosition[answer_node_id]
			route_chunks := esp.generateRoutes.AnswerGenerateRoute[answer_node_id][route_position:]

			for _, route_chunk := range route_chunks {
				if route_chunk.Type() == answernodesenumtypes.GenerateRouteChunk_TEXT {
					new_route_chunk := any(route_chunk).(*answernodesentities.TextGenerateRouteChunk)
					if !yield(&graphengineentities.NodeRunStreamChunkEvent{
						BaseNodeEvent: &graphengineentities.BaseNodeEvent{
							ID:       event.ID,
							NodeID:   event.NodeID,
							NodeType: event.NodeType,
							NodeData: event.NodeData,

							RouteNodeState:      event.RouteNodeState,
							ParallelID:          event.ParallelID,
							ParallelStartNodeID: event.ParallelStartNodeID,
						},
						ChunkContent:         new_route_chunk.Text,
						FromVariableSelector: []string{answer_node_id, "answer"},
					}) {
						return
					}
				} else {
					new_route_chunk := any(route_chunk).(*answernodesentities.VarGenerateRouteChunk)
					value_selector := new_route_chunk.ValueSelector
					if len(value_selector) == 0 {
						break
					}
					value := esp.VariablePool.Get(value_selector)

					if value == nil {
						break
					}

					text := value.Markdown()

					if text != "" {
						if !yield(&graphengineentities.NodeRunStreamChunkEvent{
							BaseNodeEvent: &graphengineentities.BaseNodeEvent{
								ID:       event.ID,
								NodeID:   event.NodeID,
								NodeType: event.NodeType,
								NodeData: event.NodeData,

								RouteNodeState:      event.RouteNodeState,
								ParallelID:          event.ParallelID,
								ParallelStartNodeID: event.ParallelStartNodeID,
							},
							ChunkContent:         text,
							FromVariableSelector: value_selector,
						}) {
							return
						}
					}
				}
				esp.routePosition[answer_node_id] += 1
			}
		}
	}
}

func (esp *AnswerStreamProcessor) getStreamOutAnswerNodeIDs(event *graphengineentities.NodeRunStreamChunkEvent) []string {
	/*
	   Is stream out support
	   :param event: queue text chunk event
	   :return:
	*/
	if len(event.FromVariableSelector) == 0 {
		return []string{}
	}
	stream_output_value_selector := event.FromVariableSelector

	stream_out_answer_node_ids := []string{}
	for answer_node_id, route_position := range esp.routePosition {
		if !slices.Contains(esp.RestNodeIDs, answer_node_id) {
			continue
		}
		// all depends on answer node id not in rest node ids
		all_not_ins := []bool{}
		for _, dep_id := range esp.generateRoutes.AnswerDependencies[answer_node_id] {
			all_not_ins = append(all_not_ins, !slices.Contains(esp.RestNodeIDs, dep_id))
		}
		if condition.AllTrue(all_not_ins) {
			if route_position >= len(esp.generateRoutes.AnswerGenerateRoute[answer_node_id]) {
				continue
			}
			route_chunk := esp.generateRoutes.AnswerGenerateRoute[answer_node_id][route_position]

			if route_chunk.Type() != answernodesenumtypes.GenerateRouteChunk_VAR {
				continue
			}
			real_route_chunk := any(route_chunk).(*answernodesentities.VarGenerateRouteChunk)
			value_selector := real_route_chunk.ValueSelector

			// check chunk node id is before current node id or equal to current node id

			if slices.Compare(value_selector, stream_output_value_selector) != 0 {
				continue
			}
			stream_out_answer_node_ids = append(stream_out_answer_node_ids, answer_node_id)
		}
	}
	return stream_out_answer_node_ids
}

func (esp *AnswerStreamProcessor) Process(generator iter.Seq[graphengineentities.GraphEngineEvent]) iter.Seq[graphengineentities.GraphEngineEvent] {
	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		for event := range generator {
			if ev, ok := any(event).(*graphengineentities.NodeRunStartedEvent); ok {
				if ev.RouteNodeState.NodeID == esp.Graph.RootNodeID && len(esp.RestNodeIDs) > 0 {
					esp.Reset()
				}
				if !yield(ev) {
					return
				}
			} else if ev, ok := any(event).(*graphengineentities.NodeRunStreamChunkEvent); ok {
				if ev.InIterationID != "" {
					if !yield(ev) {
						return
					}
					continue
				}
				var stream_out_answer_node_ids []string
				if _, ok := esp.currentStreamChunkGeneratingNodeIDs[ev.RouteNodeState.NodeID]; ok {
					stream_out_answer_node_ids = esp.currentStreamChunkGeneratingNodeIDs[ev.RouteNodeState.NodeID]
				} else {
					stream_out_answer_node_ids = esp.getStreamOutAnswerNodeIDs(ev)
					esp.currentStreamChunkGeneratingNodeIDs[ev.RouteNodeState.NodeID] = stream_out_answer_node_ids
				}
				for _, v := range stream_out_answer_node_ids {
					fmt.Println(v)
					if !yield(ev) {
						return
					}
				}
			} else if ev, ok := any(event).(*graphengineentities.NodeRunSucceededEvent); ok {
				if !yield(ev) {
					return
				}
				if _, ok := esp.currentStreamChunkGeneratingNodeIDs[ev.RouteNodeState.NodeID]; ok {
					// update esp.route_position after all stream event finished
					for _, answer_node_id := range esp.currentStreamChunkGeneratingNodeIDs[ev.RouteNodeState.NodeID] {
						esp.routePosition[answer_node_id] += 1
					}
					delete(esp.currentStreamChunkGeneratingNodeIDs, ev.RouteNodeState.NodeID)
				}
				esp.RemoveUnreachableNodes(ev)

				// generate stream outputs
				for item := range esp.GenerateStreamOutputsWhenNodeFinished(ev) {
					if !yield(item) {
						return
					}
				}
			} else if ev, ok := any(event).(*graphengineentities.NodeRunExceptionEvent); ok {
				if !yield(ev) {
					return
				}
				if _, ok := esp.currentStreamChunkGeneratingNodeIDs[ev.RouteNodeState.NodeID]; ok {
					// update esp.route_position after all stream event finished
					for _, answer_node_id := range esp.currentStreamChunkGeneratingNodeIDs[ev.RouteNodeState.NodeID] {
						esp.routePosition[answer_node_id] += 1
					}
					delete(esp.currentStreamChunkGeneratingNodeIDs, ev.RouteNodeState.NodeID)
				}
				esp.RemoveUnreachableNodes(ev)

				// generate stream outputs
				for item := range esp.GenerateStreamOutputsWhenNodeFinished(&graphengineentities.NodeRunSucceededEvent{
					BaseNodeEvent: ev.BaseNodeEvent,
				}) {
					if !yield(item) {
						return
					}
				}
			} else {
				if !yield(ev) {
					return
				}
			}
		}
	}
}
