package rag

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

type Condition[T string | []string | int | float64] struct {
	Name               string `json:"name"`
	ComparisonOperator string `json:"comparison_operator"` //SupportedComparisonOperator
	Value              T      `json:"value"`
}

type MetadataCondition[T string | []string | int | float64] struct {
	LogicalOperator string          `json:"logical_operator"` // Optional[Literal["and", "or"]] = "and"
	Conditions      []*Condition[T] `json:"conditions"`
}
