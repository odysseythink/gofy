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
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/main/plugins/services"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
	"github.com/odysseythink/mrun"
	"google.golang.org/grpc/peer"
)

type PluginsService struct {
	pbapi.UnimplementedPluginsServer
}

var (
	Plugins PluginsService
)

func (s *PluginsService) FetchPreferences(ctx context.Context, in *pbapi.FetchPreferencesRequest) (out *pbapi.FetchPreferencesReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] plugin.FetchPreferences call:%#v", p.Addr.String(), in)

	out = &pbapi.FetchPreferencesReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	permission := services.ServiceGroupApp.PluginPermission.GetPermission(in.TenantId)
	out.PermissionDict = map[string]string{
		"install_permission": string(models.TenantPluginInstallPermission_EVERYONE),
		"debug_permission":   string(models.TenantPluginDebugPermission_EVERYONE),
	}

	if permission != nil {
		out.PermissionDict["install_permission"] = string(permission.InstallPermission)
		out.PermissionDict["debug_permission"] = string(permission.DebugPermission)
	}
	auto_upgrade := services.ServiceGroupApp.PluginAutoUpgrade.GetStrategy(in.TenantId)
	auto_upgrade_dict := map[string]any{
		"strategy_setting":    models.TenantPluginAutoUpgradeStrategySetting_DISABLED,
		"upgrade_time_of_day": 0,
		"upgrade_mode":        models.TenantPluginAutoUpgradeStrategyUpgradeMode_EXCLUDE,
		"exclude_plugins":     []string{},
		"include_plugins":     []string{},
	}

	if auto_upgrade != nil {
		_ = json.Unmarshal(auto_upgrade.ExcludePlugins, &auto_upgrade.ExcludePluginList)
		_ = json.Unmarshal(auto_upgrade.IncludePlugins, &auto_upgrade.IncludePluginList)
		auto_upgrade_dict = map[string]any{
			"strategy_setting":    auto_upgrade.StrategySetting,
			"upgrade_time_of_day": auto_upgrade.UpgradeTimeOfDay,
			"upgrade_mode":        auto_upgrade.UpgradeMode,
			"exclude_plugins":     auto_upgrade.ExcludePluginList,
			"include_plugins":     auto_upgrade.IncludePluginList,
		}
	}
	bindata, _ := json.Marshal(auto_upgrade_dict)
	out.AutoUpgradeDictStr = string(bindata)

	return
}
func (s *PluginsService) FetchInstallTasks(ctx context.Context, in *pbapi.FetchInstallTasksRequest) (out *pbapi.FetchInstallTasksReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] plugin.FetchInstallTasks call:%#v", p.Addr.String(), in)

	out = &pbapi.FetchInstallTasksReply{}
	tasks := services.ServiceGroupApp.Plugin.FetchInstallTasks(in.TenantId, int(in.Page), int(in.PageSize))
	if len(tasks) > 0 {
		bindata, _ := json.Marshal(tasks)
		out.TasksStr = string(bindata)
	} else {
		out.TasksStr = "[]"
	}
	return
}

func (s *PluginsService) DispatchTextEmbeddingInvoke(ctx context.Context, in *pbapi.DispatchTextEmbeddingInvokeRequest) (out *pbapi.DispatchTextEmbeddingInvokeReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] plugin.DispatchTextEmbeddingInvoke call:%#v", p.Addr.String(), in)

	out = &pbapi.DispatchTextEmbeddingInvokeReply{}

	return
}

// go build -o app.so -buildmode=plugin main.go
func (s *PluginsService) Init(args ...any) error {
	mlog.Info("plugin init......")
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

	return nil
}

func (s *PluginsService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *PluginsService) Destroy() {

}

func (s *PluginsService) UserData() any {
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

	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Plugins})

	err = mrun.Run(&Plugins)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
