package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/cache"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/main/link/router"
	"mlib.com/gofy/server/proto/pbapi"

	"mlib.com/mlog"
)

// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	conn := cluster.Instance().GetRpcClientByModule("bot_workflow")
// 	mlog.Infof("-----")
// 	if conn != nil {
// 		mlog.Infof("-----")
// 		cli := pbapi.NewBotWorkflowClient(conn)
// 		mlog.Infof("-----")
// 		timeout_ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
// 		stream, err := cli.DraftWorkflowRun(timeout_ctx, &pbapi.DraftWorkflowRunRequest{Name: "ranwei"})
// 		mlog.Infof("-----")
// 		if err != nil {
// 			mlog.Errorf("remote call DraftWorkflowRun failed:%v", err)
// 			fmt.Fprintf(w, "remote call DraftWorkflowRun failed:%v", err)
// 		} else {
// 			for {
// 				res, err := stream.Recv()
// 				if err != nil {
// 					if err == io.EOF {
// 						mlog.Infof("remote call DraftWorkflowRun end")
// 						fmt.Fprintf(w, "remote call DraftWorkflowRun end")
// 						break
// 					} else {
// 						mlog.Errorf("remote call DraftWorkflowRun stream read failed:%v", err)
// 						fmt.Fprintf(w, "remote call DraftWorkflowRun stream read failed:%v", err)
// 						break
// 					}
// 				} else {
// 					mlog.Infof("remote call DraftWorkflowRun stream return=%#v", res)
// 					fmt.Fprintf(w, "remote call DraftWorkflowRun stream return=%#v", res)
// 				}
// 			}
// 		}
// 	} else {
// 		mlog.Errorf("get rpc client failed")
// 	}
// })
// go http.ListenAndServe(":8000", nil)

type LinkService struct {
	pbapi.UnimplementedLinkServer
	httpserver *http.Server

	wg sync.WaitGroup
}

var (
	Link LinkService
)

func (s *LinkService) RunWindowsServer() {

}

func (s *LinkService) LinkStatus(ctx context.Context, in *pbapi.LinkStatusRequest) (*pbapi.LinkStatusReply, error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] link.LinkStatus call:%s", p.Addr.String(), in.Name)
	return &pbapi.LinkStatusReply{Message: "hello " + in.Name}, nil
}

// go build -o SeatStatus.so -buildmode=plugin main.go
func (s *LinkService) Init(args ...any) error {
	mlog.Info("link init......")
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

	r := router.InitRouters()
	r.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", viper.GetIntWithDefault("system.addr", 5001))
	s.httpserver = &http.Server{
		Addr:              address,
		Handler:           r,
		ReadHeaderTimeout: 60 * time.Second,
		WriteTimeout:      60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	mlog.Info("server run success on ", "address", address)

	fmt.Printf(`
	欢迎使用 gofy
	当前版本:v0.0.0
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
`, address)
	s.wg.Add(1)
	go func() {
		mlog.Error(s.httpserver.ListenAndServe().Error())
		s.wg.Done()
	}()
	return nil
}

func (s *LinkService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *LinkService) Destroy() {
	mlog.Info("link destroy")
	if s.httpserver != nil {
		s.httpserver.Close()
	}
}

func (s *LinkService) UserData() any {
	return nil
}
