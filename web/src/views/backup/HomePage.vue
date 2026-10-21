<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'

const entries = [
  { title: 'Backup', description: t('Manual backup and restore operations for Nginx configuration files with version history.'), path: '/backup/backup-and-restore', icon: 'i-tabler:database' },
  { title: 'Auto Backup', description: t('Configure scheduled automatic backups with configurable interval, retention, and target path.'), path: '/backup/auto-backup', icon: 'i-tabler:clock' },
]
</script>

<template>
  <div>
    <h2 class="ios-title mb-6">{{ t('Backup Manage') }}</h2>
    <EntryGroup :entries="entries" />
  </div>
</template>
