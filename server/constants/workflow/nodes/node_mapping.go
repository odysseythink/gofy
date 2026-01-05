package nodes

import (
	"errors"
	"fmt"
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
	if _, ok := NODE_TYPE_CLASSES_MAPPING[n.Type()]; ok {
		log.Printf("node(%s) already regist\n", n.Type())
		return fmt.Errorf("node(%s) already regist", n.Type())
	}
	NODE_TYPE_CLASSES_MAPPING[n.Type()] = map[string]base.Noder{LATEST_VERSION: n, "1": n}
	return nil
}
