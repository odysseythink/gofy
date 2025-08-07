package tools

type ApiToolBundle struct {
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
	Parameters []*ToolParameter `json:"parameters"`
	// # author
	Author string `json:"author"`
	// # icon
	Icon *string `json:"icon"`
	// # openapi operation
	Openapi map[string]any `json:"openapi"`
}
