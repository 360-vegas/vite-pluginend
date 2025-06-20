<template>
  <div class="app-market">
    <div class="app-market-header">
      <el-tabs 
        v-model="activeTab" 
        class="app-market-tabs"
        @tab-click="handleTabClick"
      >
        <el-tab-pane name="installed" label="已安装插件">
          <template #label>
            <el-icon><Files /></el-icon>
            <span>已安装插件</span>
          </template>
        </el-tab-pane>
        <el-tab-pane name="market" label="插件市场">
          <template #label>
            <el-icon><Collection /></el-icon>
            <span>插件市场</span>
          </template>
        </el-tab-pane>
      </el-tabs>
    </div>

    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Files, Collection } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const activeTab = ref('installed')

// 根据路由更新当前标签
watch(() => route.path, (path) => {
  if (path.includes('/app-market/installed')) {
    activeTab.value = 'installed'
  } else if (path.includes('/app-market/market')) {
    activeTab.value = 'market'
  }
}, { immediate: true })

// 处理标签点击
const handleTabClick = (tab: any) => {
  router.push(`/app-market/${tab.props.name}`)
}
</script>

<style lang="scss" scoped>
.app-market {
  height: 100%;
  display: flex;
  flex-direction: column;

  .app-market-header {
    padding: 16px 20px;
    background-color: #fff;
    border-bottom: 1px solid #e4e7ed;
  }

  .app-market-tabs {
    :deep(.el-tabs__header) {
      margin: 0;
    }

    :deep(.el-tabs__nav-wrap) {
      &::after {
        display: none;
      }
    }

    :deep(.el-tabs__item) {
      .el-icon {
        margin-right: 8px;
        vertical-align: middle;
      }

      span {
        vertical-align: middle;
      }
    }
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.3s ease;
  }

  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }
}
</style> 