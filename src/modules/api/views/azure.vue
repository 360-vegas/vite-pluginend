<template>
  <div class="settings-azure">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Azure设置</span>
          <el-switch v-model="enabled" />
        </div>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item label="API密钥">
          <el-input v-model="form.apiKey" type="password" show-password placeholder="请输入API密钥" />
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="form.region" placeholder="请选择区域">
            <el-option label="东亚" value="east-asia" />
            <el-option label="东南亚" value="southeast-asia" />
            <el-option label="美国东部" value="east-us" />
            <el-option label="美国西部" value="west-us" />
          </el-select>
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
  region: 'east-asia',
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
.settings-azure {
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