<script setup lang="ts">
import { t } from '@/i18n'
import { getAbout } from '@/api/system'

const loading = ref(true)
const about = ref<any>(null)

onMounted(async () => {
  try {
    about.value = await getAbout()
  } catch { /**/ }
  loading.value = false
})
</script>

<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">{{ t('About') }}</h3>
    <ASpin v-if="loading" class="flex-center w-full py-10" />
    <ACard v-else-if="about">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div><div class="text-sm text-gray-500">Version</div><div class="font-semibold">{{ about.version }}</div></div>
        <div><div class="text-sm text-gray-500">OpenResty</div><div class="font-semibold">{{ about.openresty }}</div></div>
        <div><div class="text-sm text-gray-500">Node.js</div><div class="font-semibold">{{ about.node }}</div></div>
        <div><div class="text-sm text-gray-500">Platform</div><div class="font-semibold">{{ about.platform }} / {{ about.arch }}</div></div>
      </div>
      <div class="mt-4 text-sm text-gray-500">Build: {{ about.buildTime }}</div>
    </ACard>
  </div>
</template>