package response

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	dbengine "mlib.com/gofy/server/db_engine"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appresponserentities "mlib.com/gofy/server/entities/app/responser"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
)

// StreamResponse represents a stream response.
type StreamResponse struct {
	// SEvent appenumtypes.StreamEventType `json:"event"`
	taskID string
}

func (sr *StreamResponse) TaskID() string {
	return sr.taskID
}

// toDict converts a StreamResponse to a dictionary.
func (sr *StreamResponse) ToDict(responser appresponserentities.StreamResponser) map[string]any {
	bindata, _ := json.Marshal(responser)
	var dic map[string]any
	err := json.Unmarshal(bindata, &dic)
	if err != nil {
		mlog.Error("json.Unmarshal failed:", err)
		return nil
	}
	dic["event"] = responser.Event()
	dic["task_id"] = responser.TaskID()
	return dic
}

// ErrorStreamResponse represents an error stream response.
type ErrorStreamResponse struct {
	*StreamResponse
	Err error `json:"err"`
}

func NewErrorStreamResponse(task_id string, exp error) *ErrorStreamResponse {
	return &ErrorStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		Err: exp,
	}
}

func (rsp *ErrorStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_ERROR
}

// MessageStreamResponse represents a message stream response.
type MessageStreamResponse struct {
	*StreamResponse
	ID                   string   `json:"id"`
	Answer               string   `json:"answer"`
	FromVariableSelector []string `json:"from_variable_selector,omitempty"`
}

func NewMessageStreamResponse(task_id, id, answer string, variable_selector []string) *MessageStreamResponse {
	return &MessageStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		ID:                   id,
		Answer:               answer,
		FromVariableSelector: variable_selector,
	}
}

func (rsp *MessageStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_MESSAGE
}

// MessageAudioStreamResponse represents a message audio stream response.
type MessageAudioStreamResponse struct {
	*StreamResponse
	Audio string `json:"audio"`
}

func (rsp *MessageAudioStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_TTS_MESSAGE
}

// MessageAudioEndStreamResponse represents a message audio end stream response.
type MessageAudioEndStreamResponse struct {
	*StreamResponse
	Audio string `json:"audio"`
}

func (rsp *MessageAudioEndStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_TTS_MESSAGE_END
}

// MessageEndStreamResponse represents a message end stream response.
type MessageEndStreamResponse struct {
	*StreamResponse
	ID       string           `json:"id"`
	Metadata map[string]any   `json:"metadata"`
	Files    []map[string]any `json:"files,omitempty"`
}

func NewMessageEndStreamResponse(id, task_id string, metadata map[string]any) *MessageEndStreamResponse {
	return &MessageEndStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		ID:       id,
		Metadata: metadata,
	}
}

func (rsp *MessageEndStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_MESSAGE_END
}

// MessageFileStreamResponse represents a message file stream response.
type MessageFileStreamResponse struct {
	*StreamResponse
	ID        string `json:"id"`
	Type      string `json:"type"`
	BelongsTo string `json:"belongs_to"`
	URL       string `json:"url"`
}

func (rsp *MessageFileStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_MESSAGE_FILE
}

// MessageReplaceStreamResponse represents a message replace stream response.
type MessageReplaceStreamResponse struct {
	*StreamResponse
	Answer string `json:"answer"`
}

func NewMessageReplaceStreamResponse(task_id, answer string) *MessageReplaceStreamResponse {
	return &MessageReplaceStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		Answer: answer,
	}
}

func (rsp *MessageReplaceStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_MESSAGE_REPLACE
}

// AgentThoughtStreamResponse represents an agent thought stream response.
type AgentThoughtStreamResponse struct {
	*StreamResponse
	ID           string         `json:"id"`
	Position     int            `json:"position"`
	Thought      string         `json:"thought,omitempty"`
	Observation  string         `json:"observation,omitempty"`
	Tool         string         `json:"tool,omitempty"`
	ToolLabels   map[string]any `json:"tool_labels,omitempty"`
	ToolInput    string         `json:"tool_input,omitempty"`
	MessageFiles []any          `json:"message_files,omitempty"`
}

func NewAgentThoughtStreamResponse(id, task_id string) *AgentThoughtStreamResponse {
	return &AgentThoughtStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		ID: id,
	}
}

func (rsp *AgentThoughtStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_AGENT_THOUGHT
}

// AgentMessageStreamResponse represents an agent message stream response.
type AgentMessageStreamResponse struct {
	*StreamResponse
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

func NewAgentMessageStreamResponse(id, task_id, answer string) *AgentMessageStreamResponse {
	return &AgentMessageStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		ID:     id,
		Answer: answer,
	}
}

func (rsp *AgentMessageStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_AGENT_MESSAGE
}

type WorkflowStartStreamResponse struct {
	*StreamResponse
	/*
	   WorkflowStartStreamResponse entity
	*/
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID             string         `json:"id"`
		WorkflowID     string         `json:"workflow_id"`
		SequenceNumber int            `json:"sequence_number"`
		Inputs         map[string]any `json:"inputs" `
		CreatedAt      int64          `json:"created_at"`
	} `json:"data"`
}

func NewWorkflowStartStreamResponse(task_id string, workflow_run *models.WorkflowRun) *WorkflowStartStreamResponse {
	return &WorkflowStartStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_run.ID,
		Data: struct {
			ID             string         `json:"id"`
			WorkflowID     string         `json:"workflow_id"`
			SequenceNumber int            `json:"sequence_number"`
			Inputs         map[string]any `json:"inputs" `
			CreatedAt      int64          `json:"created_at"`
		}{
			ID:             workflow_run.ID,
			WorkflowID:     workflow_run.WorkflowID,
			SequenceNumber: workflow_run.SequenceNumber,
			Inputs:         workflow_run.InputsDict(),
			CreatedAt:      workflow_run.CreatedAt.Unix(),
		},
	}
}

func (rsp *WorkflowStartStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_WORKFLOW_STARTED
}

type WorkflowFinishStreamResponse struct {
	*StreamResponse
	/*
	   WorkflowFinishStreamResponse entity
	*/
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID              string           `json:"id"`
		WorkflowID      string           `json:"workflow_id"`
		SequenceNumber  int              `json:"sequence_number"`
		Status          string           `json:"status"`
		Outputs         map[string]any   `json:"outputs,omitempty"`
		Error           string           `json:"error,omitempty"`
		ElapsedTime     float64          `json:"elapsed_time"`
		TotalTokens     int              `json:"total_tokens"`
		TotalSteps      int              `json:"total_steps"`
		CreatedBy       map[string]any   `json:"created_by,omitempty"`
		CreatedAt       time.Time        `json:"created_at"`
		FinishedAt      time.Time        `json:"finished_at"`
		ExceptionsCount int              `json:"exceptions_count,omitempty"`
		Files           []map[string]any `json:"files,omitempty"`
	} `json:"data"`
}

func NewWorkflowFinishStreamResponse(task_id string, workflow_run *models.WorkflowRun) *WorkflowFinishStreamResponse {
	var created_by map[string]any
	if workflow_run.CreatedByRole == models.CreatedByRole_ACCOUNT {
		account := new(models.Account)
		err := dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", workflow_run.CreatedBy).First(account).Error
		if err != nil {
			mlog.Errorf("get account(%s) failed:%v", workflow_run.CreatedBy, err)
		} else {
			created_by = map[string]any{
				"id":    account.ID,
				"name":  account.Name,
				"email": account.Email,
			}
		}
	} else if workflow_run.CreatedByRole == models.CreatedByRole_END_USER {
		end_user := new(models.EndUser)
		err := dbengine.Instance().DB.Model(&models.EndUser{}).Where("id = ?", workflow_run.CreatedBy).First(end_user).Error
		if err != nil {
			mlog.Errorf("get EndUser(%s) failed:%v", workflow_run.CreatedBy, err)
		} else {
			created_by = map[string]any{
				"id":   end_user.ID,
				"user": end_user.SessionID,
			}
		}
	} else {
		panic(exceptions.NewNotImplementedError(fmt.Sprintf("unknown created_by_role: %v", workflow_run.CreatedByRole)))
	}
	return &WorkflowFinishStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_run.ID,
		Data: struct {
			ID              string           `json:"id"`
			WorkflowID      string           `json:"workflow_id"`
			SequenceNumber  int              `json:"sequence_number"`
			Status          string           `json:"status"`
			Outputs         map[string]any   `json:"outputs,omitempty"`
			Error           string           `json:"error,omitempty"`
			ElapsedTime     float64          `json:"elapsed_time"`
			TotalTokens     int              `json:"total_tokens"`
			TotalSteps      int              `json:"total_steps"`
			CreatedBy       map[string]any   `json:"created_by,omitempty"`
			CreatedAt       time.Time        `json:"created_at"`
			FinishedAt      time.Time        `json:"finished_at"`
			ExceptionsCount int              `json:"exceptions_count,omitempty"`
			Files           []map[string]any `json:"files,omitempty"`
		}{
			ID:              workflow_run.ID,
			WorkflowID:      workflow_run.WorkflowID,
			SequenceNumber:  workflow_run.SequenceNumber,
			Status:          workflow_run.Status,
			Outputs:         workflow_run.OutputsDict(),
			Error:           workflow_run.Error,
			ElapsedTime:     workflow_run.ElapsedTime,
			TotalTokens:     workflow_run.TotalTokens,
			TotalSteps:      workflow_run.TotalSteps,
			CreatedBy:       created_by,
			CreatedAt:       *workflow_run.CreatedAt,
			FinishedAt:      *workflow_run.FinishedAt,
			Files:           file.FetchFilesFromNodeOutputs(workflow_run.OutputsDict()),
			ExceptionsCount: workflow_run.ExceptionsCount,
		},
	}
}

func (rsp *WorkflowFinishStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_WORKFLOW_FINISHED
}

type NodeStartStreamResponse struct {
	*StreamResponse
	/*
	   NodeStartStreamResponse entity
	*/
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID                        string         `json:"id"`
		NodeID                    string         `json:"node_id"`
		NodeType                  string         `json:"node_type"`
		Title                     string         `json:"title"`
		Index                     int            `json:"index"`
		PredecessorNodeID         string         `json:"predecessor_node_id,omitempty"`
		Inputs                    map[string]any `json:"inputs,omitempty"`
		CreatedAt                 int64          `json:"created_at"`
		Extras                    map[string]any `json:"extras"`
		ParallelID                string         `json:"parallel_id,omitempty"`
		ParallelStartNodeID       string         `json:"parallel_start_node_id,omitempty"`
		ParentParallelID          string         `json:"parent_parallel_id,omitempty"`
		ParentParallelStartNodeID string         `json:"parent_parallel_start_node_id,omitempty"`
		IterationID               string         `json:"iteration_id,omitempty"`
		ParallelRunID             string         `json:"parallel_run_id,omitempty"`
	} `json:"data"`
}

func NewNodeStartStreamResponse(event *appqueueentities.QueueNodeStartedEvent, task_id string, workflow_node_execution *models.WorkflowNodeExecution) *NodeStartStreamResponse {
	response := &NodeStartStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_node_execution.WorkflowRunID,
		Data: struct {
			ID                        string         `json:"id"`
			NodeID                    string         `json:"node_id"`
			NodeType                  string         `json:"node_type"`
			Title                     string         `json:"title"`
			Index                     int            `json:"index"`
			PredecessorNodeID         string         `json:"predecessor_node_id,omitempty"`
			Inputs                    map[string]any `json:"inputs,omitempty"`
			CreatedAt                 int64          `json:"created_at"`
			Extras                    map[string]any `json:"extras"`
			ParallelID                string         `json:"parallel_id,omitempty"`
			ParallelStartNodeID       string         `json:"parallel_start_node_id,omitempty"`
			ParentParallelID          string         `json:"parent_parallel_id,omitempty"`
			ParentParallelStartNodeID string         `json:"parent_parallel_start_node_id,omitempty"`
			IterationID               string         `json:"iteration_id,omitempty"`
			ParallelRunID             string         `json:"parallel_run_id,omitempty"`
		}{
			ID:                        workflow_node_execution.ID,
			NodeID:                    workflow_node_execution.NodeID,
			NodeType:                  string(workflow_node_execution.NodeType),
			Title:                     workflow_node_execution.Title,
			Index:                     workflow_node_execution.Index,
			PredecessorNodeID:         workflow_node_execution.PredecessorNodeID,
			Inputs:                    workflow_node_execution.InputsDict(),
			CreatedAt:                 workflow_node_execution.CreatedAt.Unix(),
			ParallelID:                event.ParallelID,
			ParallelStartNodeID:       event.ParallelStartNodeID,
			ParentParallelID:          event.ParentParallelID,
			ParentParallelStartNodeID: event.ParentParallelStartNodeID,
			IterationID:               event.InIterationID,
			ParallelRunID:             event.ParallelModeRunID,
		},
	}

	// extras logic
	// if event.NodeType == nodesenumtypes.Node_TOOL{
	// 	event.
	// 	node_data = cast(ToolNodeData, event.node_data)
	// 	response.data.extras["icon"] = ToolManager.get_tool_icon(
	// 		tenant_id=mgr._application_generate_entity.app_config.tenant_id,
	// 		provider_type=node_data.provider_type,
	// 		provider_id=node_data.provider_id,
	// 	)
	// }
	return response
}

func (rsp *NodeStartStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_NODE_STARTED
}

func (rsp *NodeStartStreamResponse) ToIgnoreDetailDict() map[string]any {
	return map[string]any{
		"event":           rsp.Event(),
		"task_id":         rsp.TaskID,
		"workflow_run_id": rsp.WorkflowRunID,
		"data": map[string]any{
			"id":                            rsp.Data.ID,
			"node_id":                       rsp.Data.NodeID,
			"node_type":                     rsp.Data.NodeType,
			"title":                         rsp.Data.Title,
			"index":                         rsp.Data.Index,
			"predecessor_node_id":           rsp.Data.PredecessorNodeID,
			"created_at":                    rsp.Data.CreatedAt,
			"extras":                        map[string]any{},
			"parallel_id":                   rsp.Data.ParallelID,
			"parallel_start_node_id":        rsp.Data.ParallelStartNodeID,
			"parent_parallel_id":            rsp.Data.ParentParallelID,
			"parent_parallel_start_node_id": rsp.Data.ParentParallelStartNodeID,
			"iteration_id":                  rsp.Data.IterationID,
		},
	}
}

type NodeFinishStreamResponse struct {
	*StreamResponse
	/*
	   NodeFinishStreamResponse entity
	*/
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID                        string           `json:"id"`
		NodeID                    string           `json:"node_id"`
		NodeType                  string           `json:"node_type"`
		Title                     string           `json:"title"`
		Index                     int              `json:"index"`
		PredecessorNodeID         string           `json:"predecessor_node_id,omitempty"`
		Inputs                    map[string]any   `json:"inputs,omitempty"`
		ProcessData               map[string]any   `json:"process_data,omitempty"`
		Outputs                   map[string]any   `json:"outputs,omitempty"`
		Status                    string           `json:"status"`
		Error                     string           `json:"error,omitempty"`
		ElapsedTime               float64          `json:"elapsed_time"`
		ExecutionMetadata         map[string]any   `json:"execution_metadata,omitempty"`
		CreatedAt                 int64            `json:"created_at"`
		FinishedAt                int64            `json:"finished_at"`
		Files                     []map[string]any `json:"files,omitempty"`
		ParallelID                string           `json:"parallel_id,omitempty"`
		ParallelStartNodeID       string           `json:"parallel_start_node_id,omitempty"`
		ParentParallelID          string           `json:"parent_parallel_id,omitempty"`
		ParentParallelStartNodeID string           `json:"parent_parallel_start_node_id,omitempty"`
		IterationID               string           `json:"iteration_id,omitempty"`
	} `json:"data"`
}

func NewNodeFinishStreamResponse[T *appqueueentities.QueueNodeSucceededEvent | *appqueueentities.QueueNodeFailedEvent | *appqueueentities.QueueNodeInIterationFailedEvent | *appqueueentities.QueueNodeExceptionEvent](
	event T,
	task_id string,
	workflow_node_execution *models.WorkflowNodeExecution,
) *NodeFinishStreamResponse {
	rsp := &NodeFinishStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_node_execution.WorkflowRunID,
		Data: struct {
			ID                        string           `json:"id"`
			NodeID                    string           `json:"node_id"`
			NodeType                  string           `json:"node_type"`
			Title                     string           `json:"title"`
			Index                     int              `json:"index"`
			PredecessorNodeID         string           `json:"predecessor_node_id,omitempty"`
			Inputs                    map[string]any   `json:"inputs,omitempty"`
			ProcessData               map[string]any   `json:"process_data,omitempty"`
			Outputs                   map[string]any   `json:"outputs,omitempty"`
			Status                    string           `json:"status"`
			Error                     string           `json:"error,omitempty"`
			ElapsedTime               float64          `json:"elapsed_time"`
			ExecutionMetadata         map[string]any   `json:"execution_metadata,omitempty"`
			CreatedAt                 int64            `json:"created_at"`
			FinishedAt                int64            `json:"finished_at"`
			Files                     []map[string]any `json:"files,omitempty"`
			ParallelID                string           `json:"parallel_id,omitempty"`
			ParallelStartNodeID       string           `json:"parallel_start_node_id,omitempty"`
			ParentParallelID          string           `json:"parent_parallel_id,omitempty"`
			ParentParallelStartNodeID string           `json:"parent_parallel_start_node_id,omitempty"`
			IterationID               string           `json:"iteration_id,omitempty"`
		}{
			ID:                workflow_node_execution.ID,
			NodeID:            workflow_node_execution.NodeID,
			NodeType:          string(workflow_node_execution.NodeType),
			Title:             workflow_node_execution.Title,
			Index:             workflow_node_execution.Index,
			PredecessorNodeID: workflow_node_execution.PredecessorNodeID,
			Inputs:            workflow_node_execution.InputsDict(),
			ProcessData:       workflow_node_execution.ProcessDataDict(),
			Outputs:           workflow_node_execution.OutputsDict(),
			Status:            workflow_node_execution.Status,
			Error:             workflow_node_execution.Error,
			ElapsedTime:       workflow_node_execution.ElapsedTime,
			ExecutionMetadata: workflow_node_execution.ExecutionMetadataDict(),
			CreatedAt:         workflow_node_execution.CreatedAt.Unix(),
			FinishedAt:        workflow_node_execution.FinishedAt.Unix(),
			Files:             file.FetchFilesFromNodeOutputs(workflow_node_execution.OutputsDict()),
			// ParallelID               :"",
			// ParallelStartNodeID      :,
			// ParentParallelID         :,
			// ParentParallelStartNodeID:,
			// IterationID              :,
		},
	}

	switch ev := any(event).(type) {
	case *appqueueentities.QueueNodeSucceededEvent:
		rsp.Data.ParallelID = ev.ParallelID
		rsp.Data.ParallelStartNodeID = ev.ParallelStartNodeID
		rsp.Data.ParentParallelID = ev.ParentParallelID
		rsp.Data.ParentParallelStartNodeID = ev.ParentParallelStartNodeID
		rsp.Data.IterationID = ev.InIterationID
	case *appqueueentities.QueueNodeFailedEvent:
		rsp.Data.ParallelID = ev.ParallelID
		rsp.Data.ParallelStartNodeID = ev.ParallelStartNodeID
		rsp.Data.ParentParallelID = ev.ParentParallelID
		rsp.Data.ParentParallelStartNodeID = ev.ParentParallelStartNodeID
		rsp.Data.IterationID = ev.InIterationID
	case *appqueueentities.QueueNodeInIterationFailedEvent:
		rsp.Data.ParallelID = ev.ParallelID
		rsp.Data.ParallelStartNodeID = ev.ParallelStartNodeID
		rsp.Data.ParentParallelID = ev.ParentParallelID
		rsp.Data.ParentParallelStartNodeID = ev.ParentParallelStartNodeID
		rsp.Data.IterationID = ev.InIterationID
	case *appqueueentities.QueueNodeExceptionEvent:
		rsp.Data.ParallelID = ev.ParallelID
		rsp.Data.ParallelStartNodeID = ev.ParallelStartNodeID
		rsp.Data.ParentParallelID = ev.ParentParallelID
		rsp.Data.ParentParallelStartNodeID = ev.ParentParallelStartNodeID
		rsp.Data.IterationID = ev.InIterationID
	}
	return rsp
}

func (rsp *NodeFinishStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_NODE_FINISHED
}

func (rsp *NodeFinishStreamResponse) ToIgnoreDetailDict() map[string]any {
	return map[string]any{
		"event":           rsp.Event(),
		"task_id":         rsp.TaskID,
		"workflow_run_id": rsp.WorkflowRunID,
		"data": map[string]any{
			"id":                            rsp.Data.ID,
			"node_id":                       rsp.Data.NodeID,
			"node_type":                     rsp.Data.NodeType,
			"title":                         rsp.Data.Title,
			"index":                         rsp.Data.Index,
			"predecessor_node_id":           rsp.Data.PredecessorNodeID,
			"status":                        rsp.Data.Status,
			"elapsed_time":                  rsp.Data.ElapsedTime,
			"created_at":                    rsp.Data.CreatedAt,
			"finished_at":                   rsp.Data.FinishedAt,
			"parallel_id":                   rsp.Data.ParallelID,
			"parallel_start_node_id":        rsp.Data.ParallelStartNodeID,
			"parent_parallel_id":            rsp.Data.ParentParallelID,
			"parent_parallel_start_node_id": rsp.Data.ParentParallelStartNodeID,
			"iteration_id":                  rsp.Data.IterationID,
		},
	}
}

type NodeRetryStreamResponse struct {
	*StreamResponse
	/*
	   NodeRetryStreamResponse entity
	*/
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID                        string           `json:"id"`
		NodeID                    string           `json:"node_id"`
		NodeType                  string           `json:"node_type"`
		Title                     string           `json:"title"`
		Index                     int              `json:"index"`
		PredecessorNodeID         string           `json:"predecessor_node_id,omitempty"`
		Inputs                    map[string]any   `json:"inputs,omitempty"`
		ProcessData               map[string]any   `json:"process_data,omitempty"`
		Outputs                   map[string]any   `json:"outputs,omitempty"`
		Status                    string           `json:"status"`
		Error                     string           `json:"error,omitempty"`
		ElapsedTime               float64          `json:"elapsed_time"`
		ExecutionMetadata         map[string]any   `json:"execution_metadata,omitempty"`
		CreatedAt                 int64            `json:"created_at"`
		FinishedAt                int64            `json:"finished_at"`
		Files                     []map[string]any `json:"files,omitempty"`
		ParallelID                string           `json:"parallel_id,omitempty"`
		ParallelStartNodeID       string           `json:"parallel_start_node_id,omitempty"`
		ParentParallelID          string           `json:"parent_parallel_id,omitempty"`
		ParentParallelStartNodeID string           `json:"parent_parallel_start_node_id,omitempty"`
		IterationID               string           `json:"iteration_id,omitempty"`
		RetryIndex                int              `json:"retry_index"`
	} `json:"data"`
}

func NewNodeRetryStreamResponse(event *appqueueentities.QueueNodeRetryEvent, task_id string, workflow_node_execution *models.WorkflowNodeExecution) *NodeRetryStreamResponse {
	return &NodeRetryStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_node_execution.WorkflowRunID,
		Data: struct {
			ID                        string           `json:"id"`
			NodeID                    string           `json:"node_id"`
			NodeType                  string           `json:"node_type"`
			Title                     string           `json:"title"`
			Index                     int              `json:"index"`
			PredecessorNodeID         string           `json:"predecessor_node_id,omitempty"`
			Inputs                    map[string]any   `json:"inputs,omitempty"`
			ProcessData               map[string]any   `json:"process_data,omitempty"`
			Outputs                   map[string]any   `json:"outputs,omitempty"`
			Status                    string           `json:"status"`
			Error                     string           `json:"error,omitempty"`
			ElapsedTime               float64          `json:"elapsed_time"`
			ExecutionMetadata         map[string]any   `json:"execution_metadata,omitempty"`
			CreatedAt                 int64            `json:"created_at"`
			FinishedAt                int64            `json:"finished_at"`
			Files                     []map[string]any `json:"files,omitempty"`
			ParallelID                string           `json:"parallel_id,omitempty"`
			ParallelStartNodeID       string           `json:"parallel_start_node_id,omitempty"`
			ParentParallelID          string           `json:"parent_parallel_id,omitempty"`
			ParentParallelStartNodeID string           `json:"parent_parallel_start_node_id,omitempty"`
			IterationID               string           `json:"iteration_id,omitempty"`
			RetryIndex                int              `json:"retry_index"`
		}{
			ID:                        workflow_node_execution.ID,
			NodeID:                    workflow_node_execution.NodeID,
			NodeType:                  string(workflow_node_execution.NodeType),
			Title:                     workflow_node_execution.Title,
			Index:                     workflow_node_execution.Index,
			PredecessorNodeID:         workflow_node_execution.PredecessorNodeID,
			Inputs:                    workflow_node_execution.InputsDict(),
			ProcessData:               workflow_node_execution.ProcessDataDict(),
			Outputs:                   workflow_node_execution.OutputsDict(),
			Status:                    workflow_node_execution.Status,
			Error:                     workflow_node_execution.Error,
			ElapsedTime:               workflow_node_execution.ElapsedTime,
			ExecutionMetadata:         workflow_node_execution.ExecutionMetadataDict(),
			CreatedAt:                 workflow_node_execution.CreatedAt.Unix(),
			FinishedAt:                workflow_node_execution.FinishedAt.Unix(),
			Files:                     file.FetchFilesFromNodeOutputs(workflow_node_execution.OutputsDict()),
			ParallelID:                event.ParallelID,
			ParallelStartNodeID:       event.ParallelStartNodeID,
			ParentParallelID:          event.ParentParallelID,
			ParentParallelStartNodeID: event.ParentParallelStartNodeID,
			IterationID:               event.InIterationID,
			RetryIndex:                event.RetryIndex,
		},
	}
}

func (rsp *NodeRetryStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_NODE_RETRY
}

func (rsp *NodeRetryStreamResponse) ToIgnoreDetailDict() map[string]any {
	return map[string]any{
		"event":           rsp.Event(),
		"task_id":         rsp.TaskID,
		"workflow_run_id": rsp.WorkflowRunID,
		"data": map[string]any{
			"id":                            rsp.Data.ID,
			"node_id":                       rsp.Data.NodeID,
			"node_type":                     rsp.Data.NodeType,
			"title":                         rsp.Data.Title,
			"index":                         rsp.Data.Index,
			"predecessor_node_id":           rsp.Data.PredecessorNodeID,
			"status":                        rsp.Data.Status,
			"elapsed_time":                  rsp.Data.ElapsedTime,
			"created_at":                    rsp.Data.CreatedAt,
			"finished_at":                   rsp.Data.FinishedAt,
			"parallel_id":                   rsp.Data.ParallelID,
			"parallel_start_node_id":        rsp.Data.ParallelStartNodeID,
			"parent_parallel_id":            rsp.Data.ParentParallelID,
			"parent_parallel_start_node_id": rsp.Data.ParentParallelStartNodeID,
			"iteration_id":                  rsp.Data.IterationID,
			"retry_index":                   rsp.Data.RetryIndex,
		},
	}
}

// ParallelBranchStartStreamResponse entity
type ParallelBranchStartStreamResponse struct {
	*StreamResponse
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ParallelID                string `json:"parallel_id"`
		ParallelBranchID          string `json:"parallel_branch_id"`
		ParentParallelID          string `json:"parent_parallel_id,omitempty"`
		ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
		IterationID               string `json:"iteration_id,omitempty"`
		CreatedAt                 int64  `json:"created_at"`
	} `json:"data"`
}

func NewParallelBranchStartStreamResponse(event *appqueueentities.QueueParallelBranchRunStartedEvent, task_id string, workflow_run *models.WorkflowRun) *ParallelBranchStartStreamResponse {
	return &ParallelBranchStartStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_run.ID,
		Data: struct {
			ParallelID                string `json:"parallel_id"`
			ParallelBranchID          string `json:"parallel_branch_id"`
			ParentParallelID          string `json:"parent_parallel_id,omitempty"`
			ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
			IterationID               string `json:"iteration_id,omitempty"`
			CreatedAt                 int64  `json:"created_at"`
		}{
			ParallelID:                event.ParallelID,
			ParallelBranchID:          event.ParallelStartNodeID,
			ParentParallelID:          event.ParentParallelID,
			ParentParallelStartNodeID: event.ParentParallelStartNodeID,
			IterationID:               event.InIterationID,
			CreatedAt:                 time.Now().Unix(),
		},
	}
}

func (rsp *ParallelBranchStartStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_PARALLEL_BRANCH_STARTED
}

// ParallelBranchFinishedStreamResponse entity
type ParallelBranchFinishedStreamResponse struct {
	*StreamResponse
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ParallelID                string `json:"parallel_id"`
		ParallelBranchID          string `json:"parallel_branch_id"`
		ParentParallelID          string `json:"parent_parallel_id,omitempty"`
		ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
		IterationID               string `json:"iteration_id,omitempty"`
		Status                    string `json:"status"`
		Error                     string `json:"error,omitempty"`
		CreatedAt                 int64  `json:"created_at"`
	} `json:"data"`
}

func NewParallelBranchFinishedStreamResponse[T *appqueueentities.QueueParallelBranchRunSucceededEvent | *appqueueentities.QueueParallelBranchRunFailedEvent](
	event T, task_id string, workflow_run *models.WorkflowRun,
) *ParallelBranchFinishedStreamResponse {
	rsp := &ParallelBranchFinishedStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_run.ID,
		Data: struct {
			ParallelID                string `json:"parallel_id"`
			ParallelBranchID          string `json:"parallel_branch_id"`
			ParentParallelID          string `json:"parent_parallel_id,omitempty"`
			ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
			IterationID               string `json:"iteration_id,omitempty"`
			Status                    string `json:"status"`
			Error                     string `json:"error,omitempty"`
			CreatedAt                 int64  `json:"created_at"`
		}{
			CreatedAt: time.Now().Unix(),
		},
	}
	switch ev := any(event).(type) {
	case *appqueueentities.QueueParallelBranchRunSucceededEvent:
		rsp.Data.ParallelID = ev.ParallelID
		rsp.Data.ParallelBranchID = ev.ParallelStartNodeID
		rsp.Data.ParentParallelID = ev.ParentParallelID
		rsp.Data.ParentParallelStartNodeID = ev.ParentParallelStartNodeID
		rsp.Data.IterationID = ev.InIterationID
		rsp.Data.Status = "succeeded"
	case *appqueueentities.QueueParallelBranchRunFailedEvent:
		rsp.Data.ParallelID = ev.ParallelID
		rsp.Data.ParallelBranchID = ev.ParallelStartNodeID
		rsp.Data.ParentParallelID = ev.ParentParallelID
		rsp.Data.ParentParallelStartNodeID = ev.ParentParallelStartNodeID
		rsp.Data.IterationID = ev.InIterationID
		rsp.Data.Status = "failed"
		rsp.Data.Error = ev.Error
	}
	return rsp
}

func (rsp *ParallelBranchFinishedStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_PARALLEL_BRANCH_FINISHED
}

// IterationNodeStartStreamResponse entity
type IterationNodeStartStreamResponse struct {
	*StreamResponse
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID                  string         `json:"id"`
		NodeID              string         `json:"node_id"`
		NodeType            string         `json:"node_type"`
		Title               string         `json:"title"`
		CreatedAt           int64          `json:"created_at"`
		Extras              map[string]any `json:"extras"`
		Metadata            map[string]any `json:"metadata"`
		Inputs              map[string]any `json:"inputs"`
		ParallelID          string         `json:"parallel_id,omitempty"`
		ParallelStartNodeID string         `json:"parallel_start_node_id,omitempty"`
	} `json:"data"`
}

func NewIterationNodeStartStreamResponse(
	event *appqueueentities.QueueIterationStartEvent, task_id string, workflow_run *models.WorkflowRun,
) *IterationNodeStartStreamResponse {
	return &IterationNodeStartStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_run.ID,
		Data: struct {
			ID                  string         `json:"id"`
			NodeID              string         `json:"node_id"`
			NodeType            string         `json:"node_type"`
			Title               string         `json:"title"`
			CreatedAt           int64          `json:"created_at"`
			Extras              map[string]any `json:"extras"`
			Metadata            map[string]any `json:"metadata"`
			Inputs              map[string]any `json:"inputs"`
			ParallelID          string         `json:"parallel_id,omitempty"`
			ParallelStartNodeID string         `json:"parallel_start_node_id,omitempty"`
		}{
			ID:                  event.NodeID,
			NodeID:              event.NodeID,
			NodeType:            string(event.NodeType),
			Title:               event.NodeData.Title,
			CreatedAt:           time.Now().Unix(),
			Extras:              map[string]any{},
			Metadata:            event.Metadata,
			Inputs:              event.Inputs,
			ParallelID:          event.ParallelID,
			ParallelStartNodeID: event.ParallelStartNodeID,
		},
	}
}

func (rsp *IterationNodeStartStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_ITERATION_STARTED
}

// IterationNodeNextStreamResponse entity
type IterationNodeNextStreamResponse struct {
	*StreamResponse
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID                  string         `json:"id"`
		NodeID              string         `json:"node_id"`
		NodeType            string         `json:"node_type"`
		Title               string         `json:"title"`
		Index               int            `json:"index"`
		CreatedAt           int64          `json:"created_at"`
		PreIterationOutput  any            `json:"pre_iteration_output,omitempty"`
		Extras              map[string]any `json:"extras"`
		ParallelID          string         `json:"parallel_id,omitempty"`
		ParallelStartNodeID string         `json:"parallel_start_node_id,omitempty"`
		ParallelModeRunID   string         `json:"parallel_mode_run_id,omitempty"`
		Duration            float64        `json:"duration,omitempty"`
	} `json:"data"`
}

func NewIterationNodeNextStreamResponse(
	event *appqueueentities.QueueIterationNextEvent, task_id string, workflow_run *models.WorkflowRun,
) *IterationNodeNextStreamResponse {
	return &IterationNodeNextStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_run.ID,
		Data: struct {
			ID                  string         `json:"id"`
			NodeID              string         `json:"node_id"`
			NodeType            string         `json:"node_type"`
			Title               string         `json:"title"`
			Index               int            `json:"index"`
			CreatedAt           int64          `json:"created_at"`
			PreIterationOutput  any            `json:"pre_iteration_output,omitempty"`
			Extras              map[string]any `json:"extras"`
			ParallelID          string         `json:"parallel_id,omitempty"`
			ParallelStartNodeID string         `json:"parallel_start_node_id,omitempty"`
			ParallelModeRunID   string         `json:"parallel_mode_run_id,omitempty"`
			Duration            float64        `json:"duration,omitempty"`
		}{
			ID:                  event.NodeID,
			NodeID:              event.NodeID,
			NodeType:            string(event.NodeType),
			Title:               event.NodeData.Title,
			CreatedAt:           time.Now().Unix(),
			Extras:              map[string]any{},
			Index:               event.Index,
			PreIterationOutput:  event.Output,
			ParallelID:          event.ParallelID,
			ParallelStartNodeID: event.ParallelStartNodeID,
			ParallelModeRunID:   event.ParallelModeRunID,
			Duration:            event.Duration,
		},
	}
}

func (rsp *IterationNodeNextStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_ITERATION_NEXT
}

// IterationNodeCompletedStreamResponse entity
type IterationNodeCompletedStreamResponse struct {
	*StreamResponse
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID                  string         `json:"id"`
		NodeID              string         `json:"node_id"`
		NodeType            string         `json:"node_type"`
		Title               string         `json:"title"`
		Outputs             map[string]any `json:"outputs,omitempty"`
		CreatedAt           int64          `json:"created_at"`
		Extras              map[string]any `json:"extras,omitempty"`
		Inputs              map[string]any `json:"inputs,omitempty"`
		Status              string         `json:"status"`
		Error               string         `json:"error,omitempty"`
		ElapsedTime         float64        `json:"elapsed_time"`
		TotalTokens         int            `json:"total_tokens"`
		ExecutionMetadata   map[string]any `json:"execution_metadata,omitempty"`
		FinishedAt          int            `json:"finished_at"`
		Steps               int            `json:"steps"`
		ParallelID          string         `json:"parallel_id,omitempty"`
		ParallelStartNodeID string         `json:"parallel_start_node_id,omitempty"`
	} `json:"data"`
}

func NewIterationNodeCompletedStreamResponse(
	event *appqueueentities.QueueIterationCompletedEvent, task_id string, workflow_run *models.WorkflowRun,
) *IterationNodeCompletedStreamResponse {
	rsp := &IterationNodeCompletedStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		WorkflowRunID: workflow_run.ID,
		Data: struct {
			ID                  string         `json:"id"`
			NodeID              string         `json:"node_id"`
			NodeType            string         `json:"node_type"`
			Title               string         `json:"title"`
			Outputs             map[string]any `json:"outputs,omitempty"`
			CreatedAt           int64          `json:"created_at"`
			Extras              map[string]any `json:"extras,omitempty"`
			Inputs              map[string]any `json:"inputs,omitempty"`
			Status              string         `json:"status"`
			Error               string         `json:"error,omitempty"`
			ElapsedTime         float64        `json:"elapsed_time"`
			TotalTokens         int            `json:"total_tokens"`
			ExecutionMetadata   map[string]any `json:"execution_metadata,omitempty"`
			FinishedAt          int            `json:"finished_at"`
			Steps               int            `json:"steps"`
			ParallelID          string         `json:"parallel_id,omitempty"`
			ParallelStartNodeID string         `json:"parallel_start_node_id,omitempty"`
		}{
			ID:                  event.NodeID,
			NodeID:              event.NodeID,
			NodeType:            string(event.NodeType),
			Title:               event.NodeData.Title,
			CreatedAt:           time.Now().Unix(),
			Extras:              map[string]any{},
			Outputs:             event.Outputs,
			ParallelID:          event.ParallelID,
			ParallelStartNodeID: event.ParallelStartNodeID,
			Inputs:              event.Inputs,
			Status:              string(models.WorkflowNodeExecutionStatus_SUCCEEDED),
			ElapsedTime:         time.Since(event.StartAt).Seconds(),
			// TotalTokens        :,
			ExecutionMetadata: event.Metadata,
			FinishedAt:        int(time.Now().Unix()),
			Steps:             event.Steps,
		},
	}
	if event.Error != "" {
		rsp.Data.Status = string(models.WorkflowNodeExecutionStatus_FAILED)
		rsp.Data.Error = event.Error
	}
	return rsp
}
func (rsp *IterationNodeCompletedStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_ITERATION_COMPLETED
}

// TextChunkStreamResponse entity
type TextChunkStreamResponse struct {
	*StreamResponse
	/*
	   TextChunkStreamResponse entity
	*/
	Data struct {
		Text                 string   `json:"text"`
		FromVariableSelector []string `json:"from_variable_selector,omitempty"`
	} `json:"data"`
}

func NewTextChunkStreamResponse(
	task_id string, text string, from_variable_selector []string,
) *TextChunkStreamResponse {
	return &TextChunkStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
		Data: struct {
			Text                 string   `json:"text"`
			FromVariableSelector []string `json:"from_variable_selector,omitempty"`
		}{
			Text:                 text,
			FromVariableSelector: from_variable_selector,
		},
	}
}

func (rsp *TextChunkStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_TEXT_CHUNK
}

// TextReplaceStreamResponse entity
type TextReplaceStreamResponse struct {
	*StreamResponse
	/*
	   TextReplaceStreamResponse entity
	*/
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

func (rsp *TextReplaceStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_TEXT_REPLACE
}

// PingStreamResponse entity
type PingStreamResponse struct {
	*StreamResponse
	/*
	   PingStreamResponse entity
	*/
}

func NewPingStreamResponse(task_id string) *PingStreamResponse {
	return &PingStreamResponse{
		StreamResponse: &StreamResponse{
			taskID: task_id,
		},
	}
}

func (rsp *PingStreamResponse) Event() appenumtypes.StreamEventType {
	return appenumtypes.StreamEvent_PING
}

type AppStreamResponse struct {
	StreamResponser appresponserentities.StreamResponser
}

// ChatbotAppStreamResponse entity
type ChatbotAppStreamResponse struct {
	*AppStreamResponse
	/*
	   ChatbotAppStreamResponse entity
	*/
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
}

// CompletionAppStreamResponse entity
type CompletionAppStreamResponse struct {
	*AppStreamResponse
	/*
	   CompletionAppStreamResponse entity
	*/
	MessageID string `json:"message_id"`
	CreatedAt int64  `json:"created_at"`
}

// WorkflowAppStreamResponse entity
type WorkflowAppStreamResponse struct {
	*AppStreamResponse
	/*
	   WorkflowAppStreamResponse entity
	*/
	WorkflowRunID string `json:"workflow_run_id,omitempty"`
}

func NewWorkflowAppStreamResponse(workflow_run_id string, stream_response appresponserentities.StreamResponser) *WorkflowAppStreamResponse {
	return &WorkflowAppStreamResponse{
		AppStreamResponse: &AppStreamResponse{
			StreamResponser: stream_response,
		},
		WorkflowRunID: workflow_run_id,
	}
}
