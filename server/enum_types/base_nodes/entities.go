package basenodes

type DefaultValueType string

const (
	DefaultValue_STRING       DefaultValueType = "string"
	DefaultValue_NUMBER       DefaultValueType = "number"
	DefaultValue_OBJECT       DefaultValueType = "object"
	DefaultValue_ARRAY_NUMBER DefaultValueType = "array[number]"
	DefaultValue_ARRAY_STRING DefaultValueType = "array[string]"
	DefaultValue_ARRAY_OBJECT DefaultValueType = "array[object]"
	DefaultValue_ARRAY_FILES  DefaultValueType = "array[file]"
)
