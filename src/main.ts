import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { versionChecker, initVersionInfo } from './services/version-control'

// 导入全局样式
import './styles/index.scss'

// 导入导航服务
import { navigationService } from '@/services/navigation'

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 设置全局变量
window.__VUE_APP__ = app

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 初始化版本信息
initVersionInfo(import.meta.env.VITE_APP_VERSION || '1.0.0', import.meta.env.VITE_APP_BUILD_TIME || '')

// 挂载应用并注册插件
router.isReady().then(async () => {
  // 先挂载应用
  app.mount('#app')
  
  console.log('✅ 应用挂载完成')
  
  // 应用挂载后再注册插件菜单
  try {
    await navigationService.scanAndRegisterPlugins()
    console.log('✅ 插件菜单注册完成')
  } catch (error) {
    console.error('❌ 插件注册失败:', error)
  }
})

// 启动版本检查
versionChecker.start()

// 在开发环境下禁用版本检查
if (import.meta.env.DEV) {
  versionChecker.stop()
}