package callbacks

import (
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
)

type WorkflowCallback interface {
	OnEvent(event graphengineentities.GraphEngineEvent)
}
