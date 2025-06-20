import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 版本信息接口
interface VersionInfo {
  version: string
  buildTime: string
  changelog: string
  files: {
    js: string[]
    css: string[]
  }
}

// 当前版本信息
const currentVersion = ref<string>('1.0.0')
const currentBuildTime = ref<string>('')

// 检查更新间隔（毫秒）
const CHECK_INTERVAL = 5 * 60 * 1000 // 5分钟

// 版本检查器
class VersionChecker {
  private timer: number | null = null
  private isChecking = false
  private updateAvailable = false
  private swRegistration: ServiceWorkerRegistration | null = null

  // 初始化 Service Worker
  async initServiceWorker() {
    if ('serviceWorker' in navigator) {
      try {
        // 暂时禁用 Service Worker
      // this.swRegistration = await navigator.serviceWorker.register('/sw.js')
        console.log('Service Worker 注册成功')
      } catch (error) {
        console.error('Service Worker 注册失败:', error)
      }
    }
  }

  // 开始检查更新
  start() {
    if (this.timer) return
    // this.initServiceWorker() // 强制禁用 Service Worker 注册
    this.checkVersion()
    this.timer = window.setInterval(() => this.checkVersion(), CHECK_INTERVAL)
  }

  // 停止检查更新
  stop() {
    if (this.timer) {
      clearInterval(this.timer)
      this.timer = null
    }
  }

  // 检查版本
  private async checkVersion() {
    if (this.isChecking) return
    this.isChecking = true

    try {
      const response = await fetch('/version.json?t=' + Date.now())
      const data: VersionInfo = await response.json()

      if (this.isNewVersion(data.version)) {
        this.handleNewVersion(data)
      }
    } catch (error) {
      console.error('版本检查失败:', error)
    } finally {
      this.isChecking = false
    }
  }

  // 判断是否是新版本
  private isNewVersion(version: string): boolean {
    return version !== currentVersion.value
  }

  // 处理新版本
  private async handleNewVersion(info: VersionInfo) {
    if (this.updateAvailable) return
    this.updateAvailable = true

    try {
      await this.showUpdateDialog(info)
    } catch (error) {
      console.error('处理新版本失败:', error)
    }
  }

  // 显示更新对话框
  private async showUpdateDialog(info: VersionInfo) {
    try {
      await ElMessageBox.confirm(
        `发现新版本 ${info.version}，是否立即更新？\n\n更新内容：\n${info.changelog}`,
        '系统更新',
        {
          confirmButtonText: '立即更新',
          cancelButtonText: '稍后更新',
          type: 'info'
        }
      )
      await this.performUpdate(info)
    } catch {
      // 用户取消更新
    }
  }

  // 执行更新
  private async performUpdate(info: VersionInfo) {
    try {
      // 显示更新进度
      const loading = ElMessage({
        message: '正在更新系统...',
        type: 'info',
        duration: 0
      })

      // 预加载新版本资源
      await this.preloadNewVersion(info)

      // 更新成功提示
      loading.close()
      ElMessage.success('系统更新成功，新功能已生效')

      // 更新版本信息
      currentVersion.value = info.version
      currentBuildTime.value = info.buildTime

      // 重置状态
      this.updateAvailable = false
    } catch (error) {
      ElMessage.error('更新失败，请稍后重试')
      console.error('更新失败:', error)
    }
  }

  // 预加载新版本
  private async preloadNewVersion(info: VersionInfo) {
    // 预加载 JS 文件
    await Promise.all(
      info.files.js.map((file) => this.preloadFile(file))
    )

    // 预加载 CSS 文件
    await Promise.all(
      info.files.css.map((file) => this.preloadFile(file))
    )
  }

  // 预加载单个文件
  private async preloadFile(url: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const link = document.createElement('link')
      link.rel = 'preload'
      link.as = url.endsWith('.js') ? 'script' : 'style'
      link.href = url + '?t=' + Date.now()
      link.onload = () => resolve()
      link.onerror = () => reject()
      document.head.appendChild(link)
    })
  }
}

// 创建版本检查器实例
export const versionChecker = new VersionChecker()

// 初始化版本信息
export const initVersionInfo = (version: string, buildTime: string) => {
  currentVersion.value = version
  currentBuildTime.value = buildTime
}

// 导出当前版本信息
export const getCurrentVersion = () => currentVersion.value
export const getCurrentBuildTime = () => currentBuildTime.value 