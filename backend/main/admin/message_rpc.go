package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	dbengine "mlib.com/gofy/server/db_engine"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
)

func (s *AdminService) GetSuggestedQuestionMessage(ctx context.Context, in *pbapi.GetSuggestedQuestionMessageRequest) (out *pbapi.GetSuggestedQuestionMessageReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetSuggestedQuestionMessage call:%v", p.Addr.String(), in)

	out = &pbapi.GetSuggestedQuestionMessageReply{}
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
	if in.MessageId == "" {
		mlog.Error("missing message id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing message id")
		return
	}
	out.Questions = func() []string {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if _, ok := r.(*exceptions.MessageNotExistsError); ok {
					exp := httpexceptions.NewNotFound("Message not found")
					out.Exp = exp.ToPbHttpException(exp)
				} else if _, ok := r.(*exceptions.ConversationNotExistsError); ok {
					exp := httpexceptions.NewNotFound("Conversation not found")
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(*exceptions.ProviderTokenNotInitError); ok {
					exp := httpexceptions.NewProviderNotInitializeError(exp.Error())
					out.Exp = exp.ToPbHttpException(exp)
				} else if _, ok := r.(*exceptions.QuotaExceededError); ok {
					exp := httpexceptions.NewProviderQuotaExceededError()
					out.Exp = exp.ToPbHttpException(exp)
				} else if _, ok := r.(*exceptions.ModelCurrentlyNotSupportError); ok {
					exp := httpexceptions.NewProviderModelCurrentlyNotSupportError()
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(*modelruntimeexceptions.InvokeError); ok {
					exp := httpexceptions.NewCompletionRequestError(exp.Error())
					out.Exp = exp.ToPbHttpException(exp)
				} else if _, ok := r.(*exceptions.SuggestedQuestionsAfterAnswerDisabledError); ok {
					exp := httpexceptions.NewAppSuggestedQuestionsAfterAnswerDisabledError()
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		return services.ServiceGroupApp.AccountMessage.GetSuggestedQuestionsAfterAnswer(app_model, current_user, in.MessageId, appenumtypes.InvokeFrom_DEBUGGER)

	}()
	if out.Exp != nil {
		return
	}
	return
}
func (s *AdminService) ChatMessageList(ctx context.Context, in *pbapi.ChatMessageListRequest) (out *pbapi.ChatMessageListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.ChatMessageList call:%v", p.Addr.String(), in)

	out = &pbapi.ChatMessageListReply{}
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
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing conversation id")
		return
	}
	if _, err1 := uuid.FromString(in.ConversationId); err1 != nil {
		mlog.Error("conversation id is not a valid uuid:%v", err1)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("conversation id is not a valid uuid:" + err1.Error())
		return
	}
	if in.Limit < 1 || in.Limit > 100 {
		mlog.Errorf("limit[%d] must be in range[1, 100]", in.Limit)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("limit[%d] must be in range[1, 100]", in.Limit))
		return
	}
	rsp := func() *response.MessageInfiniteScrollPaginationResponse {
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
		conversation := new(models.Conversation)
		err := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ? and app_id = ?", in.ConversationId, app_model.ID).First(conversation).Error
		if err != nil {
			mlog.Errorf("get Conversation failed:%v", err)
			conversation = nil
		}

		if conversation == nil {
			panic(httpexceptions.NewNotFound("Conversation Not Exists."))
		}
		var history_messages []*models.Message
		if in.FirstId != "" {
			first_message := new(models.Message)
			err := dbengine.Instance().DB.Model(&models.Message{}).Where("conversation_id = ? and id = ?", conversation.ID, in.FirstId).First(first_message).Error
			if err != nil {
				mlog.Errorf("get Message failed:%v", err)
				first_message = nil
			}

			if first_message == nil {
				panic(httpexceptions.NewNotFound("First message not found"))
			}

			db := dbengine.Instance().DB.Model(&models.Message{})
			db = db.Where("conversation_id = ?", conversation.ID)
			db = db.Where("created_at < ?", first_message.CreatedAt)
			db = db.Where("id != ?", first_message.ID)
			err = db.Order("created_at DESC").Limit(int(in.Limit)).Find(&history_messages).Error
			if err != nil {
				mlog.Errorf("get Messages failed:%v", err)
				history_messages = nil
			}
		} else {
			err := dbengine.Instance().DB.Model(&models.Message{}).Where("conversation_id = ?", conversation.ID).Order("created_at DESC").Limit(int(in.Limit)).Find(&history_messages).Error
			if err != nil {
				mlog.Errorf("get Messages failed:%v", err)
				history_messages = nil
			}
		}
		has_more := false
		if len(history_messages) == int(in.Limit) {
			current_page_first_message := history_messages[len(history_messages)-1]
			var rest_count int64
			err := dbengine.Instance().DB.Model(&models.Message{}).Where("conversation_id = ? and created_at < ? and id != ?", conversation.ID, current_page_first_message.CreatedAt, current_page_first_message.ID).Order("created_at DESC").Count(&rest_count).Error
			if err != nil {
				mlog.Errorf("count Messages failed:%v", err)
			}

			if rest_count > 0 {
				has_more = true
			}
		}
		slices.Reverse(history_messages)
		rsp := &response.MessageInfiniteScrollPaginationResponse{
			Limit:   in.Limit,
			HasMore: has_more,
		}
		for _, history_message := range history_messages {
			if rsp.Data == nil {
				rsp.Data = make([]*response.MessageDetailResponse, 0)
			}
			rsp.Data = append(rsp.Data, response.NewMessageDetailResponse(history_message))
		}
		return rsp
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.MessageInfiniteScrollPaginationStr = string(bindata)
	return
}
func (s *AdminService) MessageFeedback(ctx context.Context, in *pbapi.MessageFeedbackRequest) (out *pbapi.MessageFeedbackReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.MessageFeedback call:%v", p.Addr.String(), in)

	out = &pbapi.MessageFeedbackReply{}
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
	if _, err1 := uuid.FromString(in.MessageId); err1 != nil {
		mlog.Errorf("message_id is invalid uuid:%v", err1)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("message_id is invalid uuid")
		return
	}
	if !slices.Contains([]string{"like", "dislike"}, in.Rating) {
		mlog.Error("rating[%s] must be \"like\" or \"dislike\"")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("rating[%s] must be \"like\" or \"dislike\"")
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
		message := new(models.Message)
		err := dbengine.Instance().DB.Model(&models.Message{}).Where("id =?", in.MessageId).Where("app_id =?", app_model.ID).First(message).Error
		if err != nil {
			mlog.Errorf("get message failed:%v", err)
			message = nil
		}
		if message == nil {
			panic(httpexceptions.NewNotFound("Message Not Exists."))
		}

		feedback := message.AdminFeedback()

		if in.Rating == "" && feedback != nil {
			dbengine.Instance().DB.Delete(feedback)
		} else if in.Rating != "" && feedback != nil {
			feedback.Rating = in.Rating
			dbengine.Instance().DB.Updates(&models.MessageFeedback{ID: feedback.ID, Rating: in.Rating})
		} else if in.Rating == "" && feedback == nil {
			panic(exceptions.NewValueError("rating cannot be None when feedback not exists"))
		} else {
			now := time.Now()
			feedback = &models.MessageFeedback{
				ID:             uuid.NewV4().String(),
				AppID:          app_model.ID,
				ConversationID: message.ConversationID,
				MessageID:      message.ID,
				Rating:         in.Rating,
				FromSource:     "admin",
				FromAccountID:  current_user.ID,
				CreatedAt:      &now,
				UpdatedAt:      &now,
			}
			dbengine.Instance().DB.Create(feedback)
		}
	}()
	if out.Exp != nil {
		return
	}
	return
}
func (s *AdminService) SetMessageAnnotation(ctx context.Context, in *pbapi.SetMessageAnnotationRequest) (out *pbapi.SetMessageAnnotationReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.SetMessageAnnotation call:%v", p.Addr.String(), in)

	out = &pbapi.SetMessageAnnotationReply{}
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
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if _, err1 := uuid.FromString(in.MessageId); err1 != nil {
		mlog.Errorf("message_id is invalid uuid:%v", err1)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("message_id is invalid uuid")
		return
	}
	if in.Question == "" {
		mlog.Error("missing question")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing question")
		return
	}
	if in.Answer == "" {
		mlog.Error("missing answer")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing answer")
		return
	}
	annotation_reply := map[string]any{}
	if in.AnnotationReplyStr != "" {
		err1 := json.Unmarshal([]byte(in.AnnotationReplyStr), &annotation_reply)
		if err1 != nil {
			mlog.Error("json unmarshal=%s to dict failed:%v", in.AnnotationReplyStr, err1)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp("annotation_reply must be dict")
			return
		}
	}
	rsp := func() *response.AnnotationResponse {
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
		args_dict := map[string]any{}
		bindata, _ := json.Marshal(in)
		json.Unmarshal(bindata, &args_dict)
		args_dict["annotation_reply"] = annotation_reply
		annotation := services.ServiceGroupApp.AppAnnotation.UpInsertAppAnnotationFromMessage(args_dict, app_model.ID, current_user)
		rsp := &response.AnnotationResponse{
			ID:       annotation.ID,
			Question: annotation.Question,
			Content:  annotation.Content,
		}
		tmpacc := annotation.Account()
		if tmpacc != nil {
			rsp.Account = &response.SimpleAccountResponse{
				ID:    tmpacc.ID,
				Name:  tmpacc.Name,
				Email: tmpacc.Email,
			}
		}
		if annotation.CreatedAt != nil {
			rsp.CreatedAt = annotation.CreatedAt.Unix()
		}
		return rsp
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.AnnotationStr = string(bindata)
	return
}
func (s *AdminService) MessageAnnotationCount(ctx context.Context, in *pbapi.MessageAnnotationCountRequest) (out *pbapi.MessageAnnotationCountReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.MessageAnnotationCount call:%v", p.Addr.String(), in)

	out = &pbapi.MessageAnnotationCountReply{}
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
		err1 := dbengine.Instance().DB.Model(&models.MessageAnnotation{}).Where("app_id =?", app_model.ID).Count(&out.Count).Error
		if err1 != nil {
			mlog.Error("count MessageAnnotation failed:", err1.Error())
		}
	}()
	if out.Exp != nil {
		return
	}
	return
}
func (s *AdminService) GetMessage(ctx context.Context, in *pbapi.GetMessageRequest) (out *pbapi.GetMessageReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetMessage call:%v", p.Addr.String(), in)

	out = &pbapi.GetMessageReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.MessageId == "" {
		mlog.Error("missing message id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing message id")
		return
	}
	rsp := func() *response.MessageDetailResponse {
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
		message := new(models.Message)
		err := dbengine.Instance().DB.Model(&models.Message{}).Where("id=? and app_id=?", in.MessageId, in.AppId).First(message).Error
		if err != nil {
			mlog.Errorf("get Message failed:%v", err)
			message = nil
		}

		if message == nil {
			panic(httpexceptions.NewNotFound("Message Not Exists."))
		}
		return response.NewMessageDetailResponse(message)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.MessageDetailStr = string(bindata)
	return
}
