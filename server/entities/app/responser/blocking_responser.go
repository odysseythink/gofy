package responser

type AppBlockingResponser interface {
	ToDict(AppBlockingResponser) map[string]any
	TaskID() string
}
