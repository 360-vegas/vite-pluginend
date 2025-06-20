# 插件依赖管理系统

## 概述

插件依赖管理系统是一个全面的解决方案，用于检查、验证和设置插件运行所需的各种依赖。该系统确保插件在安装前满足所有必要的环境要求。

## 功能特性

### 1. 多类型依赖检查

- **数据库依赖**: 检查MongoDB、MySQL、PostgreSQL等数据库连接
- **服务依赖**: 验证外部服务的可用性（API、缓存、队列等）
- **环境变量**: 检查必需的环境变量配置
- **基础依赖**: 验证Go模块、npm包等依赖包
- **权限检查**: 确认用户权限满足插件需求

### 2. 智能数据库管理

- **自动检测**: 检测现有数据库和集合
- **配置建议**: 提供数据库配置建议
- **一键创建**: 支持自动创建数据库和集合
- **索引管理**: 自动创建必要的数据库索引

### 3. 用户友好的界面

- **可视化检查结果**: 清晰显示依赖状态
- **交互式配置**: 提供图形化配置界面
- **实时反馈**: 即时显示检查和配置结果
- **建议指导**: 提供解决问题的具体建议

## 系统架构

### 后端组件

#### 1. 数据模型 (`models/plugin_dependency.go`)

```go
type PluginDependency struct {
    PluginKey     string                  
    Dependencies  []Dependency            
    Database      *DatabaseRequirement    
    Services      []ServiceRequirement    
    Environment   []EnvironmentVariable   
    Permissions   []PermissionRequirement 
}
```

#### 2. 依赖服务 (`services/dependency_service.go`)

- `CheckPluginDependencies()`: 检查插件所有依赖
- `SetupDatabase()`: 配置数据库环境
- `checkMongoDBRequirement()`: MongoDB专用检查
- `checkServiceRequirements()`: 服务可用性检查
- `checkEnvironmentVariables()`: 环境变量验证

#### 3. API端点 (`handlers/plugin_handler.go`)

- `GET /api/plugins/{key}/dependencies/check`: 检查依赖
- `POST /api/plugins/{key}/dependencies/setup`: 配置数据库

### 前端组件

#### 1. 依赖检查器 (`DependencyChecker.vue`)

- 依赖状态可视化
- 数据库配置界面
- 实时检查结果更新
- 用户交互式配置

#### 2. 市场插件页面集成

- 安装前依赖检查
- 配置向导集成
- 错误处理和用户指导

## 使用流程

### 1. 插件依赖定义

在插件目录创建 `dependencies.json` 文件：

```json
{
  "plugin_key": "example-plugin",
  "database": {
    "type": "mongodb",
    "database_name": "plugin_example",
    "required": true,
    "collections": [
      {
        "name": "example_data",
        "description": "示例数据表",
        "indexes": [...]
      }
    ]
  },
  "environment": [
    {
      "name": "API_KEY",
      "required": true,
      "type": "string",
      "description": "API密钥"
    }
  ]
}
```

### 2. 安装时检查

```javascript
// 前端调用
const result = await appsApi.checkPluginDependencies('example-plugin');

if (result.data.can_install) {
  // 可以安装
  await installPlugin();
} else {
  // 显示依赖问题和解决建议
  showDependencyChecker();
}
```

### 3. 数据库配置

```javascript
// 配置数据库
const config = {
  suggested_database_name: 'plugin_example',
  create_new_database: true
};

await appsApi.setupPluginDatabase('example-plugin', config);
```

## 配置示例

### MongoDB依赖配置

```json
{
  "database": {
    "type": "mongodb",
    "database_name": "plugin_external_links",
    "required": true,
    "collections": [
      {
        "name": "external_links",
        "description": "外链数据表",
        "indexes": [
          {
            "name": "url_index",
            "fields": { "url": 1 },
            "unique": true
          }
        ]
      }
    ]
  }
}
```

### 环境变量配置

```json
{
  "environment": [
    {
      "name": "MONGODB_URI",
      "required": true,
      "type": "url",
      "description": "MongoDB连接字符串",
      "default": "mongodb://localhost:27017"
    }
  ]
}
```

### 服务依赖配置

```json
{
  "services": [
    {
      "name": "redis",
      "type": "cache",
      "host": "localhost",
      "port": 6379,
      "required": false,
      "description": "Redis缓存服务"
    }
  ]
}
```

## API接口

### 检查依赖

**请求**: `GET /api/plugins/{key}/dependencies/check`

**响应**:
```json
{
  "success": true,
  "data": {
    "plugin_key": "example-plugin",
    "overall_status": "warning",
    "can_install": true,
    "requires_setup": true,
    "database": {
      "status": "setup_required",
      "message": "数据库不存在，需要创建",
      "setup_options": {
        "suggested_database_name": "plugin_example",
        "create_new_database": true
      }
    },
    "suggestions": [
      "需要设置数据库连接和创建必要的集合"
    ]
  }
}
```

### 配置数据库

**请求**: `POST /api/plugins/{key}/dependencies/setup`

**请求体**:
```json
{
  "suggested_database_name": "plugin_example",
  "create_new_database": true,
  "use_existing_database": false
}
```

**响应**:
```json
{
  "success": true,
  "message": "数据库设置成功"
}
```

## 状态类型

### 整体状态
- `success`: 所有依赖满足
- `warning`: 存在非关键问题
- `error`: 存在关键问题，无法安装

### 依赖状态
- `available`: 依赖可用
- `missing`: 依赖缺失
- `setup_required`: 需要配置
- `version_mismatch`: 版本不匹配
- `unreachable`: 服务不可达

## 扩展指南

### 添加新的依赖类型

1. 在 `models/plugin_dependency.go` 中定义新的依赖结构
2. 在 `dependency_service.go` 中实现检查逻辑
3. 在前端 `DependencyChecker.vue` 中添加UI显示

### 支持新的数据库类型

1. 在 `checkDatabaseRequirement()` 方法中添加新的 case
2. 实现特定数据库的检查逻辑
3. 在前端添加相应的配置界面

## 最佳实践

### 插件开发者

1. **详细的依赖说明**: 在配置文件中提供清晰的依赖描述
2. **合理的默认值**: 为环境变量提供合理的默认值
3. **向后兼容**: 考虑插件更新时的兼容性
4. **错误处理**: 提供详细的错误信息和解决建议

### 系统管理员

1. **预先配置**: 在插件安装前确保基础环境就绪
2. **定期检查**: 定期验证插件依赖的健康状态
3. **备份策略**: 为数据库依赖制定备份策略
4. **监控告警**: 监控服务依赖的可用性

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 `MONGODB_URI` 环境变量
   - 确认数据库服务运行状态
   - 验证网络连接和防火墙设置

2. **权限不足**
   - 检查数据库用户权限
   - 确认文件系统权限
   - 验证API访问权限

3. **环境变量缺失**
   - 检查 `.env` 文件配置
   - 确认环境变量名称拼写
   - 验证环境变量值格式

### 日志分析

系统提供详细的日志记录，包括：
- 依赖检查过程
- 数据库操作记录
- 错误详情和堆栈跟踪

## 未来规划

### 计划功能

1. **插件间依赖**: 支持插件之间的依赖关系
2. **版本管理**: 依赖版本冲突检测和解决
3. **自动修复**: 自动修复常见的依赖问题
4. **性能监控**: 依赖性能监控和优化建议
5. **云服务集成**: 支持云数据库和服务的自动配置

### 技术改进

1. **并行检查**: 并行检查多个依赖以提高性能
2. **缓存机制**: 缓存检查结果以减少重复检查
3. **实时更新**: 实时监控依赖状态变化
4. **插件市场**: 集成到插件市场的依赖信息展示

## 总结

插件依赖管理系统提供了一个全面、用户友好的解决方案来管理插件依赖。通过自动化的检查和配置流程，大大简化了插件的安装和维护工作，提高了系统的稳定性和用户体验。

系统的模块化设计使其易于扩展和维护，能够适应未来不断变化的技术需求。通过持续的改进和优化，该系统将成为插件生态系统的重要基础设施。 