package graphengine

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	conditionentities "github.com/odysseythink/gofy/backend/entities/workflow/condition"
)

type RunCondition struct {
	Type string //["branch_identify", "condition"]
	/*condition type*/

	BranchIdentify string
	/*branch identify like: sourceHandle, required when type is branch_identify*/

	Conditions []*conditionentities.Condition
	/*conditions to run the node, required when type is condition*/
}

func (rc *RunCondition) Hash() string {
	jsonData, err := json.Marshal(rc)
	if err != nil {
		fmt.Printf("Error marshaling RunCondition to JSON: %v\n", err)
		return ""
	}
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash)
}
