import request from '@/utils/request'

// 外链接口类型定义
export interface ExternalLink {
  id: number
  url: string
  category: string
  description: string
  status: boolean
  is_valid: boolean
  clicks: number
  created_at: string
  updated_at: string
  last_clicked_at: string
}

// 外链统计接口类型定义
export interface ExternalStatistics {
  total_links: number
  active_links: number
  expired_links: number
  invalid_links: number
  total_clicks: number
  average_clicks: number
  categories: {
    [key: string]: number
  }
  tags: {
    [key: string]: number
  }
}

// 外链趋势数据接口类型定义
export interface ExternalTrend {
  date: string
  new_links: number
  total_clicks: number
  active_links: number
}

// 点击量分布接口类型定义
export interface ClickDistribution {
  range: string
  count: number
}

// 外链查询参数接口
export interface ExternalLinkQuery {
  page: number
  per_page: number
  keyword?: string
  category?: string
  status?: string
  is_valid?: string
  sort_field?: string
  sort_order?: 'asc' | 'desc'
}

// 外链列表响应接口
export interface ExternalLinkResponse {
  data: ExternalLink[]
  meta: {
    total: number
    page: number
    per_page: number
  }
}

// 外链创建接口
export interface ExternalLinkCreate {
  url: string
  category: string
  description?: string
  status?: boolean
}

// 外链更新接口
export interface ExternalLinkUpdate extends Partial<ExternalLinkCreate> {}

// 外链统计接口
export interface ExternalLinkStats {
  total_links: number
  total_clicks: number
  active_links: number
  categories: Record<string, number>
  trends: Array<{
    date: string
    clicks: number
  }>
  top_links: ExternalLink[]
}

// 链接检查结果接口
export interface LinkCheckResult {
  url: string
  status: 'success' | 'error'
  message: string
}

class ExternalApi {
  // 获取外链列表
  getExternalLinks(params: ExternalLinkQuery) {
    return request.get<{
      data: ExternalLink[]
      meta: {
        total: number
        page: number
        per_page: number
      }
    }>('/api/external-links', { params })
  }

  // 获取单个外链
  getExternalLink(id: string) {
    return request.get<ExternalLink>(`/api/external-links/${id}`)
  }

  // 创建外链
  createExternalLink(data: ExternalLinkCreate) {
    return request.post<ExternalLink>('/api/external-links', data)
  }

  // 批量创建外链
  batchCreateExternalLinks(data: ExternalLinkCreate[]) {
    return request.post<ExternalLink[]>('/api/external-links/batch', data)
  }

  // 更新外链
  updateExternalLink(id: string, data: ExternalLinkUpdate) {
    return request.put<ExternalLink>(`/api/external-links/${id}`, data)
  }

  // 删除外链
  deleteExternalLink(id: string) {
    return request.delete(`/api/external-links/${id}`)
  }

  // 批量删除外链
  batchDeleteExternalLinks(ids: string[]) {
    return request.delete('/api/external-links/batch', { data: { ids } })
  }

  // 增加点击量
  incrementClicks(id: string) {
    return request.post(`/api/external-links/${id}/clicks`)
  }

  // 获取统计数据
  getExternalStatistics() {
    return request.get<ExternalLinkStats>('/api/external-links/statistics')
  }

  // 检测链接可用性
  checkExternalLinks() {
    return request.get<LinkCheckResult[]>('/api/external-links/check')
  }

  // 检测单个链接
  checkSingleLink(id: string) {
    return request.post<LinkCheckResult>(`/api/external-links/${id}/check`)
  }

  // 批量检测链接
  batchCheckLinks(ids: string[]) {
    return request.post<LinkCheckResult[]>('/api/external-links/batch-check', { ids })
  }

  // 获取所有外链（不分页）
  getAllExternalLinks() {
    return request.get<{
      data: ExternalLink[]
      total: number
    }>('/api/external-links/all')
  }

  // 后端批量检测链接（支持检测所有链接）
  batchCheckLinksBackend(ids: string[] = [], checkAll: boolean = false) {
    return request.post<{
      message: string
      results: LinkCheckResult[]
    }>('/api/external-links/batch-check', { 
      ids, 
      all: checkAll 
    })
  }

  // 访问外链（后台访问）
  visitExternalLink(id: string) {
    return request.post<{ content: string }>(`/api/external-links/${id}/visit`)
  }

  // 获取不可用外链
  getInvalidExternalLinks() {
    return request.get<{
      data: ExternalLink[]
      total: number
    }>('/api/external-links/invalid')
  }

  // 批量删除不可用外链
  batchDeleteInvalidExternalLinks() {
    return request.delete<{
      message: string
      deleted_count: number
    }>('/api/external-links/invalid/batch')
  }
}

export const externalApi = new ExternalApi() 