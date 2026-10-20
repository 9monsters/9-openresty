<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'

const entries = [
  { title: 'Sites List', description: t('Browse, search, edit, and delete your server block configurations. Enable or disable sites per file.'), path: '/sites/list', icon: 'i-tabler:list' },
  { title: t('Add Site'), description: t('Create a new Nginx server block with listen ports, server name, root directory, locations, and Lua directives.'), path: '/sites/add', icon: 'i-tabler:plus' },
  { title: 'Snippets', description: t('Reusable configuration snippets that can be included in multiple sites. Manage custom and built-in templates.'), path: '/snippets', icon: 'i-tabler:code-dots' },
]
</script>

<template>
  <div>
    <h2 class="ios-title mb-6">{{ t('Manage Sites') }}</h2>
    <EntryGroup :entries="entries" />
  </div>
</template>
