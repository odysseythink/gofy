package ifelse

import (
	"encoding/json"

	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	conditionentities "mlib.com/gofy/server/entities/workflow/condition"
)

type Case struct {
	CaseID          string                         `json:"case_id"`
	LogicalOperator string                         `json:"logical_operator"` //: Literal["and", "or"]
	Conditions      []*conditionentities.Condition `json:"conditions"`
}

func (c *Case) ModelDump() map[string]any {
	bindata, _ := json.Marshal(c)
	rsp := map[string]any{}
	json.Unmarshal(bindata, &rsp)
	return rsp
}

// CodeNodeData represents answer node data
type IfElseNodeData struct {
	*basenodesentities.BaseNodeData
	LogicalOperator string                         `json:"logical_operator"` //: Literal["and", "or"]
	Conditions      []*conditionentities.Condition `json:"conditions"`
	Cases           []*Case                        `json:"cases"`
}

func New() *IfElseNodeData {
	return &IfElseNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
	}
}
