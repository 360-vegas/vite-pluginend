<template>
  <div class="dependency-checker">
    <el-dialog
      v-model="visible"
      title="插件依赖检查"
      width="800px"
      :before-close="handleClose"
    >
      <div v-if="loading" class="loading-container">
        <el-loading-spinner />
        <p>正在检查依赖...</p>
      </div>

      <div v-else-if="checkResult" class="check-result">
        <!-- 总体状态 -->
        <div class="status-header">
          <el-tag 
            :type="getStatusType(checkResult.overall_status)" 
            size="large"
            effect="dark"
          >
            {{ getStatusText(checkResult.overall_status) }}
          </el-tag>
          <span class="plugin-name">{{ pluginKey }}</span>
        </div>

        <!-- 数据库依赖 -->
        <div v-if="checkResult.database" class="dependency-section">
          <h3>
            <el-icon><Coin /></el-icon>
            数据库依赖
          </h3>
          
          <el-card shadow="never" class="dependency-card">
            <div class="dependency-item">
              <div class="item-header">
                <span class="item-title">{{ checkResult.database.requirement.type.toUpperCase() }}</span>
                <el-tag :type="getStatusType(checkResult.database.status)" size="small">
                  {{ checkResult.database.status }}
                </el-tag>
              </div>
              <p class="item-message">{{ checkResult.database.message }}</p>
              
              <!-- 数据库设置选项 -->
              <div v-if="checkResult.database.setup_options" class="setup-options">
                <el-divider>数据库配置</el-divider>
                
                <el-form :model="dbSetupForm" label-width="120px">
                  <el-form-item label="数据库名称">
                    <el-input 
                      v-model="dbSetupForm.suggested_database_name"
                      placeholder="输入数据库名称"
                    />
                  </el-form-item>
                  
                  <el-form-item v-if="checkResult.database.setup_options.available_databases.length > 0" label="现有数据库">
                    <el-select v-model="selectedExistingDb" placeholder="选择现有数据库" clearable>
                      <el-option
                        v-for="db in checkResult.database.setup_options.available_databases"
                        :key="db"
                        :label="db"
                        :value="db"
                      />
                    </el-select>
                  </el-form-item>
                  
                  <el-form-item label="操作选择">
                    <el-radio-group v-model="dbSetupForm.action">
                      <el-radio value="create_new">创建新数据库</el-radio>
                      <el-radio 
                        v-if="checkResult.database.setup_options.available_databases.length > 0"
                        value="use_existing"
                      >
                        使用现有数据库
                      </el-radio>
                    </el-radio-group>
                  </el-form-item>
                </el-form>
                
                <div class="setup-actions">
                  <el-button 
                    type="primary" 
                    @click="setupDatabase"
                    :loading="settingUpDb"
                  >
                    配置数据库
                  </el-button>
                </div>
              </div>
            </div>
          </el-card>
        </div>

        <!-- 环境变量 -->
        <div v-if="checkResult.environment && checkResult.environment.length > 0" class="dependency-section">
          <h3>
            <el-icon><Setting /></el-icon>
            环境变量
          </h3>
          
          <el-card shadow="never" class="dependency-card">
            <div v-for="env in checkResult.environment" :key="env.requirement.name" class="dependency-item">
              <div class="item-header">
                <span class="item-title">{{ env.requirement.name }}</span>
                <el-tag :type="getStatusType(env.status)" size="small">
                  {{ env.status }}
                </el-tag>
              </div>
              <p class="item-message">{{ env.message }}</p>
              <p v-if="env.requirement.description" class="item-description">
                {{ env.requirement.description }}
              </p>
              <p v-if="env.current" class="current-value">
                当前值: {{ env.current }}
              </p>
            </div>
          </el-card>
        </div>

        <!-- 基础依赖 -->
        <div v-if="checkResult.dependencies && checkResult.dependencies.length > 0" class="dependency-section">
          <h3>
            <el-icon><Box /></el-icon>
            基础依赖
          </h3>
          
          <el-card shadow="never" class="dependency-card">
            <div v-for="dep in checkResult.dependencies" :key="dep.dependency.name" class="dependency-item">
              <div class="item-header">
                <span class="item-title">{{ dep.dependency.name }}</span>
                <el-tag :type="getStatusType(dep.status)" size="small">
                  {{ dep.status }}
                </el-tag>
              </div>
              <p class="item-message">{{ dep.message }}</p>
              <p v-if="dep.dependency.description" class="item-description">
                {{ dep.dependency.description }}
              </p>
            </div>
          </el-card>
        </div>

        <!-- 建议 -->
        <div v-if="checkResult.suggestions && checkResult.suggestions.length > 0" class="suggestions">
          <h3>
            <el-icon><InfoFilled /></el-icon>
            建议
          </h3>
          <ul>
            <li v-for="(suggestion, index) in checkResult.suggestions" :key="index">
              {{ suggestion }}
            </li>
          </ul>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleClose">关闭</el-button>
          <el-button 
            v-if="checkResult && checkResult.can_install"
            type="primary" 
            @click="handleInstall"
          >
            继续安装
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Coin,
  Setting,
  Box,
  InfoFilled
} from '@element-plus/icons-vue'

interface Props {
  modelValue: boolean
  pluginKey: string
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'install-confirmed'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = ref(false)
const loading = ref(false)
const settingUpDb = ref(false)
const checkResult = ref<any>(null)
const selectedExistingDb = ref('')

const dbSetupForm = ref({
  suggested_database_name: '',
  action: 'create_new'
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.pluginKey) {
    checkDependencies()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

// 检查依赖
const checkDependencies = async () => {
  loading.value = true
  try {
    const response = await fetch(`/api/plugins/${props.pluginKey}/dependencies/check`)
    const data = await response.json()
    
    if (data.success) {
      checkResult.value = data.data
      if (checkResult.value.database?.setup_options) {
        dbSetupForm.value.suggested_database_name = 
          checkResult.value.database.setup_options.suggested_database_name
      }
    } else {
      ElMessage.error('检查依赖失败: ' + data.message)
    }
  } catch (error) {
    ElMessage.error('检查依赖时发生错误')
    console.error('Error checking dependencies:', error)
  } finally {
    loading.value = false
  }
}

// 设置数据库
const setupDatabase = async () => {
  settingUpDb.value = true
  try {
    const config = {
      suggested_database_name: dbSetupForm.value.action === 'use_existing' 
        ? selectedExistingDb.value 
        : dbSetupForm.value.suggested_database_name,
      create_new_database: dbSetupForm.value.action === 'create_new',
      use_existing_database: dbSetupForm.value.action === 'use_existing'
    }

    const response = await fetch(`/api/plugins/${props.pluginKey}/dependencies/setup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(config)
    })

    const data = await response.json()
    
    if (data.success) {
      ElMessage.success('数据库配置成功')
      // 重新检查依赖
      await checkDependencies()
    } else {
      ElMessage.error('数据库配置失败: ' + data.message)
    }
  } catch (error) {
    ElMessage.error('配置数据库时发生错误')
    console.error('Error setting up database:', error)
  } finally {
    settingUpDb.value = false
  }
}

// 获取状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'success':
    case 'available':
    case 'set':
      return 'success'
    case 'warning':
    case 'setup_required':
    case 'version_mismatch':
      return 'warning'
    case 'error':
    case 'missing':
    case 'unreachable':
      return 'danger'
    default:
      return 'info'
  }
}

// 获取状态文本
const getStatusText = (status: string) => {
  switch (status) {
    case 'success':
      return '依赖检查通过'
    case 'warning':
      return '存在警告'
    case 'error':
      return '依赖检查失败'
    default:
      return '未知状态'
  }
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
  checkResult.value = null
}

// 继续安装
const handleInstall = () => {
  emit('install-confirmed')
  handleClose()
}
</script>

<style scoped>
.dependency-checker {
  .loading-container {
    text-align: center;
    padding: 40px;
  }

  .status-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 24px;
    padding: 16px;
    background-color: #f5f7fa;
    border-radius: 8px;

    .plugin-name {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }

  .dependency-section {
    margin-bottom: 24px;

    h3 {
      display: flex;
      align-items: center;
      gap: 8px;
      margin: 0 0 12px 0;
      color: #303133;
      font-size: 16px;
      font-weight: 600;
    }
  }

  .dependency-card {
    border: 1px solid #e4e7ed;
    border-radius: 8px;
  }

  .dependency-item {
    padding: 16px;

    &:not(:last-child) {
      border-bottom: 1px solid #f0f2f5;
    }

    .item-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;

      .item-title {
        font-weight: 600;
        color: #303133;
      }
    }

    .item-message {
      margin: 0 0 8px 0;
      color: #606266;
    }

    .item-description {
      margin: 0 0 8px 0;
      color: #909399;
      font-size: 12px;
    }

    .current-value {
      margin: 0;
      color: #67c23a;
      font-size: 12px;
      font-family: monospace;
    }
  }

  .setup-options {
    margin-top: 16px;
    padding-top: 16px;

    .setup-actions {
      margin-top: 16px;
    }
  }

  .suggestions {
    margin-top: 24px;
    padding: 16px;
    background-color: #ecf5ff;
    border-radius: 8px;

    h3 {
      margin: 0 0 12px 0;
      color: #409eff;
      font-size: 14px;
    }

    ul {
      margin: 0;
      padding-left: 20px;
      color: #606266;
    }

    li {
      margin-bottom: 4px;
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}
</style> 