package responser

import appenumtypes "mlib.com/gofy/server/enum_types/app"

type StreamResponser interface {
	Event() appenumtypes.StreamEventType
	ToDict(StreamResponser) map[string]any
	TaskID() string
}
