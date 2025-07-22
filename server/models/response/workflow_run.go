package response

import (
	"encoding/json"

	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

// WorkflowRunForLogResponse 对应 workflow_run_for_log_fields
type WorkflowRunForLogResponse struct {
	ID              string  `json:"id"`
	Version         string  `json:"version"`
	Status          string  `json:"status"`
	Error           string  `json:"error"`
	ElapsedTime     float64 `json:"elapsed_time"`
	TotalTokens     int     `json:"total_tokens"`
	TotalSteps      int     `json:"total_steps"`
	CreatedAt       int64   `json:"created_at"`
	FinishedAt      int64   `json:"finished_at"`
	ExceptionsCount int     `json:"exceptions_count"`
}

// WorkflowRunForListResponse 对应 workflow_run_for_list_fields
type WorkflowRunForListResponse struct {
	ID               string                 `json:"id"`
	SequenceNumber   int                    `json:"sequence_number"`
	Version          string                 `json:"version"`
	Status           string                 `json:"status"`
	ElapsedTime      float64                `json:"elapsed_time"`
	TotalTokens      int                    `json:"total_tokens"`
	TotalSteps       int                    `json:"total_steps"`
	CreatedByAccount *SimpleAccountResponse `json:"created_by_account"`
	CreatedAt        int64                  `json:"created_at"`
	FinishedAt       int64                  `json:"finished_at"`
	ExceptionsCount  int                    `json:"exceptions_count"`
	RetryIndex       int                    `json:"retry_index"`
}

// AdvancedChatWorkflowRunForListResponse 对应 advanced_chat_workflow_run_for_list_fields
type AdvancedChatWorkflowRunForListResponse struct {
	ID               string                 `json:"id"`
	ConversationID   string                 `json:"conversation_id"`
	MessageID        string                 `json:"message_id"`
	SequenceNumber   int                    `json:"sequence_number"`
	Version          string                 `json:"version"`
	Status           string                 `json:"status"`
	ElapsedTime      float64                `json:"elapsed_time"`
	TotalTokens      int                    `json:"total_tokens"`
	TotalSteps       int                    `json:"total_steps"`
	CreatedByAccount *SimpleAccountResponse `json:"created_by_account"`
	CreatedAt        int64                  `json:"created_at"`
	FinishedAt       int64                  `json:"finished_at"`
	ExceptionsCount  int                    `json:"exceptions_count"`
	RetryIndex       int                    `json:"retry_index"`
}

// AdvancedChatWorkflowRunPaginationResponse 对应 advanced_chat_workflow_run_pagination_fields
type AdvancedChatWorkflowRunPaginationResponse struct {
	Limit   int                                      `json:"limit"`
	HasMore bool                                     `json:"has_more"`
	Data    []AdvancedChatWorkflowRunForListResponse `json:"data"`
}

func NewAdvancedChatWorkflowRunPaginationResponse(args any) *AdvancedChatWorkflowRunPaginationResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(AdvancedChatWorkflowRunPaginationResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AdvancedChatWorkflowRunPaginationResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(AdvancedChatWorkflowRunPaginationResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AdvancedChatWorkflowRunPaginationResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if result, ok := args.(*InfiniteScrollPagination); ok && result != nil {
		rsp := &AdvancedChatWorkflowRunPaginationResponse{
			Limit:   result.Limit,
			HasMore: result.HasMore,
			Data:    make([]AdvancedChatWorkflowRunForListResponse, 0),
		}
		for _, v := range result.Data.([]*models.WorkflowWithMessage) {
			sv := AdvancedChatWorkflowRunForListResponse{
				ID:              v.ID,
				ConversationID:  v.ConversationID,
				MessageID:       v.MessageID,
				SequenceNumber:  v.SequenceNumber,
				Version:         v.Version,
				Status:          v.Status,
				ElapsedTime:     v.ElapsedTime,
				TotalTokens:     v.TotalTokens,
				TotalSteps:      v.TotalSteps,
				CreatedAt:       v.CreatedAt.Unix(),
				FinishedAt:      v.FinishedAt.Unix(),
				ExceptionsCount: v.ExceptionsCount,
			}
			created_by_account := v.CreatedByAccount()
			if created_by_account != nil {
				sv.CreatedByAccount = &SimpleAccountResponse{
					ID:    created_by_account.ID,
					Name:  created_by_account.Name,
					Email: created_by_account.Email,
				}
			}
			rsp.Data = append(rsp.Data, sv)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(AdvancedChatWorkflowRunPaginationResponse)
}

// WorkflowRunPaginationResponse 对应 workflow_run_pagination_fields
type WorkflowRunPaginationResponse struct {
	Limit   int                          `json:"limit"`
	HasMore bool                         `json:"has_more"`
	Data    []WorkflowRunForListResponse `json:"data"`
}

func NewWorkflowRunPaginationResponse(args any) *WorkflowRunPaginationResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(WorkflowRunPaginationResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunPaginationResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(WorkflowRunPaginationResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunPaginationResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if result, ok := args.(*InfiniteScrollPagination); ok && result != nil {
		rsp := &WorkflowRunPaginationResponse{
			Limit:   result.Limit,
			HasMore: result.HasMore,
			Data:    make([]WorkflowRunForListResponse, 0),
		}
		for _, v := range result.Data.([]*models.WorkflowRun) {
			sv := WorkflowRunForListResponse{
				ID:              v.ID,
				SequenceNumber:  v.SequenceNumber,
				Version:         v.Version,
				Status:          v.Status,
				ElapsedTime:     v.ElapsedTime,
				TotalTokens:     v.TotalTokens,
				TotalSteps:      v.TotalSteps,
				ExceptionsCount: v.ExceptionsCount,
			}
			if v.CreatedAt != nil {
				sv.CreatedAt = v.CreatedAt.Unix()
			}
			if v.FinishedAt != nil {
				sv.FinishedAt = v.FinishedAt.Unix()
			}
			created_by_account := v.CreatedByAccount()
			if created_by_account != nil {
				sv.CreatedByAccount = &SimpleAccountResponse{
					ID:    created_by_account.ID,
					Name:  created_by_account.Name,
					Email: created_by_account.Email,
				}
			}

			rsp.Data = append(rsp.Data, sv)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(WorkflowRunPaginationResponse)
}

// WorkflowRunDetailResponse 对应 workflow_run_detail_fields
type WorkflowRunDetailResponse struct {
	ID               string                 `json:"id"`
	SequenceNumber   int                    `json:"sequence_number"`
	Version          string                 `json:"version"`
	Graph            any                    `json:"graph"`  // 对应 fields.Raw
	Inputs           any                    `json:"inputs"` // 对应 fields.Raw
	Status           string                 `json:"status"`
	Outputs          any                    `json:"outputs"` // 对应 fields.Raw
	Error            string                 `json:"error"`
	ElapsedTime      float64                `json:"elapsed_time"`
	TotalTokens      int                    `json:"total_tokens"`
	TotalSteps       int                    `json:"total_steps"`
	CreatedByRole    string                 `json:"created_by_role"`
	CreatedByAccount *SimpleAccountResponse `json:"created_by_account"`
	CreatedByEndUser *SimpleEndUserResponse `json:"created_by_end_user"`
	CreatedAt        int64                  `json:"created_at"`
	FinishedAt       int64                  `json:"finished_at"`
	ExceptionsCount  int                    `json:"exceptions_count"`
}

func NewWorkflowRunDetailResponse(args any) *WorkflowRunDetailResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(WorkflowRunDetailResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunDetailResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(WorkflowRunDetailResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunDetailResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if result, ok := args.(*models.WorkflowRun); ok && result != nil {
		rsp := &WorkflowRunDetailResponse{
			ID:              result.ID,
			SequenceNumber:  result.SequenceNumber,
			Version:         result.Version,
			Graph:           result.Graph,
			Inputs:          result.Inputs,
			Status:          result.Status,
			Outputs:         result.Outputs,
			Error:           result.Error,
			ElapsedTime:     result.ElapsedTime,
			TotalTokens:     result.TotalTokens,
			TotalSteps:      result.TotalSteps,
			CreatedByRole:   string(result.CreatedByRole),
			ExceptionsCount: result.ExceptionsCount,
		}
		if result.CreatedAt != nil {
			rsp.CreatedAt = result.CreatedAt.Unix()
		}
		if result.FinishedAt != nil {
			rsp.FinishedAt = result.FinishedAt.Unix()
		}
		created_by_account := result.CreatedByAccount()
		if created_by_account != nil {
			rsp.CreatedByAccount = &SimpleAccountResponse{
				ID:    created_by_account.ID,
				Name:  created_by_account.Name,
				Email: created_by_account.Email,
			}
		}
		created_by_end_user := result.CreatedByEndUser()
		if created_by_end_user != nil {
			rsp.CreatedByEndUser = &SimpleEndUserResponse{
				ID:          created_by_end_user.ID,
				Type:        created_by_end_user.Type,
				IsAnonymous: created_by_end_user.IsAnonymous,
				SessionID:   created_by_end_user.SessionID,
			}
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(WorkflowRunDetailResponse)
}

// RetryEventField 对应 retry_event_field
type RetryEventField struct {
	ElapsedTime float64 `json:"elapsed_time"`
	Status      string  `json:"status"`
	Inputs      any     `json:"inputs"`       // 对应 fields.Raw
	ProcessData any     `json:"process_data"` // 对应 fields.Raw
	Outputs     any     `json:"outputs"`      // 对应 fields.Raw
	Metadata    any     `json:"metadata"`     // 对应 fields.Raw
	LLMUsage    any     `json:"llm_usage"`    // 对应 fields.Raw
	Error       string  `json:"error"`
	RetryIndex  int     `json:"retry_index"`
}

// WorkflowRunNodeExecutionResponse 对应 workflow_run_node_execution_fields
type WorkflowRunNodeExecutionResponse struct {
	ID                string                 `json:"id"`
	Index             int                    `json:"index"`
	PredecessorNodeID string                 `json:"predecessor_node_id"`
	NodeID            string                 `json:"node_id"`
	NodeType          string                 `json:"node_type"`
	Title             string                 `json:"title"`
	Inputs            any                    `json:"inputs"`       // 对应 fields.Raw
	ProcessData       any                    `json:"process_data"` // 对应 fields.Raw
	Outputs           any                    `json:"outputs"`      // 对应 fields.Raw
	Status            string                 `json:"status"`
	Error             string                 `json:"error"`
	ElapsedTime       float64                `json:"elapsed_time"`
	ExecutionMetadata any                    `json:"execution_metadata"` // 对应 fields.Raw
	Extras            any                    `json:"extras"`             // 对应 fields.Raw
	CreatedAt         int64                  `json:"created_at"`
	CreatedByRole     string                 `json:"created_by_role"`
	CreatedByAccount  *SimpleAccountResponse `json:"created_by_account"`
	CreatedByEndUser  *SimpleEndUserResponse `json:"created_by_end_user"`
	FinishedAt        int64                  `json:"finished_at"`
}

func NewWorkflowRunNodeExecutionResponse(args any) *WorkflowRunNodeExecutionResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(WorkflowRunNodeExecutionResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunNodeExecutionResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(WorkflowRunNodeExecutionResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunNodeExecutionResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if workflow_node_execution, ok := args.(*models.WorkflowNodeExecution); ok && workflow_node_execution != nil {
		rsp := &WorkflowRunNodeExecutionResponse{
			ID:                workflow_node_execution.ID,
			Index:             workflow_node_execution.Index,
			PredecessorNodeID: workflow_node_execution.PredecessorNodeID,
			NodeID:            workflow_node_execution.NodeID,
			NodeType:          string(workflow_node_execution.NodeType),
			Title:             workflow_node_execution.Title,
			Inputs:            workflow_node_execution.InputsDict(),
			ProcessData:       workflow_node_execution.ProcessDataDict(),
			Outputs:           workflow_node_execution.OutputsDict(),
			Status:            workflow_node_execution.Status,
			Error:             workflow_node_execution.Error,
			ElapsedTime:       workflow_node_execution.ElapsedTime,
			ExecutionMetadata: workflow_node_execution.ExecutionMetadataDict(),
			Extras:            workflow_node_execution.Extras(),
			CreatedAt:         workflow_node_execution.CreatedAt.Unix(),
			CreatedByRole:     string(workflow_node_execution.CreatedByRole),
			CreatedByAccount:  nil,
			CreatedByEndUser:  nil,
			FinishedAt:        workflow_node_execution.FinishedAt.Unix(),
		}
		created_by_acc := workflow_node_execution.CreatedByAccount()
		if created_by_acc != nil {
			rsp.CreatedByAccount = &SimpleAccountResponse{
				ID:    created_by_acc.ID,
				Name:  created_by_acc.Name,
				Email: created_by_acc.Email,
			}
		}
		created_by_end_user := workflow_node_execution.CreatedByEndUser()
		if created_by_end_user != nil {
			rsp.CreatedByEndUser = &SimpleEndUserResponse{
				ID:          created_by_end_user.ID,
				Type:        created_by_end_user.Type,
				IsAnonymous: created_by_end_user.IsAnonymous,
				SessionID:   created_by_end_user.SessionID,
			}
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(WorkflowRunNodeExecutionResponse)
}

// WorkflowRunNodeExecutionListResponse 对应 workflow_run_node_execution_list_fields
type WorkflowRunNodeExecutionListResponse struct {
	Data []WorkflowRunNodeExecutionResponse `json:"data"`
}

func NewWorkflowRunNodeExecutionListResponse(args any) *WorkflowRunNodeExecutionListResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(WorkflowRunNodeExecutionListResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunNodeExecutionListResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(WorkflowRunNodeExecutionListResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to WorkflowRunNodeExecutionListResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if results, ok := args.([]*models.WorkflowNodeExecution); ok && len(results) > 0 {
		rsp := &WorkflowRunNodeExecutionListResponse{
			Data: make([]WorkflowRunNodeExecutionResponse, 0),
		}
		for _, v := range results {
			sv := WorkflowRunNodeExecutionResponse{
				ID:                v.ID,
				Index:             v.Index,
				PredecessorNodeID: v.PredecessorNodeID,
				NodeID:            v.NodeID,
				NodeType:          string(v.NodeType),
				Title:             v.Title,
				Inputs:            v.Inputs,
				ProcessData:       v.ProcessData,
				Outputs:           v.Outputs,
				Status:            v.Status,
				Error:             v.Error,
				ElapsedTime:       v.ElapsedTime,
				ExecutionMetadata: v.ExecutionMetadata,
				Extras:            v.Extras(),
				CreatedAt:         v.CreatedAt.Unix(),
				CreatedByRole:     string(v.CreatedByRole),
				// CreatedByAccount:  v.CreatedByAccount,
				// CreatedByEndUser:  v.CreatedByEndUser,
				FinishedAt: v.FinishedAt.Unix(),
			}
			created_by_account := v.CreatedByAccount()
			if created_by_account != nil {
				sv.CreatedByAccount = &SimpleAccountResponse{
					ID:    created_by_account.ID,
					Name:  created_by_account.Name,
					Email: created_by_account.Email,
				}
			}
			created_by_end_user := v.CreatedByEndUser()
			if created_by_end_user != nil {
				sv.CreatedByEndUser = &SimpleEndUserResponse{
					ID:          created_by_end_user.ID,
					Type:        created_by_end_user.Type,
					IsAnonymous: created_by_end_user.IsAnonymous,
					SessionID:   created_by_end_user.SessionID,
				}
			}
			rsp.Data = append(rsp.Data, sv)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(WorkflowRunNodeExecutionListResponse)
}
