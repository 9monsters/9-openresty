<script setup lang="ts">
import { t } from '@/i18n'

// Simple placeholder — edit stream is a form page
const route = useRoute()
const isEdit = computed(() => route.params.name !== undefined)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ isEdit ? t('Edit Stream') : t('Add Stream') }}</h3>
    </div>
    <ACard>
      <AEmpty :description="t('Stream editing is managed through the stream list page.')" />
      <div class="text-center mt-4">
        <AButton @click="$router.push('/streams')">{{ t('Back to List') }}</AButton>
      </div>
    </ACard>
  </div>
</template>