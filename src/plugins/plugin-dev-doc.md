
## 二、插件元信息（meta.ts）规范

```ts
export default {
  name: '插件名称',
  key: 'plugin-xxx',           // 唯一标识
  path: '/plugin-xxx/xxx',     // 路由路径
  requiresMongo: true,         // 是否依赖MongoDB
  requiresMysql: false,        // 是否依赖MySQL（如需本地缓存等）
  // 其他自定义字段
}
```

## 三、插件开发流程

### 1. 插件生成

- 推荐使用自动化脚本（如 `scripts/create-plugin.cjs`）或前端一键生成页面，自动创建插件目录和基础模板。
- 自动生成唯一 key、path，页面风格与主项目一致。

### 2. 页面开发

- 在 `index.vue` 中开发插件主页面。
- 可直接使用主项目依赖（如 Element Plus）、全局样式。
- 建议页面内容自适应填满内容区，风格统一。

### 3. 路由与导航集成

- 插件的 `meta.ts` 结构与主项目导航一致，便于自动集成。
- 路由在 `src/router/index.ts` 静态注册，避免刷新 404。

---

## 四、数据存储与交互说明

1. **敏感/核心数据（如积分、用户信息）**
   - 通过后端 API 进行传递和获取，前端不直接操作数据库，保证数据安全和一致性。
   - 示例：`axios.get('/api/user/score')` 获取积分，`axios.post('/api/user/update')` 更新信息。

2. **本地数据（如日志、行为分析、缓存等）**
   - 插件可直接操作本地 MySQL/MongoDB 存储高频、海量、非敏感数据。
   - 适合场景：访问日志、行为分析、插件临时数据、离线缓存等。
   - 本地数据可定期同步到后端，或由后端定期拉取分析。

3. **数据同步与安全**
   - 插件如需将本地数据同步到后端，需通过安全 API 进行。
   - 敏感数据严禁仅存本地，必须通过后端 API 统一管理。

---

### 插件本地数据库操作示例

#### 1. 操作本地 MongoDB（Node.js 端）

```js
const { MongoClient } = require('mongodb');
const client = new MongoClient('mongodb://localhost:27017');
await client.connect();
const db = client.db('local_plugin');
await db.collection('plugin-xxx_logs').insertOne({ ... });
```

#### 2. 操作本地 MySQL（Node.js 端）

```js
const mysql = require('mysql2/promise');
const conn = await mysql.createConnection({host: 'localhost', user: 'root', database: 'local_plugin'});
await conn.execute('INSERT INTO plugin_xxx_cache (data) VALUES (?)', [data]);
```

---

### 插件开发注意事项

- 插件需区分敏感数据和本地可缓存数据，合理选择存储方式。
- 本地数据库仅用于提升体验和处理高频数据，敏感数据必须通过后端 API 统一管理。
- 插件如需本地数据库支持，需在 `meta.ts` 中声明依赖（如 `requiresMongo: true`、`requiresMysql: true`）。
- 本地数据同步到后端时，需做好数据校验和安全处理。

---

### 场景举例

- **积分**：通过后端 API 获取和提交，保证安全。
- **访问日志**：本地 MongoDB 记录，定期同步到后端分析。
- **插件缓存**：本地 MySQL 存储，提升页面响应速度。

---

## 五、数据存储与日志分析（MongoDB）

#### a. 依赖声明

- 在 `meta.ts` 中设置 `requiresMongo: true`，插件安装时自动检测和初始化 MongoDB。

#### b. 集合命名规范

- 每个插件独立集合，命名为 `${pluginKey}_logs`，如 `plugin-wailki_logs`。

#### c. 数据结构建议

```json
{
  "plugin": "plugin-wailki",
  "user": "userId",
  "url": "https://xxx",
  "result": "success",
  "score": 1,
  "timestamp": "2024-06-01T12:00:00Z"
}
```

#### d. 日志与积分

- 插件页面通过后端 API 记录访问日志、积分等数据。
- 后端统一管理 MySQL 和 MongoDB 连接，按需分流数据。

---

## 六、插件安装与数据库自动化

#### a. 安装流程

1. 读取插件 `meta.ts`，判断是否需要 MongoDB/MySQL。
2. 检查本地 MongoDB/MySQL 是否可用，不可用则自动安装或提示用户。
3. 自动创建数据库和集合/表（如未存在）。
4. 完成插件安装。

#### b. 自动化脚本示例

```js
// 检查并创建集合
const { MongoClient } = require('mongodb');
async function ensureMongoCollection(pluginKey) {
  const client = new MongoClient('mongodb://localhost:27017');
  await client.connect();
  const db = client.db('vite_pluginend');
  const colName = `${pluginKey}_logs`;
  const collections = await db.listCollections().toArray();
  if (!collections.find(c => c.name === colName)) {
    await db.createCollection(colName);
    console.log(`已为插件 ${pluginKey} 创建集合 ${colName}`);
  }
  await client.close();
}
```

#### c. Docker 环境推荐

- 推荐开发环境用 Docker 管理 MongoDB，自动拉取并启动：
  ```bash
  docker run -d --name vite-plugin-mongo -p 27017:27017 mongo:latest
  ```

---

## 七、日志分析与接口

- 后端提供如 `/api/plugin-xxx/logs`、`/api/plugin-xxx/score` 等接口，支持日志查询、积分统计。
- 支持按插件、用户、时间等维度分析。

---

## 八、插件开发注意事项（汇总）

- 插件页面风格需与主项目统一，内容区自适应填满，间距合适。
- 插件数据需与主项目隔离，避免数据冲突。
- 插件如需数据库，务必在 `meta.ts` 中声明依赖。
- 插件卸载时可选自动清理对应集合（谨慎操作）。

---

## 九、常见问题与解决

- **MongoDB/MySQL 未安装/未启动**：请按提示安装或启动数据库服务，或使用 Docker 启动。
- **端口冲突**：如 27017/3306 被占用，请修改数据库配置或容器端口。
- **数据分析需求变化**：建议日志结构预留扩展字段，便于后续分析。

---

## 十、示例：外链插件（plugin-wailki）开发流程

1. 通过脚本或页面一键生成插件骨架。
2. 在 `meta.ts` 中声明 `requiresMongo: true`。
3. 在 `index.vue` 开发页面，实现外链拼接、验证、访问、日志记录、积分统计等功能。
4. 后端自动为插件创建 MongoDB 集合，提供日志与积分接口。
5. 插件页面通过 API 读写数据，实现完整业务闭环。

---

如需详细代码模板、API 设计、数据库脚本等，可随时补充！