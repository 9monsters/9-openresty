<script setup lang="ts">
import { listNotifications, markRead, markAllRead, deleteNotification } from '@/api/notifications'
import { t } from '@/i18n'
import dayjs from 'dayjs'

const loading = ref(true)
const notifications = ref<any[]>([])
const filterType = ref('all')

const typeColors: Record<string, string> = {
  success: 'green', warning: 'orange', error: 'red', info: 'blue',
}

async function fetchData() {
  loading.value = true
  try {
    const res = await listNotifications()
    const all = res.notifications || []
    notifications.value = filterType.value === 'all' ? all : all.filter(n => n.type === filterType.value)
  } catch { /**/ }
  loading.value = false
}

function changeFilter(type: string) {
  filterType.value = type
  fetchData()
}

async function handleMarkRead(id: number) {
  try { await markRead(id); fetchData() } catch { /**/ }
}

async function handleMarkAllRead() {
  try { await markAllRead(); fetchData() } catch { /**/ }
}

async function handleDelete(id: number) {
  try { await deleteNotification(id); fetchData() } catch { /**/ }
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Notifications') }}</h3>
      <div class="flex gap-2">
        <AButton size="small" @click="handleMarkAllRead">{{ t('Mark All Read') || 'Mark All Read' }}</AButton>
      </div>
    </div>

    <ARadioGroup :value="filterType" @change="(e: any) => changeFilter(e.target.value)" class="mb-4">
      <ARadioButton value="all">{{ t('All') || 'All' }}</ARadioButton>
      <ARadioButton value="success">{{ t('Success') }}</ARadioButton>
      <ARadioButton value="warning">{{ t('Warning') }}</ARadioButton>
      <ARadioButton value="error">{{ t('Error') }}</ARadioButton>
      <ARadioButton value="info">{{ t('Info') }}</ARadioButton>
    </ARadioGroup>

    <ASpin v-if="loading" class="flex-center w-full py-10" />

    <ATimeline v-else>
      <ATimelineItem v-for="n in notifications" :key="n.id" :color="typeColors[n.type] || 'blue'">
        <div class="flex-between">
          <div>
            <div class="font-semibold" :class="{ 'text-gray-400': n.read }">
              <ATag :color="typeColors[n.type]" class="mr-2" style="vertical-align: middle">{{ n.type }}</ATag>
              {{ n.title }}
              <span v-if="!n.read" class="i-tabler:circle-filled text-blue-500 text-xs ml-1 inline-block" style="vertical-align: middle" />
            </div>
            <div class="text-sm text-gray-500 mt-1">{{ n.message }}</div>
            <div class="text-xs text-gray-400 mt-1">{{ formatTime(n.time) }}</div>
          </div>
          <div class="flex gap-2">
            <AButton v-if="!n.read" type="link" size="small" @click="handleMarkRead(n.id)">{{ t('Mark Read') || 'Read' }}</AButton>
            <AButton type="link" size="small" danger @click="handleDelete(n.id)">{{ t('Delete') }}</AButton>
          </div>
        </div>
      </ATimelineItem>
    </ATimeline>

    <AEmpty v-if="!loading && notifications.length === 0" :description="t('No Data')" />
  </div>
</template>