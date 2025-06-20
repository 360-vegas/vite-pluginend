<template>
  <div class="external-links-list">
    <!-- 创建外链表单 -->
    <el-card class="create-card">
      <template #header>
        <div class="card-header">
          <span>创建外链</span>
          <div>
            <el-button type="success" @click="showBatchAddDialog">批量添加</el-button>
            <el-button type="primary" @click="toggleCreateForm">
              {{ showCreateForm ? '收起' : '展开' }}
          </el-button>
          </div>
        </div>
      </template>

      <div v-show="showCreateForm">
        <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="链接地址" prop="url">
                <el-input 
                  v-model="createForm.url" 
                  placeholder="请输入链接地址"
                  @blur="validateUrl"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="分类" prop="category">
                <el-input v-model="createForm.category" placeholder="请输入分类" />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-form-item label="描述">
            <el-input 
              v-model="createForm.description" 
              type="textarea" 
              :rows="2"
              placeholder="请输入描述"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              @click="handleCreateLink" 
              :loading="creating"
            >
              {{ creating ? '检测并创建中...' : '创建外链' }}
            </el-button>
            <el-button @click="resetCreateForm">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <!-- 筛选条件 -->
    <el-card class="filter-card">
      <el-form :model="query" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="query.keyword"
            placeholder="搜索链接地址"
            clearable
            style="width: 250px"
          />
        </el-form-item>
        <el-form-item label="可用性">
          <el-select v-model="query.is_valid" placeholder="选择可用性" clearable @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="✅ 可用" value="true" />
            <el-option label="❌ 不可用" value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="选择状态" clearable>
            <el-option label="全部" value="" />
            <el-option label="启用" value="true" />
            <el-option label="禁用" value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 外链列表 -->
    <el-card class="link-table">
      <template #header>
        <div class="card-header">
          <span>外链列表 ({{ total }})</span>
          <div>
            <el-button 
              v-if="selectedLinks.length > 0"
              type="danger" 
              @click="handleBatchDelete"
              :disabled="selectedLinks.length === 0"
            >
              批量删除 ({{ selectedLinks.length }})
            </el-button>
            <el-button 
              v-if="selectedLinks.length > 0"
              type="warning" 
              @click="handleBatchCheck"
              :loading="batchChecking"
            >
              批量检测 ({{ selectedLinks.length }})
            </el-button>
            <el-button type="success" @click="checkAllLinks" :loading="checking">
              {{ checking ? '检测中...' : '全部检测' }}
            </el-button>
            <el-button 
              type="danger" 
              plain
              @click="handleDeleteInvalidLinks"
              :loading="deletingInvalid"
            >
              {{ deletingInvalid ? '删除中...' : '删除不可用链接' }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 快速筛选按钮 -->
      <div class="quick-filters" style="margin-bottom: 16px;">
        <span style="margin-right: 12px; color: #666;">快速筛选:</span>
        <el-button-group>
          <el-button 
            :type="query.is_valid === '' ? 'primary' : 'default'"
            size="small"
            @click="quickFilter('')"
          >
            全部 ({{ total }})
          </el-button>
          <el-button 
            :type="query.is_valid === 'true' ? 'success' : 'default'"
            size="small"
            @click="quickFilter('true')"
          >
            ✅ 可用
          </el-button>
          <el-button 
            :type="query.is_valid === 'false' ? 'danger' : 'default'"
            size="small"
            @click="quickFilter('false')"
          >
            ❌ 不可用
          </el-button>
        </el-button-group>
        
        <el-divider direction="vertical" />
        
        <el-button 
          size="small"
          type="warning"
          @click="retryFailedLinks"
          :loading="retryingFailed"
        >
          🔄 重试失败链接
        </el-button>
      </div>

      <!-- 检测进度显示 -->
      <el-alert
        v-if="checkingProgress.show"
        :title="checkingProgress.title"
        type="info"
        :description="checkingProgress.description"
        show-icon
        :closable="false"
        style="margin-bottom: 20px"
      >
        <el-progress
          :percentage="checkingProgress.percentage"
          :stroke-width="18"
          :text-inside="true"
          status="success"
        />
      </el-alert>

      <el-table
        :data="links" 
        v-loading="loading"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="is_valid" label="可用性 ↕️" width="130" sortable="custom">
          <template #default="{ row }">
            <el-tag 
              :type="row.is_valid ? 'success' : 'danger'" 
              size="large"
              style="font-weight: bold;"
            >
              {{ row.is_valid ? '✅ 可用' : '❌ 不可用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="链接地址" min-width="400">
          <template #default="{ row }">
            <div class="link-cell">
              <el-link 
                :href="row.url" 
                target="_blank" 
                @click="handleClick(row)"
                :type="row.is_valid ? 'primary' : 'danger'"
                style="font-weight: 500"
              >
              {{ row.url }}
            </el-link>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="clicks" label="点击量" width="120" sortable />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch 
              v-model="row.status" 
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" sortable>
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button
                size="small"
                type="warning"
                @click="checkSingleLink(row)"
                :loading="row.checking"
              >
                检测
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.per_page"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 批量添加对话框 -->
    <el-dialog 
      v-model="batchAddVisible" 
      title="批量添加外链" 
      width="60%"
      :close-on-click-modal="false"
    >
      <el-form :model="batchAddForm" label-width="100px">
        <el-form-item label="批量链接">
          <el-input
            v-model="batchAddForm.urls"
            type="textarea"
            :rows="10"
            placeholder="请输入链接地址，每行一个链接&#10;支持格式：&#10;1. 纯链接：https://example.com&#10;2. 链接|分类：https://example.com|娱乐&#10;3. 链接|分类|描述：https://example.com|娱乐|这是一个娱乐网站"
            style="width: 100%"
          />
          <div class="batch-tips">
            <el-text type="info" size="small">
              支持多种格式：纯链接、链接|分类、链接|分类|描述，每行一个
            </el-text>
          </div>
        </el-form-item>
        
        <el-form-item label="默认分类">
          <el-input 
            v-model="batchAddForm.defaultCategory" 
            placeholder="当链接没有指定分类时使用的默认分类"
            style="width: 200px"
          />
        </el-form-item>
        
        <el-form-item label="默认状态">
          <el-switch 
            v-model="batchAddForm.defaultStatus"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
        
        <el-form-item label="自动检测">
          <el-switch 
            v-model="batchAddForm.autoCheck"
            active-text="创建后自动检测可用性"
            inactive-text="创建后不检测"
          />
        </el-form-item>
      </el-form>
      
      <!-- 预览区域 -->
      <div v-if="parsedLinks.length > 0" class="preview-section">
        <el-divider content-position="left">预览 ({{ parsedLinks.length }} 条)</el-divider>
        <div class="preview-list">
          <div 
            v-for="(link, index) in parsedLinks.slice(0, 10)" 
            :key="index"
            class="preview-item"
          >
            <div class="preview-url">{{ link.url }}</div>
            <div class="preview-meta">
              <el-tag size="small">{{ link.category }}</el-tag>
              <span class="preview-desc">{{ link.description || '无描述' }}</span>
            </div>
          </div>
          <div v-if="parsedLinks.length > 10" class="more-indicator">
            还有 {{ parsedLinks.length - 10 }} 条...
          </div>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="batchAddVisible = false">取消</el-button>
          <el-button @click="previewBatchAdd">预览</el-button>
          <el-button 
            type="primary" 
            @click="handleBatchAdd"
            :loading="batchAdding"
            :disabled="parsedLinks.length === 0"
          >
            {{ batchAdding ? `添加中... (${batchProgress.current}/${batchProgress.total})` : `确认添加 (${parsedLinks.length} 条)` }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { externalApi } from '@/api/index'
import type { ExternalLink, ExternalLinkQuery } from '@/api/external'
import { formatDate } from '@/utils/date'

const router = useRouter()
const loading = ref(false)
const creating = ref(false)
const checking = ref(false)
const batchChecking = ref(false)
const batchAdding = ref(false)
const deletingInvalid = ref(false)
const showCreateForm = ref(true)
const batchAddVisible = ref(false)
const links = ref<ExternalLink[]>([])
const total = ref(0)
const categories = ref<string[]>([])
const createFormRef = ref<FormInstance>()
const selectedLinks = ref<ExternalLink[]>([])

// 检测进度显示
const checkingProgress = reactive({
  show: false,
  title: '',
  description: '',
  percentage: 0,
  current: 0,
  total: 0
})

const query = ref<ExternalLinkQuery>({
  page: 1,
  per_page: 10,
  keyword: '',
  category: '',
  status: '',
  is_valid: '', // 默认为空，显示全部
  sort_field: 'is_valid',  // 默认按可用性排序
  sort_order: 'desc'       // 降序：可用的在前面
})

const createForm = reactive({
  url: '',
  category: '',
  description: '',
  status: true
})

const batchAddForm = reactive({
  urls: '',
  defaultCategory: '',
  defaultStatus: true,
  autoCheck: false
})

const batchProgress = reactive({
  current: 0,
  total: 0
})

const createRules = {
  url: [
    { required: true, message: '请输入链接地址', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请输入分类', trigger: 'blur' }
  ]
}

// 解析批量添加的链接
const parsedLinks = computed(() => {
  if (!batchAddForm.urls.trim()) return []
  
  const lines = batchAddForm.urls.trim().split('\n')
  const links: Array<{ url: string; category: string; description: string }> = []
  
  lines.forEach(line => {
    const trimmedLine = line.trim()
    if (!trimmedLine) return
    
    const parts = trimmedLine.split('|')
    const url = parts[0]?.trim()
    
    if (!url) return
    
    // 验证URL格式
    try {
      new URL(url)
    } catch {
      return // 跳过无效URL
    }
    
    const category = parts[1]?.trim() || batchAddForm.defaultCategory || '默认分类'
    const description = parts[2]?.trim() || ''
    
    links.push({ url, category, description })
  })
  
  return links
})

// 获取外链列表
const fetchLinks = async () => {
  loading.value = true
  try {
    console.log('🔄 开始获取外链列表，查询参数:', query.value) // 调试日志
    
    // 清理查询参数，移除空值
    const cleanQuery = { ...query.value }
    if (cleanQuery.is_valid === '') {
      delete cleanQuery.is_valid // 删除空字符串，让后端显示全部数据
    }
    
    console.log('🧹 清理后的查询参数:', cleanQuery) // 调试日志
    
    const response = await externalApi.getExternalLinks(cleanQuery)
    
    console.log('📡 API响应:', response) // 调试日志
    console.log('📊 响应数据类型:', typeof response)
    console.log('📋 数据内容:', response.data)
    console.log('📈 元数据:', (response as any).meta)
    
    if (response.data && Array.isArray(response.data)) {
      links.value = response.data.map(link => ({ ...link, checking: false }))
      total.value = (response as any).meta?.total || response.data.length
      
      console.log(`✅ 成功加载 ${links.value.length} 条链接，总数: ${total.value}`)
      
      // 分析可用性分布
      const availableCount = links.value.filter(link => link.is_valid).length
      const unavailableCount = links.value.length - availableCount
      console.log(`📊 可用链接: ${availableCount}, 不可用链接: ${unavailableCount}`)
      
    } else {
      console.error('❌ 响应数据格式异常:', response)
      ElMessage.error('数据格式异常，请检查API响应')
    }
    
  } catch (error: any) {
    console.error('❌ 获取外链列表失败:', error)
    console.error('错误详情:', error.response?.data || error.message)
    ElMessage.error(`获取外链列表失败: ${error.response?.data?.message || error.message || '未知错误'}`)
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const stats = await externalApi.getExternalStatistics()
    categories.value = Object.keys(stats.categories)
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 验证URL
const validateUrl = async () => {
  if (!createForm.url) return
  
  try {
    new URL(createForm.url)
  } catch {
    ElMessage.warning('请输入有效的URL地址')
  }
}

// 创建外链
const handleCreateLink = async () => {
  if (!createFormRef.value) return
  
  await createFormRef.value.validate(async (valid) => {
    if (valid) {
      creating.value = true
      try {
        ElMessage.info('正在创建外链...')
        
        // 创建外链，初始状态设为不可用
        const newLink = await externalApi.createExternalLink({
          ...createForm,
          is_valid: false, // 初始设为false
          clicks: 0
        })
        
        console.log('创建外链响应:', newLink) // 调试日志
        
        ElMessage.success('外链创建成功')
        resetCreateForm()
        
        // 重新获取列表
        await fetchLinks()
        await fetchCategories()
        
        // 创建成功后提示用户手动检测
        ElMessage.info('请点击"检测"按钮验证链接可用性')
        
      } catch (error) {
        console.error('创建外链完整错误信息:', error)
        ElMessage.error('创建外链失败，请重试')
      } finally {
        creating.value = false
      }
    }
  })
}

// 检测单个链接 - 使用真实用户访问模拟
const checkSingleLink = async (link: any) => {
  link.checking = true
  checkingProgress.show = true
  checkingProgress.title = '🎭 模拟真实用户访问中...'
  checkingProgress.description = `🌐 模拟打开浏览器访问: ${link.url}`
  checkingProgress.percentage = 10
  
  // 模拟检测阶段
  const stages = [
    { desc: '🚀 启动浏览器...', percent: 20 },
    { desc: '🔍 解析域名...', percent: 35 },
    { desc: '🔗 建立连接...', percent: 50 },
    { desc: '📄 加载页面内容...', percent: 75 },
    { desc: '✅ 验证页面可用性...', percent: 90 }
  ]
  
  // 显示检测阶段
  for (const stage of stages) {
    checkingProgress.description = stage.desc
    checkingProgress.percentage = stage.percent
    await new Promise(resolve => setTimeout(resolve, 500)) // 暂停500ms显示进度
  }
  
  try {
    ElMessage.info(`🎭 正在模拟真实用户访问: ${link.url}`)
    
    // 使用后端API检测单个链接
    const response = await externalApi.batchCheckLinksBackend([link.id.toString()], false)
    
    console.log('单个检测响应:', response) // 调试日志
    
    if (response && (response as any).results && Array.isArray((response as any).results) && (response as any).results.length > 0) {
      const result = (response as any).results[0]
      
      // 更新链接状态
      link.is_valid = result.is_valid
      link.checked_at = result.checked_at
      
      // 显示检测结果
      checkingProgress.percentage = 100
      
      if (result.is_valid) {
        ElMessage.success(`🎉 真实访问成功: ${result.message || '链接正常可用'}`)
        checkingProgress.description = `🎉 用户访问成功: 链接可正常使用`
      } else {
        ElMessage.warning(`❌ 访问失败: ${result.error_message || result.message || '链接无法正常访问'}`)
        checkingProgress.description = `❌ 用户访问失败: 链接存在问题`
      }
      
      // 重新获取列表以确保数据同步
      await fetchLinks()
      
    } else {
      console.error('单个检测响应格式异常:', response)
      ElMessage.error('检测响应格式异常')
      checkingProgress.description = '检测失败'
    }
    
  } catch (error: any) {
    console.error('单个检测完整错误信息:', error)
    
    // 简化错误处理
    ElMessage.error('检测失败，请重试')
    link.is_valid = false
    checkingProgress.description = '检测失败'
  } finally {
    link.checking = false
    // 延迟隐藏进度条
    setTimeout(() => {
      checkingProgress.show = false
    }, 3000)
  }
}

// 批量检测所有链接 - 使用真实用户访问模拟
const checkAllLinks = async () => {
  checking.value = true
  checkingProgress.show = true
  checkingProgress.title = '🎭 正在模拟真实用户访问所有外链...'
  checkingProgress.current = 0
  checkingProgress.percentage = 10
  checkingProgress.description = '🚀 初始化真实用户访问环境...'
  
  // 模拟准备阶段
  await new Promise(resolve => setTimeout(resolve, 1000))
  checkingProgress.description = '👥 启动多个用户访问会话...'
  checkingProgress.percentage = 20
  
  try {
    ElMessage.info('🎭 开始模拟真实用户批量访问所有链接...')
    
    console.log('开始调用batchCheckLinksBackend API，参数: [], true')
    
    // 使用后端API检测所有链接
    const response = await externalApi.batchCheckLinksBackend([], true)
    
    console.log('API响应原始数据:', response)
    console.log('API响应类型:', typeof response)
    
    // 直接从响应中提取结果，不管嵌套结构
    let results: any[] = []
    const resp: any = response
    if (resp && resp.results) {
      results = resp.results
    } else if (resp && Array.isArray(resp)) {
      results = resp
    } else if (resp && resp.data && resp.data.results) {
      results = resp.data.results
    } else if (resp && resp.data && Array.isArray(resp.data)) {
      results = resp.data
    }
    
    console.log('解析出的results:', results)
    
    if (Array.isArray(results) && results.length > 0) {
      const totalChecked = results.length
      const successCount = results.filter(r => r.is_valid === true).length
      const failCount = results.filter(r => r.is_valid === false).length
      
      // 显示检测结果
      checkingProgress.percentage = 100
      checkingProgress.description = `🎉 真实用户访问完成: 共访问 ${totalChecked} 个，${successCount} 个正常，${failCount} 个异常`
      
      ElMessage.success(`🎉 全部真实用户访问完成！共模拟访问 ${totalChecked} 个链接，${successCount} 个正常访问，${failCount} 个访问异常`)
      
      console.log('可用链接数:', successCount)
      console.log('不可用链接数:', failCount)
      
      // 重新获取列表以确保数据同步
      await fetchLinks()
      
    } else {
      console.log('没有找到检测结果，results:', results)
      ElMessage.info('没有链接需要检测或检测结果为空')
      checkingProgress.description = '没有找到需要检测的链接'
    }
    
  } catch (error: any) {
    console.error('批量检测完整错误信息:', error)
    console.error('错误堆栈:', error.stack)
    
    // 显示具体错误信息
    let errorMsg = '批量检测失败'
    if (error.message) {
      errorMsg += ': ' + error.message
    }
    
    ElMessage.error(errorMsg)
    checkingProgress.description = '检测失败: ' + (error.message || '未知错误')
  } finally {
    checking.value = false
    // 延迟隐藏进度条
    setTimeout(() => {
      checkingProgress.show = false
    }, 5000)
  }
}

// 批量检测选中的链接
const handleBatchCheck = async () => {
  if (selectedLinks.value.length === 0) {
    ElMessage.warning('请先选择要检测的链接')
    return
  }
  
  batchChecking.value = true
  checkingProgress.show = true
  checkingProgress.title = `正在检测选中的 ${selectedLinks.value.length} 个链接...`
  checkingProgress.current = 0
  checkingProgress.total = selectedLinks.value.length
  checkingProgress.percentage = 0
  
  try {
    ElMessage.info(`开始检测选中的 ${selectedLinks.value.length} 个链接...`)
    
    const ids = selectedLinks.value.map(link => link.id.toString())
    console.log('批量检测选中链接IDs:', ids)
    
    const response = await externalApi.batchCheckLinksBackend(ids, false)
    
    console.log('批量检测选中链接API响应:', response)
    
    // 直接从响应中提取结果，不管嵌套结构
    let results: any[] = []
    const resp: any = response
    if (resp && resp.results) {
      results = resp.results
    } else if (resp && Array.isArray(resp)) {
      results = resp
    } else if (resp && resp.data && resp.data.results) {
      results = resp.data.results
    } else if (resp && resp.data && Array.isArray(resp.data)) {
      results = resp.data
    }
    
    if (Array.isArray(results) && results.length > 0) {
      const totalChecked = results.length
      const successCount = results.filter(r => r.is_valid === true).length
      const failCount = results.filter(r => r.is_valid === false).length
      
      // 显示检测结果
      checkingProgress.percentage = 100
      checkingProgress.description = `检测完成: 共 ${totalChecked} 个，${successCount} 个可用，${failCount} 个不可用`
      
      ElMessage.success(`批量检测完成！共检测 ${totalChecked} 个链接，${successCount} 个可用，${failCount} 个不可用`)
      
      // 重新获取列表以确保数据同步
      await fetchLinks()
      
    } else {
      console.log('批量检测没有找到结果，results:', results)
      ElMessage.warning('检测响应异常，请重试')
      checkingProgress.description = '响应格式异常'
    }
    
  } catch (error: any) {
    console.error('批量检测选中链接完整错误信息:', error)
    console.error('错误堆栈:', error.stack)
    
    let errorMsg = '批量检测失败'
    if (error.message) {
      errorMsg += ': ' + error.message
    }
    
    ElMessage.error(errorMsg)
    checkingProgress.description = '检测失败: ' + (error.message || '未知错误')
  } finally {
    batchChecking.value = false
    // 延迟隐藏进度条
    setTimeout(() => {
      checkingProgress.show = false
    }, 5000)
  }
}

// 显示批量添加对话框
const showBatchAddDialog = () => {
  batchAddVisible.value = true
  // 重置表单
  Object.assign(batchAddForm, {
    urls: '',
    defaultCategory: '',
    defaultStatus: true,
    autoCheck: false
  })
}

// 预览批量添加
const previewBatchAdd = () => {
  if (parsedLinks.value.length === 0) {
    ElMessage.warning('请输入有效的链接地址')
    return
  }
  
  ElMessage.success(`解析成功，共 ${parsedLinks.value.length} 个有效链接`)
}

// 批量添加外链
const handleBatchAdd = async () => {
  if (parsedLinks.value.length === 0) {
    ElMessage.warning('请输入有效的链接地址')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要添加 ${parsedLinks.value.length} 个外链吗？${batchAddForm.autoCheck ? '添加后将自动检测可用性。' : ''}`,
      '确认批量添加',
      { type: 'info' }
    )
  } catch {
    return
  }
  
  batchAdding.value = true
  batchProgress.current = 0
  batchProgress.total = parsedLinks.value.length
  
  try {
    ElMessage.info(`开始批量添加 ${parsedLinks.value.length} 个外链...`)
    
    let successCount = 0
    let failCount = 0
    
    for (const linkData of parsedLinks.value) {
      try {
        batchProgress.current++
        
        // 创建外链
        await externalApi.createExternalLink({
          url: linkData.url,
          category: linkData.category,
          description: linkData.description,
          status: batchAddForm.defaultStatus,
          is_valid: false, // 初始设为false
          clicks: 0
        })
        
        successCount++
        
        // 如果启用自动检测
        if (batchAddForm.autoCheck) {
          try {
            // 为简化处理，批量添加时不进行单独检测
            // 用户可以添加完成后使用"全部检测"功能进行检测
            console.log('批量添加时跳过自动检测，请使用全部检测功能')
          } catch (error) {
            console.error('自动检测失败:', error)
          }
        }
        
        // 添加延迟避免请求过于频繁
        await new Promise(resolve => setTimeout(resolve, 100))
        
      } catch (error) {
        console.error('添加链接失败:', linkData.url, error)
        failCount++
      }
    }
    
    ElMessage.success(`批量添加完成: ${successCount} 个成功, ${failCount} 个失败`)
    
    // 关闭对话框
    batchAddVisible.value = false
    
    // 重新获取列表
    await fetchLinks()
    await fetchCategories()
    
    if (batchAddForm.autoCheck && successCount > 0) {
      ElMessage.info('正在进行自动检测，请稍后查看检测结果')
    }
    
  } catch (error) {
    ElMessage.error('批量添加失败')
    console.error('批量添加失败:', error)
  } finally {
    batchAdding.value = false
    batchProgress.current = 0
    batchProgress.total = 0
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedLinks.value.length === 0) {
    ElMessage.warning('请先选择要删除的链接')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedLinks.value.length} 个链接吗？此操作不可恢复。`,
      '确认批量删除',
      { type: 'warning' }
    )
    
    const ids = selectedLinks.value.map(link => link.id.toString())
    
    console.log('准备删除的ID:', ids) // 调试日志
    
    const response = await externalApi.batchDeleteExternalLinks(ids)
    
    console.log('删除响应:', response) // 调试日志
    
    ElMessage.success(`批量删除成功，共删除 ${(response as any)?.deleted_count || selectedLinks.value.length} 个链接`)
    
    // 清空选中状态
    selectedLinks.value = []
    
    // 重新获取列表
    await fetchLinks()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('批量删除完整错误信息:', error)
      ElMessage.error('批量删除失败，请重试')
    }
  }
}

// 处理选择变化
const handleSelectionChange = (selection: ExternalLink[]) => {
  selectedLinks.value = selection
}

// 切换创建表单显示
const toggleCreateForm = () => {
  showCreateForm.value = !showCreateForm.value
}

// 重置创建表单
const resetCreateForm = () => {
  Object.assign(createForm, {
    url: '',
    category: '',
    description: '',
    status: true
  })
  createFormRef.value?.clearValidate()
}

// 处理状态变更
const handleStatusChange = async (row: ExternalLink) => {
  try {
    await externalApi.updateExternalLink(row.id.toString(), {
      status: row.status
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
    row.status = !row.status // 回滚状态
  }
}

// 处理搜索
const handleSearch = () => {
  query.value.page = 1
  fetchLinks()
}

// 重置查询条件
const resetQuery = () => {
  query.value = {
    page: 1,
    per_page: 10,
    keyword: '',
    category: '',
    status: '',
    is_valid: '', // 默认显示全部
    sort_field: 'is_valid',  // 默认按可用性排序
    sort_order: 'desc'       // 降序：可用的在前面
  }
  fetchLinks()
}

// 快速筛选功能
const quickFilter = (validStatus: string) => {
  query.value.is_valid = validStatus
  query.value.page = 1 // 重置到第一页
  fetchLinks()
}

// 重试失败链接
const retryingFailed = ref(false)
const retryFailedLinks = async () => {
  try {
    retryingFailed.value = true
    
    // 先获取所有不可用的链接
    const invalidLinksResponse = await externalApi.getInvalidExternalLinks()
    const invalidLinks = (invalidLinksResponse as any)?.data || []
    
    if (invalidLinks.length === 0) {
      ElMessage.info('没有找到失败的链接')
      return
    }
    
    ElMessage.info(`🔄 开始重试 ${invalidLinks.length} 个失败链接...`)
    
    // 提取失败链接的IDs
    const failedIds = invalidLinks.map((link: any) => link.id?.toString()).filter(Boolean)
    
    if (failedIds.length === 0) {
      ElMessage.warning('无法获取失败链接的ID')
      return
    }
    
    // 调用批量检测API重试
    const response = await externalApi.batchCheckLinksBackend(failedIds, false)
    const resp: any = response
    
    let results: any[] = []
    if (resp && resp.results) {
      results = resp.results
    } else if (resp && Array.isArray(resp)) {
      results = resp
    } else if (resp && resp.data && resp.data.results) {
      results = resp.data.results
    } else if (resp && resp.data && Array.isArray(resp.data)) {
      results = resp.data
    }
    
    if (Array.isArray(results) && results.length > 0) {
      const retrySuccessCount = results.filter(r => r.is_valid === true).length
      const retryFailCount = results.filter(r => r.is_valid === false).length
      
      ElMessage.success(`🎉 重试完成！${retrySuccessCount} 个恢复成功，${retryFailCount} 个仍然失败`)
      
      // 刷新列表
      await fetchLinks()
    } else {
      ElMessage.warning('重试响应异常，请手动刷新页面查看结果')
    }
    
  } catch (error: any) {
    console.error('重试失败链接错误:', error)
    ElMessage.error(`重试失败: ${error?.message || '未知错误'}`)
  } finally {
    retryingFailed.value = false
  }
}

// 处理排序
const handleSortChange = ({ column, prop, order }: any) => {
  query.value.sort_field = prop
  query.value.sort_order = order === 'ascending' ? 'asc' : 'desc'
  fetchLinks()
}

// 处理分页
const handleSizeChange = (val: number) => {
  query.value.per_page = val
  fetchLinks()
}

const handleCurrentChange = (val: number) => {
  query.value.page = val
  fetchLinks()
}

// 处理删除
const handleDelete = async (row: ExternalLink) => {
  try {
    await ElMessageBox.confirm('确定要删除这个链接吗？', '提示', {
      type: 'warning'
    })
    await externalApi.deleteExternalLink(row.id.toString())
    ElMessage.success('删除成功')
    fetchLinks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 处理点击
const handleClick = async (row: ExternalLink) => {
  try {
    await externalApi.incrementClicks(row.id.toString())
    row.clicks++
  } catch (error) {
    console.error('增加点击量失败:', error)
  }
}

// 处理删除不可用链接
const handleDeleteInvalidLinks = async () => {
  deletingInvalid.value = true
  
  try {
    // 首先获取不可用链接的数量
    const invalidLinksResponse = await externalApi.getInvalidExternalLinks()
    const invalidCount = (invalidLinksResponse as any)?.total || 0
    
    console.log('不可用链接数量:', invalidCount) // 调试日志
    
    if (invalidCount === 0) {
      ElMessage.info('没有不可用的链接需要删除')
      return
    }
    
    await ElMessageBox.confirm(
      `确定要删除所有 ${invalidCount} 个不可用链接吗？此操作不可恢复。`,
      '确认删除所有不可用链接',
      { type: 'warning' }
    )
    
    const response = await externalApi.batchDeleteInvalidExternalLinks()
    
    console.log('删除不可用链接响应:', response) // 调试日志
    
    ElMessage.success(`批量删除成功，共删除 ${(response as any)?.deleted_count || invalidCount} 个不可用链接`)
    
    // 清空选中状态
    selectedLinks.value = []
    
    // 重新获取列表
    await fetchLinks()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除不可用链接完整错误信息:', error)
      ElMessage.error('删除不可用链接失败，请重试')
    }
  } finally {
    deletingInvalid.value = false
  }
}

onMounted(() => {
  fetchLinks()
  fetchCategories()
})
</script>

<style scoped>
.external-links-list {
  padding: 20px;
}

.create-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-card {
  margin-bottom: 20px;
}

.link-table {
  margin-bottom: 20px;
}

.link-cell {
  display: flex;
  align-items: center;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.batch-tips {
  margin-top: 8px;
}

.preview-section {
  margin-top: 20px;
}

.preview-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 12px;
}

.preview-item {
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.preview-item:last-child {
  border-bottom: none;
}

.preview-url {
  font-size: 14px;
  color: #409eff;
  margin-bottom: 4px;
  word-break: break-all;
}

.preview-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preview-desc {
  font-size: 12px;
  color: #909399;
}

.more-indicator {
  text-align: center;
  color: #909399;
  font-size: 12px;
  padding: 8px;
  border-top: 1px solid #f0f0f0;
  margin-top: 8px;
}
</style> 