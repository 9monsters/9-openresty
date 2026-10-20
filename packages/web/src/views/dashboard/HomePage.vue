<script setup lang="ts">
import { t } from '@/i18n'
import EntryGroup from '@/components/EntryGroup.vue'
import { getDashboard, type DashboardData } from '@/api/dashboard'
import { getNginxStatus } from '@/api/nginx'
import { useSettingsStore } from '@/stores/settings'
import dayjs from 'dayjs'

const settings = useSettingsStore()

const dash = ref<DashboardData | null>(null)
const version = ref('')

onMounted(async () => {
  try {
    const info = await getNginxStatus()
    version.value = info.version || ''
    settings.nginxStatus = info.status as 'running' | 'stopped' | 'unknown'
  } catch { /* the shell keeps showing an unknown state */ }
  try {
    dash.value = await getDashboard()
    if (!version.value) version.value = dash.value.version || ''
  } catch { /* readouts stay empty rather than showing zeros */ }
})

const entries = computed(() => [
  { title: t('Server'), description: t('Server hardware and OS metrics — CPU, memory, disk, and uptime of the host machine.'), path: '/dashboard/server', icon: 'i-tabler:server' },
  { title: t('Nginx'), description: t('Real-time Nginx performance metrics — active connections, request rates, and connection status breakdown.'), path: '/dashboard/nginx', icon: 'i-tabler:activity' },
  { title: t('Sites'), description: t('Quick navigation to all configured sites with search, SSL status, and enabled state.'), path: '/dashboard/sites', icon: 'i-tabler:world-www' },
])

const state = computed(() =>
  settings.nginxStatus === 'running' ? 'ok' : settings.nginxStatus === 'stopped' ? 'bad' : 'warn',
)

const stateLabel = computed(() =>
  settings.nginxStatus === 'running' ? t('Running')
    : settings.nginxStatus === 'stopped' ? t('Stopped')
      : t('Unknown'),
)

const stats = computed(() => dash.value?.stats ?? null)

const serverTime = computed(() =>
  dash.value?.serverTime ? dayjs(dash.value.serverTime).format('YYYY-MM-DD HH:mm:ss') : '—',
)

function num(value?: number) {
  return typeof value === 'number' ? value.toLocaleString('en-US') : '—'
}
</script>

<template>
  <div>
    <!-- The box's own state is the most characteristic thing this product has,
         so it opens the page instead of a generic stat-card row. -->
    <section class="ios-group hero">
      <div class="hero-head">
        <span class="ios-dot" :data-state="state" :class="{ 'is-live': state === 'ok' }" />
        <h2 class="ios-title hero-state">{{ stateLabel }}</h2>
        <span v-if="version" class="ios-mono hero-version">{{ version }}</span>
      </div>

      <p class="ios-footnote hero-meta">
        {{ t('Server Time') }}
        <span class="ios-mono">{{ serverTime }}</span>
      </p>

      <dl class="ios-instrument">
        <div><dt>{{ t('Active Connections') }}</dt><dd>{{ num(stats?.activeConnections) }}</dd></div>
        <div><dt>{{ t('Accepts') }}</dt><dd>{{ num(stats?.accepts) }}</dd></div>
        <div><dt>{{ t('Handled') }}</dt><dd>{{ num(stats?.handled) }}</dd></div>
        <div><dt>{{ t('Requests') }}</dt><dd>{{ num(stats?.requests) }}</dd></div>
      </dl>
    </section>

    <p class="ios-section-label section">{{ t('Navigate') }}</p>
    <EntryGroup :entries="entries" />
  </div>
</template>

<style scoped>
.hero { padding: 20px 18px 6px; }

.hero-head {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.hero-state { margin: 0; }

.hero-version {
  margin-left: auto;
  color: var(--label-2);
  background: var(--fill);
  padding: 4px 10px;
  border-radius: var(--r-pill);
}

.hero-meta { margin: 6px 0 0 20px; }

/* the status dot indents the title by 8px + 12px, so the first readout follows */
.hero :deep(.ios-instrument > *:first-child) { padding-left: 20px; }

.section { margin: 26px 0 8px 18px; }
</style>
