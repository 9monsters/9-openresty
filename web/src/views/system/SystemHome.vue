<script setup lang="ts">
import { getSelfCheck, getUpgradeInfo, getAbout } from '@/api/system'
import { t } from '@/i18n'

const loading = ref(true)
const checks = ref<any[]>([])
const upgrade = ref<any>(null)
const about = ref<any>(null)

async function fetchData() {
  loading.value = true
  try {
    const [c, u, a] = await Promise.all([getSelfCheck(), getUpgradeInfo(), getAbout()])
    checks.value = c.checks
    upgrade.value = u
    about.value = a
  } catch { /**/ }
  loading.value = false
}

async function handleRefreshSelfCheck() {
  try { const c = await getSelfCheck(); checks.value = c.checks } catch { /**/ }
}

const statusIcons: Record<string, string> = {
  ok: 'i-tabler:circle-check text-green-500',
  warning: 'i-tabler:alert-triangle text-orange-500',
  error: 'i-tabler:circle-x text-red-500',
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('System') }}</h3>
      <AButton size="small" @click="fetchData">{{ t('Refresh') }}</AButton>
    </div>

    <ASpin v-if="loading" class="flex-center w-full py-10" />

    <div v-else class="space-y-6">
      <!-- Self Check -->
      <ACard :title="t('Self Check')">
        <template #extra>
          <AButton size="small" @click="handleRefreshSelfCheck">{{ t('Refresh') }}</AButton>
        </template>
        <div class="space-y-2">
          <div v-for="check in checks" :key="check.name" class="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded">
            <span :class="statusIcons[check.status] || 'i-tabler:question-mark'" class="text-lg" />
            <div class="flex-1">
              <div class="font-semibold">{{ check.name }}</div>
              <div class="text-sm text-gray-500">{{ check.detail }}</div>
            </div>
            <ATag :color="check.status === 'ok' ? 'green' : check.status === 'warning' ? 'orange' : 'red'">{{ check.status }}</ATag>
          </div>
        </div>
      </ACard>

      <!-- Upgrade -->
      <ACard :title="t('Upgrade')" v-if="upgrade">
        <div class="flex items-center gap-6">
          <div>
            <div class="text-sm text-gray-500">{{ t('Current Version') || 'Current Version' }}</div>
            <div class="text-xl font-bold">{{ upgrade.currentVersion }}</div>
          </div>
          <span class="i-tabler:arrow-right text-2xl text-gray-400" />
          <div>
            <div class="text-sm text-gray-500">{{ t('Latest Version') || 'Latest Version' }}</div>
            <div class="text-xl font-bold">{{ upgrade.latestVersion }}</div>
          </div>
          <AButton v-if="upgrade.updateAvailable" type="primary">{{ t('Upgrade') }}</AButton>
          <ATag v-else color="green">{{ t('Up to date') || 'Up to date' }}</ATag>
        </div>
      </ACard>

      <!-- About -->
      <ACard :title="t('About')" v-if="about">
        <div class="grid grid-cols-4 gap-4">
          <div><div class="text-sm text-gray-500">{{ t('Version') || 'Version' }}</div><div class="font-semibold">{{ about.version }}</div></div>
          <div><div class="text-sm text-gray-500">Node.js</div><div class="font-semibold">{{ about.node }}</div></div>
          <div><div class="text-sm text-gray-500">OpenResty</div><div class="font-semibold">{{ about.openresty }}</div></div>
          <div><div class="text-sm text-gray-500">{{ t('Platform') || 'Platform' }}</div><div class="font-semibold">{{ about.platform }} / {{ about.arch }}</div></div>
        </div>
        <div class="mt-2 text-sm text-gray-500">{{ t('Build Time') || 'Build Time' }}: {{ about.buildTime }}</div>
      </ACard>
    </div>
  </div>
</template>