<template>
  <div class="installed-plugins">
    <div class="plugins-header">
      <h2>已安装插件 ({{ plugins.length }})</h2>
      <div class="header-actions">
        <el-button type="primary" @click="refreshPlugins" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 插件卡片展示 -->
    <div class="plugins-grid" v-loading="loading">
      <el-card 
        v-for="plugin in plugins" 
        :key="plugin.key" 
        class="plugin-card"
        shadow="hover"
      >
        <template #header>
          <div class="plugin-header">
            <div class="plugin-info">
              <el-icon class="plugin-icon"><component :is="plugin.icon || 'Connection'" /></el-icon>
              <div class="plugin-name">
                <span class="name">{{ plugin.name }}</span>
                <span class="version">v{{ plugin.version }}</span>
              </div>
            </div>
            <el-tag :type="getStatusType(plugin.status)" size="small">
              {{ getStatusText(plugin.status) }}
            </el-tag>
          </div>
        </template>
        
        <div class="plugin-content">
          <p class="plugin-description">{{ plugin.description || '暂无描述' }}</p>
          
          <div class="plugin-details">
            <div class="detail-item">
              <span class="label">作者:</span>
              <span class="value">{{ plugin.author || '未知' }}</span>
            </div>
            <div class="detail-item">
              <span class="label">页面数:</span>
              <span class="value">{{ plugin.pages?.length || 0 }} 个</span>
            </div>
            <div class="detail-item">
              <span class="label">分类:</span>
              <span class="value">{{ plugin.category || '其他' }}</span>
            </div>
          </div>
          
          <!-- 页面列表 -->
          <div class="plugin-pages" v-if="plugin.pages && plugin.pages.length > 0">
            <div class="pages-title">包含页面:</div>
            <div class="pages-list">
              <el-tag 
                v-for="page in plugin.pages" 
                :key="page.key"
                size="small"
                class="page-tag"
                @click="navigateToPage(page)"
              >
                <el-icon><component :is="page.icon || 'Document'" /></el-icon>
                {{ page.title }}
              </el-tag>
            </div>
          </div>
        </div>
        
        <template #footer>
          <div class="plugin-actions">
            <el-button 
              size="small" 
              type="primary"
              @click="navigateToPlugin(plugin)"
              :disabled="plugin.status !== 'installed'"
            >
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button 
              size="small" 
              :type="plugin.status === 'installed' ? 'warning' : 'success'"
              @click="togglePlugin(plugin)"
            >
              <el-icon><component :is="plugin.status === 'installed' ? 'VideoPause' : 'VideoPlay'" /></el-icon>
              {{ plugin.status === 'installed' ? '禁用' : '启用' }}
            </el-button>
            <el-dropdown @command="handleCommand">
              <el-button size="small">
                <el-icon><More /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="{ action: 'config', plugin }">
                    <el-icon><Setting /></el-icon>
                    配置
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'export', plugin }">
                    <el-icon><Download /></el-icon>
                    导出
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'uninstall', plugin }" divided>
                    <el-icon><Delete /></el-icon>
                    卸载
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-card>
    </div>

    <!-- 空状态 -->
    <el-empty 
      v-if="!loading && plugins.length === 0"
      description="暂无已安装的插件"
      :image-size="100"
    >
      <el-button type="primary" @click="$router.push('/plugin-generator')">
        创建插件
      </el-button>
    </el-empty>

    <!-- 配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      :title="currentPlugin ? `配置插件: ${currentPlugin.name}` : '配置插件'"
      width="50%"
    >
      <template v-if="currentPlugin">
        <el-form :model="pluginConfig" label-width="120px">
          <el-form-item label="插件状态">
            <el-switch
              v-model="pluginEnabled"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
          <el-form-item label="插件信息">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="名称">{{ currentPlugin.name }}</el-descriptions-item>
              <el-descriptions-item label="版本">{{ currentPlugin.version }}</el-descriptions-item>
              <el-descriptions-item label="作者">{{ currentPlugin.author }}</el-descriptions-item>
              <el-descriptions-item label="分类">{{ currentPlugin.category }}</el-descriptions-item>
            </el-descriptions>
          </el-form-item>
        </el-form>
      </template>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="configDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="savePluginConfig">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { 
  Refresh, 
  Connection, 
  View, 
  VideoPause, 
  VideoPlay, 
  More, 
  Setting, 
  Download, 
  Delete,
  Document
} from '@element-plus/icons-vue'
import { navigationService } from '@/services/navigation'
import { appsApi } from '@/api'

interface PluginInfo {
  key: string
  name: string
  version: string
  author: string
  description: string
  category: string
  icon: string
  status: 'installed' | 'inactive' | 'uninstalled'
  pages?: Array<{
    key: string
    title: string
    path: string
    name?: string
    icon: string
  }>
}

const router = useRouter()
const loading = ref(false)
const plugins = ref<PluginInfo[]>([])
const configDialogVisible = ref(false)
const currentPlugin = ref<PluginInfo | null>(null)
const pluginConfig = ref<Record<string, any>>({})
const pluginEnabled = ref(false)

// 从导航系统和插件meta文件加载插件信息
async function loadInstalledPlugins() {
  const pluginModules = import.meta.glob('../../plugins/*/meta.ts', { eager: true })
  const loadedPlugins: PluginInfo[] = []

  for (const [path, mod] of Object.entries(pluginModules)) {
    try {
      const meta = (mod as any).default
      if (!meta || !meta.mainNav) continue

      const pluginKey = meta.mainNav.key
      const plugin: PluginInfo = {
        key: pluginKey,
        name: meta.name || meta.mainNav.title,
        version: meta.version || '1.0.0',
        author: meta.author || '未知',
        description: meta.description || '暂无描述',
        category: meta.category || '其他',
        icon: meta.mainNav.icon || 'Connection',
        status: 'installed', // 默认为已安装
        pages: meta.pages?.map((page: any) => ({
          key: page.key,
          title: page.title,
          path: page.path,
          name: page.name,
          icon: page.icon || 'Document'
        })) || []
      }
      
      loadedPlugins.push(plugin)
      console.log(`✅ 加载插件: ${plugin.name}`)
    } catch (error) {
      console.error(`❌ 加载插件失败 ${path}:`, error)
    }
  }

  return loadedPlugins
}

// 刷新插件列表
async function refreshPlugins() {
  loading.value = true
  try {
    plugins.value = await loadInstalledPlugins()
    // 同时刷新导航系统
    await navigationService.refreshPlugins()
    ElMessage.success(`插件列表已更新，共找到 ${plugins.value.length} 个插件`)
  } catch (error) {
    console.error('刷新插件失败:', error)
    ElMessage.error('更新插件列表失败')
  } finally {
    loading.value = false
  }
}

// 获取状态类型
function getStatusType(status: string) {
  const types: Record<string, string> = {
    installed: 'success',
    inactive: 'warning',
    uninstalled: 'info'
  }
  return types[status] || 'info'
}

// 获取状态文本
function getStatusText(status: string) {
  const texts: Record<string, string> = {
    installed: '已启用',
    inactive: '已禁用',
    uninstalled: '未安装'
  }
  return texts[status] || status
}

// 切换插件状态
async function togglePlugin(plugin: PluginInfo) {
  try {
    const enabled = plugin.status !== 'installed'
    await appsApi.togglePlugin(plugin.key, enabled)
    plugin.status = enabled ? 'installed' : 'inactive'
    ElMessage.success(`${plugin.name} ${enabled ? '已启用' : '已禁用'}`)
    
    // 刷新导航菜单
    if (enabled) {
      await navigationService.refreshPlugins()
    } else {
      navigationService.removePluginMenu(plugin.key)
    }
  } catch (error) {
    console.error('切换插件状态失败:', error)
    ElMessage.error('操作失败')
  }
}

// 导航到插件主页
function navigateToPlugin(plugin: PluginInfo) {
  if (plugin.pages && plugin.pages.length > 0) {
    navigateToPage(plugin.pages[0])
  } else {
    ElMessage.warning('该插件暂无可访问的页面')
  }
}

// 导航到指定页面
function navigateToPage(page: any) {
  try {
    if (page.name) {
      router.push({ name: page.name })
    } else if (page.path) {
      router.push(page.path)
    }
  } catch (error) {
    console.error('导航失败:', error)
    ElMessage.error('页面跳转失败')
  }
}

// 处理下拉菜单命令
async function handleCommand(command: { action: string; plugin: PluginInfo }) {
  const { action, plugin } = command
  
  switch (action) {
    case 'config':
      configurePlugin(plugin)
      break
    case 'export':
      await exportPlugin(plugin)
      break
    case 'uninstall':
      await uninstallPlugin(plugin)
      break
  }
}

// 配置插件
function configurePlugin(plugin: PluginInfo) {
  currentPlugin.value = plugin
  pluginEnabled.value = plugin.status === 'installed'
  pluginConfig.value = {}
  configDialogVisible.value = true
}

// 保存插件配置
function savePluginConfig() {
  if (currentPlugin.value) {
    currentPlugin.value.status = pluginEnabled.value ? 'installed' : 'inactive'
    ElMessage.success('配置已保存')
  }
  configDialogVisible.value = false
}

// 导出插件
async function exportPlugin(plugin: PluginInfo) {
  try {
    const response = await appsApi.exportPlugin(plugin.key)
    const blob = new Blob([response.data], { type: 'application/zip' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${plugin.key}.zip`
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
    ElMessage.success(`插件 ${plugin.name} 导出成功`)
  } catch (error) {
    console.error('导出插件失败:', error)
    ElMessage.error('导出插件失败')
  }
}

// 卸载插件
async function uninstallPlugin(plugin: PluginInfo) {
  try {
    const result = await ElMessageBox.confirm(
      `确定要卸载插件 "${plugin.name}" 吗？这将删除所有相关数据。`,
      '确认卸载',
      {
        confirmButtonText: '卸载',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    if (result) {
      await appsApi.deletePlugin(plugin.key)
      // 从列表中移除
      const index = plugins.value.findIndex(p => p.key === plugin.key)
      if (index > -1) {
        plugins.value.splice(index, 1)
      }
      // 从导航中移除
      navigationService.removePluginMenu(plugin.key)
      ElMessage.success(`插件 ${plugin.name} 已卸载`)
    }
  } catch (error) {
    console.error('卸载插件失败:', error)
    ElMessage.error('卸载插件失败')
  }
}

// 初始化
onMounted(() => {
  refreshPlugins()
})
</script>

<style lang="scss" scoped>
.installed-plugins {
  padding: 20px;
  height: 100%;
  overflow: auto;

  .plugins-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      color: #303133;
    }
  }

  .plugins-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 20px;
    margin-bottom: 20px;
  }

  .plugin-card {
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
      .plugin-description {
        margin: 0 0 15px 0;
        color: #606266;
        line-height: 1.4;
      }
      
      .plugin-details {
        margin-bottom: 15px;
        
        .detail-item {
          display: flex;
          justify-content: space-between;
          margin-bottom: 8px;
          font-size: 14px;
          
          .label {
            color: #909399;
          }
          
          .value {
            color: #303133;
            font-weight: 500;
          }
        }
      }
      
      .plugin-pages {
        .pages-title {
          font-size: 14px;
          color: #606266;
          margin-bottom: 8px;
        }
        
        .pages-list {
          .page-tag {
            margin: 2px 4px 2px 0;
            cursor: pointer;
            
            .el-icon {
              margin-right: 4px;
            }
            
            &:hover {
              opacity: 0.8;
            }
          }
        }
      }
    }
    
    .plugin-actions {
      display: flex;
      gap: 8px;
      
      .el-button {
        flex: 1;
      }
    }
  }
}
</style> 