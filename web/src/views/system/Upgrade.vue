<script setup lang="ts">
import { t } from '@/i18n'
import { getUpgradeInfo } from '@/api/system'

const loading = ref(true)
const upgrade = ref<any>(null)

onMounted(async () => {
  try {
    upgrade.value = await getUpgradeInfo()
  } catch { /**/ }
  loading.value = false
})
</script>

<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">{{ t('Upgrade') }}</h3>
    <ASpin v-if="loading" class="flex-center w-full py-10" />
    <ACard v-else-if="upgrade">
      <div class="flex items-center gap-6">
        <div>
          <div class="text-sm text-gray-500">Current Version</div>
          <div class="text-xl font-bold">{{ upgrade.currentVersion }}</div>
        </div>
        <span class="i-tabler:arrow-right text-2xl text-gray-400" />
        <div>
          <div class="text-sm text-gray-500">Latest Version</div>
          <div class="text-xl font-bold">{{ upgrade.latestVersion }}</div>
        </div>
        <AButton v-if="upgrade.updateAvailable" type="primary">{{ t('Upgrade') }}</AButton>
        <ATag v-else color="green">Up to date</ATag>
      </div>
    </ACard>
  </div>
</template>