<script setup lang="ts">
import { getSystemHistory, type SystemHistoryPoint } from '@/api/stats'
import { t } from '@/i18n'
import VChart from 'vue-echarts'

const loading = ref(true)
const history = ref<SystemHistoryPoint[]>([])
const interval = ref(5)
const intervals = [1, 5, 10] as const
let timer: ReturnType<typeof setInterval> | null = null

async function fetchData() {
  try {
    const r = await getSystemHistory()
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

function formatBytes(bytes: number): string {
  return (bytes / (1024 ** 3)).toFixed(1) + ' GB'
}

const cpuChartOption = computed(() => ({
  backgroundColor: 'transparent',
  grid: { top: 8, right: 12, bottom: 20, left: 40 },
  xAxis: { type: 'category', data: history.value.map(p => formatTime(p.time)), axisLabel: { color: 'var(--label-2)', fontSize: 10 }, axisLine: { lineStyle: { color: 'var(--separator)' } } },
  yAxis: [
    { type: 'value', name: 'Load', min: 0, axisLabel: { color: 'var(--label-2)' }, splitLine: { lineStyle: { color: 'var(--separator)' } } },
  ],
  tooltip: { trigger: 'axis' },
  series: [
    { name: 'Load 1m', type: 'line', data: history.value.map(p => +p.loadAvg1.toFixed(2)), smooth: true, symbol: 'none', lineStyle: { color: 'var(--tint)', width: 1.5 } },
    { name: 'Load 5m', type: 'line', data: history.value.map(p => +p.loadAvg5.toFixed(2)), smooth: true, symbol: 'none', lineStyle: { color: '#ff9500', width: 1.5 } },
    { name: 'Load 15m', type: 'line', data: history.value.map(p => +p.loadAvg15.toFixed(2)), smooth: true, symbol: 'none', lineStyle: { color: 'var(--label-2)', width: 1, type: 'dashed' } },
  ],
}))

const memChartOption = computed(() => ({
  backgroundColor: 'transparent',
  grid: { top: 8, right: 12, bottom: 20, left: 45 },
  xAxis: { type: 'category', data: history.value.map(p => formatTime(p.time)), axisLabel: { color: 'var(--label-2)', fontSize: 10 }, axisLine: { lineStyle: { color: 'var(--separator)' } } },
  yAxis: { type: 'value', axisLabel: { color: 'var(--label-2)', formatter: (v: number) => formatBytes(v) }, splitLine: { lineStyle: { color: 'var(--separator)' } } },
  tooltip: { trigger: 'axis', valueFormatter: (v: number) => formatBytes(v) },
  series: [{
    name: t('Memory Used', '已用内存'), type: 'line', data: history.value.map(p => p.memoryUsed), smooth: true,
    symbol: 'none', areaStyle: { color: 'rgba(10,132,255,0.12)' }, lineStyle: { color: 'var(--tint)', width: 1.5 },
  }],
}))

const latest = computed(() => history.value.length > 0 ? history.value[history.value.length - 1] : null)

function percent(part: number, total: number): number {
  if (!total) return 0
  return Math.round((part / total) * 100)
}

onMounted(() => startPolling())
onUnmounted(() => stopPolling())
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Server Status') }}</h3>
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
      <!-- Gauges row -->
      <div class="ios-group panel gauges" v-if="latest">
        <div class="gauge-item">
          <div class="ios-caption">{{ t('CPU Load', 'CPU 负载') }}</div>
          <div class="ios-instrument gauge-value">{{ latest.loadAvg1.toFixed(2) }}</div>
          <div class="ios-footnote">{{ latest.cpuCount }} cores</div>
        </div>
        <div class="gauge-item">
          <div class="ios-caption">{{ t('Memory', '内存') }}</div>
          <div class="ios-instrument gauge-value">{{ percent(latest.memoryUsed, latest.memoryTotal) }}%</div>
          <div class="ios-footnote">{{ formatBytes(latest.memoryUsed) }} / {{ formatBytes(latest.memoryTotal) }}</div>
        </div>
        <div class="gauge-item">
          <div class="ios-caption">{{ t('Uptime', '运行时间') }}</div>
          <div class="ios-instrument gauge-value">{{ Math.floor(latest.uptimeSeconds / 3600) }}h</div>
          <div class="ios-footnote">{{ Math.floor(latest.uptimeSeconds / 86400) }}d</div>
        </div>
      </div>

      <!-- CPU Load chart -->
      <div class="ios-group chart-card">
        <p class="ios-headline chart-title">{{ t('CPU Load Average', 'CPU 负载') }}</p>
        <VChart :option="cpuChartOption" autoresize style="height: 240px" />
      </div>

      <!-- Memory chart -->
      <div class="ios-group chart-card">
        <p class="ios-headline chart-title">{{ t('Memory Usage', '内存使用') }}</p>
        <VChart :option="memChartOption" autoresize style="height: 200px" />
      </div>

      <p v-if="history.length === 0" class="ios-footnote empty">
        {{ t('Collecting data...', '数据采集中...') }}
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
.empty { margin: 14px 4px 0; }
</style>