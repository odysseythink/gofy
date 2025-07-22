package response

import (
	"encoding/json"

	"github.com/spf13/cast"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type AnnotationResponse struct {
	ID        string                 `json:"id"`
	Question  string                 `json:"question"`
	Content   string                 `json:"content"`
	Account   *SimpleAccountResponse `json:"account"`
	CreatedAt int64                  `json:"created_at"`
}

func NewAnnotationResponse(args any) *AnnotationResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(AnnotationResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AnnotationResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(AnnotationResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AnnotationResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if ma, ok := args.(*models.MessageAnnotation); ok && ma != nil {
		rsp := &AnnotationResponse{
			ID:       ma.ID,
			Question: ma.Question,
			Content:  ma.Content,
		}
		tmpacc := ma.Account()
		if tmpacc != nil {
			rsp.Account = &SimpleAccountResponse{
				ID:    tmpacc.ID,
				Name:  tmpacc.Name,
				Email: tmpacc.Email,
			}
		}
		if ma.CreatedAt != nil {
			rsp.CreatedAt = ma.CreatedAt.Unix()
		}
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(AnnotationResponse)
}

type AnnotationHitHistoryResponse struct {
	AnnotationID            string                 `json:"annotation_id"`
	AnnotationCreateAccount *SimpleAccountResponse `json:"annotation_create_account"`
	CreatedAt               int64                  `json:"created_at"`
}

type FeedbackResponse struct {
	Rating        string                 `json:"rating"`
	Content       string                 `json:"content"`
	FromSource    string                 `json:"from_source"`
	FromEndUserID string                 `json:"from_end_user_id"`
	FromAccount   *SimpleAccountResponse `json:"from_account"`
}
type MessageFileResponse struct {
	ID             string `json:"id"`
	Filename       string `json:"filename"`
	Type           string `json:"type"`
	URL            string `json:"url"`
	MimeType       string `json:"mime_type"`
	Size           int    `json:"size"`
	TransferMethod string `json:"transfer_method"`
	BelongsTo      string `json:"belongs_to"` //default="user"),
}
type AgentThoughtResponse struct {
	ID          string   `json:"id"`
	ChainID     string   `json:"chain_id"`
	MessageID   string   `json:"message_id"`
	Position    int      `json:"position"`
	Thought     string   `json:"thought"`
	Tool        string   `json:"tool"`
	ToolLabels  any      `json:"tool_labels"`
	ToolInput   string   `json:"tool_input"`
	CreatedAt   int64    `json:"created_at"`
	Observation string   `json:"observation"`
	Files       []string `json:"files"`
}

type MessageDetailResponse struct {
	ID                      string                        `json:"id"`
	ConversationID          string                        `json:"conversation_id"`
	Inputs                  map[string]any                `json:"inputs"`
	Query                   string                        `json:"query"`
	Message                 string                        `json:"message"`
	MessageTokens           int                           `json:"message_tokens"`
	Answer                  string                        `json:"answer"` //"re_sign_file_url_answer"),
	AnswerTokens            int                           `json:"answer_tokens"`
	ProviderResponseLatency float64                       `json:"provider_response_latency"`
	FromSource              string                        `json:"from_source"`
	FromEndUserID           string                        `json:"from_end_user_id"`
	FromAccountID           string                        `json:"from_account_id"`
	Feedbacks               []*FeedbackResponse           `json:"feedbacks"`
	WorkflowRunID           string                        `json:"workflow_run_id"`
	Annotation              *AnnotationResponse           `json:"annotation"`
	AnnotationHitHistory    *AnnotationHitHistoryResponse `json:"annotation_hit_history"`
	CreatedAt               int64                         `json:"created_at"`
	AgentThoughts           []*AgentThoughtResponse       `json:"agent_thoughts"`
	MessageFiles            []*MessageFileResponse        `json:"message_files"`
	Metadata                any                           `json:"metadata"` //"message_metadata_dict"),
	Status                  string                        `json:"status"`
	Error                   string                        `json:"error"`
	ParentMessageID         string                        `json:"parent_message_id"`
}

func NewMessageDetailResponse(args any) *MessageDetailResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(MessageDetailResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to MessageDetailResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(MessageDetailResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to MessageDetailResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if msg, ok := args.(*models.Message); ok && msg != nil {
		rsp := &MessageDetailResponse{
			ID:                      msg.ID,
			ConversationID:          msg.ConversationID,
			Inputs:                  msg.Inputs(),
			Query:                   msg.Query,
			Message:                 msg.MessageJson,
			MessageTokens:           msg.MessageTokens,
			Answer:                  msg.Answer,
			AnswerTokens:            msg.AnswerTokens,
			ProviderResponseLatency: msg.ProviderResponseLatency,
			FromSource:              msg.FromSource,
			FromEndUserID:           msg.FromEndUserID,
			FromAccountID:           msg.FromAccountID,
			WorkflowRunID:           msg.WorkflowRunID,
			// MessageFiles:            msg.MessageFiles(),
			Metadata:        msg.MessageMetadataDict(),
			Status:          msg.Status,
			Error:           msg.Error,
			ParentMessageID: msg.ParentMessageID,
		}
		for _, feedback := range msg.Feedbacks() {
			if rsp.Feedbacks == nil {
				rsp.Feedbacks = make([]*FeedbackResponse, 0)
			}
			tmp := &FeedbackResponse{
				Rating:        feedback.Rating,
				Content:       feedback.Content,
				FromSource:    feedback.FromSource,
				FromEndUserID: feedback.FromEndUserID,
			}
			from_account := feedback.FromAccount()
			if from_account != nil {
				tmp.FromAccount = &SimpleAccountResponse{
					ID:    from_account.ID,
					Name:  from_account.Name,
					Email: from_account.Email,
				}
			}
			rsp.Feedbacks = append(rsp.Feedbacks, tmp)
		}
		annotation := msg.Annotation()
		if annotation != nil {
			rsp.Annotation = &AnnotationResponse{
				ID:       annotation.ID,
				Question: annotation.Question,
				Content:  annotation.Content,
			}
			account := annotation.Account()
			if account != nil {
				rsp.Annotation.Account = &SimpleAccountResponse{
					ID:    account.ID,
					Name:  account.Name,
					Email: account.Email,
				}
			}
			if annotation.CreatedAt != nil {
				rsp.Annotation.CreatedAt = annotation.CreatedAt.Unix()
			}
		}
		annotation_hit_history := msg.AnnotationHitHistory()
		if annotation_hit_history != nil {
			rsp.AnnotationHitHistory = &AnnotationHitHistoryResponse{
				AnnotationID: annotation_hit_history.ID,
			}
			if annotation_hit_history.CreatedAt != nil {
				rsp.AnnotationHitHistory.CreatedAt = annotation_hit_history.CreatedAt.Unix()
			}
			annotation_create_account := annotation_hit_history.AnnotationCreateAccount()
			if annotation_create_account != nil {
				rsp.AnnotationHitHistory.AnnotationCreateAccount = &SimpleAccountResponse{
					ID:    annotation_create_account.ID,
					Name:  annotation_create_account.Name,
					Email: annotation_create_account.Email,
				}
			}
		}
		if msg.CreatedAt != nil {
			rsp.CreatedAt = msg.CreatedAt.Unix()
		}
		for _, agent_thought := range msg.AgentThoughts() {
			if rsp.AgentThoughts == nil {
				rsp.AgentThoughts = make([]*AgentThoughtResponse, 0)
			}
			tmp := &AgentThoughtResponse{
				ID:          agent_thought.ID,
				ChainID:     agent_thought.MessageChainID,
				MessageID:   agent_thought.MessageID,
				Position:    agent_thought.Position,
				Thought:     agent_thought.Thought,
				Tool:        agent_thought.Tool,
				ToolLabels:  agent_thought.ToolLabels,
				ToolInput:   agent_thought.ToolInput,
				Observation: agent_thought.Observation,
			}
			for _, v := range agent_thought.Files() {
				if tmp.Files == nil {
					tmp.Files = make([]string, 0)
				}
				tmp.Files = append(tmp.Files, cast.ToString(v))
			}
			if agent_thought.CreatedAt != nil {
				tmp.CreatedAt = agent_thought.CreatedAt.Unix()
			}
			rsp.AgentThoughts = append(rsp.AgentThoughts, tmp)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(MessageDetailResponse)
}

type SimpleModelConfigResponse struct {
	Model     any    `json:"model"` //attribute="model_dict"
	PrePrompt string `json:"pre_prompt"`
}

type FeedbackStatResponse struct {
	Like    int64 `json:"like"`
	Dislike int64 `json:"dislike"`
}
type StatusCountResponse struct {
	Success        int `json:"success"`
	Failed         int `json:"failed"`
	OartialSuccess int `json:"partial_success"`
}
type ModelConfigResponse struct {
	OpeningStatement   string `json:"opening_statement"`
	SuggestedQuestions any    `json:"suggested_questions"`
	Model              any    `json:"model"`
	UserInputForm      any    `json:"user_input_form"`
	PrePrompt          string `json:"pre_prompt"`
	AgentMode          any    `json:"agent_mode"`
}
type ConversationWithSummaryResponse struct {
	ID                   string                     `json:"id"`
	Status               string                     `json:"status"`
	FromSource           string                     `json:"from_source"`
	FromEndUserID        string                     `json:"from_end_user_id"`
	FromEndUserSessionID string                     `json:"from_end_user_session_id"`
	FromAccountID        string                     `json:"from_account_id"`
	FromAccountName      string                     `json:"from_account_name"`
	Name                 string                     `json:"name"`
	Summary              string                     `json:"summary"` //attribute="summary_or_query"
	ReadAt               int64                      `json:"read_at"`
	CreatedAt            int64                      `json:"created_at"`
	UpdatedAt            int64                      `json:"updated_at"`
	Annotated            bool                       `json:"annotated"`
	ModelConfig          *SimpleModelConfigResponse `json:"model_config"`
	MessageCount         int64                      `json:"message_count"`
	UserFeedbackStats    *FeedbackStatResponse      `json:"user_feedback_stats"`
	AdminFeedbackStats   *FeedbackStatResponse      `json:"admin_feedback_stats"`
	StatusCount          *StatusCountResponse       `json:"status_count"`
}
type ConversationWithSummaryPaginationResponse struct {
	Page    int                                `json:"page"`
	Limit   int                                `json:"limit"`
	Total   int                                `json:"total"`
	HasMore bool                               `json:"has_more"`
	Data    []*ConversationWithSummaryResponse `json:"data"`
}

func NewConversationWithSummaryPaginationResponse(args any) *ConversationWithSummaryPaginationResponse {
	rsp := new(ConversationWithSummaryPaginationResponse)
	if real_args, ok := args.(string); ok && real_args != "" {
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to ConversationWithSummaryPaginationResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to ConversationWithSummaryPaginationResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return rsp
}

type ConversationDetailResponse struct {
	ID                 string                `json:"id"`
	Status             string                `json:"status"`
	FromSource         string                `json:"from_source"`
	FromEndUserID      string                `json:"from_end_user_id"`
	FromAccountID      string                `json:"from_account_id"`
	CreatedAt          int64                 `json:"created_at"`
	UpdatedAt          int64                 `json:"updated_at"`
	Annotated          bool                  `json:"annotated"`
	Introduction       string                `json:"introduction"`
	ModelConfig        *ModelConfigResponse  `json:"model_config"`
	MessageCount       int64                 `json:"message_count"`
	UserFeedbackStats  *FeedbackStatResponse `json:"user_feedback_stats"`
	AdminFeedbackStats *FeedbackStatResponse `json:"admin_feedback_stats"`
}

func NewConversationDetailResponse(args any) *ConversationDetailResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(ConversationDetailResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to ConversationDetailResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(ConversationDetailResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to ConversationDetailResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if conversation, ok := args.(*models.Conversation); ok && conversation != nil {
		rsp := &ConversationDetailResponse{
			ID:                 conversation.ID,
			Status:             conversation.Status,
			FromSource:         conversation.FromSource,
			FromEndUserID:      conversation.FromEndUserID,
			FromAccountID:      conversation.FromAccountID,
			Annotated:          conversation.Annotated(),
			Introduction:       conversation.Introduction,
			MessageCount:       conversation.MessageCount(),
			UserFeedbackStats:  &FeedbackStatResponse{},
			AdminFeedbackStats: &FeedbackStatResponse{},
		}
		rsp.UserFeedbackStats.Like, rsp.UserFeedbackStats.Dislike = conversation.UserFeedbackStats()
		rsp.AdminFeedbackStats.Like, rsp.AdminFeedbackStats.Dislike = conversation.AdminFeedbackStats()
		if conversation.CreatedAt != nil {
			rsp.CreatedAt = conversation.CreatedAt.Unix()
		}
		if conversation.CreatedAt != nil {
			rsp.CreatedAt = conversation.CreatedAt.Unix()
		}
		model_config := conversation.ModelConfig()
		if len(model_config) > 0 {
			bindata, _ := json.Marshal(model_config)
			rsp.ModelConfig = &ModelConfigResponse{}
			err := json.Unmarshal(bindata, rsp.ModelConfig)
			if err != nil {
				mlog.Errorf("json unmarshal failed:%v", err)
				rsp.ModelConfig = nil
			}
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(ConversationDetailResponse)
}

type MessageInfiniteScrollPaginationResponse struct {
	Limit   int32                    `json:"limit"`
	HasMore bool                     `json:"has_more"`
	Data    []*MessageDetailResponse `json:"data"`
}

func NewMessageInfiniteScrollPaginationResponse(args any) *MessageInfiniteScrollPaginationResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(MessageInfiniteScrollPaginationResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to MessageInfiniteScrollPaginationResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(MessageInfiniteScrollPaginationResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to MessageInfiniteScrollPaginationResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(MessageInfiniteScrollPaginationResponse)
}
