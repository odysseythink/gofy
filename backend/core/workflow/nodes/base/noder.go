package base

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/graph"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
	nodesentities "github.com/odysseythink/gofy/backend/entities/nodes"
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
)

type NodeBeaner interface {
	GetID() string
	SetID(string)
	GetTenantID() string
	SetTenantID(string)
	GetAppID() string
	SetAppID(string)
	GetWorkflowType() models.WorkflowType
	SetWorkflowType(models.WorkflowType)
	GetWorkflowID() string
	SetWorkflowID(string)
	GetGraphConfig() map[string]any
	SetGraphConfig(map[string]any)
	GetUserID() string
	SetUserID(string)
	GetUserFrom() models.UserFrom
	SetUserFrom(models.UserFrom)
	GetInvokeFrom() appenumtypes.InvokeFrom
	SetInvokeFrom(appenumtypes.InvokeFrom)
	GetWorkflowCallDepth() int
	SetWorkflowCallDepth(int)
	GetGraph() *graph.Graph
	SetGraph(*graph.Graph)
	GetGraphRuntimeState() *graphengineentities.GraphRuntimeState
	SetGraphRuntimeState(*graphengineentities.GraphRuntimeState)
	GetPreviousNodeID() string
	SetPreviousNodeID(string)
	GetNodeID() string
	SetNodeID(string)
	GetBaseNodeData() *basenodesentities.BaseNodeData
	GetNodeData() any
	ShouldContinueOnError() bool
	ShouldRetry() bool
}
type Noder interface {
	GetDefaultConfig(filters map[string]any) map[string]any
	Type() nodesenumtypes.NodeType
	Run() (*workflowentities.NodeRunResult, iter.Seq[any])
	RunIter(Noder) iter.Seq[any]
	NodeBeaner
}

type SpecificNoder[T nodesentities.GenericNodeData] interface {
	ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data T) map[string][]string
}
