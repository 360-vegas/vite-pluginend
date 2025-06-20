# Vite Pluginend Backend

这是 Vite Pluginend 项目的后端服务，使用 Golang 和 Gin 框架开发。

## 环境要求

- Go 1.21 或更高版本
- MongoDB 4.4 或更高版本
- Redis（可选，用于缓存）

## 项目结构

```
backend/
├── cmd/
│   └── server/          # 主程序入口
├── internal/
│   ├── api/            # API 处理器
│   ├── models/         # 数据模型
│   ├── services/       # 业务逻辑
│   └── middleware/     # 中间件
├── pkg/
│   ├── database/       # 数据库连接
│   └── utils/          # 工具函数
└── go.mod              # Go 模块文件
```

## 快速开始

1. 克隆项目并进入后端目录：
   ```bash
   cd backend
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 配置环境变量：
   - 复制 `.env.example` 为 `.env`
   - 根据需要修改配置

4. 启动服务：
   ```bash
   go run cmd/server/main.go
   ```

## API 文档

### 插件管理

- `POST /api/create-plugin` - 创建新插件
- `GET /api/plugins/:key` - 获取插件信息
- `GET /api/plugins` - 获取所有插件
- `GET /api/pack-plugin` - 打包插件

## 开发指南

### 添加新的 API 端点

1. 在 `internal/api/handlers` 中创建新的处理器
2. 在 `internal/services` 中实现业务逻辑
3. 在 `cmd/server/main.go` 中注册路由

### 数据库操作

使用 MongoDB 作为主数据库，所有数据库操作都在 `pkg/database` 中定义。

### 中间件

常用的中间件包括：
- CORS
- 认证
- 日志
- 错误处理

## 部署

1. 构建二进制文件：
   ```bash
   go build -o server cmd/server/main.go
   ```

2. 运行服务：
   ```bash
   ./server
   ```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| PORT | 服务端口 | 3001 |
| GIN_MODE | Gin 运行模式 | debug |
| MONGODB_URI | MongoDB 连接 URI | mongodb://localhost:27017 |
| MONGODB_DB | MongoDB 数据库名 | vite_pluginend |
| JWT_SECRET | JWT 密钥 | - |
| JWT_EXPIRE | JWT 过期时间 | 24h |
| UPLOAD_DIR | 上传文件目录 | uploads |
| MAX_UPLOAD_SIZE | 最大上传大小 | 10485760 |

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT 