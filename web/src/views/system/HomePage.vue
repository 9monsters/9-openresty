<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'

const entries = [
  { title: 'Self Check', description: t('Run diagnostic checks on the Nginx binary, configuration syntax, permissions, and port availability.'), path: '/system/self_check', icon: 'i-tabler:stethoscope' },
  { title: 'Upgrade', description: t('Check for new versions and upgrade OpenResty Admin to the latest release.'), path: '/system/upgrade', icon: 'i-tabler:arrow-up-circle' },
  { title: 'About', description: t('View version information for OpenResty Admin, Nginx, Node.js, and the host platform.'), path: '/system/about', icon: 'i-tabler:info-circle' },
]
</script>

<template>
  <div>
    <h2 class="ios-title mb-6">{{ t('System') }}</h2>
    <EntryGroup :entries="entries" />
  </div>
</template>
