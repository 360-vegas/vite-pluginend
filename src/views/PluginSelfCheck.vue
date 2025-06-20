<template>
    <div>
      <el-button @click="check">开始自检</el-button>
      <el-table :data="tableData" style="width: 100%; margin-top: 20px;">
        <el-table-column prop="item" label="检测项"/>
        <el-table-column prop="status" label="状态"/>
        <el-table-column prop="desc" label="说明"/>
      </el-table>
      <el-card style="margin-top: 20px;">
        <div style="font-weight:bold;">自检日志信息：</div>
        <el-input
          type="textarea"
          :rows="10"
          v-model="logs"
          readonly
          style="margin-top: 10px;"
        />
      </el-card>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import axios from 'axios'
  
  const tableData = ref([])
  const logs = ref('')
  
  const check = async () => {
    const res = await axios.get('/api/self-check')
    tableData.value = [
      { item: '插件目录写权限', status: res.data.writePermission ? '✔' : '✖', desc: res.data.writePermission ? '可写' : '不可写' },
      { item: 'Node.js', status: res.data.node ? '✔' : '✖', desc: res.data.node ? '可用' : '不可用' },
      { item: 'npm', status: res.data.npm ? '✔' : '✖', desc: res.data.npm ? '可用' : '不可用' },
      { item: '依赖', status: res.data.dependencies ? '✔' : '✖', desc: res.data.dependencies ? '已安装' : '未安装' },
      { item: 'Vite脚本', status: res.data.viteScript ? '✔' : '✖', desc: res.data.viteScript ? '存在' : '缺失' },
      { item: 'MongoDB', status: res.data.mongodb ? '✔' : '⚠', desc: res.data.mongodb ? '已运行' : '未运行' },
      { item: 'MySQL', status: res.data.mysql ? '✔' : '⚠', desc: res.data.mysql ? '已运行' : '未运行' },
      { item: 'create-plugin.cjs', status: res.data.createPluginScript ? '✔' : '⚠', desc: res.data.createPluginScript ? '存在' : '缺失' }
    ]
    logs.value = res.data.logs || ''
  }
  </script>