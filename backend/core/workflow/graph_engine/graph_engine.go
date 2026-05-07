package graphengine

import (
	"fmt"
	"iter"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/core/workflow/graph"
	conditionhandlers "github.com/odysseythink/gofy/backend/core/workflow/graph_engine/condition_handlers"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	answergeneraterouter "github.com/odysseythink/gofy/backend/core/workflow/nodes_stream_processor/answer"
	endgeneraterouter "github.com/odysseythink/gofy/backend/core/workflow/nodes_stream_processor/end"
	streamprocessor "github.com/odysseythink/gofy/backend/core/workflow/stream_processor"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
	evententities "github.com/odysseythink/gofy/backend/entities/nodes/event"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	graphengineenumtypes "github.com/odysseythink/gofy/backend/enum_types/graph_engine"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	workflowenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/utils"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

type GraphEngine struct {
	threadPoolID      string
	isMainThreadPool  bool
	graph             *graph.Graph
	initParams        *graphengineentities.GraphInitParams
	GraphRuntimeState *graphengineentities.GraphRuntimeState
	maxExecutionSteps int
	maxExecutionTime  int
}

func NewGraphEngine(
	tenantID string,
	appID string,
	workflowType models.WorkflowType,
	workflowID string,
	userID string,
	userFrom models.UserFrom,
	invokeFrom appenumtypes.InvokeFrom,
	callDepth int,
	graph *graph.Graph,
	graphConfig map[string]interface{},
	variablePool *workflowentities.VariablePool,
	maxExecutionSteps int,
	maxExecutionTime int,
) *GraphEngine {
	ge := &GraphEngine{
		threadPoolID:     uuid.NewV4().String(),
		isMainThreadPool: true,
	}
	ge.graph = graph
	ge.initParams = &graphengineentities.GraphInitParams{
		TenantID:     tenantID,
		AppID:        appID,
		WorkflowType: workflowType,
		WorkflowID:   workflowID,
		GraphConfig:  graphConfig,
		UserID:       userID,
		UserFrom:     userFrom,
		InvokeFrom:   invokeFrom,
		CallDepth:    callDepth,
	}
	ge.GraphRuntimeState = &graphengineentities.GraphRuntimeState{
		VariablePool: variablePool,
		StartAt:      time.Now(),
		NodeRunState: &graphengineentities.RuntimeRouteState{
			Routes:           make(map[string][]string),
			NodeStateMapping: make(map[string]*graphengineentities.RouteNodeState),
		},
	}
	ge.maxExecutionSteps = maxExecutionSteps
	ge.maxExecutionTime = maxExecutionTime

	return ge
}

func (ge *GraphEngine) Run() iter.Seq[graphengineentities.GraphEngineEvent] {
	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		handle_exceptions := []string{}
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(*exceptions.GraphRunFailedError); ok {
					yield(&graphengineentities.GraphRunFailedEvent{
						Error:           exp.Error(),
						ExceptionsCount: len(handle_exceptions),
					})
				} else if exp, ok := r.(error); ok {
					mlog.Error("Unknown Error when graph running")
					yield(&graphengineentities.GraphRunFailedEvent{
						Error:           exp.Error(),
						ExceptionsCount: len(handle_exceptions),
					})
					panic(exp)
				} else {
					panic(r)
				}
			}
		}()
		if !yield(&graphengineentities.GraphRunStartedEvent{}) {
			return
		}

		var stream_processor streamprocessor.StreamProcessor
		if ge.initParams.WorkflowType == models.Workflow_CHAT {
			stream_processor = answergeneraterouter.NewAnswerStreamProcessor(ge.graph, ge.GraphRuntimeState.VariablePool)
		} else {
			stream_processor = endgeneraterouter.NewEndStreamProcessor(ge.graph, ge.GraphRuntimeState.VariablePool)
		}
		generator := stream_processor.Process(ge._run(ge.graph.RootNodeID, "", "", "", handle_exceptions))
		for item := range generator {
			func() {
				defer func() {
					if r := recover(); r != nil {
						mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
						if exp, ok := r.(error); ok {
							mlog.Errorf("Graph run failed:%#v", exp)
							yield(&graphengineentities.GraphRunFailedEvent{
								Error:           exp.Error(),
								ExceptionsCount: len(handle_exceptions),
							})
							return
						} else {
							panic(r)
						}
					}
				}()
				if !yield(item) {
					return
				}

				if ev, ok := any(item).(*graphengineentities.NodeRunFailedEvent); ok {
					gevent := &graphengineentities.GraphRunFailedEvent{
						Error:           ev.RouteNodeState.FailedReason,
						ExceptionsCount: len(handle_exceptions),
					}
					if gevent.Error == "" {
						gevent.Error = "Unknown error."
					}

					yield(gevent)
					return
				} else if ev, ok := any(item).(*graphengineentities.NodeRunSucceededEvent); ok {
					if ev.NodeType == nodesenumtypes.Node_END {
						if ev.RouteNodeState.NodeRunResult != nil && ev.RouteNodeState.NodeRunResult.Outputs != nil {
							ge.GraphRuntimeState.Outputs = ev.RouteNodeState.NodeRunResult.Outputs
						} else {
							ge.GraphRuntimeState.Outputs = make(map[string]any)
						}
					} else if ev.NodeType == nodesenumtypes.Node_ANSWER {
						mlog.Debugf("------ge.GraphRuntimeState=%#v", ge.GraphRuntimeState)
						if ge.GraphRuntimeState.Outputs == nil {
							ge.GraphRuntimeState.Outputs = make(map[string]any)
						}
						if _, ok := ge.GraphRuntimeState.Outputs["answer"]; !ok {
							ge.GraphRuntimeState.Outputs["answer"] = ""
						}
						answer := ""
						if ev.RouteNodeState.NodeRunResult != nil && ev.RouteNodeState.NodeRunResult.Outputs != nil {

							if _, ok := ev.RouteNodeState.NodeRunResult.Outputs["answer"]; ok {
								if _, ok := ev.RouteNodeState.NodeRunResult.Outputs["answer"].(string); ok {
									answer = ev.RouteNodeState.NodeRunResult.Outputs["answer"].(string)
								}
							}

						}
						ge.GraphRuntimeState.Outputs["answer"] = ge.GraphRuntimeState.Outputs["answer"].(string) + "\n" + answer

						ge.GraphRuntimeState.Outputs["answer"] = strings.ReplaceAll(ge.GraphRuntimeState.Outputs["answer"].(string), " ", "")
					}
				}
			}()
		}
		mlog.Debugf("------graph engine Run return")
		// count exceptions to determine partial success
		if len(handle_exceptions) > 0 {
			yield(&graphengineentities.GraphRunPartialSucceededEvent{
				ExceptionsCount: len(handle_exceptions),
				Outputs:         ge.GraphRuntimeState.Outputs,
			})
		} else {
			// trigger graph run success event
			yield(&graphengineentities.GraphRunSucceededEvent{
				Outputs: ge.GraphRuntimeState.Outputs,
			})
		}
	}
}

func (ge *GraphEngine) _run_parallel_branches(
	edge_mappings []*graphengineentities.GraphEdge,
	in_parallel_id string,
	parallel_start_node_id string,
	handle_exceptions []string,
) iter.Seq[any] {
	return func(yield func(any) bool) {
		// if nodes has no run conditions, parallel run all nodes
		parallel_id := ""
		if _, ok := ge.graph.NodeParallelMapping[edge_mappings[0].TargetNodeID]; ok {
			parallel_id = ge.graph.NodeParallelMapping[edge_mappings[0].TargetNodeID]
		}
		if parallel_id == "" {
			node_id := edge_mappings[0].TargetNodeID
			var node_config map[string]any
			if _, ok := ge.graph.NodeIDConfigMapping[node_id]; ok {
				node_config = ge.graph.NodeIDConfigMapping[node_id]
			}
			if node_config == nil {
				mlog.Errorf("Node %s related parallel not found or incorrectly connected to multiple parallel branches.", node_id)
				panic(exceptions.NewGraphRunFailedError(fmt.Sprintf("Node %s related parallel not found or incorrectly connected to multiple parallel branches.", node_id)))
			}
			mlog.Debugf("------node_config=%#v", node_config)
			node_title := ""
			if _, ok := node_config["data"]; ok {
				if _, ok := node_config["data"].(map[string]any); ok {
					if _, ok := node_config["data"].(map[string]any)["title"]; ok {
						if _, ok := node_config["data"].(map[string]any)["title"].(string); ok {
							node_title = node_config["data"].(map[string]any)["title"].(string)
						}
					}
				}
			}
			mlog.Errorf("Node %s related parallel not found or incorrectly connected to multiple parallel branches.", node_title)
			panic(exceptions.NewGraphRunFailedError(fmt.Sprintf("Node %s related parallel not found or incorrectly connected to multiple parallel branches.", node_title)))
		}
		var parallel *graphengineentities.GraphParallel
		if _, ok := ge.graph.ParallelMapping[parallel_id]; ok {
			parallel = ge.graph.ParallelMapping[parallel_id]
		}

		if parallel == nil {
			mlog.Errorf("Parallel %s not found.", parallel_id)
			panic(exceptions.NewGraphRunFailedError(fmt.Sprintf("Parallel %s not found.", parallel_id)))
		}
		// run parallel nodes, run in new thread and use queue to get results
		q := make(chan any, 1024)

		// Create a list to store the threads
		// futures := []any{}
		var futures sync.WaitGroup

		// new thread
		parallel_num := 0
		for _, edge := range edge_mappings {
			if v, ok := ge.graph.NodeParallelMapping[edge.TargetNodeID]; !ok || v != parallel_id {
				continue
			}
			futures.Add(1)
			go func(target string) {
				parallel_num++
				ge._run_parallel_node(q, parallel_id, target, in_parallel_id, parallel_start_node_id, handle_exceptions)
				futures.Done()
			}(edge.TargetNodeID)

		}
		succeeded_count := 0
		for {
			timer := time.NewTimer(1 * time.Second)
			select {
			case event := <-q:
				if event == nil {
					goto END
				}
				if !yield(event) {
					return
				}
				if ev, ok := event.(*graphengineentities.ParallelBranchRunSucceededEvent); ok {
					if ev.ParallelID == parallel_id {
						succeeded_count += 1
						if succeeded_count == parallel_num {
							// q.put(None)
						}
						break
					}
				} else if ev, ok := event.(*graphengineentities.ParallelBranchRunFailedEvent); ok {
					panic(exceptions.NewGraphRunFailedError(ev.Error))
				}
			case <-timer.C:
				goto END
			}
		}
	END:
		mlog.Debugf("parallel event queue monitor exit")
		// wait all threads
		futures.Wait()
		mlog.Debugf("all parallel branch exit")

		// get final node id
		final_node_id := parallel.EndToNodeID
		if final_node_id != "" {
			if !yield(final_node_id) {
				return
			}
		}
	}
}

// // func create_copy(self){
// // 	/*
// // 	create a graph engine copy
// // 	:return: with a new variable pool instance of graph engine
// // 	*/
// // 	new_instance = copy(self)
// // 	new_instance.graph_runtime_state = copy(ge.graph_runtime_state)
// // 	new_instance.graph_runtime_state.variable_pool = deepcopy(ge.graph_runtime_state.variable_pool)
// // 	return new_instance
// // }

func (ge *GraphEngine) runNode(
	node_instance base.Noder,
	route_node_state *graphengineentities.RouteNodeState,
	parallel_id string,
	parallel_start_node_id string,
	parent_parallel_id string,
	parent_parallel_start_node_id string,
	handle_exceptions []string,
) iter.Seq[graphengineentities.GraphEngineEvent] {
	/*
		Run node
	*/
	// trigger node run start event
	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		if !yield(&graphengineentities.NodeRunStartedEvent{
			BaseNodeEvent: &graphengineentities.BaseNodeEvent{
				ID:                        node_instance.GetID(),
				NodeID:                    node_instance.GetNodeID(),
				NodeType:                  node_instance.Type(),
				NodeData:                  node_instance.GetBaseNodeData(),
				RouteNodeState:            route_node_state,
				ParallelID:                parallel_id,
				ParallelStartNodeID:       parallel_start_node_id,
				ParentParallelID:          parent_parallel_id,
				ParentParallelStartNodeID: parent_parallel_start_node_id,
			},
			PredecessorNodeID: node_instance.GetPreviousNodeID(),
		}) {
			return
		}

		max_retries := node_instance.GetBaseNodeData().RetryConfig.MaxRetries
		retry_interval := node_instance.GetBaseNodeData().RetryConfig.RetryIntervalSeconds()
		retries := 0
		should_continue_retry := true
		for should_continue_retry && retries <= max_retries {
			func() {
				defer func() {
					if r := recover(); r != nil {
						if _, ok := r.(*exceptions.GenerateTaskStoppedError); ok {
							// trigger node run failed event
							route_node_state.Status = graphengineenumtypes.RouteNodeStateStatus_FAILED
							route_node_state.FailedReason = "Workflow stopped."
							yield(&graphengineentities.NodeRunFailedEvent{
								BaseNodeEvent: &graphengineentities.BaseNodeEvent{
									ID:                        node_instance.GetID(),
									NodeID:                    node_instance.GetNodeID(),
									NodeType:                  node_instance.Type(),
									NodeData:                  node_instance.GetBaseNodeData(),
									RouteNodeState:            route_node_state,
									ParallelID:                parallel_id,
									ParallelStartNodeID:       parallel_start_node_id,
									ParentParallelID:          parent_parallel_id,
									ParentParallelStartNodeID: parent_parallel_start_node_id,
								},
								Error: "Workflow stopped.",
							})
							return
						} else if exp, ok := r.(error); ok {
							mlog.Errorf("Node %s run failed", node_instance.GetBaseNodeData().Title)
							panic(exp)
						} else {
							panic(r)
						}
					}
				}()
				// run node
				retry_start_at := time.Now()
				generator := node_instance.RunIter(node_instance)
				for item := range generator {
					if reflect.TypeOf(item).Implements(reflect.TypeOf((*graphengineentities.GraphEngineEvent)(nil)).Elem()) {
						if ev, ok := item.(*graphengineentities.BaseIterationEvent); ok {
							// add parallel info to iteration event
							ev.ParallelID = parallel_id
							ev.ParallelStartNodeID = parallel_start_node_id
							ev.ParentParallelID = parent_parallel_id
							ev.ParentParallelStartNodeID = parent_parallel_start_node_id
							item = any(ev)
						}
						if !yield(item.(graphengineentities.GraphEngineEvent)) {
							return
						}
					} else {
						if ev, ok := item.(*evententities.RunCompletedEvent); ok {
							run_result := ev.RunResult
							mlog.Debugf("---------run_result=%#v", run_result)
							mlog.Debugf("---------ShouldRetry=%#v, retries=%d, max_retries=%d", node_instance.ShouldRetry(), retries, max_retries)
							if run_result.Status == models.WorkflowNodeExecutionStatus_FAILED {
								if retries == max_retries &&
									node_instance.Type() == nodesenumtypes.Node_HTTP_REQUEST &&
									run_result.Outputs != nil &&
									!node_instance.ShouldContinueOnError() {
									run_result.Status = models.WorkflowNodeExecutionStatus_SUCCEEDED
								}
								if node_instance.ShouldRetry() && retries < max_retries {
									retries += 1
									route_node_state.NodeRunResult = run_result
									e := &graphengineentities.NodeRunRetryEvent{
										NodeRunStartedEvent: &graphengineentities.NodeRunStartedEvent{
											BaseNodeEvent: &graphengineentities.BaseNodeEvent{
												ID:                        node_instance.GetID(),
												NodeID:                    node_instance.GetNodeID(),
												NodeType:                  node_instance.Type(),
												NodeData:                  node_instance.GetBaseNodeData(),
												RouteNodeState:            route_node_state,
												ParallelID:                parallel_id,
												ParallelStartNodeID:       parallel_start_node_id,
												ParentParallelID:          parent_parallel_id,
												ParentParallelStartNodeID: parent_parallel_start_node_id,
											},
											PredecessorNodeID: node_instance.GetPreviousNodeID(),
										},
										Error:      run_result.Error,
										RetryIndex: retries,
										StartAt:    retry_start_at,
									}
									if e.Error == "" {
										e.Error = "Unknown error"
									}
									if !yield(e) {
										return
									}
									time.Sleep(time.Duration(retry_interval) * time.Second)
									continue
								}
							}
							route_node_state.SetFinished(run_result)

							if run_result.Status == models.WorkflowNodeExecutionStatus_FAILED {
								mlog.Debug("--------ShouldContinueOnError=", node_instance.ShouldContinueOnError())
								if node_instance.ShouldContinueOnError() {
									// if run failed, handle error
									run_result = ge.handleContinueOnError(
										node_instance,
										ev.RunResult,
										ge.GraphRuntimeState.VariablePool,
										handle_exceptions,
									)
									route_node_state.NodeRunResult = run_result
									route_node_state.Status = graphengineenumtypes.RouteNodeStateStatus_EXCEPTION
									if len(run_result.Outputs) > 0 {
										for variable_key, variable_value := range run_result.Outputs {
											// append variables to variable pool recursively
											ge.appendVariablesRecursively(
												node_instance.GetNodeID(),
												[]string{variable_key},
												variable_value,
											)
										}
									}
									e := &graphengineentities.NodeRunExceptionEvent{
										BaseNodeEvent: &graphengineentities.BaseNodeEvent{
											ID:                        node_instance.GetID(),
											NodeID:                    node_instance.GetNodeID(),
											NodeType:                  node_instance.Type(),
											NodeData:                  node_instance.GetBaseNodeData(),
											RouteNodeState:            route_node_state,
											ParallelID:                parallel_id,
											ParallelStartNodeID:       parallel_start_node_id,
											ParentParallelID:          parent_parallel_id,
											ParentParallelStartNodeID: parent_parallel_start_node_id,
										},
										Error: run_result.Error,
									}
									if e.Error == "" {
										e.Error = "System Error"
									}
									if !yield(e) {
										return
									}
									should_continue_retry = false
								} else {
									e := &graphengineentities.NodeRunFailedEvent{
										BaseNodeEvent: &graphengineentities.BaseNodeEvent{
											ID:                        node_instance.GetID(),
											NodeID:                    node_instance.GetNodeID(),
											NodeType:                  node_instance.Type(),
											NodeData:                  node_instance.GetBaseNodeData(),
											RouteNodeState:            route_node_state,
											ParallelID:                parallel_id,
											ParallelStartNodeID:       parallel_start_node_id,
											ParentParallelID:          parent_parallel_id,
											ParentParallelStartNodeID: parent_parallel_start_node_id,
										},
										Error: route_node_state.FailedReason,
									}
									if e.Error == "" {
										e.Error = "Unknown error."
									}
									if !yield(e) {
										return
									}
								}
								should_continue_retry = false
							} else if run_result.Status == models.WorkflowNodeExecutionStatus_SUCCEEDED {
								if _, ok := ge.graph.EdgeMapping[node_instance.GetNodeID()]; ok && node_instance.ShouldContinueOnError() {
									run_result.EdgeSourceHandle = string(nodesenumtypes.FailBranchSourceHandle_SUCCESS)
								}

								if len(run_result.Metadata) > 0 {
									if _, ok := run_result.Metadata[workflowenumtypes.NodeRunMetadataKey_TOTAL_TOKENS]; ok {
										// plus state total_tokens
										ge.GraphRuntimeState.TotalTokens += run_result.Metadata[workflowenumtypes.NodeRunMetadataKey_TOTAL_TOKENS].(int) /* type: ignore[arg-type]*/
									}
								}
								if run_result.LLmUsage != nil {
									// use the latest usage
									if ge.GraphRuntimeState.LLMUsage == nil {
										ge.GraphRuntimeState.LLMUsage = run_result.LLmUsage
									} else {
										ge.GraphRuntimeState.LLMUsage.Plus(run_result.LLmUsage)
									}
								}
								// append node output variables to variable pool
								if len(run_result.Outputs) > 0 {
									for variable_key, variable_value := range run_result.Outputs {
										// append variables to variable pool recursively
										ge.appendVariablesRecursively(
											node_instance.GetNodeID(),
											[]string{variable_key},
											variable_value,
										)
									}
								}

								// When setting metadata, convert to dict first
								if run_result.Metadata == nil {
									run_result.Metadata = make(map[workflowenumtypes.NodeRunMetadataKey]any)
								}
								if parallel_id != "" && parallel_start_node_id != "" {
									metadata_dict := run_result.Metadata
									metadata_dict[workflowenumtypes.NodeRunMetadataKey_PARALLEL_ID] = parallel_id
									metadata_dict[workflowenumtypes.NodeRunMetadataKey_PARALLEL_START_NODE_ID] = parallel_start_node_id
									if parent_parallel_id != "" && parent_parallel_start_node_id != "" {
										metadata_dict[workflowenumtypes.NodeRunMetadataKey_PARENT_PARALLEL_ID] = parent_parallel_id
										metadata_dict[workflowenumtypes.NodeRunMetadataKey_PARENT_PARALLEL_START_NODE_ID] = parent_parallel_start_node_id
									}
									run_result.Metadata = metadata_dict
								}
								if !yield(&graphengineentities.NodeRunSucceededEvent{
									BaseNodeEvent: &graphengineentities.BaseNodeEvent{
										ID:                        node_instance.GetID(),
										NodeID:                    node_instance.GetNodeID(),
										NodeType:                  node_instance.Type(),
										NodeData:                  node_instance.GetBaseNodeData(),
										RouteNodeState:            route_node_state,
										ParallelID:                parallel_id,
										ParallelStartNodeID:       parallel_start_node_id,
										ParentParallelID:          parent_parallel_id,
										ParentParallelStartNodeID: parent_parallel_start_node_id,
									},
								}) {
									return
								}
								should_continue_retry = false
							}
							mlog.Debugf("---------ev=%#v return", ev)
							return
						} else if ev, ok := item.(*evententities.RunStreamChunkEvent); ok {
							if !yield(&graphengineentities.NodeRunStreamChunkEvent{
								BaseNodeEvent: &graphengineentities.BaseNodeEvent{
									ID:                        node_instance.GetID(),
									NodeID:                    node_instance.GetNodeID(),
									NodeType:                  node_instance.Type(),
									NodeData:                  node_instance.GetBaseNodeData(),
									RouteNodeState:            route_node_state,
									ParallelID:                parallel_id,
									ParallelStartNodeID:       parallel_start_node_id,
									ParentParallelID:          parent_parallel_id,
									ParentParallelStartNodeID: parent_parallel_start_node_id,
								},
								ChunkContent:         ev.ChunkContent,
								FromVariableSelector: ev.FromVariableSelector,
							}) {
								return
							}
						} else if ev, ok := item.(*evententities.RunRetrieverResourceEvent); ok {
							if !yield(&graphengineentities.NodeRunRetrieverResourceEvent{
								BaseNodeEvent: &graphengineentities.BaseNodeEvent{
									ID:                        node_instance.GetID(),
									NodeID:                    node_instance.GetNodeID(),
									NodeType:                  node_instance.Type(),
									NodeData:                  node_instance.GetBaseNodeData(),
									RouteNodeState:            route_node_state,
									ParallelID:                parallel_id,
									ParallelStartNodeID:       parallel_start_node_id,
									ParentParallelID:          parent_parallel_id,
									ParentParallelStartNodeID: parent_parallel_start_node_id,
								},
								RetrieverResources: ev.RetrieverResources,
								Context:            ev.Context,
							}) {
								return
							}
						}
					}
				}
			}()
		}
	}
}

func (ge *GraphEngine) _run_parallel_node(
	q chan any,
	parallel_id string,
	parallel_start_node_id string,
	parent_parallel_id string,
	parent_parallel_start_node_id string,
	handle_exceptions []string,
) {
	/*
		Run parallel nodes
	*/

	defer func() {
		if r := recover(); r != nil {
			if exp, ok := r.(*exceptions.GraphRunFailedError); ok {
				q <- &graphengineentities.ParallelBranchRunFailedEvent{
					BaseParallelBranchEvent: &graphengineentities.BaseParallelBranchEvent{
						ParallelID:                parallel_id,
						ParallelStartNodeID:       parallel_start_node_id,
						ParentParallelID:          parent_parallel_id,
						ParentParallelStartNodeID: parent_parallel_start_node_id,
					},
					Error: exp.Error(),
				}
			} else if exp, ok := r.(error); ok {
				mlog.Error("Unknown Error when generating in parallel")
				q <- &graphengineentities.ParallelBranchRunFailedEvent{
					BaseParallelBranchEvent: &graphengineentities.BaseParallelBranchEvent{
						ParallelID:                parallel_id,
						ParallelStartNodeID:       parallel_start_node_id,
						ParentParallelID:          parent_parallel_id,
						ParentParallelStartNodeID: parent_parallel_start_node_id,
					},
					Error: exp.Error(),
				}
			} else {
				panic(r)
			}
		}
	}()
	q <- &graphengineentities.ParallelBranchRunStartedEvent{
		BaseParallelBranchEvent: &graphengineentities.BaseParallelBranchEvent{
			ParallelID:                parallel_id,
			ParallelStartNodeID:       parallel_start_node_id,
			ParentParallelID:          parent_parallel_id,
			ParentParallelStartNodeID: parent_parallel_start_node_id,
		},
	}
	// run node
	generator := ge._run(
		parallel_start_node_id,
		parallel_id,
		parent_parallel_id,
		parent_parallel_start_node_id,
		handle_exceptions,
	)

	for item := range generator {
		q <- item
	}

	// trigger graph run success event
	q <- &graphengineentities.ParallelBranchRunSucceededEvent{
		BaseParallelBranchEvent: &graphengineentities.BaseParallelBranchEvent{
			ParallelID:                parallel_id,
			ParallelStartNodeID:       parallel_start_node_id,
			ParentParallelID:          parent_parallel_id,
			ParentParallelStartNodeID: parent_parallel_start_node_id,
		},
	}
}

func (ge *GraphEngine) _run(
	start_node_id string,
	in_parallel_id string,
	parent_parallel_id string,
	parent_parallel_start_node_id string,
	handle_exceptions []string,
) iter.Seq[graphengineentities.GraphEngineEvent] {
	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		parallel_start_node_id := ""
		if in_parallel_id != "" {
			parallel_start_node_id = start_node_id
		}
		next_node_id := start_node_id
		var previous_route_node_state *graphengineentities.RouteNodeState
		for {
			// max steps reached
			if ge.GraphRuntimeState.NodeRunSteps > ge.maxExecutionSteps {
				panic(exceptions.NewGraphRunFailedError(fmt.Sprintf("Max steps %d reached.", ge.maxExecutionSteps)))
			}
			// or max execution time reached
			if ge.isTimedOut(ge.GraphRuntimeState.StartAt, ge.maxExecutionTime) {
				panic(exceptions.NewGraphRunFailedError(fmt.Sprintf("Max execution time %ds reached.", ge.maxExecutionTime)))
			}
			// init route node state
			route_node_state := ge.GraphRuntimeState.NodeRunState.CreateNodeState(next_node_id)

			// get node config
			node_id := route_node_state.NodeID
			node_config := map[string]any{}
			if _, ok := ge.graph.NodeIDConfigMapping[node_id]; ok {
				node_config = ge.graph.NodeIDConfigMapping[node_id]
			}
			if len(node_config) == 0 {
				panic(exceptions.NewGraphRunFailedError(fmt.Sprintf("Node %s config not found.", node_id)))
			}
			// convert to specific node
			var node_type nodesenumtypes.NodeType
			if _, ok := node_config["data"]; ok {
				if _, ok := node_config["data"].(map[string]any); ok {
					if _, ok := node_config["data"].(map[string]any)["type"]; ok {
						if _, ok := node_config["data"].(map[string]any)["type"].(string); ok {
							node_type = nodesenumtypes.NodeType(node_config["data"].(map[string]any)["type"].(string))
						}
					}
				}
			}

			// node_version = node_config.get("data", {}).get("version", "1")
			// node_cls = NODE_TYPE_CLASSES_MAPPING[node_type][node_version]
			previous_node_id := ""
			if previous_route_node_state != nil {
				previous_node_id = previous_route_node_state.NodeID
			}

			// init workflow run state
			node_instance := nodes.NewNode( // type: ignore
				route_node_state.ID,
				node_config,
				ge.initParams,
				ge.graph,
				ge.GraphRuntimeState,
				previous_node_id,
				node_type,
			)
			func() {
				defer func() {
					if r := recover(); r != nil {
						if exp, ok := r.(error); ok {
							route_node_state.Status = graphengineenumtypes.RouteNodeStateStatus_FAILED
							route_node_state.FailedReason = exp.Error()
							event := &graphengineentities.NodeRunFailedEvent{
								BaseNodeEvent: &graphengineentities.BaseNodeEvent{
									// ID:                        node_instance.GetID(),
									NodeID:   next_node_id,
									NodeType: node_type,
									// NodeData:                  node_instance.GetBaseNodeData(),
									RouteNodeState:            route_node_state,
									ParallelID:                in_parallel_id,
									ParallelStartNodeID:       parallel_start_node_id,
									ParentParallelID:          parent_parallel_id,
									ParentParallelStartNodeID: parent_parallel_start_node_id,
								},
								Error: exp.Error(),
							}
							if node_instance != nil {
								event.ID = node_instance.GetID()
								event.NodeData = node_instance.GetBaseNodeData()
							} else {
								event.ID = route_node_state.ID
							}
							yield(event)
						}
						panic(r)
					}
				}()
				// run node
				generator := ge.runNode(
					node_instance,
					route_node_state,
					in_parallel_id,
					parallel_start_node_id,
					parent_parallel_id,
					parent_parallel_start_node_id,
					handle_exceptions,
				)

				for item := range generator {
					mlog.Debugf("----item=%#v", item)
					if ev, ok := item.(*graphengineentities.NodeRunStartedEvent); ok {
						ge.GraphRuntimeState.NodeRunSteps += 1
						ev.RouteNodeState.Index = ge.GraphRuntimeState.NodeRunSteps
						item = ev
					}
					if !yield(item) {
						return
					}
				}
				mlog.Debugf("----run next......")

				ge.GraphRuntimeState.NodeRunState.NodeStateMapping[route_node_state.ID] = route_node_state

				// append route
				if previous_route_node_state != nil {
					ge.GraphRuntimeState.NodeRunState.AddRoute(previous_route_node_state.ID, route_node_state.ID)
				}
			}()

			if _, ok := ge.graph.NodeIDConfigMapping[next_node_id]; ok {
				if _, ok := ge.graph.NodeIDConfigMapping[next_node_id]["data"]; ok {
					if _, ok := ge.graph.NodeIDConfigMapping[next_node_id]["data"].(map[string]any); ok {
						if _, ok := ge.graph.NodeIDConfigMapping[next_node_id]["data"].(map[string]any)["type"]; ok {
							if _, ok := ge.graph.NodeIDConfigMapping[next_node_id]["data"].(map[string]any)["type"].(string); ok {
								node_type = nodesenumtypes.NodeType(ge.graph.NodeIDConfigMapping[next_node_id]["data"].(map[string]any)["type"].(string))
							}
						}
					}
				}
			}
			mlog.Debugf("----next node type=%s ", node_type)
			if node_type == nodesenumtypes.Node_END {
				mlog.Debugf("----end return")
				break
			}
			previous_route_node_state = route_node_state

			// get next node ids
			var edge_mappings []*graphengineentities.GraphEdge
			if _, ok := ge.graph.EdgeMapping[next_node_id]; ok {
				edge_mappings = ge.graph.EdgeMapping[next_node_id]
			}

			if len(edge_mappings) == 0 {
				break
			}
			if len(edge_mappings) == 1 {
				edge := edge_mappings[0]
				if previous_route_node_state.Status == graphengineenumtypes.RouteNodeStateStatus_EXCEPTION &&
					node_instance.GetBaseNodeData().ErrorStrategy == nodesenumtypes.ErrorStrategy_FAIL_BRANCH &&
					edge.RunCondition == nil {
					break
				}
				if edge.RunCondition != nil {
					if !(&conditionhandlers.ConditionManager{}).GetConditionHandler(
						ge.initParams,
						ge.graph,
						edge.RunCondition,
					).Check(
						ge.GraphRuntimeState,
						previous_route_node_state,
					) {
						break
					}
				}
				next_node_id = edge.TargetNodeID
			} else {
				final_node_id := ""
				run_condition_any_none_empty := true
				run_condition_empty_flags := []bool{}
				for _, edge := range edge_mappings {
					if edge.RunCondition == nil {
						run_condition_any_none_empty = false
					}
					run_condition_empty_flags = append(run_condition_empty_flags, edge.RunCondition != nil)
				}
				mlog.Debugf("run_condition_empty_flags=%#v", run_condition_empty_flags)
				if run_condition_any_none_empty {
					// if nodes has run conditions, get node id which branch to take based on the run condition results
					condition_edge_mappings := map[string][]*graphengineentities.GraphEdge{}
					for _, edge := range edge_mappings {
						mlog.Debugf("--------edge=%#v", edge)
						mlog.Debugf("--------edge.RunCondition=%#v", edge.RunCondition)
						if edge.RunCondition != nil {
							run_condition_hash := edge.RunCondition.Hash()
							if _, ok := condition_edge_mappings[run_condition_hash]; !ok {
								condition_edge_mappings[run_condition_hash] = make([]*graphengineentities.GraphEdge, 0)
							}
							condition_edge_mappings[run_condition_hash] = append(condition_edge_mappings[run_condition_hash], edge)
						}
					}
					for _, sub_edge_mappings := range condition_edge_mappings {
						if len(sub_edge_mappings) == 0 {
							continue
						}
						edge := sub_edge_mappings[0]
						if edge.RunCondition == nil {
							mlog.Warningf("Edge %s run condition is None", edge.TargetNodeID)
							continue
						}

						if !(&conditionhandlers.ConditionManager{}).GetConditionHandler(
							ge.initParams,
							ge.graph,
							edge.RunCondition,
						).Check(
							ge.GraphRuntimeState,
							previous_route_node_state,
						) {
							continue
						}
						mlog.Debugf("---matched")
						if len(sub_edge_mappings) == 1 {
							final_node_id = edge.TargetNodeID
						} else {
							parallel_generator := ge._run_parallel_branches(
								sub_edge_mappings,
								in_parallel_id,
								parallel_start_node_id,
								handle_exceptions,
							)

							for parallel_result := range parallel_generator {
								if _, ok := parallel_result.(string); ok {
									final_node_id = parallel_result.(string)
								} else {
									if !yield(parallel_result.(graphengineentities.GraphEngineEvent)) {
										return
									}
								}
							}
						}
						break
					}
					mlog.Debugf("final_node_id=%s", final_node_id)
					if final_node_id == "" {
						break
					}
					next_node_id = final_node_id
				} else if node_instance.GetBaseNodeData().ErrorStrategy == nodesenumtypes.ErrorStrategy_FAIL_BRANCH &&
					node_instance.ShouldContinueOnError() &&
					previous_route_node_state.Status == graphengineenumtypes.RouteNodeStateStatus_EXCEPTION {
					break
				} else {
					parallel_generator := ge._run_parallel_branches(
						edge_mappings,
						in_parallel_id,
						parallel_start_node_id,
						handle_exceptions,
					)
					for generated_item := range parallel_generator {
						if _, ok := generated_item.(string); ok {
							final_node_id = generated_item.(string)
						} else {
							if !yield(generated_item.(graphengineentities.GraphEngineEvent)) {
								return
							}
						}
					}
					mlog.Debugf("final_node_id=%s", final_node_id)
					if final_node_id == "" {
						break
					}
					next_node_id = final_node_id
					mlog.Debugf("next_node_id=%s", next_node_id)
				}
			}
			if in_parallel_id != "" {
				tmp := ""
				if _, ok := ge.graph.NodeParallelMapping[next_node_id]; ok {
					tmp = ge.graph.NodeParallelMapping[next_node_id]
				}
				if tmp != in_parallel_id {
					break
				}
			}
		}
	}
}

func (ge *GraphEngine) appendVariablesRecursively(node_id string, variable_key_list []string, variable_value any /*workflowentities.VariableValue*/) {
	/*
		Append variables recursively
		:param node_id: node id
		:param variable_key_list: variable key list
		:param variable_value: variable value
		:return:
	*/
	variable_key_list = append([]string{node_id}, variable_key_list...)
	ge.GraphRuntimeState.VariablePool.Add(variable_key_list, variable_value)

	// if variable_value is a dict, then recursively append variables
	if _, ok := variable_value.(map[string]any); ok {
		for key, value := range variable_value.(map[string]any) {
			// construct new key list
			new_key_list := append(variable_key_list, key)
			ge.appendVariablesRecursively(node_id, new_key_list, value)
		}
	}
}
func (ge *GraphEngine) isTimedOut(start_at time.Time, max_execution_time int) bool {
	/*
		Check timeout
		:param start_at: start time
		:param max_execution_time: max execution time
		:return:
	*/
	return time.Since(start_at).Seconds() > float64(max_execution_time)
}

func (ge *GraphEngine) handleContinueOnError(
	node_instance base.Noder,
	error_result *workflowentities.NodeRunResult,
	variable_pool *workflowentities.VariablePool,
	handle_exceptions []string,
) *workflowentities.NodeRunResult {
	/*
		handle continue on error when ge._should_continue_on_error is true

		:param    error_result (NodeRunResult): error run result
		:param    variable_pool (VariablePool): variable pool
		:return:  excption run result
	*/
	// add error message and error type to variable pool
	variable_pool.Add([]string{node_instance.GetNodeID(), "error_message"}, error_result.Error)
	variable_pool.Add([]string{node_instance.GetNodeID(), "error_type"}, error_result.ErrorType)
	// add error message to handle_exceptions
	handle_exceptions = append(handle_exceptions, error_result.Error)
	node_error_args := &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_EXCEPTION,
		Error:  error_result.Error,
		Inputs: error_result.Inputs,
		Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
			workflowenumtypes.NodeRunMetadataKey_ERROR_STRATEGY: node_instance.GetBaseNodeData().ErrorStrategy,
		},
	}

	if node_instance.GetBaseNodeData().ErrorStrategy == nodesenumtypes.ErrorStrategy_DEFAULT_VALUE {
		node_error_args.Outputs = node_instance.GetBaseNodeData().DefaultValueDict()
		node_error_args.Outputs["error_message"] = error_result.Error
		node_error_args.Outputs["error_type"] = error_result.ErrorType
		return node_error_args
	} else if node_instance.GetBaseNodeData().ErrorStrategy == nodesenumtypes.ErrorStrategy_FAIL_BRANCH {
		if _, ok := ge.graph.EdgeMapping[node_instance.GetNodeID()]; ok {
			node_error_args.EdgeSourceHandle = string(nodesenumtypes.FailBranchSourceHandle_FAILED)
		}
		node_error_args.Outputs["error_message"] = error_result.Error
		node_error_args.Outputs["error_type"] = error_result.ErrorType
		return node_error_args
	}
	return error_result
}
