<template>
  <div class="sidebar-container">
    <!-- 主导航 -->
    <div class="main-sidebar">
      <div class="logo">
        <img src="/vite.svg" alt="logo" />
      </div>
      
      <div class="main-menu">
        <!-- 动态菜单项 + 硬编码调试菜单 -->
        <div 
          v-for="item in allMainMenus" 
          :key="item.key"
          class="main-menu-item"
          :class="{ active: currentMain === item.key }"
          @click="handleMainClick(item.key)"
        >
          <el-tooltip 
            :content="item.title" 
            placement="right"
            :show-after="500"
          >
            <div class="menu-item-content">
              <el-icon><component :is="item.icon" /></el-icon>
              <span class="menu-title">{{ item.title }}</span>
            </div>
          </el-tooltip>
        </div>
      </div>
    </div>

    <!-- 次导航 -->
    <div class="sub-sidebar hide-scrollbar" v-show="!collapsed">
      <div class="sub-header">
        <div class="breadcrumb">
          <el-icon><HomeFilled /></el-icon>
          <span>{{ currentMainTitle }}</span>
          <span v-if="currentSubTitle" class="separator">/</span>
          <span class="sub-title">{{ currentSubTitle }}</span>
        </div>
      </div>
      
      <el-menu
        :default-active="$route.path"
        background-color="#fff"
        text-color="#333"
        active-text-color="#fff"
      >
        <template v-for="item in currentSubMenus" :key="item.key">
          <el-sub-menu v-if="item.children && item.children.length > 0" :index="item.key">
            <template #title>
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item 
              v-for="child in item.children" 
              :key="child.key" 
              :index="child.key"
              @click="handleMenuItemClick(child)"
            >
              <el-icon><component :is="child.icon" /></el-icon>
              <template #title>{{ child.title }}</template>
            </el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else :index="item.key" @click="handleMenuItemClick(item)">
            <el-icon><component :is="item.icon" /></el-icon>
            <template #title>{{ item.title }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNavigationStore } from '@/stores/navigation'
import {
  HomeFilled,
  Grid,
  Link,
  Connection,
  Share,
  Setting,
  Plus,
  DataLine,
  TrendCharts,
  User,
  ChatDotSquare,
  Files,
  Collection,
  Document,
  Crop,
  UserFilled,
  Lock,
  Monitor,
  ChatDotRound,
  Cpu,
  Box,
  StarFilled,
  Upload
} from '@element-plus/icons-vue'

defineProps<{
  collapsed: boolean
}>()

const route = useRoute()
const router = useRouter()
const currentMain = ref('home')

// 获取导航store（使用try-catch确保安全）
let navigationStore: any = null
try {
  navigationStore = useNavigationStore()
  console.log('✅ NavigationStore 初始化成功')
} catch (error) {
  console.error('❌ NavigationStore 初始化失败:', error)
}

// 硬编码的调试菜单项
const debugMenuItem = {
  key: 'debug',
  title: '调试',
  icon: 'Monitor',
  path: '/debug'
}

// 合并动态菜单和调试菜单
const allMainMenus = computed(() => {
  try {
    const dynamicMenus = navigationStore?.sortedMainMenus || []
    return [...dynamicMenus, debugMenuItem]
  } catch (error) {
    console.error('❌ 获取菜单失败:', error)
    return [debugMenuItem] // 失败时至少显示调试菜单
  }
})

// 当前主菜单标题
const currentMainTitle = computed(() => {
  try {
    const menu = allMainMenus.value.find((item: any) => item.key === currentMain.value)
    return menu ? menu.title : ''
  } catch (error) {
    console.error('❌ 获取主菜单标题失败:', error)
    return ''
  }
})

// 当前子菜单
const currentSubMenus = computed(() => {
  try {
    // 调试菜单的子菜单
    if (currentMain.value === 'debug') {
      return [{
        key: 'debug-page',
        title: '调试页面',
        icon: 'Monitor',
        path: '/debug'
      }]
    }
    
    // 动态菜单的子菜单
    return navigationStore?.subMenusMap[currentMain.value] || []
  } catch (error) {
    console.error('❌ 获取子菜单失败:', error)
    return []
  }
})

// 当前子菜单标题
const currentSubTitle = computed(() => {
  try {
    const currentPath = route.path
    const subMenus = currentSubMenus.value
    const currentSubMenu = subMenus.find((item: any) => item.path === currentPath)
    return currentSubMenu?.title || ''
  } catch (error) {
    console.error('❌ 获取子菜单标题失败:', error)
    return ''
  }
})

// 处理主菜单点击
const handleMainClick = (key: string) => {
  try {
    currentMain.value = key
    
    // 特殊处理调试菜单
    if (key === 'debug') {
      router.push('/debug')
      return
    }
    
    // 处理动态菜单
    const subMenus = navigationStore?.subMenusMap[key]
    if (subMenus && subMenus.length > 0) {
      handleMenuItemClick(subMenus[0])
    } else {
      // 如果没有子菜单，直接跳转到主路由
      const mainMenu = allMainMenus.value.find((menu: any) => menu.key === key)
      if (mainMenu && mainMenu.path) {
        handleMenuItemClick(mainMenu)
      }
    }
  } catch (error) {
    console.error('❌ 处理主菜单点击失败:', error)
  }
}

// 处理菜单项点击
const handleMenuItemClick = (item: any) => {
  try {
    console.log('点击菜单项:', item)
    if (item.name) {
      router.push({ name: item.name })
    } else if (item.path) {
      router.push(item.path)
    }
  } catch (error) {
    console.error('❌ 处理菜单项点击失败:', error)
  }
}

// 初始化当前主菜单
const initCurrentMain = () => {
  try {
    const path = route.path
    
    // 优先检查调试路径
    if (path.includes('/debug')) {
      currentMain.value = 'debug'
      return
    }
    
    // 检查动态菜单路径
    if (navigationStore?.getCurrentMainKey) {
      const mainKey = navigationStore.getCurrentMainKey(path)
      currentMain.value = mainKey
    } else {
      currentMain.value = 'home'
    }
  } catch (error) {
    console.error('❌ 初始化主菜单失败:', error)
    currentMain.value = 'home'
  }
}

// 监听路由变化
watch(
  () => route.path,
  () => {
    initCurrentMain()
  }
)

// 初始化
onMounted(() => {
  console.log('🔍 完整版Sidebar已挂载, 当前路由:', route.path)
  initCurrentMain()
})
</script>

<style lang="scss" scoped>
.sidebar-container {
  display: flex;
  height: 100%;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;

  .main-sidebar {
    width: 64px;
    height: 100%;
    background-color: #001529;
    display: flex;
    flex-direction: column;
    position: relative;
    z-index: 10;
    padding: 4px;
    
    .logo {
      width: 56px;
      height: 56px;
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 0 auto;
      border-radius: 4px;
      
      img {
        width: 24px;
        height: 24px;
      }
    }
    
    .main-menu {
      flex: 1;
      overflow-y: auto;
      overflow-x: hidden;
      
      &::-webkit-scrollbar {
        width: 0;
        background: transparent;
      }
      
      .main-menu-item {
        height: 56px;
        width: 56px;
        display: flex;
        justify-content: center;
        align-items: center;
        cursor: pointer;
        color: rgba(255, 255, 255, 0.65);
        transition: all 0.3s;
        border-radius: 4px;
        margin: 4px auto;
        
        &:hover {
          color: #fff;
          background-color: rgba(24, 144, 255, 0.1);
        }
        
        &.active {
          color: #fff;
          background-color: #409eff;
        }

        .menu-item-content {
          display: flex;
          flex-direction: column;
          align-items: center;
          
          .el-icon {
            font-size: 16px;
            margin-bottom: 4px;
          }
          
          .menu-title {
            font-size: 12px;
            line-height: 1;
          }
        }
      }
    }
  }

  .sub-sidebar {
    width: 220px;
    height: 100%;
    background-color: #fff;
    border-right: 1px solid #e6e6e6;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .sub-header {
      height: 56px;
      display: flex;
      align-items: center;
      padding: 0 16px;
      border-bottom: 1px solid #e6e6e6;
      background-color: #f8f8f8;
      flex-shrink: 0;

      .breadcrumb {
        display: flex;
        align-items: center;
        font-size: 14px;
        color: #666;

        .el-icon {
          margin-right: 8px;
          font-size: 16px;
          color: #409eff;
        }

        .separator {
          margin: 0 8px;
          color: #999;
        }

        .sub-title {
          color: #333;
          font-weight: 500;
        }
      }
    }

    :deep(.el-menu) {
      flex: 1;
      border-right: none;
      overflow-y: auto;
      overflow-x: hidden;
      padding: 4px;
      margin-top: 5px;
      background-color: transparent;

      .el-menu-item {
        height: 56px;
        line-height: 56px;
        margin-bottom: 4px;
        border-radius: 4px;
        padding: 0 16px !important;

        &:hover {
          background-color: #f5f7fa !important;
          color: #409eff !important;
        }
        
        &.is-active {
          background-color: #ecf5ff !important;
          color: #409eff !important;
          border-right: none;

          .el-icon {
            color: #409eff !important;
          }
        }

        .el-icon {
          font-size: 16px;
          margin-right: 12px;
        }
      }

      &::-webkit-scrollbar {
        width: 0;
        background: transparent;
      }
    }
  }
}
</style> 