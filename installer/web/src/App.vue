<template>
  <div class="installer-app">
    <div class="installer-container">
      <!-- 头部 -->
      <div class="header">
        <h1>🚀 插件系统安装向导</h1>
        <p>欢迎安装跨平台插件系统</p>
      </div>

      <!-- 主要内容区域 -->
      <div class="main-content">
        <!-- 欢迎页面 -->
        <WelcomePage 
          v-if="currentStep === 'welcome'"
          @next="handleWelcomeNext"
        />

        <!-- 依赖检查页面 -->
        <DependencyCheck 
          v-if="currentStep === 'dependencies'"
          @next="handleDependencyNext"
          @back="currentStep = 'welcome'"
        />

        <!-- 数据库配置页面 -->
        <DatabaseConfig 
          v-if="currentStep === 'database'"
          @next="handleDatabaseNext"
          @back="currentStep = 'dependencies'"
        />

        <!-- 安装进度页面 -->
        <InstallationProgress 
          v-if="currentStep === 'installation'"
          @complete="handleInstallComplete"
          @back="currentStep = 'database'"
        />

        <!-- 完成页面 -->
        <CompletePage 
          v-if="currentStep === 'complete'"
          @finish="handleFinish"
        />
      </div>

      <!-- 底部 -->
      <div class="footer">
        <p>版本 1.0.0 | 支持 Windows、macOS、Linux</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import WelcomePage from './components/WelcomePage.vue'
import DependencyCheck from './components/DependencyCheck.vue'
import DatabaseConfig from './components/DatabaseConfig.vue'
import InstallationProgress from './components/InstallationProgress.vue'
import CompletePage from './components/CompletePage.vue'

type Step = 'welcome' | 'dependencies' | 'database' | 'installation' | 'complete'

const currentStep = ref<Step>('welcome')

const handleWelcomeNext = () => {
  currentStep.value = 'dependencies'
}

const handleDependencyNext = () => {
  currentStep.value = 'database'
}

const handleDatabaseNext = () => {
  currentStep.value = 'installation'
}

const handleInstallComplete = () => {
  currentStep.value = 'complete'
}

const handleFinish = () => {
  // 关闭安装程序或跳转到应用
  window.close()
}
</script>

<style scoped>
.installer-app {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.installer-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  max-width: 800px;
  width: 100%;
  overflow: hidden;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 30px;
  text-align: center;
}

.header h1 {
  margin: 0 0 10px 0;
  font-size: 28px;
  font-weight: 600;
}

.header p {
  margin: 0;
  opacity: 0.9;
  font-size: 16px;
}

.main-content {
  min-height: 500px;
  padding: 30px;
}

.footer {
  background: #f8f9fa;
  padding: 15px 30px;
  text-align: center;
  border-top: 1px solid #e9ecef;
}

.footer p {
  margin: 0;
  color: #6c757d;
  font-size: 14px;
}
</style> 