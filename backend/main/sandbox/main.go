package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/gofy/backend/cluster"
	"github.com/odysseythink/gofy/backend/main/sandbox/global"
	"github.com/odysseythink/gofy/backend/main/sandbox/runner/python"
	runnertypes "github.com/odysseythink/gofy/backend/main/sandbox/runner/types"
	"github.com/odysseythink/gofy/backend/main/sandbox/services"
	pbexceptions "github.com/odysseythink/gofy/backend/proto/exceptions"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
	"github.com/odysseythink/mrun"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/peer"
)

func (s *SandboxService) sandbox_user_init() error {
	// create sandbox user
	user := confy.GetWithDefault[string]("sandbox_user", "sandbox")
	uid := confy.GetWithDefault[int]("sandbox_user_uid", 65537)

	// check if user exists
	_, err := exec.Command("id", user).Output()
	if err != nil {
		// create user
		output, err := exec.Command("bash", "-c", "useradd -u "+strconv.Itoa(uid)+" "+user).Output()
		if err != nil {
			mlog.Errorf("failed to create user: %v, %v", err, string(output))
			return fmt.Errorf("failed to create user: %v, %v", err, string(output))
		}
	}

	// get gid of sandbox user and setgid
	gid, err := exec.Command("id", "-g", user).Output()
	if err != nil {
		mlog.Errorf("failed to get gid of user: %v", err)
		return fmt.Errorf("failed to get gid of user: %v", err)
	}

	global.SANDBOX_GROUP_ID, err = strconv.Atoi(strings.TrimSpace(string(gid)))
	if err != nil {
		mlog.Errorf("failed to convert gid: %v", err)
		return fmt.Errorf("failed to convert gid: %v", err)
	}
	return nil
}

type SandboxService struct {
	pbapi.UnimplementedSandboxServer
}

var (
	Sandbox SandboxService
	limiter *rate.Limiter
)

func (s *SandboxService) Health(ctx context.Context, in *pbapi.HealthRequest) (out *pbapi.HealthReply, err error) {
	if !limiter.Allow() {
		return nil, errors.New("too many requests")
	}
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] sandbox.Health call:%#v", p.Addr.String(), in)

	out = &pbapi.HealthReply{Msg: "ok"}

	return
}
func (s *SandboxService) Run(ctx context.Context, in *pbapi.RunRequest) (out *pbapi.RunReply, err error) {
	if !limiter.Allow() {
		return nil, errors.New("too many requests")
	}
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] sandbox.Run call:%#v", p.Addr.String(), in)

	out = &pbapi.RunReply{}
	switch in.Language {
	case "python3":
		rsp := services.ServiceGroupApp.Python.RunPython3Code(in.Code, in.Preload, &runnertypes.RunnerOptions{
			EnableNetwork: in.EnableNetwork,
		})
		out.Code = int32(rsp.Code)
		out.Message = rsp.Message
		if rsp.Data != nil {
			out.Error = rsp.Data.(*services.RunCodeResponse).Stderr
			out.Stdout = rsp.Data.(*services.RunCodeResponse).Stdout
		}
	default:
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Message: "unsupported language",
		}
	}
	return
}
func (s *SandboxService) GetDependencies(ctx context.Context, in *pbapi.DependenciesRequest) (out *pbapi.DependenciesReply, err error) {
	if !limiter.Allow() {
		return nil, errors.New("too many requests")
	}
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] sandbox.GetDependencies call:%#v", p.Addr.String(), in)

	switch in.Language {
	case "python3":
		out = services.ServiceGroupApp.Python.ListPython3Dependencies()
	default:
		out = &pbapi.DependenciesReply{}
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Message: "unsupported language",
		}
	}
	return
}
func (s *SandboxService) UpdateDependencies(ctx context.Context, in *pbapi.DependenciesRequest) (out *pbapi.DependenciesReply, err error) {
	if !limiter.Allow() {
		return nil, errors.New("too many requests")
	}
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] sandbox.UpdateDependencies call:%#v", p.Addr.String(), in)

	switch in.Language {
	case "python3":
		out = services.ServiceGroupApp.Python.UpdateDependencies()
	default:
		out = &pbapi.DependenciesReply{}
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Message: "unsupported language",
		}
	}
	return
}
func (s *SandboxService) RefreshDependencies(ctx context.Context, in *pbapi.DependenciesRequest) (out *pbapi.DependenciesReply, err error) {
	if !limiter.Allow() {
		return nil, errors.New("too many requests")
	}
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] sandbox.RefreshDependencies call:%#v", p.Addr.String(), in)

	switch in.Language {
	case "python3":
		out = services.ServiceGroupApp.Python.RefreshPython3Dependencies()
	default:
		out = &pbapi.DependenciesReply{}
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Message: "unsupported language",
		}
	}
	return
}

// go build -o app.so -buildmode=sandbox main.go
func (s *SandboxService) Init(args ...any) error {
	mlog.Info("sandbox init......")

	err := s.sandbox_user_init()
	if err != nil {
		mlog.Errorf("user init failed:%v", err)
		return fmt.Errorf("user init failed:%v", err)
	}
	err = global.SetupRunnerDependencies()
	if err != nil {
		mlog.Errorf("failed to setup runner dependencies: %v", err)
		return fmt.Errorf("failed to setup runner dependencies: %v", err)
	}
	err = python.Setup()
	if err != nil {
		mlog.Errorf("failed to setup python: %v", err)
		return fmt.Errorf("failed to setup python: %v", err)
	}
	err = cache.Instance().Init()
	if err != nil {
		mlog.Errorf("cahe init failed:%v", err)
		return fmt.Errorf("cahe init failed:%v", err)
	}

	limiter = rate.NewLimiter(rate.Limit(confy.GetWithDefault[int]("max_requests", 10000)), confy.GetWithDefault[int]("max_requests", 10000)) // 每秒最多10000个请求，桶大小为10000

	return nil
}

func (s *SandboxService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *SandboxService) Destroy() {

}

func (s *SandboxService) UserData() any {
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

	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Sandbox})

	err = mrun.Run(&Sandbox)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
