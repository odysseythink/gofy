package workflow

import (
	"encoding/json"
	"iter"
	"maps"

	"github.com/odysseythink/gofy/backend/core/app/generator_response_convertes/base"
	appresponseentities "github.com/odysseythink/gofy/backend/entities/app/response"
)

type WorkflowAppGenerateResponseConvert struct {
	*base.BaseAppGeneratorResponseConvert[*appresponseentities.WorkflowAppBlockingResponse, iter.Seq[*appresponseentities.WorkflowAppStreamResponse]]
}

func (convert *WorkflowAppGenerateResponseConvert) ConvertBlockingFullResponse(blocking_response *appresponseentities.WorkflowAppBlockingResponse) map[string]any {
	/*
	   Convert blocking full response.
	   :param blocking_response: blocking response
	   :return:
	*/

	return blocking_response.ToDict(blocking_response)
}
func (convert *WorkflowAppGenerateResponseConvert) ConvertBlockingSimpleResponse(blocking_response *appresponseentities.WorkflowAppBlockingResponse) map[string]any {
	/*
	   Convert blocking simple response.
	   :param blocking_response: blocking response
	   :return:
	*/
	return convert.ConvertBlockingFullResponse(blocking_response)
}
func (convert *WorkflowAppGenerateResponseConvert) ConvertStreamFullResponse(stream_response iter.Seq[*appresponseentities.WorkflowAppStreamResponse]) iter.Seq[string] {
	/*
	   Convert stream full response.
	   :param stream_response: stream response
	   :return:
	*/
	return func(yield func(string) bool) {
		for chunk := range stream_response {
			if asp, ok := any(chunk).(*appresponseentities.WorkflowAppStreamResponse); ok {
				sub_stream_response := asp.StreamResponser
				if _, ok := any(sub_stream_response).(*appresponseentities.PingStreamResponse); ok {
					if !yield("ping") {
						return
					}
					continue
				}
				response_chunk := map[string]any{
					"event":           sub_stream_response.Event(),
					"workflow_run_id": asp.WorkflowRunID,
				}
				if sasp, ok := any(sub_stream_response).(*appresponseentities.ErrorStreamResponse); ok {
					data := convert.ErrorToStreamResponse(sasp.Err)
					// for k, v := range data {
					// 	response_chunk[k] = v
					// }
					maps.Copy(response_chunk, data)
				} else {
					// for k, v := range sub_stream_response.ToDict(sub_stream_response) {
					// 	response_chunk[k] = v
					// }
					maps.Copy(response_chunk, sub_stream_response.ToDict(sub_stream_response))
				}
				bindata, _ := json.Marshal(response_chunk)
				if !yield(string(bindata)) {
					return
				}
			}
		}
	}
}
func (convert *WorkflowAppGenerateResponseConvert) ConvertStreamSimpleResponse(stream_response iter.Seq[*appresponseentities.WorkflowAppStreamResponse]) iter.Seq[string] {
	/*
	   Convert stream simple response.
	   :param stream_response: stream response
	   :return:
	*/
	return func(yield func(language string) bool) {
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
				"workflow_run_id": chunk.WorkflowRunID,
			}
			if sasp, ok := any(sub_stream_response).(*appresponseentities.ErrorStreamResponse); ok {
				data := convert.ErrorToStreamResponse(sasp.Err)
				// for k, v := range data {
				// 	response_chunk[k] = v
				// }
				maps.Copy(response_chunk, data)
			} else if sasp, ok := any(sub_stream_response).(*appresponseentities.NodeStartStreamResponse); ok {
				// for k, v := range sasp.ToIgnoreDetailDict() {
				// 	response_chunk[k] = v
				// }
				maps.Copy(response_chunk, sasp.ToIgnoreDetailDict())
			} else if sasp, ok := any(sub_stream_response).(*appresponseentities.NodeFinishStreamResponse); ok {
				// for k, v := range sasp.ToIgnoreDetailDict() {
				// 	response_chunk[k] = v
				// }
				maps.Copy(response_chunk, sasp.ToIgnoreDetailDict())
			} else {
				for k, v := range sub_stream_response.ToDict(sub_stream_response) {
					response_chunk[k] = v
				}
				maps.Copy(response_chunk, sub_stream_response.ToDict(sub_stream_response))
			}
			bindata, _ := json.Marshal(response_chunk)
			if !yield(string(bindata)) {
				return
			}

		}
	}
}
