package advancedchat

import (
	"encoding/json"
	"iter"
	"maps"

	"github.com/odysseythink/gofy/backend/core/app/generator_response_convertes/base"
	appresponseentities "github.com/odysseythink/gofy/backend/entities/app/response"
)

type AdvancedChatAppGeneratorResponseConvert struct {
	*base.BaseAppGeneratorResponseConvert[*appresponseentities.ChatbotAppBlockingResponse, iter.Seq[*appresponseentities.ChatbotAppStreamResponse]]
}

func (convert *AdvancedChatAppGeneratorResponseConvert) ConvertBlockingFullResponse(blocking_response *appresponseentities.ChatbotAppBlockingResponse) map[string]any {
	/*
		Convert blocking full response.
		:param blocking_response: blocking response
		:return:
	*/

	return map[string]any{
		"event":           "message",
		"task_id":         blocking_response.TaskID(),
		"id":              blocking_response.Data.ID,
		"message_id":      blocking_response.Data.MessageID,
		"conversation_id": blocking_response.Data.ConversationID,
		"mode":            blocking_response.Data.Mode,
		"answer":          blocking_response.Data.Answer,
		"metadata":        blocking_response.Data.Metadata,
		"created_at":      blocking_response.Data.CreatedAt,
	}

}
func (convert *AdvancedChatAppGeneratorResponseConvert) ConvertBlockingSimpleResponse(blocking_response *appresponseentities.ChatbotAppBlockingResponse) map[string]any {
	/*
		Convert blocking simple response.
		:param blocking_response: blocking response
		:return:
	*/
	response := convert.ConvertBlockingFullResponse(blocking_response)
	metadata := response["metadata"].(map[string]any)
	response["metadata"] = convert.GetSimpleMetadata(metadata)

	return response

}
func (convert *AdvancedChatAppGeneratorResponseConvert) ConvertStreamFullResponse(
	stream_response iter.Seq[*appresponseentities.ChatbotAppStreamResponse],
) iter.Seq[string] {
	/*
		Convert stream full response.
		:param stream_response: stream response
		:return:
	*/
	return func(yield func(string) bool) {
		for chunk := range stream_response {
			sub_stream_response := chunk.StreamResponser

			if _, ok := any(sub_stream_response).(*appresponseentities.PingStreamResponse); ok {
				if !yield("ping") {
					return
				}
				continue
			}
			response_chunk := map[string]any{
				"event":           sub_stream_response.Event(),
				"conversation_id": chunk.ConversationID,
				"message_id":      chunk.MessageID,
				"created_at":      chunk.CreatedAt,
			}

			if real_sub_stream_response, ok := any(sub_stream_response).(*appresponseentities.ErrorStreamResponse); ok {
				maps.Copy(response_chunk, convert.ErrorToStreamResponse(real_sub_stream_response.Err))
			} else {
				maps.Copy(response_chunk, sub_stream_response.ToDict(sub_stream_response))
			}
			bindata, _ := json.Marshal(response_chunk)
			if !yield(string(bindata)) {
				return
			}
		}
	}
}
func (convert *AdvancedChatAppGeneratorResponseConvert) ConvertStreamSimpleResponse(
	stream_response iter.Seq[*appresponseentities.ChatbotAppStreamResponse],
) iter.Seq[string] {
	/*
		Convert stream simple response.
		:param stream_response: stream response
		:return:
	*/
	return func(yield func(string) bool) {
		for chunk := range stream_response {
			sub_stream_response := chunk.StreamResponser

			if _, ok := any(sub_stream_response).(*appresponseentities.PingStreamResponse); ok {
				if !yield("ping") {
					return
				}
				continue
			}
			response_chunk := map[string]any{
				"event":           sub_stream_response.Event(),
				"conversation_id": chunk.ConversationID,
				"message_id":      chunk.MessageID,
				"created_at":      chunk.CreatedAt,
			}
			if _, ok := any(sub_stream_response).(*appresponseentities.MessageEndStreamResponse); ok {
				sub_stream_response_dict := sub_stream_response.ToDict(sub_stream_response)
				metadata := map[string]any{}
				if _, ok := sub_stream_response_dict["metadata"]; ok {
					if _, ok := sub_stream_response_dict["metadata"].(map[string]any); ok {
						metadata = sub_stream_response_dict["metadata"].(map[string]any)
					}
				}

				sub_stream_response_dict["metadata"] = convert.GetSimpleMetadata(metadata)
				maps.Copy(response_chunk, convert.GetSimpleMetadata(metadata))
			} else if real_sub_stream_response, ok := any(sub_stream_response).(*appresponseentities.ErrorStreamResponse); ok {
				maps.Copy(response_chunk, convert.ErrorToStreamResponse(real_sub_stream_response.Err))
			} else if real_sub_stream_response, ok := any(sub_stream_response).(*appresponseentities.NodeStartStreamResponse); ok {
				maps.Copy(response_chunk, real_sub_stream_response.ToIgnoreDetailDict())
			} else if real_sub_stream_response, ok := any(sub_stream_response).(*appresponseentities.NodeFinishStreamResponse); ok {
				maps.Copy(response_chunk, real_sub_stream_response.ToIgnoreDetailDict())
			} else {
				maps.Copy(response_chunk, sub_stream_response.ToDict(sub_stream_response))
			}
			bindata, _ := json.Marshal(response_chunk)
			if !yield(string(bindata)) {
				return
			}
		}
	}
}
