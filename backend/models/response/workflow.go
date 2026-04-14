package response

import (
	"encoding/json"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/models"
)

// type EnvironmentVariableField any

// func (ef EnvironmentVariableField) MarshalJSON() ([]byte, error) {
// 	if ef == nil {
// 		return nil, errors.New("nil EnvironmentVariableField")
// 	}
// 	// Mask secret variables values in environment_variables

// 	if _, ok := ef.(*variables.Variable); ok {
// 		dic := map[string]any{
// 			"id":         ef.(*variables.Variable).ID,
// 			"name":       ef.(*variables.Variable).Name,
// 			"value":      ef.(*variables.Variable).Value,
// 			"value_type": ef.(*variables.Variable).ValueType,
// 		}
// 		return json.Marshal(dic)
// 	} else if _, ok := ef.(variables.Variable); ok {
// 		dic := map[string]any{
// 			"id":         ef.(*variables.Variable).ID,
// 			"name":       ef.(*variables.Variable).Name,
// 			"value":      ef.(*variables.Variable).Value,
// 			"value_type": ef.(*variables.Variable).ValueType,
// 		}
// 		return json.Marshal(dic)
// 	} else if _, ok := ef.(map[string]any); ok && ef.(map[string]any) != nil {
// 		if _, ok := ef.(map[string]any)["value_type"]; ok {
// 			if _, ok := ef.(map[string]any)["value_type"].(string); ok {
// 				if !slices.Contains(enumtypes.ENVIRONMENT_VARIABLE_SUPPORTED_TYPES, enumtypes.VariableType(ef.(map[string]any)["value_type"].(string))) {
// 					return nil, fmt.Errorf("Unsupported environment variable value type: %v", ef.(map[string]any)["value_type"])
// 				}
// 			} else {
// 				return nil, fmt.Errorf("Unsupported environment variable value type: %v", ef.(map[string]any)["value_type"])
// 			}
// 		} else {
// 			return nil, errors.New("no environment variable value_type")
// 		}
// 		return json.Marshal(ef)
// 	} else if _, ok := ef.(*time.Time); ok && ef.(map[string]any) != nil {
// 		return []byte(ef.(*time.Time).Format(time.DateTime)), nil
// 	} else if _, ok := ef.(*time.Time); ok && ef.(map[string]any) != nil {
// 		return []byte(ef.(time.Time).Format(time.DateTime)), nil
// 	} else {
// 		return nil, fmt.Errorf("Unsupported EnvironmentVariableField(%#v)", ef)
// 	}
// }

type EnvironmentVariableResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ValueType string `json:"value_type"` //": fields.String(attribute="value_type.value"),
	Value     any    `json:"value"`      //": fields.Raw,
}

type ConversationVariableResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ValueType   string `json:"value_type"` //": fields.String(attribute="value_type.value"),
	Value       any    `json:"value"`      //": fields.Raw,
	Description string `json:"description"`
}

type WorkflowPartialFields struct {
	ID        string `json:"id"`
	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt int64  `json:"updated_at"`
}
type WorkflowResponse struct {
	ID                    string                          `json:"id"`
	Graph                 map[string]any                  `json:"graph"`      //": fields.Raw(attribute="graph_dict"),
	Features              map[string]any                  `json:"features"`   //": fields.Raw(attribute="features_dict"),
	Hash                  string                          `json:"hash"`       //: fields.String(attribute="unique_hash"),
	Version               string                          `json:"version"`    //: fields.String(attribute="version"),
	CreatedBy             *SimpleAccountResponse          `json:"created_by"` //": fields.Nested(simple_account_fields, attribute="created_by_account"),
	CreatedAt             int64                           `json:"created_at"`
	UpdatedBy             *SimpleAccountResponse          `json:"updated_by"` //": fields.Nested(simple_account_fields, attribute="updated_by_account", allow_null=True),
	UpdatedAt             int64                           `json:"updated_at"`
	ToolPublished         bool                            `json:"tool_published"`
	EnvironmentVariables  []*EnvironmentVariableResponse  `json:"environment_variables"`  //": fields.List(EnvironmentVariableField()),
	ConversationVariables []*ConversationVariableResponse `json:"conversation_variables"` // ": fields.List(fields.Nested(conversation_variable_fields)),
}

func NewWorkflowResponse(args any) *WorkflowResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(WorkflowResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(WorkflowResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if wf, ok := args.(*models.Workflow); ok && wf != nil {
		ret := &WorkflowResponse{
			ID:                    wf.ID,
			Graph:                 wf.GraphDict(),
			Features:              wf.FeaturesDict(),
			Hash:                  wf.UniqueHash(),
			Version:               wf.Version,
			CreatedBy:             nil,
			CreatedAt:             wf.CreatedAt.Unix(),
			UpdatedBy:             nil,
			UpdatedAt:             wf.UpdatedAt.Unix(),
			ToolPublished:         wf.ToolPublished(),
			EnvironmentVariables:  []*EnvironmentVariableResponse{},  //wf.EnvironmentVariables(),
			ConversationVariables: []*ConversationVariableResponse{}, //wf.ConversationVariables(),
		}
		if wf.CreatedByAccount != nil {
			ret.CreatedBy = &SimpleAccountResponse{
				ID:    wf.CreatedByAccount.ID,
				Name:  wf.CreatedByAccount.Name,
				Email: wf.CreatedByAccount.Email,
			}
		}
		if wf.UpdatedByAccount != nil {
			ret.UpdatedBy = &SimpleAccountResponse{
				ID:    wf.UpdatedByAccount.ID,
				Name:  wf.UpdatedByAccount.Name,
				Email: wf.UpdatedByAccount.Email,
			}
		}
		for _, v := range wf.GetEnvironmentVariables() {
			if ret.EnvironmentVariables == nil {
				ret.EnvironmentVariables = make([]*EnvironmentVariableResponse, 0)
			}
			ret.EnvironmentVariables = append(ret.EnvironmentVariables, &EnvironmentVariableResponse{
				ID:        v.GetID(),
				Name:      v.GetName(),
				ValueType: string(v.ValueType()),
				Value:     v.GetValue(),
			})
		}
		for _, v := range wf.GetConversationVariables() {
			mlog.Debugf("------conversation variable=%#v", v)
			if ret.ConversationVariables == nil {
				ret.ConversationVariables = make([]*ConversationVariableResponse, 0)
			}
			ret.ConversationVariables = append(ret.ConversationVariables, &ConversationVariableResponse{
				ID:          v.GetID(),
				Name:        v.GetName(),
				ValueType:   string(v.ValueType()),
				Value:       v.GetValue(),
				Description: v.GetDescription(),
			})
		}
		return ret
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(WorkflowResponse)
}

type InfiniteScrollPagination struct {
	Data    any  `json:"data"`
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
}

type WorkflowAppLogPartialResponse struct {
	ID               string                     `json:"id"`
	WorkflowRun      *WorkflowRunForLogResponse `json:"workflow_run"`
	CreatedFrom      string                     `json:"created_from"`
	CreatedByRole    string                     `json:"created_by_role"`
	CreatedByAccount *SimpleAccountResponse     `json:"created_by_account"`  //attribute="created_by_account"
	CreatedByEndUser *SimpleEndUserResponse     `json:"created_by_end_user"` //attribute="created_by_end_user"
	CreatedAt        int64                      `json:"created_at"`
}

type WorkflowAppLogPaginationResponse struct {
	Data    []*WorkflowAppLogPartialResponse `json:"data"`
	Limit   int                              `json:"limit"`
	Page    int                              `json:"page"`
	Total   int64                            `json:"total"`
	HasMore bool                             `json:"has_more"`
}

func NewWorkflowAppLogPaginationResponse(args any) *WorkflowAppLogPaginationResponse {
	rsp := new(WorkflowAppLogPaginationResponse)
	if real_args, ok := args.(string); ok && real_args != "" {
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowAppLogPaginationResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowAppLogPaginationResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return rsp
}
