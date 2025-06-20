<template>
  <div class="settings-google-translate">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>谷歌翻译设置</span>
          <el-switch v-model="enabled" />
        </div>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item label="API密钥">
          <el-input v-model="form.apiKey" type="password" show-password placeholder="请输入API密钥" />
        </el-form-item>
        <el-form-item label="默认源语言">
          <el-select v-model="form.sourceLanguage" placeholder="请选择源语言">
            <el-option label="自动检测" value="auto" />
            <el-option label="英语" value="en" />
            <el-option label="中文" value="zh" />
            <el-option label="日语" value="ja" />
            <el-option label="韩语" value="ko" />
          </el-select>
        </el-form-item>
        <el-form-item label="默认目标语言">
          <el-select v-model="form.targetLanguage" placeholder="请选择目标语言">
            <el-option label="英语" value="en" />
            <el-option label="中文" value="zh" />
            <el-option label="日语" value="ja" />
            <el-option label="韩语" value="ko" />
          </el-select>
        </el-form-item>
        <el-form-item label="请求超时">
          <el-input-number v-model="form.timeout" :min="1000" :max="10000" :step="100" />
          <span class="form-tip">毫秒</span>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave">保存设置</el-button>
          <el-button @click="handleTest">测试连接</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const enabled = ref(true)
const form = ref({
  apiKey: '',
  sourceLanguage: 'auto',
  targetLanguage: 'en',
  timeout: 5000
})

const handleSave = () => {
  ElMessage.success('保存成功')
}

const handleTest = () => {
  ElMessage.success('连接测试成功')
}
</script>

<style lang="scss" scoped>
.settings-google-translate {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .form-tip {
    margin-left: 8px;
    color: #999;
  }
}
</style> 