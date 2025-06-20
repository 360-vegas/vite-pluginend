import { defineStore } from 'pinia'

interface MenuItem {
  key: string
  title: string
  path?: string
  name?: string
  icon?: string
  permission?: string
  children?: MenuItem[]
}

interface NavigationState {
  mainMenus: MenuItem[]
  subMenusMap: Record<string, MenuItem[]>
}

export const useNavigationStore = defineStore('navigation', {
  state: (): NavigationState => ({
    mainMenus: [
      // 系统核心菜单 - 固定不变
      {
        key: 'home',
        title: '插件管理',
        icon: 'Grid',
        path: '/plugin-generator'
      },
      {
        key: 'app-market',
        title: '应用市场',
        icon: 'Collection',
        path: '/app-market'
      }
    ],
    subMenusMap: {
      // 系统核心子菜单
      home: [
        {
          key: 'plugin-generator',
          title: '插件生成器',
          path: '/plugin-generator',
          name: 'PluginGenerator',
          icon: 'Plus'
        },
        {
          key: 'plugin-package',
          title: '插件打包',
          path: '/plugin-package',
          name: 'PluginPackage',
          icon: 'Box'
        }
      ],
      'app-market': [
        {
          key: 'installed-plugins',
          title: '已安装插件',
          path: '/app-market/installed',
          name: 'InstalledPlugins',
          icon: 'Files'
        },
        {
          key: 'market-plugins',
          title: '插件市场',
          path: '/app-market/market',
          name: 'MarketPlugins',
          icon: 'Collection'
        }
      ]
    }
  }),

  getters: {
    // 获取排序后的主菜单
    sortedMainMenus: (state) => {
      return [...state.mainMenus].sort((a, b) => {
        // 系统菜单优先级排序
        const systemOrder = ['home', 'app-market']
        const aIndex = systemOrder.indexOf(a.key)
        const bIndex = systemOrder.indexOf(b.key)
        
        if (aIndex !== -1 && bIndex !== -1) {
          return aIndex - bIndex
        }
        if (aIndex !== -1) return -1
        if (bIndex !== -1) return 1
        
        // 插件菜单按标题排序
        return a.title.localeCompare(b.title)
      })
    },

    // 获取当前路径对应的主菜单
    getCurrentMainKey: (state) => (currentPath: string) => {
      for (const [key, subMenus] of Object.entries(state.subMenusMap)) {
        if (subMenus.some(menu => menu.path === currentPath)) {
          return key
        }
      }
      return 'home' // 默认返回首页
    }
  },

  actions: {
    // 注册插件菜单 - 统一入口
    registerPluginMenu(menu: MenuItem) {
      // 检查是否已存在
      const existingIndex = this.mainMenus.findIndex(item => item.key === menu.key)
      if (existingIndex !== -1) {
        // 更新现有菜单
        this.mainMenus[existingIndex] = menu
      } else {
        // 添加新菜单
        this.mainMenus.push(menu)
      }

      // 注册子菜单
      if (menu.children && menu.children.length > 0) {
        this.subMenusMap[menu.key] = menu.children
      }
    },

    // 批量注册插件菜单
    registerPluginMenus(menus: MenuItem[]) {
      menus.forEach(menu => this.registerPluginMenu(menu))
    },

    // 移除插件菜单
    removePluginMenu(key: string) {
      // 不允许删除系统核心菜单
      const systemKeys = ['home', 'app-market']
      if (systemKeys.includes(key)) {
        console.warn(`不能删除系统核心菜单: ${key}`)
        return
      }

      const index = this.mainMenus.findIndex(item => item.key === key)
      if (index !== -1) {
        this.mainMenus.splice(index, 1)
        delete this.subMenusMap[key]
      }
    },

    // 清空所有插件菜单（保留系统菜单）
    clearPluginMenus() {
      this.mainMenus = this.mainMenus.filter(menu => 
        ['home', 'app-market'].includes(menu.key)
      )
      
      // 清理对应的子菜单映射
      Object.keys(this.subMenusMap).forEach(key => {
        if (!['home', 'app-market'].includes(key)) {
          delete this.subMenusMap[key]
        }
      })
    }
  }
})