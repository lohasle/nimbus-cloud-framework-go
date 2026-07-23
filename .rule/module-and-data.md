# 模块与数据规则

- 未经用户确认，不新增、删除、合并或拆分服务及模块。
- Gateway、System、Infra、Member、Pay、Business 六进程边界必须保留。
- `application`、`im`、`app` 由 Business 承载，默认只提供 Health。
- 数据库固定为 MySQL 8.4；Redis 是基础设施依赖，不替代主数据库。
- 每个服务只迁移和初始化自己负责的数据表，重复启动必须幂等。
