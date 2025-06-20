<template>
  <div class="debug-page">
    <h1>🔍 系统调试页面</h1>
    <p>如果你能看到这个页面，说明基本路由和组件加载是正常的。</p>
    
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="debug-info">
          <template #header>
            <span>系统信息</span>
          </template>
          <div class="info-list">
            <div><strong>当前路由:</strong> {{ $route.path }}</div>
            <div><strong>Vue版本:</strong> Vue 3</div>
            <div><strong>Element Plus:</strong> 已加载</div>
            <div><strong>时间:</strong> {{ currentTime }}</div>
            <div><strong>开发服务器:</strong> 
              <el-tag type="success" size="small">运行中</el-tag>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="navigation-status">
          <template #header>
            <span>导航状态</span>
            <el-button size="small" @click="checkNavigationStatus" :loading="checking">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </template>
          <div class="status-list">
            <div class="status-item">
              <span>主菜单数量:</span>
              <el-tag size="small">{{ navigationStatus.mainMenuCount }}</el-tag>
            </div>
            <div class="status-item">
              <span>插件菜单数量:</span>
              <el-tag size="small">{{ navigationStatus.pluginMenuCount }}</el-tag>
            </div>
            <div class="status-item">
              <span>导航服务:</span>
              <el-tag :type="navigationStatus.serviceReady ? 'success' : 'danger'" size="small">
                {{ navigationStatus.serviceReady ? '正常' : '异常' }}
              </el-tag>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="actions" style="margin-top: 20px;">
      <template #header>
        <span>快速操作</span>
      </template>
      <el-space wrap>
        <el-button type="primary" @click="testAlert">测试弹窗</el-button>
        <el-button type="success" @click="goToGenerator">插件生成器</el-button>
        <el-button type="info" @click="goToAppMarket">应用市场</el-button>
        <el-button type="warning" @click="checkNavigationStatus">检查导航</el-button>
        <el-button @click="showMainMenus">查看主菜单</el-button>
        <el-button @click="showPluginInfo">查看插件信息</el-button>
      </el-space>
    </el-card>

    <el-card v-if="showMenuDetail" class="menu-detail" style="margin-top: 20px;">
      <template #header>
        <span>菜单详情</span>
        <el-button size="small" @click="showMenuDetail = false">
          <el-icon><Close /></el-icon>
          关闭
        </el-button>
      </template>
      <el-tabs>
        <el-tab-pane label="主菜单" name="main">
          <pre>{{ JSON.stringify(navigationStatus.mainMenus, null, 2) }}</pre>
        </el-tab-pane>
        <el-tab-pane label="子菜单映射" name="sub">
          <pre>{{ JSON.stringify(navigationStatus.subMenusMap, null, 2) }}</pre>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-card v-if="showPluginDetail" class="plugin-detail" style="margin-top: 20px;">
      <template #header>
        <span>插件详情</span>
        <el-button size="small" @click="showPluginDetail = false">
          <el-icon><Close /></el-icon>
          关闭
        </el-button>
      </template>
      <el-table :data="pluginList" stripe>
        <el-table-column prop="key" label="插件Key" width="150" />
        <el-table-column prop="name" label="插件名称" width="150" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="author" label="作者" width="120" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="pagesCount" label="页面数" width="100" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Refresh, Close } from '@element-plus/icons-vue'
import { useNavigationStore } from '@/stores/navigation'
import { navigationService } from '@/services/navigation'

const router = useRouter()
const currentTime = ref('')
const checking = ref(false)
const showMenuDetail = ref(false)
const showPluginDetail = ref(false)
const pluginList = ref<any[]>([])

const navigationStatus = reactive({
  mainMenuCount: 0,
  pluginMenuCount: 0,
  serviceReady: false,
  mainMenus: [] as any[],
  subMenusMap: {} as any
})

const updateTime = () => {
  currentTime.value = new Date().toLocaleString()
}

const testAlert = () => {
  ElMessage.success('Element Plus 组件工作正常！')
}

const goToGenerator = () => {
  router.push('/plugin-generator')
}

const goToAppMarket = () => {
  router.push('/app-market')
}

const checkNavigationStatus = async () => {
  checking.value = true
  try {
    const navigationStore = useNavigationStore()
    
    // 获取导航状态
    const navState = navigationService.getNavigationState()
    navigationStatus.mainMenus = navState.mainMenus
    navigationStatus.subMenusMap = navState.subMenusMap
    navigationStatus.mainMenuCount = navState.mainMenus.length
    
    // 计算插件菜单数量（排除系统菜单）
    const systemMenuKeys = ['home', 'app-market', 'debug']
    navigationStatus.pluginMenuCount = navState.mainMenus.filter(
      menu => !systemMenuKeys.includes(menu.key)
    ).length
    
    navigationStatus.serviceReady = true
    
    console.log('🔍 导航状态检查完成:', navigationStatus)
    ElMessage.success('导航状态检查完成')
  } catch (error) {
    console.error('❌ 导航状态检查失败:', error)
    navigationStatus.serviceReady = false
    ElMessage.error('导航状态检查失败')
  } finally {
    checking.value = false
  }
}

const showMainMenus = () => {
  showMenuDetail.value = true
  checkNavigationStatus()
}

const showPluginInfo = async () => {
  try {
    // 扫描插件目录
    const pluginModules = import.meta.glob('../plugins/*/meta.ts', { eager: true })
    const plugins = []
    
    for (const [path, mod] of Object.entries(pluginModules)) {
      const meta = (mod as any).default
      if (meta && meta.mainNav) {
        plugins.push({
          key: meta.mainNav.key,
          name: meta.name || meta.mainNav.title,
          version: meta.version || '1.0.0',
          author: meta.author || '未知',
          description: meta.description || '暂无描述',
          pagesCount: meta.pages?.length || 0
        })
      }
    }
    
    pluginList.value = plugins
    showPluginDetail.value = true
    ElMessage.success(`找到 ${plugins.length} 个插件`)
  } catch (error) {
    console.error('❌ 获取插件信息失败:', error)
    ElMessage.error('获取插件信息失败')
  }
}

onMounted(() => {
  updateTime()
  setInterval(updateTime, 1000)
  console.log('🔍 调试页面已加载')
  
  // 自动检查一次导航状态
  setTimeout(() => {
    checkNavigationStatus()
  }, 1000)
})
</script>

<style scoped>
.debug-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.debug-info .info-list div,
.navigation-status .status-list .status-item {
  margin: 8px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.actions {
  margin-top: 20px;
}

.menu-detail pre,
.plugin-detail {
  max-height: 400px;
  overflow: auto;
}

.menu-detail pre {
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
}

.el-card .el-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 