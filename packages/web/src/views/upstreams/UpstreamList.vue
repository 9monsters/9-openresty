<script setup lang="ts">
import { listUpstreams } from '@/api/upstreams'
import type { Upstream } from '@/types'

const upstreams = ref<Upstream[]>([])
const loading = ref(true)

async function fetch() {
  loading.value = true
  try {
    const res = await listUpstreams()
    upstreams.value = res.upstreams
  } finally {
    loading.value = false
  }
}

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { up: 'green', down: 'red', checking: 'orange' }
  return map[status] || 'default'
}

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">Upstream 管理</h3>
      <AButton type="primary">
        <span class="i-tabler:plus mr-1" /> 添加 Upstream
      </AButton>
    </div>

    <ACard v-for="us in upstreams" :key="us.name" :title="us.name" class="mb-4 card-hover">
      <template #extra>
        <ATag color="blue">{{ us.strategy }}</ATag>
      </template>

      <ATable :data-source="us.servers" :pagination="false" row-key="address">
        <ATableColumn title="服务器地址" data-index="address" key="address" />
        <ATableColumn title="权重" data-index="weight" key="weight" width="80" />
        <ATableColumn title="状态" key="status" width="100">
          <template #default="{ record }: { record: any }">
            <ATag :color="getStatusColor(record.status)">{{ record.status }}</ATag>
          </template>
        </ATableColumn>
      </ATable>
    </ACard>

    <AEmpty v-if="!loading && upstreams.length === 0" description="暂无 Upstream 配置" />
    <ASpin v-if="loading" class="flex-center py-10" />
  </div>
</template>
