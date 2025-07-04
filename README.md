# BoCloudStore
A cloud native high concurrency storage platform based on go.

## directory
BoCloudStore/  
├── cmd  
│   └── cloudStore-server  
│       └── main.go              # 主入口文件  
├── internal  
│   ├── api  
│   │   ├── handler              # HTTP 请求处理器  
│   │   │   ├── upload_handler.go  
│   │   │   ├── download_handler.go  
│   │   │   └── auth_handler.go  
│   │   ├── middleware           # Gin 中间件  
│   │   │   ├── jwt_auth.go  
│   │   │   ├── rate_limiter.go  
│   │   │   └── request_logger.go  
│   │   └── router.go            # 路由配置  
│   ├── service  
│   │   ├── upload_service.go    # 分块上传业务逻辑  
│   │   ├── metadata_service.go  # 元数据管理  
│   │   └── auth_service.go      # 认证授权服务  
│   ├── storage  
│   │   ├── minio_client.go      # MinIO 存储操作封装  
│   │   ├── redis_cache.go       # Redis 缓存操作  
│   │   └── postgres_store.go    # PostgreSQL 元数据存储  
│   ├── config  
│   │   └── config.go            # 配置加载与解析  
│   ├── util  
│   │   ├── crypto_util.go       # 加密工具  
│   │   ├── file_util.go         # 文件处理工具  
│   │   └── logger_util.go       # 日志工具  
│   └── model  
│       ├── file_model.go        # 文件元数据结构体  
│       └── user_model.go        # 用户模型  
├── pkg  
│   ├── minio                    # MinIO SDK 封装（可选）  
│   └── virusscan                # 病毒扫描模块  
├── deployments  
│   ├── docker  
│   │   └── Dockerfile           # Docker 构建文件  
│   ├── kubernetes  
│   │   ├── deployment.yaml      # K8s 部署配置  
│   │   ├── service.yaml  
│   │   └── hpa.yaml             # 自动扩缩容配置  
│   └── helm  
│       └── cloudstor            # Helm Chart  
├── configs  
│   ├── config.yaml              # 主配置文件  
│   └── config.dev.yaml          # 开发环境配置  
├── scripts  
│   ├── build.sh                 # 构建脚本  
│   ├── deploy.sh                # 部署脚本  
│   └── migrate.sql              # 数据库迁移脚本  
├── test  
│   ├── integration              # 集成测试  
│   └── unit                     # 单元测试  
├── docs  
│   ├── api.md                   # API 文档  
│   └── design.md                # 架构设计文档  
├── .gitignore  
├── go.mod  
├── go.sum  
└── Makefile                     # 构建管理

## another
BoCloudStore/  
├── cmd/    
│   └── main.go           # 项目入口文件  
├── internal/  
│   ├── access/           # 接入层相关代码  
│   │   ├── handlers/     # HTTP 请求处理函数  
│   │   │   └── upload.go # 文件上传处理  
│   │   ├── middleware/   # 中间件代码  
│   │   │   └── auth.go   # JWT 鉴权中间件  
│   │   └── router.go     # 路由配置  
│   ├── business/         # 业务逻辑层相关代码  
│   │   └── chunk_upload/ # 分块上传服务  
│   │       └── service.go # 分块上传服务实现  
│   ├── metadata/         # 元数据管理相关代码  
│   │   ├── db/           # 数据库操作  
│   │   │   └── postgres.go # PostgreSQL 操作  
│   │   ├── cache/        # 缓存操作  
│   │   │   └── redis.go  # Redis 操作  
│   │   └── service.go    # 元数据管理服务  
│   ├── storage/          # 对象存储引擎相关代码  
│   │   └── minio.go      # MinIO 操作  
│   ├── security/         # 安全与治理模块相关代码  
│   │   ├── antivirus/    # 病毒扫描  
│   │   │   └── clamav.go # ClamAV 引擎调用  
│   │   └── authz/        # 权限控制  
│   │       └── rbac.go   # RBAC 模型实现  
│   └── observability/    # 运维支撑相关代码  
│       ├── monitoring/   # 监控指标收集  
│       │   └── prometheus.go # Prometheus 集成  
│       └── logging/      # 日志记录  
│           └── loki.go   # Loki 集成  
├── pkg/  
│   ├── config/           # 配置文件处理   
│   │   └── config.go     # 配置加载  
│   ├── pool/             # 协程池实现  
│   │   └── worker_pool.go # 协程池代码  
│   └── utils/            # 通用工具函数  
│       └── utils.go      # 工具函数实现  
├── test/  
│   └── unit/             # 单元测试代码  
│       ├── access_test.go # 接入层单元测试  
│       ├── business_test.go # 业务逻辑层单元测试  
│       ├── metadata_test.go # 元数据管理单元测试  
│       ├── storage_test.go # 对象存储引擎单元测试  
│       └── security_test.go # 安全与治理模块单元测试  
├── deploy/  
│   └── kubernetes/       # Kubernetes 部署文件  
│       ├── helm/         # Helm Chart 目录  
│       │   ├── Chart.yaml # Helm Chart 元数据  
│       │   ├── values.yaml # Helm Chart 配置值  
│       │   └── templates/ # Kubernetes 资源模板  
│       │       ├── deployment.yaml # Deployment 模板  
│       │       └── service.yaml    # Service 模板  
│       └── hpa.yaml      # Horizontal Pod Autoscaler 配置  
├── docs/  
│   ├── api/              # API 文档  
│   │   └── swagger.yaml  # Swagger API 定义  
│   └── architecture/     # 架构设计文档  
│       └── design.md     # 架构设计说明  
├── .gitignore            # Git 忽略文件  
├── go.mod                # Go 模块文件  
└── go.sum                # Go 模块依赖文件   