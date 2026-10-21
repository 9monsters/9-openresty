<script setup lang="ts">
import { t } from '@/i18n'
import { getSelfCheck } from '@/api/system'

const loading = ref(true)
const checks = ref<any[]>([])

async function fetchData() {
  loading.value = true
  try {
    const res = await getSelfCheck()
    checks.value = res.checks
  } catch { /**/ }
  loading.value = false
}

const statusIcons: Record<string, string> = {
  ok: 'i-tabler:circle-check text-green-500',
  warning: 'i-tabler:alert-triangle text-orange-500',
  error: 'i-tabler:circle-x text-red-500',
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Self Check') }}</h3>
      <AButton size="small" @click="fetchData">{{ t('Refresh') }}</AButton>
    </div>
    <ASpin v-if="loading" class="flex-center w-full py-10" />
    <ACard v-else>
      <div class="space-y-2">
        <div v-for="check in checks" :key="check.name" class="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded">
          <span :class="statusIcons[check.status] || 'i-tabler:question-mark'" class="text-lg" />
          <div class="flex-1">
            <div class="font-semibold">{{ check.name }}</div>
            <div class="text-sm text-gray-500">{{ check.detail }}</div>
          </div>
          <ATag :color="check.status === 'ok' ? 'green' : check.status === 'warning' ? 'orange' : 'red'">{{ check.status }}</ATag>
        </div>
      </div>
    </ACard>
  </div>
</template>