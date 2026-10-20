<script setup lang="ts">
import { getDashboard } from '@/api/dashboard'
import { getNginxStatus } from '@/api/nginx'
import { getSystemInfo } from '@/api/system'
import { t } from '@/i18n'
import dayjs from 'dayjs'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart, GaugeChart, LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, TitleComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([BarChart, GaugeChart, LineChart, GridComponent, TooltipComponent, TitleComponent, LegendComponent, CanvasRenderer])

interface DashData {
  status: string; version: string; isOpenResty: boolean
  stats: { activeConnections: number; accepts: number; handled: number; requests: number; reading: number; writing: number; waiting: number } | null
  serverTime: string
}

const loading = ref(true)
const dashboard = ref<DashData | null>(null)
const nginxInfo = ref<any>(null)
const systemInfo = ref<any>(null)

onMounted(async () => {
  try {
    const [dash, nginx, sys] = await Promise.all([
      getDashboard(), getNginxStatus(), getSystemInfo().catch(() => null),
    ])
    dashboard.value = dash as any
    nginxInfo.value = nginx
    systemInfo.value = sys
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})

// Gauge chart for active connections
const connectionsGauge = computed(() => ({
  series: [{
    type: 'gauge',
    startAngle: 210, endAngle: -30,
    center: ['50%', '60%'], radius: '90%',
    min: 0, max: 200,
    axisLine: { lineStyle: { width: 12, color: [[0.3, '#52c41a'], [0.7, '#faad14'], [1, '#ff4d4f']] } },
    pointer: { length: '60%', width: 6, itemStyle: { color: 'auto' } },
    axisTick: { distance: -12, length: 6, lineStyle: { width: 1, color: '#999' } },
    splitLine: { distance: -22, length: 16, lineStyle: { width: 2, color: '#999' } },
    axisLabel: { distance: -40, fontSize: 10, color: '#999' },
    detail: { valueAnimation: true, fontSize: 28, fontWeight: 'bold', offsetCenter: [0, '70%'], formatter: '{value}' },
    data: [{ value: dashboard.value?.stats?.activeConnections ?? 0, name: t('Active Connections') }],
  }],
}))

// Bar chart: accepts / handled / requests
const requestBar = computed(() => ({
  tooltip: { trigger: 'axis' },
  legend: { bottom: 0, textStyle: { fontSize: 11 } },
  grid: { left: 50, right: 20, top: 20, bottom: 30 },
  xAxis: { type: 'category', data: [t('Accepts') || 'Accepts', t('Handled') || 'Handled', t('Requests') || 'Requests'] },
  yAxis: { type: 'value' },
  series: [{
    type: 'bar', name: t('Count') || 'Count',
    data: [
      { value: dashboard.value?.stats?.accepts ?? 0, itemStyle: { color: '#1890ff' } },
      { value: dashboard.value?.stats?.handled ?? 0, itemStyle: { color: '#722ed1' } },
      { value: dashboard.value?.stats?.requests ?? 0, itemStyle: { color: '#fa8c16' } },
    ],
    barWidth: '40%', borderRadius: 4,
  }],
}))

// Pie: Reading / Writing / Waiting
const rwPie = computed(() => ({
  tooltip: { trigger: 'item' },
  legend: { bottom: 0, textStyle: { fontSize: 11 } },
  series: [{
    type: 'pie',
    radius: ['45%', '75%'],
    center: ['50%', '45%'],
    label: { show: false },
    emphasis: { label: { show: true, fontWeight: 'bold' } },
    data: [
      { value: dashboard.value?.stats?.reading ?? 0, name: t('Reading'), itemStyle: { color: '#1890ff' } },
      { value: dashboard.value?.stats?.writing ?? 0, name: t('Writing'), itemStyle: { color: '#52c41a' } },
      { value: dashboard.value?.stats?.waiting ?? 0, name: t('Waiting'), itemStyle: { color: '#fa8c16' } },
    ],
  }],
}))

// System info
const memPercent = computed(() => systemInfo.value?.memory?.usagePercent ?? 0)
</script>

<template>
  <div v-if="loading" class="flex-center py-20"><ASpin size="large" /></div>

  <div v-else>
    <!-- Status header -->
    <div class="flex-between mb-6">
      <div>
        <h2 class="text-2xl font-bold m-0">
          {{ nginxInfo?.isOpenResty ? 'OpenResty' : 'Nginx' }}
          <ATag :color="nginxInfo?.status === 'running' ? 'green' : 'red'" class="ml-2 align-middle">
            {{ nginxInfo?.status === 'running' ? t('Running') : t('Stopped') }}
          </ATag>
        </h2>
        <p class="text-gray-500 mt-1">{{ nginxInfo?.version }}</p>
      </div>
      <div class="text-right text-gray-400 text-sm">
        {{ t('Server Time') }}: {{ dayjs(dashboard?.serverTime).format('YYYY-MM-DD HH:mm:ss') }}
      </div>
    </div>

    <!-- Charts row -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-6">
      <!-- Gauge: Active Connections -->
      <ACard :title="t('Active Connections')" size="small" class="col-span-1">
        <VChart :option="connectionsGauge" autoresize style="height:260px" />
      </ACard>

      <!-- Bar: Accepts / Handled / Requests -->
      <ACard :title="t('Request Statistics') || 'Request Statistics'" size="small" class="col-span-1">
        <VChart :option="requestBar" autoresize style="height:260px" />
      </ACard>

      <!-- Pie: Reading / Writing / Waiting -->
      <ACard :title="t('Connection Details')" size="small" class="col-span-1">
        <VChart :option="rwPie" autoresize style="height:260px" />
      </ACard>
    </div>

    <!-- System info + Memory bar -->
    <div v-if="systemInfo" class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
      <ACard :title="t('Server')" size="small">
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div><span class="text-gray-500">{{ t('Hostname') }}:</span> {{ systemInfo.hostname }}</div>
          <div><span class="text-gray-500">{{ t('Platform') }}:</span> {{ systemInfo.platform }} / {{ systemInfo.arch }}</div>
          <div><span class="text-gray-500">CPU:</span> {{ systemInfo.cpu?.model?.split('@')[0]?.trim() }}</div>
          <div><span class="text-gray-500">Node.js:</span> {{ systemInfo.nodejs }}</div>
        </div>
      </ACard>
      <ACard :title="t('Memory')" size="small">
        <div class="text-center text-2xl font-bold mb-2">{{ memPercent }}%</div>
        <AProgress :percent="memPercent" :strokeColor="memPercent > 80 ? '#ff4d4f' : '#1890ff'" />
        <div class="text-center text-sm text-gray-500 mt-2">
          {{ ((systemInfo.memory?.used ?? 0) / (1024**3)).toFixed(1) }} GB / {{ ((systemInfo.memory?.total ?? 0) / (1024**3)).toFixed(1) }} GB
        </div>
      </ACard>
    </div>

    <!-- Quick actions -->
    <ACard :title="t('Quick Actions')">
      <ASpace>
        <AButton type="primary" @click="$router.push('/sites/list')">
          <span class="i-tabler:world-www mr-1" /> {{ t('Manage Sites') }}
        </AButton>
        <AButton @click="$router.push('/config')">
          <span class="i-tabler:file-code mr-1" /> {{ t('Manage Configs') }}
        </AButton>
        <AButton @click="$router.push('/nginx_log/access')">
          <span class="i-tabler:file-text mr-1" /> {{ t('Nginx Log') }}
        </AButton>
        <AButton @click="$router.push('/lua')">
          <span class="i-tabler:code mr-1" /> {{ t('Lua Manager') }}
        </AButton>
      </ASpace>
    </ACard>
  </div>
</template>