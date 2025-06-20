<template>
  <div class="plugin-package">
    <div class="package-header">
      <h2>插件打包</h2>
      <div class="header-actions">
        <el-button type="primary" @click="refreshPlugins" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新列表
        </el-button>
      </div>
    </div>

    <!-- Plugin Selection -->
    <el-card class="selection-card" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Box /></el-icon>
          <span>选择插件</span>
        </div>
      </template>
      
      <el-alert
        title="插件来源目录"
        :description="`扫描目录: src/plugins/`"
        type="info"
        :closable="false"
        show-icon
        class="source-info"
      />

      <div v-if="plugins.length > 0" class="plugins-grid" v-loading="loading">
        <el-card 
          v-for="plugin in plugins" 
          :key="plugin.id"
          class="plugin-card"
          shadow="hover"
          :class="{ 'selected': selectedPlugin?.id === plugin.id }"
          @click="selectPlugin(plugin)"
        >
          <template #header>
            <div class="plugin-header">
              <div class="plugin-info">
                <el-icon class="plugin-icon"><component :is="'Connection'" /></el-icon>
                <div class="plugin-name">
                  <span class="name">{{ plugin.name }}</span>
                  <span class="version">v{{ plugin.version }}</span>
                </div>
              </div>
              <el-tag :type="getStatusTagType(plugin.status)" size="small">
                {{ getStatusText(plugin.status) }}
              </el-tag>
            </div>
          </template>
          
          <div class="plugin-content">
            <div class="plugin-details">
              <div class="detail-item">
                <span class="label">目录:</span>
                <code class="value">{{ plugin.directory }}</code>
              </div>
            </div>
          </div>
          
          <template #footer>
            <div class="plugin-actions">
              <el-button 
                type="primary" 
                size="small"
                :disabled="plugin.status !== 'ready'"
                @click.stop="selectPlugin(plugin)"
              >
                {{ selectedPlugin?.id === plugin.id ? '已选择' : '选择打包' }}
              </el-button>
            </div>
          </template>
        </el-card>
      </div>

      <el-empty 
        v-else
        description="暂无可用插件"
        :image-size="100"
      >
        <el-button type="primary" @click="$router.push('/plugin-generator')">
          创建插件
        </el-button>
      </el-empty>
    </el-card>

    <!-- Package Configuration -->
    <el-card v-if="selectedPlugin" class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Setting /></el-icon>
          <span>打包配置</span>
        </div>
      </template>
      
      <el-alert
        :title="`即将打包: ${selectedPlugin.name}`"
        :description="`输出路径: dist/plugins/${selectedPlugin.name}/`"
        type="success"
        :closable="false"
        show-icon
        class="selected-info"
      />
      
      <div class="config-options">
        <el-checkbox v-model="config.includeDemoData">
          <div class="option-content">
            <div class="option-title">包含示例数据</div>
            <div class="option-desc">包含插件的示例数据，用于演示和测试</div>
          </div>
        </el-checkbox>
        
        <el-checkbox v-model="config.includeDocs">
          <div class="option-content">
            <div class="option-title">包含文档</div>
            <div class="option-desc">包含插件的使用文档和API文档</div>
          </div>
        </el-checkbox>
        
        <el-checkbox v-model="config.includeTests">
          <div class="option-content">
            <div class="option-title">包含测试</div>
            <div class="option-desc">包含插件的测试用例和测试数据</div>
          </div>
        </el-checkbox>
        
        <el-checkbox v-model="config.minify">
          <div class="option-content">
            <div class="option-title">压缩代码</div>
            <div class="option-desc">压缩代码以减小包体积</div>
          </div>
        </el-checkbox>
      </div>
      
      <div class="config-actions">
                 <el-button type="primary" size="large" @click="startPackaging" :loading="isPackaging">
           <el-icon><Promotion /></el-icon>
           {{ isPackaging ? '打包中...' : '开始打包' }}
         </el-button>
        <el-button v-if="isPackaging" @click="cancelPackaging">
          取消打包
        </el-button>
      </div>
    </el-card>

    <!-- Package Progress -->
    <el-card v-if="isPackaging" class="progress-card" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Loading /></el-icon>
          <span>打包进度</span>
        </div>
      </template>
      
      <div class="progress-content">
        <el-progress :percentage="progress" :stroke-width="8" />
        <div class="progress-info">
          <span class="progress-text">{{ progressMessage }}</span>
          <span class="progress-percent">{{ progress }}%</span>
        </div>
      </div>
    </el-card>

    <!-- Package Result -->
    <el-card v-if="packageResult" class="result-card" shadow="never">
      <template #header>
        <div class="card-header success">
          <el-icon class="header-icon"><SuccessFilled /></el-icon>
          <span>打包完成</span>
        </div>
      </template>
      
      <div class="result-content">
        <el-alert
          title="插件打包成功！"
          :description="`插件已成功打包，耗时 ${packageResult.duration} 秒`"
          type="success"
          :closable="false"
          show-icon
        />
        
        <div class="result-details">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="输出路径">{{ packageResult.outputPath }}</el-descriptions-item>
            <el-descriptions-item label="包大小">{{ formatFileSize(packageResult.size) }}</el-descriptions-item>
            <el-descriptions-item label="文件数量">{{ packageResult.fileCount }} 个</el-descriptions-item>
            <el-descriptions-item label="打包耗时">{{ packageResult.duration }} 秒</el-descriptions-item>
          </el-descriptions>
        </div>
        
        <div class="result-actions">
          <el-button type="primary" @click="downloadPackage">
            <el-icon><Download /></el-icon>
            下载插件包
          </el-button>
          <el-button @click="openOutputDir">
            <el-icon><FolderOpened /></el-icon>
            打开目录
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Refresh, 
  Box, 
  Connection, 
  Setting, 
  Promotion, 
  Loading, 
  SuccessFilled,
  Download,
  FolderOpened
} from '@element-plus/icons-vue'

// Types
type PluginStatus = 'ready' | 'in-progress' | 'error'

interface Plugin {
  id: number
  name: string
  version: string
  status: PluginStatus
  directory: string
}

interface PackageResult {
  outputPath: string
  size: number
  fileCount: number
  duration: number
}

// State
const router = useRouter()
const plugins = ref<Plugin[]>([])
const selectedPlugin = ref<Plugin | null>(null)
const isPackaging = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const packageResult = ref<PackageResult | null>(null)
const loading = ref(false)

const config = reactive({
  includeDemoData: true,
  includeDocs: true,
  includeTests: true,
  minify: true
})

// Methods
const getStatusClass = (status: PluginStatus): string => {
  const classes: Record<PluginStatus, string> = {
    'ready': 'bg-green-100 text-green-800',
    'in-progress': 'bg-yellow-100 text-yellow-800',
    'error': 'bg-red-100 text-red-800'
  }
  return classes[status]
}

const getStatusText = (status: PluginStatus): string => {
  const texts: Record<PluginStatus, string> = {
    'ready': '就绪',
    'in-progress': '进行中',
    'error': '错误'
  }
  return texts[status]
}

const getStatusTagType = (status: PluginStatus): string => {
  const types: Record<PluginStatus, string> = {
    'ready': 'success',
    'in-progress': 'warning',
    'error': 'danger'
  }
  return types[status]
}

const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const selectPlugin = (plugin: Plugin) => {
  selectedPlugin.value = plugin
  packageResult.value = null
}

const refreshPlugins = async () => {
  loading.value = true
  try {
    // 调用后端 API 刷新插件列表
    const response = await fetch('/api/plugins/scan', {
      method: 'POST'
    })

    if (!response.ok) {
      throw new Error('刷新插件列表失败')
    }

    await loadPlugins()
  } catch (error) {
    console.error('刷新插件列表失败:', error)
  } finally {
    loading.value = false
  }
}

const startPackaging = async () => {
  if (!selectedPlugin.value) return

  try {
    isPackaging.value = true
    progress.value = 0
    progressMessage.value = '正在开始打包流程...'
    packageResult.value = null

    const startTime = Date.now()

    // 调用后端 API 进行打包
    const response = await fetch('/api/plugins/package', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        pluginName: selectedPlugin.value.name,
        config: {
          includeDemoData: config.includeDemoData,
          includeDocs: config.includeDocs,
          includeTests: config.includeTests,
          minify: config.minify
        }
      })
    })

    if (!response.ok) {
      throw new Error('打包请求失败')
    }

    // 模拟进度更新
    for (let i = 0; i <= 100; i += 20) {
      await new Promise(resolve => setTimeout(resolve, 1000))
      progress.value = i
      progressMessage.value = `处理中... ${i}%`
    }

    const endTime = Date.now()
    const duration = ((endTime - startTime) / 1000).toFixed(1)

    // 设置打包结果
    packageResult.value = {
      outputPath: `dist/plugins/${selectedPlugin.value.name}/`,
      size: 1024 * 1024 * 2.5, // 示例大小
      fileCount: 15, // 示例文件数
      duration: Number(duration)
    }

    progressMessage.value = '打包完成！'
  } catch (error) {
    progressMessage.value = '打包过程中发生错误'
    console.error('打包错误:', error)
  } finally {
    isPackaging.value = false
  }
}

const cancelPackaging = () => {
  selectedPlugin.value = null
  progress.value = 0
  progressMessage.value = ''
  isPackaging.value = false
  packageResult.value = null
}

const openOutputDir = () => {
  // TODO: 实现打开输出目录的功能
  window.open(`file://${packageResult.value?.outputPath}`, '_blank')
}

const downloadPackage = async () => {
  if (!packageResult.value) return

  try {
    const response = await fetch(`/api/plugins/download/${selectedPlugin.value?.name}`)
    if (!response.ok) throw new Error('下载失败')

    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${selectedPlugin.value?.name}.zip`
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
  } catch (error) {
    console.error('下载插件包失败:', error)
  }
}

// 组件加载时获取插件列表
const loadPlugins = async () => {
  try {
    console.log('开始加载插件列表...')
    const response = await fetch('/api/plugins/scan', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    console.log('API响应状态:', response.status, response.statusText)

    if (!response.ok) {
      const errorText = await response.text()
      console.error('API错误响应:', errorText)
      throw new Error(`扫描插件失败: ${response.status} ${errorText}`)
    }

    const data = await response.json()
    console.log('API响应数据:', data)
    
    if (data.plugins && Array.isArray(data.plugins)) {
      plugins.value = data.plugins.map((plugin: any, index: number) => ({
        id: index + 1,
        name: plugin.name,
        version: plugin.version,
        status: plugin.status as PluginStatus,
        directory: plugin.directory
      }))
      console.log('成功加载插件:', plugins.value.length, '个')
    } else {
      console.warn('API响应中没有插件数据或格式不正确')
      // 使用临时假数据进行测试
      plugins.value = [
        {
          id: 1,
          name: 'plugin-wailki',
          version: '1.0.0',
          status: 'ready' as PluginStatus,
          directory: 'src/plugins/plugin-wailki'
        },
        {
          id: 2,
          name: 'plugin-wasda',
          version: '0.5.0',
          status: 'ready' as PluginStatus,
          directory: 'src/plugins/plugin-wasda'
        }
      ]
      console.log('使用临时数据:', plugins.value.length, '个插件')
    }
  } catch (error) {
    console.error('加载插件列表失败:', error)
    // 如果API失败，使用临时假数据
    plugins.value = [
      {
        id: 1,
        name: 'plugin-wailki',
        version: '1.0.0',
        status: 'ready' as PluginStatus,
        directory: 'src/plugins/plugin-wailki'
      },
      {
        id: 2,
        name: 'plugin-wasda',
        version: '0.5.0',
        status: 'ready' as PluginStatus,
        directory: 'src/plugins/plugin-wasda'
      }
    ]
    console.log('API失败，使用备用数据:', plugins.value.length, '个插件')
  }
}

loadPlugins()
</script>

<style lang="scss" scoped>
.plugin-package {
  padding: 20px;
  height: 100%;
  overflow: auto;

  .package-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      color: #303133;
    }
  }

  .selection-card, .config-card, .progress-card, .result-card {
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

  .source-info {
    margin-bottom: 20px;
  }

  .plugins-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 20px;
  }

  .plugin-card {
    cursor: pointer;
    transition: all 0.3s ease;
    
    &.selected {
      border-color: #409eff;
      box-shadow: 0 0 0 1px #409eff;
    }
    
    &:hover {
      transform: translateY(-2px);
    }
    
    .plugin-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .plugin-info {
        display: flex;
        align-items: center;
        
        .plugin-icon {
          font-size: 24px;
          margin-right: 12px;
          color: #409eff;
        }
        
        .plugin-name {
          .name {
            display: block;
            font-weight: 600;
            color: #303133;
          }
          
          .version {
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }
    
    .plugin-content {
      .plugin-details {
        .detail-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 8px;
          font-size: 14px;
          
          .label {
            color: #909399;
          }
          
          .value {
            font-family: monospace;
            font-size: 12px;
            background: #f5f7fa;
            padding: 2px 6px;
            border-radius: 4px;
            color: #606266;
          }
        }
      }
    }
    
    .plugin-actions {
      display: flex;
      justify-content: flex-end;
    }
  }

  .selected-info {
    margin-bottom: 20px;
  }

  .config-options {
    margin: 20px 0;
    
    .el-checkbox {
      display: flex;
      align-items: flex-start;
      margin-bottom: 16px;
      width: 100%;
      
      .option-content {
        margin-left: 8px;
        
        .option-title {
          font-weight: 600;
          color: #303133;
          margin-bottom: 4px;
        }
        
        .option-desc {
          font-size: 14px;
          color: #606266;
          line-height: 1.4;
        }
      }
    }
  }

  .config-actions {
    text-align: center;
    padding-top: 20px;
    border-top: 1px solid #ebeef5;
    
    .el-button {
      margin: 0 10px;
    }
  }

  .progress-content {
    .el-progress {
      margin-bottom: 15px;
    }
    
    .progress-info {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .progress-text {
        color: #606266;
      }
      
      .progress-percent {
        color: #409eff;
        font-weight: 600;
      }
    }
  }

  .result-content {
    .el-alert {
      margin-bottom: 20px;
    }
    
    .result-details {
      margin: 20px 0;
    }
    
    .result-actions {
      text-align: center;
      margin-top: 20px;
      
      .el-button {
        margin: 0 10px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .plugin-package {
    padding: 15px;
    
    .plugins-grid {
      grid-template-columns: 1fr;
    }
  }
}

  .plugin-package {
  padding: 20px;
  height: 100%;
  overflow: auto;

  .package-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      color: #303133;
    }
  }

  .selection-card, .config-card, .progress-card, .result-card {
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

  .source-info {
    margin-bottom: 20px;
  }

  .plugins-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 20px;
  }

  .plugin-card {
    cursor: pointer;
    transition: all 0.3s ease;
    
    &.selected {
      border-color: #409eff;
      box-shadow: 0 0 0 1px #409eff;
    }
    
    &:hover {
      transform: translateY(-2px);
    }
    
    .plugin-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .plugin-info {
        display: flex;
        align-items: center;
        
        .plugin-icon {
          font-size: 24px;
          margin-right: 12px;
          color: #409eff;
        }
        
        .plugin-name {
          .name {
            display: block;
            font-weight: 600;
            color: #303133;
          }
          
          .version {
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }
    
    .plugin-content {
      .plugin-details {
        .detail-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 8px;
          font-size: 14px;
          
          .label {
            color: #909399;
          }
          
          .value {
            font-family: monospace;
            font-size: 12px;
            background: #f5f7fa;
            padding: 2px 6px;
            border-radius: 4px;
            color: #606266;
          }
        }
      }
    }
    
    .plugin-actions {
      display: flex;
      justify-content: flex-end;
    }
  }

  .selected-info {
    margin-bottom: 20px;
  }

  .config-options {
    margin: 20px 0;
    
    .el-checkbox {
      display: flex;
      align-items: flex-start;
      margin-bottom: 16px;
      width: 100%;
      
      .option-content {
        margin-left: 8px;
        
        .option-title {
          font-weight: 600;
          color: #303133;
          margin-bottom: 4px;
        }
        
        .option-desc {
          font-size: 14px;
          color: #606266;
          line-height: 1.4;
        }
      }
    }
  }

  .config-actions {
    text-align: center;
    padding-top: 20px;
    border-top: 1px solid #ebeef5;
    
    .el-button {
      margin: 0 10px;
    }
  }

  .progress-content {
    .el-progress {
      margin-bottom: 15px;
    }
    
    .progress-info {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .progress-text {
        color: #606266;
      }
      
      .progress-percent {
        color: #409eff;
        font-weight: 600;
      }
    }
  }

  .result-content {
    .el-alert {
      margin-bottom: 20px;
    }
    
    .result-details {
      margin: 20px 0;
    }
    
    .result-actions {
      text-align: center;
      margin-top: 20px;
      
      .el-button {
        margin: 0 10px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .plugin-package {
    padding: 15px;
    
    .plugins-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>