<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'

const entries = [
  { title: 'Certificates List', description: t('Browse, import, and manage SSL/TLS certificates. View expiry dates, domains, and auto-renew status.'), path: '/certificates/list', icon: 'i-tabler:certificate' },
  { title: 'ACME User', description: t('Manage ACME (Automatic Certificate Management Environment) accounts for automatic certificate issuance and renewal.'), path: '/certificates/acme_users', icon: 'i-tabler:user-shield' },
]
</script>

<template>
  <div>
    <h2 class="ios-title mb-6">{{ t('Certificates') }}</h2>
    <EntryGroup :entries="entries" />
  </div>
</template>
