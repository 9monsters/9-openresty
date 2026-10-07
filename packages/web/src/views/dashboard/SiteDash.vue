<script setup lang="ts">
import { listSites } from '@/api/sites'
const sites = ref<any[]>([])
onMounted(async () => { try { sites.value = ((await listSites()) as any).sites } catch {} })
</script>
<template>
  <div>
    <ATable :data-source="sites" row-key="name" size="small" :pagination="false">
      <ATableColumn title="站点" data-index="name" />
      <ATableColumn title="域名" width="200"><template #default="{ record }"><ATag v-for="n in record.serverName" :key="n" class="mr-1" size="small">{{ n }}</ATag></template></ATableColumn>
      <ATableColumn title="SSL" width="60"><template #default="{ record }"><span :class="record.ssl ? 'i-tabler:lock text-green-500' : 'i-tabler:lock-open text-gray-400'" /></template></ATableColumn>
    </ATable>
  </div>
</template>
