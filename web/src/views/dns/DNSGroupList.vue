<script setup lang="ts">
import { listDnsGroups, createDnsGroup, updateDnsGroup, deleteDnsGroup } from '@/api/dns'
import { t } from '@/i18n'

const loading = ref(true)
const groups = ref<any[]>([])
const visible = ref(false)
const editing = ref<any>(null)
const form = reactive({ name: '', strategy: 'round_robin' })

async function fetch() {
  loading.value = true
  try { const r = await listDnsGroups(); groups.value = r.groups } catch { /**/ }
  loading.value = false
}
function add() { editing.value = null; form.name = ''; form.strategy = 'round_robin'; visible.value = true }
function edit(r: any) { editing.value = r; form.name = r.name; form.strategy = r.strategy; visible.value = true }
async function save() {
  try {
    if (editing.value) await updateDnsGroup(editing.value.id, { ...form })
    else await createDnsGroup({ ...form })
    visible.value = false; fetch()
  } catch { /**/ }
}
async function del(r: any) { try { await deleteDnsGroup(r.id); fetch() } catch { /**/ } }

const columns = [
  { title: 'Name', dataIndex: 'name', key: 'name' },
  { title: 'Strategy', dataIndex: 'strategy', key: 'strategy' },
  { title: 'Domains', key: 'domains' },
  { title: '', key: 'actions', width: 150 },
]

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Groups') }}</h3>
      <AButton type="primary" size="small" @click="add"><span class="i-tabler:plus mr-1" /> Add</AButton>
    </div>
    <ATable :data-source="groups" :loading="loading" :columns="columns" row-key="id" :pagination="false" size="small">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'strategy'">
          <ATag>{{ record.strategy }}</ATag>
        </template>
        <template v-if="column.key === 'domains'">
          <ATag v-for="d in record.domains" :key="d" class="mr-1">{{ d }}</ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="edit(record)">{{ t('Edit') }}</AButton>
          <AButton type="link" size="small" danger @click="del(record)">{{ t('Delete') }}</AButton>
        </template>
      </template>
    </ATable>
    <AModal v-model:open="visible" :title="(editing ? 'Edit' : 'Add') + ' Group'" @ok="save">
      <AForm layout="vertical">
        <AFormItem label="Name" required><AInput v-model:value="form.name" /></AFormItem>
        <AFormItem label="Strategy">
          <ASelect v-model:value="form.strategy">
            <ASelectOption value="round_robin">Round Robin</ASelectOption>
            <ASelectOption value="weighted">Weighted</ASelectOption>
            <ASelectOption value="geo">Geographic</ASelectOption>
          </ASelect>
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>