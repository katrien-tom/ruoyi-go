# AGENTS.md

## Scope

本文件适用于整个 `ruoyi-go` 仓库，面向 AI coding agent 与人工开发者。

规则优先级：

1. 用户或任务中的直接指令
2. 本文件
3. 个人习惯或默认偏好

开始改动前，先判断需求落在哪一层，再在该层及必要的相邻层内做最小变更。除非任务明确要求，不要顺手引入第二套分层方式或无关重构。

## Project Architecture

当前仓库是一个按模块拆分的 Gin + GORM + Redis 服务，主要目录职责如下：

- `cmd/main.go`
  - 只负责启动装配。
  - 当前顺序是：初始化日志 -> 读取 `config.yml` -> 初始化 MySQL -> 初始化 Redis -> 注入 `internal/app` -> 初始化 Gin 路由。
  - 不放业务逻辑。

- `internal/app`
  - 全局基础依赖容器。
  - 当前持有 `DB` 和 `Redis`，用于模块装配时共享基础设施依赖。
  - 不放业务状态或业务方法。

- `internal/router`
  - 负责创建 Gin 实例、注册通用中间件、统一注册业务模块。
  - 新模块应通过实现 `Module` 接口接入，而不是在 `cmd/main.go` 中零散注册路由。

- `internal/modules/<name>`
  - 业务模块目录，当前已有 `auth` 和 `user`。
  - 模块内优先按 `router`、`handler`、`service`、`repository`、`entity`、`request`、`response` 分层。

- `internal/middleware`
  - 负责请求级横切逻辑，例如 `request_id`、访问日志、panic recovery、JWT 校验、权限控制。
  - 同类逻辑优先放在这里，不要在 handler 中重复实现。

- `pkg/*`
  - 放跨模块复用的基础设施能力，例如数据库初始化、JWT、日志、统一响应、请求参数校验。
  - `pkg` 只能向下依赖第三方库或标准库，不能反向依赖 `internal/modules`。

- `sql/`
  - 放建表或初始化脚本。

- `config.yml`
  - 本地运行配置文件。
  - 变更配置结构时，代码与配置示例要同步更新。

依赖方向应保持为：

`cmd -> internal/router -> internal/modules/* -> pkg / internal/app`

额外限制：

- `handler` 不直接访问数据库或缓存。
- `repository` 不依赖 Gin 上下文。
- `service` 不直接写 HTTP 响应。

## Request Flow

标准请求链路：

`HTTP Request -> router group -> middleware -> handler -> service -> repository -> response`

各层职责如下：

- `middleware`
  - 处理通用横切逻辑。
  - 典型场景包括链路 ID、日志、panic 恢复、鉴权、权限。

- `handler`
  - 处理 HTTP 适配。
  - 负责参数绑定、基础校验入口、调用 service、将错误映射为统一响应。

- `service`
  - 处理业务规则。
  - 负责事务边界、跨仓储协作、缓存与 JWT 等基础能力协调。

- `repository`
  - 处理 GORM 查询和持久化。
  - 只表达数据访问，不承载 HTTP 语义和复杂业务分支。

- `response`
  - 通过 `pkg/response` 输出统一响应体。
  - 不要在业务层手写不一致的 JSON 结构。

必须遵守：

- handler 使用 `c.Request.Context()` 向下传递上下文。
- repository 使用 `WithContext(ctx)` 执行数据库操作。
- service 使用领域错误表达业务失败，handler 负责将其映射到 `pkg/response` 的错误码。

## Module Convention

新增业务模块时，默认创建到 `internal/modules/<module-name>/`，并优先采用以下文件组织：

- `router.go`
  - 模块装配和路由注册。
  - 在这里组合 repository、service、handler，并暴露 `Register` 方法。

- `handler.go`
  - HTTP 入口层。
  - 只做参数绑定、调用 service、统一响应。

- `service.go`
  - 业务规则层。
  - 放领域错误、业务判断、跨依赖编排。

- `repository.go`
  - 持久化层。
  - 封装 GORM 查询、写入、更新。

- `entity.go`
  - 数据表实体映射。

- `request.go`
  - 请求结构体。

- `response.go`
  - 响应结构体。

补充规则：

- 如果模块不访问数据库，可以没有 `repository.go`。
- 如果接口入参或出参非常简单，可以暂时不拆 `request.go` 或 `response.go`。
- 不要把大量请求/响应结构长期堆在 `handler.go`。
- 新模块接入时，统一在 `internal/router/router.go` 的模块列表中注册。

## Code Style

### General

- 遵循 `gofmt`。
- 导出标识符使用 PascalCase，包内私有标识符使用 camelCase。
- 包名保持简短、小写、无下划线。
- 错误处理优先使用显式 `if err != nil` 早返回。

### Layering

- `context.Context` 作为第一个参数向下传递到 repository 和外部依赖。
- handler 不自行创建脱离请求生命周期的上下文，优先使用 `c.Request.Context()`。
- 请求参数校验优先复用 `pkg/validation`，避免在各模块重复手写同类 `Validate()` 逻辑。
- repository 默认使用 `db.WithContext(ctx)`。
- handler 不直接处理 GORM、Redis、JWT 的底层细节。
- repository 不返回 HTTP 层错误语义。
- service 不直接调用 `c.JSON` 或拼接 HTTP 响应。

### Response And Errors

- 统一使用 `pkg/response` 返回结果。
- 新增业务错误码时，优先扩展 `pkg/response/code.go`。
- service 中定义可判定的领域错误，handler 再映射为响应码和消息。
- 不要随意修改公共响应字段名或既有错误码语义。

### Logging And Constants

- 统一使用 `pkg/logger` 记录结构化日志。
- 避免在正式代码中使用 `fmt.Println` 做调试输出。
- 公共常量、缓存 key 前缀、可复用错误值应集中定义，避免 magic string 分散。

### Comments

- 注释优先解释“为什么”或“约束是什么”。
- 对类型、复杂边界逻辑、非显然设计可保留简洁注释。
- 不要写逐行翻译代码动作的低价值注释。

## Collaboration Rules

### Run And Test

- 默认启动命令：`go run ./cmd`
- 默认测试命令：`go test ./...`
- 只改单一模块时，至少跑受影响包测试。
- 修改公共层、路由、中间件、鉴权、统一响应结构、`pkg/*` 时，跑全量 `go test ./...`。
- 提交前对变更过的 Go 文件执行 `gofmt`。

### Change Boundaries

- 不在 `cmd/main.go` 中堆业务逻辑。
- 不让 `pkg` 反向依赖 `internal/modules`。
- 不在 handler 中直接操作数据库、Redis 或 JWT 细节。
- 不重复实现已有的中间件、日志或统一响应逻辑。
- 不随意破坏既有路由路径和响应契约，除非任务明确要求。
- 不提交临时调试代码、测试密钥、硬编码账号密码。

### Config And Schema Sync

- 变更配置项时，同步更新 `config.yml` 约定。
- 变更表结构时，同步更新 `sql/` 中对应脚本或说明。
- 接口契约变更时，明确影响范围，并确认调用方是否需要同步调整。

### Refactoring

- 优先做与当前任务直接相关的最小重构。
- 不进行与当前需求无关的大规模目录搬迁或风格统一。
- 如果发现现有实现与本规约不一致，可在当前改动范围内做最小修正。
- 如果项目的真实结构已经演进，应在同一任务中同步更新本文件，而不是让规约长期失真。

## Change Checklist

1. 改动是否落在正确分层中。
2. 是否复用了已有的 middleware、logger、response 能力。
3. 是否正确传递了 `context.Context`。
4. 是否避免在 handler 中写持久化逻辑。
5. 行为变更时，是否补充或更新了测试。
6. 涉及配置或表结构时，是否同步更新了 `config.yml` 或 `sql/`。
7. 是否完成了格式化、必要验证和差异检查。
