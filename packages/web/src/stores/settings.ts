import { defineStore } from 'pinia'

export const useSettingsStore = defineStore('settings', () => {
  const theme = ref<'light' | 'dark'>('light')
  const sidebarCollapsed = ref(false)
  const serverName = ref('OpenResty Admin')
  const nginxStatus = ref<'running' | 'stopped' | 'unknown'>('unknown')

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTheme(t: 'light' | 'dark') {
    theme.value = t
    document.documentElement.classList.toggle('dark', t === 'dark')
  }

  return { theme, sidebarCollapsed, serverName, nginxStatus, toggleSidebar, setTheme }
})
