<script setup lang="ts">
import { listAccessLists } from '@/api/access-lists'
const lists = ref<any[]>([])
const loading = ref(true)
onMounted(async () => {
  try { lists.value = ((await listAccessLists()) as any).lists } catch {}
  loading.value = false
})
</script>
<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">访问控制列表</h3>
    <ACard v-for="list in lists" :key="list.name" :title="list.name" size="small" class="mb-3">
      <ATable :data-source="list.rules" :pagination="false" size="small">
        <ATableColumn title="指令" data-index="directive" width="80"><template #default="{ record }"><ATag :color="record.directive === 'allow' ? 'green' : 'red'">{{ record.directive }}</ATag></template></ATableColumn>
        <ATableColumn title="地址" data-index="address" />
      </ATable>
    </ACard>
    <ASpin v-if="loading" class="flex-center py-10" />
  </div>
</template>
