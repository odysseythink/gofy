package callbacks

import (
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
)

type WorkflowCallback interface {
	OnEvent(event graphengineentities.GraphEngineEvent)
}
