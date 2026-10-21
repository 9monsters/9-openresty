<script setup lang="ts">
import { listUpstreams, createUpstream, updateUpstream, deleteUpstream } from '@/api/upstreams'
import { t } from '@/i18n'

const loading = ref(true)
const upstreams = ref<any[]>([])
const modalVisible = ref(false)
const editing = ref<any>(null)
const formData = reactive({ name: '', strategy: 'round_robin', servers: '' })

async function fetchData() {
  loading.value = true
  try {
    const res = await listUpstreams()
    upstreams.value = res.upstreams
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  editing.value = null
  formData.name = ''
  formData.strategy = 'round_robin'
  formData.servers = ''
  modalVisible.value = true
}

function handleEdit(record: any) {
  editing.value = record
  formData.name = record.name
  formData.strategy = record.strategy || 'round_robin'
  formData.servers = (record.servers || []).map((s: any) => `${s.address} weight=${s.weight}`).join('\n')
  modalVisible.value = true
}

async function handleDelete(record: any) {
  try {
    await deleteUpstream(record.id)
    fetchData()
  } catch { /**/ }
}

async function handleSubmit() {
  try {
    const servers = formData.servers
      .split('\n')
      .filter((l: string) => l.trim())
      .map((l: string) => {
        const parts = l.trim().split(/\s+/)
        const address = parts[0]
        const weight = parts.find((p: string) => p.startsWith('weight='))?.split('=')[1] || '1'
        return { address, weight: Number(weight), maxFails: 3, failTimeout: '30s', status: 'up' }
      })

    if (editing.value) {
      await updateUpstream(editing.value.id, { name: formData.name, strategy: formData.strategy, servers })
    } else {
      await createUpstream({ name: formData.name, strategy: formData.strategy, servers })
    }
    modalVisible.value = false
    fetchData()
  } catch { /**/ }
}

const serverColumns = [
  { title: 'Address', dataIndex: 'address', key: 'address' },
  { title: 'Weight', dataIndex: 'weight', key: 'weight', width: 80 },
  { title: 'Max Fails', dataIndex: 'maxFails', key: 'maxFails', width: 100 },
  { title: 'Timeout', dataIndex: 'failTimeout', key: 'failTimeout', width: 100 },
  { title: 'Status', dataIndex: 'status', key: 'status', width: 80 },
]

function getStatusTag(status: string) {
  if (status === 'up') return 'green'
  if (status === 'down') return 'red'
  return 'orange'
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Upstream') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Add') }}
      </AButton>
    </div>

    <ASpin v-if="loading" class="flex-center w-full py-10" />

    <div v-else class="space-y-4">
      <ACard v-for="u in upstreams" :key="u.id" :title="u.name" size="small">
        <template #extra>
          <ATag>{{ u.strategy }}</ATag>
          <AButton size="small" type="link" @click="handleEdit(u)">{{ t('Edit') }}</AButton>
          <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(u)">
            <AButton size="small" type="link" danger>{{ t('Delete') }}</AButton>
          </APopconfirm>
        </template>

        <ATable :data-source="u.servers" :pagination="false" size="small" :columns="serverColumns" row-key="id">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <ATag :color="getStatusTag(record.status)">{{ record.status }}</ATag>
            </template>
          </template>
        </ATable>
      </ACard>

      <AEmpty v-if="upstreams.length === 0" :description="t('No Data')" />
    </div>

    <AModal v-model:open="modalVisible" :title="(editing ? t('Edit') : t('Add')) + ' Upstream'" @ok="handleSubmit" width="600px">
      <AForm layout="vertical">
        <AFormItem :label="t('Name') || 'Name'" required>
          <AInput v-model:value="formData.name" placeholder="backend_api" />
        </AFormItem>
        <AFormItem :label="t('Strategy') || 'Strategy'">
          <ASelect v-model:value="formData.strategy">
            <ASelectOption value="round_robin">round_robin</ASelectOption>
            <ASelectOption value="ip_hash">ip_hash</ASelectOption>
            <ASelectOption value="least_conn">least_conn</ASelectOption>
            <ASelectOption value="random">random</ASelectOption>
          </ASelect>
        </AFormItem>
        <AFormItem :label="t('Servers') || 'Servers'">
          <ATextarea v-model:value="formData.servers" :rows="6" placeholder="10.0.0.1:8080 weight=1&#10;10.0.0.2:8080 weight=2" style="font-family: monospace" />
          <div class="text-xs text-gray-400 mt-1">Format: address:port weight=N (one per line)</div>
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>