package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	pbexceptions "mlib.com/gofy/server/proto/exceptions"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/mlog"
)

func (s *AdminService) GetChatConversationPagination(ctx context.Context, in *pbapi.GetChatConversationPaginationRequest) (out *pbapi.GetChatConversationPaginationReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetChatConversationPagination call:%#v", p.Addr.String(), in)

	out = &pbapi.GetChatConversationPaginationReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	var startAt time.Time
	if in.Start != "" {
		startAt, err = time.Parse("2006-01-02 15:04", in.Start)
		if err != nil {
			mlog.Errorf("start=%s is invalid:%v", in.Start, err)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("start=%s is invalid:%v", in.Start, err))
			return
		}
	}
	var endAt time.Time
	if in.End != "" {
		endAt, err = time.Parse("2006-01-02 15:04", in.End)
		if err != nil {
			mlog.Errorf("end=%s is invalid:%v", in.End, err)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("end=%s is invalid:%v", in.End, err))
			return
		}
	}
	if !slices.Contains([]string{"annotated", "not_annotated", "all"}, in.AnnotationStatus) {
		mlog.Errorf("annotation_status=%s must be \"annotated\", \"not_annotated\" or \"all\"", in.AnnotationStatus)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("annotation_status=%s must be \"annotated\", \"not_annotated\" or \"all\"", in.AnnotationStatus))
		return
	}
	if in.MessageCountGte != 0 && (in.MessageCountGte < 1 || in.MessageCountGte > 99999) {
		mlog.Errorf("message_count_gte=%d is beyond range(1, 99999)", in.MessageCountGte)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("message_count_gte=%d is beyond range(1, 99999)", in.MessageCountGte))
		return
	}
	if in.Page < 1 || in.Page > 99999 {
		mlog.Errorf("page=%d is beyond range(1, 99999)", in.Page)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("page=%d is beyond range(1, 99999)", in.Page))
		return
	}
	if in.Limit < 1 || in.Page > 100 {
		mlog.Errorf("limit=%d is beyond range(1, 100)", in.Limit)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("limit=%d is beyond range(1, 100)", in.Limit))
		return
	}
	if !slices.Contains([]string{"created_at", "-created_at", "updated_at", "-updated_at"}, in.SortBy) {
		mlog.Errorf("sort_by=%s must be \"created_at\", \"-created_at\", \"updated_at\", \"-updated_at\"", in.SortBy)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("sort_by=%s must be \"created_at\", \"-created_at\", \"updated_at\", \"-updated_at\"", in.SortBy))
		return
	}

	rsp := func() *response.ConversationWithSummaryPaginationResponse {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		query := dbengine.Instance().DB.Model(&models.Conversation{}).Where("app_id =?", app_model.ID)

		if in.Keyword != "" {
			keyword_filter := "%" + in.Keyword + "%"
			subquery := dbengine.Instance().DB.Model(&models.Conversation{}).
				Select("id as conversation_id", "end_user.session_id as from_end_user_session_id").
				Joins("LEFT JOIN end_users ON conversations.from_end_user_id = end_users.id")

			query = query.Joins("JOIN messages ON messages.conversation_id = conversations.id").
				Joins("JOIN (?) AS subquery ON subquery.conversation_id = conversations.id", subquery).
				Where(
					"messages.query LIKE ? OR messages.answer LIKE ? OR "+
						"conversations.name LIKE ? OR conversations.introduction LIKE ? OR "+
						"subquery.from_end_user_session_id LIKE ?",
					keyword_filter, keyword_filter, keyword_filter, keyword_filter, keyword_filter,
				).
				Group("conversations.id")
		}

		if !startAt.IsZero() {
			switch in.SortBy {
			case "updated_at":
				fallthrough
			case "-updated_at":
				query = query.Where("conversations.updated_at >= ?", startAt)
			case "created_at":
				fallthrough
			case "-created_at":
				query = query.Where("conversations.created_at >= ?", startAt)
			}
		}
		if !endAt.IsZero() {
			endAt.Add(59 * time.Second)
			switch in.SortBy {
			case "updated_at":
				fallthrough
			case "-updated_at":
				query = query.Where("conversations.updated_at <= ?", endAt)
			case "created_at":
				fallthrough
			case "-created_at":
				query = query.Where("conversations.created_at <= ?", endAt)
			}
		}
		switch in.AnnotationStatus {
		case "annotated":
			query = query.Joins("JOIN message_annotations ON message_annotations.conversation_id = conversations.id")
		case "not_annotated":
			query = query.Joins("LEFT JOIN message_annotations ON message_annotations.conversation_id = conversations.id").
				Group("conversations.id").
				Having("COUNT(message_annotations.id) = 0")
		}

		if in.MessageCountGte >= 1 {
			query = query.Joins("JOIN messages ON messages.conversation_id = conversations.id").
				Group("conversations.id").
				Having("COUNT(messages.id) >= ?", in.MessageCountGte)
		}
		if app_model.Mode == models.AppMode_ADVANCED_CHAT {
			query = query.Where("invoke_from != ?", "debugger")
		}
		// Handle sorting
		switch in.SortBy {
		case "created_at":
			query = query.Order("created_at ASC")
		case "-created_at":
			query = query.Order("created_at DESC")
		case "updated_at":
			query = query.Order("updated_at ASC")
		case "-updated_at":
			query = query.Order("updated_at DESC")
		default:
			query = query.Order("created_at DESC")
		}
		// Execute query and count total
		var total int64
		err = query.Count(&total).Error
		if err != nil {
			mlog.Errorf("count query failed:%v", err)
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusBadRequest,
				Code:    "count_conversation_error",
				Message: "count conversation failed:" + err.Error(),
			}
			return nil
		}
		var conversations []*models.Conversation
		err = query.Find(&conversations).Error
		if err != nil {
			mlog.Errorf("query failed:%v", err)
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusBadRequest,
				Code:    "query_conversation_error",
				Message: "query conversation failed:" + err.Error(),
			}
			return nil
		}
		rsp := &response.ConversationWithSummaryPaginationResponse{
			Page:    int(in.Page),
			Limit:   int(in.Limit),
			Total:   int(total),
			HasMore: true,
			Data:    make([]*response.ConversationWithSummaryResponse, 0),
		}
		for _, conversation := range conversations {
			tmpconver := &response.ConversationWithSummaryResponse{
				ID:                   conversation.ID,
				Status:               conversation.Status,
				FromSource:           conversation.FromSource,
				FromEndUserID:        conversation.FromEndUserID,
				FromEndUserSessionID: conversation.FromEndUserSessionID(),
				FromAccountID:        conversation.FromAccountID,
				FromAccountName:      conversation.FromAccountName(),
				Name:                 conversation.Name,
				Summary:              conversation.Summary,
				Annotated:            conversation.Annotated(),
				ModelConfig: &response.SimpleModelConfigResponse{
					Model: conversation.ModelConfig(),
				},
				MessageCount:       conversation.MessageCount(),
				UserFeedbackStats:  &response.FeedbackStatResponse{},
				AdminFeedbackStats: &response.FeedbackStatResponse{},
				StatusCount:        &response.StatusCountResponse{},
			}
			tmpconver.UserFeedbackStats.Like, tmpconver.UserFeedbackStats.Dislike = conversation.UserFeedbackStats()
			tmpconver.AdminFeedbackStats.Like, tmpconver.AdminFeedbackStats.Dislike = conversation.AdminFeedbackStats()
			if conversation.ReadAt != nil {
				tmpconver.ReadAt = conversation.ReadAt.Unix()
			}
			if conversation.CreatedAt != nil {
				tmpconver.CreatedAt = conversation.CreatedAt.Unix()
			}
			if conversation.UpdatedAt != nil {
				tmpconver.UpdatedAt = conversation.UpdatedAt.Unix()
			}
			status_count := conversation.StatusCount()
			if len(status_count) > 0 {
				bindata, _ := json.Marshal(status_count)
				json.Unmarshal(bindata, tmpconver.StatusCount)
			}
			rsp.Data = append(rsp.Data, tmpconver)
		}
		return rsp
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.ConversationWithSummaryPaginationStr = string(bindata)
	return
}
func (s *AdminService) ChatConversationDetail(ctx context.Context, in *pbapi.ChatConversationDetailRequest) (out *pbapi.ChatConversationDetailReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.ChatConversationDetail call:%#v", p.Addr.String(), in)

	out = &pbapi.ChatConversationDetailReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.ConversationId == "" {
		mlog.Error("missing conversation id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	rsp := func() *response.ConversationDetailResponse {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		conversation := services.ServiceGroupApp.AccountConversation.GetConversation1(current_user, app_model, in.ConversationId)
		return response.NewConversationDetailResponse(conversation)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.ConversationDetailStr = string(bindata)
	return
}
func (s *AdminService) DelChatConversation(ctx context.Context, in *pbapi.DelChatConversationRequest) (out *pbapi.DelChatConversationReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.DelChatConversation call:%#v", p.Addr.String(), in)

	out = &pbapi.DelChatConversationReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.ConversationId == "" {
		mlog.Error("missing conversation id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		var count int64
		err := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ? and app_id=?", in.ConversationId, app_model.ID).Count(&count).Error
		if err != nil {
			mlog.Errorf("get Conversation from mysql failed:%v", err)
		}
		if count == 0 {
			panic(httpexceptions.NewNotFound("Conversation Not Exists."))
		}
		dbengine.Instance().DB.Updates(&models.Conversation{ID: in.ConversationId, IsDeleted: true})
	}()
	if out.Exp != nil {
		return
	}
	return
}
