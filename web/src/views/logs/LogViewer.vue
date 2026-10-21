<script setup lang="ts">
import { getAccessLog, getErrorLog } from '@/api/logs'

const activeTab = ref<'access' | 'error'>('access')
const logContent = ref('')
const loading = ref(false)
const lines = ref(200)
const autoRefresh = ref(false)
let refreshTimer: ReturnType<typeof setInterval> | null = null

async function fetchLogs() {
  loading.value = true
  try {
    const res = activeTab.value === 'access' ? await getAccessLog(lines.value) : await getErrorLog(lines.value)
    logContent.value = res.content
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function toggleAutoRefresh() {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    refreshTimer = setInterval(fetchLogs, 5000)
  } else if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

function highlightLine(line: string): boolean {
  return line.includes('error') || line.includes('Error') || line.includes('ERROR')
    || (activeTab.value === 'error' && (line.includes('warn') || line.includes('crit')))
}

watch(activeTab, fetchLogs)
onMounted(fetchLogs)
onUnmounted(() => { if (refreshTimer) clearInterval(refreshTimer) })
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">日志查看</h3>
      <ASpace>
        <ASelect v-model:value="lines" style="width:100px" @change="fetchLogs">
          <ASelectOption :value="100">100 行</ASelectOption>
          <ASelectOption :value="200">200 行</ASelectOption>
          <ASelectOption :value="500">500 行</ASelectOption>
          <ASelectOption :value="1000">1000 行</ASelectOption>
        </ASelect>
        <AButton @click="toggleAutoRefresh" :type="autoRefresh ? 'primary' : 'default'">
          <span :class="autoRefresh ? 'i-tabler:player-pause' : 'i-tabler:player-play'" /> 
          {{ autoRefresh ? '停止刷新' : '自动刷新' }}
        </AButton>
        <AButton @click="fetchLogs">
          <span class="i-tabler:refresh" />
        </AButton>
      </ASpace>
    </div>

    <ATabs v-model:activeKey="activeTab">
      <ATabPane key="access" tab="访问日志 (Access Log)" />
      <ATabPane key="error" tab="错误日志 (Error Log)" />
    </ATabs>

    <div
      v-if="loading && !logContent"
      class="flex-center py-20"
    >
      <ASpin size="large" />
    </div>

    <div
      v-else
      class="bg-gray-900 text-green-400 font-mono text-xs p-4 rounded-lg overflow-auto max-h-[calc(100vh-16rem)]"
      style="white-space: pre-wrap; word-break: break-all;"
    >
      <template v-if="logContent">
        <div v-for="(line, i) in logContent.split('\n')" :key="i"
          :class="highlightLine(line) ? 'text-red-400' : ''"
        >
          {{ line }}
        </div>
      </template>
      <div v-else class="text-gray-500">暂无日志内容</div>
    </div>
  </div>
</template>
