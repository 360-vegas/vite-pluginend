const CACHE_NAME = 'app-cache-v1'
const VERSION_CHECK_INTERVAL = 5 * 60 * 1000 // 5分钟检查一次

// 安装 Service Worker
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll([
        '/',
        '/index.html',
        '/static/js/app.js',
        '/static/css/app.css'
      ])
    })
  )
})

// 激活 Service Worker
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME)
          .map((name) => caches.delete(name))
      )
    })
  )
})

// 拦截网络请求
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      // 返回缓存的响应或发起网络请求
      return response || fetch(event.request)
    })
  )
})

// 定期检查版本更新
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'CHECK_VERSION') {
    checkVersionUpdate()
  }
})

// 检查版本更新
async function checkVersionUpdate() {
  try {
    const response = await fetch('/version.json')
    const data = await response.json()
    
    // 比较版本
    if (data.version !== self.version) {
      // 通知主线程有新版本
      self.clients.matchAll().then((clients) => {
        clients.forEach((client) => {
          client.postMessage({
            type: 'UPDATE_AVAILABLE',
            version: data.version,
            changelog: data.changelog
          })
        })
      })
    }
  } catch (error) {
    console.error('版本检查失败:', error)
  }
} 