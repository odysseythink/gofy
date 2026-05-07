package manageres

import (
	"github.com/odysseythink/gofy/backend/core/manageres/conversation"
	datasetmanager "github.com/odysseythink/gofy/backend/core/manageres/dataset"
	modelmanager "github.com/odysseythink/gofy/backend/core/manageres/model_manager"
	workflownodeexecutionmanager "github.com/odysseythink/gofy/backend/core/manageres/workflow_node_execution"
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
