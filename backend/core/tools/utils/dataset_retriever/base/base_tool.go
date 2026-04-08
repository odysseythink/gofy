package base

import (
	"encoding/json"

	idxtoolcbhandler "mlib.com/gofy/server/core/callback_handler/index_tool"
	"mlib.com/mlog"
)

type DatasetRetrieverToolor interface {
	Run(query string) string
	GetName() string
	SetName(string)
	GetDescription() string
	SetDescription(string)
	GetTenantID() string
	SetTenantID(string)
	GetTopK() int
	SetTopK(int)
	GetScoreThreshold() float64
	SetScoreThreshold(float64)
	GetHitCallbacks() []*idxtoolcbhandler.DatasetIndexToolCallbackHandler
	SetHitCallbacks([]*idxtoolcbhandler.DatasetIndexToolCallbackHandler)
	GetReturnResource() bool
	SetReturnResource(bool)
	GetRetrieverFrom() string
	SetRetrieverFrom(string)
	GetModelConfig() map[string]any
	SetModelConfig(map[string]any)
}
type DatasetRetrieverBaseTool struct {
	Name           string                                              `json:"name"`
	Description    string                                              `json:"description"`
	TenantID       string                                              `json:"tenant_id"`
	TopK           int                                                 `json:"top_k"`
	ScoreThreshold float64                                             `json:"score_threshold"`
	HitCallbacks   []*idxtoolcbhandler.DatasetIndexToolCallbackHandler `json:"hit_callbacks"`
	ReturnResource bool                                                `json:"return_resource"`
	RetrieverFrom  string                                              `json:"retriever_from"`
	ModelConfig    map[string]any                                      `json:"model_config"`
}

func NewDatasetRetrieverBaseTool(args map[string]any) *DatasetRetrieverBaseTool {
	if len(args) == 0 {
		return &DatasetRetrieverBaseTool{
			Name:         "dataset",
			Description:  "use this to retrieve a dataset. ",
			TopK:         2,
			HitCallbacks: make([]*idxtoolcbhandler.DatasetIndexToolCallbackHandler, 0),
		}
	} else {
		ret := new(DatasetRetrieverBaseTool)
		bindata, _ := json.Marshal(args)
		if err := json.Unmarshal(bindata, &ret); err != nil {
			mlog.Errorf("json unmarshal %s to DatasetRetrieverBaseTool failed:%v", bindata, err)
			return nil
		}
		return ret
	}
}

func (tool *DatasetRetrieverBaseTool) GetName() string {
	return tool.Name
}
func (tool *DatasetRetrieverBaseTool) SetName(val string) {
	tool.Name = val
}
func (tool *DatasetRetrieverBaseTool) GetDescription() string {
	return tool.Description
}
func (tool *DatasetRetrieverBaseTool) SetDescription(val string) {
	tool.Description = val
}
func (tool *DatasetRetrieverBaseTool) GetTenantID() string {
	return tool.TenantID
}
func (tool *DatasetRetrieverBaseTool) SetTenantID(val string) {
	tool.TenantID = val
}
func (tool *DatasetRetrieverBaseTool) GetTopK() int {
	return tool.TopK
}
func (tool *DatasetRetrieverBaseTool) SetTopK(val int) {
	tool.TopK = val
}
func (tool *DatasetRetrieverBaseTool) GetScoreThreshold() float64 {
	return tool.ScoreThreshold
}
func (tool *DatasetRetrieverBaseTool) SetScoreThreshold(val float64) {
	tool.ScoreThreshold = val
}
func (tool *DatasetRetrieverBaseTool) GetHitCallbacks() []*idxtoolcbhandler.DatasetIndexToolCallbackHandler {
	return tool.HitCallbacks
}
func (tool *DatasetRetrieverBaseTool) SetHitCallbacks(val []*idxtoolcbhandler.DatasetIndexToolCallbackHandler) {
	tool.HitCallbacks = val
}
func (tool *DatasetRetrieverBaseTool) GetReturnResource() bool {
	return tool.ReturnResource
}
func (tool *DatasetRetrieverBaseTool) SetReturnResource(val bool) {
	tool.ReturnResource = val
}
func (tool *DatasetRetrieverBaseTool) GetRetrieverFrom() string {
	return tool.RetrieverFrom
}
func (tool *DatasetRetrieverBaseTool) SetRetrieverFrom(val string) {
	tool.RetrieverFrom = val
}
func (tool *DatasetRetrieverBaseTool) GetModelConfig() map[string]any {
	return tool.ModelConfig
}
func (tool *DatasetRetrieverBaseTool) SetModelConfig(val map[string]any) {
	tool.ModelConfig = val
}
