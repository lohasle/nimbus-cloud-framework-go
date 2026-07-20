# Nimbus Cloud Framework Go

Nimbus Framework 的 Go 1.26 微服务版，统一采用 Gin 1.12、新 UI 与 MySQL 8.4。

## 进程边界

| 进程 | 端口 | 职责 |
| --- | ---: | --- |
| `nimbus-gateway` | 58080 | 统一入口、Trace ID、动态服务解析与反向代理 |
| `nimbus-system` | 58081 | 租户、运营用户、认证与权限 |
| `nimbus-infra` | 58082 | 基础设施扩展边界（当前 health） |
| `nimbus-pay` | 58085 | 支付中心扩展边界（当前 health） |
| `nimbus-member` | 58087 | 会员中心扩展边界（当前 health） |
| `nimbus-business` | 58090 | application / IM / APP 聚合扩展边界（当前 health） |

Nacos 负责注册发现；没有配置 Nacos 时，网关自动回退到本地静态地址，便于开发和故障恢复。默认数据库为 MySQL 8.4。

## 启动

```bash
./scripts/init-local.sh
make test build
./scripts/start-cloud.sh
./scripts/status-cloud.sh
```

默认租户 `Nimbus Framework`，账号 `admin / admin123`。生产环境必须替换 JWT 密钥和初始化口令。

OpenAPI 文档由 Swagger 注释生成在 `docs/`；系统服务运行时可访问 `http://localhost:58081/swagger/index.html`。

