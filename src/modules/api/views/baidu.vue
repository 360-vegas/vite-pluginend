<template>
  <div class="settings-baidu">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>百度设置</span>
          <el-switch v-model="enabled" />
        </div>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item label="API密钥">
          <el-input v-model="form.apiKey" type="password" show-password placeholder="请输入API密钥" />
        </el-form-item>
        <el-form-item label="Secret Key">
          <el-input v-model="form.secretKey" type="password" show-password placeholder="请输入Secret Key" />
        </el-form-item>
        <el-form-item label="服务选择">
          <el-checkbox-group v-model="form.services">
            <el-checkbox label="translate">翻译服务</el-checkbox>
            <el-checkbox label="tts">语音合成</el-checkbox>
            <el-checkbox label="asr">语音识别</el-checkbox>
          </el-checkbox-group>
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

const enabled = ref(false)
const form = ref({
  apiKey: '',
  secretKey: '',
  services: ['translate']
})

const handleSave = () => {
  ElMessage.success('保存成功')
}

const handleTest = () => {
  ElMessage.success('连接测试成功')
}
</script>

<style lang="scss" scoped>
.settings-baidu {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .el-checkbox-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
}
</style> 