import { definePlugin } from '@/plugins/types'
import ExternalLinks from './pages/ExternalLinks.vue'
import ExternalLinkForm from './pages/ExternalLinkForm.vue'
import ExternalLinkStats from './pages/ExternalLinkStats.vue'
import ExternalLinkPublish from './pages/ExternalLinkPublish.vue'

export default definePlugin({
  id: 'plugin-wailki',
  name: '外链管理',
  description: '管理外部链接和跳转规则',
  version: '1.0.0',
  author: 'System',
  icon: 'Link',
  status: 'installed',
  route: {
    path: '',
    name: 'PluginWailkiExternalLinks',
    component: ExternalLinks,
    meta: {
      title: '外链管理',
      icon: 'Link',
      menu: true
    },
    children: [
      {
        path: 'publish',
        name: 'PluginWailkiPublishExternalLink',
        component: ExternalLinkPublish,
        meta: {
          title: '发布外链',
          icon: 'Upload',
          menu: true
        }
      },
      {
        path: 'create',
        name: 'PluginWailkiCreateExternalLink',
        component: ExternalLinkForm,
        meta: {
          title: '创建外链',
          icon: 'Plus',
          menu: false
        }
      },
      {
        path: 'edit/:id',
        name: 'PluginWailkiEditExternalLink',
        component: ExternalLinkForm,
        meta: {
          title: '编辑外链',
          icon: 'Edit',
          menu: false
        }
      },
      {
        path: 'stats',
        name: 'PluginWailkiExternalLinkStats',
        component: ExternalLinkStats,
        meta: {
          title: '外链统计',
          icon: 'DataLine',
          menu: true
        }
      }
    ]
  }
}) 