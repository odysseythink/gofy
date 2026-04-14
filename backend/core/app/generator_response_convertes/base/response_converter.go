package base

import (
	"fmt"
	"iter"
	"reflect"
	"slices"

	"github.com/odysseythink/mlog"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appresponseentities "mlib.com/gofy/server/entities/app/response"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
)

type BaseAppGeneratorResponseConvert[T1 interface {
	*appresponseentities.WorkflowAppBlockingResponse | *appresponseentities.ChatbotAppBlockingResponse
}, T2 iter.Seq[*appresponseentities.ChatbotAppStreamResponse] | iter.Seq[*appresponseentities.WorkflowAppStreamResponse]] struct {
}

func (c *BaseAppGeneratorResponseConvert[T1, T2]) ConvertBlocking(
	cls appgeneratorentities.AppGeneratorResponseConverter[T1, T2],
	response T1, /*Union[AppBlockingResponse, Generator[AppStreamResponse, Any, None]]*/
	invoke_from appenumtypes.InvokeFrom,
) map[string]any { //map[string]any, iter.Seq[string]
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_DEBUGGER, appenumtypes.InvokeFrom_SERVICE_API}, invoke_from) {
		return cls.ConvertBlockingFullResponse(response)
	} else {
		return cls.ConvertBlockingSimpleResponse(response)
	}
}

func (c *BaseAppGeneratorResponseConvert[T1, T2]) ConvertStream(
	cls appgeneratorentities.AppGeneratorResponseConverter[T1, T2],
	response T2, /*Union[AppBlockingResponse, Generator[AppStreamResponse, Any, None]]*/
	invoke_from appenumtypes.InvokeFrom,
) iter.Seq[string] { //map[string]any, iter.Seq[string]
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_DEBUGGER, appenumtypes.InvokeFrom_SERVICE_API}, invoke_from) {
		if realrsp, ok := any(response).(iter.Seq[*appresponseentities.WorkflowAppStreamResponse]); ok {
			iterfunc := func(iter.Seq[*appresponseentities.WorkflowAppStreamResponse]) iter.Seq[string] {
				return func(yield func(string) bool) {
					for chunk := range cls.ConvertStreamFullResponse(response) {
						if chunk == "ping" {
							if !yield(fmt.Sprintf("event: %s\n\n", chunk)) {
								return
							}
						} else {
							if !yield(fmt.Sprintf("data: %s\n\n", chunk)) {
								return
							}
						}
					}
				}
			}
			return iterfunc(realrsp)
		} else if realrsp, ok := any(response).(iter.Seq[*appresponseentities.ChatbotAppStreamResponse]); ok {
			iterfunc := func(iter.Seq[*appresponseentities.ChatbotAppStreamResponse]) iter.Seq[string] {
				return func(yield func(string) bool) {
					for chunk := range cls.ConvertStreamFullResponse(response) {
						if chunk == "ping" {
							if !yield(fmt.Sprintf("event: %s\n\n", chunk)) {
								return
							}
						} else {
							if !yield(fmt.Sprintf("data: %s\n\n", chunk)) {
								return
							}
						}
					}
				}
			}
			return iterfunc(realrsp)
		}

	} else {
		if realrsp, ok := any(response).(iter.Seq[*appresponseentities.WorkflowAppStreamResponse]); ok {
			iterfunc := func(iter.Seq[*appresponseentities.WorkflowAppStreamResponse]) iter.Seq[string] {
				return func(yield func(string) bool) {
					for chunk := range cls.ConvertStreamSimpleResponse(response) {
						if chunk == "ping" {
							if !yield(fmt.Sprintf("event: %s\n\n", chunk)) {
								return
							}
						} else {
							if !yield(fmt.Sprintf("data: %s\n\n", chunk)) {
								return
							}
						}
					}
				}
			}
			return iterfunc(realrsp)
		} else if realrsp, ok := any(response).(iter.Seq[*appresponseentities.ChatbotAppStreamResponse]); ok {
			iterfunc := func(iter.Seq[*appresponseentities.ChatbotAppStreamResponse]) iter.Seq[string] {
				return func(yield func(string) bool) {
					for chunk := range cls.ConvertStreamSimpleResponse(response) {
						if chunk == "ping" {
							if !yield(fmt.Sprintf("event: %s\n\n", chunk)) {
								return
							}
						} else {
							if !yield(fmt.Sprintf("data: %s\n\n", chunk)) {
								return
							}
						}
					}
				}
			}
			return iterfunc(realrsp)
		}
	}
	return nil
}

func (c *BaseAppGeneratorResponseConvert[T1, T2]) GetSimpleMetadata(metadata map[string]any) map[string]any {
	/*
		Get simple metadata.
		:param metadata: metadata
		:return:
	*/
	// show_retrieve_source
	updated_resources := []map[string]any{}
	if _, ok := metadata["retriever_resources"]; ok {
		if _, ok := metadata["retriever_resources"].([]any); ok {
			for _, resource := range metadata["retriever_resources"].([]any) {
				if _, ok := resource.(map[string]any); ok {
					tmp := map[string]any{
						"segment_id":    "",
						"position":      resource.(map[string]any)["position"],
						"document_name": resource.(map[string]any)["document_name"],
						"score":         resource.(map[string]any)["score"],
						"content":       resource.(map[string]any)["content"],
					}

					if _, ok := resource.(map[string]any)["segment_id"]; ok {
						if _, ok := resource.(map[string]any)["segment_id"].(string); ok {
							tmp["segment_id"] = resource.(map[string]any)["segment_id"].(string)
						}
					}
					updated_resources = append(updated_resources, tmp)
				}
			}
			metadata["retriever_resources"] = updated_resources
		}
	}
	// show annotation reply
	delete(metadata, "annotation_reply")

	// show usage
	delete(metadata, "usage")

	return metadata
}

func (c *BaseAppGeneratorResponseConvert[T1, T2]) ErrorToStreamResponse(e error) map[string]any {
	/*
		Error to stream response.
		:param e: exception
		:return:
	*/
	error_responses := map[string]map[string]any{
		"ValueError":                {"code": "invalid_param", "status": 400},
		"ProviderTokenNotInitError": {"code": "provider_not_initialize", "status": 400},
		"QuotaExceededError": {
			"code":    "provider_quota_exceeded",
			"message": "Your quota for Gofy Hosted Model Provider has been exhausted. Please go to Settings -> Model Provider to complete your own provider credentials.",
			"status":  400,
		},
		"ModelCurrentlyNotSupportError": {"code": "model_currently_not_support", "status": 400},
		"InvokeError":                   {"code": "completion_request_error", "status": 400},
	}

	// Determine the response based on the type of exception
	data := map[string]any{}
	for k, v := range error_responses {
		if reflect.ValueOf(e).Type().Name() == k {
			data = v
		}
	}
	if len(data) > 0 {
		data["message"] = e.Error()
	} else {
		mlog.Errorf(e.Error())
		data = map[string]any{
			"code":    "internal_server_error",
			"message": "Internal Server Error, please contact support.",
			"status":  500,
		}
	}
	return data
}
