/**
 * 导航服务 - 统一管理系统导航
 */

import { useNavigationStore } from '@/stores/navigation'

export interface PluginMeta {
  name: string
  version: string
  author: string
  description: string
  category: string
  tags: string
  mainNav: {
    key: string
    title: string
    icon: string
    path: string
    permission?: string
  }
  subNav: Array<{
    key: string
    title: string
    icon: string
    path: string
    permission?: string
  }>
  pages: Array<{
    key: string
    title: string
    path: string
    name?: string
    icon: string
    permission?: string
    description: string
    component: () => Promise<any>
  }>
}

/**
 * 导航服务类
 */
export class NavigationService {
  /**
   * 获取导航store实例
   */
  private getNavigationStore() {
    return useNavigationStore()
  }

  /**
   * 扫描并注册所有插件
   */
  async scanAndRegisterPlugins() {
    try {
      // 动态导入所有插件的 meta 文件
      const pluginMetaModules = import.meta.glob('../plugins/*/meta.ts', { eager: true })
      
      const registeredMenus = []
      const navigationStore = this.getNavigationStore()
      
      for (const [path, mod] of Object.entries(pluginMetaModules)) {
        const meta = (mod as any).default as PluginMeta
        
        if (!meta || !meta.mainNav) {
          console.warn(`插件 ${path} 缺少必要的 mainNav 配置`)
          continue
        }

        // 转换为导航菜单项
        const menuItem = this.transformPluginToMenuItem(meta)
        
        // 注册到导航系统
        navigationStore.registerPluginMenu(menuItem)
        registeredMenus.push(menuItem)
        
        console.log(`✅ 注册插件菜单: ${meta.mainNav.title}`)
      }

      console.log(`🎉 成功注册 ${registeredMenus.length} 个插件菜单`)
      return registeredMenus
      
    } catch (error) {
      console.error('❌ 插件扫描注册失败:', error)
      return []
    }
  }

  /**
   * 转换插件元信息为菜单项
   */
  private transformPluginToMenuItem(meta: PluginMeta) {
    return {
      key: meta.mainNav.key,
      title: meta.mainNav.title,
      icon: meta.mainNav.icon,
      path: meta.mainNav.path,
      permission: meta.mainNav.permission,
      children: meta.subNav?.map(item => {
        const correspondingPage = meta.pages?.find(page => page.key === item.key)
        return {
          key: item.key,
          title: item.title,
          icon: item.icon,
          path: item.path,
          name: correspondingPage?.name,
          permission: item.permission
        }
      }) || []
    }
  }

  /**
   * 清空所有插件菜单
   */
  clearPluginMenus() {
    const navigationStore = this.getNavigationStore()
    navigationStore.clearPluginMenus()
    console.log('🧹 已清空所有插件菜单')
  }

  /**
   * 移除指定插件的菜单
   */
  removePluginMenu(pluginKey: string) {
    const navigationStore = this.getNavigationStore()
    navigationStore.removePluginMenu(pluginKey)
    console.log(`🗑️ 移除插件菜单: ${pluginKey}`)
  }

  /**
   * 重新扫描插件
   */
  async refreshPlugins() {
    this.clearPluginMenus()
    return await this.scanAndRegisterPlugins()
  }

  /**
   * 获取当前导航状态
   */
  getNavigationState() {
    const navigationStore = this.getNavigationStore()
    return {
      mainMenus: navigationStore.sortedMainMenus,
      subMenusMap: navigationStore.subMenusMap
    }
  }
}

// 导出单例
export const navigationService = new NavigationService() 