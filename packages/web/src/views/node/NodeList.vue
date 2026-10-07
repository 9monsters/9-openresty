<script setup lang="ts">
import { listNodes } from '@/api/nodes'
const nodes = ref<any[]>([])
onMounted(async () => { try { nodes.value = ((await listNodes()) as any).nodes } catch {} })
</script>
<template>
  <div>
    <div class="flex-between mb-4"><h3 class="text-xl font-semibold m-0">节点管理</h3><AButton type="primary" size="small"><span class="i-tabler:plus" /> 添加节点</AButton></div>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <ACard v-for="node in nodes" :key="node.id" size="small" class="card-hover">
        <div class="flex-between mb-2"><span class="font-semibold">{{ node.name }}</span><ATag :color="node.status === 'online' ? 'green' : 'red'">{{ node.status }}</ATag></div>
        <div class="text-xs text-gray-400 space-y-1"><div>{{ node.host }}:{{ node.port }}</div><div>{{ node.version }}</div><div>Uptime: {{ node.uptime }}</div></div>
      </ACard>
    </div>
  </div>
</template>
