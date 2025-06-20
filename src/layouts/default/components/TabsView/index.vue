<template>
  <div class="tabs-view">
    <div class="tabs-nav">
      <el-tabs
        v-model="activeTab"
        type="card"
        @tab-remove="removeTab"
        @tab-click="clickTab"
      >
        <el-tab-pane
          v-for="item in visitedViews"
          :key="item.path"
          :label="item.title"
          :name="item.path"
          :closable="item.path !== '/dashboard'"
        >
          <template #label>
            <span class="tab-label">
              <el-icon v-if="item.path === '/dashboard'"><HomeFilled /></el-icon>
              {{ item.title }}
            </span>
          </template>
        </el-tab-pane>
      </el-tabs>
    </div>

    <div class="tabs-actions">
      <!-- 应用切换按钮 -->
      <el-dropdown 
        class="app-switch-dropdown" 
        trigger="click"
        @command="handleCommand"
      >
        <el-button class="app-switch">
          <el-icon><Grid /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="refresh">
              <el-icon><Refresh /></el-icon>
              <span>刷新</span>
            </el-dropdown-item>
            <el-dropdown-item command="closeOthers">
              <el-icon><Close /></el-icon>
              <span>关闭其他</span>
            </el-dropdown-item>
            <el-dropdown-item command="closeLeft">
              <el-icon><Back /></el-icon>
              <span>关闭左侧</span>
            </el-dropdown-item>
            <el-dropdown-item command="closeRight">
              <el-icon><Right /></el-icon>
              <span>关闭右侧</span>
            </el-dropdown-item>
            <el-dropdown-item command="closeAll">
              <el-icon><Close /></el-icon>
              <span>关闭全部</span>
            </el-dropdown-item>
            <el-dropdown-item divided command="settings">
              <el-icon><Setting /></el-icon>
              <span>标签设置</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    
    <!-- 右键菜单 -->
    <ul v-show="contextMenuVisible" class="contextmenu" :style="contextMenuStyle">
      <li @click="refreshSelectedTag">刷新页面</li>
      <li @click="closeSelectedTag">关闭当前</li>
      <li @click="closeOthersTags">关闭其他</li>
      <li @click="closeAllTags">关闭所有</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  HomeFilled,
  Grid,
  Refresh,
  Close,
  Back,
  Right,
  Setting
} from '@element-plus/icons-vue'

interface TabView {
  title: string
  path: string
  icon?: string
}

const route = useRoute()
const router = useRouter()

const visitedViews = ref<TabView[]>([])
const activeTab = ref('/dashboard')
const contextMenuVisible = ref(false)
const contextMenuStyle = ref({
  left: '0px',
  top: '0px'
})

// 在组件挂载时初始化首页标签
onMounted(() => {
  visitedViews.value = [{
    title: '首页1',
    path: '/dashboard',
    icon: 'HomeFilled'
  }]
})

// 监听路由变化
watch(
  () => route.fullPath,
  (newPath) => {
    // 如果是首页相关路径，只激活不添加
    if (newPath === '/' || newPath === '/dashboard') {
      activeTab.value = '/dashboard'
      return
    }
    
    // 添加新标签
    const title = route.meta?.title as string
    if (title && !visitedViews.value.some(v => v.path === newPath)) {
      visitedViews.value.push({
        title,
        path: newPath,
        icon: route.meta?.icon as string
      })
    }
    activeTab.value = newPath
  }
)

// 点击标签
const clickTab = (tab: any) => {
  const path = tab.props.name
  if (path === activeTab.value) return
  router.push(path)
}

// 移除标签
const removeTab = (targetPath: string) => {
  const tabs = visitedViews.value
  let activePath = activeTab.value
  
  if (activePath === targetPath) {
    tabs.forEach((tab, index) => {
      if (tab.path === targetPath) {
        const nextTab = tabs[index + 1] || tabs[index - 1]
        if (nextTab) {
          activePath = nextTab.path
        }
      }
    })
  }
  
  activeTab.value = activePath
  visitedViews.value = tabs.filter(tab => tab.path !== targetPath)
  router.push(activePath)
}

// 右键菜单相关方法
const refreshSelectedTag = () => {
  contextMenuVisible.value = false
  // 实现页面刷新逻辑
}

const closeSelectedTag = () => {
  contextMenuVisible.value = false
  removeTab(activeTab.value)
}

const closeOthersTags = () => {
  contextMenuVisible.value = false
  const currentTab = visitedViews.value.find(tab => tab.path === activeTab.value)
  if (currentTab) {
    visitedViews.value = [currentTab]
  }
}

const closeAllTags = () => {
  contextMenuVisible.value = false
  visitedViews.value = []
  router.push('/dashboard')
}

// 处理下拉菜单命令
const handleCommand = (command: string) => {
  switch (command) {
    case 'refresh':
      window.location.reload()
      break
    case 'closeOthers':
      closeOthersTags()
      break
    case 'closeLeft':
      // 关闭左侧标签
      break
    case 'closeRight':
      // 关闭右侧标签
      break
    case 'closeAll':
      closeAllTags()
      break
    case 'settings':
      // 打开标签设置
      break
  }
}
</script>

<style lang="scss" scoped>
.tabs-view {
  height: 40px;
  padding: 0;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  
  .tabs-nav {
    flex: 1;
    overflow: hidden;
    
    :deep(.el-tabs) {
      height: 100%;

      .el-tabs__header {
        margin: 0;
        height: 100%;
        border: none;
        padding: 0 16px;
        
        .el-tabs__nav-wrap {
          height: 100%;

          &::after {
            display: none;
          }

          .el-tabs__nav-scroll {
            height: 100%;

            .el-tabs__nav {
              height: 100%;
              border: none;
              padding: 4px 0;
            }
          }
        }

        .el-tabs__item {
          height: 32px;
          line-height: 32px;
          border: none !important;
          background: transparent;
          transition: all 0.3s;
          padding: 0 16px;
          font-size: 13px;
          color: #666;
          border-radius: 4px;

          &:hover {
            color: #409eff;
          }

          &.is-active {
            background-color: #ecf5ff;
            color: #409eff;
          }

          .tab-label {
            display: inline-flex;
            align-items: center;
            
            .el-icon {
              margin-right: 6px;
              font-size: 14px;
            }
          }

          .el-icon-close {
            margin-left: 4px;
            
            &:hover {
              background-color: #909399;
              color: #fff;
            }
          }

          & + .el-tabs__item {
            margin-left: 4px;
          }
        }
      }
    }
  }

  .tabs-actions {
    padding-right: 16px;
    
    .app-switch-dropdown {
      :deep(.el-dropdown-menu__item) {
        display: flex;
        align-items: center;
        padding: 8px 16px;
        
        .el-icon {
          margin-right: 8px;
          font-size: 16px;
        }
      }

      .app-switch {
        padding: 6px;
        border: none;
        background: transparent;
        color: #666;
        font-size: 16px;
        cursor: pointer;

        &:hover, &:focus {
          color: #409eff;
          background-color: #f5f7fa;
        }
      }
    }
  }

  .contextmenu {
    position: fixed;
    z-index: 100;
    width: 120px;
    margin: 0;
    padding: 8px 0;
    background: #fff;
    border: 1px solid #eee;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    
    li {
      padding: 8px 16px;
      font-size: 14px;
      cursor: pointer;
      
      &:hover {
        background: #f5f7fa;
      }
    }
  }
}
</style>
