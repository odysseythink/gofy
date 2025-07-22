//go:build !windows
// +build !windows

package cluster

import (
	"errors"
	"fmt"
	"plugin"

	"mlib.com/gofy/server/utils"
	"mlib.com/mrun"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mlib.com/mlog"
)

func (s *server) argsInit(args ...any) error {
	if len(args) < 1 {
		mlog.Errorf("invalid arg")
		return errors.New("invalid arg")
	}
	if provide, ok := args[0].(ServiceDiscoveryProvider); !ok || provide == nil {
		mlog.Errorf("args[0] must be ServiceDiscoveryProvider")
		return errors.New("args[0] must be ServiceDiscoveryProvider")
	} else {
		port := uint16(viper.GetInt("cluster.port"))
		if port == 0 {
			mlog.Warning("no port define, means no need cluster, derict return")
			return nil
		}

		if !provide.Available() {
			mlog.Errorf("service discovery provider isn't available")
			return errors.New("service discovery provider isn't available")
		}
		s.service_discovery_provider = provide

		// 服务类型
		s.ModuleName = viper.GetString("ServerCom.ModuleName")
		if s.ModuleName == "" {
			mlog.Errorf("ModuleName is not exist")
			return fmt.Errorf("ModuleName is not exist")
		}
		mlog.Infof("ModuleName %s", s.ModuleName)
		plug, err := plugin.Open(s.ModuleName + ".so")
		if err != nil {
			mlog.Errorf("load Module(%s) failed:%v", s.ModuleName, err)
			return fmt.Errorf("load Module(%s) failed:%v", s.ModuleName, err)
		}

		symbol, err := plug.Lookup(GetExportModulename(s.ModuleName))
		if err != nil {
			mlog.Errorf("Module(%s) Lookup WorkModule failed:%v", s.ModuleName, err)
			return fmt.Errorf("Module(%s) Lookup WorkModule failed:%v", s.ModuleName, err)
		}
		if _, ok := symbol.(mrun.IModule); !ok {
			mlog.Errorf("Module(%s) must implement mrun.IModule interface", s.ModuleName)
			return fmt.Errorf("Module(%s) must implement mrun.IModule interface", s.ModuleName)
		}
		s.rpcService = symbol
		// 创建gRPC服务器
		grpcPanicRecoveryHandler := func(p any) (err error) {
			mlog.Errorf("recovered from panic:%v\n%s", p, utils.GetCurrentGoroutineStack())
			return status.Errorf(codes.Internal, "%v", p)
		}
		s.grpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			),
			grpc.ChainStreamInterceptor(
				recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			),
		)
		err = s.grpcServer.RegisterServiceWithoutDesc(s.rpcService, "pbapi."+s.ModuleName)
		if err != nil {
			mlog.Errorf("RegisterServiceWithoutDesc failed: %v", err)
			s.grpcServer = nil
			return fmt.Errorf("RegisterServiceWithoutDesc failed: %v", err)
		}
		return nil
	}
}
