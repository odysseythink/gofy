package end

import (
	"iter"
	"slices"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/workflow/graph"
	streamprocessor "mlib.com/gofy/server/core/workflow/stream_processor"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	endnodesentities "mlib.com/gofy/server/entities/nodes/end"
	workflowentities "mlib.com/gofy/server/entities/workflow"
)

// EndStreamProcessor implements StreamProcessor
type EndStreamProcessor struct {
	*streamprocessor.BaseStreamProcessor
	endStreamParam                      *endnodesentities.EndStreamParam
	routePosition                       map[string]int
	currentStreamChunkGeneratingNodeIDs map[string][]string
	hasOutput                           bool
	outputNodeIDs                       map[string]struct{}
}

// NewEndStreamProcessor creates a new EndStreamProcessor
func NewEndStreamProcessor(graph *graph.Graph, variablePool *workflowentities.VariablePool) *EndStreamProcessor {
	esp := &EndStreamProcessor{
		BaseStreamProcessor: &streamprocessor.BaseStreamProcessor{
			Graph:        graph,
			VariablePool: variablePool,
		},

		endStreamParam:                      graph.EndStreamParam,
		routePosition:                       make(map[string]int),
		currentStreamChunkGeneratingNodeIDs: make(map[string][]string),
		hasOutput:                           false,
		outputNodeIDs:                       make(map[string]struct{}),
	}

	// Initialize routePosition
	for endNodeID := range esp.endStreamParam.EndStreamVariableSelectorMapping {
		esp.routePosition[endNodeID] = 0
	}

	esp.RestNodeIDs = make([]string, len(graph.NodeIDs))
	copy(esp.RestNodeIDs, graph.NodeIDs)

	return esp
}
func (esp *EndStreamProcessor) Process(generator iter.Seq[graphengineentities.GraphEngineEvent]) iter.Seq[graphengineentities.GraphEngineEvent] {
	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		for event := range generator {
			if event != nil {
				if e, ok := event.(*graphengineentities.NodeRunStartedEvent); ok {
					if e.RouteNodeState.NodeID == esp.Graph.RootNodeID && len(esp.RestNodeIDs) > 0 {
						esp.Reset()
					}
					if !yield(event) {
						return
					}
				} else if e, ok := event.(*graphengineentities.NodeRunStreamChunkEvent); ok {
					if e.InIterationID != "" {
						if esp.hasOutput {
							if _, ok := esp.outputNodeIDs[e.NodeID]; !ok {
								e.ChunkContent = "\n" + e.ChunkContent
							}
						}
						esp.outputNodeIDs[e.NodeID] = struct{}{}
						esp.hasOutput = true
						if !yield(event) {
							return
						}
						continue
					}
					var stream_out_end_node_ids []string
					if _, ok := esp.currentStreamChunkGeneratingNodeIDs[e.RouteNodeState.NodeID]; ok {
						stream_out_end_node_ids = esp.currentStreamChunkGeneratingNodeIDs[e.RouteNodeState.NodeID]
					} else {
						stream_out_end_node_ids = esp.getStreamOutEndNodeIDs(e)
						esp.currentStreamChunkGeneratingNodeIDs[e.RouteNodeState.NodeID] = stream_out_end_node_ids
					}
					if len(stream_out_end_node_ids) > 0 {
						if _, ok := esp.outputNodeIDs[e.NodeID]; !ok && esp.hasOutput {
							e.ChunkContent = "\n" + e.ChunkContent
						}
						esp.outputNodeIDs[e.NodeID] = struct{}{}
						esp.hasOutput = true
						if !yield(event) {
							return
						}
					}
				} else if e, ok := event.(*graphengineentities.NodeRunSucceededEvent); ok {
					if !yield(event) {
						return
					}
					if _, ok := esp.currentStreamChunkGeneratingNodeIDs[e.RouteNodeState.NodeID]; ok {
						// update esp.route_position after all stream event finished
						for _, end_node_id := range esp.currentStreamChunkGeneratingNodeIDs[e.RouteNodeState.NodeID] {
							esp.routePosition[end_node_id] += 1
						}
						delete(esp.currentStreamChunkGeneratingNodeIDs, e.RouteNodeState.NodeID)
					}
					// remove unreachable nodes
					esp.RemoveUnreachableNodes(e)

					// generate stream outputs
					for se := range esp.GenerateStreamOutputsWhenNodeFinished(e) {
						if !yield(se) {
							return
						}
					}
				} else {
					if !yield(event) {
						return
					}
				}
			}
		}
		mlog.Debugf("------end stream processor return")
	}
}

// Reset resets the processor
func (esp *EndStreamProcessor) Reset() {
	esp.routePosition = make(map[string]int)
	for endNodeID := range esp.endStreamParam.EndStreamVariableSelectorMapping {
		esp.routePosition[endNodeID] = 0
	}
	esp.RestNodeIDs = slices.Clone(esp.Graph.NodeIDs)
	esp.currentStreamChunkGeneratingNodeIDs = make(map[string][]string)
}

// GenerateStreamOutputsWhenNodeFinished generates stream outputs when a node finishes
func (esp *EndStreamProcessor) GenerateStreamOutputsWhenNodeFinished(event *graphengineentities.NodeRunSucceededEvent) iter.Seq[graphengineentities.GraphEngineEvent] {
	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		for end_node_id, position := range esp.routePosition {
			// all depends on end node id not in rest node ids
			if event.RouteNodeState.NodeID != end_node_id {
				not_all_in := true
				if _, ok := esp.endStreamParam.EndDependencies[end_node_id]; ok {
					for _, dep_id := range esp.endStreamParam.EndDependencies[end_node_id] {
						if slices.Contains(esp.RestNodeIDs, dep_id) {
							not_all_in = false
							break
						}
					}
				}
				if !slices.Contains(esp.RestNodeIDs, end_node_id) || not_all_in {
					continue
				}
			}

			route_position := esp.routePosition[end_node_id]

			position = 0
			value_selectors := [][]string{}
			for _, current_value_selectors := range esp.endStreamParam.EndStreamVariableSelectorMapping[end_node_id] {
				if position >= route_position {
					value_selectors = append(value_selectors, current_value_selectors)
				}
				position += 1
			}
			for _, value_selector := range value_selectors {
				if len(value_selector) == 0 {
					continue
				}
				value := esp.VariablePool.Get(value_selector)

				if value == nil {
					break
				}
				text := value.Markdown()

				if text != "" {
					current_node_id := value_selector[0]
					if _, ok := esp.outputNodeIDs[current_node_id]; !ok && esp.hasOutput {
						text = "\n" + text
					}
					esp.outputNodeIDs[current_node_id] = struct{}{}
					esp.hasOutput = true
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
				esp.routePosition[end_node_id] += 1
			}
		}
	}
}

func (esp *EndStreamProcessor) getStreamOutEndNodeIDs(event *graphengineentities.NodeRunStreamChunkEvent) []string {
	/*
	   Is stream out support
	   :param event: queue text chunk event
	   :return:
	*/
	stream_output_value_selector := event.FromVariableSelector
	if len(stream_output_value_selector) == 0 {
		return []string{}
	}
	stream_out_end_node_ids := []string{}
	for end_node_id, route_position := range esp.routePosition {
		if !slices.Contains(esp.RestNodeIDs, end_node_id) {
			continue
		}
		not_all_in := true
		if _, ok := esp.endStreamParam.EndDependencies[end_node_id]; ok {
			for _, dep_id := range esp.endStreamParam.EndDependencies[end_node_id] {
				if slices.Contains(esp.RestNodeIDs, dep_id) {
					not_all_in = false
					break
				}
			}
		}
		// all depends on end node id not in rest node ids
		if not_all_in {
			if route_position >= len(esp.endStreamParam.EndStreamVariableSelectorMapping[end_node_id]) {
				continue
			}
			position := 0
			var value_selector []string
			for _, current_value_selectors := range esp.endStreamParam.EndStreamVariableSelectorMapping[end_node_id] {
				if position == route_position {
					value_selector = current_value_selectors
					break
				}
				position += 1
			}
			if value_selector == nil {
				continue
			}
			// check chunk node id is before current node id or equal to current node id
			if slices.Compare(value_selector, stream_output_value_selector) != 0 {
				continue
			}
			stream_out_end_node_ids = append(stream_out_end_node_ids, end_node_id)
		}
	}
	return stream_out_end_node_ids
}
