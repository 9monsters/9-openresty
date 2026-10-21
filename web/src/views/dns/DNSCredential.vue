<script setup lang="ts">
import { listDnsCredentials, createDnsCredential, updateDnsCredential, deleteDnsCredential } from '@/api/dns'
import { t } from '@/i18n'

const loading = ref(true)
const credentials = ref<any[]>([])
const visible = ref(false)
const editing = ref<any>(null)
const form = reactive({ name: '', provider: 'cloudflare', email: '', accessKeyId: '' })

const columns = [
  { title: 'Name', dataIndex: 'name', key: 'name' },
  { title: 'Provider', dataIndex: 'provider', key: 'provider' },
  { title: t('Status'), dataIndex: 'status', key: 'status', width: 100 },
  { title: t('Actions'), key: 'actions', width: 180 },
]

async function fetch() {
  loading.value = true
  try { const r = await listDnsCredentials(); credentials.value = r.credentials } catch { /**/ }
  loading.value = false
}

function add() { editing.value = null; form.name = ''; form.provider = 'cloudflare'; form.email = ''; form.accessKeyId = ''; visible.value = true }
function edit(r: any) { editing.value = r; form.name = r.name; form.provider = r.provider; form.email = r.email || ''; form.accessKeyId = r.accessKeyId || ''; visible.value = true }
async function save() {
  try {
    if (editing.value) await updateDnsCredential(editing.value.id, { ...form })
    else await createDnsCredential({ ...form })
    visible.value = false; fetch()
  } catch { /**/ }
}
async function del(r: any) { try { await deleteDnsCredential(r.id); fetch() } catch { /**/ } }

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('DNS Credentials') }}</h3>
      <AButton type="primary" size="small" @click="add"><span class="i-tabler:plus mr-1" /> Add</AButton>
    </div>
    <ATable :dataSource="credentials" :columns="columns" :loading="loading" rowKey="id" :pagination="false" size="small">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'"><ATag :color="record.status === 'active' ? 'green' : 'red'">{{ record.status }}</ATag></template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="edit(record)">{{ t('Edit') }}</AButton>
          <AButton type="link" size="small" danger @click="del(record)">{{ t('Delete') }}</AButton>
        </template>
      </template>
    </ATable>
    <AModal v-model:open="visible" :title="(editing ? 'Edit' : 'Add') + ' Credential'" @ok="save">
      <AForm layout="vertical">
        <AFormItem label="Name" required><AInput v-model:value="form.name" /></AFormItem>
        <AFormItem label="Provider"><AInput v-model:value="form.provider" /></AFormItem>
        <AFormItem label="Email"><AInput v-model:value="form.email" /></AFormItem>
        <AFormItem label="API Key"><AInput v-model:value="form.accessKeyId" /></AFormItem>
      </AForm>
    </AModal>
  </div>
</template>