module github.com/odysseythink/gofy/backend

go 1.25.7

// ====== mlib.com packages → local lib/ ======
replace (
	mlib.com/gofy/server/utils/jieba => ./utils/jieba
	mlib.com/zkmgr => ./lib/zkmgr
)

// ====== Patched vendor packages (custom grpc with RegisterServiceWithoutDesc etc.) ======
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
    github.com/robfig/cron/v3 v3.0.0
	github.com/odysseythink/mlog v0.0.2
	github.com/odysseythink/mrun v0.0.2
	github.com/amikos-tech/chroma-go v0.4.1

	// Utilities
	github.com/anaskhan96/soup v1.2.5

	// Web framework
	github.com/gin-gonic/gin v1.10.0
	github.com/go-playground/validator/v10 v10.30.1

	// Authentication
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.2.0
	github.com/pkoukk/tiktoken-go v0.1.7

	// Redis
	github.com/redis/go-redis/v9 v9.7.0

	// UUID
	github.com/satori/go.uuid v1.2.0
	github.com/seccomp/libseccomp-golang v0.11.1
	github.com/shopspring/decimal v1.4.0
	github.com/spf13/viper v1.19.0
	github.com/unrolled/secure v1.16.0
	github.com/xuri/excelize/v2 v2.9.0

	// etcd
	go.etcd.io/etcd/client/v3 v3.5.17

	// Crypto & text
	golang.org/x/crypto v0.46.0
	golang.org/x/exp/errors v0.0.0-20260312153236-7ab1446f8b90
	golang.org/x/sync v0.19.0
	golang.org/x/text v0.32.0
	golang.org/x/time v0.9.0

	// gRPC
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.10

	// YAML
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/datatypes v1.2.5
	gorm.io/driver/mysql v1.5.7
	gorm.io/gen v0.3.26

	// Database
	gorm.io/gorm v1.25.12
	// Self-owned mlib.com packages
	github.com/odysseythink/confy v0.0.1
	mlib.com/zkmgr v0.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/amikos-tech/chroma-go-local v0.3.4 // indirect
	github.com/amikos-tech/pure-onnx v0.0.1 // indirect
	github.com/amikos-tech/pure-tokenizers v0.1.5 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/ebitengine/purego v0.10.0 // indirect

	// HTML parsing
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect

	// Protobuf (legacy)
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xuri/efp v0.0.0-20240408161823-9ad904a10d6d // indirect
	github.com/xuri/nfp v0.0.0-20240318013403-ab9948c2c4a7 // indirect
	go.etcd.io/etcd/api/v3 v3.6.1 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/exp v0.0.0-20241217172543-b2144cdd0a67 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gorm.io/hints v1.1.0 // indirect
	gorm.io/plugin/dbresolver v1.5.0 // indirect
)

require (
	github.com/jackc/pgx/v5 v5.6.0
	gorm.io/driver/postgres v1.6.0
	mlib.com/gofy/server/utils/jieba v0.0.0-00010101000000-000000000000
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	go.etcd.io/etcd/client/pkg/journal v0.0.0-00010101000000-000000000000 // indirect
)
