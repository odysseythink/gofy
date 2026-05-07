package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/gofy/backend/cluster"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	_ "github.com/odysseythink/gofy/backend/core/model_runtime/model_provides"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	_ "github.com/odysseythink/gofy/backend/events"
	eventhandlers "github.com/odysseythink/gofy/backend/events/event_handlers"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/mlog"
	"github.com/odysseythink/mrun"
	"google.golang.org/grpc/peer"
)

type ToolsService struct {
	pbapi.UnimplementedToolsServer
}

var (
	Tools ToolsService
)

func (s *ToolsService) GetToolLabelsList(ctx context.Context, in *pbapi.GetToolLabelsListRequest) (out *pbapi.GetToolLabelsListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetToolLabelsList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetToolLabelsListReply{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		tool_labels_list := services.ServiceGroupApp.Tools.ListToolLabels()
		bindata, _ := json.Marshal(tool_labels_list)
		out.ToolLabelsListStr = string(bindata)
	}()
	return
}
func (s *ToolsService) GetToolProviderList(ctx context.Context, in *pbapi.GetToolProviderListRequest) (out *pbapi.GetToolProviderListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetToolProviderList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetToolProviderListReply{}

	return
}
func (s *ToolsService) GetToolList(ctx context.Context, in *pbapi.GetToolListRequest) (out *pbapi.GetToolListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetToolList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetToolListReply{}

	return
}
func (s *ToolsService) GetMCPToolList(ctx context.Context, in *pbapi.GetMCPToolListRequest) (out *pbapi.GetMCPToolListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetMCPToolList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetMCPToolListReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	tools := services.ServiceGroupApp.Tools.RetrieveMCPTools(in.TenantId, false)
	datas := []map[string]any{}
	for _, tool := range tools {
		datas = append(datas, tool.ToDict())
	}
	bindata, _ := json.Marshal(datas)
	out.ToolsStr = string(bindata)
	return
}

// go build -o app.so -buildmode=plugin main.go
func (s *ToolsService) Init(args ...any) error {
	mlog.Info("tools init......")
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

func (s *ToolsService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *ToolsService) Destroy() {

}

func (s *ToolsService) UserData() any {
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

	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Tools})

	err = mrun.Run(&Tools)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
