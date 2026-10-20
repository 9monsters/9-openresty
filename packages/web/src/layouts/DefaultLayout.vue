<script setup lang="ts">
import type { Component } from 'vue'
import { RouterLink } from 'vue-router'
import { useSettingsStore } from '@/stores/settings'
import { t } from '@/i18n'
import { getNginxStatus } from '@/api/nginx'
import { getServerName } from '@/api/settings'
import { useUserStore } from '@/stores/user'
import {
  HomeOutlined, CloudOutlined, ShareAltOutlined, ClusterOutlined,
  SafetyOutlined, FileOutlined, SafetyCertificateOutlined, FileTextOutlined,
  CodeOutlined, StarOutlined, CloudServerOutlined,
  DatabaseOutlined, BellOutlined, SettingOutlined, ClockCircleOutlined,
  InfoCircleOutlined, UsergroupAddOutlined,
} from '@antdv-next/icons'
import Logo from '@/components/Logo.vue'

const router = useRouter()
const route = useRoute()
const settings = useSettingsStore()
const user = useUserStore()

// --- Responsive ---
const drawerVisible = ref(false)
const collapsed = ref(false)
const hideLayoutSidebar = ref(false)

let resizeTimer: ReturnType<typeof setTimeout> | null = null
function _init() {
  collapsed.value = document.body.clientWidth < 1080
  hideLayoutSidebar.value = document.body.clientWidth < 760
}
function init() {
  if (resizeTimer) clearTimeout(resizeTimer)
  resizeTimer = setTimeout(_init, 50)
}
onMounted(() => { init(); addEventListener('resize', init) })
onUnmounted(() => removeEventListener('resize', init))

// --- Scroll edge effect: the bar earns its hairline once content slides under it ---
const scrolled = ref(false)
function onScroll() {
  scrolled.value = window.scrollY > 6
}
onMounted(() => { onScroll(); addEventListener('scroll', onScroll, { passive: true }) })
onUnmounted(() => removeEventListener('scroll', onScroll))

// --- Sidebar menu from routes ---
interface SidebarEntry {
  path: string
  name: string
  meta: { icon?: Component; hiddenInSidebar?: boolean; hideChildren?: boolean; name: string; modules?: string[] }
  children: SidebarEntry[]
}

const sidebars = computed(() => {
  return (router.options.routes[0]?.children || []) as unknown as SidebarEntry[]
})

const visible = computed(() => {
  return sidebars.value.filter(s => !s.meta?.hiddenInSidebar).map(s => ({
    ...s,
    children: (s.children || []).filter(c => !(c.meta as any)?.hiddenInSidebar),
  }))
})

// Icons map — must match route meta.icon strings
const iconMap: Record<string, Component> = {
  HomeOutlined, CloudOutlined, ShareAltOutlined, ClusterOutlined,
  SafetyOutlined, FileOutlined, SafetyCertificateOutlined, FileTextOutlined,
  CodeOutlined, StarOutlined, CloudServerOutlined,
  DatabaseOutlined, BellOutlined, SettingOutlined, ClockCircleOutlined,
  InfoCircleOutlined, UsergroupAddOutlined,
}

const openKeys = ref<string[]>([])
const selectedKeys = ref<string[]>([])

// Every key the sidebar can actually highlight. Routes such as `Dashboard Home`
// are hidden in the sidebar, so selection has to fall back to their group.
const menuKeys = computed(() => {
  const keys = new Set<string>()
  for (const group of visible.value) {
    keys.add(group.name)
    for (const child of group.children) keys.add(child.name)
  }
  return keys
})

function syncSelection() {
  const names = route.matched.map(r => r.name as string).filter(Boolean).reverse()
  const selected = names.find(n => menuKeys.value.has(n)) ?? (route.name as string)
  selectedKeys.value = [selected]

  const owner = visible.value.find(g => g.name === selected || g.children.some(c => c.name === selected))
  openKeys.value = owner && owner.children.length ? [owner.name] : []
}

// immediate: without it the group for the landing route never opens.
watch(route, syncSelection, { immediate: true })
watch(visible, syncSelection)

function translate(name: string | (() => string)): string {
  return typeof name === 'function' ? name() : t(name)
}

const menuItems = computed(() => {
  return visible.value.map(s => {
    const icon = s.meta?.icon ? iconMap[s.meta.icon as unknown as string] : undefined
    if (!s.children.length) {
      return {
        key: s.name,
        icon,
        label: translate(s.meta?.name),
        title: translate(s.meta?.name),
        onClick: () => router.push(`/${s.path}`),
      }
    }
    return {
      key: s.name,
      icon,
      label: translate(s.meta?.name),
      title: translate(s.meta?.name),
      children: s.children.map(c => ({
        key: c.name,
        label: h(RouterLink, { to: `/${s.path}/${c.path}` }, () => translate((c.meta as any)?.name)),
      })),
    }
  })
})

// --- Server identity ---
const { serverName } = storeToRefs(settings)
onMounted(async () => {
  try {
    const r = await getServerName()
    serverName.value = r.name
  } catch { /* ignore */ }
  try {
    const data = await getNginxStatus()
    settings.nginxStatus = data.status as 'running' | 'stopped' | 'unknown'
  } catch { /* ignore */ }
})

// --- Page title + rail state ---
provide('sidebarCollapsed', collapsed)

const pageTitle = computed(() => {
  const name = route.meta?.name
  return typeof name === 'string' ? t(name) : (typeof name === 'function' ? name() : t('Dashboard'))
})

const breadcrumbTrail = computed(() =>
  route.matched
    .filter(r => r.meta?.title)
    .map(r => t(r.meta.title as string)),
)

// --- Theme ---
function toggleTheme() {
  const next = settings.preferenceTheme === 'dark' ? 'light'
    : settings.preferenceTheme === 'light' ? 'auto' : 'dark'
  settings.setPreferenceTheme(next)
  settings.applyTheme()
}
const themeIcon = computed(() => {
  switch (settings.preferenceTheme) {
    case 'dark': return 'i-tabler:moon'
    case 'light': return 'i-tabler:sun'
    default: return 'i-tabler:sun-moon'
  }
})
const themeLabel = computed(() => t(settings.preferenceTheme === 'dark' ? 'Dark' : settings.preferenceTheme === 'light' ? 'Light' : 'Auto'))

function toggleLanguage() {
  settings.setLanguage(settings.language === 'zh-CN' ? 'en-US' : 'zh-CN')
}

const statusState = computed(() => settings.nginxStatus === 'running' ? 'ok' : settings.nginxStatus === 'stopped' ? 'bad' : 'warn')
const statusLabel = computed(() => settings.nginxStatus === 'running' ? t('Running') : settings.nginxStatus === 'stopped' ? t('Stopped') : t('Unknown'))
</script>

<template>
  <div class="shell">
    <!-- Mobile sheet -->
    <div class="rail-mobile">
      <ADrawer
        v-model:open="drawerVisible"
        :closable="false"
        placement="left"
        :width="286"
        :styles="{ body: { padding: 0 }, header: { display: 'none' } }"
      >
        <div class="rail-body">
          <Logo />
          <AMenu
            v-model:open-keys="openKeys"
            v-model:selected-keys="selectedKeys"
            mode="inline"
            :items="menuItems"
          />
        </div>
      </ADrawer>
    </div>

    <!-- Desktop glass rail -->
    <aside v-if="!hideLayoutSidebar" class="rail ios-glass" :class="{ 'is-collapsed': collapsed }">
      <Logo />
      <AMenu
        class="rail-menu"
        v-model:open-keys="openKeys"
        v-model:selected-keys="selectedKeys"
        mode="inline"
        :inline-collapsed="collapsed"
        :items="menuItems"
      />
      <div class="rail-foot">
        <span class="ios-caption rail-host">{{ serverName }}</span>
        <button class="rail-tool" type="button" :aria-label="t('Collapse')" @click="collapsed = !collapsed">
          <span :class="collapsed ? 'i-tabler:layout-sidebar-left-expand' : 'i-tabler:layout-sidebar-left-collapse'" />
        </button>
      </div>
    </aside>

    <div class="column" :class="{ 'is-wide': hideLayoutSidebar }">
      <!-- Glass bar. It is chrome, so it floats over content and gains its
           hairline only once content actually slides beneath it. -->
      <header class="bar ios-glass" :class="{ 'is-scrolled': scrolled }">
        <button v-if="hideLayoutSidebar" class="bar-tool" type="button" aria-label="Menu" @click="drawerVisible = true">
          <span class="i-tabler:menu-2" />
        </button>

        <nav class="bar-trail" aria-label="Breadcrumb">
          <Transition name="bar-swap" mode="out-in">
            <span v-if="scrolled" key="title" class="bar-title">{{ pageTitle }}</span>
            <span v-else key="trail" class="ios-footnote">{{ breadcrumbTrail.length ? breadcrumbTrail.join(' / ') : pageTitle }}</span>
          </Transition>
        </nav>

        <div class="bar-end">
          <button class="bar-tool bar-lang" type="button" :aria-label="settings.language === 'zh-CN' ? 'Switch to English' : '切换到中文'" @click="toggleLanguage">
            {{ settings.language === 'zh-CN' ? 'EN' : '中' }}
          </button>
          <button class="bar-tool" type="button" :aria-label="themeLabel" @click="toggleTheme">
            <span :class="themeIcon" />
          </button>
          <span class="ios-chip" :data-state="statusState">
            <span class="ios-dot" :data-state="statusState" :class="{ 'is-live': statusState === 'ok' }" />
            {{ statusLabel }}
          </span>
          <button class="bar-tool" type="button" :aria-label="t('Logout', '退出')" @click="user.logout(); router.push('/login')">
            <span class="i-tabler:logout" />
          </button>
        </div>
      </header>

      <main class="content">
        <div class="page">
          <RouterView v-slot="{ Component, route: r }">
            <Transition name="page" mode="out-in">
              <component :is="Component" :key="r.path" />
            </Transition>
          </RouterView>
        </div>

        <footer class="foot ios-caption">
          {{ serverName }} · {{ new Date().getFullYear() }}
        </footer>
      </main>
    </div>
  </div>
</template>

<style scoped>
.shell {
  display: flex;
  align-items: flex-start;
  gap: 22px;
  min-height: 100vh;
  padding: 14px;
}

/* ------------------------------------------------------------------- rail */

.rail {
  position: sticky;
  top: 14px;
  z-index: 12;
  display: flex;
  flex-direction: column;
  flex: none;
  width: 244px;
  height: calc(100vh - 28px);
  border-radius: var(--r-rail);
  padding: 10px 8px 8px;
  transition: width var(--dur) var(--ease-ios);
  overflow: hidden;
}
.rail.is-collapsed { width: 76px; }

.rail-menu {
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  margin-top: 8px;
}

.rail-foot {
  flex: none;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 6px 2px;
  border-top: 0.5px solid var(--separator);
}
.rail-host {
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.rail.is-collapsed .rail-host { display: none; }

.rail-tool,
.bar-tool {
  flex: none;
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  padding: 0;
  border: 0;
  border-radius: var(--r-pill);
  background: var(--fill);
  color: var(--label-2);
  font-size: 16px;
  cursor: pointer;
  transition: background var(--dur) var(--ease-ios), color var(--dur) var(--ease-ios), transform 120ms var(--ease-ios);
}
.rail-tool:hover,
.bar-tool:hover { background: var(--fill-strong); color: var(--label); }
.rail-tool:active,
.bar-tool:active { transform: scale(0.92); }
.rail.is-collapsed .rail-foot { justify-content: center; }

.rail-body { display: flex; flex-direction: column; height: 100%; }

.rail-mobile :deep(.ant-drawer-body) { background: var(--glass-strong); }

/* -------------------------------------------------------------------- bar */

.column {
  flex: 1 1 auto;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.bar {
  position: sticky;
  top: 14px;
  z-index: 11;
  display: flex;
  align-items: center;
  gap: 12px;
  height: 54px;
  padding: 0 10px 0 16px;
  border-radius: var(--r-bar);
  transition: box-shadow var(--dur) var(--ease-ios), border-color var(--dur) var(--ease-ios), background var(--dur) var(--ease-ios);
}
.bar:not(.is-scrolled) {
  border-color: transparent;
  box-shadow: none;
  background: transparent;
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
}

.bar-trail { flex: 1 1 auto; min-width: 0; overflow: hidden; }
.bar-trail .ios-footnote { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.bar-title {
  font-size: 17px;
  font-weight: 650;
  letter-spacing: -0.016em;
  color: var(--label);
}
.bar-swap-enter-active { transition: opacity 180ms var(--ease-ios), transform 180ms var(--ease-ios); }
.bar-swap-leave-active { transition: opacity 90ms var(--ease-ios); }
.bar-swap-enter-from { opacity: 0; transform: translateY(4px); }
.bar-swap-leave-to { opacity: 0; }

.bar-end { flex: none; display: flex; align-items: center; gap: 8px; margin-left: auto; }
.bar-lang { font-size: 12px; font-weight: 700; letter-spacing: 0; }

/* ---------------------------------------------------------------- content */

.content {
  flex: 1 1 auto;
  display: flex;
  flex-direction: column;
  gap: 26px;
  padding: 6px 6px 0;
}

.page {
  width: 100%;
  max-width: 1180px;
  margin: 0 auto;
}

.foot {
  max-width: 1180px;
  width: 100%;
  margin: 0 auto;
  padding: 8px 4px 28px;
  border-top: 0.5px solid var(--separator);
}

/* The shell owns the bar title; views own their own headings, which the
   design layer promotes to the iOS large title. */
.page :deep(:is(h2, h3).text-2xl),
.page :deep(:is(h2, h3).text-xl) {
  font-size: 34px;
  font-weight: 700;
  letter-spacing: -0.022em;
  line-height: 1.14;
  color: var(--label);
}

/* --------------------------------------------------------------- motion */

.page-enter-active { transition: opacity 220ms var(--ease-ios), transform 220ms var(--ease-ios); }
.page-leave-active { transition: opacity 120ms var(--ease-ios); }
.page-enter-from { opacity: 0; transform: translateY(6px); }
.page-leave-to { opacity: 0; }

@media (max-width: 760px) {
  .shell { padding: 10px; gap: 0; }
  .bar { top: 10px; height: 50px; padding: 0 8px 0 12px; }
  .content { padding: 4px 2px 0; }
}
</style>
