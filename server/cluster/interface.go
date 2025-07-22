package cluster

import (
	"strings"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	gOnce     sync.Once
	gInstance *Cluster
)

func Instance() *Cluster {
	gOnce.Do(func() {
		gInstance = &Cluster{}
	})
	return gInstance
}

func GetExportModulename(modulename string) string {
	if modulename == "" {
		return ""
	}
	c := cases.Title(language.Dutch)
	ret := ""
	tmps := strings.Split(modulename, "_")
	for _, v := range tmps {
		ret += c.String(v)
	}
	return ret
}
