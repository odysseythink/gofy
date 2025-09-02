package condition

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
		"in",
		"not in",
		"all of",
		// for number
		"=",
		"≠",
		">",
		"<",
		"≥",
		"≤",
		"null",
		"not null",
		// for file
		"exists",
		"not exists",
	}
)

type LogicalOperatorType string

const (
	LogicalOperator_OR  LogicalOperatorType = "or"
	LogicalOperator_AND LogicalOperatorType = "and"
)

type ConditionValueType interface {
	~string | []string | any
}

type SubCondition struct {
	Key                string
	ComparisonOperator string
	Value              any
}
type SubVariableCondition struct {
	LogicalOperator string          //: Literal["and", "or"]
	Conditions      []*SubCondition // = Field(default=list)

}
type Condition struct {
	ID                   string                `json:"id"`
	VariableSelector     []string              `json:"variable_selector"`
	ComparisonOperator   string                `json:"comparison_operator"`
	Value                any                   `json:"value"`
	SubVariableCondition *SubVariableCondition `json:"sub_variable_condition"`
}
