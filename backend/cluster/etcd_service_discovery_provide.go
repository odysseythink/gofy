package cluster

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

type etcdServiceDiscoveryProvide struct {
	cli          *clientv3.Client
	lease        *clientv3.LeaseGrantResponse
	endpoint_mgr endpoints.Manager
	leaseTimer   *time.Timer
}

func (provide *etcdServiceDiscoveryProvide) Init(args ...any) error {
	cli, err := clientv3.New(clientv3.Config{
		Username:    confy.Get[string]("etcd.user_name"),
		Password:    confy.Get[string]("etcd.password"),
		Endpoints:   strings.Split(confy.Get[string]("etcd.endpoints"), ","),
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		log.Printf("[E]create clientv3 client failed: %v\n", err)
		return fmt.Errorf("create clientv3 client failed: %v", err)
	}

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	ttl := confy.GetWithDefault[int64]("etcd.ttl", 10)
	lease, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		cli.Close()
		log.Printf("[E]etcd client grant failed: %v\n", err)
		return fmt.Errorf("etcd client grant failed: %v", err)
	}
	provide.cli = cli
	provide.lease = lease
	return nil
}

func (provide *etcdServiceDiscoveryProvide) RunOnce(ctx context.Context) error {
	if provide.leaseTimer == nil {
		provide.leaseTimer = time.NewTimer(5 * time.Second)
	}
	select {
	case <-provide.leaseTimer.C:
		// 续约操作
		provide.cli.KeepAliveOnce(ctx, provide.lease.ID)
		// mlog.Debugf("keep alive resp: %+v", resp)
		provide.leaseTimer.Reset(5 * time.Second)
	default:
		return nil
	}
	return nil
}

func (provide *etcdServiceDiscoveryProvide) Destroy() {
	if provide.cli != nil {
		provide.cli.Close()
	}
}

func (provide *etcdServiceDiscoveryProvide) UserData() any {
	return provide
}

func (provide *etcdServiceDiscoveryProvide) RootKey() string {
	return confy.Get[string]("etcd.ROOTKEY")
}
func (provide *etcdServiceDiscoveryProvide) Register(service_name, ip_port_addr, service_status string) error {
	endpoint_mgr, err := endpoints.NewManager(provide.cli, fmt.Sprintf("%s/%s", provide.RootKey(), service_name))
	if err != nil {
		mlog.Errorf("create endpoint manager failed: %v", err)
		return fmt.Errorf("create endpoint manager failed: %v", err)
	}
	provide.endpoint_mgr = endpoint_mgr
	err = provide.endpoint_mgr.AddEndpoint(context.Background(), fmt.Sprintf("%s/%s/%s", provide.RootKey(), service_name, ip_port_addr), endpoints.Endpoint{Addr: ip_port_addr, Metadata: service_status}, clientv3.WithLease(provide.lease.ID))
	if err != nil {
		mlog.Errorf("etcd add endpoint failed: %v", err)
		return fmt.Errorf("etcd add endpoint failed: %v", err)
	}

	return nil
}
func (provide *etcdServiceDiscoveryProvide) Available() bool {
	return provide.cli != nil && provide.lease != nil
}

func (provide *etcdServiceDiscoveryProvide) GetTarget(model_name string) string {
	return fmt.Sprintf("etcd:///%s/%s", provide.RootKey(), model_name)
}

func (provide *etcdServiceDiscoveryProvide) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// Refer to https://github.com/grpc/grpc-go/blob/16d3df80f029f57cff5458f1d6da6aedbc23545d/clientconn.go#L1587-L1611
	// fmt.Println("------", target)
	endpoint := target.URL.Path
	if endpoint == "" {
		endpoint = target.URL.Opaque
		if endpoint == "" {
			endpoint = filepath.Join(target.URL.Scheme, target.URL.Host)
		}
	}
	endpoint = strings.TrimPrefix(endpoint, "/")
	r := &etcdResolver{
		cli:    provide.cli,
		target: endpoint,
		cc:     cc,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())

	em, err := endpoints.NewManager(r.cli, r.target)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "resolver: failed to new endpoint manager: %s", err)
	}
	r.wch, err = em.NewWatchChannel(r.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resolver: failed to new watch channer: %s", err)
	}

	r.wg.Add(1)
	go r.watch()
	return r, nil
}

func (provide *etcdServiceDiscoveryProvide) Scheme() string {
	return "etcd"
}

type etcdResolver struct {
	cli    *clientv3.Client
	target string
	cc     resolver.ClientConn
	wch    endpoints.WatchChannel
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (r *etcdResolver) watch() {
	defer r.wg.Done()

	allUps := make(map[string]*endpoints.Update)
	for {
		select {
		case <-r.ctx.Done():
			return
		case ups, ok := <-r.wch:
			if !ok {
				return
			}

			for _, up := range ups {
				switch up.Op {
				case endpoints.Add:
					allUps[up.Key] = up
				case endpoints.Delete:
					delete(allUps, up.Key)
				}
			}

			eps := convertToGRPCEndpoint(allUps)
			r.cc.UpdateState(resolver.State{Endpoints: eps})
		}
	}
}

func convertToGRPCEndpoint(ups map[string]*endpoints.Update) []resolver.Endpoint {
	var eps []resolver.Endpoint
	for _, up := range ups {
		fmt.Printf("------up=%#v\n", up)
		ep := resolver.Endpoint{
			Addresses: []resolver.Address{
				{
					Addr:     up.Endpoint.Addr,
					Metadata: up.Endpoint.Metadata,
				},
			},
		}
		eps = append(eps, ep)
	}
	return eps
}

// ResolveNow is a no-op here.
// It's just a hint, resolver can ignore this if it's not necessary.
func (r *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *etcdResolver) Close() {
	r.cancel()
	r.wg.Wait()
}
