package variable

type VariableType string

const (
	Variable_NUMBER  VariableType = "number"
	Variable_STRING  VariableType = "string"
	Variable_BOOLEAN VariableType = "boolean"
	Variable_OBJECT  VariableType = "object"
	Variable_SECRET  VariableType = "secret"
	Variable_FILE    VariableType = "file"

	Variable_ARRAY_ANY     VariableType = "array[any]"
	Variable_ARRAY_STRING  VariableType = "array[string]"
	Variable_ARRAY_NUMBER  VariableType = "array[number]"
	Variable_ARRAY_BOOLEAN VariableType = "array[boolean]"
	Variable_ARRAY_OBJECT  VariableType = "array[object]"
	Variable_ARRAY_FILE    VariableType = "array[file]"

	Variable_NONE  VariableType = "none"
	Variable_GROUP VariableType = "group"
)

func (s VariableType) Valid() bool {
	return s == Variable_NUMBER ||
		s == Variable_STRING ||
		s == Variable_OBJECT ||
		s == Variable_SECRET ||
		s == Variable_FILE ||
		s == Variable_ARRAY_ANY ||
		s == Variable_ARRAY_STRING ||
		s == Variable_ARRAY_NUMBER ||
		s == Variable_ARRAY_BOOLEAN ||
		s == Variable_ARRAY_OBJECT ||
		s == Variable_ARRAY_FILE ||
		s == Variable_NONE ||
		s == Variable_GROUP
}

var (
	ENVIRONMENT_VARIABLE_SUPPORTED_TYPES = []VariableType{Variable_STRING, Variable_NUMBER, Variable_SECRET}
)
