<script setup lang="ts">
import { listDdnsEntries, createDdnsEntry, updateDdnsEntry, deleteDdnsEntry } from '@/api/dns'
import { t } from '@/i18n'

const loading = ref(true)
const entries = ref<any[]>([])
const visible = ref(false)
const editing = ref<any>(null)
const form = reactive({ domain: '', credential: '' })

async function fetch() {
  loading.value = true
  try { const r = await listDdnsEntries(); entries.value = r.ddns } catch { /**/ }
  loading.value = false
}
function add() { editing.value = null; form.domain = ''; form.credential = ''; visible.value = true }
function edit(r: any) { editing.value = r; form.domain = r.domain; form.credential = r.credential; visible.value = true }
async function save() {
  try {
    if (editing.value) await updateDdnsEntry(editing.value.id, { ...form })
    else await createDdnsEntry({ ...form })
    visible.value = false; fetch()
  } catch { /**/ }
}
async function del(r: any) { try { await deleteDdnsEntry(r.id); fetch() } catch { /**/ } }

const columns = [
  { title: 'Domain', dataIndex: 'domain', key: 'domain' },
  { title: 'IP', dataIndex: 'ip', key: 'ip' },
  { title: 'Credential', dataIndex: 'credential', key: 'credential' },
  { title: 'Enabled', dataIndex: 'enabled', key: 'enabled', width: 80 },
  { title: '', key: 'actions', width: 150 },
]

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">DDNS</h3>
      <AButton type="primary" size="small" @click="add"><span class="i-tabler:plus mr-1" /> Add</AButton>
    </div>
    <ATable :data-source="entries" :loading="loading" :columns="columns" row-key="id" :pagination="false" size="small">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'enabled'">
          <ATag :color="record.enabled ? 'green' : 'red'">{{ record.enabled ? 'Yes' : 'No' }}</ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="edit(record)">{{ t('Edit') }}</AButton>
          <AButton type="link" size="small" danger @click="del(record)">{{ t('Delete') }}</AButton>
        </template>
      </template>
    </ATable>
    <AModal v-model:open="visible" :title="(editing ? 'Edit' : 'Add') + ' DDNS'" @ok="save">
      <AForm layout="vertical">
        <AFormItem label="Domain" required><AInput v-model:value="form.domain" placeholder="home.example.com" /></AFormItem>
        <AFormItem label="Credential"><AInput v-model:value="form.credential" /></AFormItem>
      </AForm>
    </AModal>
  </div>
</template>