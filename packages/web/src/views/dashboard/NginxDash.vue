<script setup lang="ts">
import { getNginxHistory, type NginxHistoryPoint } from '@/api/stats'
import { t } from '@/i18n'
import VChart from 'vue-echarts'

const loading = ref(true)
const history = ref<NginxHistoryPoint[]>([])
const interval = ref(5)
const intervals = [1, 5, 10] as const
let timer: ReturnType<typeof setInterval> | null = null

async function fetchData() {
  try {
    const r = await getNginxHistory()
    history.value = r.history || []
  } catch { /* ignore */ }
  loading.value = false
}

function startPolling() {
  stopPolling()
  fetchData()
  timer = setInterval(fetchData, interval.value * 1000)
}
function stopPolling() {
  if (timer) { clearInterval(timer); timer = null }
}
function setIntervalAndRestart(s: number) {
  interval.value = s
  startPolling()
}

function formatTime(iso: string): string {
  const d = new Date(iso)
  return d.toLocaleTimeString('en-US', { hour12: false })
}
function num(v?: number): string {
  return typeof v === 'number' ? v.toLocaleString('en-US') : '—'
}

const requestsChartOption = computed(() => ({
  backgroundColor: 'transparent',
  grid: { top: 8, right: 12, bottom: 20, left: 50 },
  xAxis: { type: 'category', data: history.value.map(p => formatTime(p.time)), axisLabel: { color: 'var(--label-2)', fontSize: 10 }, axisLine: { lineStyle: { color: 'var(--separator)' } } },
  yAxis: { type: 'value', axisLabel: { color: 'var(--label-2)' }, splitLine: { lineStyle: { color: 'var(--separator)' } } },
  tooltip: { trigger: 'axis' },
  series: [
    { name: t('Accepts', '总接受连接'), type: 'line', data: history.value.map(p => p.accepts), smooth: true, symbol: 'none', lineStyle: { color: 'var(--tint)', width: 1.5 } },
    { name: t('Handled', '已处理连接'), type: 'line', data: history.value.map(p => p.handled), smooth: true, symbol: 'none', lineStyle: { color: '#34c759', width: 1.5 } },
    { name: t('Requests', '总请求数'), type: 'line', data: history.value.map(p => p.requests), smooth: true, symbol: 'none', lineStyle: { color: '#ff9500', width: 1.5 } },
  ],
}))

const connectionsChartOption = computed(() => ({
  backgroundColor: 'transparent',
  grid: { top: 8, right: 12, bottom: 20, left: 40 },
  xAxis: { type: 'category', data: history.value.map(p => formatTime(p.time)), axisLabel: { color: 'var(--label-2)', fontSize: 10 }, axisLine: { lineStyle: { color: 'var(--separator)' } } },
  yAxis: { type: 'value', min: 0, axisLabel: { color: 'var(--label-2)' }, splitLine: { lineStyle: { color: 'var(--separator)' } } },
  tooltip: { trigger: 'axis' },
  series: [
    { name: t('Active Connections', '活跃连接数'), type: 'line', data: history.value.map(p => p.activeConnections), smooth: true, symbol: 'none', areaStyle: { color: 'rgba(10,132,255,0.12)' }, lineStyle: { color: 'var(--tint)', width: 1.5 } },
  ],
}))

const ioChartOption = computed(() => ({
  backgroundColor: 'transparent',
  grid: { top: 8, right: 12, bottom: 20, left: 40 },
  xAxis: { type: 'category', data: history.value.map(p => formatTime(p.time)), axisLabel: { color: 'var(--label-2)', fontSize: 10 }, axisLine: { lineStyle: { color: 'var(--separator)' } } },
  yAxis: { type: 'value', min: 0, axisLabel: { color: 'var(--label-2)' }, splitLine: { lineStyle: { color: 'var(--separator)' } } },
  tooltip: { trigger: 'axis' },
  legend: { textStyle: { color: 'var(--label-2)', fontSize: 11 } },
  series: [
    { name: t('Reading', '读取中'), type: 'line', data: history.value.map(p => p.reading), smooth: true, symbol: 'none', lineStyle: { color: '#34c759', width: 1.5 } },
    { name: t('Writing', '写入中'), type: 'line', data: history.value.map(p => p.writing), smooth: true, symbol: 'none', lineStyle: { color: 'var(--tint)', width: 1.5 } },
    { name: t('Waiting', '等待中'), type: 'line', data: history.value.map(p => p.waiting), smooth: true, symbol: 'none', lineStyle: { color: '#ff9500', width: 1.5 } },
  ],
}))

const latest = computed(() => history.value.length > 0 ? history.value[history.value.length - 1] : null)

onMounted(() => startPolling())
onUnmounted(() => stopPolling())
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Nginx Performance') }}</h3>
      <div class="flex gap-2">
        <button
          v-for="s in intervals" :key="s"
          class="interval-btn" :class="{ active: interval === s }"
          @click="setIntervalAndRestart(s)"
        >{{ s }}s</button>
      </div>
    </div>

    <ASpin v-if="loading" class="flex-center w-full py-20" />

    <template v-else>
      <!-- KPI strip -->
      <div class="ios-group panel gauges" v-if="latest">
        <div class="gauge-item">
          <div class="ios-caption">{{ t('Active Connections', '活跃连接数') }}</div>
          <div class="ios-instrument gauge-value">{{ num(latest.activeConnections) }}</div>
        </div>
        <div class="gauge-item">
          <div class="ios-caption">{{ t('Requests', '总请求数') }}</div>
          <div class="ios-instrument gauge-value">{{ num(latest.requests) }}</div>
        </div>
        <div class="gauge-item">
          <div class="ios-caption">{{ t('Accepts', '总接受连接') }}</div>
          <div class="ios-instrument gauge-value">{{ num(latest.accepts) }}</div>
        </div>
        <div class="gauge-item">
          <div class="ios-caption">{{ t('Handled', '已处理连接') }}</div>
          <div class="ios-instrument gauge-value">{{ num(latest.handled) }}</div>
        </div>
      </div>

      <!-- Requests chart -->
      <div class="ios-group chart-card">
        <p class="ios-headline chart-title">{{ t('Request Statistics', '请求统计') }}</p>
        <VChart :option="requestsChartOption" autoresize style="height: 240px" />
      </div>

      <!-- Connections chart -->
      <div class="ios-group chart-card">
        <p class="ios-headline chart-title">{{ t('Active Connections', '活跃连接数') }}</p>
        <VChart :option="connectionsChartOption" autoresize style="height: 200px" />
      </div>

      <!-- I/O chart -->
      <div class="ios-group chart-card">
        <p class="ios-headline chart-title">{{ t('Connection Details', '连接状态详情') }}</p>
        <VChart :option="ioChartOption" autoresize style="height: 200px" />
      </div>

      <p v-if="history.length === 0" class="ios-footnote empty">
        {{ t('No Nginx stats available. Enable stub_status for live data.', '无 Nginx 统计数据，启用 stub_status 后可见实时数据。') }}
      </p>
    </template>
  </div>
</template>

<style scoped>
.interval-btn {
  padding: 4px 12px;
  border: 0.5px solid var(--separator);
  border-radius: var(--r-pill);
  background: var(--fill);
  color: var(--label-2);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--dur) var(--ease-ios);
}
.interval-btn.active {
  background: var(--tint);
  color: #fff;
  border-color: var(--tint);
}
.gauges {
  display: flex;
  gap: 0;
  padding: 12px 0;
  margin-bottom: 14px;
}
.gauge-item {
  flex: 1;
  text-align: center;
  padding: 8px 12px;
}
.gauge-item + .gauge-item {
  border-left: 0.5px solid var(--separator);
}
.gauge-value {
  font-size: 28px;
  margin: 4px 0;
}
.chart-card {
  padding: 14px 16px 8px;
  margin-bottom: 14px;
}
.chart-title {
  margin: 0 0 4px 4px;
  color: var(--label);
}
.empty { margin: 14px 4px 0; color: var(--label-2); }
</style>