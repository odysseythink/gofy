package graphengine

type RouteNodeStateStatusType string

const (
	RouteNodeStateStatus_RUNNING   RouteNodeStateStatusType = "running"
	RouteNodeStateStatus_SUCCESS   RouteNodeStateStatusType = "success"
	RouteNodeStateStatus_FAILED    RouteNodeStateStatusType = "failed"
	RouteNodeStateStatus_PAUSED    RouteNodeStateStatusType = "paused"
	RouteNodeStateStatus_EXCEPTION RouteNodeStateStatusType = "exception"
)
