package cluster

import (
	"context"
	"errors"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
)

// peerServer 服务发现
type peerServer struct {
	modulename                 string
	grpcClientConn             *grpc.ClientConn
	service_discovery_provider ServiceDiscoveryProvider
}

// newPeerServer  新建发现服务
// func newPeerServer(modulename string) *peerServer {
// 	if modulename == "" {
// 		mlog.Errorf("modulename must provided")
// 		return nil
// 	}
// 	d := &peerServer{
// 		modulename: modulename,
// 		// grpcClientConn: conn,
// 	}
// 	resolver.Register(d)
// 	var err error
// 	fmt.Println("------------zkRootKeyStr=", mgr.zkRootKeyStr)
// 	fmt.Println("------------servertype=", servertype)
// 	d.grpcClientConn, err = grpc.Dial(mgr.zkRootKeyStr+"://8.8.8.8/"+servertype, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
// 	if err != nil {
// 		mlog.Errorf("连接服务端失败: %s", err)
// 		return nil
// 	}

// 	return d
// }

func (ps *peerServer) Init(args ...any) error {
	if len(args) < 2 {
		mlog.Errorf("invalid arg")
		return errors.New("invalid arg")
	}
	if provide, ok := args[0].(ServiceDiscoveryProvider); !ok || provide == nil {
		mlog.Errorf("args[0] must be ServiceDiscoveryProvider")
		return errors.New("args[0] must be ServiceDiscoveryProvider")
	} else {
		if !provide.Available() {
			mlog.Errorf("service discovery provider isn't available")
			return errors.New("service discovery provider isn't available")
		}
		ps.service_discovery_provider = provide
		if modulename, ok := args[1].(string); !ok || modulename == "" {
			mlog.Errorf("args[1] must provide a modulename")
			return errors.New("args[1] must provide a modulename")
		} else {
			ps.modulename = modulename
			var err error
			target := ps.service_discovery_provider.GetTarget(ps.modulename)
			ps.grpcClientConn, err = grpc.NewClient(target, grpc.WithResolvers(ps.service_discovery_provider), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(timeout.UnaryClientInterceptor(time.Duration(confy.GetWithDefault[int]("cluster.client_connect_max_timeout", 30))*time.Second)), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
			if err != nil {
				mlog.Warningf("连接服务端失败: %s", err)
				return nil
			}
		}
	}
	return nil
}

func (ps *peerServer) RunOnce(ctx context.Context) error {
	if ps.grpcClientConn == nil {
		var err error
		target := ps.service_discovery_provider.GetTarget(ps.modulename)
		ps.grpcClientConn, err = grpc.NewClient(target, grpc.WithResolvers(ps.service_discovery_provider), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(timeout.UnaryClientInterceptor(time.Duration(confy.GetWithDefault[int]("cluster.client_connect_max_timeout", 30))*time.Second)), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
		if err != nil {
			mlog.Warningf("连接服务端失败: %s", err)
			return nil
		}
	}
	return nil
}

func (ps *peerServer) Destroy() {

}

func (ps *peerServer) UserData() any {
	return ps
}
