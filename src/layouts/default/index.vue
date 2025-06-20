<template>
  <div class="layout-container" :class="{ collapsed }">
    <!-- 左侧菜单 -->
    <Sidebar :collapsed="collapsed" />
    
    <div class="layout-main">
      <!-- 顶部导航 -->
      <Header class="layout-header" :collapsed="collapsed" @toggle="toggleSidebar" />
      
      <!-- 主内容区域 -->
      <div class="layout-content hide-scrollbar">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './components/Sidebar/index.vue'
import Header from './components/Header/index.vue'

const collapsed = ref(false)

const toggleSidebar = () => {
  collapsed.value = !collapsed.value
}
</script>

<style lang="scss" scoped>
.layout-container {
  display: flex;
  width: 100%;
  height: 100vh;
  
  .layout-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    margin-left: 284px; // 64px (主导航) + 220px (子导航)
    
    .layout-header {
      height: 57px;
      background: #fff;
      box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    }
    
    .layout-content {
      flex: 1;
      padding: 0;
      overflow: auto;
      background: transparent;
      box-sizing: border-box;
      width: 100%;
    }
  }

  &.collapsed {
    .layout-main {
      margin-left: 64px;
    }
  }
}
</style>
