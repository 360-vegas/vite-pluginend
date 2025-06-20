export default {
  name: '132',
  version: '1.0.0',
  author: '123',
  description: '123',
  category: 'tools',
  mainNav: {
    key: '132',
    title: '132',
    icon: 'Box',
    path: '/132'
  },
  subNav: [
    {
      key: 'index',
      title: '首页',
      icon: 'Document',
      path: '/132/index'
    },
  ],
  pages: [
    {
      key: 'index',
      title: '首页',
      path: '/132/index',
      name: '132-index',
      icon: 'Document',
      description: '首页页面',
      component: () => import('./pages/index.vue')
    },
  ]
}
