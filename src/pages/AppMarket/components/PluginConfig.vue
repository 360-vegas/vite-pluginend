<template>
  <div class="plugin-config">
    <el-form :model="config" label-width="120px">
      <el-form-item label="插件状态">
        <el-switch
          v-model="config.enabled"
          :active-text="'启用'"
          :inactive-text="'禁用'"
        />
      </el-form-item>

      <template v-if="hasCustomConfig">
        <el-divider>插件配置</el-divider>
        <component 
          :is="configComponent"
          v-model="config.custom"
          @update:modelValue="updateCustomConfig"
        />
      </template>

      <el-divider>高级设置</el-divider>
      
      <el-form-item label="自动更新">
        <el-switch
          v-model="config.autoUpdate"
          :active-text="'开启'"
          :inactive-text="'关闭'"
        />
      </el-form-item>

      <el-form-item label="更新通道">
        <el-radio-group v-model="config.updateChannel">
          <el-radio label="stable">稳定版</el-radio>
          <el-radio label="beta">测试版</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="日志级别">
        <el-select v-model="config.logLevel">
          <el-option label="错误" value="error" />
          <el-option label="警告" value="warn" />
          <el-option label="信息" value="info" />
          <el-option label="调试" value="debug" />
        </el-select>
      </el-form-item>
    </el-form>

    <div class="form-actions">
      <el-button @click="$emit('cancel')">取消</el-button>
      <el-button type="primary" @click="saveConfig">保存</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Plugin } from '../../../modules/apps/types'

const props = defineProps<{
  plugin: Plugin
}>()

const emit = defineEmits<{
  (e: 'save', config: any): void
  (e: 'cancel'): void
}>()

// 配置状态
const config = ref({
  enabled: true,
  autoUpdate: true,
  updateChannel: 'stable',
  logLevel: 'info',
  custom: null as any
})

// 是否有自定义配置组件
const hasCustomConfig = computed(() => {
  return !!props.plugin.id
})

// 动态配置组件
const configComponent = computed(() => {
  if (!props.plugin.id) return null
  // 动态导入插件的配置组件
  return `plugin-${props.plugin.id}-config`
})

// 更新自定义配置
function updateCustomConfig(value: any) {
  config.value.custom = value
}

// 保存配置
function saveConfig() {
  emit('save', config.value)
}

// 加载现有配置
onMounted(async () => {
  try {
    // 这里应该从存储中加载插件的配置
    // 暂时使用默认值
    config.value = {
      enabled: true,
      autoUpdate: true,
      updateChannel: 'stable',
      logLevel: 'info',
      custom: null
    }
  } catch (error) {
    console.error('Failed to load plugin config:', error)
  }
})
</script>

<style scoped>
.plugin-config {
  padding: 20px;
}

.form-actions {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-divider__text) {
  font-size: 14px;
  color: #666;
}
</style> 