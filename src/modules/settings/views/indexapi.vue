<template>
  <div class="api-overview">
    <!-- 数据概览卡片 -->
    <el-row :gutter="8">
      <el-col :span="6" v-for="(card, index) in statsCards" :key="index">
        <el-card shadow="hover">
          <div class="stats-card">
            <el-icon :size="40" :class="card.color">
              <component :is="card.icon" />
            </el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ card.value }}</div>
              <div class="stats-label">{{ card.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- API 服务状态和调用记录 -->
    <el-row :gutter="8" class="mt-4">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>服务状态</span>
              <el-tag type="success">正常</el-tag>
            </div>
          </template>
          <div class="service-status">
            <div class="service-item" v-for="(service, index) in services" :key="index">
              <div class="service-info">
                <span class="service-name">{{ service.name }}</span>
                <el-tag :type="service.status === 'normal' ? 'success' : 'danger'" size="small">
                  {{ service.status === 'normal' ? '正常' : '异常' }}
                </el-tag>
              </div>
              <div class="service-stats">
                <div class="stat-item">
                  <span class="label">成功率</span>
                  <el-progress 
                    :percentage="service.successRate" 
                    :status="service.successRate < 90 ? 'exception' : 'success'"
                  />
                </div>
                <div class="stat-item">
                  <span class="label">响应时间</span>
                  <span class="value">{{ service.responseTime }}ms</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近调用</span>
              <el-button type="primary" link>查看更多</el-button>
            </div>
          </template>
          <el-table :data="recentCalls" style="width: 100%">
            <el-table-column prop="api" label="API" />
            <el-table-column prop="method" label="方法" width="100">
              <template #default="{ row }">
                <el-tag :type="row.method === 'GET' ? 'success' : 'warning'" size="small">
                  {{ row.method }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="调用时间" width="180" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 调用趋势 -->
    <el-card shadow="hover" class="mt-4">
      <template #header>
        <div class="card-header">
          <span>调用趋势</span>
          <el-radio-group v-model="timeRange" size="small">
            <el-radio-button label="today">今日</el-radio-button>
            <el-radio-button label="week">本周</el-radio-button>
            <el-radio-button label="month">本月</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <div class="chart-placeholder" style="height: 300px">
        <el-empty description="调用趋势图" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  Connection,
  Monitor,
  Timer,
  DataLine
} from '@element-plus/icons-vue'

const timeRange = ref('today')

// 数据概览卡片
const statsCards = [
  { label: '总调用次数', value: '128,546', icon: 'Connection', color: 'text-primary' },
  { label: '今日调用', value: '12,345', icon: 'Monitor', color: 'text-success' },
  { label: '平均响应', value: '85ms', icon: 'Timer', color: 'text-warning' },
  { label: '成功率', value: '99.9%', icon: 'DataLine', color: 'text-danger' }
]

// API 服务状态
const services = [
  { 
    name: '谷歌翻译', 
    status: 'normal',
    successRate: 99.8,
    responseTime: 150
  },
  { 
    name: 'OpenAI', 
    status: 'normal',
    successRate: 99.5,
    responseTime: 800
  },
  { 
    name: '火山引擎', 
    status: 'normal',
    successRate: 99.9,
    responseTime: 120
  },
  { 
    name: 'Azure AI', 
    status: 'normal',
    successRate: 99.7,
    responseTime: 200
  }
]

// 最近调用记录
const recentCalls = [
  { api: '谷歌翻译', method: 'POST', time: '2024-02-22 15:30:00' },
  { api: 'OpenAI Chat', method: 'POST', time: '2024-02-22 15:28:30' },
  { api: '图像识别', method: 'POST', time: '2024-02-22 15:25:12' },
  { api: '语音合成', method: 'GET', time: '2024-02-22 15:20:45' }
]
</script>

<style lang="scss" scoped>
.api-overview {
  padding: 6px;
  box-sizing: border-box;
  width: 100%;

  .el-row {
    margin: 0 -8px 8px;
    
    &:last-child {
      margin-bottom: 0;
    }
  }

  .el-col {
    padding: 0 8px;
    box-sizing: border-box;
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

  .stats-card {
    display: flex;
    align-items: center;
    padding: 8px;

    .stats-info {
      margin-left: 16px;

      .stats-value {
        font-size: 22px;
        font-weight: bold;
        line-height: 1;
        margin-bottom: 6px;
      }

      .stats-label {
        font-size: 13px;
        color: #666;
      }
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .service-status {
    .service-item {
      padding: 16px;
      border-bottom: 1px solid #f0f0f0;

      &:last-child {
        border-bottom: none;
      }

      .service-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;

        .service-name {
          font-weight: bold;
          color: #333;
        }
      }

      .service-stats {
        .stat-item {
          display: flex;
          align-items: center;
          margin-top: 8px;

          .label {
            width: 80px;
            color: #666;
          }

          .value {
            color: #333;
          }

          :deep(.el-progress) {
            width: 200px;
            margin-left: 8px;
          }
        }
      }
    }
  }

  .chart-placeholder {
    min-height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .mt-4 {
    margin-top: 8px;
  }

  // 文字颜色
  .text-primary { color: var(--el-color-primary); }
  .text-success { color: var(--el-color-success); }
  .text-warning { color: var(--el-color-warning); }
  .text-danger { color: var(--el-color-danger); }
}
</style> 