export default {
  name: '外链',
  version: '1.0.0',
  author: '外链',
  description: '外链',
  category: '外链',
  tags: "外链",
  mainNav: {
  "key": "plugin-wailki",
  "title": "外链",
  "icon": "StarFilled",
  "path": "/plugin-wailki/index",
  "permission": "plugin_wailki_access"
},
  subNav: [
  {
    "key": "index",
    "title": "外链列表",
    "icon": "Link",
    "path": "/plugin-wailki/index",
    "permission": "plugin_wailki_index_access"
  },
  {
    "key": "publish",
    "title": "发布外链",
    "icon": "Upload",
    "path": "/plugin-wailki/index/publish",
    "permission": "plugin_wailki_publish_access"
  },
  {
    "key": "stats",
    "title": "外链统计",
    "icon": "DataLine",
    "path": "/plugin-wailki/index/stats",
    "permission": "plugin_wailki_stats_access"
  }
],
  pages: [
    {
      key: 'index',
      title: '外链列表',
      path: '/plugin-wailki/index',
      icon: 'Link',
      permission: 'plugin_wailki_index_access',
      description: '管理所有外链',
      component: () => import('./pages/ExternalLinks.vue')
    },
    {
      key: 'publish',
      title: '发布外链',
      path: '/plugin-wailki/index/publish',
      name: 'PluginWailkiPublishExternalLink',
      icon: 'Upload',
      permission: 'plugin_wailki_publish_access',
      description: '发布外链到各个平台',
      component: () => import('./pages/ExternalLinkPublish.vue')
    },
    {
      key: 'stats',
      title: '外链统计',
      path: '/plugin-wailki/index/stats',
      name: 'PluginWailkiExternalLinkStats',
      icon: 'DataLine',
      permission: 'plugin_wailki_stats_access',
      description: '外链点击统计',
      component: () => import('./pages/ExternalLinkStats.vue')
    }
  ]
}
