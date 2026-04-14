package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/mrun"
	"mlib.com/gofy/server/utils"

	"github.com/odysseythink/mlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetIP() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		mlog.Errorf("net.Interfaces() failed:%v", err)
		return "0.0.0.0"
	}
	for _, intf := range netInterfaces {
		// fmt.Println("------name=", intf.Name)
		// fmt.Printf("------.Flags=%d\n", intf.Flags&net.FlagLoopback)
		if (intf.Flags&net.FlagUp) != 0 && (intf.Flags&net.FlagLoopback) == 0 {
			addrs, _ := intf.Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						// fmt.Println("-----------",ipnet.IP.String())
						return ipnet.IP.To4().String()
					}
				}
			}
		}
	}
	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	mlog.Errorf("InterfaceAddrs failed:%v", err)
	// 	return ""
	// }
	// fmt.Println("")
	// for _, address := range addrs {
	// 	// 检查ip地址判断是否回环地址
	// 	if ipnet, ok := address.(*net.IPNet); ok {
	// 		fmt.Println("----", ipnet, "--IsLoopback=", ipnet.IP.IsLoopback(), "--IsMulticast=", ipnet.IP.IsMulticast(), "------IsInterfaceLocalMulticast=", ipnet.IP.IsInterfaceLocalMulticast(), "------IsLinkLocalMulticast=", ipnet.IP.IsLinkLocalMulticast(), "----IsLinkLocalUnicast=", ipnet.IP.IsLinkLocalUnicast(), "---IsGlobalUnicast", ipnet.IP.IsGlobalUnicast())
	// 	}
	// }
	// for _, address := range addrs {
	// 	// 检查ip地址判断是否回环地址
	// 	if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 		if ipnet.IP.To4() != nil {
	// 			// fmt.Println(ipnet.IP.String())
	// 			return ipnet.IP.To4().String()
	// 		}
	// 	}
	// }
	return "0.0.0.0"
}

type server struct {
	IpPort           string `json:"ip_port"`              // ip和端口，用于标识唯一的服务
	Ip               string `json:"ip"`                   // 服务IP
	Port             uint16 `json:"port"`                 // 服务端口
	ModuleName       string `json:"server_type"`          // 服务类型,见 SERVER_TYPE
	MsgType          int    `json:"msg_type"`             // 服务处理消息类型,从加载的工作模块获取
	ReqServerTypeStr string `json:"req_server_type"`      // 请求服务消息类型列表
	GreyServerFlag   uint16 `json:"grey_server_flag"`     // 灰度升级服务标志
	RecvNum          int    `json:"req_total_num"`        // 接收请求数
	RejectNum        int    `json:"req_total_reject_num"` // 拒绝请求数
	FailedNum        int    `json:"req_total_failed_num"` // 失败请求数
	LastRecvNum      int    `json:"req_num"`              // 最新接收请求数
	LastRejectNum    int    `json:"req_reject_num"`       // 最新拒绝请求数
	LastFailedNum    int    `json:"req_failed_num"`       // 最新失败请求数
	AvgTime          uint   `json:"req_avg_time"`         // 最近N个消息的平均处理时间
	HostName         string `json:"hostname"`             //本地服务器主机名
	// VirtualNodeNumer     int    `json:"VirtualNodeNumer"`     //虚拟节点数
	// Weight               int    `json:"Weight"` //权重比例
	LastRefreshTime int64          `json:"-"` // 服务状态刷新时间 毫秒
	IsConnected     bool           `json:"-"` // 是否连接
	ConnectId       int            `json:"-"` // 连接ID
	ConnectStatus   bool           `json:"-"` // true:成功，false:失败
	ReqModulesStr   string         `json:"-"`
	rpcService      any            `json:"-"`
	grpcServer      *grpc.Server   `json:"-"`
	lis             net.Listener   `json:"-"`
	wg              sync.WaitGroup `json:"-"`
	// listmgr              *ServerListMgr `json:"-"`
	service_discovery_provider ServiceDiscoveryProvider `json:"-"`
	peerServers                mrun.ModuleMgr           `json:"-"`
}

func decodeServerStatus(strJsonData string, info *server) error {
	if strJsonData == "" {
		mlog.Errorf("empty input json string")
		return fmt.Errorf("empty input json string")
	}
	if info == nil {
		mlog.Errorf("info must provied")
		return fmt.Errorf("info must provied")
	}

	type TmpSt struct {
		ServerInfo *server `json:"server_info"`
		Interval   int     `json:"interval"`
	}
	tmp := &TmpSt{ServerInfo: info}
	err := json.Unmarshal([]byte(strJsonData), tmp)
	if err != nil {
		mlog.Errorf("unmarshal json failed:%v", err)
		return fmt.Errorf("unmarshal json failed:%v", err)
	}
	return nil
}

func encodeServerStatus(info *server) string {
	if info == nil {
		mlog.Error("encodeServerStatus failed:empty input ServerInfo struct")
		return ""
	}
	type TmpSt struct {
		ServerInfo *server `json:"server_info"`
	}
	tmp := TmpSt{ServerInfo: info}
	data, err := json.Marshal(tmp)
	if err != nil {
		mlog.Errorf("encodeServerStatus marshal json failed:%v", err)
		return ""
	}
	return string(data)
}

func (s *server) GetRpcClientByModule(modulename string) *grpc.ClientConn {
	mods := s.peerServers.GetModulesByAlias(modulename)
	if len(mods) == 0 {
		mlog.Warningf("module(%s) not exist in module list", modulename)
		return nil
	}
	return mods[0].(*peerServer).grpcClientConn
}

func (s *server) Init(args ...any) error {
	err := s.argsInit(args...)
	if err != nil {
		return err
	}

	//本地服务器主机名
	s.HostName = GetIP()
	if s.HostName == "" {
		mlog.Errorf("null hostname")
		return fmt.Errorf("null hostname")
	}
	mlog.Infof("Server HostName %s", s.HostName)

	// //虚拟节点数
	// s.serverInfo.VirtualNodeNumer = cfginstance.GetInt("ServerCom.VirtualNodeNum")
	// if s.serverInfo.VirtualNodeNumer == 0 {
	// 	s.serverInfo.VirtualNodeNumer = 128
	// }
	// mlog.Infof("VirtualNodeNum %d", s.serverInfo.VirtualNodeNumer)

	// //权重
	// s.serverInfo.Weight = cfginstance.GetInt("ServerCom.Weight")
	// if s.serverInfo.Weight == 0 {
	// 	s.serverInfo.Weight = 1
	// }
	// mlog.Infof("Weight %d", s.serverInfo.Weight)

	s.Ip = confy.Get[string]("ServerCom.IP")
	if s.Ip == "" {
		s.Ip = s.HostName
	}
	mlog.Infof("IP %s", s.Ip)

	s.Port = uint16(confy.Get[int]("cluster.port"))
	if s.Port == 0 {
		mlog.Warning("no port define, means no need cluster, derict return")
		return nil
	}
	mlog.Infof("Port %d", s.Port)
	// 服务IP和PORT
	s.IpPort = fmt.Sprintf("%s:%d", s.Ip, s.Port)
	mlog.Infof("Server IpPort %s", s.IpPort)

	// 请求服务类型列表
	s.ReqModulesStr = confy.Get[string]("cluster.req_modules")
	if s.ReqModulesStr == "" {
		mlog.Errorf("cluster.req_modules is not exist")
		return fmt.Errorf("cluster.req_modules is not exist")
	}
	mlog.Infof("ReqModulesStr: %s", s.ReqModulesStr)

	// 监听本地端口
	s.lis, err = net.Listen("tcp", s.IpPort)
	if err != nil {
		mlog.Errorf("Listen(%s) failed: %v", s.IpPort, err)
		return fmt.Errorf("Listen(%s) failed: %v", s.IpPort, err)
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// reflection.Register(s)
		mlog.Debugf("start grpc serve")
		err = s.grpcServer.Serve(s.lis)
		mlog.Debugf("grpc serve return")
		if err != nil {
			mlog.Errorf("开启服务失败: %s", err)
			return
		}
	}()
	err = s.service_discovery_provider.Register(s.ModuleName, s.IpPort, encodeServerStatus(s))
	if err != nil {
		mlog.Errorf("service discovery provider Register failed: %v", err)
		s.Destroy()
		return fmt.Errorf("service discovery provider Register failed: %v", err)
	}
	for v := range strings.SplitSeq(s.ReqModulesStr, ",") {
		if v == s.ModuleName {
			continue
		}

		if s.peerServers.GetModulesByAlias(v) == nil { // 不存在
			ps := &peerServer{}
			err = s.peerServers.Register(ps, []mrun.ModuleMgrOption{mrun.NewModuleAliasOption(v)}, s.service_discovery_provider, v)
			if err != nil {
				mlog.Errorf("new PeerServer(%s) failed:%v", v, err)
				s.Destroy()
				return fmt.Errorf("new PeerServer(%s) failed:%v", v, err)
			}
		}
	}
	err = s.peerServers.Init()
	if err != nil {
		mlog.Errorf("peerServers init failed:%v", err)
		s.Destroy()
		return fmt.Errorf("peerServers init failed:%v", err)
	}

	return nil
}

func (s *server) RunOnce(ctx context.Context) error {
	return nil
}

func (s *server) Destroy() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
		s.grpcServer = nil
		s.wg.Wait()
	}
	if s.lis != nil {
		s.lis.Close()
		s.lis = nil
	}
}

func (s *server) UserData() any {
	return s
}

func (s *server) argsInit(args ...any) error {
	if len(args) < 2 {
		mlog.Errorf("invalid arg")
		return errors.New("invalid arg")
	}
	if provide, ok := args[0].(ServiceDiscoveryProvider); !ok || provide == nil {
		mlog.Errorf("args[0] must be ServiceDiscoveryProvider")
		return errors.New("args[0] must be ServiceDiscoveryProvider")
	} else {

		port := uint16(confy.Get[int]("cluster.port"))
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
		s.ModuleName = confy.Get[string]("ServerCom.ModuleName")
		if s.ModuleName == "" {
			mlog.Errorf("ModuleName is not exist")
			return fmt.Errorf("ModuleName is not exist")
		}
		mlog.Infof("ModuleName %s", s.ModuleName)
		s.rpcService = args[1]
		if s.rpcService == nil {
			mlog.Errorf("rpc service implement instance must be provided")
			return errors.New("rpc service implement instance must be provided")
		}

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
		err := s.grpcServer.RegisterServiceWithoutDesc(s.rpcService, "pbapi."+s.ModuleName)
		if err != nil {
			mlog.Errorf("RegisterServiceWithoutDesc failed: %v", err)
			s.grpcServer = nil
			return fmt.Errorf("RegisterServiceWithoutDesc failed: %v", err)
		}
		return nil
	}
}
