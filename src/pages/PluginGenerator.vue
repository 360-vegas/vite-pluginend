<template>
  <div class="plugin-generator">
    <div class="generator-header">
      <h2>插件生成器</h2>
      <div class="header-actions">
        <el-button type="primary" @click="resetGenerator" v-if="showResult">
          <el-icon><Refresh /></el-icon>
          重新生成
        </el-button>
      </div>
    </div>

    <!-- 表单区域 -->
    <el-card class="form-card" shadow="never" v-if="!showResult">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><DocumentAdd /></el-icon>
          <span>创建新插件</span>
        </div>
      </template>
      
      <el-form :model="formData" label-width="120px" class="generator-form">
        <div class="form-grid">
          <div class="form-section">
            <h3>基本信息</h3>
            <el-form-item label="插件名称" required>
              <el-input 
                v-model="formData.name" 
                placeholder="请输入插件名称"
                clearable
                @input="onNameInput"
              />
            </el-form-item>
            
            <el-form-item label="插件标识符" required>
              <el-input 
                v-model="formData.key" 
                placeholder="plugin-example"
                clearable
                @input="onKeyInput"
              >
                <template #prepend>plugin-</template>
              </el-input>
              <div class="form-tip">
                标识符用于生成文件夹名称，只能包含小写字母、数字和连字符
              </div>
            </el-form-item>
            
            <el-form-item label="插件描述">
              <el-input 
                v-model="formData.description" 
                type="textarea"
                :rows="3"
                placeholder="请输入插件描述"
              />
            </el-form-item>
            
            <el-form-item label="作者">
              <el-input 
                v-model="formData.author" 
                placeholder="请输入作者名称"
                clearable
              />
            </el-form-item>
            
            <el-form-item label="版本">
              <el-input 
                v-model="formData.version" 
                placeholder="1.0.0"
                clearable
              />
            </el-form-item>
            
            <el-form-item label="分类">
              <el-select v-model="formData.category" placeholder="选择插件分类">
                <el-option label="工具类" value="tools" />
                <el-option label="娱乐类" value="entertainment" />
                <el-option label="效率类" value="productivity" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
          </div>
          
          <div class="form-section">
            <h3>页面配置</h3>
            <el-form-item label="插件页面">
              <div class="pages-manager">
                <div v-for="(page, index) in formData.pages" :key="index" class="page-item">
                  <el-input 
                    v-model="page.title" 
                    placeholder="页面标题"
                    style="width: 45%"
                  />
                  <el-input 
                    v-model="page.key" 
                    placeholder="页面路径"
                    style="width: 45%"
                  />
                  <el-button 
                    type="danger" 
                    size="small" 
                    @click="removePage(index)"
                    :disabled="formData.pages.length <= 1"
                  >
                    删除
                  </el-button>
                </div>
                <el-button type="primary" plain @click="addPage" size="small">
                  <el-icon><DocumentAdd /></el-icon>
                  添加页面
                </el-button>
              </div>
            </el-form-item>
            
            <el-form-item label="功能特性">
              <el-checkbox-group v-model="formData.features">
                <el-checkbox label="api">API 接口</el-checkbox>
                <el-checkbox label="storage">本地存储</el-checkbox>
                <el-checkbox label="router">路由导航</el-checkbox>
                <el-checkbox label="components">自定义组件</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            
            <el-form-item label="开发选项">
              <el-checkbox-group>
                <el-checkbox v-model="formData.withTests">包含测试文件</el-checkbox>
                <el-checkbox v-model="formData.withDocs">包含文档模板</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </div>
        </div>
        
        <div class="form-actions">
          <el-button type="primary" size="large" @click="generatePlugin" :loading="generating">
            <el-icon><Tools /></el-icon>
            {{ generating ? '生成中...' : '生成插件' }}
          </el-button>
        </div>
      </el-form>
    </el-card>

    <!-- 开发指南 -->
    <el-card class="guide-card" shadow="never" v-if="!showResult">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Guide /></el-icon>
          <span>开发指南</span>
        </div>
      </template>
      
      <div class="guide-content">
        <el-steps :active="4" direction="vertical" finish-status="success">
          <el-step title="填写插件基本信息" description="包括名称、描述、作者等信息" />
          <el-step title="配置页面信息" description="设置页面名称和路径" />
          <el-step title="选择功能特性" description="根据需要选择要包含的功能" />
          <el-step title="生成插件代码" description="自动生成完整的插件文件结构" />
          <el-step title="开发和测试" description="在生成的代码基础上进行开发" />
        </el-steps>
      </div>
    </el-card>

    <!-- 生成结果 -->
    <div v-if="showResult" class="result-section">
      <el-card shadow="never">
        <template #header>
          <div class="card-header success">
            <el-icon class="header-icon"><SuccessFilled /></el-icon>
            <span>插件生成成功！</span>
          </div>
        </template>
        
        <div class="result-content">
          <el-alert
            title="插件已成功生成"
            :description="`插件 ${formData.name} 已生成完成，文件已保存到 src/plugins/${generatedPluginKey}/ 目录`"
            type="success"
            :closable="false"
            show-icon
          />
          
          <div class="file-structure">
            <h4>生成的文件结构：</h4>
            <el-tree
              :data="fileStructure"
              :props="{ children: 'children', label: 'name' }"
              default-expand-all
              class="file-tree"
            >
              <template #default="{ node, data }">
                <span class="file-node">
                  <el-icon v-if="data.type === 'folder'"><Folder /></el-icon>
                  <el-icon v-else><Document /></el-icon>
                  {{ data.name }}
                </span>
              </template>
            </el-tree>
          </div>
          
          <div class="next-steps">
            <h4>下一步操作：</h4>
            <el-steps :active="0" direction="horizontal">
              <el-step title="开发插件" description="编辑生成的文件，实现具体功能" />
              <el-step title="测试插件" description="在开发环境中测试插件功能" />
              <el-step title="打包发布" description="使用插件打包功能生成分发包" />
            </el-steps>
          </div>
          
          <div class="action-buttons">
            <el-button type="primary" @click="openPluginFolder">
              <el-icon><FolderOpened /></el-icon>
              查看文件
            </el-button>
            <el-button @click="goToPackage">
              <el-icon><Box /></el-icon>
              去打包
            </el-button>
            <el-button @click="resetGenerator">
              <el-icon><Refresh /></el-icon>
              重新生成
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 进度对话框 -->
    <el-dialog
      v-model="generating"
      title="正在生成插件..."
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <div class="progress-content">
        <el-progress :percentage="progress" :stroke-width="6" />
        <p class="progress-text">{{ currentStep }}</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  DocumentAdd, 
  Tools, 
  Guide, 
  SuccessFilled, 
  Refresh,
  Folder,
  Document,
  FolderOpened,
  Box
} from '@element-plus/icons-vue'
import { appsApi } from '../api'

const router = useRouter()

// 表单数据
const formData = ref({
  name: '',
  key: '',
  version: '1.0.0',
  description: '',
  pages: [{ title: '首页', key: 'index' }],
  withTests: true,
  withDocs: true,
  author: '',
  category: 'tools',
  pageName: '',
  pagePath: '/plugin-demo',
  features: ['api', 'storage', 'router', 'components']
})

// 状态
const generating = ref(false)
const packaging = ref(false)
const generated = ref(false)
const packaged = ref(false)
const progressVisible = ref(false)
const progress = ref(0)
const progressStatus = ref<'success' | 'exception'>('success')
const progressMessage = ref('')
const progressTitle = ref('生成进度')
const showResult = ref(false)
const generatedPluginKey = ref('')
const currentStep = ref('')

// 文件结构
const fileStructure = ref([
  {
    name: 'plugin-demo',
    type: 'folder',
    children: [
      { name: 'index.ts', type: 'file' },
      { name: 'meta.ts', type: 'file' },
      { name: 'pages', type: 'folder', children: [
        { name: 'index.vue', type: 'file' }
      ]},
      { name: 'components', type: 'folder', children: [] },
      { name: 'assets', type: 'folder', children: [] }
    ]
  }
])

// 计算预览URL
const previewUrl = computed(() => {
  if (!formData.value.key) return ''
  return `http://localhost:3000/plugin-${formData.value.key}`
})

// 监听插件名称输入，自动生成标识符
function onNameInput(value: string) {
  // 只在key为空或者key是由name自动生成的情况下才自动更新
  if (!formData.value.key || isAutoGeneratedKey.value) {
    const newKey = value
      .toLowerCase()
      .replace(/\s+/g, '-')
      .replace(/[^a-z0-9-]/g, '')
      .replace(/^-+|-+$/g, '') // 移除开头和结尾的连字符
    formData.value.key = newKey
    isAutoGeneratedKey.value = true
  }
}

// 监听标识符输入，确保格式正确
function onKeyInput(value: string) {
  const cleanValue = value
    .toLowerCase()
    .replace(/[^a-z0-9-]/g, '')
    .replace(/^-+|-+$/g, '') // 移除开头和结尾的连字符
  formData.value.key = cleanValue
  // 用户手动修改了key，标记为非自动生成
  isAutoGeneratedKey.value = false
}

// 标记key是否为自动生成
const isAutoGeneratedKey = ref(true)

// 添加页面
function addPage() {
  const title = `页面 ${formData.value.pages.length + 1}`
  const key = `page-${formData.value.pages.length + 1}`
  formData.value.pages.push({ title, key })
}

// 删除页面
function removePage(index: number) {
  formData.value.pages.splice(index, 1)
}

// 生成插件
async function generatePlugin() {
  if (!validateForm()) return

  generating.value = true
  progressVisible.value = true
  progressTitle.value = '生成进度'
  progress.value = 0
  progressStatus.value = 'success'
  progressMessage.value = '准备生成插件...'

  try {
    // 1. 创建插件目录
    await updateProgress(20, '创建插件目录...')
    await createPluginDirectory()

    // 2. 生成插件文件
    await updateProgress(40, '生成插件文件...')
    await generatePluginFiles()

    // 3. 生成页面文件
    await updateProgress(60, '生成页面文件...')
    await generatePageFiles()

    // 4. 生成其他文件
    await updateProgress(80, '生成附加文件...')
    await generateAdditionalFiles()

    // 完成
    await updateProgress(100, '插件生成完成！')
    generated.value = true
    showResult.value = true
    generatedPluginKey.value = `plugin-${formData.value.key}`
    progressVisible.value = false
    showMessage('插件生成成功！', 'success')
  } catch (error) {
    progressStatus.value = 'exception'
    progressMessage.value = `生成失败: ${error}`
    showMessage(`生成插件失败: ${error}`, 'error')
  } finally {
    generating.value = false
  }
}

// 打包插件
async function packagePlugin() {
  packaging.value = true
  progressVisible.value = true
  progressTitle.value = '打包进度'
  progress.value = 0
  progressStatus.value = 'success'
  progressMessage.value = '准备打包...'

  try {
    // 1. 构建插件
    await updateProgress(30, '构建插件...')
    await buildPlugin()

    // 2. 打包文件
    await updateProgress(60, '打包文件...')
    await createPackage()

    // 3. 生成安装脚本
    await updateProgress(90, '生成安装脚本...')
    await generateInstallScript()

    // 完成
    await updateProgress(100, '打包完成！')
    packaged.value = true
    showMessage('插件打包成功！', 'success')
  } catch (error) {
    progressStatus.value = 'exception'
    progressMessage.value = `打包失败: ${error}`
    showMessage(`打包插件失败: ${error}`, 'error')
  } finally {
    packaging.value = false
  }
}

// 更新进度
async function updateProgress(value: number, message: string) {
  progress.value = value
  progressMessage.value = message
  currentStep.value = message
  await new Promise(resolve => setTimeout(resolve, 500))
}

// 创建插件目录
async function createPluginDirectory() {
  const response = await appsApi.generatePlugin({
    ...formData.value,
    pages: formData.value.pages.map(page => ({
      title: page.title,
      key: page.key,
      type: 'list'
    }))
  })
  console.log('API 响应:', response)
  // 后端返回格式：{success: true, message: '插件生成成功', data: {...}}
  if (!(response as any).success) {
    throw new Error((response as any).error || (response as any).message || '创建插件目录失败')
  }
}

// 生成插件文件
async function generatePluginFiles() {
  // 这个步骤已经在 createPluginDirectory 中完成
  return Promise.resolve()
}

// 生成页面文件
async function generatePageFiles() {
  // 这个步骤已经在 createPluginDirectory 中完成
  return Promise.resolve()
}

// 生成附加文件
async function generateAdditionalFiles() {
  // 这个步骤已经在 createPluginDirectory 中完成
  return Promise.resolve()
}

// 构建插件
async function buildPlugin() {
  const { data } = await appsApi.packagePlugin(formData.value.key)
  if (!data.success) {
    throw new Error(data.error || '构建插件失败')
  }
}

// 创建插件包
async function createPackage() {
  // 这个步骤已经在 buildPlugin 中完成
  return Promise.resolve()
}

// 生成安装脚本
async function generateInstallScript() {
  // 这个步骤已经在 buildPlugin 中完成
  return Promise.resolve()
}

// 前往应用市场
function goToMarket() {
  router.push('/app-market')
  progressVisible.value = false
}

// 重置生成器
function resetGenerator() {
  // 重置表单数据
  formData.value = {
    name: '',
    key: '',
    version: '1.0.0',
    description: '',
    pages: [{ title: '首页', key: 'index' }],
    withTests: true,
    withDocs: true,
    author: '',
    category: 'tools',
    pageName: '',
    pagePath: '/plugin-demo',
    features: ['api', 'storage', 'router', 'components']
  }
  
  // 重置状态
  generating.value = false
  packaging.value = false
  generated.value = false
  packaged.value = false
  progressVisible.value = false
  progress.value = 0
  progressStatus.value = 'success'
  progressMessage.value = ''
  progressTitle.value = '生成进度'
  showResult.value = false
  generatedPluginKey.value = ''
  currentStep.value = ''
  isAutoGeneratedKey.value = true
  
  showMessage('生成器已重置', 'success')
}

// 消息提示函数
function showMessage(message: string, type: 'success' | 'error' | 'warning' = 'success') {
  // 使用浏览器原生通知或者简单的alert
  if (type === 'error') {
    alert(`错误: ${message}`)
  } else {
    console.log(`${type}: ${message}`)
  }
}

// 表单验证
function validateForm() {
  if (!formData.value.name.trim()) {
    showMessage('请输入插件名称', 'error')
    return false
  }
  if (!formData.value.key.trim()) {
    showMessage('请输入插件标识符', 'error')
    return false
  }
  if (!/^[a-z0-9-]+$/.test(formData.value.key)) {
    showMessage('插件标识符只能包含小写字母、数字和连字符', 'error')
    return false
  }
  if (!formData.value.author.trim()) {
    showMessage('请输入作者名称', 'error')
    return false
  }
  if (!formData.value.description.trim()) {
    showMessage('请输入插件描述', 'error')
    return false
  }
  if (!formData.value.pages.length) {
    showMessage('请至少添加一个页面', 'error')
    return false
  }
  if (formData.value.pages.some(page => !page.title.trim())) {
    showMessage('页面标题不能为空', 'error')
    return false
  }
  if (formData.value.pages.some(page => !page.key.trim())) {
    showMessage('页面路径不能为空', 'error')
    return false
  }
  // 检查页面key是否重复
  const pageKeys = formData.value.pages.map(page => page.key)
  const uniqueKeys = [...new Set(pageKeys)]
  if (pageKeys.length !== uniqueKeys.length) {
    showMessage('页面路径不能重复', 'error')
    return false
  }
  return true
}

// 打开插件文件夹
function openPluginFolder() {
  const pluginPath = `src/plugins/${generatedPluginKey.value}`
  showMessage(`插件文件位于: ${pluginPath}`, 'success')
  // 在实际项目中，这里可以调用 API 或使用 Electron API 来打开文件夹
}

// 跳转到打包页面
function goToPackage() {
  router.push('/plugin-package')
}
</script>

<style lang="scss" scoped>
.plugin-generator {
  padding: 20px;
  height: 100%;
  overflow: auto;

  .generator-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      color: #303133;
    }
  }

  .form-card, .guide-card {
    margin-bottom: 20px;
    
    .card-header {
      display: flex;
      align-items: center;
      
      .header-icon {
        font-size: 18px;
        margin-right: 8px;
        color: #409eff;
      }
      
      &.success {
        color: #67c23a;
        
        .header-icon {
          color: #67c23a;
        }
      }
    }
  }

  .generator-form {
    .form-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 40px;
      margin-bottom: 30px;
    }
    
    .form-section {
      h3 {
        margin: 0 0 20px 0;
        color: #303133;
        font-size: 16px;
        border-bottom: 1px solid #ebeef5;
        padding-bottom: 10px;
      }
    }
    
    .form-tip {
      font-size: 12px;
      color: #909399;
      margin-top: 5px;
      line-height: 1.4;
    }
    
    .pages-manager {
      .page-item {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-bottom: 10px;
        
        .el-input {
          flex: 1;
        }
        
        .el-button {
          flex-shrink: 0;
        }
      }
    }
    
    .form-actions {
      text-align: center;
      padding-top: 20px;
      border-top: 1px solid #ebeef5;
    }
  }

  .guide-content {
    .el-steps {
      max-width: 600px;
    }
  }

  .result-section {
    .result-content {
      .el-alert {
        margin-bottom: 20px;
      }
      
      .file-structure {
        margin: 20px 0;
        
        h4 {
          margin: 0 0 15px 0;
          color: #303133;
        }
        
        .file-tree {
          background: #f5f7fa;
          border-radius: 4px;
          padding: 15px;
          
          .file-node {
            display: flex;
            align-items: center;
            
            .el-icon {
              margin-right: 6px;
              font-size: 14px;
            }
          }
        }
      }
      
      .next-steps {
        margin: 20px 0;
        
        h4 {
          margin: 0 0 15px 0;
          color: #303133;
        }
      }
      
      .action-buttons {
        text-align: center;
        margin-top: 30px;
        
        .el-button {
          margin: 0 10px;
        }
      }
    }
  }

  .progress-content {
    text-align: center;
    
    .progress-text {
      margin-top: 15px;
      color: #606266;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .plugin-generator {
    padding: 15px;
    
    .generator-form .form-grid {
      grid-template-columns: 1fr;
      gap: 20px;
    }
  }
}
</style>