# 架构说明

```text
Vue 运营后台
      |
Gateway :58080
      |
 +----+----+------+--------+----------+
 |         |      |        |          |
System   Infra   Pay     Member    Business
:58081  :58082  :58085  :58087    :58090
                              application / im / app Health
      |
 MySQL 8.4 · Redis 7.4 · Nacos 3
```

- Gateway 负责统一入口、Trace ID、Nacos 动态解析和反向代理。
- System 负责运营用户、RBAC、租户、字典、OAuth2、通知、邮件和短信。
- Infra 负责参数、文件、日志、数据源、定时任务和 Redis 监控。
- Member 与 Pay 保持现有业务中心边界。
- Business 只承载三个预留模块的健康检查。

所有服务共享接口契约，但各自拥有迁移和启动入口；跨服务地址不得硬编码到业务代码。
