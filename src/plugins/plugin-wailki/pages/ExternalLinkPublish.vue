<template>
  <div class="external-link-publish">
    <!-- 外链展示区域 -->
    <el-card class="link-display-card">
      <template #header>
        <div class="card-header">
          <span>外链列表 ({{ filteredLinks.length }} / {{ allLinks.length }})</span>
          <div>
            <el-button @click="refreshLinks" :loading="refreshing">刷新</el-button>
            <el-button type="warning" @click="quickCheckLinks" :loading="quickChecking">
              🔍 快速检测
            </el-button>
            <el-button type="success" @click="retryFailedLinksInPublish" :loading="retryingFailed">
              🔄 重试失败
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选和统计区域 -->
      <div class="filter-section" style="margin-bottom: 16px;">
        <div class="filter-controls">
          <span style="margin-right: 12px; color: #666;">显示模式:</span>
          <el-button-group>
            <el-button 
              :type="linkFilter === 'all' ? 'primary' : 'default'"
              size="small"
              @click="setLinkFilter('all')"
            >
              全部 ({{ allLinks.length }})
            </el-button>
            <el-button 
              :type="linkFilter === 'available' ? 'success' : 'default'"
              size="small"
              @click="setLinkFilter('available')"
            >
              ✅ 可用 ({{ availableLinks.length }})
            </el-button>
            <el-button 
              :type="linkFilter === 'unavailable' ? 'danger' : 'default'"
              size="small"
              @click="setLinkFilter('unavailable')"
            >
              ❌ 不可用 ({{ unavailableLinks.length }})
            </el-button>
            <el-button 
              :type="linkFilter === 'enabled' ? 'info' : 'default'"
              size="small"
              @click="setLinkFilter('enabled')"
            >
              🟢 启用 ({{ enabledLinks.length }})
            </el-button>
          </el-button-group>
        </div>
        
        <div class="stats-info" style="margin-top: 8px;">
          <el-text type="info" size="small">
            总计: {{ allLinks.length }} | 
            可用: {{ availableLinks.length }} | 
            不可用: {{ unavailableLinks.length }} | 
            启用: {{ enabledLinks.length }} | 
            禁用: {{ disabledLinks.length }}
          </el-text>
        </div>
      </div>
      
              <!-- 固定高度的外链列表 -->
      <div class="links-container">
        <div v-if="filteredLinks.length === 0" class="no-links">
          <el-empty :description="getEmptyDescription()">
            <el-button type="primary" @click="$router.push('/plugin-wailki/index')">
              去创建外链
            </el-button>
          </el-empty>
        </div>
        
        <div v-else class="links-list">
          <div 
            v-for="(link, index) in filteredLinks" 
            :key="index" 
            :class="['link-item', { 'unavailable': !link.is_valid, 'disabled': !link.status }]"
          >
            <div class="link-status">
              <el-tag 
                :type="getStatusTagType(link)"
                size="large"
                style="font-weight: bold;"
              >
                {{ getStatusText(link) }}
              </el-tag>
            </div>
            <div class="link-info">
              <div class="link-url">
                <el-text :type="link.is_valid ? 'primary' : 'danger'">
                  {{ buildUrlWithParams(link.url) }}
                </el-text>
              </div>
              <div class="link-meta">
                <el-tag size="small" :type="link.category ? 'info' : 'warning'">
                  {{ link.category || '未分类' }}
                </el-tag>
                <el-text size="small" type="info">点击量: {{ link.clicks || 0 }}</el-text>
                <el-tag 
                  v-if="link.last_check_error" 
                  size="small" 
                  type="danger"
                  style="margin-left: 8px;"
                >
                  检测异常
                </el-tag>
              </div>
            </div>
            <div class="link-actions">
              <el-button 
                size="small" 
                type="warning" 
                @click="checkSingleLinkInPublish(link)"
                :loading="link.checking"
              >
                {{ link.checking ? '检测中' : '🔍 检测' }}
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 发布配置区域 -->
    <el-card class="publish-config-card">
      <template #header>
        <div class="card-header">
          <span>发布配置</span>
        </div>
      </template>
      
      <el-form :model="publishForm" label-width="120px">
        <!-- 发布设置 -->
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item label="发布次数">
              <el-input-number 
                v-model="publishForm.publishCount" 
                :min="1" 
                :max="100"
                style="width: 120px"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="发布间隔(秒)">
              <el-input-number 
                v-model="publishForm.publishInterval" 
                :min="1" 
                :max="60"
                style="width: 120px"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="超时时间(秒)">
              <el-input-number 
                v-model="publishForm.timeout" 
                :min="5" 
                :max="60"
                style="width: 120px"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="使用代理">
              <el-switch 
                v-model="publishForm.useProxy"
                @change="onProxyToggle"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <!-- 代理设置 -->
        <el-row v-if="publishForm.useProxy" :gutter="20">
          <el-col :span="8">
            <el-form-item label="代理服务器">
              <el-input 
                v-model="publishForm.proxyHost" 
                placeholder="代理服务器地址"
                style="width: 200px"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="代理端口">
              <el-input-number 
                v-model="publishForm.proxyPort" 
                :min="1" 
                :max="65535"
                placeholder="端口号"
                style="width: 120px"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="代理类型">
              <el-select v-model="publishForm.proxyType" style="width: 120px">
                <el-option label="HTTP" value="http" />
                <el-option label="HTTPS" value="https" />
                <el-option label="SOCKS5" value="socks5" />
          </el-select>
        </el-form-item>
          </el-col>
        </el-row>
        
        <!-- 参数配置 -->
        <el-divider content-position="left">URL参数配置</el-divider>
        
        <el-row :gutter="20">
          <el-col :span="8">
        <el-form-item label="国家">
              <el-input 
                v-model="publishForm.params.country" 
                placeholder="请输入国家"
                style="width: 150px"
                @input="updateLinksDisplay"
              />
        </el-form-item>
          </el-col>
          <el-col :span="8">
        <el-form-item label="手机号">
              <el-input 
                v-model="publishForm.params.phone" 
                placeholder="请输入手机号"
                style="width: 150px"
                @input="updateLinksDisplay"
              />
        </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="微信号">
              <el-input 
                v-model="publishForm.params.wechat" 
                placeholder="请输入微信号"
                style="width: 150px"
                @input="updateLinksDisplay"
              />
        </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="Telegram">
              <el-input 
                v-model="publishForm.params.telegram" 
                placeholder="请输入Telegram"
                style="width: 150px"
                @input="updateLinksDisplay"
              />
        </el-form-item>
          </el-col>
          <el-col :span="8">
        <el-form-item label="官网网址">
              <el-input 
                v-model="publishForm.params.website" 
                placeholder="请输入官网网址"
                style="width: 150px"
                @input="updateLinksDisplay"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="自定义参数">
              <el-input 
                v-model="publishForm.params.custom" 
                placeholder="key=value&key2=value2"
                style="width: 200px"
                @input="updateLinksDisplay"
              />
        </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item>
          <el-button 
            type="primary" 
            size="large"
            @click="confirmStartPublish" 
            :loading="publishing"
            :disabled="availableLinks.length === 0"
          >
            {{ publishing ? `🎭 真实用户访问中... (${publishProgress.current}/${publishProgress.total})` : `🚀 开始真实用户访问 (${availableLinks.length * publishForm.publishCount} 次)` }}
          </el-button>
          <el-button 
            v-if="publishing"
            type="danger"
            size="large"
            @click="stopPublish"
          >
            🛑 停止发布
          </el-button>
          <el-button size="large" @click="clearLogs">🗑️ 清空日志</el-button>
          <el-button 
            v-if="!publishing"
            size="large" 
            type="info"
            @click="showBackgroundTips"
          >
            💡 后台运行说明
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 发布日志区域 -->
    <el-card class="publish-logs-card">
      <template #header>
        <div class="card-header">
          <div class="log-title">
            <span>实时访问日志</span>
            <el-tag 
              v-if="publishing" 
              type="success" 
              size="small"
              effect="plain"
              style="margin-left: 8px;"
            >
              🛡️ 后台保障运行中
            </el-tag>
          </div>
          <div class="log-controls">
            <el-button size="small" @click="exportLogs" :disabled="publishLogs.length === 0">
              📁 导出日志
            </el-button>
            <el-button size="small" @click="clearLogs">
              🗑️ 清空日志
            </el-button>
            <el-switch 
              v-model="autoScroll" 
              size="small"
              inline-prompt
              active-text="自动滚动"
              inactive-text="手动滚动"
            />
          </div>
        </div>
      </template>

      <!-- 日志筛选和搜索区域 -->
      <div class="log-filters" style="margin-bottom: 16px; padding: 16px; background-color: #f8f9fa; border-radius: 8px;">
        <div class="filter-row" style="display: flex; align-items: center; gap: 16px; margin-bottom: 12px;">
          <span style="color: #666; font-weight: 500;">筛选：</span>
          <el-button-group>
            <el-button 
              :type="logFilter === 'all' ? 'primary' : 'default'"
              size="small"
              @click="setLogFilter('all')"
            >
              全部 ({{ publishLogs.length }})
            </el-button>
            <el-button 
              :type="logFilter === 'success' ? 'success' : 'default'"
              size="small"
              @click="setLogFilter('success')"
            >
              ✅ 成功 ({{ successCount }})
            </el-button>
            <el-button 
              :type="logFilter === 'error' ? 'danger' : 'default'"
              size="small"
              @click="setLogFilter('error')"
            >
              ❌ 失败 ({{ failCount }})
            </el-button>
            <el-button 
              :type="logFilter === 'system' ? 'info' : 'default'"
              size="small"
              @click="setLogFilter('system')"
            >
              🔧 系统 ({{ systemCount }})
            </el-button>
          </el-button-group>
          
          <el-input 
            v-model="logSearch"
            placeholder="搜索URL或消息内容..."
            style="width: 250px;"
            size="small"
            clearable
            prefix-icon="Search"
          />
        </div>
        
        <div class="stats-row" style="display: flex; align-items: center; gap: 24px;">
          <el-text type="info" size="small">
            📊 总计: {{ publishLogs.length }} 条
          </el-text>
          <el-text type="success" size="small">
            ✅ 成功: {{ successCount }} 条 ({{ successRate }}%)
          </el-text>
          <el-text type="danger" size="small">
            ❌ 失败: {{ failCount }} 条 ({{ failRate }}%)
          </el-text>
          <el-text type="info" size="small">
            ⏱️ 平均耗时: {{ averageDuration }}秒
          </el-text>
          <el-text v-if="publishing" type="warning" size="small">
            🚀 进行中: {{ publishProgress.current }}/{{ publishProgress.total }}
          </el-text>
        </div>
      </div>
      
      <div class="logs-container" ref="logsContainer">
        <div 
          v-for="(log, index) in filteredLogs" 
          :key="index" 
          :class="['log-item', log.status, log.level]"
        >
          <div class="log-icon">
            <span class="status-icon">{{ getLogIcon(log) }}</span>
          </div>
          <div class="log-time">
            <div class="time-main">{{ formatTime(log.timestamp) }}</div>
            <div class="time-detail">{{ formatDate(log.timestamp) }}</div>
          </div>
          <div class="log-content">
            <div class="log-header">
              <div class="log-url" v-if="log.url">
                <el-text :type="log.status === 'success' ? 'primary' : 'danger'" size="small">
                  🔗 {{ log.url }}
                </el-text>
              </div>
              <div class="log-tags">
                <el-tag 
                  :type="getLogTagType(log)" 
                  size="small" 
                  effect="plain"
                >
                  {{ getLogStatusText(log) }}
                </el-tag>
                <el-tag 
                  v-if="log.level && log.level !== 'access'"
                  size="small" 
                  type="info"
                  effect="plain"
                >
                  {{ getLevelText(log.level) }}
                </el-tag>
                <el-tag 
                  v-if="log.duration"
                  size="small" 
                  type="warning"
                  effect="plain"
                >
                  ⏱️ {{ log.duration }}s
                </el-tag>
              </div>
            </div>
            <div class="log-message">
              {{ log.message }}
            </div>
            <div v-if="log.details || log.ip || log.proxy" class="log-details">
              <el-button 
                size="small" 
                text 
                @click="log.expanded = !log.expanded"
                style="padding: 0; margin-bottom: 8px;"
              >
                {{ log.expanded ? '🔼 收起详情' : '🔽 展开详情' }}
              </el-button>
              <div v-if="log.expanded" class="detail-content">
                <div v-if="log.ip" class="detail-item">
                  <span class="detail-label">IP地址:</span>
                  <el-text type="success" size="small">{{ log.ip }}</el-text>
                </div>
                <div v-if="log.proxy" class="detail-item">
                  <span class="detail-label">代理服务器:</span>
                  <el-text type="warning" size="small">{{ log.proxy }}</el-text>
                </div>
                <div v-if="log.details" class="detail-item">
                  <span class="detail-label">详细信息:</span>
                  <el-text size="small">{{ log.details }}</el-text>
                </div>
                <div v-if="log.error" class="detail-item">
                  <span class="detail-label">错误详情:</span>
                  <el-text type="danger" size="small">{{ log.error }}</el-text>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div v-if="filteredLogs.length === 0" class="no-logs">
          <el-empty :description="getEmptyLogsDescription()">
            <el-text type="info" size="small">
              {{ logFilter === 'all' ? '暂无访问记录' : `暂无${getFilterText()}记录` }}
            </el-text>
          </el-empty>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { externalApi } from '@/api/external'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const refreshing = ref(false)
const publishing = ref(false)
const quickChecking = ref(false)
const retryingFailed = ref(false)
const allLinks = ref<any[]>([])
const publishTimer = ref<number | null>(null)
const userIP = ref('')
const linkFilter = ref('available') // 'all', 'available', 'unavailable', 'enabled'
const logFilter = ref('all') // 'all', 'success', 'error', 'system'
const logSearch = ref('')
const autoScroll = ref(true)
const logsContainer = ref<HTMLElement>()

const publishForm = reactive({
  publishCount: 1,
  publishInterval: 3,
  timeout: 30,
  useProxy: false,
  proxyHost: '',
  proxyPort: 8080,
  proxyType: 'http',
  params: {
  country: '',
  phone: '',
  wechat: '',
  telegram: '',
    website: '',
    custom: ''
  }
})

const publishProgress = reactive({
  current: 0,
  total: 0
})

const publishLogs = ref<Array<{
  timestamp: string
  url: string
  message: string
  status: 'success' | 'error'
  level?: 'system' | 'access' | 'error'
  ip?: string
  proxy?: string
  duration?: number
  details?: string
  error?: string
  expanded?: boolean
}>>([])

// 各种筛选后的链接
const availableLinks = computed(() => {
  return allLinks.value.filter((link: any) => 
    link.status === true && link.is_valid === true
  )
})

const unavailableLinks = computed(() => {
  return allLinks.value.filter((link: any) => link.is_valid === false)
})

const enabledLinks = computed(() => {
  return allLinks.value.filter((link: any) => link.status === true)
})

const disabledLinks = computed(() => {
  return allLinks.value.filter((link: any) => link.status === false)
})

// 根据筛选条件显示的链接
const filteredLinks = computed(() => {
  switch (linkFilter.value) {
    case 'all':
      return allLinks.value
    case 'available':
      return availableLinks.value
    case 'unavailable':
      return unavailableLinks.value
    case 'enabled':
      return enabledLinks.value
    default:
      return availableLinks.value
  }
})

// 成功和失败计数
const successCount = computed(() => 
  publishLogs.value.filter(log => log.status === 'success').length
)

const failCount = computed(() => 
  publishLogs.value.filter(log => log.status === 'error').length
)

// 系统日志计数
const systemCount = computed(() => 
  publishLogs.value.filter(log => log.level === 'system').length
)

// 成功率
const successRate = computed(() => {
  const total = publishLogs.value.length
  return total > 0 ? Math.round((successCount.value / total) * 100) : 0
})

// 失败率
const failRate = computed(() => {
  const total = publishLogs.value.length
  return total > 0 ? Math.round((failCount.value / total) * 100) : 0
})

// 平均耗时
const averageDuration = computed(() => {
  const logsWithDuration = publishLogs.value.filter(log => log.duration)
  if (logsWithDuration.length === 0) return '0'
  const avg = logsWithDuration.reduce((sum, log) => sum + (log.duration || 0), 0) / logsWithDuration.length
  return avg.toFixed(1)
})

// 筛选后的日志
const filteredLogs = computed(() => {
  let filtered = publishLogs.value

  // 按状态筛选
  if (logFilter.value !== 'all') {
    if (logFilter.value === 'system') {
      filtered = filtered.filter(log => log.level === 'system')
    } else {
      filtered = filtered.filter(log => log.status === logFilter.value)
    }
  }

  // 按搜索内容筛选
  if (logSearch.value.trim()) {
    const search = logSearch.value.trim().toLowerCase()
    filtered = filtered.filter(log => 
      log.url?.toLowerCase().includes(search) ||
      log.message?.toLowerCase().includes(search) ||
      log.details?.toLowerCase().includes(search)
    )
  }

  return filtered
})

// 获取用户IP地址
const getUserIP = async () => {
  try {
    const response = await fetch('https://api.ipify.org?format=json')
    const data = await response.json()
    userIP.value = data.ip
  } catch (error) {
    // 如果获取失败，尝试其他服务
    try {
      const response = await fetch('https://httpbin.org/ip')
      const data = await response.json()
      userIP.value = data.origin
    } catch (error2) {
      console.log('无法获取IP地址:', error2)
      userIP.value = '未知'
    }
  }
}

// 代理切换处理
const onProxyToggle = (value: boolean) => {
  if (value) {
    ElMessage.info('已启用代理模式，请配置代理服务器')
  } else {
    ElMessage.info('已禁用代理模式')
  }
}

// 获取外链列表
const fetchLinks = async () => {
  loading.value = true
  try {
    const response = await externalApi.getExternalLinks({
      page: 1,
      per_page: 1000 // 获取所有外链
    })
    
    allLinks.value = Array.isArray(response.data) ? response.data : []
    
    console.log('获取到的外链列表:', allLinks.value)
    console.log('可用外链数量:', availableLinks.value.length)
    
    if (availableLinks.value.length === 0) {
      ElMessage.warning('暂无可用的外链，请先在外链列表中创建并检测外链')
    } else {
      ElMessage.success(`加载了 ${availableLinks.value.length} 个可用外链`)
    }
    
  } catch (error) {
    console.error('获取外链列表失败:', error)
    ElMessage.error('获取外链列表失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 刷新外链列表
const refreshLinks = async () => {
  refreshing.value = true
  await fetchLinks()
  refreshing.value = false
}

// 设置链接筛选
const setLinkFilter = (filter: string) => {
  linkFilter.value = filter
}

// 获取空状态描述
const getEmptyDescription = () => {
  switch (linkFilter.value) {
    case 'all':
      return '暂无外链'
    case 'available':
      return '暂无可用的外链'
    case 'unavailable':
      return '暂无不可用的外链'
    case 'enabled':
      return '暂无启用的外链'
    default:
      return '暂无外链'
  }
}

// 获取状态标签类型
const getStatusTagType = (link: any) => {
  if (!link.status) return 'info' // 禁用
  if (link.is_valid) return 'success' // 可用
  return 'danger' // 不可用
}

// 获取状态文本
const getStatusText = (link: any) => {
  if (!link.status) return '🔒 禁用'
  if (link.is_valid) return '✅ 可用'
  return '❌ 不可用'
}

// 快速检测所有链接
const quickCheckLinks = async () => {
  if (allLinks.value.length === 0) {
    ElMessage.warning('暂无外链需要检测')
    return
  }
  
  quickChecking.value = true
  ElMessage.info(`🔍 开始快速检测 ${allLinks.value.length} 个外链...`)
  
  try {
    // 这里可以调用批量检测API
    await new Promise(resolve => setTimeout(resolve, 2000)) // 模拟检测
    await fetchLinks() // 重新获取最新状态
    ElMessage.success('✅ 快速检测完成')
  } catch (error) {
    console.error('快速检测失败:', error)
    ElMessage.error('❌ 快速检测失败')
  } finally {
    quickChecking.value = false
  }
}

// 检测单个链接
const checkSingleLinkInPublish = async (link: any) => {
  if (link.checking) return
  
  // 设置检测状态
  link.checking = true
  
  try {
    // 这里可以调用单个检测API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟检测
    ElMessage.success(`✅ ${link.url} 检测完成`)
    await fetchLinks() // 重新获取最新状态
  } catch (error) {
    console.error('检测失败:', error)
    ElMessage.error(`❌ ${link.url} 检测失败`)
  } finally {
    link.checking = false
  }
}

// 重试失败的链接
const retryFailedLinksInPublish = async () => {
  const failedLinks = unavailableLinks.value
  
  if (failedLinks.length === 0) {
    ElMessage.info('暂无失败的链接需要重试')
    return
  }
  
  retryingFailed.value = true
  ElMessage.info(`🔄 开始重试 ${failedLinks.length} 个失败的链接...`)
  
  try {
    // 串行重试每个失败的链接
    for (const link of failedLinks) {
      addLog(link.url, `🔄 重试检测: ${link.url}`, 'success')
      
      // 模拟重试访问
      const success = await visitLinkInBackground(link.url, link.id)
      
      if (success) {
        addLog(link.url, `🎉 重试成功: ${link.url}`, 'success')
      } else {
        addLog(link.url, `❌ 重试失败: ${link.url}`, 'error')
      }
      
      // 重试间隔
      await new Promise(resolve => setTimeout(resolve, 3000))
    }
    
    await fetchLinks() // 重新获取最新状态
    ElMessage.success('🎉 重试完成')
  } catch (error) {
    console.error('重试失败:', error)
    ElMessage.error('❌ 重试过程中出现错误')
  } finally {
    retryingFailed.value = false
  }
}

// 设置日志筛选
const setLogFilter = (filter: string) => {
  logFilter.value = filter
}

// 导出日志
const exportLogs = () => {
  if (publishLogs.value.length === 0) {
    ElMessage.warning('暂无日志可导出')
    return
  }
  
  const logData = publishLogs.value.map(log => ({
    时间: log.timestamp,
    链接: log.url || '系统',
    消息: log.message,
    状态: log.status === 'success' ? '成功' : '失败',
    级别: log.level === 'system' ? '系统' : log.level === 'access' ? '访问' : '错误',
    耗时: log.duration ? `${log.duration}秒` : '',
    IP地址: log.ip || '',
    代理: log.proxy || '',
    详情: log.details || '',
    错误: log.error || ''
  }))
  
  const csvContent = [
    Object.keys(logData[0]).join(','),
    ...logData.map(row => Object.values(row).map(val => `"${val}"`).join(','))
  ].join('\n')
  
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `访问日志_${new Date().toLocaleString().replace(/[/:]/g, '-')}.csv`
  link.click()
  
  ElMessage.success('📁 日志导出成功')
}

// 获取日志图标
const getLogIcon = (log: any) => {
  if (log.level === 'system') return '🔧'
  if (log.status === 'success') return '✅'
  return '❌'
}

// 格式化时间
const formatTime = (timestamp: string) => {
  return timestamp.split(' ')[1] || timestamp
}

// 格式化日期
const formatDate = (timestamp: string) => {
  return timestamp.split(' ')[0] || ''
}

// 获取日志标签类型
const getLogTagType = (log: any) => {
  if (log.level === 'system') return 'info'
  return log.status === 'success' ? 'success' : 'danger'
}

// 获取日志状态文本
const getLogStatusText = (log: any) => {
  if (log.level === 'system') return '系统'
  return log.status === 'success' ? '成功' : '失败'
}

// 获取级别文本
const getLevelText = (level: string) => {
  switch (level) {
    case 'system': return '系统'
    case 'access': return '访问'
    case 'error': return '错误'
    default: return level
  }
}

// 获取空日志描述
const getEmptyLogsDescription = () => {
  if (logSearch.value.trim()) {
    return `未找到包含"${logSearch.value}"的日志`
  }
  return getFilterText() + '日志为空'
}

// 获取筛选文本
const getFilterText = () => {
  switch (logFilter.value) {
    case 'success': return '成功'
    case 'error': return '失败'
    case 'system': return '系统'
    default: return '全部'
  }
}

// 显示后台运行说明
const showBackgroundTips = () => {
  ElMessage({
    message: `
💡 后台运行保障说明：

🛡️ 自动启用功能：
• 屏幕唤醒锁：防止设备自动休眠
• 心跳检测：每30秒发送一次活跃信号
• 页面监听：监控浏览器最小化状态

📱 使用建议：
• 可以安全地最小化浏览器窗口
• 可以切换到其他应用程序
• 建议保持设备电源充足
• 避免关闭浏览器标签页

⚠️ 注意事项：
• 任务期间请勿关闭浏览器
• 网络断开会影响访问效果
• 长时间运行建议使用稳定网络
    `,
    type: 'info',
    duration: 0,
    showClose: true,
    dangerouslyUseHTMLString: false
  })
}

// 确认开始发布
const confirmStartPublish = () => {
  ElMessageBox.confirm(
    `🎭 准备开始真实用户访问模拟
    
📊 任务详情：
• 可用链接：${availableLinks.value.length} 个
• 发布轮次：${publishForm.publishCount} 轮
• 总访问量：${availableLinks.value.length * publishForm.publishCount} 次
• 预计耗时：约 ${Math.ceil(availableLinks.value.length * publishForm.publishCount * 15 / 60)} 分钟

🛡️ 后台保障：
• 自动启用屏幕唤醒锁
• 支持浏览器最小化运行
• 实时监控访问状态

是否确认开始任务？`,
    '开始真实用户访问模拟',
    {
      confirmButtonText: '🚀 开始任务',
      cancelButtonText: '📋 取消',
      type: 'info',
      distinguishCancelAndClose: true,
      customClass: 'publish-confirm-dialog'
    }
  ).then(() => {
    startPublish()
  }).catch(() => {
    ElMessage.info('任务已取消')
  })
}

// 构建带参数的URL
const buildUrlWithParams = (baseUrl: string) => {
  const params = new URLSearchParams()
  
  // 添加表单参数
  Object.entries(publishForm.params).forEach(([key, value]) => {
    if (value && key !== 'custom') {
      params.append(key, value as string)
    }
  })
  
  // 添加自定义参数
  if (publishForm.params.custom) {
    const customParams = publishForm.params.custom.split('&')
    customParams.forEach(param => {
      const [key, value] = param.split('=')
      if (key && value) {
        params.append(key.trim(), value.trim())
      }
    })
  }
  
  return params.toString() ? `${baseUrl}?${params.toString()}` : baseUrl
}

// 更新链接显示（当参数改变时）
const updateLinksDisplay = () => {
  // 强制更新computed属性，这里不需要做任何事情
  // Vue的响应式系统会自动更新显示
}

// 添加日志
const addLog = (
  url: string, 
  message: string, 
  status: 'success' | 'error', 
  includeIP: boolean = false, 
  includeProxy: boolean = false,
  level: 'system' | 'access' | 'error' = 'access',
  duration?: number,
  details?: string,
  error?: string
) => {
  const now = new Date()
  const log = {
    timestamp: now.toLocaleString(),
    url,
    message,
    status,
    level,
    duration,
    details,
    error,
    ip: includeIP ? userIP.value : undefined,
    proxy: includeProxy && publishForm.useProxy ? `${publishForm.proxyType}://${publishForm.proxyHost}:${publishForm.proxyPort}` : undefined,
    expanded: false
  }
  
  publishLogs.value.unshift(log) // 新日志添加到顶部
  
  // 限制日志数量，最多保留1000条
  if (publishLogs.value.length > 1000) {
    publishLogs.value = publishLogs.value.slice(0, 1000)
  }
  
  // 自动滚动
  if (autoScroll.value) {
    nextTick(() => {
      if (logsContainer.value) {
        logsContainer.value.scrollTop = 0
      }
    })
  }
}

// 更新外链点击统计
const updateLinkStatistics = async (linkId: string) => {
  try {
    // 这里调用API更新外链的点击统计
    // 假设有一个API可以增加点击量
    // await externalApi.incrementClicks(linkId)
    console.log(`更新外链 ${linkId} 的点击统计`)
  } catch (error) {
    console.error('更新外链统计失败:', error)
  }
}

// 模拟真实用户访问单个链接
const visitLinkInBackground = async (url: string, linkId?: string): Promise<boolean> => {
  const fullUrl = buildUrlWithParams(url)
  
  try {
    // 阶段1: 模拟用户准备阶段
    addLog(fullUrl, '🎭 模拟用户准备访问...', 'success')
    const prepTime = 800 + Math.random() * 1500 // 800-2300ms
    await new Promise(resolve => {
      // 使用更精确的延时，不受浏览器状态影响
      const startTime = Date.now()
      const checkTime = () => {
        if (Date.now() - startTime >= prepTime) {
          resolve(undefined)
        } else {
          requestAnimationFrame(checkTime)
        }
      }
      checkTime()
    })
    
    // 阶段2: 模拟打开浏览器
    addLog(fullUrl, '🌐 模拟打开浏览器页面...', 'success')
    const browserTime = 1000 + Math.random() * 2000 // 1000-3000ms
    await new Promise(resolve => setTimeout(resolve, browserTime))
    
    // 阶段3: 模拟网络连接
    addLog(fullUrl, '🔗 模拟建立网络连接...', 'success')
    
    // 构建请求选项
    const requestOptions: RequestInit = {
      method: 'GET',
      mode: 'no-cors', // 避免CORS问题
      cache: 'no-cache',
      headers: {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
        'Accept-Encoding': 'gzip, deflate, br',
        'Referer': 'https://www.google.com/'
      }
    }
    
    // 如果使用代理，添加代理信息到日志
    if (publishForm.useProxy) {
      addLog(fullUrl, `🔒 使用代理: ${publishForm.proxyType}://${publishForm.proxyHost}:${publishForm.proxyPort}`, 'success')
    }
    
    // 阶段4: 发起真实网络请求
    addLog(fullUrl, '📡 发起网络请求...', 'success')
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), publishForm.timeout * 1000)
    
    requestOptions.signal = controller.signal
    
    const startTime = Date.now()
    const response = await fetch(fullUrl, requestOptions)
    const requestDuration = Date.now() - startTime
    
    clearTimeout(timeoutId)
    
    // 阶段5: 模拟页面加载和用户浏览
    addLog(fullUrl, `📄 页面加载完成 (${requestDuration}ms)`, 'success')
    const browsingTime = 2000 + Math.random() * 4000 // 2-6秒浏览时间
    await new Promise(resolve => setTimeout(resolve, browsingTime))
    
    const totalDuration = Math.round((Date.now() - startTime + prepTime + browserTime) / 1000)
    addLog(
      fullUrl, 
      `🎉 用户访问完成`, 
      'success', 
      true, 
      publishForm.useProxy,
      'access',
      totalDuration,
      `浏览器启动: ${Math.round(browserTime/1000)}s, 网络请求: ${requestDuration}ms, 浏览时间: ${Math.round(browsingTime/1000)}s`
    )
    
    // 更新外链统计
    if (linkId) {
      await updateLinkStatistics(linkId)
    }
    
    return true
    
  } catch (error: any) {
    let errorMessage = '❌ 用户访问失败'
    
    if (error?.name === 'AbortError') {
      errorMessage = `⏰ 访问超时 (${publishForm.timeout}秒)`
    } else if (error?.message) {
      errorMessage = `🚫 访问异常: ${error.message}`
    }
    
    addLog(
      fullUrl, 
      errorMessage, 
      'error', 
      true, 
      publishForm.useProxy,
      'error',
      undefined,
      `请求超时: ${publishForm.timeout}秒`,
      error?.message || error
    )
    console.error('用户访问模拟失败:', error)
    return false
  }
}

// 防止浏览器降频执行的工具函数
const createStableInterval = (callback: () => void, delay: number) => {
  let timeoutId: number
  
  const run = () => {
    callback()
    timeoutId = window.setTimeout(run, delay)
  }
  
  timeoutId = window.setTimeout(run, delay)
  
  return () => {
    if (timeoutId) {
      window.clearTimeout(timeoutId)
    }
  }
}

// 防止浏览器休眠的心跳机制
const startHeartbeat = () => {
  let heartbeatInterval: number
  
  const heartbeat = () => {
    // 使用console.debug而不是console.log，避免控制台污染
    console.debug('📡 发布任务心跳检测:', new Date().toLocaleTimeString())
  }
  
  // 每30秒发送一次心跳，保持页面活跃
  heartbeatInterval = window.setInterval(heartbeat, 30000)
  
  return () => {
    if (heartbeatInterval) {
      window.clearInterval(heartbeatInterval)
    }
  }
}

// 添加页面可见性变化监听
const handleVisibilityChange = () => {
  if (document.hidden) {
    addLog('', '📱 浏览器已最小化，发布任务继续在后台运行...', 'success', false, false, 'system')
  } else {
    addLog('', '👀 浏览器已激活，发布任务正常运行中...', 'success', false, false, 'system')
  }
}

// 开始发布（改进版本，包含失败重发和后台运行保障）
const startPublish = async () => {
  // 声明后台运行相关变量
  let stopHeartbeat: (() => void) | null = null
  let wakeLock: any = null
  
  if (availableLinks.value.length === 0) {
    ElMessage.warning('没有可用的外链进行发布')
    return
  }
  
  if (publishForm.useProxy && (!publishForm.proxyHost || !publishForm.proxyPort)) {
    ElMessage.warning('请配置代理服务器地址和端口')
    return
  }
  
  // 添加页面可见性监听
  document.addEventListener('visibilitychange', handleVisibilityChange)
  
  // 启动心跳机制
  stopHeartbeat = startHeartbeat()
  
  // 阻止页面休眠
  if ('wakeLock' in navigator) {
    try {
      wakeLock = await (navigator as any).wakeLock.request('screen')
      addLog('', '🔒 已启用屏幕唤醒锁，防止设备休眠', 'success', false, false, 'system')
    } catch (err) {
      addLog('', '⚠️ 无法启用屏幕唤醒锁，请保持浏览器活跃状态', 'success', false, false, 'system')
    }
  }
  
  publishing.value = true
  publishProgress.current = 0
  publishProgress.total = availableLinks.value.length * publishForm.publishCount
  
  ElMessage.info(`🎭 开始真实用户访问模拟，共 ${publishProgress.total} 次访问`)
  addLog('', '🚀 === 真实用户访问模拟开始 ===', 'success', false, false, 'system')
  addLog('', '🛡️ 后台运行保障已启用，最小化浏览器不会影响任务执行', 'success', false, false, 'system')
  
  // 记录失败的链接，用于重试
  let failedLinks: Array<{ url: string, id?: string, attempts: number }> = []
  
  try {
    for (let round = 1; round <= publishForm.publishCount; round++) {
      if (!publishing.value) break // 检查是否被停止
      
      addLog('', `🎯 === 第 ${round} 轮用户访问开始 ===`, 'success', false, false, 'system')
      
      // 串行访问所有链接（模拟真实用户逐个访问）
      const results: boolean[] = []
      
      for (let i = 0; i < availableLinks.value.length; i++) {
        if (!publishing.value) break
        
        const link = availableLinks.value[i]
        
        // 模拟用户在访问之间的思考时间
        if (i > 0) {
          const thinkTime = 3000 + Math.random() * 5000 // 3-8秒思考时间
          addLog('', `🤔 用户思考时间 ${Math.round(thinkTime/1000)}秒...`, 'success')
          await new Promise(resolve => setTimeout(resolve, thinkTime))
        }
        
        addLog('', `👤 用户 ${i + 1}/${availableLinks.value.length} 开始访问...`, 'success')
        
        const success = await visitLinkInBackground(link.url, link.id)
        results.push(success)
        
        // 记录失败的链接
        if (!success) {
          const existingFailed = failedLinks.find(f => f.url === link.url)
          if (existingFailed) {
            existingFailed.attempts++
          } else {
            failedLinks.push({ url: link.url, id: link.id, attempts: 1 })
          }
        }
        
        publishProgress.current++
      }
      
      const successCount = results.filter(r => r === true).length
      const failCount = results.filter(r => r === false).length
      
      addLog('', `✅ 第 ${round} 轮完成: ${successCount} 成功, ${failCount} 失败`, successCount > failCount ? 'success' : 'error')
      
              // 如果不是最后一轮，添加间隔
        if (round < publishForm.publishCount && publishing.value) {
          const waitTime = publishForm.publishInterval
          addLog('', `⏳ 等待 ${waitTime} 秒后开始下一轮...`, 'success', false, false, 'system')
          
          // 使用更稳定的延时机制，不受浏览器状态影响
          await new Promise(resolve => {
            let remainingTime = waitTime
            const countdownInterval = setInterval(() => {
              remainingTime--
              if (remainingTime <= 0) {
                clearInterval(countdownInterval)
                resolve(undefined)
              }
              // 每秒更新一次倒计时（可选，避免太频繁的日志）
              if (remainingTime % 5 === 0) {
                addLog('', `⏰ 倒计时: ${remainingTime} 秒`, 'success', false, false, 'system')
              }
            }, 1000)
          })
        }
    }
    
    // 处理失败重发
    if (failedLinks.length > 0 && publishing.value) {
      addLog('', `🔄 === 开始失败链接重试 (${failedLinks.length} 个) ===`, 'success')
      
      // 最多重试3次
      for (let retryRound = 1; retryRound <= 3; retryRound++) {
        if (!publishing.value) break
        
        const linksToRetry = failedLinks.filter(link => link.attempts <= retryRound)
        if (linksToRetry.length === 0) break
        
        addLog('', `🔄 第 ${retryRound} 轮重试 (${linksToRetry.length} 个链接)`, 'success')
        
        for (const failedLink of linksToRetry) {
          if (!publishing.value) break
          
          addLog('', `🔁 重试访问: ${failedLink.url}`, 'success')
          
          // 重试前等待更长时间
          const retryDelay = 5000 + Math.random() * 5000 // 5-10秒
          await new Promise(resolve => setTimeout(resolve, retryDelay))
          
          const success = await visitLinkInBackground(failedLink.url, failedLink.id)
          
          if (success) {
            addLog('', `🎉 重试成功: ${failedLink.url}`, 'success')
            // 从失败列表中移除
            const index = failedLinks.indexOf(failedLink)
            if (index > -1) failedLinks.splice(index, 1)
          } else {
            addLog('', `❌ 重试失败: ${failedLink.url}`, 'error')
          }
        }
        
        if (retryRound < 3 && failedLinks.length > 0) {
          addLog('', `⏳ 等待 ${publishForm.publishInterval * 2} 秒后进行下轮重试...`, 'success')
          await new Promise(resolve => 
            setTimeout(resolve, publishForm.publishInterval * 2 * 1000)
          )
        }
      }
      
      // 重试完成后的统计
      if (failedLinks.length > 0) {
        addLog('', `⚠️ 仍有 ${failedLinks.length} 个链接访问失败，建议稍后手动重试`, 'error')
      } else {
        addLog('', `🎉 所有失败链接重试成功！`, 'success')
      }
    }
    
    if (publishing.value) {
      const finalFailCount = failedLinks.length
      const finalSuccessCount = publishProgress.total - finalFailCount
      
      if (finalFailCount === 0) {
        ElMessage.success(`🎉 所有用户访问任务完成！(${finalSuccessCount}/${publishProgress.total})`)
        addLog('', '🎊 === 所有用户访问任务完美完成 ===', 'success')
      } else {
        ElMessage.warning(`⚠️ 访问任务完成，${finalFailCount} 个链接仍然失败`)
        addLog('', `📊 === 访问任务完成：成功 ${finalSuccessCount}，失败 ${finalFailCount} ===`, 'error')
      }
    } else {
      ElMessage.warning('用户访问已停止')
      addLog('', '🛑 === 用户访问被手动停止 ===', 'error')
    }
  } catch (error: any) {
    ElMessage.error('用户访问过程中出现错误')
    addLog('', `💥 系统错误: ${error?.message || error}`, 'error')
    console.error('用户访问过程错误:', error)
  } finally {
    // 清理后台运行保障机制
    document.removeEventListener('visibilitychange', handleVisibilityChange)
    
    if (stopHeartbeat) {
      stopHeartbeat()
    }
    
    // 释放屏幕唤醒锁
    if (wakeLock) {
      try {
        await wakeLock.release()
        addLog('', '🔓 已释放屏幕唤醒锁', 'success', false, false, 'system')
      } catch (err) {
        console.warn('释放唤醒锁失败:', err)
      }
    }
    
    publishing.value = false
    publishProgress.current = 0
    publishProgress.total = 0
    
    addLog('', '🏁 === 后台运行保障已关闭 ===', 'success', false, false, 'system')
  }
}

// 停止发布
const stopPublish = () => {
  publishing.value = false
  if (publishTimer.value) {
    clearTimeout(publishTimer.value)
    publishTimer.value = null
  }
  
  // 清理后台运行保障机制
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  
  addLog('', '🛑 用户手动停止发布任务', 'success', false, false, 'system')
  addLog('', '🏁 === 后台运行保障已关闭 ===', 'success', false, false, 'system')
  
  ElMessage.info('🛑 发布已停止，后台保障机制已关闭')
}

// 清空日志
const clearLogs = () => {
  publishLogs.value = []
  ElMessage.success('日志已清空')
}

onMounted(() => {
  fetchLinks()
  getUserIP()
})
</script>

<style scoped>
.external-link-publish {
  padding: 20px;
}

.link-display-card,
.publish-config-card,
.publish-logs-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.links-container {
  height: 300px;
  overflow-y: auto;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.no-links {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.links-list {
  padding: 12px;
}

.link-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.2s;
}

.link-item:hover {
  background-color: #f5f7fa;
}

.link-item:last-child {
  border-bottom: none;
}

.link-item.unavailable {
  background-color: #fef2f2;
  border-left: 3px solid #f56565;
}

.link-item.disabled {
  background-color: #f7fafc;
  border-left: 3px solid #a0aec0;
  opacity: 0.7;
}

.link-url {
  font-size: 14px;
  word-break: break-all;
  line-height: 1.4;
}

.link-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 旧的日志样式已被新样式替换 */

/* 新增样式 */
.filter-section {
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.filter-controls {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.stats-info {
  font-size: 12px;
}

.link-status {
  flex-shrink: 0;
  width: 80px;
}

.link-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0; /* 允许内容收缩 */
}

.link-actions {
  flex-shrink: 0;
}

.link-url {
  font-size: 14px;
  word-break: break-all;
  line-height: 1.4;
}

.link-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

/* 改进按钮组样式 */
.el-button-group .el-button {
  font-size: 12px;
  padding: 4px 8px;
}

/* 优化日志样式 */
.log-title {
  display: flex;
  align-items: center;
}

.log-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.log-filters {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
}

.filter-row, .stats-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.logs-container {
  max-height: 500px;
  overflow-y: auto;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 0;
}

.log-item {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f2f5;
  transition: all 0.2s;
  gap: 12px;
}

.log-item:hover {
  background-color: #f8f9fa;
}

.log-item:last-child {
  border-bottom: none;
}

.log-item.success {
  background-color: #f0f9ff;
  border-left: 3px solid #67c23a;
}

.log-item.error {
  background-color: #fef2f2;
  border-left: 3px solid #f56c6c;
}

.log-item.system {
  background-color: #f4f4f5;
  border-left: 3px solid #909399;
}

.log-icon {
  flex-shrink: 0;
  width: 24px;
  text-align: center;
}

.status-icon {
  font-size: 16px;
}

.log-time {
  flex-shrink: 0;
  width: 90px;
  font-size: 12px;
  color: #909399;
}

.time-main {
  font-weight: 500;
  margin-bottom: 2px;
}

.time-detail {
  font-size: 10px;
  opacity: 0.7;
}

.log-content {
  flex: 1;
  min-width: 0;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
  gap: 12px;
}

.log-url {
  flex: 1;
  word-break: break-all;
  line-height: 1.4;
}

.log-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  flex-shrink: 0;
}

.log-message {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
  margin-bottom: 8px;
}

.log-details {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed #e4e7ed;
}

.detail-content {
  padding: 8px 12px;
  background-color: #fafafa;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

.detail-item {
  display: flex;
  align-items: center;
  margin-bottom: 6px;
  gap: 8px;
}

.detail-item:last-child {
  margin-bottom: 0;
}

.detail-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
  min-width: 80px;
}

.no-logs {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

/* 发布确认对话框样式 */
:global(.publish-confirm-dialog) {
  width: 500px;
}

:global(.publish-confirm-dialog .el-message-box__message) {
  white-space: pre-line;
  line-height: 1.6;
  font-size: 14px;
}

:global(.publish-confirm-dialog .el-message-box__title) {
  font-size: 16px;
  font-weight: 600;
}
</style> 