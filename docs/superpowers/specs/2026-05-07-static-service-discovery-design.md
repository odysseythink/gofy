# 静态服务发现方案设计

**日期**: 2026-05-07
**状态**: 设计已确认，待实现
**作者**: Kimi Code CLI

---

## 1. 背景与问题

当前 `backend/cluster` 模块通过 `ServiceDiscoveryProvider` 接口支持两种服务发现后端：

- **ZooKeeper** —— 多节点集群部署时使用
- **etcd** —— 多节点集群部署时使用

当单节点部署时（本地开发、小规模测试），仍必须启动 ZooKeeper 或 etcd 集群才能正常运行。这带来了不必要的运维负担，且每个服务需要配置大量连接参数（HOST、ROOTKEY、认证等）。

**核心诉求**：单节点部署时零外部依赖启动，服务间仍能正常通信；多节点部署时继续使用现有 zk/etcd 机制。

---

## 2. 目标

1. 新增一种不依赖外部协调服务的静态服务发现 Provider
2. 单节点时，服务启动无需 etcd/zookeeper 即可工作
3. 服务间 RPC 配置极简化：知道模块名即可自动定位对端地址
4. 多节点场景无缝切换回 zk/etcd，改动仅限于配置
5. 不破坏现有 etcd/zookeeper Provider 的任何行为

---

## 3. 方案对比

| 方案 | 原理 | 优点 | 缺点 | 结论 |
|------|------|------|------|------|
| A. 同进程方法调用 | 所有服务编译为单一进程 | 无网络开销 | 与当前"独立可执行文件"架构冲突 | ❌ 不适用 |
| B. 约定式端口 | 代码内置模块名→端口映射 | 最简单、零外部依赖、全平台兼容 | 新增模块需更新代码映射表 | ✅ 选用 |
| C. Unix Domain Socket | 同机通过 `.sock` 文件通信 | 比 TCP 更快、无端口冲突 | Windows 兼容性差 | ❌ 不选用 |
| D. 文件系统注册表 | 各服务将地址写入共享临时文件 | 可动态增删 | 需处理文件锁和变更通知，比 B 复杂 | ❌ 不选用 |

**选用方案 B（约定式端口）**。理由：与现有 gRPC + TCP 架构完全兼容；实现最简单；Windows/Linux/macOS 全平台支持；单节点与多节点切换成本最低（改一行配置即可）。

---

## 4. 详细设计

### 4.1 架构关系

```
cluster.go Init()
├── 读取 cluster.service_discovery_provide
│   ├── "zookeeper" → zkServiceDiscoveryProvide
│   ├── "etcd"      → etcdServiceDiscoveryProvide
│   └── "static"    → staticServiceDiscoveryProvide  ← 新增
└── 注入到 server / peerServer
```

### 4.2 新增文件

```
backend/cluster/
├── static_service_discovery_provide.go   # 静态 Provider 实现
└── static_service_discovery_provide_test.go  # 单元测试
```

### 4.3 端口约定机制

#### 4.3.1 默认端口映射（代码内置）

在 `static_service_discovery_provide.go` 中内置默认映射表，覆盖项目已有模块：

```go
var defaultModulePorts = map[string]int{
    "api":    19001,
    "worker": 19002,
    "model":  19003,
    // 后续新增模块时在此添加一行
}
```

新增模块只需在代码里增加一行端口约定，**不需要修改每个服务的配置文件**。

#### 4.3.2 配置项

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `cluster.static_host` | string | `127.0.0.1` | 单节点时各服务所在主机地址 |
| `cluster.static_ports` | map[string]int | `{}` | 可选，覆盖默认端口映射 |

配置覆盖示例：

```yaml
cluster:
  static_host: "127.0.0.1"
  static_ports:
    api: 20001  # 覆盖默认 19001
```

#### 4.3.3 地址解析逻辑

```go
func (p *staticServiceDiscoveryProvide) GetTarget(model_name string) string {
    return fmt.Sprintf("static:///%s", model_name)
}

// resolver.Build() 中：
moduleName := strings.TrimPrefix(target.URL.Path, "/")
port, ok := p.portMap[moduleName]
if !ok {
    return nil, fmt.Errorf("static: unknown module %s, add it to defaultModulePorts", moduleName)
}
addr := fmt.Sprintf("%s:%d", p.host, port)
cc.UpdateState(resolver.State{
    Addresses: []resolver.Address{{Addr: addr}},
})
```

### 4.4 静态 Resolver 实现

```go
type staticServiceDiscoveryProvide struct {
    host    string
    portMap map[string]int
}

func (p *staticServiceDiscoveryProvide) Init(args ...any) error {
    p.host = confy.GetWithDefault[string]("cluster.static_host", "127.0.0.1")
    
    // 深拷贝默认映射，再用配置覆盖
    p.portMap = make(map[string]int, len(defaultModulePorts))
    for k, v := range defaultModulePorts {
        p.portMap[k] = v
    }
    if confy.InConfig("cluster.static_ports") {
        for k, v := range confy.Get[map[string]int]("cluster.static_ports") {
            p.portMap[k] = v
        }
    }
    return nil
}

func (p *staticServiceDiscoveryProvide) Register(service_name, ip_port_addr, service_status string) error {
    mlog.Infof("[static] register service=%s addr=%s (no-op)", service_name, ip_port_addr)
    return nil
}

func (p *staticServiceDiscoveryProvide) Available() bool {
    return true // 静态 Provider 始终可用
}

func (p *staticServiceDiscoveryProvide) RootKey() string {
    return "" // 静态模式下不使用根键
}

func (p *staticServiceDiscoveryProvide) Build(
    target resolver.Target,
    cc resolver.ClientConn,
    opts resolver.BuildOptions,
) (resolver.Resolver, error) {
    moduleName := strings.TrimPrefix(target.URL.Path, "/")
    port, ok := p.portMap[moduleName]
    if !ok {
        return nil, fmt.Errorf("static: unknown module %s, add it to defaultModulePorts", moduleName)
    }
    addr := fmt.Sprintf("%s:%d", p.host, port)
    cc.UpdateState(resolver.State{
        Addresses: []resolver.Address{{Addr: addr}},
    })
    return &staticResolver{}, nil
}

func (p *staticServiceDiscoveryProvide) Scheme() string {
    return "static"
}

type staticResolver struct{}

func (r *staticResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *staticResolver) Close()                                {}
```

**关键特性**：
- `Register()` 为无操作，仅打印日志，**不访问任何外部系统**
- `Available()` 始终返回 `true`
- 无需 `RunOnce()` 续约/心跳逻辑
- 静态地址不变，resolver 解析一次后不再更新

### 4.5 cluster.go 集成

在 `cluster.go` 的 `Init()` 方法中增加 `"static"` 分支：

```go
func (c *Cluster) Init(args ...any) error {
    // ...
    provider_name := confy.Get[string]("cluster.service_discovery_provide")
    switch provider_name {
    case "zookeeper":
        // 现有逻辑
    case "etcd":
        // 现有逻辑
    case "static":
        static_provider := &staticServiceDiscoveryProvide{}
        c.subsMgr.Register(static_provider, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)})
        c.service_discovery_provider = static_provider
    default:
        mlog.Errorf("unsupported provider=%s", provider_name)
        return fmt.Errorf("unsupported provider=%s", provider_name)
    }
    // ... 后续逻辑不变
}
```

### 4.6 配置示例

#### 单节点部署（api-service）

```yaml
ServerCom:
  ModuleName: api

cluster:
  service_discovery_provide: static
  port: 19001
  req_modules: worker,model
```

#### 单节点部署（worker-service）

```yaml
ServerCom:
  ModuleName: worker

cluster:
  service_discovery_provide: static
  port: 19002
  req_modules: api,model
```

#### 多节点部署（切换回 ZooKeeper）

```yaml
cluster:
  service_discovery_provide: zookeeper
  port: 19001
  req_modules: worker,model

ZooKeeper:
  HOST: "zk1:2181,zk2:2181"
  ROOTKEY: "gofy"
  CLUSTER: "prod"
  TIMEOUT: 5000
```

**对比**：单节点时不再需要配置 `ZooKeeper.*` 或 `etcd.*` 的任何连接参数。

---

## 5. 边界情况处理

| 场景 | 处理方式 |
|------|----------|
| 模块名不在 `defaultModulePorts` 中 | `Build()` 返回显式错误，日志提示 `"unknown module xxx, add it to defaultModulePorts"` |
| `cluster.port` 与约定端口不一致 | 服务正常启动（本机可绑定任意端口），但 peerServer 仍按约定端口连接对端。建议保持一致 |
| 对端服务尚未启动 | gRPC 连接进入 `TRANSIENT_FAILURE`，自动重连，等服务启动后恢复 |
| 从单节点切换到多节点 | 仅修改 `service_discovery_provide` 配置项，其余不变 |
| 同一机器运行多个相同模块实例 | 静态 Provider 不支持（单节点场景默认一个模块一个实例），多实例需使用 zk/etcd |

---

## 6. 测试计划

### 6.1 单元测试

- `TestStaticResolver_Build`：验证各模块名正确解析为预期地址
- `TestStaticResolver_UnknownModule`：验证未知模块返回错误
- `TestStaticResolver_ConfigOverride`：验证配置覆盖默认端口

### 6.2 集成测试

- 启动两个最小化 gRPC 服务（echo 服务），分别使用 static provider
- 验证 service A 能通过模块名调用 service B
- 验证 service B 下线后 gRPC 连接进入重连，上线后自动恢复

### 6.3 回归测试

- 验证 etcd Provider 分支逻辑未被改动
- 验证 zookeeper Provider 分支逻辑未被改动
- 验证配置 `service_discovery_provide` 为 etcd/zk 时行为完全一致

---

## 7. 实现文件清单

| 文件 | 动作 | 说明 |
|------|------|------|
| `backend/cluster/static_service_discovery_provide.go` | 新增 | 静态 Provider 核心实现 |
| `backend/cluster/static_service_discovery_provide_test.go` | 新增 | 单元测试 |
| `backend/cluster/cluster.go` | 修改 | 增加 `"static"` 分支 |

---

## 8. 未来扩展（可选）

- **环境变量支持**：允许通过 `GOFY_STATIC_PORTS=api:19001,worker:19002` 覆盖端口，进一步减少配置文件
- **自动端口分配**：如果 `defaultModulePorts` 未命中，基于模块名哈希自动分配端口，避免手动维护映射表
- **健康检查集成**：静态模式下可选启用 TCP 端口探测，提前发现对端不可达
