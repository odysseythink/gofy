package base

import (
	"iter"

	"mlib.com/gofy/server/core/workflow/graph"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	nodesentities "mlib.com/gofy/server/entities/nodes"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
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
	ExtractVarSelectorToVarMapping(
		Noder,
		graph_config map[string]any,
		config map[string]any,
	) map[string][]string
	ExtractVarSelectorToVarMapping1(
		graph_config map[string]any,
		node_id string,
		node_data map[string]any,
	) map[string][]string
	NodeBeaner
}

type SpecificNoder[T nodesentities.GenericNodeData] interface {
	ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data T) map[string][]string
}
