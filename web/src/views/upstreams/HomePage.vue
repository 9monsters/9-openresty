<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'

const entries = [
  { title: 'Upstream Groups', description: t('View and manage upstream groups with their backend servers, weights, and health status.'), path: '/upstream/list', icon: 'i-tabler:server' },
]
</script>

<template>
  <div>
    <h2 class="ios-title mb-6">{{ t('Upstream') }}</h2>
    <EntryGroup :entries="entries" />
  </div>
</template>
