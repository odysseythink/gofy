package response

import (
	"encoding/json"

	appresponserentities "mlib.com/gofy/server/entities/app/responser"
	"mlib.com/gofy/server/models"
)

type AppBlockingResponse struct {
	/*
	   AppBlockingResponse entity
	*/

	taskID string
}

func (rsp *AppBlockingResponse) TaskID() string {
	return rsp.taskID
}
func (rsp *AppBlockingResponse) ToDict(responser appresponserentities.AppBlockingResponser) map[string]any {
	bindata, _ := json.Marshal(responser)
	var rspmap map[string]any
	json.Unmarshal(bindata, &rspmap)
	rspmap["task_id"] = responser.TaskID()
	return rspmap
}

// ChatbotAppBlockingResponse entity
type ChatbotAppBlockingResponse struct {
	*AppBlockingResponse
	/*
	   ChatbotAppBlockingResponse entity
	*/
	Data struct {
		ID             string         `json:"id"`
		Mode           string         `json:"mode"`
		ConversationID string         `json:"conversation_id"`
		MessageID      string         `json:"message_id"`
		Answer         string         `json:"answer"`
		Metadata       map[string]any `json:"metadata"`
		CreatedAt      int64          `json:"created_at"`
	} `json:"data"`
}

func NewChatbotAppBlockingResponse(
	task_id, message_id, conversation_id, answer string, conversation_mode models.AppMode, created_at int64, metadata map[string]any,
) *ChatbotAppBlockingResponse {
	return &ChatbotAppBlockingResponse{
		AppBlockingResponse: &AppBlockingResponse{
			taskID: task_id,
		},
		Data: struct {
			ID             string         `json:"id"`
			Mode           string         `json:"mode"`
			ConversationID string         `json:"conversation_id"`
			MessageID      string         `json:"message_id"`
			Answer         string         `json:"answer"`
			Metadata       map[string]any `json:"metadata"`
			CreatedAt      int64          `json:"created_at"`
		}{
			ID:             message_id,
			Mode:           string(conversation_mode),
			MessageID:      message_id,
			Answer:         answer,
			CreatedAt:      created_at,
			Metadata:       metadata,
			ConversationID: conversation_id,
		},
	}
}

// CompletionAppBlockingResponse entity
type CompletionAppBlockingResponse struct {
	*AppBlockingResponse
	/*
	   CompletionAppBlockingResponse entity
	*/
	Data struct {
		ID        string         `json:"id"`
		Mode      string         `json:"mode"`
		MessageID string         `json:"message_id"`
		Answer    string         `json:"answer"`
		Metadata  map[string]any `json:"metadata"`
		CreatedAt int64          `json:"created_at"`
	} `json:"data"`
}

func NewCompletionAppBlockingResponse(
	task_id, message_id, answer string, conversation_mode models.AppMode, created_at int64, metadata map[string]any,
) *CompletionAppBlockingResponse {
	return &CompletionAppBlockingResponse{
		AppBlockingResponse: &AppBlockingResponse{
			taskID: task_id,
		},
		Data: struct {
			ID        string         `json:"id"`
			Mode      string         `json:"mode"`
			MessageID string         `json:"message_id"`
			Answer    string         `json:"answer"`
			Metadata  map[string]any `json:"metadata"`
			CreatedAt int64          `json:"created_at"`
		}{
			ID:        message_id,
			Mode:      string(conversation_mode),
			MessageID: message_id,
			Answer:    answer,
			CreatedAt: created_at,
			Metadata:  metadata,
		},
	}
}

// WorkflowAppBlockingResponse entity
type WorkflowAppBlockingResponse struct {
	*AppBlockingResponse
	/*
	   WorkflowAppBlockingResponse entity
	*/
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID          string         `json:"id"`
		WorkflowID  string         `json:"workflow_id"`
		Status      string         `json:"status"`
		Outputs     map[string]any `json:"outputs,omitempty"`
		Error       string         `json:"error,omitempty"`
		ElapsedTime float64        `json:"elapsed_time"`
		TotalTokens int            `json:"total_tokens"`
		TotalSteps  int            `json:"total_steps"`
		CreatedAt   int64          `json:"created_at"`
		FinishedAt  int64          `json:"finished_at"`
	} `json:"data"`
}

func NewWorkflowAppBlockingResponse(
	stream_response *WorkflowFinishStreamResponse, task_id string,
) *WorkflowAppBlockingResponse {
	return &WorkflowAppBlockingResponse{
		AppBlockingResponse: &AppBlockingResponse{
			taskID: task_id,
		},
		WorkflowRunID: stream_response.Data.ID,
		Data: struct {
			ID          string         `json:"id"`
			WorkflowID  string         `json:"workflow_id"`
			Status      string         `json:"status"`
			Outputs     map[string]any `json:"outputs,omitempty"`
			Error       string         `json:"error,omitempty"`
			ElapsedTime float64        `json:"elapsed_time"`
			TotalTokens int            `json:"total_tokens"`
			TotalSteps  int            `json:"total_steps"`
			CreatedAt   int64          `json:"created_at"`
			FinishedAt  int64          `json:"finished_at"`
		}{
			ID:          stream_response.Data.ID,
			WorkflowID:  stream_response.Data.WorkflowID,
			Status:      stream_response.Data.Status,
			Outputs:     stream_response.Data.Outputs,
			Error:       stream_response.Data.Error,
			ElapsedTime: stream_response.Data.ElapsedTime,
			TotalTokens: stream_response.Data.TotalTokens,
			TotalSteps:  stream_response.Data.TotalSteps,
			CreatedAt:   stream_response.Data.CreatedAt.Unix(),
			FinishedAt:  stream_response.Data.FinishedAt.Unix(),
		},
	}
}
