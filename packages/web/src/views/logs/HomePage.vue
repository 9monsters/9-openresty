<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'

const entries = [
  { title: 'Access Logs', description: t('View real-time Nginx access logs with configurable line count and auto-refresh.'), path: '/nginx_log/access', icon: 'i-tabler:file-text' },
  { title: 'Error Logs', description: t('Monitor Nginx error logs with highlighted error levels and live refresh.'), path: '/nginx_log/error', icon: 'i-tabler:file-alert' },
  { title: 'Log List', description: t('Browse all available log files on the filesystem with sizes and modification times.'), path: '/nginx_log/list', icon: 'i-tabler:files' },
]
</script>

<template>
  <div>
    <h2 class="ios-title mb-6">{{ t('Nginx Log') }}</h2>
    <EntryGroup :entries="entries" />
  </div>
</template>
