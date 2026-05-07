package generator

import (
	"iter"

	appresponseentities "github.com/odysseythink/gofy/backend/entities/app/response"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
)

type AppGeneratorResponseConverter[T1 interface {
	*appresponseentities.WorkflowAppBlockingResponse | *appresponseentities.ChatbotAppBlockingResponse
}, T2 iter.Seq[*appresponseentities.ChatbotAppStreamResponse] | iter.Seq[*appresponseentities.WorkflowAppStreamResponse]] interface {
	ConvertBlocking(AppGeneratorResponseConverter[T1, T2], T1 /*Union[AppBlockingResponse, iter.Seq[AppStreamResponse]]*/, appenumtypes.InvokeFrom) map[string]any //map[string]any, iter.Seq[string]
	ConvertStream(AppGeneratorResponseConverter[T1, T2], T2 /*Union[AppBlockingResponse, iter.Seq[AppStreamResponse]]*/, appenumtypes.InvokeFrom) iter.Seq[string] //map[string]any, iter.Seq[string]

	ConvertBlockingFullResponse(blocking_response T1) map[string]any
	ConvertBlockingSimpleResponse(blocking_response T1) map[string]any
	ConvertStreamFullResponse(T2) iter.Seq[string]
	ConvertStreamSimpleResponse(T2) iter.Seq[string]
}
