package workflow

import (
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"

	wfexceptions "mlib.com/gofy/server/core/exceptions/workflow"
	"mlib.com/gofy/server/core/workflow/callbacks"
	"mlib.com/gofy/server/core/workflow/graph"
	graphengine "mlib.com/gofy/server/core/workflow/graph_engine"
	"mlib.com/gofy/server/core/workflow/nodes"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	"mlib.com/gofy/server/core/workflow/utils/condition"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	filefactory "mlib.com/gofy/server/factories/file_factory"
)

type WorkflowEntry struct {
	GraphEngine *graphengine.GraphEngine
}

func NewWorkflowEntry(
	tenant_id string,
	app_id string,
	workflow_id string,
	workflow_type models.WorkflowType,
	graph_config map[string]any,
	gf *graph.Graph,
	user_id string,
	user_from models.UserFrom,
	invoke_from appenumtypes.InvokeFrom,
	call_depth int,
	variable_pool *workflowentities.VariablePool,
) *WorkflowEntry {
	/*
		Init workflow entry
		:param tenant_id: tenant id
		:param app_id: app id
		:param workflow_id: workflow id
		:param workflow_type: workflow type
		:param graph_config: workflow graph config
		:param graph: workflow graph
		:param user_id: user id
		:param user_from: user from
		:param invoke_from: invoke from
		:param call_depth: call depth
		:param variable_pool: variable pool
		:param thread_pool_id: thread pool id
	*/
	// check call depth
	workflow_call_max_depth := viper.GetIntWithDefault("workflow.call_max_depth", 5)
	if call_depth > workflow_call_max_depth {
		panic(exceptions.NewValueError(fmt.Sprintf("Max workflow call depth %d reached.", workflow_call_max_depth)))
	}
	// init workflow run state
	wfe := &WorkflowEntry{}
	wfe.GraphEngine = graphengine.NewGraphEngine(
		tenant_id,
		app_id,
		workflow_type,
		workflow_id,
		user_id,
		user_from,
		invoke_from,
		call_depth,
		gf,
		graph_config,
		variable_pool,
		viper.GetIntWithDefault("workflow.max_execution_steps", 500),
		viper.GetIntWithDefault("workflow.max_execution_time", 1200),
	)
	return wfe
}

func (wfe *WorkflowEntry) Run(
	cbs []callbacks.WorkflowCallback,
) iter.Seq[graphengineentities.GraphEngineEvent] {
	/*
		:param callbacks: workflow callbacks
	*/

	return func(yield func(graphengineentities.GraphEngineEvent) bool) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if exp, ok := r.(error); ok {
						for _, cb := range cbs {
							cb.OnEvent(&graphengineentities.GraphRunFailedEvent{
								Error: exp.Error(),
							})
						}
					} else {
						panic(r)
					}
				}
			}()
			// try:
			// run workflow
			generator := wfe.GraphEngine.Run()
			for ev := range generator {
				if !yield(ev) {
					return
				}
			}
			mlog.Debugf("------ workflow run return")
		}()
	}
	// except GenerateTaskStoppedError:
	// 	pass
	// except Exception as e:
	// 	logger.exception("Unknown Error when workflow entry running")
	// 	if callbacks:
	// 		for callback in callbacks:
	// 			callback.on_event(event=GraphRunFailedEvent(error=str(e)))
	// 	return
}
func (wfe *WorkflowEntry) SingleStepRun(
	wf *models.Workflow,
	node_id string,
	user_id string,
	user_inputs map[string]any,
) (base.Noder, iter.Seq[any]) {
	/*
		Single step run workflow node
		:param workflow: Workflow instance
		:param node_id: node id
		:param user_id: user id
		:param user_inputs: user inputs
		:return:
	*/
	// fetch node info from workflow graph
	workflow_graph := wf.GraphDict()
	if len(workflow_graph) == 0 {
		mlog.Errorf("workflow graph not found")
		panic(exceptions.NewValueError("workflow graph not found"))
	}
	mlog.Debugf("------workflow_graph=%#v", workflow_graph)
	var node_configs []map[string]any
	if _, ok := workflow_graph["nodes"]; ok {
		if _, ok := workflow_graph["nodes"].([]map[string]any); ok {
			node_configs = workflow_graph["nodes"].([]map[string]any)
		} else if _, ok := workflow_graph["nodes"].([]any); ok {
			for _, v := range workflow_graph["nodes"].([]any) {
				if _, ok := v.(map[string]any); ok {
					node_configs = append(node_configs, v.(map[string]any))
				}
			}
		}
	}
	if len(node_configs) == 0 {
		mlog.Errorf("nodes not found in workflow graph")
		panic(exceptions.NewValueError("nodes not found in workflow graph"))
	}
	var node_config map[string]any
	for _, node := range node_configs {
		id := ""
		if _, ok := node["id"]; ok {
			if _, ok := node["id"].(string); ok {
				id = node["id"].(string)
			}
		}
		if id == node_id {
			node_config = node
			break
		}
	}
	if len(node_config) == 0 {
		mlog.Errorf("node id not found in workflow graph")
		panic(exceptions.NewValueError("node id not found in workflow graph"))
	}
	// init variable pool
	variable_pool := workflowentities.NewVariablePool(nil, nil, wf.GetEnvironmentVariables(), nil)

	// init graph
	gf := graph.NewGraph(wf.GraphDict(), "")

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

	// init workflow run state
	node_instance := nodes.NewNode(
		uuid.NewV4().String(),
		node_config,
		&graphengineentities.GraphInitParams{
			TenantID:     wf.TenantID,
			AppID:        wf.AppID,
			WorkflowType: wf.Type,
			WorkflowID:   wf.ID,
			GraphConfig:  wf.GraphDict(),
			UserID:       user_id,
			UserFrom:     models.UserFrom_ACCOUNT,
			InvokeFrom:   appenumtypes.InvokeFrom_DEBUGGER,
			CallDepth:    0,
		},
		gf,
		&graphengineentities.GraphRuntimeState{VariablePool: variable_pool, StartAt: time.Now()},
		"",
		node_type,
	)
	variable_mapping := nodes.ExtractVariableSelectorToVariableMappingByNoder(wf.GraphDict(), node_instance.GetNodeID(), node_instance)

	if variable_mapping == nil {
		variable_mapping = make(map[string][]string)
	}
	mlog.Debugf("------variable_mapping=%#v", variable_mapping)
	mlog.Debugf("------user_inputs=%#v", user_inputs)
	wfe.MappingUserInputsToVariablePool(
		variable_mapping,
		user_inputs,
		variable_pool,
		wf.TenantID,
	)

	return node_instance, func(yield func(any) bool) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					panic(wfexceptions.NewWorkflowNodeRunFailedError(node_instance, err.Error()))
				} else {
					panic(r)
				}
			}
		}()
		generator := node_instance.RunIter(node_instance)
		for item := range generator {
			if !yield(item) {
				return
			}
		}
	}
}

func (wfe *WorkflowEntry) MappingUserInputsToVariablePool(
	variable_mapping map[string][]string,
	user_inputs map[string]any,
	variable_pool *workflowentities.VariablePool,
	tenant_id string,
) {
	for node_variable, variable_selector := range variable_mapping {
		mlog.Debugf("------node_variable=%#v, variable_selector=%#v", node_variable, variable_selector)
		// fetch node id and variable key from node_variable
		node_variable_list := strings.Split(node_variable, ".")
		if len(node_variable_list) < 1 {
			panic(exceptions.NewValueError(fmt.Sprintf("Invalid node variable %s", node_variable)))
		}
		node_variable_key := strings.Join(node_variable_list[1:], ".")
		_, ok1 := user_inputs[node_variable_key]
		_, ok2 := user_inputs[node_variable]
		if (!ok1 && !ok2) && variable_pool.Get(variable_selector) == nil {
			panic(exceptions.NewValueError(fmt.Sprintf("Variable key %s not found in user inputs.", node_variable)))
		}
		// environment variable already exist in variable pool, not from user inputs
		if variable_pool.Get(variable_selector) != nil {
			continue
		}
		// fetch variable node id from variable selector
		variable_node_id := variable_selector[0]
		variable_key_list := variable_selector[1:]

		// get input value
		var input_value any
		if _, ok := user_inputs[node_variable]; !ok {
			input_value = user_inputs[node_variable_key]
		} else {
			input_value = user_inputs[node_variable]
		}
		if _, ok := input_value.(map[string]any); ok {
			if _, ok := input_value.(map[string]any)["type"]; ok {
				if _, ok := input_value.(map[string]any)["transfer_method"]; ok {
					f := filefactory.BuildFromMapping(input_value.(map[string]any), tenant_id, nil)
					input_value = f
					// variable_pool.Add(append([]string{variable_node_id}, variable_key_list...), f)
				}
			}
		}

		if _, ok := input_value.([]map[string]any); ok {
			is_file_mappings := []bool{}
			for _, v := range input_value.([]map[string]any) {
				is_true := false
				if _, ok := v["type"]; ok {
					if _, ok := v["transfer_method"]; ok {
						is_true = true
					}
				}
				is_file_mappings = append(is_file_mappings, is_true)
			}
			if condition.AllTrue(is_file_mappings) {
				fs := filefactory.BuildFromMappings(input_value.([]map[string]any), tenant_id, nil)
				input_value = fs
				//
			}
		}
		variable_pool.Add(append([]string{variable_node_id}, variable_key_list...), input_value)
	}
}

func (wfe *WorkflowEntry) HandleSpecialValues(value map[string]any) map[string]any {
	result := wfe.handleSpecialValues(value)
	if v, ok := result.(map[string]any); ok || v == nil {
		return result.(map[string]any)
	}
	var rsp map[string]any
	bindata, _ := json.Marshal(result)
	json.Unmarshal(bindata, &rsp)
	return rsp
}
func (wfe *WorkflowEntry) handleSpecialValues(value any) any {
	if value == nil {
		return value
	}
	if _, ok := value.(map[string]any); ok {
		res := map[string]any{}
		for k, v := range value.(map[string]any) {
			res[k] = wfe.handleSpecialValues(v)
		}
		return res
	}
	if _, ok := value.([]any); ok {
		res_list := []any{}
		for _, item := range value.([]any) {
			res_list = append(res_list, wfe.handleSpecialValues(item))
		}
		return res_list
	}
	if _, ok := value.(*file.File); ok {
		return value.(*file.File).ToDict()
	}
	return value
}
