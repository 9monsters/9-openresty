<script setup lang="ts">
import { useSettingsStore } from '@/stores/settings'
import { getNginxStatus } from '@/api/nginx'

const router = useRouter()
const route = useRoute()
const settings = useSettingsStore()

// Build menu from routes
interface MenuItem {
  key: string
  label: string
  icon?: string
  children?: MenuItem[]
}

const menuItems = computed<MenuItem[]>(() => {
  const mainRoute = router.options.routes.find(r => r.path === '/')
  if (!mainRoute?.children) return []

  // Dashboard sub-items
  const dashboardItems: MenuItem[] = []
  const mainItems: MenuItem[] = []

  for (const r of mainRoute.children) {
    if (r.meta?.hidden) continue

    const path = r.path as string
    if (path.startsWith('dashboard')) {
      if (path === 'dashboard') continue // skip dashboard itself, use its children
      dashboardItems.push({
        key: path,
        label: r.meta?.title as string || path,
      })
      continue
    }

    mainItems.push({
      key: path,
      label: r.meta?.title as string || path,
      icon: r.meta?.icon as string,
    })
  }

  // Insert Dashboard with sub-items at the top
  return [
    {
      key: 'dashboard',
      label: '仪表盘',
      icon: 'i-tabler:dashboard',
      children: dashboardItems,
    },
    ...mainItems,
  ]
})

const selectedKeys = ref<string[]>([route.path.startsWith('/') ? route.path.slice(1) : route.path || 'dashboard'])
const openKeys = ref<string[]>(route.path.startsWith('/dashboard') ? ['dashboard'] : [])

watch(() => route.path, (p) => {
  const normalized = p.startsWith('/') ? p.slice(1) : p
  selectedKeys.value = [normalized]
  if (normalized.startsWith('dashboard')) {
    if (!openKeys.value.includes('dashboard')) openKeys.value = ['dashboard']
  }
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
        v-model:openKeys="openKeys"
        theme="dark"
        mode="inline"
        @click="({ key }) => router.push('/' + key)"
      >
        <template v-for="item in menuItems" :key="item.key">
          <ASubMenu v-if="item.children && item.children.length > 0" :key="item.key" :title="item.label">
            <template #icon><span :class="item.icon" /></template>
            <AMenuItem v-for="child in item.children" :key="child.key">
              {{ child.label }}
            </AMenuItem>
          </ASubMenu>
          <AMenuItem v-else :key="item.key">
            <template #icon><span :class="item.icon" /></template>
            {{ item.label }}
          </AMenuItem>
        </template>
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