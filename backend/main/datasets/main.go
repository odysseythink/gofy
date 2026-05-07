package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/gofy/backend/cluster"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	providermanager "github.com/odysseythink/gofy/backend/core/manageres/provider_manager"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	pluginentities "github.com/odysseythink/gofy/backend/entities/plugin"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/mlog"
	"github.com/odysseythink/mrun"
	"google.golang.org/grpc/peer"
)

type DatasetsService struct {
	pbapi.UnimplementedDatasetsServer
}

var (
	Datasets DatasetsService
)

func (s *DatasetsService) DatasetList(ctx context.Context, in *pbapi.DatasetListRequest) (out *pbapi.DatasetListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] datasets.DatasetList call:%#v", p.Addr.String(), in)
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Limit <= 0 {
		in.Limit = 20
	}

	out = &pbapi.DatasetListReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	var datasets []*models.Dataset
	var total int64
	if len(in.Ids) > 0 {
		datasets, total = services.ServiceGroupApp.Dataset.GetDatasetsByIDs(in.Ids, current_user.CurrentTenantID())
	} else {
		datasets, total = services.ServiceGroupApp.Dataset.GetDatasets(in.Page, in.Limit, current_user.CurrentTenantID(), current_user, in.Keyword, in.TagIds, in.IncludeAll)
	}
	// check embedding setting
	configurations := (&providermanager.ProviderManager{}).GetConfigurations(current_user.CurrentTenantID())

	embedding_models := (&providermanager.ProviderConfigurationsManager{}).GetModels(configurations, "", modelruntimeenumtypes.Model_TEXT_EMBEDDING, true)

	model_names := []string{}
	for _, embedding_model := range embedding_models {
		model_names = append(model_names, fmt.Sprintf("%s:%s", embedding_model.Model, embedding_model.Provider.Provider))
	}
	data := []*models.DatasetDetailFields{}
	for _, v := range datasets {
		data = append(data, models.NewDatasetDetailFields(v))
	}

	for idx, item := range data {
		// convert embedding_model_provider to plugin standard format
		if item.IndexingTechnique == "high_quality" && item.EmbeddingModelProvider != "" {
			item.EmbeddingModelProvider = pluginentities.NewModelProviderID(item.EmbeddingModelProvider, false).String()
			item_model := fmt.Sprintf("%s:%s", item.EmbeddingModel, item.EmbeddingModelProvider)
			if slices.Contains(model_names, item_model) {
				item.EmbeddingAvailable = true
			} else {
				item.EmbeddingAvailable = false
			}
		} else {
			item.EmbeddingAvailable = true
		}
		if item.Permission == "partial_members" {
			item.PartialMemberList = services.ServiceGroupApp.Dataset.GetDatasetPartialMemberList(item.ID)
		} else {
			item.PartialMemberList = []string{}
		}
		data[idx] = item
	}
	bindata, _ := json.Marshal(data)
	out.DatasetsStr = string(bindata)
	out.HasMore = len(datasets) == int(in.Limit)
	out.Limit = in.Limit
	out.Page = in.Page
	out.Total = total
	return
}
func (s *DatasetsService) ExternalKnowledgeApiList(ctx context.Context, in *pbapi.ExternalKnowledgeApiListRequest) (out *pbapi.ExternalKnowledgeApiListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] datasets.ExternalKnowledgeApiList call:%#v", p.Addr.String(), in)
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Limit <= 0 {
		in.Limit = 20
	}

	out = &pbapi.ExternalKnowledgeApiListReply{Limit: in.Limit, Page: in.Page}
	if in.CurrentTenantId == "" {
		mlog.Errorf("CurrentTenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("CurrentTenantId not provide")
		return
	}
	external_knowledge_apis, total := services.ServiceGroupApp.Dataset.GetExternalKnowledgeAPIs(in.Page, in.Limit, in.CurrentTenantId, in.Keyword)
	datas := []map[string]any{}
	for _, item := range external_knowledge_apis {
		datas = append(datas, item.ToDict())
	}
	bindata, _ := json.Marshal(datas)
	out.DatasStr = string(bindata)
	out.HasMore = len(external_knowledge_apis) == int(in.Limit)
	out.Total = total
	return
}

// go build -o app.so -buildmode=plugin main.go
func (s *DatasetsService) Init(args ...any) error {
	mlog.Info("datasets init......")
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

func (s *DatasetsService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *DatasetsService) Destroy() {

}

func (s *DatasetsService) UserData() any {
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

	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Datasets})

	err = mrun.Run(&Datasets)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
