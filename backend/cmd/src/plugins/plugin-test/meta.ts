export default {
  name: '????',
  version: '1.0.0',
  author: '????',
  description: '????????',
  category: '??',
  mainNav: {
    key: 'test',
    title: '????',
    icon: 'Box',
    path: '/test'
  },
  subNav: [
    {
      key: 'index',
      title: '??',
      icon: 'Document',
      path: '/test/index'
    },
  ],
  pages: [
    {
      key: 'index',
      title: '??',
      path: '/test/index',
      name: 'test-index',
      icon: 'Document',
      description: '??页面',
      component: () => import('./pages/index.vue')
    },
  ]
}
