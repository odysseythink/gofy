package condition

type SupportedComparisonOperator string

const (
	// for string or array
	OperatorContains    SupportedComparisonOperator = "contains"
	OperatorNotContains SupportedComparisonOperator = "not contains"
	OperatorStartWith   SupportedComparisonOperator = "start with"
	OperatorEndWith     SupportedComparisonOperator = "end with"
	OperatorIs          SupportedComparisonOperator = "is"
	OperatorIsNot       SupportedComparisonOperator = "is not"
	OperatorEmpty       SupportedComparisonOperator = "empty"
	OperatorNotEmpty    SupportedComparisonOperator = "not empty"
	OperatorIn          SupportedComparisonOperator = "in"
	OperatorNotIn       SupportedComparisonOperator = "not in"
	OperatorAllOf       SupportedComparisonOperator = "all of"
	// for number
	OperatorEqual              SupportedComparisonOperator = "="
	OperatorNotEqual           SupportedComparisonOperator = "≠"
	OperatorGreaterThan        SupportedComparisonOperator = ">"
	OperatorLessThan           SupportedComparisonOperator = "<"
	OperatorGreaterThanOrEqual SupportedComparisonOperator = "≥"
	OperatorLessThanOrEqual    SupportedComparisonOperator = "≤"
	OperatorNull               SupportedComparisonOperator = "null"
	OperatorNotNull            SupportedComparisonOperator = "not null"
	// for file
	OperatorExists    SupportedComparisonOperator = "exists"
	OperatorNotExists SupportedComparisonOperator = "not exists"
)
