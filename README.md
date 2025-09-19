# Go RSS 订阅器

这是一个使用Go语言开发的RSS订阅器后端服务，能够帮助用户订阅、管理和获取各类网站的RSS内容。

## 目录结构

项目采用模块化设计，主要包括以下几个部分：

```
├── main.go              # 主程序入口
├── rss.go               # RSS解析相关功能
├── scraper.go           # RSS源定时抓取功能
├── handle_*.go          # 各类API处理函数
├── json.go              # JSON响应处理
├── middleware_auth.go   # 认证中间件
├── internal/
│   ├── auth/            # 认证相关代码
│   └── db/              # 数据库操作代码（自动生成）
├── sql/
│   ├── query/           # SQL查询定义
│   └── schema/          # 数据库表结构定义
├── .env                 # 环境变量配置
└── sqlc.yaml            # sqlc配置文件
```

## 技术栈

- **Go语言**：主要开发语言
- **PostgreSQL**：数据库存储
- **sqlc**：SQL查询代码自动生成
- **goose**：数据库迁移管理
- **chi**：HTTP路由框架
- **godotenv**：环境变量管理

## 功能介绍

### 1. 用户管理
- 用户注册
- 用户认证（基于API Key）

### 2. RSS源管理
- 添加新的RSS源
- 获取所有RSS源

### 3. 订阅管理
- 关注RSS源
- 取消关注RSS源
- 获取用户已关注的RSS源

### 4. 文章管理
- 定时自动抓取RSS源内容
- 获取用户订阅的文章列表

## RSS技术流程

### RSS订阅和抓取流程

1. **添加RSS源**：用户通过API提交RSS源的名称和URL(`feeds`)
2. **关注RSS源**：用户关注感兴趣的RSS源，建立用户与RSS源之间的关联`feed_follows`
3. **定时抓取**：系统后台定时启动爬虫任务，并发抓取多个RSS源
4. **解析内容**：解析RSS XML格式，提取文章标题、链接、描述和发布时间
5. **存储文章**：将抓取到的文章信息存储到数据库`post`中
6. **获取文章**：用户可以获取自己订阅的所有RSS源的最新文章(`feed_follows`和`post`联表查询)

### 核心代码流程

1. **RSS解析** (`rss.go`)
   ```go
   // 将URL转换为RSSFeed对象
   func urlToFeed(url string) (RSSFeed, error) {
       // 发送HTTP请求获取RSS内容
       // 解析XML内容并映射到RSSFeed结构体
       // 返回解析后的RSSFeed对象
   }
   ```

2. **定时抓取** (`scraper.go`)
   ```go
   // 启动定时抓取任务
   func startScraping(query *db.Queries, concurrency int, timeBetweenRequest time.Duration) {
       // 创建定时器，按指定间隔执行抓取任务
       // 获取需要抓取的RSS源列表
       // 并发抓取每个RSS源的内容
   }
   
   // 抓取单个RSS源
   func scrapeFeed(wg *sync.WaitGroup, query *db.Queries, feed db.Feed) {
       // 标记RSS源为已抓取
       // 解析RSS内容
       // 提取文章信息并存储到数据库
   }
   ```

3. **主程序流程** (`main.go`)
   ```go
   func main() {
       // 加载环境变量
       // 连接数据库
       // 启动定时抓取任务
       // 设置API路由
       // 启动HTTP服务器
   }
   ```

## 数据库和SQL管理

### 数据库架构

项目使用PostgreSQL作为数据库，包含以下主要表：

| 表名             | 字段名          | 数据类型        | 约束                             | 描述                   |
| ---------------- | --------------- | --------------- | -------------------------------- | ---------------------- |
| **users**        |                 |                 |                                  | 用户信息表             |
|                  | id              | UUID 或 SERIAL  | PRIMARY KEY                      | 用户唯一标识           |
|                  | username        | VARCHAR(255)    | NOT NULL, UNIQUE                 | 用户名                 |
|                  | password        | VARCHAR(255)    |                                  | 密码（当前项目未使用） |
|                  | created_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 创建时间               |
|                  | updated_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 更新时间               |
|                  | api_key         | VARCHAR(255)    | NOT NULL, UNIQUE                 | API密钥，用于认证      |
| **feeds**        |                 |                 |                                  | RSS源表                |
|                  | id              | UUID 或 SERIAL  | PRIMARY KEY                      | RSS源唯一标识          |
|                  | name            | VARCHAR(255)    | NOT NULL                         | RSS源名称              |
|                  | url             | VARCHAR(500)    | NOT NULL, UNIQUE                 | RSS源URL               |
|                  | created_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 创建时间               |
|                  | updated_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 更新时间               |
|                  | user_id         | UUID 或 INTEGER | FOREIGN KEY REFERENCES users(id) | 创建者用户ID           |
|                  | last_fetched_at | TIMESTAMP       |                                  | 最后抓取时间           |
| **feed_follows** |                 |                 |                                  | 用户订阅表             |
|                  | id              | UUID 或 SERIAL  | PRIMARY KEY                      | 订阅关系唯一标识       |
|                  | user_id         | UUID 或 INTEGER | FOREIGN KEY REFERENCES users(id) | 用户ID                 |
|                  | feed_id         | UUID 或 INTEGER | FOREIGN KEY REFERENCES feeds(id) | RSS源ID                |
|                  | created_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 创建时间               |
|                  | updated_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 更新时间               |
| **posts**        |                 |                 |                                  | 文章表                 |
|                  | id              | UUID 或 SERIAL  | PRIMARY KEY                      | 文章唯一标识           |
|                  | title           | VARCHAR(500)    | NOT NULL                         | 文章标题               |
|                  | url             | VARCHAR(1000)   | NOT NULL, UNIQUE                 | 文章链接               |
|                  | description     | TEXT            |                                  | 文章描述               |
|                  | published_at    | TIMESTAMP       |                                  | 发布时间               |
|                  | feed_id         | UUID 或 INTEGER | FOREIGN KEY REFERENCES feeds(id) | 所属RSS源ID            |
|                  | created_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 创建时间               |
|                  | updated_at      | TIMESTAMP       | NOT NULL DEFAULT NOW()           | 更新时间               |

### SQL生成和管理

项目使用**sqlc**和**goose**工具来管理数据库操作：

1. **sqlc**：自动生成Go代码
   - 在`sql/query`目录下定义SQL查询
   - sqlc根据SQL查询自动生成对应的Go代码
   - 生成的代码位于`internal/db`目录
   - 配置文件：`sqlc.yaml`

   **使用方法**：
   
   ```bash
   # 安装sqlc（如果尚未安装）
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   
   # 生成代码
   sqlc generate
   ```
   
2. **goose**：数据库迁移管理
  
   - 在`sql/schema`目录下定义迁移文件
   - 文件名格式：`XXX_name.sql`，其中XXX是迁移序号
   - 每个文件包含`-- +goose Up`和`-- +goose Down`两部分
   - Up部分定义应用迁移的SQL语句
   - Down部分定义回滚迁移的SQL语句
   
   **使用方法**：
   ```bash
   # 安装goose（如果尚未安装）
   go install github.com/pressly/goose/v3/cmd/goose@latest
   
   # 应用所有迁移
   goose -dir sql/schema postgres "postgres://postgres:123456@localhost:5432/go-rss?sslmode=disable" up
   
   # 回滚上一个迁移
   goose -dir sql/schema postgres "postgres://postgres:123456@localhost:5432/go-rss?sslmode=disable" down
   ```

## API接口

### 基础接口
- `GET /v1/healthz`：检查服务是否正常运行
- `GET /v1/error`：测试错误处理

### 用户接口
- `POST /v1/users`：创建新用户
- `GET /v1/users`：获取当前用户信息（需要认证）

### RSS源接口
- `POST /v1/feeds`：添加新的RSS源（需要认证）
- `GET /v1/feeds`：获取所有RSS源

### 订阅接口
- `POST /v1/feed_follows`：关注RSS源（需要认证）
- `GET /v1/feed_follows`：获取用户关注的所有RSS源（需要认证）
- `DELETE /v1/feed_follows/{feedID}`：取消关注RSS源（需要认证）

### 文章接口
- `GET /v1/posts`：获取用户订阅的文章（需要认证）

## 认证方式

项目使用API Key进行认证：
- 用户注册时系统自动生成唯一的API Key
- API请求时需要在`Authorization`请求头中提供API Key
- 未提供或提供无效API Key的请求将被拒绝

## 环境配置

项目需要以下环境变量（配置在`.env`文件中）：
- `PORT`：服务器端口号
- `DB_URL`：PostgreSQL数据库连接字符串

## 运行项目

### 前提条件
- 已安装Go（1.24.3或更高版本）
- 已安装PostgreSQL
- 已创建数据库
- 已应用数据库迁移

### 启动服务

```bash
# 安装依赖
go mod download

# 运行项目
go run .
```

## 开发指南

### 添加新的数据库查询
1. 在`sql/query`目录下添加或修改SQL文件
2. 运行`sqlc generate`生成对应的Go代码

### 修改数据库结构
1. 在`sql/schema`目录下创建新的迁移文件
2. 运行`goose up`应用迁移

## 注意事项
- 本项目目前仅实现了后端API，没有前端界面
- 密码字段已在数据库中定义，但当前认证机制基于API Key
- RSS抓取频率可在`main.go`中的`startScraping`函数调用处调整

## 下一步改进方向
1. 添加前端界面
2. 增强用户认证和授权机制
3. 优化RSS抓取性能和错误处理
4. 添加文章全文搜索功能
5. 实现文章阅读状态跟踪
6. 引入Redis？