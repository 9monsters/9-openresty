<script setup lang="ts">
import { listSnippets, listTemplates } from '@/api/snippets'
const snippets = ref<any[]>([])
const templates = ref<any[]>([])
onMounted(async () => {
  try {
    const [s, t] = await Promise.all([listSnippets(), listTemplates()])
    snippets.value = (s as any).snippets; templates.value = (t as any).templates
  } catch {}
})
</script>
<template>
  <div>
    <div class="flex-between mb-4"><h3 class="text-xl font-semibold m-0">代码片段</h3><AButton type="primary" size="small"><span class="i-tabler:plus" /> 创建片段</AButton></div>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <ACard title="用户片段" size="small">
        <div v-for="s in snippets" :key="s.name" class="flex-between py-2 border-b last:border-0 dark:border-gray-700">
          <div><div class="font-medium text-sm">{{ s.name }}</div><div class="text-xs text-gray-400">{{ s.description }}</div></div>
          <ATag class="text-xs">{{ s.usage }}</ATag>
        </div>
      </ACard>
      <ACard title="内置模板" size="small">
        <div v-for="t in templates" :key="t.name" class="flex-between py-2 border-b last:border-0 dark:border-gray-700">
          <div><div class="font-medium text-sm">{{ t.name }}</div><div class="text-xs text-gray-400">{{ t.description }}</div></div>
          <ATag color="blue" class="text-xs">{{ t.category }}</ATag>
        </div>
      </ACard>
    </div>
  </div>
</template>
