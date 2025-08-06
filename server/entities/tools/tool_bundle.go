package tools

type ApiToolBundle[T1 float64 | int | string, T2 float64 | int] struct {
	// """
	// This class is used to store the schema information of an api based tool.
	//  such as the url, the method, the parameters, etc.
	// """

	// # server_url
	ServerURL string `json:"server_url"`
	// # method
	Method string `json:"method"`
	// # summary
	Summary *string `json:"summary"`
	// # operation_id
	OperationID *string `json:"operation_id"`
	// # parameters
	Parameters []*ToolParameter[T1, T2] `json:"parameters"`
	// # author
	Author string `json:"author"`
	// # icon
	Icon *string `json:"icon"`
	// # openapi operation
	Openapi map[string]any `json:"openapi"`
}
