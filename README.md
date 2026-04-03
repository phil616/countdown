# 该项目被升级为OpenClaw的数据中心，不再维护
# 计时器 — 运维计时管理系统

面向运维工程师的专项计时管理系统，提供倒计时/正计时、项目聚合、TODO待办、通知提醒等功能。

## 功能特性

- **计时单元管理**：支持时间型（倒计时/正计时）和数值型（目标计数/累计计数）
- **项目管理**：将相关计时单元归属到项目中统一管理
- **TODO待办**：轻量级待办事项管理，支持分组和批量操作
- **通知提醒**：后台定时扫描，自动生成到期/超期提醒
- **单用户认证**：JWT + API Token 双认证方式
- **RESTful API**：完整的 REST API，支持脚本和第三方集成

## 技术栈

| 层面 | 技术 |
|------|------|
| 后端 | Go 1.22+, Gin, GORM, SQLite |
| 前端 | Vue 3, Vuetify 3, Pinia, TypeScript |
| 构建 | Vite 5, Docker |

## 快速开始

### 开发模式

**后端：**

```bash
cd ops-timer-backend
go mod tidy
go run ./cmd/server/ --config config.yaml
```

默认管理员账户：`admin` / `admin123`

**前端：**

```bash
cd ops-timer-frontend
npm install
npm run dev
```

前端开发服务器运行在 `http://localhost:5173`，API 请求自动代理到后端 `http://localhost:8080`。

### Docker 部署

```bash
docker compose up -d
```

访问 `http://localhost` 即可使用。

## 项目结构

```
cd-v2/
├── ops-timer-backend/       # Go 后端
│   ├── cmd/server/          # 程序入口
│   ├── internal/
│   │   ├── api/             # HTTP 处理层
│   │   ├── service/         # 业务逻辑层
│   │   ├── repository/      # 数据访问层
│   │   ├── model/           # 数据模型
│   │   ├── dto/             # 数据传输对象
│   │   ├── config/          # 配置管理
│   │   └── pkg/             # 工具包
│   ├── config.yaml          # 配置文件
│   └── Dockerfile
├── ops-timer-frontend/      # Vue 3 前端
│   ├── src/
│   │   ├── api/             # API 请求层
│   │   ├── components/      # 通用组件
│   │   ├── views/           # 页面视图
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── router/          # 路由配置
│   │   ├── types/           # TypeScript 类型
│   │   └── utils/           # 工具函数
│   └── Dockerfile
└── docker-compose.yml
```

## API 认证

**JWT 认证：**

```bash
# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 使用 Token
curl http://localhost:8080/api/v1/units \
  -H "Authorization: Bearer <token>"
```

**API Token 认证：**

```bash
curl http://localhost:8080/api/v1/units \
  -H "X-API-Token: <your_api_token>"
```

## API 文档（OpenAPI）

后端 HTTP 接口的 **OpenAPI 3.0** 规范位于：

- [`ops-timer-backend/docs/openapi.yaml`](ops-timer-backend/docs/openapi.yaml) — 路径、请求/响应模型、认证（JWT / `X-API-Token`）、错误码
- [`ops-timer-backend/docs/README.md`](ops-timer-backend/docs/README.md) — 导入 Postman、Swagger UI、代码生成说明

第三方与自动化工具可直接导入该文件进行对接。

## 配置说明

默认读取工作目录下的 `config.yaml`（可用 `-config` 指定路径）。**若该文件不存在**，则不从磁盘加载 YAML，仅使用环境变量 `TIMER_*` 与程序内默认值（适合 Docker 只注入环境变量、不挂载配置文件）。

参考模板：`ops-timer-backend/config.yaml`；**环境变量键名示例**（全量 `TIMER_*`）：[`ops-timer-backend/env.example`](ops-timer-backend/env.example)。

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `server.port` | 服务端口 | 8080 |
| `database.driver` | 数据库驱动 | sqlite |
| `database.dsn` | 数据库连接 | ./data/timer.db |
| `auth.jwt_secret` | JWT 密钥 | 需修改 |
| `auth.jwt_expiry_hours` | JWT 过期时间（小时） | 24 |
| `scheduler.notification_scan_interval` | 通知扫描间隔 | 1h |

有配置文件时，各项均可被环境变量覆盖；无配置文件时则完全依赖环境变量。前缀均为 `TIMER_`，嵌套键中的 `.` 写成 `_`，如 `TIMER_AUTH_JWT_SECRET`、`TIMER_DATABASE_DSN`。
