# SPEC-002 Java Cloud 底座迁移到 Go

## 状态

- 当前：已完成
- 源工程：`nimbus-cloud-framework/backend`
- 目标工程：`nimbus-cloud-framework-go/backend`
- 数据库基线：MySQL 8.4

## 目标与边界

1. Go Cloud 必须迁移 Java Cloud 的 System、Infra、Member、Pay 通用中心能力，保持前端 `/admin-api` 契约兼容。
2. `application`、`im`、`app` 是自有预留边界，本阶段仅保留 Health，不预建业务实体。
3. 根目录固定为 `frontend/`、`backend/`；后端在 `internal/modules/` 按中心分模块，在 `cmd/` 按独立进程提供启动入口。
4. Gateway 只负责统一入口、认证上下文传播、Trace ID 与服务路由；服务通过 Nacos 注册发现并保留本地回退。
5. 所有公开接口必须包含 Swagger 注释。

## 目标结构

```text
nimbus-cloud-framework-go/
├── frontend/
└── backend/
    ├── cmd/{gateway,system,infra,member,pay,business}/
    ├── internal/platform/
    ├── internal/modules/{system,infra,member,pay,application,im,app}/
    ├── migrations/
    └── docs/
```

## 验收

- MySQL 8.4、Nacos 与全部服务可以从空环境初始化。
- Gateway、System、Infra、Member、Pay 与 Business 独立进程正常启动。
- 登录、首页及所有已展示菜单无 404/服务器错误。
- `go test ./...`、Swagger、全部二进制和前端生产构建通过。
- 不以 Health 空壳冒充 System、Infra、Member、Pay 的 Java 能力迁移完成。

## 当前已验证闭环

- Gateway、System、Infra、Member、Pay、Business 六个进程独立启动并通过健康检查。
- System 已实现运营用户、RBAC、菜单、组织、租户、字典、审计日志、OAuth2、通知、邮件和短信管理。
- Infra 已实现参数、文件、访问/错误日志、数据源、定时任务、任务日志及 Redis 监控；Member 与 Pay 保持既有管理闭环。
- Application、IM、App 仅由 Business 进程提供 Health，不建立业务表。
- Swagger 已生成 101 条路径；Go 全量测试/编译和前端类型检查/生产构建通过。

未在当前菜单开放的 Java 扩展能力（代码生成、第三方支付真实适配、钱包/转账扩展等）留作后续按需迁移，不提供伪实现。
