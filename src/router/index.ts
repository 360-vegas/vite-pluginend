import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

// 1. 扫描所有插件 meta
const pluginMetaModules = import.meta.glob('../plugins/*/meta.ts', { eager: true })

// 2. 生成插件页面的路由数组（累加所有插件页面）
const pluginRoutes: RouteRecordRaw[] = []
Object.entries(pluginMetaModules).forEach(([path, mod]: any) => {
  const meta = mod.default
  if (Array.isArray(meta.pages)) {
    meta.pages.forEach((page: any) => {
      const pluginDirMatch = path.match(/plugins[\\/](.*)[\\/]/)
      if (!pluginDirMatch) return
      const pluginDir = pluginDirMatch[1]
      const fileName = page.key.charAt(0).toUpperCase() + page.key.slice(1)
      pluginRoutes.push({
        path: page.path,
        name: page.name, // 用唯一 name
        component: page.component, // 直接使用 meta.ts 中定义的组件引用
        meta: {
          title: page.title,
          icon: page.icon,
          permission: page.permission,
          description: page.description
        }
      })
    })
  }
})

// 3. 合并插件页面到 / 路由的 children
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Root',
    component: () => import('@/layouts/default/index.vue'),
    redirect: '/plugin-generator', // 恢复默认路由
    children: [
      {
        path: 'debug',
        name: 'DebugPage',
        component: () => import('../pages/DebugPage.vue'),
        meta: {
          title: '调试页面',
          icon: 'Monitor'
        }
      },
      {
        path: 'plugin-generator',
        name: 'PluginGenerator',
        component: () => import('../pages/PluginGenerator.vue'),
        meta: {
          title: '插件生成器',
          icon: 'Plus'
        }
      },
      {
        path: 'plugin-package',
        name: 'PluginPackage',
        component: () => import('../pages/PluginPackage.vue'),
        meta: {
          title: '插件打包',
          icon: 'Box'
        }
      },
      {
        path: 'app-market',
        name: 'AppMarket',
        component: () => import('@/pages/AppMarket.vue'),
        redirect: '/app-market/installed',
        meta: { 
          title: '应用市场',
          icon: 'Collection'
        },
        children: [
          {
            path: 'installed',
            name: 'InstalledPlugins',
            component: () => import('@/pages/AppMarket/InstalledPlugins.vue'),
            meta: { 
              title: '已安装插件',
              icon: 'Files'
            }
          },
          {
            path: 'market',
            name: 'MarketPlugins',
            component: () => import('@/pages/AppMarket/MarketPlugins.vue'),
            meta: { 
              title: '插件市场',
              icon: 'Collection'
            }
          }
        ]
      },
      ...pluginRoutes // 插件页面全部静态注册为 children
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/modules/components/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局路由守卫
router.beforeEach((to, from, next) => {
  console.log('🚦 路由跳转:', from.path, '->', to.path)
  
  // 如果路由不存在，重定向到首页
  if (to.matched.length === 0) {
    console.log('❌ 路由不存在，重定向到首页')
    next('/')
    return
  }
  next()
})

export { router }
export default router