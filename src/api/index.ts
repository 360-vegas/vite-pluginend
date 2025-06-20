import request from '@/utils/request'
import type { ExternalStatistics } from './external'

// 外链相关接口
interface ExternalLinkResponse {
  data: any[]
  meta: {
    total: number
    current_page: number
    per_page: number
    last_page: number
  }
}

interface ExternalLinkMeta {
  total: number
  current_page: number
  per_page: number
  last_page: number
}

interface ExternalLinkApiResponse {
  data: any[]
  meta: ExternalLinkMeta
}

export const externalApi = {
  // 获取外链列表
  getExternalLinks(params: {
    page?: number
    per_page?: number
    keyword?: string
    category?: string
    sort_field?: string
    sort_order?: string
  }): Promise<ExternalLinkResponse> {
    return request({
      url: '/api/external-links',
      method: 'get',
      params
    }).then(response => {
      // 处理响应数据结构
      if (response && response.data) {
        return {
          data: response.data,
          meta: (response as any).meta || {
            total: response.data.length,
            current_page: params.page || 1,
            per_page: params.per_page || 10,
            last_page: 1
          }
        }
      }
      
      // 如果响应结构不符合预期，尝试直接返回
      if (Array.isArray(response)) {
        return {
          data: response,
          meta: {
            total: response.length,
            current_page: params.page || 1,
            per_page: params.per_page || 10,
            last_page: 1
          }
        }
      }
      
      // 如果都不符合，返回空数据
      return {
        data: [],
        meta: {
          total: 0,
          current_page: 1,
          per_page: 10,
          last_page: 1
        }
      }
    }).catch(error => {
      console.error('获取外链列表失败:', error)
      if (error.response) {
        console.error('错误响应:', error.response)
      }
      throw error
    })
  },

  // 创建外链
  createExternalLink(data: any) {
    return request({
      url: '/api/external-links',
      method: 'post',
      data
    })
  },

  // 获取单个外链
  getExternalLink(id: string) {
    return request({
      url: `/api/external-links/${id}`,
      method: 'get'
    })
  },

  // 更新外链
  updateExternalLink(id: string, data: any) {
    return request({
      url: `/api/external-links/${id}`,
      method: 'put',
      data
    })
  },

  // 删除外链
  deleteExternalLink(id: string) {
    return request({
      url: `/api/external-links/${id}`,
      method: 'delete'
    })
  },

  // 批量删除外链
  batchDeleteExternalLinks(ids: string[]) {
    console.log('批量删除API调用参数:', { ids })
    return request({
      url: '/api/external-links/batch',
      method: 'delete',
      data: { ids }
    })
  },

  // 获取所有外链（不分页）
  getAllExternalLinks() {
    return request({
      url: '/api/external-links/all',
      method: 'get'
    })
  },

  // 后端批量检测外链
  batchCheckLinksBackend(ids: string[] = [], checkAll: boolean = false) {
    console.log('API调用参数:', { ids, checkAll })
    return request({
      url: '/api/external-links/batch-check',
      method: 'post',
      data: { ids, all: checkAll },
      timeout: 120000 // 批量检测设置为2分钟超时
    })
  },

  // 获取外链统计数据
  getExternalStatistics() {
    return request({
      url: '/api/external-links/statistics',
      method: 'get'
    }).then(response => {
      return response.data || response
    })
  },

  // 获取外链趋势数据
  getExternalTrends(params: { period: 'day' | 'week' | 'month', limit?: number }) {
    return request({
      url: '/api/external-links/trends',
      method: 'get',
      params
    })
  },

  // 获取点击量分布
  getClickDistribution(params: { group_size?: number }) {
    return request({
      url: '/api/external-links/click-distribution',
      method: 'get',
      params
    })
  },

  // 获取不可用外链
  getInvalidExternalLinks() {
    return request({
      url: '/api/external-links/invalid',
      method: 'get'
    })
  },

  // 批量删除不可用外链
  batchDeleteInvalidExternalLinks() {
    return request({
      url: '/api/external-links/invalid/batch',
      method: 'delete'
    })
  },

  // 增加点击量
  incrementClicks(id: string) {
    return request({
      url: `/api/external-links/${id}/clicks`,
      method: 'post'
    })
  }
}

// 应用相关接口
export const appsApi = {
  // 获取应用列表
  getAppList() {
    return request({
      url: '/api/apps',
      method: 'get'
    })
  },

  // 获取应用详情
  getAppDetail(id: string) {
    return request({
      url: `/api/apps/${id}`,
      method: 'get'
    })
  },

  // 安装应用
  installApp(id: string) {
    return request({
      url: `/api/apps/${id}/install`,
      method: 'post'
    })
  },

  // 卸载应用
  uninstallApp(id: string) {
    return request({
      url: `/api/apps/${id}/uninstall`,
      method: 'post'
    })
  },

  // 启动应用
  startApp(id: string) {
    return request({
      url: `/api/apps/${id}/start`,
      method: 'post'
    })
  },

  // 停止应用
  stopApp(id: string) {
    return request({
      url: `/api/apps/${id}/stop`,
      method: 'post'
    })
  },

  // 获取应用统计数据
  getAppStats() {
    return request({
      url: '/api/apps/stats',
      method: 'get'
    })
  },

  // 上传插件
  uploadPlugin(data: FormData) {
    return request({
      url: '/api/plugins/upload',
      method: 'post',
      data,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }).then(response => {

      
      if (!response) {
        throw new Error('API 响应为空')
      }
      
      // 从响应中提取插件信息
      const pluginData = response.data || response
      return {
        data: {
          id: pluginData.id || 'hello-world', // 使用插件ID作为默认值
          name: pluginData.name || 'Hello World', // 使用插件名称作为默认值
          version: pluginData.version || '1.0.0',
          description: pluginData.description || 'A simple hello world plugin for testing',
          author: pluginData.author || 'Test Author',
          icon: pluginData.icon || 'Box',
          status: 'inactive',
          features: pluginData.features || [
            'Display hello world message',
            'Show current time',
            'Basic counter'
          ],
          dependencies: pluginData.dependencies || [
            { name: 'vue', version: '^3.0.0' },
            { name: 'element-plus', version: '^2.0.0' }
          ],
          route: pluginData.route || '/hello-world',
          component: pluginData.component || './src/plugins/hello-world-plugin/src/index.ts',
          order: pluginData.order || 0
        }
      }
    }).catch(error => {
      console.error('上传插件失败:', error)
      if (error.response) {
        console.error('错误响应:', error.response)
        throw new Error((error.response as any).data?.message || '上传插件失败')
      }
      throw error
    })
  },

  // 获取插件列表
  listPlugins: () => {
    return request({
      url: '/api/plugins',
      method: 'get'
    })
  },

  // 获取插件环境要求
  getPluginRequirements: (key: string) => {
    return request({
      url: `/api/plugins/${key}/requirements`,
      method: 'get'
    })
  },

  // 检查环境
  checkEnvironment: (key: string) => {
    return request({
      url: `/api/plugins/${key}/check-environment`,
      method: 'get'
    })
  },

  // 修复环境问题
  fixEnvironment: (key: string) => {
    return request({
      url: `/api/plugins/${key}/fix-environment`,
      method: 'post'
    })
  },

  // 构建插件
  buildPlugin: (data: any) => {
    return request({
      url: '/api/plugins/build',
      method: 'post',
      data
    })
  },

  // 保存插件配置
  savePluginConfig: (data: any) => {
    return request({
      url: '/api/plugins/config',
      method: 'post',
      data
    })
  },

  // 打包插件
  packagePlugin: (data: any) => {
    return request({
      url: '/api/plugins/package',
      method: 'post',
      data
    })
  },

  installPlugin: (pluginId: string) => {
    return request({
      url: `/api/plugins/${pluginId}/install`,
      method: 'post'
    })
  },
  
  uninstallPlugin: (pluginId: string) => {
    return request({
      url: `/api/plugins/${pluginId}/uninstall`,
      method: 'post'
    })
  },

  // 生成插件
  generatePlugin: (data: any) => {
    return request({
      url: '/api/create-plugin',
      method: 'post',
      data
    })
  },

  // 扫描插件
  scanPlugins: () => {
    return request({
      url: '/api/plugins/scan',
      method: 'post'
    })
  },

  // 获取插件详情
  getPluginDetail: (pluginKey: string) => {
    return request({
      url: `/api/plugins/${pluginKey}`,
      method: 'get'
    })
  },

  // 启用/禁用插件
  togglePlugin: (pluginKey: string, enabled: boolean) => {
    return request({
      url: `/api/plugins/${pluginKey}/toggle`,
      method: 'post',
      data: { enabled }
    })
  },

  // 删除插件
  deletePlugin: (pluginKey: string) => {
    return request({
      url: `/api/plugins/${pluginKey}`,
      method: 'delete'
    })
  },

  // 导出插件
  exportPlugin: (pluginKey: string) => {
    return request({
      url: `/api/plugins/${pluginKey}/export`,
      method: 'get',
      responseType: 'blob'
    })
  },

  // 下载插件包
  downloadPlugin: (pluginKey: string) => {
    return request({
      url: `/api/plugins/${pluginKey}/download`,
      method: 'get',
      responseType: 'blob'
    })
  },

  // 安装插件包
  installPluginPackage: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    
    return request({
      url: '/api/plugins/install',
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 检查插件依赖
  checkPluginDependencies: (pluginKey: string) => {
    return request({
      url: `/api/plugins/${pluginKey}/dependencies/check`,
      method: 'get'
    })
  },

  // 设置插件数据库
  setupPluginDatabase: (pluginKey: string, config: any) => {
    return request({
      url: `/api/plugins/${pluginKey}/dependencies/setup`,
      method: 'post',
      data: config
    })
  }
}

export default {
  externalApi,
  appsApi
} 