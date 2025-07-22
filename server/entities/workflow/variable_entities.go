package workflow

type VariableSelector struct {
	/*
	   Variable Selector.
	*/

	Variable      string   `json:"variable"`
	ValueSelector []string `json:"value_selector"`
}
