<script setup lang="ts">
import { listSites } from '@/api/sites'
import { t } from '@/i18n'

const router = useRouter()
const loading = ref(true)
const sites = ref<any[]>([])
const searchText = ref('')

const columns = [
  { title: t('Site Name') || 'Site Name', dataIndex: 'name', key: 'name' },
  { title: t('Domains') || 'Domains', key: 'domains' },
  { title: 'SSL', key: 'ssl', width: 60 },
  { title: t('Listen') || 'Listen', key: 'listen', width: 120 },
  { title: t('Status'), key: 'status', width: 100 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await listSites()
    sites.value = (res.sites || []).filter((s: any) => {
      if (!searchText.value) return true
      const q = searchText.value.toLowerCase()
      return (s.name || '').toLowerCase().includes(q)
    })
  } catch { /**/ }
  loading.value = false
}

watch(searchText, fetchData)

function goToSite(record: any) {
  router.push(`/sites/${encodeURIComponent(record.name)}/edit`)
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Site Navigation') }}</h3>
      <AInputSearch v-model:value="searchText" :placeholder="t('Search')" style="width: 250px" />
    </div>

    <ATable
      :dataSource="sites"
      :columns="columns"
      :loading="loading"
      rowKey="name"
      :pagination="false"
      @row-click="goToSite"
      :customRow="() => ({ style: { cursor: 'pointer' } })"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'domains'">
          <ATag v-for="d in (record.domains || [])" :key="d" class="mr-1">{{ d }}</ATag>
          <span v-if="!record.domains?.length" class="text-gray-400">-</span>
        </template>
        <template v-if="column.key === 'ssl'">
          <span v-if="record.ssl" class="i-tabler:lock text-green-500" />
          <span v-else class="i-tabler:lock-open text-gray-400" />
        </template>
        <template v-if="column.key === 'listen'">
          <ATag v-for="p in (record.listenPorts || [])" :key="p" class="mr-1">{{ p }}</ATag>
        </template>
        <template v-if="column.key === 'status'">
          <ATag :color="record.enabled !== false ? 'green' : 'red'">
            {{ record.enabled !== false ? t('Enabled') : t('Disabled') }}
          </ATag>
        </template>
      </template>
    </ATable>
  </div>
</template>