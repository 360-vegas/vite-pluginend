export interface Plugin {
  id: string
  name: string
  description: string
  version: string
  author: string
  status: 'active' | 'inactive'
  route?: {
    path: string
    name: string
    component: any
    meta?: {
      title: string
      icon?: string
      order?: number
    }
  }
} 