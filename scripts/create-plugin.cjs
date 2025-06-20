const fs = require('fs');
const path = require('path');
const pinyin = require('pinyin');

function toCamel(str) {
  return str.replace(/-([a-z])/g, (g) => g[1].toUpperCase())
}

function toPinyin(str) {
  // 兼容非字符串输入
  if (!str || typeof str !== 'string') return '';
  // pinyin v2.x 返回数组
  return pinyin(str, { style: pinyin.STYLE_NORMAL })
    .flat()
    .join('-')
    .replace(/[^a-zA-Z0-9-]/g, '')
    .toLowerCase()
}

function genMetaTs(options) {
  const { name, key, version = '1.0.0', description = '', withTests = false, withDocs = false } = options
  const pages = options.pages || [{ title: '首页' }]

  // 自动生成 key 和 path
  pages.forEach(p => {
    if (!p.key) {
      p.key = toPinyin(p.title)
    }
    if (!p.path) {
      p.path = `/plugin-${key}/${p.key}`
    }
  })

  const mainNav = {
    key: `plugin-${key}`,
    title: name,
    icon: 'Document',
    path: `/plugin-${key}/${pages[0].key}`,
    permission: `plugin_${key}_access`
  }

  const subNav = pages.map(p => ({
    key: p.key,
    title: p.title,
    icon: 'Document',
    path: `/plugin-${key}/${p.key}`,
    permission: `plugin_${key}_${p.key}_access`
  }))

  return `import { definePlugin } from '../../plugins/types'

export default definePlugin({
  name: '${name}',
  version: '${version}',
  description: '${description}',
  mainNav: ${JSON.stringify(mainNav, null, 2)},
  subNav: ${JSON.stringify(subNav, null, 2)},
  pages: [
    ${pages.map(p => `{
      key: '${p.key}',
      title: '${p.title}',
      path: '${p.path}',
      permission: '${p.permission || `plugin_${key}_${p.key}_access`}',
      component: () => import('./pages/${p.key}.vue')
    }`).join(',\n    ')}
  ]
})
`
}

function genIndexTs(options) {
  const { key } = options
  return `import meta from './meta'
export default meta
`
}

function genPageTemplate(page, options) {
  const { key } = options
  return `<template>
  <div class="plugin-page">
    <div class="page-header">
      <h2>${page.title}</h2>
      <p class="description">这是 ${page.title} 页面的描述文本</p>
    </div>

    <div class="page-content">
      <el-card class="content-section">
        <template #header>
          <div class="card-header">
            <span>基础信息</span>
            <el-button type="primary">操作按钮</el-button>
          </div>
        </template>
        
        <el-form label-width="100px">
          <el-form-item label="标题">
            <el-input v-model="form.title" placeholder="请输入标题" />
          </el-form-item>
          <el-form-item label="描述">
            <el-input 
              v-model="form.description" 
              type="textarea" 
              :rows="3"
              placeholder="请输入描述"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSubmit">提交</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const form = ref({
  title: '',
  description: ''
})

const handleSubmit = () => {
  ElMessage.success('提交成功')
}

const handleReset = () => {
  form.value = {
    title: '',
    description: ''
  }
}
</script>

<style scoped>
.plugin-page {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.description {
  color: #666;
  margin: 0;
}

.page-content {
  display: grid;
  gap: 24px;
}

.content-section {
  background: #fff;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-size: 16px;
  font-weight: 500;
}
</style>
`
}

function genReadme(options) {
  const { name, key, description, version } = options
  return `# ${name}

${description}

## 版本

${version}

## 安装

1. 解压插件包到 plugins 目录
2. 运行以下命令安装依赖：

\`\`\`bash
cd plugin-${key}
npm install
\`\`\`

## 使用

插件提供以下功能页面：

${options.pages.map(p => `- ${p.title}`).join('\n')}

## 开发

\`\`\`bash
# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
\`\`\`

## 许可证

MIT
`
}

function createPlugin(options) {
  if (!options.key) {
    console.error('缺少插件标识符')
    process.exit(1)
  }

  const { key } = options
  const dir = path.join(__dirname, '../src/plugins', `plugin-${key}`)

  try {
    // 检查插件是否已存在
    if (fs.existsSync(dir)) {
      console.error('插件已存在')
      process.exit(1)
    }

    // 创建目录结构
    fs.mkdirSync(dir, { recursive: true })
    fs.mkdirSync(path.join(dir, 'pages'))
    fs.mkdirSync(path.join(dir, 'components'))
    fs.mkdirSync(path.join(dir, 'assets'))

    // 生成文件
    fs.writeFileSync(path.join(dir, 'meta.ts'), genMetaTs(options))
    fs.writeFileSync(path.join(dir, 'index.ts'), genIndexTs(options))
    
    // 生成页面
    options.pages.forEach(page => {
      fs.writeFileSync(
        path.join(dir, 'pages', `${page.key}.vue`),
        genPageTemplate(page, options)
      )
    })

    // 生成 README
    if (options.withDocs) {
      fs.writeFileSync(path.join(dir, 'README.md'), genReadme(options))
    }

    // 生成测试文件
    if (options.withTests) {
      fs.mkdirSync(path.join(dir, '__tests__'))
      options.pages.forEach(page => {
        fs.writeFileSync(
          path.join(dir, '__tests__', `${page.key}.test.ts`),
          `import { describe, it, expect } from 'vitest'

describe('${page.title}', () => {
  it('should work', () => {
    expect(true).toBe(true)
  })
})
`
        )
      })
    }

    console.log(`插件 ${key} 创建成功！`)
    process.exit(0)
  } catch (error) {
    console.error('创建插件失败:', error.message)
    process.exit(1)
  }
}

if (require.main === module) {
  const options = JSON.parse(process.argv[2])
  createPlugin(options)
}

module.exports = {
  createPlugin,
  genMetaTs,
  genPageTemplate
} 
