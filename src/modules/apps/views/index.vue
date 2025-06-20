<template>
  <div class="apps-overview">
    <!-- 应用分类卡片 -->
    <el-row :gutter="16">
      <el-col :span="6" v-for="(category, index) in appCategories" :key="index">
        <el-card shadow="hover" class="category-card">
          <div class="category-header">
            <el-icon :size="24" :class="category.color">
              <component :is="category.icon" />
            </el-icon>
            <span class="category-title">{{ category.title }}</span>
          </div>
          <div class="category-stats">
            <div class="stat-item">
              <span class="label">总数</span>
              <span class="value">{{ category.total }}</span>
            </div>
            <div class="stat-item">
              <span class="label">运行中</span>
              <span class="value success">{{ category.running }}</span>
            </div>
            <div class="stat-item">
              <span class="label">已停止</span>
              <span class="value warning">{{ category.stopped }}</span>
            </div>
          </div>
          <div class="category-actions">
            <el-button type="primary" link @click="handleCreateApp(category.type)">
              创建{{ category.title }}
            </el-button>
            <el-button type="primary" link @click="handleViewAll(category.type)">
              查看全部
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 应用统计图表 -->
    <el-row :gutter="16" class="mt-4">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>应用运行状态</span>
            </div>
          </template>
          <div class="chart-container">
            <!-- 这里可以添加饼图或柱状图展示应用状态分布 -->
            <div class="chart-placeholder">应用状态统计图表</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>应用类型分布</span>
            </div>
          </template>
          <div class="chart-container">
            <!-- 这里可以添加饼图展示应用类型分布 -->
            <div class="chart-placeholder">应用类型分布图表</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近活动 -->
    <el-card shadow="hover" class="mt-4">
      <template #header>
        <div class="card-header">
          <span>最近活动</span>
        </div>
      </template>
      <el-timeline>
        <el-timeline-item
          v-for="(activity, index) in recentActivities"
          :key="index"
          :timestamp="activity.time"
          :type="activity.type"
        >
          {{ activity.content }}
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { appsApi } from '@/api'
import { ElMessage } from 'element-plus'
import {
  Grid,
  Link,
  Connection,
  Share
} from '@element-plus/icons-vue'

// 应用分类配置
const appCategories = ref([
  {
    type: 'sites',
    title: '站群应用',
    icon: 'Grid',
    color: 'text-primary',
    total: 0,
    running: 0,
    stopped: 0
  },
  {
    type: 'external',
    title: '外链应用',
    icon: 'Link',
    color: 'text-success',
    total: 0,
    running: 0,
    stopped: 0
  },
  {
    type: 'spider',
    title: '蜘蛛池应用',
    icon: 'Connection',
    color: 'text-warning',
    total: 0,
    running: 0,
    stopped: 0
  },
  {
    type: 'proxy',
    title: '反代应用',
    icon: 'Share',
    color: 'text-danger',
    total: 0,
    running: 0,
    stopped: 0
  }
])

// 最近活动数据
const recentActivities = ref([])

// 获取应用统计数据
const fetchAppStats = async () => {
  try {
    const data = await appsApi.getAppStats()
    // 更新应用分类数据
    appCategories.value = appCategories.value.map(category => {
      const stats = data.categories[category.type] || { total: 0, running: 0, stopped: 0 }
      return {
        ...category,
        ...stats
      }
    })
    // 更新最近活动
    recentActivities.value = data.recentActivities || []
  } catch (error) {
    console.error('获取应用统计数据失败:', error)
    ElMessage.error('获取应用统计数据失败')
  }
}

// 处理创建应用
const handleCreateApp = (type: string) => {
  // TODO: 实现创建应用的逻辑
}

// 处理查看全部
const handleViewAll = (type: string) => {
  // TODO: 实现查看全部的逻辑
}

// 页面加载时获取数据
onMounted(() => {
  fetchAppStats()
})
</script>

<style lang="scss" scoped>
.apps-overview {
  padding: 6px;
  box-sizing: border-box;
  width: 100%;

  .el-row {
    margin: 0 -8px 16px;
  }

  .el-col {
    padding: 0 8px;
  }

  .category-card {
    .category-header {
      display: flex;
      align-items: center;
      margin-bottom: 16px;

      .el-icon {
        margin-right: 8px;
      }

      .category-title {
        font-size: 16px;
        font-weight: 500;
      }
    }

    .category-stats {
      display: flex;
      justify-content: space-between;
      margin-bottom: 16px;

      .stat-item {
        display: flex;
        flex-direction: column;
        align-items: center;

        .label {
          font-size: 12px;
          color: #666;
          margin-bottom: 4px;
        }

        .value {
          font-size: 18px;
          font-weight: bold;

          &.success {
            color: #67c23a;
          }

          &.warning {
            color: #e6a23c;
          }
        }
      }
    }

    .category-actions {
      display: flex;
      justify-content: space-between;
    }
  }

  .el-card {
    margin-bottom: 8px;
  }

  :deep(.el-card__header) {
    padding: 8px 12px;
  }

  :deep(.el-card__body) {
    padding: 12px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .chart-container {
    height: 300px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f7fa;
    border-radius: 4px;

    .chart-placeholder {
      color: #909399;
      font-size: 14px;
    }
  }

  .mt-4 {
    margin-top: 16px;
  }
}
</style> 