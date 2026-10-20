<script setup lang="ts">
import { theme } from 'antdv-next'
import zhCN from 'antdv-next/locale/zh_CN'
import enUS from 'antdv-next/locale/en_US'
import { useSettingsStore } from '@/stores/settings'

const settings = useSettingsStore()
const isDark = computed(() => settings.theme === 'dark')

const lang = computed(() => settings.language === 'zh-CN' ? zhCN : enUS)

const fontFamily = "-apple-system, BlinkMacSystemFont, 'SF Pro Text', 'Segoe UI', system-ui, sans-serif"

// Bridge the CSS palette into antdv's own tokens so component internals
// (popovers, dropdowns, tables) match the Liquid Glass surfaces.
const designTokens = computed(() => ({
  colorPrimary: isDark.value ? '#0a84ff' : '#007aff',
  colorInfo: isDark.value ? '#0a84ff' : '#007aff',
  colorSuccess: '#34c759',
  colorWarning: '#ff9f0a',
  colorError: '#ff3b30',
  colorBgBase: isDark.value ? '#000000' : '#f2f2f7',
  colorBgContainer: isDark.value ? '#1c1c1e' : '#ffffff',
  colorBgElevated: isDark.value ? '#1c1c1e' : '#ffffff',
  colorText: isDark.value ? 'rgba(255, 255, 255, 0.92)' : 'rgba(0, 0, 0, 0.9)',
  colorTextSecondary: isDark.value ? 'rgba(235, 235, 245, 0.6)' : 'rgba(60, 60, 67, 0.62)',
  colorBorder: isDark.value ? 'rgba(84, 84, 88, 0.5)' : 'rgba(60, 60, 67, 0.16)',
  colorBorderSecondary: isDark.value ? 'rgba(84, 84, 88, 0.32)' : 'rgba(60, 60, 67, 0.1)',
  borderRadius: 12,
  borderRadiusLG: 18,
  borderRadiusSM: 10,
  fontFamily,
  fontSize: 15,
  controlHeight: 34,
  wireframe: false,
}))

// Apply theme on mount + follow the system preference.
const media = window.matchMedia('(prefers-color-scheme: dark)')
function handleMediaChange() {
  settings.applyTheme()
}
media.addEventListener('change', handleMediaChange)
onUnmounted(() => media.removeEventListener('change', handleMediaChange))

settings.applyTheme()
</script>

<template>
  <AConfigProvider
    :theme="{
      algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
      token: designTokens,
    }"
    :locale="lang"
    :button="{ autoInsertSpace: false }"
  >
    <AApp>
      <RouterView />
    </AApp>
  </AConfigProvider>
</template>
