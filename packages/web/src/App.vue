<script setup lang="ts">
import { theme } from 'antdv-next'
import zhCN from 'antdv-next/locale/zh_CN'
import { useSettingsStore } from './stores/settings'

const settings = useSettingsStore()
const isDark = computed(() => settings.theme === 'dark')
</script>

<template>
  <AConfigProvider
    :theme="{ algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm }"
    :locale="zhCN"
    :button="{ autoInsertSpace: false }"
  >
    <AApp>
      <RouterView />
    </AApp>
  </AConfigProvider>
</template>

<style>
body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }

::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #c1c1c1; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #a1a1a1; }

.dark ::-webkit-scrollbar-thumb { background: #4a4a4a; }
.dark ::-webkit-scrollbar-thumb:hover { background: #5a5a5a; }
</style>
