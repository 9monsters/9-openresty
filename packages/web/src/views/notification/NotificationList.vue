<script setup lang="ts">
import { listNotifications } from '@/api/notifications'
import dayjs from 'dayjs'
const notifications = ref<any[]>([])
onMounted(async () => { try { notifications.value = ((await listNotifications()) as any).notifications } catch {} })
const typeColor = (t: string) => ({ success: 'green', warning: 'orange', error: 'red', info: 'blue' } as any)[t] || 'default'
</script>
<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">通知</h3>
    <ATimeline>
      <ATimelineItem v-for="n in notifications" :key="n.id" :color="typeColor(n.type)">
        <template #dot><span :class="`i-tabler:${n.type === 'error' ? 'alert-triangle' : n.type === 'warning' ? 'alert-circle' : n.type === 'success' ? 'circle-check' : 'info-circle'}`" /></template>
        <div class="font-medium">{{ n.title }}</div>
        <div class="text-sm text-gray-400">{{ n.message }}</div>
        <div class="text-xs text-gray-400 mt-1">{{ dayjs(n.time).format('YYYY-MM-DD HH:mm') }}</div>
      </ATimelineItem>
    </ATimeline>
  </div>
</template>
