# Agent 上下文

Nimbus Cloud Framework Go 是面向 App 后台与运营后台的 Go 微服务脚手架。

## 不可变边界

- Gateway、System、Infra、Member、Pay、Business 六个 Go 进程。
- Nacos 服务注册发现与本地静态回退。
- MySQL 8.4 主数据库，Redis 基础设施监控。
- `system`、`infra`、`member`、`pay` 为现有功能中心。
- `application`、`im`、`app` 只提供 Health 示例，不预建业务实体。

`muse-app-go` 与 `nimbus-framework-go` 可作为公共 Go 能力参考，但同步时必须保留 Cloud 的包名、服务边界、网关和注册发现。
