<template>
  <div class="external-link-stats">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <span>总链接数</span>
            </div>
          </template>
          <div class="stat-value">{{ stats.total_links }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <span>总点击量</span>
            </div>
          </template>
          <div class="stat-value">{{ stats.total_clicks }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <span>活跃链接</span>
            </div>
          </template>
          <div class="stat-value">{{ stats.active_links }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <span>分类数量</span>
            </div>
          </template>
          <div class="stat-value">{{ Object.keys(stats.categories).length }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>点击趋势</span>
              <el-radio-group v-model="trendPeriod" size="small">
                <el-radio-button label="week">周</el-radio-button>
                <el-radio-button label="month">月</el-radio-button>
                <el-radio-button label="year">年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <v-chart :option="trendChartOption" autoresize />
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>分类分布</span>
            </div>
          </template>
          <div class="chart-container">
            <v-chart :option="categoryChartOption" autoresize />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="chart-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>热门链接</span>
            </div>
          </template>
          <el-table :data="topLinks" border>
            <el-table-column prop="url" label="链接地址" min-width="300" />
            <el-table-column prop="category" label="分类" width="120" />
            <el-table-column prop="clicks" label="点击量" width="100" sortable />
            <el-table-column prop="last_clicked_at" label="最后点击时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.last_clicked_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import { externalApi } from '@/api/external'
import { formatDate } from '@/utils/date'

use([
  CanvasRenderer,
  LineChart,
  PieChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent
])

const stats = ref({
  total_links: 0,
  total_clicks: 0,
  active_links: 0,
  categories: {},
  trends: [],
  top_links: []
})

const trendPeriod = ref('week')
const topLinks = ref([])

// 获取统计数据
const fetchStats = async () => {
  try {
    const data = await externalApi.getExternalStatistics()
    stats.value = data
    topLinks.value = data.top_links
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 趋势图配置
const trendChartOption = computed(() => ({
  tooltip: {
    trigger: 'axis'
  },
  xAxis: {
    type: 'category',
    data: stats.value.trends.map(item => item.date)
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      name: '点击量',
      type: 'line',
      smooth: true,
      data: stats.value.trends.map(item => item.clicks)
    }
  ]
}))

// 分类分布图配置
const categoryChartOption = computed(() => ({
  tooltip: {
    trigger: 'item'
  },
  legend: {
    orient: 'vertical',
    left: 'left'
  },
  series: [
    {
      name: '分类分布',
      type: 'pie',
      radius: '50%',
      data: Object.entries(stats.value.categories).map(([name, value]) => ({
        name,
        value
      }))
    }
  ]
}))

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.external-link-stats {
  padding: 20px;
}

.stat-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  color: #409EFF;
}

.chart-row {
  margin-bottom: 20px;
}

.chart-container {
  height: 400px;
}
</style> 