<script setup lang="ts">
import { useSettingsStore } from '@/stores/settings'
import { getNginxStatus } from '@/api/nginx'

const router = useRouter()
const route = useRoute()
const settings = useSettingsStore()

// Sidebar menu items
const menuItems = computed(() => {
  return router.options.routes
    .find(r => r.path === '/')
    ?.children
    ?.filter(r => !r.meta?.hidden)
    .map(r => ({
      key: r.path,
      label: r.meta?.title as string,
      icon: r.meta?.icon as string,
    })) || []
})

const selectedKeys = computed(() => {
  const path = route.path.startsWith('/') ? route.path.slice(1) : route.path
  return [path.split('/')[0] || 'dashboard']
})

// Fetch nginx status
onMounted(async () => {
  try {
    const data = await getNginxStatus()
    settings.nginxStatus = data.status as 'running' | 'stopped' | 'unknown'
  } catch { /* ignore */ }
})
</script>

<template>
  <ALayout class="min-h-screen">
    <!-- Sidebar -->
    <ALayoutSider
      v-model:collapsed="settings.sidebarCollapsed"
      collapsible
      theme="dark"
      :style="{ overflow: 'auto', height: '100vh', position: 'fixed', left: 0, top: 0, bottom: 0, zIndex: 10 }"
    >
      <div class="flex-center h-16 text-white text-lg font-bold gap-2 px-4">
        <span class="i-tabler:brand-openai text-2xl" />
        <span v-if="!settings.sidebarCollapsed">OpenResty Admin</span>
      </div>
      <AMenu
        v-model:selectedKeys="selectedKeys"
        theme="dark"
        mode="inline"
        @click="({ key }) => router.push('/' + key)"
      >
        <AMenuItem v-for="item in menuItems" :key="item.key">
          <template #icon>
            <span :class="item.icon" />
          </template>
          {{ item.label }}
        </AMenuItem>
      </AMenu>
    </ALayoutSider>

    <!-- Main content -->
    <ALayout :style="{ marginLeft: settings.sidebarCollapsed ? '80px' : '200px', transition: 'margin-left 0.2s' }">
      <!-- Header -->
      <ALayoutHeader class="bg-white dark:bg-gray-800 px-6 flex-between border-b border-gray-200 dark:border-gray-700 h-16">
        <div class="flex items-center gap-3">
          <span class="text-lg font-semibold">{{ route.meta?.title || 'OpenResty Admin' }}</span>
        </div>
        <div class="flex items-center gap-4">
          <ATag :color="settings.nginxStatus === 'running' ? 'green' : 'red'">
            <template #icon>
              <span :class="settings.nginxStatus === 'running' ? 'i-tabler:circle-check' : 'i-tabler:circle-x'" />
            </template>
            Nginx: {{ settings.nginxStatus === 'running' ? '运行中' : '已停止' }}
          </ATag>
          <AButton
            type="text"
            @click="settings.setTheme(settings.theme === 'dark' ? 'light' : 'dark')"
          >
            <span :class="settings.theme === 'dark' ? 'i-tabler:sun' : 'i-tabler:moon'" class="text-lg" />
          </AButton>
        </div>
      </ALayoutHeader>

      <!-- Page content -->
      <ALayoutContent class="m-6 p-6 bg-white dark:bg-gray-800 rounded-lg min-h-[calc(100vh-8rem)]">
        <RouterView />
      </ALayoutContent>

      <!-- Footer -->
      <ALayoutFooter class="text-center text-gray-400 text-sm py-4">
        OpenResty Admin &copy; 2026
      </ALayoutFooter>
    </ALayout>
  </ALayout>
</template>

<style scoped>
:deep(.ant-layout-sider) {
  box-shadow: 2px 0 8px rgba(0,0,0,0.15);
}
</style>