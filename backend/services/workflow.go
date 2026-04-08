package services

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"time"

	uuid "github.com/satori/go.uuid"
	advancedchatcfgmgr "mlib.com/gofy/server/core/app/config_manageres/advanced_chat"
	workflowcfgmgr "mlib.com/gofy/server/core/app/config_manageres/workflow"
	"mlib.com/gofy/server/core/exceptions"
	wfexceptions "mlib.com/gofy/server/core/exceptions/workflow"
	"mlib.com/gofy/server/core/variables"
	coreworkflow "mlib.com/gofy/server/core/workflow"
	"mlib.com/gofy/server/core/workflow/nodes"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	dbengine "mlib.com/gofy/server/db_engine"
	nodesevententities "mlib.com/gofy/server/entities/nodes/event"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/events"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type WorkflowService struct {
}

func (s *WorkflowService) GetDraftWorkflow(app *models.App) *models.Workflow {
	// """
	// Get draft workflow
	// """
	// fetch draft workflow by app_model
	wf := &models.Workflow{}
	err := dbengine.Instance().DB.Model(&models.Workflow{}).Where("tenant_id = ? and app_id = ? and version = 'draft'", app.TenantID, app.ID).Preload("Tenant").Preload("App").Preload("CreatedByAccount").Preload("UpdatedByAccount").First(wf).Error
	if err != nil {
		mlog.Errorf("get Workflow by tenant_id = %s and app_id = %s and version = 'draft' failed:%v", app.TenantID, app.ID, err)
		wf = nil
	}
	return wf
}
func (s *WorkflowService) GetDefaultBlockConfigs() []map[string]any {
	default_block_configs := []map[string]any{}
	for _, v := range nodes.NODE_TYPE_CLASSES_MAPPING {
		cfg := v[nodes.LATEST_VERSION].GetDefaultConfig(nil)
		if len(cfg) > 0 {
			default_block_configs = append(default_block_configs, cfg)
		}
	}
	return default_block_configs
}
func (s *WorkflowService) GetDefaultBlockConfig(node_type nodesenumtypes.NodeType, filters map[string]any) map[string]any {
	if !node_type.Valid() {
		mlog.Errorf("invalid node_type=%v", node_type)
		return nil
	}
	// return default block config
	if _, ok := nodes.NODE_TYPE_CLASSES_MAPPING[node_type]; !ok {
		mlog.Errorf("node_type=%s not exist in NODE_TYPE_CLASSES_MAPPING", node_type)
		return nil
	}

	node_class := nodes.NODE_TYPE_CLASSES_MAPPING[node_type][nodes.LATEST_VERSION]
	default_config := node_class.GetDefaultConfig(filters)

	return default_config
}
func (s *WorkflowService) GetPublishedWorkflow(app *models.App) *models.Workflow {
	if app.WorkflowID == "" {
		mlog.Errorf("invalid workflow_id")
		return nil
	}
	ret := &models.Workflow{}
	err := dbengine.Instance().DB.Model(&models.Workflow{}).Where("tenant_id = ? and app_id = ? and id = ?", app.TenantID, app.ID, app.WorkflowID).First(ret).Error
	if err != nil {
		mlog.Errorf("get Workflow by tenant_id = %s and app_id = %s and id = '%s failed:%v", app.TenantID, app.ID, app.WorkflowID, err)
		return nil
	}
	return ret
}

func (s *WorkflowService) validate_features_structure(app_model *models.App, features map[string]any) map[string]any {
	if app_model.Mode == models.AppMode_ADVANCED_CHAT {
		return advancedchatcfgmgr.New().ConfigValidate(
			app_model.TenantID, features, true,
		)
	} else if app_model.Mode == models.AppMode_WORKFLOW {
		return workflowcfgmgr.New().ConfigValidate(
			app_model.TenantID, features, true,
		)
	} else {
		panic(exceptions.NewValueError(fmt.Sprintf("Invalid app mode: {%s}", app_model.Mode)))
	}
}

func (s *WorkflowService) SyncDraftWorkflow(
	app_model *models.App,
	graph map[string]any,
	features map[string]any,
	unique_hash string,
	account *models.Account,
	environment_variables []variables.Variabler,
	conversation_variables []variables.Variabler,
) *models.Workflow {
	// fetch draft workflow by app_model
	mlog.Debugf("-----app=%#v", app_model)
	wf := s.GetDraftWorkflow(app_model)
	mlog.Debugf("-----workflow=%#v", wf)

	// if wf != nil && wf.UniqueHash() == unique_hash {
	// 	panic(appexceptions.NewWorkflowHashNotEqualError(""))
	// }
	// validate features structure
	s.validate_features_structure(app_model, features)

	// create draft workflow if not found
	if wf == nil {
		wf = &models.Workflow{
			ID:        uuid.NewV4().String(),
			TenantID:  app_model.TenantID,
			AppID:     app_model.ID,
			Version:   "draft",
			CreatedBy: account.ID,
		}
		if app_model.Mode == models.AppMode_WORKFLOW {
			wf.Type = models.WorkflowType(models.AppMode_WORKFLOW)
		} else {
			wf.Type = models.WorkflowType(models.AppMode_CHAT)
		}
		tmpbin, _ := json.Marshal(graph)
		wf.Graph = string(tmpbin)
		tmpbin, _ = json.Marshal(features)
		wf.FeaturesStr = string(tmpbin)
		wf.SetEnvironmentVariables(environment_variables)
		wf.SetConversationVariables(conversation_variables)
		dbengine.Instance().Create(wf)
		// update draft workflow if found
	} else {
		now := time.Now()
		tmpbin, _ := json.Marshal(graph)
		wf.Graph = string(tmpbin)
		tmpbin, _ = json.Marshal(features)
		wf.FeaturesStr = string(tmpbin)
		wf.SetEnvironmentVariables(environment_variables)
		wf.SetConversationVariables(conversation_variables)
		wf.UpdatedBy = account.ID
		wf.UpdatedAt = &now
		dbengine.Instance().DB.Model(&models.Workflow{}).Where("id = ?", wf.ID).Updates(models.Workflow{
			Graph:                    wf.Graph,
			FeaturesStr:              wf.FeaturesStr,
			EnvironmentVariablesStr:  wf.EnvironmentVariablesStr,
			ConversationVariablesStr: wf.ConversationVariablesStr,
			UpdatedBy:                account.ID,
			UpdatedAt:                &now,
		})
	}

	// trigger app workflow events
	events.Instance.AppDraftWorkflowWasSyncedSig.Emit(app_model, wf)

	// return draft workflow
	return wf
}

func (s *WorkflowService) Get(id string) (wf *models.Workflow, err error) {
	wf = &models.Workflow{}
	err = dbengine.Instance().DB.Model(&models.Workflow{}).Where("id = ?", id).Preload("Tenant").Preload("App").Preload("CreatedByAccount").Preload("UpdatedByAccount").First(wf).Error
	if err != nil {
		wf = nil
	}
	return
}
func (s *WorkflowService) RunDraftWorkflowNode(
	app_model *models.App, node_id string, user_inputs map[string]any, account *models.Account,
) (workflow_node_execution *models.WorkflowNodeExecution) {
	/*
		Run draft workflow node
	*/
	// fetch draft workflow by app_model
	draft_workflow := s.GetDraftWorkflow(app_model)
	if draft_workflow == nil {
		panic(exceptions.NewValueError("Workflow not initialized"))
	}
	// run draft workflow node
	start_at := time.Now()
	var run_succeeded bool
	var node_run_result *workflowentities.NodeRunResult
	var errmsg string
	var node_instance base.Noder
	defer func() {
		if r := recover(); r != nil {
			if exp, ok := r.(*wfexceptions.WorkflowNodeRunFailedError); ok {
				if workflow_node_execution != nil {
					if node_instance == nil {
						node_instance = exp.NodeInstance
					}
					run_succeeded = false
					node_run_result = nil
					errmsg = exp.Errmsg
				}
			}
			panic(r)
		}
		if run_succeeded && node_run_result != nil {
			// create workflow node execution
			var inputs map[string]any
			if len(node_run_result.Inputs) > 0 {
				inputs = (&coreworkflow.WorkflowEntry{}).HandleSpecialValues(node_run_result.Inputs)
			}
			var process_data map[string]any
			if len(node_run_result.ProcessData) > 0 {
				process_data = (&coreworkflow.WorkflowEntry{}).HandleSpecialValues(node_run_result.ProcessData)
			}
			var outputs map[string]any
			if len(node_run_result.Outputs) > 0 {
				outputs = (&coreworkflow.WorkflowEntry{}).HandleSpecialValues(node_run_result.Outputs)
			}

			bindata, _ := json.Marshal(inputs)
			workflow_node_execution.Inputs = string(bindata)
			bindata, _ = json.Marshal(process_data)
			workflow_node_execution.ProcessData = string(bindata)
			bindata, _ = json.Marshal(outputs)
			workflow_node_execution.Outputs = string(bindata)
			bindata, _ = json.Marshal(node_run_result.Metadata)
			workflow_node_execution.ExecutionMetadata = string(bindata)
			if node_run_result.Status == models.WorkflowNodeExecutionStatus_SUCCEEDED {
				workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_SUCCEEDED)
			} else if node_run_result.Status == models.WorkflowNodeExecutionStatus_EXCEPTION {
				workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_EXCEPTION)
				workflow_node_execution.Error = node_run_result.Error
			}
		} else {
			// create workflow node execution
			workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_FAILED)
			workflow_node_execution.Error = errmsg
		}
		dbengine.Instance().DB.Create(workflow_node_execution)
	}()
	node_instance, generator := (&coreworkflow.WorkflowEntry{}).SingleStepRun(
		draft_workflow,
		node_id,
		account.ID,
		user_inputs,
	)

	for event := range generator {
		if real_event, ok := any(event).(*nodesevententities.RunCompletedEvent); ok {
			node_run_result = real_event.RunResult

			// sign output files
			node_run_result.Outputs = (&coreworkflow.WorkflowEntry{}).HandleSpecialValues(node_run_result.Outputs)
			break
		}
	}
	mlog.Debugf("------node_run_result.Inputs=%#v", node_run_result.Inputs)
	if node_run_result == nil {
		mlog.Errorf("Node run failed with no run result")
		panic(exceptions.NewValueError("Node run failed with no run result"))
	}
	// single step debug mode error handling return
	if node_run_result.Status == models.WorkflowNodeExecutionStatus_FAILED && node_instance.ShouldContinueOnError() {
		node_run_result = &workflowentities.NodeRunResult{
			Status: models.WorkflowNodeExecutionStatus_EXCEPTION,
			Error:  node_run_result.Error,
			Inputs: node_run_result.Inputs,
			Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
				workflowenumtypes.NodeRunMetadataKey_ERROR_STRATEGY: node_instance.GetBaseNodeData().ErrorStrategy,
			},
		}

		if node_instance.GetBaseNodeData().ErrorStrategy == nodesenumtypes.ErrorStrategy_DEFAULT_VALUE {
			node_run_result.Outputs = maps.Clone(node_instance.GetBaseNodeData().DefaultValueDict())

			node_run_result.Outputs["error_message"] = node_run_result.Error
			node_run_result.Outputs["error_type"] = node_run_result.ErrorType

		} else {
			node_run_result.Outputs = map[string]any{
				"error_message": node_run_result.Error,
				"error_type":    node_run_result.ErrorType,
			}
		}
	}
	run_succeeded = slices.Contains([]models.WorkflowNodeExecutionStatus{models.WorkflowNodeExecutionStatus_SUCCEEDED, models.WorkflowNodeExecutionStatus_EXCEPTION}, node_run_result.Status)
	if !run_succeeded {
		errmsg = node_run_result.Error
	} else {
		errmsg = ""
	}

	// except WorkflowNodeRunFailedError as e:
	// 	node_instance = e.node_instance
	// 	run_succeeded = False
	// 	node_run_result = None
	// 	error = e.error
	now := time.Now()
	workflow_node_execution = &models.WorkflowNodeExecution{
		ID:            uuid.NewV4().String(),
		TenantID:      app_model.TenantID,
		AppID:         app_model.ID,
		WorkflowID:    draft_workflow.ID,
		TriggeredFrom: string(models.WorkflowNodeExecutionTriggeredFrom_SINGLE_STEP),
		Index:         1,
		NodeID:        node_id,
		ElapsedTime:   time.Since(start_at).Seconds(),
		CreatedAt:     &now,
		CreatedByRole: models.CreatedByRole_ACCOUNT,
		CreatedBy:     account.ID,
		FinishedAt:    &now,
	}
	if node_instance != nil {
		workflow_node_execution.NodeType = node_instance.Type()
		workflow_node_execution.Title = node_instance.GetBaseNodeData().Title
	}

	return workflow_node_execution
}
func (s *WorkflowService) PublishWorkflow(app *models.App, account *models.Account, draft_workflow *models.Workflow) *models.Workflow {
	if draft_workflow == nil {
		draft_workflow = s.GetDraftWorkflow(app)
	}
	if draft_workflow == nil {
		panic(exceptions.NewValueError("No valid workflow found."))
	}
	now := time.Now()
	// create new workflow
	wf := &models.Workflow{
		ID:                       uuid.NewV4().String(),
		TenantID:                 app.TenantID,
		AppID:                    app.ID,
		Type:                     draft_workflow.Type,
		Version:                  now.Format("2006-01-02 15:04:05.000000"),
		Graph:                    draft_workflow.Graph,
		FeaturesStr:              draft_workflow.FeaturesStr,
		CreatedBy:                account.ID,
		CreatedAt:                &now,
		UpdatedBy:                account.ID,
		UpdatedAt:                &now,
		EnvironmentVariablesStr:  draft_workflow.EnvironmentVariablesStr,
		ConversationVariablesStr: draft_workflow.ConversationVariablesStr,
		// tenant_id=app_model.tenant_id,
		// app_id=app_model.id,
		// type=draft_workflow.type,
		// version=str(datetime.now(UTC).replace(tzinfo=None)),
		// graph=draft_workflow.graph,
		// features=draft_workflow.features,
		// created_by=account.id,
		// environment_variables=draft_workflow.environment_variables,
		// conversation_variables=draft_workflow.conversation_variables,
	}
	dbengine.Instance().DB.Create(wf)

	app.WorkflowID = wf.ID
	dbengine.Instance().DB.Updates(&models.App{ID: app.ID, WorkflowID: wf.ID})

	// trigger app workflow events
	events.Instance.AppPublishedWorkflowWasUpdatedSig.Emit(app, wf)

	// return new workflow
	return wf
}

func (s *WorkflowService) RestorePublishedWorkflowToDraft(appModel *models.App, workflowID string, account *models.Account) (*models.Workflow, error) {
	// Find the published workflow by ID
	var published models.Workflow
	if err := dbengine.Instance().DB.Where("id = ? AND app_id = ?", workflowID, appModel.ID).First(&published).Error; err != nil {
		return nil, fmt.Errorf("published workflow not found")
	}

	// Get or create draft
	draft := s.GetDraftWorkflow(appModel)
	if draft == nil {
		return nil, fmt.Errorf("draft workflow not found")
	}

	// Copy published content to draft
	draft.Graph = published.Graph
	draft.FeaturesStr = published.FeaturesStr
	draft.EnvironmentVariablesStr = published.EnvironmentVariablesStr
	draft.ConversationVariablesStr = published.ConversationVariablesStr
	draft.UpdatedBy = account.ID

	now := time.Now()
	draft.UpdatedAt = &now

	if err := dbengine.Instance().DB.Save(draft).Error; err != nil {
		return nil, err
	}
	return draft, nil
}

func (s *WorkflowService) GetAllPublishedWorkflow(appModel *models.App, page, limit int, userID string, namedOnly bool) ([]*models.Workflow, bool) {
	var workflows []*models.Workflow
	query := dbengine.Instance().DB.Where("app_id = ? AND version != ?", appModel.ID, "draft")
	if userID != "" {
		query = query.Where("created_by = ?", userID)
	}
	if namedOnly {
		query = query.Where("marked_name != ''")
	}

	query.Order("version DESC").Offset((page - 1) * limit).Limit(limit + 1).Find(&workflows)

	hasMore := len(workflows) > limit
	if hasMore {
		workflows = workflows[:limit]
	}
	return workflows, hasMore
}

func (s *WorkflowService) UpdateWorkflow(workflowID string, tenantID string, appID string, graph map[string]any, features map[string]any, uniqueHash string) (*models.Workflow, error) {
	var workflow models.Workflow
	if err := dbengine.Instance().DB.Where("id = ? AND tenant_id = ? AND app_id = ?", workflowID, tenantID, appID).First(&workflow).Error; err != nil {
		return nil, fmt.Errorf("workflow not found")
	}

	graphBytes, _ := json.Marshal(graph)
	featuresBytes, _ := json.Marshal(features)

	updates := map[string]any{
		"graph":       string(graphBytes),
		"features":    string(featuresBytes),
		"unique_hash": uniqueHash,
		"updated_at":  time.Now(),
	}

	if err := dbengine.Instance().DB.Model(&workflow).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (s *WorkflowService) DeleteWorkflow(workflowID string, tenantID string) error {
	return dbengine.Instance().DB.Where("id = ? AND tenant_id = ?", workflowID, tenantID).Delete(&models.Workflow{}).Error
}

func (s *WorkflowService) GetNodeLastRun(appModel *models.App, workflow *models.Workflow, nodeID string) *models.WorkflowNodeExecution {
	var execution models.WorkflowNodeExecution
	if err := dbengine.Instance().DB.Where(
		"app_id = ? AND workflow_id = ? AND node_id = ? AND triggered_from = ?",
		appModel.ID, workflow.ID, nodeID, "single-step",
	).Order("created_at DESC").First(&execution).Error; err != nil {
		return nil
	}
	return &execution
}

func (s *WorkflowService) IsWorkflowExist(appModel *models.App) bool {
	var count int64
	dbengine.Instance().DB.Model(&models.Workflow{}).Where("app_id = ?", appModel.ID).Count(&count)
	return count > 0
}

func (s *WorkflowService) ValidateGraphStructure(graph map[string]any) error {
	nodes, ok := graph["nodes"]
	if !ok {
		return fmt.Errorf("graph must contain nodes")
	}
	edges, ok := graph["edges"]
	if !ok {
		return fmt.Errorf("graph must contain edges")
	}

	nodeList, ok := nodes.([]any)
	if !ok || len(nodeList) == 0 {
		return fmt.Errorf("graph must have at least one node")
	}

	_, ok = edges.([]any)
	if !ok {
		return fmt.Errorf("edges must be an array")
	}

	return nil
}

// === Human Input Form Support ===

func (s *WorkflowService) GetHumanInputFormPreview(appModel *models.App, nodeID string) (map[string]any, error) {
	draft := s.GetDraftWorkflow(appModel)
	if draft == nil {
		return nil, fmt.Errorf("draft workflow not found")
	}
	// Parse graph to find the human input node
	var graph map[string]any
	json.Unmarshal([]byte(draft.Graph), &graph)
	nodes, _ := graph["nodes"].([]any)
	for _, n := range nodes {
		node, _ := n.(map[string]any)
		if nodeData, ok := node["data"].(map[string]any); ok {
			if id, _ := node["id"].(string); id == nodeID {
				return map[string]any{
					"node_id":   nodeID,
					"node_type": node["type"],
					"form_data": nodeData,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("node %s not found", nodeID)
}

func (s *WorkflowService) GetPublishedWorkflowByID(appModel *models.App, workflowID string) *models.Workflow {
	var workflow models.Workflow
	if err := dbengine.Instance().DB.Where("id = ? AND app_id = ? AND version != ?", workflowID, appModel.ID, "draft").First(&workflow).Error; err != nil {
		return nil
	}
	return &workflow
}

// === Credential Validation ===

func (s *WorkflowService) ValidateWorkflowCredentials(workflow *models.Workflow) error {
	if workflow == nil {
		return fmt.Errorf("workflow is nil")
	}
	var graph map[string]any
	if err := json.Unmarshal([]byte(workflow.Graph), &graph); err != nil {
		return fmt.Errorf("invalid graph JSON: %w", err)
	}
	nodes, _ := graph["nodes"].([]any)
	for _, n := range nodes {
		node, _ := n.(map[string]any)
		nodeType, _ := node["type"].(string)
		if nodeType == "llm" {
			data, _ := node["data"].(map[string]any)
			provider, _ := data["model"].(map[string]any)["provider"].(string)
			model, _ := data["model"].(map[string]any)["name"].(string)
			if provider == "" || model == "" {
				nodeID, _ := node["id"].(string)
				return fmt.Errorf("LLM node %s has no model configured", nodeID)
			}
		}
	}
	return nil
}

func (s *WorkflowService) ConvertToWorkflow(appModel *models.App, account *models.Account) (*models.App, error) {
	if appModel.Mode == models.AppMode_WORKFLOW || appModel.Mode == models.AppMode_ADVANCED_CHAT {
		return nil, fmt.Errorf("app is already a workflow type")
	}
	// Update app mode
	dbengine.Instance().DB.Model(appModel).Update("mode", "workflow")
	return appModel, nil
}
