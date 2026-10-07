<script setup lang="ts">
import { selfCheck, checkUpgrade, getAbout } from '@/api/system'
const checks = ref<any[]>([])
const upgrade = ref<any>(null)
const about = ref<any>(null)
onMounted(async () => {
  try {
    const [c, u, a] = await Promise.all([selfCheck(), checkUpgrade(), getAbout()])
    checks.value = (c as any).checks; upgrade.value = (u as any); about.value = (a as any)
  } catch {}
})
const statusColor = (s: string) => ({ ok: '#52c41a', warning: '#faad14', error: '#ff4d4f' } as any)[s] || '#999'
</script>
<template>
  <div class="max-w-3xl">
    <h3 class="text-xl font-semibold mb-4">系统</h3>
    <ACard title="自检" size="small" class="mb-4">
      <div v-for="c in checks" :key="c.name" class="flex-between py-2 border-b last:border-0 dark:border-gray-700">
        <div><span class="font-medium">{{ c.name }}</span><div class="text-xs text-gray-400">{{ c.detail }}</div></div>
        <span :style="{ color: statusColor(c.status), fontSize: '18px' }" :class="`i-tabler:${c.status === 'ok' ? 'circle-check' : c.status === 'warning' ? 'alert-triangle' : 'circle-x'}`" />
      </div>
    </ACard>
    <ACard title="版本更新" size="small" class="mb-4">
      <template v-if="upgrade">
        <div class="flex-between"><span>当前版本</span><span class="font-mono">{{ upgrade.currentVersion }}</span></div>
        <div class="flex-between mt-2"><span>最新版本</span><span class="font-mono text-green-600">{{ upgrade.latestVersion }}</span></div>
        <AButton v-if="upgrade.updateAvailable" type="primary" size="small" class="mt-3">升级到 {{ upgrade.latestVersion }}</AButton>
      </template>
    </ACard>
    <ACard title="关于" size="small">
      <template v-if="about">
        <div class="grid grid-cols-2 gap-2 text-sm">
          <div class="text-gray-500">应用版本</div><div class="font-mono">{{ about.version }}</div>
          <div class="text-gray-500">Node.js</div><div>{{ about.node }}</div>
          <div class="text-gray-500">OpenResty</div><div>{{ about.openresty }}</div>
          <div class="text-gray-500">平台</div><div>{{ about.platform }} {{ about.arch }}</div>
        </div>
      </template>
    </ACard>
  </div>
</template>
