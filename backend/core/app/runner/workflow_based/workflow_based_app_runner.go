package workflowbased

import (
	"encoding/json"
	"slices"

	"github.com/odysseythink/gofy/backend/core/app/runner/base"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	coreworkflow "github.com/odysseythink/gofy/backend/core/workflow"
	"github.com/odysseythink/gofy/backend/core/workflow/graph"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appqueueentities "github.com/odysseythink/gofy/backend/entities/app/queue"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
	nodesentities "github.com/odysseythink/gofy/backend/entities/nodes"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	workflowenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type WorkflowBasedAppRunner[T interface {
	*appqueueentities.MessageQueueMessage | *appqueueentities.WorkflowQueueMessage
}] struct {
	*base.AppRunner[T]
	queueManager appqueueentities.AppQueueManager[T]
}

func New[T interface {
	*appqueueentities.MessageQueueMessage | *appqueueentities.WorkflowQueueMessage
}](queue_manager appqueueentities.AppQueueManager[T]) *WorkflowBasedAppRunner[T] {
	return &WorkflowBasedAppRunner[T]{
		queueManager: queue_manager,
	}
}

func (r *WorkflowBasedAppRunner[T]) InitGraph(graph_config map[string]any) *graph.Graph {
	/*
	   Init graph
	*/
	_, ok1 := graph_config["nodes"]
	_, ok2 := graph_config["edges"]
	if !ok1 || !ok2 {
		mlog.Errorf("nodes or edges not found in workflow graph")
		panic(exceptions.NewValueError("nodes or edges not found in workflow graph"))
	}
	if _, ok := graph_config["nodes"].([]any); !ok {
		mlog.Errorf("nodes in workflow graph must be a list")
		panic(exceptions.NewValueError("nodes in workflow graph must be a list"))
	}
	if _, ok := graph_config["edges"].([]any); !ok {
		mlog.Errorf("edges in workflow graph must be a list")
		panic(exceptions.NewValueError("edges in workflow graph must be a list"))
	}
	// init graph
	graph := graph.NewGraph(graph_config, "")
	if graph == nil {
		mlog.Errorf("create graph failed")
		panic(exceptions.NewValueError("graph not found in workflow"))
	}

	return graph

}
func (r *WorkflowBasedAppRunner[T]) GetGraphAndVariablePoolOfSingleIteration(
	wf *models.Workflow,
	node_id string,
	user_inputs map[string]any,
) (*graph.Graph, *workflowentities.VariablePool) {
	/*
	   Get variable pool of single iteration
	*/
	// fetch workflow graph
	graph_config := wf.GraphDict()
	if graph_config == nil {
		mlog.Errorf("workflow graph not found")
		panic(exceptions.NewValueError("workflow graph not found"))
	}
	// graph_config = cast(dict[str, Any], graph_config)
	_, ok1 := graph_config["nodes"]
	_, ok2 := graph_config["edges"]
	if !ok1 || !ok2 {
		mlog.Errorf("nodes or edges not found in workflow graph")
		panic(exceptions.NewValueError("nodes or edges not found in workflow graph"))
	}
	if _, ok := graph_config["nodes"].([]map[string]any); !ok {
		mlog.Errorf("nodes in workflow graph must be a list")
		panic(exceptions.NewValueError("nodes in workflow graph must be a list"))
	}
	if _, ok := graph_config["edges"].([]map[string]any); !ok {
		mlog.Errorf("edges in workflow graph must be a list")
		panic(exceptions.NewValueError("edges in workflow graph must be a list"))
	}
	// filter nodes only in iteration
	node_configs := []map[string]any{}
	for _, node := range graph_config["nodes"].([]map[string]any) {
		id := ""
		if _, ok := node["id"]; ok {
			if _, ok := node["id"].(string); ok {
				id = node["id"].(string)
			}
		}
		data := map[string]any{}
		if _, ok := node["data"]; ok {
			if _, ok := node["data"].(map[string]any); ok {
				data = node["data"].(map[string]any)
			}
		}
		iteration_id := ""
		if _, ok := data["iteration_id"]; ok {
			if _, ok := data["iteration_id"].(string); ok {
				iteration_id = data["iteration_id"].(string)
			}
		}
		if id == node_id || iteration_id == id {
			node_configs = append(node_configs, node)
		}
	}

	graph_config["nodes"] = node_configs
	node_ids := []string{}
	for _, node := range node_configs {
		id := ""
		if _, ok := node["id"]; ok {
			if _, ok := node["id"].(string); ok {
				id = node["id"].(string)
			}
		}
		node_ids = append(node_ids, id)
	}

	// filter edges only in iteration
	edge_configs := []map[string]any{}
	for _, edge := range graph_config["edges"].([]map[string]any) {
		source := ""
		if _, ok := edge["source"]; ok {
			if _, ok := edge["source"].(string); ok {
				source = edge["source"].(string)
			}
		}
		target := ""
		if _, ok := edge["target"]; ok {
			if _, ok := edge["target"].(string); ok {
				target = edge["target"].(string)
			}
		}
		if (source == "" || slices.Contains(node_ids, source)) && (target == "" || slices.Contains(node_ids, target)) {
			edge_configs = append(edge_configs, edge)
		}
	}

	graph_config["edges"] = edge_configs

	// init graph
	graph := graph.NewGraph(graph_config, node_id)
	if graph == nil {
		mlog.Errorf("create graph failed")
		panic(exceptions.NewValueError("graph not found in workflow"))
	}

	// fetch node config from node id
	var iteration_node_config map[string]any
	for _, node := range node_configs {
		id := ""
		if _, ok := node["id"]; ok {
			if _, ok := node["id"].(string); ok {
				id = node["id"].(string)
			}
		}
		if id == node_id {
			iteration_node_config = node
			break
		}
	}
	if len(iteration_node_config) == 0 {
		panic(exceptions.NewValueError("iteration node id not found in workflow graph"))
	}
	// Get node class
	var node_type nodesenumtypes.NodeType
	if _, ok := iteration_node_config["data"]; ok {
		if _, ok := iteration_node_config["data"].(map[string]any); ok {
			if _, ok := iteration_node_config["data"].(map[string]any)["type"]; ok {
				if v, ok := iteration_node_config["data"].(map[string]any)["type"].(string); ok {
					node_type, _ = nodesenumtypes.ParseNodeType(v)
				}
			}
		}
	}
	if _, ok := iteration_node_config["id"]; ok {
		if _, ok := iteration_node_config["id"].(string); ok {
			node_id = iteration_node_config["id"].(string)
		}
	}
	// node_version := "1"
	// if _, ok := iteration_node_config["data"]; ok {
	// 	if _, ok := iteration_node_config["data"].(map[string]any); ok {
	// 		if _, ok := iteration_node_config["data"].(map[string]any)["version"]; ok {
	// 			if _, ok := iteration_node_config["data"].(map[string]any)["version"].(string); ok {
	// 				node_version = iteration_node_config["data"].(map[string]any)["version"].(string)
	// 			}
	// 		}
	// 	}
	// }

	// init variable pool
	variable_pool := workflowentities.NewVariablePool(nil, nil, wf.GetEnvironmentVariables(), nil)

	variable_mapping := nodes.ExtractVariableSelectorToVariableMapping(wf.GraphDict(), node_id, nodesentities.NewNodeDataByNodeType(node_type))

	(&coreworkflow.WorkflowEntry{}).MappingUserInputsToVariablePool(
		variable_mapping,
		user_inputs,
		variable_pool,
		wf.TenantID,
	)

	return graph, variable_pool

}

func (r *WorkflowBasedAppRunner[T]) GetWorkflow(app_model *models.App, workflow_id string) *models.Workflow {
	/*
	   Get workflow
	*/
	// fetch workflow by workflow_id
	wf := new(models.Workflow)
	err := dbengine.Instance().DB.Model(&models.Workflow{}).Where("tenant_id = ? and app_id = ? and id = ?", app_model.TenantID, app_model.ID, workflow_id).Preload("Tenant").Preload("App").Preload("CreatedByAccount").Preload("UpdatedByAccount").First(wf).Error
	if err != nil {
		mlog.Errorf("get workflow failed:%v", err)
		return nil
	}

	// return workflow
	return wf

}
func (r *WorkflowBasedAppRunner[T]) PublishEvent(event appqueueentities.AppQueueEventer) {
	r.queueManager.Publish(event, appenumtypes.PublishFrom_APPLICATION_MANAGER)
}
func (r *WorkflowBasedAppRunner[T]) HandleEvent(workflow_entry *coreworkflow.WorkflowEntry, event graphengineentities.GraphEngineEvent) {
	/*
	   Handle event
	   :param workflow_entry: workflow entry
	   :param event: event
	*/
	mlog.Debugf("------event=%#v", event)
	if any(event) == nil || workflow_entry == nil {
		return
	}
	switch data := any(event).(type) {
	case *graphengineentities.GraphRunStartedEvent:
		r.PublishEvent(&appqueueentities.QueueWorkflowStartedEvent{GraphRuntimeState: workflow_entry.GraphEngine.GraphRuntimeState})
	case *graphengineentities.GraphRunSucceededEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueWorkflowSucceededEvent{Outputs: data.Outputs})
	case *graphengineentities.GraphRunPartialSucceededEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueWorkflowPartialSuccessEvent{Outputs: data.Outputs, ExceptionsCount: data.ExceptionsCount})
	case *graphengineentities.GraphRunFailedEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueWorkflowFailedEvent{Error: data.Error, ExceptionsCount: data.ExceptionsCount})
	case *graphengineentities.NodeRunRetryEvent:
		if data == nil {
			return
		}
		node_run_result := data.RouteNodeState.NodeRunResult
		inputs := map[string]any{}
		process_data := map[string]any{}
		outputs := map[string]any{}
		execution_metadata := map[workflowenumtypes.NodeRunMetadataKey]any{}
		if node_run_result != nil {
			inputs = node_run_result.Inputs
			process_data = node_run_result.ProcessData
			outputs = node_run_result.Outputs
			execution_metadata = node_run_result.Metadata
		}
		r.PublishEvent(&appqueueentities.QueueNodeRetryEvent{
			QueueNodeStartedEvent: &appqueueentities.QueueNodeStartedEvent{
				NodeExecutionID:           data.ID,
				NodeID:                    data.NodeID,
				NodeType:                  data.NodeType,
				NodeData:                  data.NodeData,
				ParallelID:                data.ParallelID,
				ParallelStartNodeID:       data.ParallelStartNodeID,
				ParentParallelID:          data.ParentParallelID,
				ParentParallelStartNodeID: data.ParentParallelStartNodeID,
				StartAt:                   data.StartAt,
				NodeRunIndex:              data.RouteNodeState.Index,
				PredecessorNodeID:         data.PredecessorNodeID,
				InIterationID:             data.InIterationID,
				ParallelModeRunID:         data.ParallelModeRunID,
			},
			Inputs:            inputs,
			ProcessData:       process_data,
			Outputs:           outputs,
			Error:             data.Error,
			ExecutionMetadata: execution_metadata,
			RetryIndex:        data.RetryIndex,
		})
	case *graphengineentities.NodeRunStartedEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueNodeStartedEvent{
			NodeExecutionID:           data.ID,
			NodeID:                    data.NodeID,
			NodeType:                  data.NodeType,
			NodeData:                  data.NodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.RouteNodeState.StartAt,
			NodeRunIndex:              data.RouteNodeState.Index,
			PredecessorNodeID:         data.PredecessorNodeID,
			InIterationID:             data.InIterationID,
			ParallelModeRunID:         data.ParallelModeRunID,
		})
	case *graphengineentities.NodeRunSucceededEvent:
		if data == nil {
			return
		}
		node_run_result := data.RouteNodeState.NodeRunResult
		inputs := map[string]any{}
		process_data := map[string]any{}
		outputs := map[string]any{}
		execution_metadata := map[workflowenumtypes.NodeRunMetadataKey]any{}
		if node_run_result != nil {
			inputs = node_run_result.Inputs
			process_data = node_run_result.ProcessData
			outputs = node_run_result.Outputs
			execution_metadata = node_run_result.Metadata
		}
		r.PublishEvent(&appqueueentities.QueueNodeSucceededEvent{
			NodeExecutionID:           data.ID,
			NodeID:                    data.NodeID,
			NodeType:                  data.NodeType,
			NodeData:                  data.NodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.RouteNodeState.StartAt,
			Inputs:                    inputs,
			ProcessData:               process_data,
			Outputs:                   outputs,
			ExecutionMetadata:         execution_metadata,
			InIterationID:             data.InIterationID,
		})
	case *graphengineentities.NodeRunFailedEvent:
		if data == nil {
			return
		}
		qevent := &appqueueentities.QueueNodeFailedEvent{
			NodeExecutionID:           data.ID,
			NodeID:                    data.NodeID,
			NodeType:                  data.NodeType,
			NodeData:                  data.NodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.RouteNodeState.StartAt,
			InIterationID:             data.InIterationID,
		}
		if data.RouteNodeState.NodeRunResult != nil {
			qevent.Inputs = data.RouteNodeState.NodeRunResult.Inputs
			qevent.ProcessData = data.RouteNodeState.NodeRunResult.ProcessData
			qevent.Outputs = data.RouteNodeState.NodeRunResult.Outputs
			if data.RouteNodeState.NodeRunResult.Error != "" {
				qevent.Error = data.RouteNodeState.NodeRunResult.Error
			} else {
				qevent.Error = "Unknown error"
			}
			qevent.ExecutionMetadata = data.RouteNodeState.NodeRunResult.Metadata
		} else {
			qevent.Inputs = map[string]any{}
			qevent.ProcessData = map[string]any{}
			qevent.Outputs = map[string]any{}
			qevent.Error = "Unknown error"
			qevent.ExecutionMetadata = map[workflowenumtypes.NodeRunMetadataKey]any{}
		}

		r.PublishEvent(qevent)
	case *graphengineentities.NodeRunExceptionEvent:
		if data == nil {
			return
		}
		qevent := &appqueueentities.QueueNodeExceptionEvent{
			NodeExecutionID:           data.ID,
			NodeID:                    data.NodeID,
			NodeType:                  data.NodeType,
			NodeData:                  data.NodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.RouteNodeState.StartAt,
			InIterationID:             data.InIterationID,
		}
		if data.RouteNodeState.NodeRunResult != nil {
			qevent.Inputs = data.RouteNodeState.NodeRunResult.Inputs
			qevent.ProcessData = data.RouteNodeState.NodeRunResult.ProcessData
			qevent.Outputs = data.RouteNodeState.NodeRunResult.Outputs
			if data.RouteNodeState.NodeRunResult.Error != "" {
				qevent.Error = data.RouteNodeState.NodeRunResult.Error
			} else {
				qevent.Error = "Unknown error"
			}
			qevent.ExecutionMetadata = data.RouteNodeState.NodeRunResult.Metadata
		} else {
			qevent.Inputs = map[string]any{}
			qevent.ProcessData = map[string]any{}
			qevent.Outputs = map[string]any{}
			qevent.Error = "Unknown error"
			qevent.ExecutionMetadata = map[workflowenumtypes.NodeRunMetadataKey]any{}
		}
		r.PublishEvent(qevent)
	case *graphengineentities.NodeInIterationFailedEvent:
		if data == nil {
			return
		}
		qevent := &appqueueentities.QueueNodeInIterationFailedEvent{
			NodeExecutionID:           data.ID,
			NodeID:                    data.NodeID,
			NodeType:                  data.NodeType,
			NodeData:                  data.NodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.RouteNodeState.StartAt,
			InIterationID:             data.InIterationID,
			Error:                     data.Error,
		}
		if data.RouteNodeState.NodeRunResult != nil {
			qevent.Inputs = data.RouteNodeState.NodeRunResult.Inputs
			qevent.ProcessData = data.RouteNodeState.NodeRunResult.ProcessData
			qevent.Outputs = data.RouteNodeState.NodeRunResult.Outputs
			qevent.ExecutionMetadata = data.RouteNodeState.NodeRunResult.Metadata
		} else {
			qevent.Inputs = map[string]any{}
			qevent.ProcessData = map[string]any{}
			qevent.Outputs = map[string]any{}
			qevent.ExecutionMetadata = map[workflowenumtypes.NodeRunMetadataKey]any{}
		}
		r.PublishEvent(qevent)
	case *graphengineentities.NodeRunStreamChunkEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueTextChunkEvent{
			Text:                 data.ChunkContent,
			FromVariableSelector: data.FromVariableSelector,
			InIterationID:        data.InIterationID,
		})
	case *graphengineentities.NodeRunRetrieverResourceEvent:
		if data == nil {
			return
		}
		var rsrc []*ragentities.RetrievalSourceMetadata
		if b, err := json.Marshal(data.RetrieverResources); err == nil {
			_ = json.Unmarshal(b, &rsrc)
		}
		r.PublishEvent(&appqueueentities.QueueRetrieverResourcesEvent{RetrieverResources: rsrc, InIterationID: data.InIterationID})
	case *graphengineentities.ParallelBranchRunStartedEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueParallelBranchRunStartedEvent{
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			InIterationID:             data.InIterationID,
		})
	case *graphengineentities.ParallelBranchRunSucceededEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueParallelBranchRunSucceededEvent{
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			InIterationID:             data.InIterationID,
		})
	case *graphengineentities.ParallelBranchRunFailedEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueParallelBranchRunFailedEvent{
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			InIterationID:             data.InIterationID,
			Error:                     data.Error,
		})
	case *graphengineentities.IterationRunStartedEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueIterationStartEvent{
			NodeExecutionID:           data.IterationID,
			NodeID:                    data.IterationNodeID,
			NodeType:                  data.IterationNodeType,
			NodeData:                  data.IterationNodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.StartAt,
			NodeRunIndex:              workflow_entry.GraphEngine.GraphRuntimeState.NodeRunSteps,
			Inputs:                    data.Inputs,
			PredecessorNodeID:         data.PredecessorNodeID,
			Metadata:                  data.Metadata,
		})
	case *graphengineentities.IterationRunNextEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueIterationNextEvent{
			NodeExecutionID:           data.IterationID,
			NodeID:                    data.IterationNodeID,
			NodeType:                  data.IterationNodeType,
			NodeData:                  data.IterationNodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			Index:                     data.Index,
			NodeRunIndex:              workflow_entry.GraphEngine.GraphRuntimeState.NodeRunSteps,
			Output:                    data.PreIterationOutput,
			ParallelModeRunID:         data.ParallelModeRunID,
			Duration:                  data.Duration,
		})
	case *graphengineentities.IterationRunSucceededEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueIterationCompletedEvent{
			NodeExecutionID:           data.IterationID,
			NodeID:                    data.IterationNodeID,
			NodeType:                  data.IterationNodeType,
			NodeData:                  data.IterationNodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.StartAt,
			NodeRunIndex:              workflow_entry.GraphEngine.GraphRuntimeState.NodeRunSteps,
			Inputs:                    data.Inputs,
			Outputs:                   data.Outputs,
			Metadata:                  data.Metadata,
			Steps:                     data.Steps,
		})
	case *graphengineentities.IterationRunFailedEvent:
		if data == nil {
			return
		}
		r.PublishEvent(&appqueueentities.QueueIterationCompletedEvent{
			NodeExecutionID:           data.IterationID,
			NodeID:                    data.IterationNodeID,
			NodeType:                  data.IterationNodeType,
			NodeData:                  data.IterationNodeData,
			ParallelID:                data.ParallelID,
			ParallelStartNodeID:       data.ParallelStartNodeID,
			ParentParallelID:          data.ParentParallelID,
			ParentParallelStartNodeID: data.ParentParallelStartNodeID,
			StartAt:                   data.StartAt,
			NodeRunIndex:              workflow_entry.GraphEngine.GraphRuntimeState.NodeRunSteps,
			Inputs:                    data.Inputs,
			Outputs:                   data.Outputs,
			Metadata:                  data.Metadata,
			Steps:                     data.Steps,
			Error:                     data.Error,
		})
	}
}
