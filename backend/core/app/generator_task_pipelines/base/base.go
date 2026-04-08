package base

import (
	"time"

	"mlib.com/gofy/server/core/exceptions"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	dbengine "mlib.com/gofy/server/db_engine"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appresponseentities "mlib.com/gofy/server/entities/app/response"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

// BaseGeneratorTaskPipeline 类的 Go 实现
type BaseGeneratorTaskPipeline[T1 interface {
	*appgeneratorentities.WorkflowAppGenerateEntity | *appgeneratorentities.EasyUIBasedAppGenerateEntity | *appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity | *appgeneratorentities.AdvancedChatAppGenerateEntity
}, T2 interface {
	*appqueueentities.MessageQueueMessage | *appqueueentities.WorkflowQueueMessage
}] struct {
	ApplicationGenerateEntity T1
	QueueManager              appqueueentities.AppQueueManager[T2]
	StartAt                   time.Time
	// _outputModerationHandler   *OutputModeration
	Stream bool
}

// NewBasedGenerateTaskPipeline 是构造函数
func New[T1 interface {
	*appgeneratorentities.WorkflowAppGenerateEntity | *appgeneratorentities.EasyUIBasedAppGenerateEntity | *appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity | *appgeneratorentities.AdvancedChatAppGenerateEntity
}, T2 interface {
	*appqueueentities.MessageQueueMessage | *appqueueentities.WorkflowQueueMessage
}](
	applicationGenerateEntity T1,
	queueManager appqueueentities.AppQueueManager[T2],
	stream bool,
) *BaseGeneratorTaskPipeline[T1, T2] {
	return &BaseGeneratorTaskPipeline[T1, T2]{
		ApplicationGenerateEntity: applicationGenerateEntity,
		QueueManager:              queueManager,
		StartAt:                   time.Now(),
		// _outputModerationHandler:   initOutputModeration(applicationGenerateEntity),
		Stream: stream,
	}
}

// HandleError 处理错误
func (p *BaseGeneratorTaskPipeline[T1, T2]) HandleError(event *appqueueentities.QueueErrorEvent, messageID string) error {
	mlog.Debugf("error: %v", event.Err)
	var err error

	if _, ok := any(event.Err).(*modelruntimeexceptions.InvokeAuthorizationError); ok {
		err = modelruntimeexceptions.NewInvokeAuthorizationError("Incorrect API key provided")
	} else {
		err = event.Err
	}

	if messageID == "" {
		return err
	}
	message := new(models.Message)
	err1 := dbengine.Instance().DB.Model(&models.Message{}).Where("id = ?", messageID).First(message).Error
	if err != nil {
		mlog.Errorf("get message from database failed:%v", err1)
		return err
	}
	errDesc := p.ErrorToDesc(err)
	message.Status = "error"
	message.Error = errDesc
	dbengine.Instance().DB.Save(&models.Message{ID: messageID, Status: "error", Error: errDesc})
	return err
}

// ErrorToDesc 将错误转换为描述
func (p *BaseGeneratorTaskPipeline[T1, T2]) ErrorToDesc(e error) string {
	switch e := e.(type) {
	case *exceptions.QuotaExceededError:
		return "Your quota for Dify Hosted Model Provider has been exhausted. Please go to Settings -> Model Provider to complete your own provider credentials."
	default:
		if e.Error() != "" {
			return e.Error()
		}
		return "Internal Server Error, please contact support."
	}
}

// ErrorToStreamResponse 将错误转换为流响应
func (p *BaseGeneratorTaskPipeline[T1, T2]) ErrorToStreamResponse(exp error) *appresponseentities.ErrorStreamResponse {
	task_id := ""
	switch real_entity := any(p.ApplicationGenerateEntity).(type) {
	case *appgeneratorentities.EasyUIBasedAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.ChatAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.CompletionAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.AdvancedChatAppGenerateEntity:
		task_id = real_entity.TaskID
	}
	return appresponseentities.NewErrorStreamResponse(task_id, exp)
}

// PingStreamResponse 返回 Ping 流响应
func (p *BaseGeneratorTaskPipeline[T1, T2]) PingStreamResponse() *appresponseentities.PingStreamResponse {
	task_id := ""
	switch real_entity := any(p.ApplicationGenerateEntity).(type) {
	case *appgeneratorentities.EasyUIBasedAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.ChatAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.CompletionAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.AdvancedChatAppGenerateEntity:
		task_id = real_entity.TaskID
	}
	return appresponseentities.NewPingStreamResponse(task_id)
}

// // InitOutputModeration 初始化输出审核
// func initOutputModeration(applicationGenerateEntity *AppGenerateEntity) *OutputModeration {
// 	appConfig := applicationGenerateEntity.AppConfig
// 	sensitiveWordAvoidance := appConfig.SensitiveWordAvoidance

// 	if sensitiveWordAvoidance != nil {
// 		return &OutputModeration{
// 			TenantID: appConfig.TenantID,
// 			AppID:    appConfig.AppID,
// 			Rule: ModerationRule{
// 				Type:   sensitiveWordAvoidance.Type,
// 				Config: sensitiveWordAvoidance.Config,
// 			},
// 			QueueManager: applicationGenerateEntity.AppConfig.QueueManager,
// 		}
// 	}
// 	return nil
// }

// // HandleOutputModerationWhenTaskFinished 在任务完成时处理输出审核
// func (p *BaseGeneratorTaskPipeline) HandleOutputModerationWhenTaskFinished(completion string) *string {
// 	if p._outputModerationHandler != nil {
// 		// 假设 stopThread 是一个方法
// 		p._outputModerationHandler.StopThread()

// 		completion = p._outputModerationHandler.ModerationCompletion(completion, false)

// 		p._outputModerationHandler = nil

// 		return &completion
// 	}
// 	return nil
// }
