package modelruntime

type LLMMode string

const (
	/*
	   Enum class for large language model mode.
	*/

	LLMMode_COMPLETION LLMMode = "completion"
	LLMMode_CHAT       LLMMode = "chat"
)

type LLMUsage struct {
	/*
	   Model class for llm usage.
	*/

	PromptTokens        int     `json:"prompt_tokens"`
	PromptUnitPrice     float64 `json:"prompt_unit_price"`
	PromptPriceUnit     float64 `json:"prompt_price_unit"`
	PromptPrice         float64 `json:"prompt_price"`
	CompletionTokens    int     `json:"completion_tokens"`
	CompletionUnitPrice float64 `json:"completion_unit_price"`
	CompletionPriceUnit float64 `json:"completion_price_unit"`
	CompletionPrice     float64 `json:"completion_price"`
	TotalTokens         int     `json:"total_tokens"`
	TotalPrice          float64 `json:"total_price"`
	Currency            string  `json:"currency"`
	Latency             float64 `json:"latency"`
}

func NewLLMUsage() *LLMUsage {
	return &LLMUsage{
		PromptTokens:        0,
		PromptUnitPrice:     0.0,
		PromptPriceUnit:     0.0,
		PromptPrice:         0.0,
		CompletionTokens:    0,
		CompletionUnitPrice: 0.0,
		CompletionPriceUnit: 0.0,
		CompletionPrice:     0.0,
		TotalTokens:         0,
		TotalPrice:          0.0,
		Currency:            "USD",
		Latency:             0,
	}
}

// func (lu *LLMUsage) empty_usage()
//         return cls(
//             prompt_tokens=0,
//             prompt_unit_price=Decimal("0.0"),
//             prompt_price_unit=Decimal("0.0"),
//             prompt_price=Decimal("0.0"),
//             completion_tokens=0,
//             completion_unit_price=Decimal("0.0"),
//             completion_price_unit=Decimal("0.0"),
//             completion_price=Decimal("0.0"),
//             total_tokens=0,
//             total_price=Decimal("0.0"),
//             currency="USD",
//             latency=0.0,
//         )

func (lu *LLMUsage) Plus(other *LLMUsage) *LLMUsage {
	/*
	   Add two LLMUsage instances together.

	   :param other: Another LLMUsage instance to add
	   :return: A new LLMUsage instance with summed values
	*/
	if other == nil {
		return lu
	}
	if lu.TotalTokens == 0 {
		return other
	} else {
		return &LLMUsage{
			PromptTokens:        lu.PromptTokens + other.PromptTokens,
			PromptUnitPrice:     other.PromptUnitPrice,
			PromptPriceUnit:     other.PromptPriceUnit,
			PromptPrice:         lu.PromptPrice + other.PromptPrice,
			CompletionTokens:    lu.CompletionTokens + other.CompletionTokens,
			CompletionUnitPrice: other.CompletionUnitPrice,
			CompletionPriceUnit: other.CompletionPriceUnit,
			CompletionPrice:     lu.CompletionPrice + other.CompletionPrice,
			TotalTokens:         lu.TotalTokens + other.TotalTokens,
			TotalPrice:          lu.TotalPrice + other.TotalPrice,
			Currency:            other.Currency,
			Latency:             lu.Latency + other.Latency,
		}
	}
}

// def __add__(self, other: "LLMUsage") -> "LLMUsage":
//     """
//     Overload the + operator to add two LLMUsage instances.

//     :param other: Another LLMUsage instance to add
//     :return: A new LLMUsage instance with summed values
//     """
//     return lu.plus(other)

// LLMResult represents the result of an LLM.
type LLMResult struct {
	ID                string `json:"id"`
	Model             string `json:"model"`
	PromptMessages    []PromptMessager
	Message           *AssistantPromptMessage `json:"message"`
	Usage             *LLMUsage               `json:"usage"`
	SystemFingerprint string                  `json:"system_fingerprint"`
}

// LLMResultChunkDelta represents a delta for an LLM result chunk.
type LLMResultChunkDelta struct {
	Index        int
	Message      *AssistantPromptMessage `json:"message"`
	Usage        *LLMUsage               `json:"usage"`
	FinishReason string                  `json:"finish_reason"`
}

// LLMResultChunk represents a chunk of an LLM result.
type LLMResultChunk struct {
	Model             string
	PromptMessages    []PromptMessager
	SystemFingerprint string `json:"system_fingerprint"`
	Delta             *LLMResultChunkDelta
}

// NumTokensResult represents the result of the number of tokens.
type NumTokensResult struct {
	*PriceInfo
	Tokens int
}
