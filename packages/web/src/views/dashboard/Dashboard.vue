<script setup lang="ts">
import { getDashboard, type DashboardData } from '@/api/dashboard'
import { getNginxStatus } from '@/api/nginx'
import dayjs from 'dayjs'

const dashboard = ref<DashboardData | null>(null)
const loading = ref(true)
const nginxInfo = ref<{ status: string; version: string; isOpenResty: boolean } | null>(null)

onMounted(async () => {
  try {
    const [dash, nginx] = await Promise.all([getDashboard(), getNginxStatus()])
    dashboard.value = dash
    nginxInfo.value = nginx
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})

const statCards = computed(() => [
  { title: '活跃连接数', value: dashboard.value?.stats?.activeConnections ?? '-', unit: '', color: 'text-blue-600', icon: 'i-tabler:plug-connected' },
  { title: '总接受连接', value: dashboard.value?.stats?.accepts ?? '-', unit: '', color: 'text-green-600', icon: 'i-tabler:arrow-big-down-lines' },
  { title: '已处理连接', value: dashboard.value?.stats?.handled ?? '-', unit: '', color: 'text-purple-600', icon: 'i-tabler:checks' },
  { title: '总请求数', value: dashboard.value?.stats?.requests ?? '-', unit: '', color: 'text-orange-600', icon: 'i-tabler:http-get' },
])
</script>

<template>
  <div v-if="loading" class="flex-center py-20">
    <ASpin size="large" />
  </div>

  <div v-else>
    <!-- Status header -->
    <div class="flex-between mb-6">
      <div>
        <h2 class="text-2xl font-bold m-0">
          {{ nginxInfo?.isOpenResty ? '🚀 OpenResty' : 'Nginx' }}
          <ATag :color="nginxInfo?.status === 'running' ? 'green' : 'red'" class="ml-2">
            {{ nginxInfo?.status === 'running' ? '运行中' : '已停止' }}
          </ATag>
        </h2>
        <p class="text-gray-500 mt-1">{{ nginxInfo?.version }}</p>
      </div>
      <div class="text-right text-gray-400 text-sm">
        服务器时间: {{ dayjs(dashboard?.serverTime).format('YYYY-MM-DD HH:mm:ss') }}
      </div>
    </div>

    <!-- Stat cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <ACard v-for="card in statCards" :key="card.title" class="card-hover">
        <div class="flex-between">
          <div>
            <div class="text-gray-500 text-sm">{{ card.title }}</div>
            <div :class="[card.color, 'text-2xl font-bold mt-1']">
              {{ card.value }}{{ card.unit }}
            </div>
          </div>
          <span :class="[card.icon, card.color, 'text-4xl opacity-30']" />
        </div>
      </ACard>
    </div>

    <!-- Connection detail -->
    <ACard title="连接状态详情" class="mb-6">
      <div class="grid grid-cols-2 lg:grid-cols-3 gap-6">
        <div class="text-center p-4">
          <div class="text-3xl font-bold text-blue-600">{{ dashboard?.stats?.reading ?? '-' }}</div>
          <div class="text-gray-500 mt-1">Reading</div>
        </div>
        <div class="text-center p-4">
          <div class="text-3xl font-bold text-green-600">{{ dashboard?.stats?.writing ?? '-' }}</div>
          <div class="text-gray-500 mt-1">Writing</div>
        </div>
        <div class="text-center p-4">
          <div class="text-3xl font-bold text-orange-600">{{ dashboard?.stats?.waiting ?? '-' }}</div>
          <div class="text-gray-500 mt-1">Waiting</div>
        </div>
      </div>
    </ACard>

    <!-- Quick actions -->
    <ACard title="快速操作">
      <ASpace>
        <AButton type="primary" @click="$router.push('/sites')">
          <span class="i-tabler:world-www mr-1" /> 管理站点
        </AButton>
        <AButton @click="$router.push('/config')">
          <span class="i-tabler:file-code mr-1" /> 配置文件
        </AButton>
        <AButton @click="$router.push('/logs')">
          <span class="i-tabler:file-text mr-1" /> 查看日志
        </AButton>
        <AButton @click="$router.push('/lua')">
          <span class="i-tabler:code mr-1" /> Lua 代码
        </AButton>
      </ASpace>
    </ACard>
  </div>
</template>
