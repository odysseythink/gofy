package rag

import (
	"encoding/json"

	"github.com/odysseythink/mlog"
)

var (
	SupportedComparisonOperator = []string{
		// for string or array
		"contains",
		"not contains",
		"start with",
		"end with",
		"is",
		"is not",
		"empty",
		"not empty",
		// for number
		"=",
		"≠",
		">",
		"<",
		"≥",
		"≤",
		// for time
		"before",
		"after",
	}
)

type Condition struct {
	Name               string `json:"name"`
	ComparisonOperator string `json:"comparison_operator"` //SupportedComparisonOperator
	Value              any    `json:"value"`               // [T string | []string | int | float64]
}

type MetadataCondition struct {
	LogicalOperator string       `json:"logical_operator"` // Optional[Literal["and", "or"]] = "and"
	Conditions      []*Condition `json:"conditions"`
}

func (data *MetadataCondition) ToDict() map[string]any {
	bindata, _ := json.Marshal(data)
	var ret map[string]any
	if err := json.Unmarshal(bindata, &ret); err != nil {
		mlog.Errorf("json unmarshal %s failed:%v", string(bindata), err)
		return nil
	}
	return ret
}
