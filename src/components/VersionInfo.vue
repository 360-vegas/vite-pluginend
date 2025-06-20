<template>
  <div class="version-info">
    <el-tooltip
      :content="`版本: ${version}\n构建时间: ${buildTime}`"
      placement="bottom"
    >
      <el-tag size="small" type="info">
        v{{ version }}
      </el-tag>
    </el-tooltip>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { getCurrentVersion, getCurrentBuildTime } from '@/services/version-control'

const version = computed(() => getCurrentVersion())
const buildTime = computed(() => getCurrentBuildTime())
</script>

<style lang="scss" scoped>
.version-info {
  display: inline-block;
  margin-right: 10px;
}
</style> 