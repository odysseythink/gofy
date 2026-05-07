package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/gofy/backend/cluster"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/main/link/router"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mrun"
	"google.golang.org/grpc/peer"

	"github.com/odysseythink/mlog"
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

	// Serve frontend static files with SPA fallback
	frontendAssets, err := FrontendAssets()
	if err != nil {
		mlog.Errorf("failed to load frontend assets: %v", err)
	} else {
		rfs := frontendAssets.(fs.ReadFileFS)
		r.NoRoute(func(c *gin.Context) {
			reqPath := c.Request.URL.Path

			// API paths must never fall through to SPA — return JSON 404 so the
			// frontend can distinguish missing endpoints from unauthenticated SPA loads.
			if strings.HasPrefix(reqPath, "/console/api/") || strings.HasPrefix(reqPath, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    "not_found",
					"message": fmt.Sprintf("endpoint %s not registered", reqPath),
				})
				return
			}

			path := strings.TrimPrefix(reqPath, "/")

			// Try to serve the exact file from embedded assets.
			if path != "" {
				if data, fErr := rfs.ReadFile(path); fErr == nil {
					ctype := mime.TypeByExtension(filepath.Ext(path))
					if ctype == "" {
						ctype = http.DetectContentType(data)
					}
					c.Data(http.StatusOK, ctype, data)
					return
				}
			}

			// SPA fallback: write index.html directly (avoid http.ServeFile's
			// auto-redirect on "/index.html" which would loop on NoRoute).
			data, fErr := rfs.ReadFile("index.html")
			if fErr != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		})
	}

	address := fmt.Sprintf(":%d", confy.GetWithDefault[int]("system.addr", 5001))
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
	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Link})

	err = mrun.Run(&Link)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
