# 测试策略

1. 单元测试：Token 类型、刷新轮换、菜单树、任务调度。
2. 路由契约：当前菜单依赖的路径在对应服务中存在。
3. 服务集成：六个进程注册到 Nacos，并可通过 Gateway 访问。
4. 数据集成：MySQL 初始化幂等、Redis 监控可用、租户隔离正确。
5. 前端检查：TypeScript、Lint、生产构建和浏览器页面验收。

统一接口冒烟使用 `scripts/test-functional.sh`。
