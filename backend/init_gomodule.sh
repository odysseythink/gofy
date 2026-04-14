#!/bin/bash
#
# Go Modules 初始化脚本
# 将项目从 GOPATH 模式 (GO111MODULE=off) 迁移到 Go Modules 模式
#
# 使用方法:
#   cd gofy/backend
#   bash init_gomodule.sh
#
# 前提条件:
#   - GOPATH 下的源码包可用 (mlib.com/*, github.com/odysseythink/*)
#   - 原始 vendor/ 目录存在 (含魔改过的 grpc 等)
#
set -e

GOPATH_SRC="${GOPATH:-$HOME/workspace/go_work/gopath}/src"
BACKEND_DIR="$(cd "$(dirname "$0")" && pwd)"

cd "$BACKEND_DIR"
echo "=== Go Modules 初始化: $BACKEND_DIR ==="

# ============================================================
# Step 1: 复制自研依赖包到 lib/
# ============================================================
echo "[1/7] 复制自研依赖包到 lib/ ..."

mkdir -p lib

# mlib.com/* 自研包 (从 GOPATH 复制)
cp -r "$GOPATH_SRC/mlib.com/zkmgr"   lib/zkmgr


# ============================================================
# Step 2: 修复 lib/ 内部的 import 路径
# ============================================================
echo "[2/7] 修复 lib/ 内部 import 路径 ..."

echo "  已清理嵌套 go.mod"

# ============================================================
# Step 3: 为每个 lib/ 包创建 go.mod
# ============================================================
echo "[3/7] 为 lib/ 包创建 go.mod ..."

GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+' | sed 's/go//')


go $GO_VERSION
EOF


go $GO_VERSION

require (
	github.com/go-openapi/errors v0.22.0
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	golang.org/x/net v0.33.0
	golang.org/x/sync v0.10.0
)

EOF


go $GO_VERSION

require (
	github.com/golang/protobuf v1.5.4
)

EOF

go $GO_VERSION

require (
	github.com/fsnotify/fsnotify v1.8.0
	go.etcd.io/etcd/api/v3 v3.5.17
	go.etcd.io/etcd/client/v3 v3.5.17
	gopkg.in/yaml.v3 v3.0.1
)
EOF

go $GO_VERSION

EOF

cat > lib/zkmgr/go.mod << EOF
module mlib.com/zkmgr

go $GO_VERSION

require (
	github.com/odysseythink/mrun v0.0.2
	github.com/odysseythink/mlog v0.0.2
)


EOF

echo "  已创建 6 个 go.mod"

# ============================================================
# Step 4: 将魔改过的 vendor 包移到 lib/patched/
# ============================================================
echo "[4/7] 将魔改的 vendor 包移到 lib/patched/ ..."

if [ -d vendor ] && [ ! -d lib/patched ]; then
    cp -r vendor lib/patched
    echo "  已复制 vendor/ -> lib/patched/"
elif [ -d lib/patched ]; then
    echo "  lib/patched/ 已存在, 跳过"
else
    echo "  错误: vendor/ 不存在, 请先从 git 恢复: git checkout -- vendor/"
    exit 1
fi

# 为 etcd/client/pkg/journal 创建 go.mod (它没有自己的)
if [ ! -f lib/patched/go.etcd.io/etcd/client/pkg/journal/go.mod ]; then
    cat > lib/patched/go.etcd.io/etcd/client/pkg/journal/go.mod << EOF
module go.etcd.io/etcd/client/pkg/journal

go 1.23.0
EOF
    echo "  已为 etcd journal 创建 go.mod"
fi

# ============================================================
# Step 5: 删除旧 vendor, 创建主 go.mod
# ============================================================
echo "[5/7] 创建主 go.mod ..."

# 删除旧 vendor (已备份到 lib/patched/)
rm -rf vendor

cat > go.mod << GOMOD
module backend

go $GO_VERSION

// ====== 自研包 → lib/ ======
replace (
	mlib.com/zkmgr => ./lib/zkmgr
)

// ====== 魔改的第三方包 → lib/patched/ (含 grpc RegisterServiceWithoutDesc 等) ======
replace (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 => ./lib/patched/github.com/grpc-ecosystem/go-grpc-middleware/v2
	github.com/grpc-ecosystem/grpc-gateway/v2 => ./lib/patched/github.com/grpc-ecosystem/grpc-gateway/v2
	go.etcd.io/etcd/api/v3 => ./lib/patched/go.etcd.io/etcd/api/v3
	go.etcd.io/etcd/client/pkg/journal => ./lib/patched/go.etcd.io/etcd/client/pkg/journal
	go.etcd.io/etcd/client/pkg/v3 => ./lib/patched/go.etcd.io/etcd/client/pkg/v3
	go.etcd.io/etcd/client/v3 => ./lib/patched/go.etcd.io/etcd/client/v3
	google.golang.org/genproto/googleapis/api => ./lib/patched/google.golang.org/genproto/googleapis/api
	google.golang.org/genproto/googleapis/rpc => ./lib/patched/google.golang.org/genproto/googleapis/rpc
	google.golang.org/grpc => ./lib/patched/google.golang.org/grpc
	google.golang.org/protobuf => ./lib/patched/google.golang.org/protobuf
)

require (
	github.com/amikos-tech/chroma-go v0.4.1
	github.com/anaskhan96/soup v1.2.5
	github.com/gin-gonic/gin v1.10.0
	github.com/go-playground/validator/v10 v10.30.1
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.2.0
	github.com/odysseythink/confy v0.0.1
	github.com/odysseythink/mlog v0.0.2
	github.com/odysseythink/mrun v0.0.2
	github.com/pkoukk/tiktoken-go v0.1.7
	github.com/redis/go-redis/v9 v9.7.0
	github.com/satori/go.uuid v1.2.0
	github.com/seccomp/libseccomp-golang v0.11.1
	github.com/shopspring/decimal v1.4.0
	github.com/spf13/viper v1.19.0
	github.com/unrolled/secure v1.16.0
	github.com/xuri/excelize/v2 v2.9.0
	go.etcd.io/etcd/client/v3 v3.5.17
	golang.org/x/crypto v0.46.0
	golang.org/x/exp/errors v0.0.0-20260312153236-7ab1446f8b90
	golang.org/x/sync v0.19.0
	golang.org/x/text v0.32.0
	golang.org/x/time v0.9.0
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.10
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/datatypes v1.2.5
	gorm.io/driver/mysql v1.5.7
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.12
	mlib.com/zkmgr v0.0.0
)
GOMOD

echo "  已创建 go.mod (module backend)"

# ============================================================
# Step 6: 删除嵌套 go.mod 冲突文件
# ============================================================
echo "[6/7] 清理冲突文件 ..."

rm -f utils/jieba/go.mod 2>/dev/null || true
echo "  已删除 utils/jieba/go.mod"

# ============================================================
# Step 7: 解析依赖, 生成 vendor
# ============================================================
echo "[7/7] 解析依赖 (go mod tidy + go mod vendor) ..."

go mod tidy
echo "  go mod tidy 完成 (go.sum 已生成)"

go mod vendor
echo "  go mod vendor 完成 (vendor/modules.txt 已生成)"

# ============================================================
# 验证
# ============================================================
echo ""
echo "=== 验证 ==="
echo "go.mod 行数: $(wc -l < go.mod)"
echo "go.sum 行数: $(wc -l < go.sum)"
echo "vendor/modules.txt 行数: $(wc -l < vendor/modules.txt)"

echo ""
echo "测试编译 FrameWorkServer ..."
go build -o /dev/null main/FrameWorkServer/main.go && echo "  ✓ FrameWorkServer 编译通过"

echo "测试编译 link.so ..."
go build -buildmode=plugin -o /dev/null main/link/main.go main/link/frontend.go && echo "  ✓ link.so 编译通过"

echo ""
echo "=== Go Modules 初始化完成 ==="
echo ""
echo "目录结构:"
echo "  lib/                  # 自研依赖包 (replace 指向)"
echo "    zkmgr/              # ZooKeeper 管理"
echo "    patched/            # 魔改的第三方包"
echo "      google.golang.org/grpc/  # 含 RegisterServiceWithoutDesc"
echo "  vendor/               # go mod vendor 自动生成"
echo "  go.mod                # 模块定义"
echo "  go.sum                # 依赖校验"
echo ""
echo "构建命令不变: ./build"
