import type { App } from 'vue'
import meta from './meta'

export default {
  install(app: App) {
    // 插件安装逻辑
    console.log('132 插件已安装')
  },
  meta
}

export { meta }
