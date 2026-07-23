# SPEC-003 System/Infra 公共能力对齐

状态：已完成（2026-07-23）

## 目标

- 将已验收的 Nimbus 单体 System/Infra 通用实现同步到 Cloud。
- 保持六进程、Nacos、Gateway、本地静态回退和 MySQL 8.4。
- 保留 Member、Pay；Application、IM、App 继续只提供 Health。

## 验收

- System 20 个管理页面、Infra 9 个管理页面通过网关访问。
- Access Token 与 Refresh Token 分离并支持轮换。
- 定时任务、错误日志、文件、数据源和 Redis 监控接口可用。
- Go 测试、六服务构建、前端静态检查和统一接口冒烟通过。
