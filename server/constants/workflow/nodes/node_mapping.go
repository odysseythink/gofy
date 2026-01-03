package nodes

import (
	"errors"
	"log"

	"mlib.com/gofy/server/core/workflow/nodes/base"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
)

const (
	LATEST_VERSION = "latest"
)

var (
	NODE_TYPE_CLASSES_MAPPING = map[nodesenumtypes.NodeType]map[string]base.Noder{}
)

func Regist(n base.Noder) error {
	if n == nil {
		log.Println("missing noder")
		return errors.New("missing noder")
	}

}
