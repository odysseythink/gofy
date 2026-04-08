package dashscope

type ResponseToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ResponseToolCall struct {
	Index    int                      `json:"index"`
	ID       string                   `json:"id"`
	Type     string                   `json:"type"`
	Function ResponseToolCallFunction `json:"function"`
}

type ResponseMessage struct {
	Role             string              `json:"role"`
	Content          any                 `json:"content"`
	ReasoningContent string              `json:"reasoning_content"`
	ToolCalls        []*ResponseToolCall `json:"tool_calls"`
}

type ResponsChoice struct {
	FinishReason string           `json:"finish_reason"`
	Message      *ResponseMessage `json:"message"`
}

type ResponseOutput struct {
	Text         string           `json:"text"`
	FinishReason string           `json:"finish_reason"`
	Choices      []*ResponsChoice `json:"choices"`
}

type ResponseUsage struct {
	InputTokens        int `json:"input_tokens"`
	InputTokensDetails struct {
		TextTokens  int `json:"text_tokens"`
		VideoTokens int `json:"video_tokens"`
		ImageTokens int `json:"image_tokens"`
	} `json:"input_tokens_details"`
	OutputTokens        int `json:"output_tokens"`
	OutputTokensDetails struct {
		TextTokens int `json:"text_tokens"`
	} `json:"output_tokens_details"`
	TotalTokens         int `json:"total_tokens"`
	ImageTokens         int `json:"image_tokens"`
	VideoTokens         int `json:"video_tokens"`
	AudioTokens         int `json:"audio_tokens"`
	PromptTokensDetails struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
}
type Response struct {
	StatusCode int             `json:"status_code"`
	RequestID  string          `json:"request_id"`
	Code       string          `json:"code"`
	Message    string          `json:"message"`
	Output     *ResponseOutput `json:"output"`
	Usage      *ResponseUsage  `json:"usage"`
}
