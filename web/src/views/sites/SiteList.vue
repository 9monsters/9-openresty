<script setup lang="ts">
import { listSites, deleteSite } from '@/api/sites'
import type { ParsedSite } from '@/types'

const router = useRouter()
const sites = ref<ParsedSite[]>([])
const loading = ref(true)

const columns = [
  { title: '站点名称', dataIndex: 'name', key: 'name' },
  { title: '监听端口', key: 'listen', width: 200 },
  { title: 'SSL', key: 'ssl', width: 80, align: 'center' as const },
  { title: 'Lua 指令数', key: 'lua', width: 100, align: 'center' as const },
  { title: '操作', key: 'action', width: 200 },
]

async function fetchSites() {
  loading.value = true
  try {
    const res = await listSites()
    sites.value = res.sites
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleDelete(site: ParsedSite) {
  try {
    await deleteSite(site.name)
    await fetchSites()
    message.success('站点已删除')
  } catch (e) {
    message.error('删除失败')
  }
}

onMounted(fetchSites)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">站点管理</h3>
      <AButton type="primary" @click="router.push('/sites/add')">
        <span class="i-tabler:plus mr-1" /> 添加站点
      </AButton>
    </div>

    <ATable
      :columns="columns"
      :data-source="sites"
      :loading="loading"
      row-key="name"
      :pagination="{ pageSize: 20 }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'listen'">
          <ATag v-for="l in record.listen" :key="l" class="mr-1" color="blue">{{ l }}</ATag>
        </template>
        <template v-else-if="column.key === 'ssl'">
          <span v-if="record.ssl" class="i-tabler:lock text-green-500 text-lg" />
          <span v-else class="i-tabler:lock-open text-gray-400 text-lg" />
        </template>
        <template v-else-if="column.key === 'lua'">
          <ATag :color="record.luaDirectives.length > 0 ? 'purple' : 'default'">
            {{ record.luaDirectives.length }}
          </ATag>
        </template>
        <template v-else-if="column.key === 'action'">
          <ASpace>
            <AButton size="small" @click="router.push(`/sites/${encodeURIComponent(record.name)}/edit`)">
              <span class="i-tabler:edit mr-1" /> 编辑
            </AButton>
            <AButton size="small" danger @click="handleDelete(record)">
              <span class="i-tabler:trash mr-1" /> 删除
            </AButton>
          </ASpace>
        </template>
      </template>
    </ATable>
  </div>
</template>
