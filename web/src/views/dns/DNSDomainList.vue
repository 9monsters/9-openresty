<script setup lang="ts">
import { listDnsDomains, createDnsDomain, updateDnsDomain, deleteDnsDomain } from '@/api/dns'
import { t } from '@/i18n'

const loading = ref(true)
const domains = ref<any[]>([])
const visible = ref(false)
const editing = ref<any>(null)
const form = reactive({ name: '', provider: '' })

async function fetch() {
  loading.value = true
  try { const r = await listDnsDomains(); domains.value = r.domains } catch { /**/ }
  loading.value = false
}
function add() { editing.value = null; form.name = ''; form.provider = ''; visible.value = true }
function edit(r: any) { editing.value = r; form.name = r.name; form.provider = r.provider; visible.value = true }
async function save() {
  try {
    if (editing.value) await updateDnsDomain(editing.value.id, { ...form })
    else await createDnsDomain({ ...form })
    visible.value = false; fetch()
  } catch { /**/ }
}
async function del(r: any) { try { await deleteDnsDomain(r.id); fetch() } catch { /**/ } }

const columns = [
  { title: 'Domain', dataIndex: 'name', key: 'name' },
  { title: 'Provider', dataIndex: 'provider', key: 'provider' },
  { title: 'Records', dataIndex: 'recordCount', key: 'recordCount', width: 80 },
  { title: 'Status', dataIndex: 'status', key: 'status', width: 100 },
  { title: '', key: 'actions', width: 150 },
]

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('DNS Domains') }}</h3>
      <AButton type="primary" size="small" @click="add"><span class="i-tabler:plus mr-1" /> Add</AButton>
    </div>
    <ATable :data-source="domains" :loading="loading" :columns="columns" row-key="id" :pagination="false" size="small">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <ATag color="blue">{{ record.status }}</ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="edit(record)">{{ t('Edit') }}</AButton>
          <AButton type="link" size="small" danger @click="del(record)">{{ t('Delete') }}</AButton>
        </template>
      </template>
    </ATable>
    <AModal v-model:open="visible" :title="(editing ? 'Edit' : 'Add') + ' Domain'" @ok="save">
      <AForm layout="vertical">
        <AFormItem label="Domain" required><AInput v-model:value="form.name" placeholder="example.com" /></AFormItem>
        <AFormItem label="Provider"><AInput v-model:value="form.provider" /></AFormItem>
      </AForm>
    </AModal>
  </div>
</template>