# 运维计时管理系统 — 后端文档

本目录包含 **HTTP API** 的规范说明，供前端、脚本、第三方集成与自动化工具使用。

## OpenAPI 规范

| 文件 | 说明 |
|------|------|
| [`openapi.yaml`](./openapi.yaml) | OpenAPI **3.0.3** 完整描述（路径、请求/响应模型、认证方式、错误码） |

### 使用方式

1. **导入 Postman / Insomnia / Apifox**  
   选择「Import」→「OpenAPI」→ 选择本仓库的 `ops-timer-backend/docs/openapi.yaml`。

2. **Swagger UI（本地预览）**  
   ```bash
   npx @redocly/cli preview-docs ops-timer-backend/docs/openapi.yaml
   ```
   或使用 Docker：
   ```bash
   docker run -p 8080:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v "%cd%/ops-timer-backend/docs:/docs" swaggerapi/swagger-ui
   ```
   浏览器访问 `http://localhost:8080`（Windows 下路径请按实际调整）。

3. **代码生成**  
   使用 [OpenAPI Generator](https://openapi-generator.tech/) 等工具，根据 `openapi.yaml` 生成各语言客户端 SDK。

   ```bash
   docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
     -i /local/ops-timer-backend/docs/openapi.yaml \
     -g go \
     -o /local/out/go-client
   ```

4. **CI 校验**  
   可在流水线中执行 `npx @redocly/cli lint ops-timer-backend/docs/openapi.yaml` 校验规范文件。

## 基础约定

- **Base URL**：默认 `http://localhost:8080`（部署时替换为实际域名与端口）。
- **API 前缀**：业务接口为 `/api/v1`；健康检查为 `/health`。
- **认证**（需登录的接口）：
  - **Bearer JWT**：`Authorization: Bearer <access_token>`（登录接口返回的 `token`）。
  - **API Token**：`X-API-Token: <api_token>`（用户可在系统内生成/轮换，适合脚本与集成）。
- **响应格式**：JSON。成功时 `code` 为 `0`；失败时 HTTP 状态码与业务 `code` 见 OpenAPI 文档中的「错误码」说明。
- **分页列表**：查询参数 `page`（默认 1）、`page_size`（默认 20）；响应体在 `meta` 中携带 `total`、`total_pages` 等。

## 与配置的关系

- 服务监听地址、CORS、JWT、数据库等见项目根目录 `README.md` 与 `config.yaml` / 环境变量前缀 `TIMER_`。
- OAuth 相关接口仅在后台启用 OAuth 时可用；详见 OpenAPI 中 `/auth/oauth/*` 说明。

## 维护说明

- 新增或变更 API 时，请**同步更新** `openapi.yaml`，并保持与 `internal/api/router` 及 DTO 一致。
- 版本号：在 `openapi.yaml` 的 `info.version` 中与发布版本对齐（或采用独立 API 版本号策略并在文档中说明）。
