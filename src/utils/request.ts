import axios from 'axios'
import type { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const service: AxiosInstance = axios.create({
  // 在开发环境下使用相对路径（通过Vite代理），生产环境使用完整URL
  baseURL: import.meta.env.DEV ? '' : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'),
  timeout: Number(import.meta.env.VITE_API_TIMEOUT) || 15000, // 增加请求超时时间到15秒
  headers: {
    'Content-Type': 'application/json'
  }
})

console.log('开发环境:', import.meta.env.DEV);
console.log('环境变量 VITE_API_BASE_URL:', import.meta.env.VITE_API_BASE_URL);
console.log('实际使用的 Axios baseURL:', service.defaults.baseURL);

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 在这里可以添加token等认证信息
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    console.log('原始响应:', response)
    console.log('响应数据:', response.data)
    
    // 检查响应状态码
    if (response.status < 200 || response.status >= 300) {
      console.error('响应状态码错误:', response.status)
      return Promise.reject(new Error(`HTTP错误: ${response.status}`))
    }
    
    // 检查响应数据
    if (!response.data) {
      console.error('响应数据为空')
      return Promise.reject(new Error('响应数据为空'))
    }
    
    // 检查响应格式
    if (typeof response.data !== 'object') {
      console.error('响应格式错误:', typeof response.data)
      return Promise.reject(new Error('响应格式错误'))
    }
    
    console.log('处理后返回的数据:', response.data)
    return response.data
  },
  error => {
    console.error('响应错误:', error)
    
    if (error.response) {
      // 服务器返回了错误状态码
      console.error('错误响应:', error.response)
      const status = error.response.status
      const data = error.response.data
      
      switch (status) {
        case 401:
          console.error('未授权错误')
          return Promise.reject(new Error('未授权，请重新登录'))
        case 403:
          console.error('权限错误')
          return Promise.reject(new Error('没有权限访问此资源'))
        case 404:
          console.error('资源不存在')
          return Promise.reject(new Error('请求的资源不存在'))
        case 500:
          console.error('服务器内部错误:', data)
          return Promise.reject(new Error(data?.message || '服务器内部错误'))
        default:
          console.error('其他错误:', status)
          return Promise.reject(new Error(data?.message || `请求失败: ${status}`))
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      console.error('请求错误:', error.request)
      return Promise.reject(new Error('网络请求失败，请检查网络连接'))
    } else {
      // 请求配置出错
      console.error('配置错误:', error.message)
      return Promise.reject(new Error('请求配置错误'))
    }
  }
)

export default service 