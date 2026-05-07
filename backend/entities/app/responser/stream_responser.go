package responser

import appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"

type StreamResponser interface {
	Event() appenumtypes.StreamEventType
	ToDict(StreamResponser) map[string]any
	TaskID() string
}
