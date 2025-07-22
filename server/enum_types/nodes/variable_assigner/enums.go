package variableassigner

type OperationType string

const (
	Operation_OVER_WRITE OperationType = "over-write"
	Operation_CLEAR      OperationType = "clear"
	Operation_APPEND     OperationType = "append"
	Operation_EXTEND     OperationType = "extend"
	Operation_SET        OperationType = "set"
	Operation_ADD        OperationType = "+="
	Operation_SUBTRACT   OperationType = "-="
	Operation_MULTIPLY   OperationType = "*="
	Operation_DIVIDE     OperationType = "/="
)

type InputType string

const (
	Input_VARIABLE InputType = "variable"
	Input_CONSTANT InputType = "constant"
)
