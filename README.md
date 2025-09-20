# Go RSS 订阅器

基于 **Go + PostgreSQL** 的 RSS 订阅管理服务，支持用户订阅、管理和获取最新文章。

## 功能特性

- **用户管理**：注册、API Key 认证
- **RSS源管理**：添加、查看RSS源
- **订阅管理**：关注/取消关注RSS源，查看已关注源
- **文章管理**：后台定时抓取RSS源，获取订阅文章列表
- **认证方式**：API Key（通过 `Authorization` 请求头传入）

## 技术栈

- **Go**：主要语言
- **PostgreSQL**：数据存储
- **sqlc**：SQL代码自动生成
- **goose**：数据库迁移管理
- **chi**：HTTP路由框架
- **godotenv**：环境变量加载

## 项目结构

```
├── main.go              # 主程序入口
├── rss.go               # RSS解析
├── scraper.go           # 定时抓取任务
├── handle_*.go          # API处理函数
├── middleware_auth.go   # 认证中间件
├── internal/            # 内部模块（认证、数据库）
├── sql/                 # SQL查询与迁移文件
├── .env                 # 环境变量配置
└── sqlc.yaml            # sqlc配置文件
```

## API 速览

- 健康检查：`GET /v1/healthz`
- 用户：`POST /v1/users` 注册 ｜ `GET /v1/users` 获取当前用户
- RSS源：`POST /v1/feeds` 添加 ｜ `GET /v1/feeds` 获取全部
- 订阅：`POST /v1/feed_follows` 关注 ｜ `DELETE /v1/feed_follows/{id}` 取消关注
- 文章：`GET /v1/posts` 获取订阅文章

## 快速开始

### 前提

- Go ≥ 1.24.3
- PostgreSQL 已安装并创建数据库

### 运行

```bash
# 安装依赖
go mod download

# 生成SQL代码
sqlc generate

# 应用数据库迁移
goose -dir sql/schema postgres "postgres://postgres:123456@localhost:5432/go-rss?sslmode=disable" up

# 启动服务
go run .
```

### 环境变量（.env）

```env
PORT=8080
DB_URL=postgres://user:password@localhost:5432/go-rss?sslmode=disable
```



## 开发指南

- **新增查询**：在 `sql/query` 下编写SQL → `sqlc generate`（自动生成对应Go代码，配置文件`sqlc.yaml`）
- **修改表结构**：在 `sql/schema` 下添加迁移文件 → `goose up`



## 简易前端
浏览器打开 `http://localhost:8080` 即可访问。
![简易前端截图](image.png)



## 下一步计划

- 学习[feeds (github.com/gorilla/feeds)](https://pkg.go.dev/github.com/gorilla/feeds)和[gofeed (github.com/mmcdole/gofeed)](https://pkg.go.dev/github.com/mmcdole/gofeed)
- 开发一个RSS平台，前端用AntDesignPro及其组件，后端需要有一个公共的用户中心服务，包括rss等其他服务都共用同一个用户登陆接口。使用gRPC作为跨服务调用。

