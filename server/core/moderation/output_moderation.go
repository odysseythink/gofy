package moderation

import (
	"context"
	"sync"
)

// 定义所需的类型和结构体
type ModerationRule struct {
	Type   string
	Config map[string]interface{}
}

type OutputModeration struct {
	TenantID string
	AppID    string
	Rule     ModerationRule
	// QueueManager  appqueueentities.AppQueueManager
	Thread        sync.WaitGroup
	ThreadRunning bool
	Buffer        string
	IsFinalChunk  bool
	FinalOutput   *string
	ModelConfig   map[string]any
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

// ShouldDirectOutput 判断是否直接输出
func (o *OutputModeration) ShouldDirectOutput() bool {
	return o.FinalOutput != nil
}

// GetFinalOutput 获取最终输出
func (o *OutputModeration) GetFinalOutput() string {
	if o.FinalOutput != nil {
		return *o.FinalOutput
	}
	return ""
}

// AppendNewToken 添加新令牌
func (o *OutputModeration) AppendNewToken(token string) {
	o.Buffer += token

	o.Thread.Add(1)
	// go o.StartThread(viper.GetInt("moderation.buffer_size"))

}

// // ModerationCompletion 处理完成的审核
// func (o *OutputModeration) ModerationCompletion(completion string, publicEvent bool) string {
// 	o.Buffer = completion
// 	o.IsFinalChunk = true

// 	result := o.Moderation(o.TenantID, o.AppID, completion)

// 	if result == nil || !result.Flagged {
// 		return completion
// 	}

// 	var finalOutput string
// 	if result.Action == DIRECT_OUTPUT {
// 		finalOutput = result.PresetResponse
// 	} else {
// 		finalOutput = result.Text
// 	}

// 	if publicEvent {
// 		o.QueueManager.Publish(QueueMessageReplaceEvent{Text: finalOutput}, TASK_PIPELINE)
// 	}

// 	return finalOutput
// }

// // StartThread 启动线程
// func (o *OutputModeration) StartThread(bufferSize int) {
// 	defer o.Thread.Done()

// 	bufferSize = bufferSize
// 	if bufferSize <= 0 {
// 		bufferSize = viper.GetInt("moderation.buffer_size")
// 	}

// 	o.worker(flaskApp, bufferSize)
// }

// // StopThread 停止线程
// func (o *OutputModeration) StopThread() {
// 	if o.Thread != nil && o.ThreadRunning {
// 		o.ThreadRunning = false
// 		o.Thread.Wait()
// 	}
// }

// // Worker 工作线程
// func (o *OutputModeration) worker(bufferSize int) {
// 	ctx := flaskApp.AppContext()
// 	ctx.Push()

// 	currentLength := 0
// 	for o.ThreadRunning {
// 		moderationBuffer := o.Buffer
// 		bufferLength := len(moderationBuffer)
// 		if !o.IsFinalChunk {
// 			chunkLength := bufferLength - currentLength
// 			if chunkLength >= 0 && chunkLength < bufferSize {
// 				time.Sleep(1 * time.Second)
// 				continue
// 			}
// 		}

// 		currentLength = bufferLength

// 		result := o.Moderation(o.TenantID, o.AppID, moderationBuffer)

// 		if result == nil || !result.Flagged {
// 			continue
// 		}

// 		var finalOutput string
// 		if result.Action == DIRECT_OUTPUT {
// 			finalOutput = result.PresetResponse
// 			o.FinalOutput = &finalOutput
// 		} else {
// 			finalOutput = result.Text + o.Buffer[len(moderationBuffer):]
// 		}

// 		if o.ThreadRunning {
// 			o.QueueManager.Publish(QueueMessageReplaceEvent{Text: finalOutput}, TASK_PIPELINE)
// 		}

// 		if result.Action == DIRECT_OUTPUT {
// 			break
// 		}
// 	}
// 	ctx.Pop()
// }

// // Moderation 审核
// func (o *OutputModeration) Moderation(tenantID string, appID string, moderationBuffer string) *ModerationOutputsResult {
// 	moderationFactory := ModerationFactory{
// 		Name:     o.Rule.Type,
// 		AppID:    appID,
// 		TenantID: tenantID,
// 		Config:   o.Rule.Config,
// 	}

// 	result, err := moderationFactory.ModerationForOutputs(moderationBuffer)
// 	if err != nil {
// 		logger.Exception("Moderation Output error, app_id: %s", appID)
// 		return nil
// 	}
// 	return &result
// }
