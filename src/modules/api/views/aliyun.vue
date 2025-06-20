<template>
  <div class="settings-aliyun">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>阿里云设置</span>
          <el-switch v-model="enabled" />
        </div>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item label="AccessKey">
          <el-input v-model="form.accessKey" placeholder="请输入AccessKey" />
        </el-form-item>
        <el-form-item label="AccessSecret">
          <el-input v-model="form.accessSecret" type="password" show-password placeholder="请输入AccessSecret" />
        </el-form-item>
        <el-form-item label="服务选择">
          <el-checkbox-group v-model="form.services">
            <el-checkbox label="oss">对象存储</el-checkbox>
            <el-checkbox label="sms">短信服务</el-checkbox>
            <el-checkbox label="voice">语音服务</el-checkbox>
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
  accessKey: '',
  accessSecret: '',
  services: []
})

const handleSave = () => {
  ElMessage.success('保存成功')
}

const handleTest = () => {
  ElMessage.success('连接测试成功')
}
</script>

<style lang="scss" scoped>
.settings-aliyun {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .el-checkbox-group {
    display: flex;
    gap: 20px;
  }
}
</style> 