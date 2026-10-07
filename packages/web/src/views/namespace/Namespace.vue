<script setup lang="ts">
import { listNamespaces } from '@/api/namespaces'
const namespaces = ref<any[]>([])
onMounted(async () => { try { namespaces.value = ((await listNamespaces()) as any).namespaces } catch {} })
</script>
<template>
  <div>
    <div class="flex-between mb-4"><h3 class="text-xl font-semibold m-0">命名空间</h3><AButton type="primary" size="small"><span class="i-tabler:plus" /> 添加</AButton></div>
    <ATable :data-source="namespaces" row-key="name" :pagination="false">
      <ATableColumn title="名称" data-index="name" />
      <ATableColumn title="描述" data-index="description" />
      <ATableColumn title="站点数" data-index="siteCount" width="80" />
      <ATableColumn title="状态" data-index="enabled" width="80"><template #default="{ record }"><ATag :color="record.enabled ? 'green' : 'default'">{{ record.enabled ? '启用' : '禁用' }}</ATag></template></ATableColumn>
    </ATable>
  </div>
</template>
