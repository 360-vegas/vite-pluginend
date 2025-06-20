<template>
  <div class="market-plugins">
    <!-- 插件安装区域 -->
    <el-card class="install-section" shadow="hover">
      <template #header>
        <div class="section-header">
          <el-icon><Upload /></el-icon>
          <span>安装插件包</span>
        </div>
      </template>
      
      <div class="install-content">
        <el-upload
          ref="uploadRef"
          class="upload-demo"
          drag
          accept=".zip"
          :auto-upload="false"
          :on-change="handleFileChange"
          :file-list="fileList"
          :limit="1"
          :on-exceed="handleExceed"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将插件包拖拽到此处或 <em>点击选择文件</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              仅支持 .zip 格式的插件包，文件大小不超过 50MB
            </div>
          </template>
        </el-upload>
        
        <div class="install-actions" v-if="selectedFile">
          <el-button 
            type="primary" 
            @click="installLocalPlugin"
            :loading="installing"
          >
            <el-icon><Download /></el-icon>
            安装插件
          </el-button>
          <el-button @click="clearFile">
            <el-icon><Close /></el-icon>
            清除
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 在线插件市场 -->
    <el-card class="market-section" shadow="hover">
      <template #header>
        <div class="section-header">
          <el-icon><Shop /></el-icon>
          <span>在线插件市场</span>
          <el-tag size="small" type="info">即将上线</el-tag>
        </div>
      </template>
      
      <el-table :data="marketPlugins" style="width: 100%">
        <el-table-column prop="name" label="插件名称" />
        <el-table-column prop="version" label="最新版本" width="100" />
        <el-table-column prop="desc" label="简介" />
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button size="small" type="primary" @click="installOnlinePlugin(scope.row)" disabled>
              安装
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 依赖检查对话框 -->
    <DependencyChecker
      v-model="showDependencyChecker"
      :plugin-key="pendingPluginKey"
      @install-confirmed="proceedWithInstallation"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Upload, 
  UploadFilled, 
  Download, 
  Close, 
  Shop 
} from '@element-plus/icons-vue'
import { appsApi } from '@/api'
import type { UploadFile, UploadFiles } from 'element-plus'
import DependencyChecker from './components/DependencyChecker.vue'

const marketPlugins = ref([
  { name: '流程引擎', version: '1.2.0', desc: '可视化流程设计与自动化' },
  { name: '报表工具', version: '0.8.1', desc: '多维度数据报表生成' },
])

const fileList = ref<UploadFiles>([])
const selectedFile = ref<File | null>(null)
const installing = ref(false)
const uploadRef = ref()
const showDependencyChecker = ref(false)
const pendingPluginKey = ref('')

// 处理文件选择
const handleFileChange = (file: UploadFile) => {
  console.log('📁 选择文件:', file.name)
  
  if (file.raw) {
    // 验证文件大小 (50MB)
    const maxSize = 50 * 1024 * 1024
    if (file.raw.size > maxSize) {
      ElMessage.error('文件大小不能超过 50MB')
      fileList.value = []
      selectedFile.value = null
      return
    }
    
    selectedFile.value = file.raw
    ElMessage.success(`已选择文件: ${file.name}`)
  }
}

// 处理文件数量超限
const handleExceed = () => {
  ElMessage.warning('只能选择一个插件包文件')
}

// 安装本地插件 - 先检查依赖
const installLocalPlugin = async () => {
  if (!selectedFile.value) {
    ElMessage.error('请先选择插件包文件')
    return
  }

  // 这里可以先尝试提取插件key来检查依赖
  // 为了简化演示，我们先直接安装，后续可以增强
  await performInstallation()
}

// 执行实际的安装操作
const performInstallation = async () => {
  if (!selectedFile.value) {
    ElMessage.error('请先选择插件包文件')
    return
  }

  try {
    installing.value = true
    
    console.log('🚀 开始安装插件包:', selectedFile.value.name)
    
    const response = await appsApi.installPluginPackage(selectedFile.value) as any
    
    console.log('✅ 插件安装响应:', response)
    
    if (response.success) {
      ElMessage.success(response.message || '插件安装成功！')
      
      // 显示安装成功的详情
      ElMessageBox.alert(
        `插件 "${response.data?.pluginKey}" 安装成功！页面将刷新以加载新插件。`,
        '安装成功',
        {
          confirmButtonText: '确定',
          type: 'success',
          callback: () => {
            // 清除文件选择
            clearFile()
            // 刷新页面
            window.location.reload()
          }
        }
      )
    } else {
      throw new Error(response.message || '安装失败')
    }
  } catch (error: any) {
    console.error('❌ 插件安装失败:', error)
    ElMessage.error(error.message || '安装插件失败，请检查插件包格式')
  } finally {
    installing.value = false
  }
}

// 清除文件选择
const clearFile = () => {
  fileList.value = []
  selectedFile.value = null
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

// 安装在线插件（暂未实现）
const installOnlinePlugin = (plugin: any) => {
  ElMessage.info(`在线插件安装功能正在开发中，敬请期待！`)
}

// 检查插件依赖
const checkPluginDependencies = async (pluginKey: string) => {
  pendingPluginKey.value = pluginKey
  showDependencyChecker.value = true
}

// 确认安装后继续
const proceedWithInstallation = async () => {
  await performInstallation()
}
</script>

<style lang="scss" scoped>
.market-plugins {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.install-section {
  margin-bottom: 24px;

  .section-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    
    .el-icon {
      color: #409eff;
    }
  }
}

.install-content {
  .upload-demo {
    margin-bottom: 16px;
    
    :deep(.el-upload-dragger) {
      border: 2px dashed #dcdfe6;
      border-radius: 8px;
      padding: 40px;
      transition: all 0.3s;
      
      &:hover {
        border-color: #409eff;
        background-color: #ecf5ff;
      }
    }
    
    :deep(.el-icon--upload) {
      font-size: 48px;
      color: #c0c4cc;
      margin-bottom: 16px;
    }
    
    :deep(.el-upload__text) {
      color: #606266;
      font-size: 14px;
      
      em {
        color: #409eff;
        font-style: normal;
      }
    }
    
    :deep(.el-upload__tip) {
      font-size: 12px;
      color: #909399;
      margin-top: 8px;
    }
  }
  
  .install-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
    margin-top: 20px;
  }
}

.market-section {
  .section-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    
    .el-icon {
      color: #67c23a;
    }
    
    .el-tag {
      margin-left: auto;
    }
  }
  
  :deep(.el-table) {
    .el-button {
      margin-right: 8px;
      
      &:last-child {
        margin-right: 0;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .market-plugins {
    padding: 16px;
  }
  
  .install-content {
    .upload-demo :deep(.el-upload-dragger) {
      padding: 20px;
    }
    
    .install-actions {
      flex-direction: column;
      align-items: center;
    }
  }
}
</style> 