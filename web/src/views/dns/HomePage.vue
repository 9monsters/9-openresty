<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'

const entries = [
  { title: 'Credentials', description: t('Manage DNS provider API credentials — Cloudflare, Aliyun DNS, and other supported providers.'), path: '/dns/credentials', icon: 'i-tabler:key' },
  { title: 'DNS Domains', description: t('Manage DNS domain zones, view record counts, and manage DNS records per domain.'), path: '/dns/domains', icon: 'i-tabler:world' },
  { title: 'Groups', description: t('Configure DNS groups with routing strategies (round-robin, weighted, geographic).'), path: '/dns/groups', icon: 'i-tabler:folders' },
  { title: 'DDNS', description: t('Dynamic DNS entries that automatically update IP addresses for home and dynamic endpoints.'), path: '/dns/ddns', icon: 'i-tabler:refresh' },
]
</script>

<template>
  <div>
    <h2 class="ios-title mb-6">{{ t('DNS') }}</h2>
    <EntryGroup :entries="entries" />
  </div>
</template>
