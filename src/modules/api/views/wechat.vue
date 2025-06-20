<template>
  <div class="settings-wechat">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>微信设置</span>
          <el-switch v-model="enabled" />
        </div>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item label="AppID">
          <el-input v-model="form.appId" placeholder="请输入AppID" />
        </el-form-item>
        <el-form-item label="AppSecret">
          <el-input v-model="form.appSecret" type="password" show-password placeholder="请输入AppSecret" />
        </el-form-item>
        <el-form-item label="Token">
          <el-input v-model="form.token" placeholder="请输入Token" />
        </el-form-item>
        <el-form-item label="EncodingAESKey">
          <el-input v-model="form.encodingAESKey" placeholder="请输入EncodingAESKey" />
        </el-form-item>
        <el-form-item label="消息加密">
          <el-switch v-model="form.encrypt" />
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
  appId: '',
  appSecret: '',
  token: '',
  encodingAESKey: '',
  encrypt: false
})

const handleSave = () => {
  ElMessage.success('保存成功')
}

const handleTest = () => {
  ElMessage.success('连接测试成功')
}
</script>

<style lang="scss" scoped>
.settings-wechat {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style> 