import { App } from 'vue'

declare global {
  interface Window {
    __VUE_APP__: App
  }
} 