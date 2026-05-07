package cluster

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc/resolver"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
)

type zkServiceDiscoveryProvide struct {
	zkHostStr        string
	zkRootKeyStr     string
	zkClusterStr     string
	zkSessionTimeout int
	zkClient         IZooClient
}

func (provide *zkServiceDiscoveryProvide) Init(args ...interface{}) error {

	if !confy.InConfig("cluster") || !confy.InConfig("cluster.port") || confy.Get[int]("cluster.port") == 0 {
		mlog.Warning("no port define, means no need cluster, derict return")
		return nil
	}
	// 服务IP

	if !confy.InConfig("ZooKeeper") || confy.Get[string]("ZooKeeper.HOST") == "" {
		mlog.Errorf("%s don't contain zk host config", confy.ConfigFileUsed())
		return fmt.Errorf("%s don't contain zk host config", confy.ConfigFileUsed())
	}
	provide.zkHostStr = confy.Get[string]("ZooKeeper.HOST")
	mlog.Info("ZooKeeper HOST: ", provide.zkHostStr)

	// rootkey 根节点
	provide.zkRootKeyStr = confy.Get[string]("ZooKeeper.ROOTKEY")
	if provide.zkRootKeyStr == "" {
		mlog.Errorf("%s don't contain zk ROOTKEY config", confy.ConfigFileUsed())
		return fmt.Errorf("%s don't contain zk ROOTKEY config", confy.ConfigFileUsed())
	}
	mlog.Info("ZooKeeper ROOTKEY: ", provide.zkRootKeyStr)

	// 集群标识
	// 服务path： rootKey/serverType/Cluster/IP:port
	provide.zkClusterStr = confy.Get[string]("ZooKeeper.CLUSTER")
	mlog.Info("ZooKeeper ROOTKEY: ", provide.zkClusterStr)

	// 超时时间
	provide.zkSessionTimeout = confy.Get[int]("ZooKeeper.TIMEOUT")
	if provide.zkSessionTimeout > MAX_ZK_RECV_TIMEOUT || provide.zkSessionTimeout < MIN_ZK_RECV_TIMEOUT {
		provide.zkSessionTimeout = DEFAULT_RECV_TIMEOUT
	}
	mlog.Info("ZooKeeper timeout: ", provide.zkSessionTimeout)
	hosts := strings.Split(provide.zkHostStr, ",")
	var err error
	for iLoop := 0; iLoop < 5; iLoop++ {
		provide.zkClient, err = NewZooClient(hosts)
		if err != nil {
			mlog.Warning("connect zookeeper failed:", err)
			time.Sleep(100 * time.Microsecond)
			continue
		}
	}
	if provide.zkClient == nil {
		mlog.Errorf("connect zookeeper failed:%v", err)
		return fmt.Errorf("connect zookeeper failed:%v", err)
	}
	return nil
}

func (provide *zkServiceDiscoveryProvide) RunOnce(ctx context.Context) error {

	return nil
}

func (provide *zkServiceDiscoveryProvide) Destroy() {
	if provide.zkClient != nil {
		provide.zkClient.Stop()
	}
}

func (provide *zkServiceDiscoveryProvide) UserData() interface{} {
	return provide
}

func (provide *zkServiceDiscoveryProvide) RootKey() string {
	return confy.Get[string]("ZooKeeper.ROOTKEY")
}
func (provide *zkServiceDiscoveryProvide) Register(service_name, ip_port_addr, service_status string) error {
	err := provide.zkClient.RegisterDirect(fmt.Sprintf("/%s/%s", provide.RootKey(), service_name), ip_port_addr, service_status)
	if err != nil {
		mlog.Errorf("zookeeper RegisterDirect failed: %v", err)
		return fmt.Errorf("zookeeper RegisterDirect failed: %v", err)
	}
	return nil
}

func (provide *zkServiceDiscoveryProvide) Available() bool {
	return provide.zkClient != nil
}

func (provide *zkServiceDiscoveryProvide) GetTarget(model_name string) string {
	return fmt.Sprintf("zookeeper:///%s/%s", provide.RootKey(), model_name)
}

// Build 为给定目标创建一个新的`resolver`，当调用`grpc.Dial()`时执行
func (provide *zkServiceDiscoveryProvide) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	mlog.Infof("Build target:%#v", target)
	endpoint := target.URL.Path
	if endpoint == "" {
		endpoint = target.URL.Opaque
		if endpoint == "" {
			endpoint = filepath.Join(target.URL.Scheme, target.URL.Host)
		}
	}

	r := &zkResolver{
		cli:        provide.zkClient,
		cc:         cc,
		serverList: make(map[string]*server),
		endpoint:   endpoint,
	}
	err := r.watch()
	if err != nil {
		mlog.Errorf("resolver watch failed:%v", err)
		return nil, fmt.Errorf("resolver watch failed:%v", err)
	}

	return r, nil
}

func (provide *zkServiceDiscoveryProvide) Scheme() string {
	return "zookeeper"
}

type zkResolver struct {
	cc         resolver.ClientConn
	cli        IZooClient
	serverList map[string]*server //服务列表
	endpoint   string
}

func (r *zkResolver) watch() error {
	mlog.Debugf("------endpoint=%s", r.endpoint)
	resp, _, err := r.cli.GetNodesAndEntries(r.endpoint)
	if err != nil {
		mlog.Errorf("get nodes and entries failed:%v", err)
		return err
	}
	mlog.Infof("resp=%#v", resp)

	for k, v := range resp {
		info := &server{}
		err = decodeServerStatus(v, info)
		if err != nil {
			mlog.Errorf("decodeServerStatus failed:%v", err)
			return fmt.Errorf("decodeServerStatus failed:%v", err)
		}
		r.serverList[k] = info
	}
	r.cc.UpdateState(resolver.State{Addresses: r.getServices()})
	// ps.cc.NewAddress()
	//监视前缀，修改变更的server
	r.cli.AddInstancer(r.endpoint,
		func(ctx context.Context, parent string, new map[string]string) error {
			r.serverList = make(map[string]*server)
			for k, v := range new {
				info := &server{}
				err = decodeServerStatus(v, info)
				if err != nil {
					mlog.Errorf("decodeServerStatus failed:%v", err)
					continue
				}
				r.serverList[k] = info
			}
			r.cc.UpdateState(resolver.State{Addresses: r.getServices()})
			// ps.cc.NewAddress(ps.getServices())
			return nil
		}, func(ctx context.Context, err error) {

		}, func(ctx context.Context, parent, node string, data string) error {
			return nil
		})
	return nil
}
func (r *zkResolver) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0, len(r.serverList))

	for _, v := range r.serverList {
		addr := &resolver.Address{}
		addr.Addr = v.IpPort
		addrs = append(addrs, *addr)
	}
	return addrs
}
func (r *zkResolver) Close() {
	mlog.Infof("Close")
}
func (r *zkResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	mlog.Infof("ResolveNow")
}
