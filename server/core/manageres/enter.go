package manageres

import (
	"mlib.com/gofy/server/core/manageres/conversation"
	datasetmanager "mlib.com/gofy/server/core/manageres/dataset"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	workflownodeexecutionmanager "mlib.com/gofy/server/core/manageres/workflow_node_execution"
)

type ManagerGroup struct {
	WorkflowNodeExecution *workflownodeexecutionmanager.WorkflowNodeExecutionManager
	Model                 *modelmanager.ModelManager
	Dataset               *datasetmanager.DatasetManager
	Conversation          *conversation.ConversationManager
}

var Instance = ManagerGroup{
	WorkflowNodeExecution: &workflownodeexecutionmanager.WorkflowNodeExecutionManager{},
	Model:                 modelmanager.NewModelManager(),
	Dataset:               &datasetmanager.DatasetManager{},
	Conversation:          conversation.New(),
}
