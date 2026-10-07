<script setup lang="ts">
import { listCertificates } from '@/api/certificates'
import type { Certificate } from '@/types'
import dayjs from 'dayjs'

const certificates = ref<Certificate[]>([])
const loading = ref(true)

const columns = [
  { title: '证书名称', dataIndex: 'name', key: 'name' },
  { title: '域名', key: 'domains' },
  { title: '颁发者', dataIndex: 'issuer', key: 'issuer' },
  { title: '到期时间', key: 'expiry' },
  { title: '自动续期', key: 'autoRenew', width: 100, align: 'center' as const },
  { title: '状态', key: 'status', width: 100 },
]

async function fetch() {
  loading.value = true
  try {
    const res = await listCertificates()
    certificates.value = res.certificates
  } finally {
    loading.value = false
  }
}

const getStatusTag = (status: string) => {
  const map: Record<string, { color: string; text: string }> = {
    valid: { color: 'green', text: '有效' },
    expiring: { color: 'orange', text: '即将过期' },
    expired: { color: 'red', text: '已过期' },
  }
  return map[status] || { color: 'default', text: status }
}

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">证书管理</h3>
      <AButton type="primary">
        <span class="i-tabler:plus mr-1" /> 添加证书
      </AButton>
    </div>

    <ATable :columns="columns" :data-source="certificates" :loading="loading" row-key="name">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'domains'">
          <ATag v-for="d in record.domains" :key="d" class="mr-1">{{ d }}</ATag>
        </template>
        <template v-else-if="column.key === 'expiry'">
          <span :class="{ 'text-orange-500': record.status === 'expiring', 'text-red-500': record.status === 'expired' }">
            {{ dayjs(record.notAfter).format('YYYY-MM-DD') }}
          </span>
        </template>
        <template v-else-if="column.key === 'autoRenew'">
          <span v-if="record.autoRenew" class="i-tabler:refresh text-green-500" />
          <span v-else class="i-tabler:refresh-off text-gray-400" />
        </template>
        <template v-else-if="column.key === 'status'">
          <ATag :color="getStatusTag(record.status).color">{{ getStatusTag(record.status).text }}</ATag>
        </template>
      </template>
    </ATable>
  </div>
</template>
