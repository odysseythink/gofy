package workflow

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
)

type WorkflowNodeRunFailedError struct {
	NodeInstance base.Noder
	Errmsg       string
}

func NewWorkflowNodeRunFailedError(node_instance base.Noder, errmsg string) *WorkflowNodeRunFailedError {
	return &WorkflowNodeRunFailedError{
		NodeInstance: node_instance,
		Errmsg:       errmsg,
	}
}

func (e *WorkflowNodeRunFailedError) Error() string {
	return fmt.Sprintf("Node %s run failed: %s", e.NodeInstance.GetBaseNodeData().Title, e.Errmsg)
}
