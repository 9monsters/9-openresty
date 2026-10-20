import { defineStore } from 'pinia'
import { currentLocale, setLocale } from '@/i18n'
import type { Locale } from '@/i18n'

type PreferenceTheme = 'auto' | 'light' | 'dark'

function persist(key: string, val: any) {
  try { localStorage.setItem(`openresty-${key}`, JSON.stringify(val)) } catch { /* ignore */ }
}
function load(key: string, fallback: any) {
  try { const v = localStorage.getItem(`openresty-${key}`); return v ? JSON.parse(v) : fallback } catch { return fallback }
}

export const useSettingsStore = defineStore('settings', () => {
  const theme = ref<'light' | 'dark'>(
    load('theme', 'light')
  )
  const preferenceTheme = ref<PreferenceTheme>(
    load('preference-setting', 'auto')
  )
  const sidebarCollapsed = ref(load('sidebar-collapsed', false))
  const serverName = ref('OpenResty Admin')
  const nginxStatus = ref<'running' | 'stopped' | 'unknown'>('unknown')
  const language = ref<Locale>(load('language', 'zh-CN'))

  function setTheme(t: 'light' | 'dark') {
    theme.value = t
    document.documentElement.classList.toggle('dark', t === 'dark')
    persist('theme', t)
  }

  function setPreferenceTheme(t: PreferenceTheme) {
    preferenceTheme.value = t
    persist('preference-setting', t)
  }

  function applyTheme() {
    if (preferenceTheme.value === 'auto') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      setTheme(prefersDark ? 'dark' : 'light')
    } else {
      setTheme(preferenceTheme.value)
    }
  }

  function setLanguage(lang: Locale) {
    language.value = lang
    setLocale(lang)
    persist('language', lang)
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
    persist('sidebar-collapsed', sidebarCollapsed.value)
  }

  // Initialize language from persisted value
  setLocale(language.value)

  return {
    theme, preferenceTheme, sidebarCollapsed, serverName, nginxStatus, language,
    setTheme, setPreferenceTheme, applyTheme, setLanguage, toggleSidebar,
  }
})