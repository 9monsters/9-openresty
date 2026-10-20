<script setup lang="ts">
import { t } from '@/i18n'
import { getLogFiles, readLogByPath } from '@/api/logs'
import dayjs from 'dayjs'

const loading = ref(true)
const logFiles = ref<any[]>([])

async function fetch() {
  loading.value = true
  try {
    const r = await getLogFiles()
    logFiles.value = r.files || []
  } catch { /**/ }
  loading.value = false
}

function formatSize(bytes: number) { return `${(bytes / 1024).toFixed(1)} KB` }

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Log List') }}</h3>
      <AButton size="small" @click="fetch">{{ t('Refresh') }}</AButton>
    </div>
    <ASpin v-if="loading" class="flex-center w-full py-10" />
    <ACard v-else>
      <div class="space-y-2">
        <div v-for="f in logFiles" :key="f.name" class="flex-between p-3 bg-gray-50 dark:bg-gray-800 rounded">
          <div>
            <div class="font-semibold">{{ f.name }}</div>
            <div class="text-xs text-gray-500">{{ f.path }}</div>
          </div>
          <div class="text-right text-sm text-gray-500">
            <div>{{ formatSize(f.size) }}</div>
            <div>{{ dayjs(f.mtime).format('YYYY-MM-DD HH:mm') }}</div>
          </div>
        </div>
        <AEmpty v-if="logFiles.length === 0" :description="t('No Data')" />
      </div>
    </ACard>
  </div>
</template>