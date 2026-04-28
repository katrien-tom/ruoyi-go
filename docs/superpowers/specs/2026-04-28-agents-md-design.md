# ruoyi-go `AGENTS.md` 设计说明

## 1. 背景

当前仓库已经形成了清晰但仍处于早期阶段的 Go Web 项目骨架：

- `cmd/main.go` 负责启动装配
- `internal/app` 维护全局基础依赖
- `internal/router` 统一注册中间件和业务模块
- `internal/modules/*` 按模块拆分 `handler/service/repository/entity`
- `internal/middleware` 承载请求级横切逻辑
- `pkg/*` 提供日志、数据库、JWT、统一响应等通用能力

仓库尚无 `AGENTS.md`，因此缺少一份同时面向 AI agent 与人工开发者的仓库内执行规约。现阶段最重要的不是补一份泛化说明，而是沉淀一份能直接指导改动边界、分层方式和协作检查的规范文件。

## 2. 目标

本次新增的 `AGENTS.md` 应满足以下目标：

1. 作为仓库根目录的执行规约，覆盖整个项目。
2. 同时服务 AI agent 与人工开发者，但以“必须遵守的工程规则”为主，而不是 README 风格介绍。
3. 基于仓库当前真实结构描述项目架构，避免写成与代码脱节的模板。
4. 明确新增功能、修复缺陷、重构代码时的目录职责、依赖方向和改动边界。
5. 提供统一的代码风格要求和提交前检查项，降低模块间风格漂移。

## 3. 非目标

以下内容不属于本次 `AGENTS.md` 的目标范围：

- 详细业务需求说明
- API 文档或接口清单
- 运维部署手册
- 替代 `README` 的完整项目介绍
- 大规模重构现有代码以适配规约

## 4. 目标读者与文档定位

`AGENTS.md` 的主要读者包括：

- 在仓库内执行改动的 AI coding agent
- 日常维护仓库的人工开发者

文档定位为“仓库内执行规约”。其内容应强调：

- 先判断改动属于哪一层，再落代码
- 默认沿用当前项目分层，不引入第二套风格
- 对公共层、响应契约、配置与 SQL 的改动应进行额外检查

## 5. 文档结构设计

`AGENTS.md` 将采用“执行规约优先”的结构：

1. `Scope`
2. `Project Architecture`
3. `Request Flow`
4. `Module Convention`
5. `Code Style`
6. `Collaboration Rules`
7. `Change Checklist`

这样安排的原因是：

- 开头先说明适用范围和优先级，便于 AI 与开发者快速对齐行为。
- 中间部分基于真实代码结构说明怎么改、在哪改。
- 结尾部分收束为协作规则和自检清单，提升执行性。

## 6. 项目架构说明内容

`AGENTS.md` 的架构章节将围绕“目录职责 + 依赖方向 + 扩展方式”展开。

### 6.1 目录职责

- `cmd/main.go`
  - 负责程序启动装配。
  - 当前启动顺序为：初始化日志 -> 读取配置 -> 初始化 MySQL -> 初始化 Redis -> 注入全局依赖 -> 初始化路由。
  - 不承载业务逻辑。

- `internal/app`
  - 当前作为全局基础依赖容器，持有 `DB` 与 `Redis`。
  - 只允许存放基础设施级依赖，不放业务状态和业务方法。

- `internal/router`
  - 负责 Gin 实例初始化、通用中间件注册以及模块注册。
  - 新模块通过实现 `Module` 接口接入统一路由注册流程。

- `internal/modules/<name>`
  - 负责具体业务模块。
  - 当前 `auth`、`user` 模块已经体现标准分层：`router`、`handler`、`service`、`repository`、`entity`、`dto`、`vo`。

- `internal/middleware`
  - 负责请求 ID、日志、panic recovery、JWT、权限控制等横切逻辑。
  - 业务模块不应在 handler 中重复实现同类逻辑。

- `pkg/*`
  - 负责跨模块复用的基础设施能力，例如日志、数据库、JWT、统一响应。
  - 不应包含具体业务语义，也不应反向依赖 `internal/modules`。

- `sql/`
  - 负责数据库初始化或变更脚本。

- `config.yml`
  - 负责本地运行配置。
  - 新增配置项时，应同步更新对应说明或示例。

### 6.2 依赖方向

`AGENTS.md` 会明确以下依赖方向：

`cmd -> internal/router -> internal/modules/* -> pkg / internal/app`

并补充限制：

- `pkg` 不得依赖 `internal/modules`
- `handler` 不得直接访问数据库或缓存
- `repository` 不得依赖 Gin 上下文
- `service` 不直接输出 HTTP 响应

### 6.3 扩展方式

新增模块时，应优先沿用当前模块化方式：

1. 在 `internal/modules/<module-name>/` 下创建对应文件。
2. 在模块内完成依赖装配。
3. 通过 `internal/router/router.go` 中的模块列表统一注册。

不应为单个模块引入新的路由管理方式或新的分层体系。

## 7. 请求流与分层约束

`AGENTS.md` 将定义标准请求链路：

`HTTP Request -> router group -> middleware -> handler -> service -> repository -> response`

并为每一层写明责任边界：

- `middleware`
  - 负责通用横切逻辑。
- `handler`
  - 负责参数绑定、调用 service、错误映射、统一响应。
- `service`
  - 负责业务规则、事务边界、跨仓储协作、缓存/JWT 等基础能力协调。
- `repository`
  - 负责 GORM 查询和持久化。
- `response`
  - 负责统一返回格式，避免各模块自行拼接 JSON 结构。

## 8. 模块组织约定

`AGENTS.md` 将约定新模块默认采用以下文件组织：

- `router.go`：模块装配与路由注册
- `handler.go`：HTTP 适配层
- `service.go`：业务逻辑层
- `repository.go`：持久化访问层
- `entity.go`：数据库实体映射
- `dto.go`：请求结构
- `vo.go`：响应结构

补充规则：

- 如果模块不访问数据库，可以暂时没有 `repository.go`
- 如果接口入参/出参非常简单，可以暂时不拆 `dto.go` 或 `vo.go`
- 但不要将大量请求/响应结构长期堆在 `handler.go`

## 9. 代码风格要求

`AGENTS.md` 的代码风格部分将贴合当前仓库而不是只写 Go 通用规则。

### 9.1 通用 Go 风格

- 遵循 `gofmt`
- 导出标识符使用 PascalCase，包内私有标识符使用 camelCase
- 包名保持简短、小写、无下划线
- 错误处理采用显式早返回

### 9.2 分层风格

- `context.Context` 作为第一个参数向下传递到 repository 和外部依赖
- handler 使用 `c.Request.Context()`，不自行创建脱离请求生命周期的上下文
- repository 查询默认使用 `WithContext(ctx)`
- handler 不直连数据库、Redis、JWT 细节
- repository 不返回 HTTP 语义错误
- service 使用领域错误表达业务失败，handler 将其映射到统一响应码

### 9.3 公共能力约定

- 统一使用 `pkg/logger` 记录结构化日志
- 统一使用 `pkg/response` 输出响应体
- 新增业务错误码优先扩展 `pkg/response/code.go`
- 涉及鉴权时，优先复用现有 JWT 中间件和上下文字段约定
- 公共常量、错误值、缓存 key 前缀应集中定义，避免 magic string 分散

### 9.4 注释风格

- 注释优先解释“为什么”
- 类型、复杂边界逻辑、非显然约束可以保留简洁注释
- 避免逐行解释代码动作的低信息量注释

## 10. 协作规则

`AGENTS.md` 将包含一组默认必须遵守的协作规则。

### 10.1 本地运行与测试

- 默认启动命令：`go run ./cmd`
- 默认测试命令：`go test ./...`
- 修改单一模块时，至少执行受影响包测试
- 修改公共层、路由、中间件、鉴权、统一响应结构时，应执行全量测试
- 提交前应确认代码已按 `gofmt` 格式化

### 10.2 改动边界

- 不在 `cmd/main.go` 中堆放业务逻辑
- 不让 `pkg` 反向依赖业务模块
- 不在 handler 中直接操作数据库、Redis 或 JWT 细节
- 不重复实现已有中间件、日志或响应逻辑
- 不随意变更统一响应体字段名或响应码语义
- 不提交临时调试代码、测试密钥或硬编码账号密码

### 10.3 配置与数据结构同步

- 变更配置项时，同步更新 `config.yml` 约定
- 变更表结构时，同步更新 `sql/` 脚本或相关说明
- 涉及接口契约变更时，确保调用方能感知并评估影响

### 10.4 重构规则

- 重构优先在当前任务相关范围内进行
- 不进行与当前需求无关的大规模目录搬迁或风格统一
- 若发现现有实现偏离规约，可在不扩大风险的前提下做最小修正

## 11. 变更检查清单

`AGENTS.md` 末尾将附一个简短 checklist，用于每次改动自检：

1. 改动是否落在正确分层中
2. 是否复用统一响应、日志、中间件能力
3. 是否正确传递 `context.Context`
4. 是否避免在 handler 中写持久化逻辑
5. 是否补充或更新必要测试
6. 是否同步更新配置或 SQL
7. 是否完成 `gofmt` 与对应测试

## 12. 预期结果

完成后，仓库根目录将新增一份 `AGENTS.md`，具备以下特征：

- 不是泛化模板，而是准确反映 `ruoyi-go` 当前分层结构
- 同时适用于 AI 与人工开发者
- 对“新增功能如何放置、已有代码如何修改、提交前如何检查”给出清晰约束
- 在项目规模继续增长前，先统一工程协作方式
