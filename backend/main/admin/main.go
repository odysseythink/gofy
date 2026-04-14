package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mrun"
	"mlib.com/gofy/server/cache"
	"mlib.com/gofy/server/cluster"
	dbengine "mlib.com/gofy/server/db_engine"
	eventhandlers "mlib.com/gofy/server/events/event_handlers"

	_ "mlib.com/gofy/server/events"
	"mlib.com/gofy/server/proto/pbapi"

	"github.com/odysseythink/mlog"
	_ "mlib.com/gofy/server/core/model_runtime/model_provides"
)

type AdminService struct {
	pbapi.UnimplementedAdminServer
}

var (
	Admin AdminService
)

// go build -o admin.so -buildmode=plugin main.go
func (s *AdminService) Init(args ...any) error {
	mlog.Info("admin init......")
	err := cache.Instance().Init()
	if err != nil {
		mlog.Errorf("cahe init failed:%v", err)
		return fmt.Errorf("cahe init failed:%v", err)
	}
	err = dbengine.Instance().Init()
	if err != nil {
		mlog.Errorf("mysql init failed:%v", err)
		return fmt.Errorf("mysql init failed:%v", err)
	}
	eventhandlers.Init()

	return nil
}

func (s *AdminService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *AdminService) Destroy() {

}

func (s *AdminService) UserData() any {
	return nil
}

func main() {
	var cfgfile string
	flag.StringVar(&cfgfile, "c", "", "choose config file.")
	flag.Parse()
	if cfgfile == "" {
		log.Println("usage: ./server -c config.yml")
		return
	}
	confy.SetConfigFile(cfgfile)
	confy.SetConfigType("yaml")
	err := confy.ReadInConfig()
	if err != nil {
		log.Printf("read config file(%s) failed: %v\n", cfgfile, err)
		return
	}
	{
		logpath := confy.GetWithDefault[string]("log.path", "logs")
		loglevel := confy.GetWithDefault[uint32]("log.log_level", 1)
		log.Println("******loglevel=", loglevel)
		if loglevel >= 4 {
			loglevel = 1
		}
		mlog.SetLogLevel(loglevel)
		mlog.SetLogDir(logpath)
	}
	defer mlog.Flush()
	confy.WatchConfig()

	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Admin})

	err = mrun.Run(&Admin)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
