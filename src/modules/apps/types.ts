export interface Plugin {
  id: string
  name: string
  description: string
  version: string
  author: string
  icon?: string
  route?: {
    path: string
    name: string
    component: any
    meta?: {
      title: string
      icon: string
      menu?: boolean
    }
    children?: Array<{
      path: string
      name: string
      component: any
      meta?: {
        title: string
        icon: string
        menu?: boolean
      }
    }>
  }
  status: 'installed' | 'uninstalled' | 'inactive'
} 