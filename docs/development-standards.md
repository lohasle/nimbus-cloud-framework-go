# 开发规范

- 包路径固定为 `github.com/lohasle/nimbus-cloud-framework-go`。
- 新接口同步 Swagger 和路由契约测试。
- 跨服务调用通过网关或服务发现，不在业务模块硬编码地址。
- 数据写入检查错误；金额、余额、积分和多表变更使用事务。
- 前端保留现有表格和页面组件体系，品牌改动集中在 Logo、文案、登录页和 Loading。

验证命令：

```bash
cd backend
go test ./...
make build

cd ../frontend
pnpm ts:check
pnpm lint
pnpm build:local
```
