<script setup lang="ts">
import { listNodes, createNode, updateNode, deleteNode, testNode } from '@/api/nodes'
import { t } from '@/i18n'

const loading = ref(true)
const nodes = ref<any[]>([])
const modalVisible = ref(false)
const editing = ref<any>(null)
const testing = ref<number | null>(null)
const formData = reactive({ name: '', host: '', port: 9001, environment: 'production' })

const columns = [
  { title: t('Name') || 'Name', dataIndex: 'name', key: 'name' },
  { title: t('Host') || 'Host', dataIndex: 'host', key: 'host' },
  { title: t('Port') || 'Port', dataIndex: 'port', key: 'port', width: 80 },
  { title: t('Status'), key: 'status', width: 100 },
  { title: t('Version') || 'Version', dataIndex: 'version', key: 'version', ellipsis: true },
  { title: t('Uptime') || 'Uptime', dataIndex: 'uptime', key: 'uptime', width: 100 },
  { title: t('Actions'), key: 'actions', width: 250 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await listNodes()
    nodes.value = res.nodes
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  editing.value = null
  formData.name = ''
  formData.host = ''
  formData.port = 9001
  formData.environment = 'production'
  modalVisible.value = true
}

function handleEdit(record: any) {
  editing.value = record
  formData.name = record.name
  formData.host = record.host
  formData.port = record.port
  formData.environment = record.environment || 'production'
  modalVisible.value = true
}

async function handleDelete(record: any) {
  try {
    await deleteNode(record.id)
    fetchData()
  } catch { /**/ }
}

async function handleTest(record: any) {
  testing.value = record.id
  try {
    await testNode(record.id)
    message.success(`${record.name}: ${t('reachable') || 'reachable'}`)
  } catch {
    message.error(`${record.name}: ${t('unreachable') || 'unreachable'}`)
  }
  testing.value = null
}

async function handleSubmit() {
  try {
    if (editing.value) {
      await updateNode(editing.value.id, { ...formData })
    } else {
      await createNode({ ...formData })
    }
    modalVisible.value = false
    fetchData()
  } catch { /**/ }
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Nodes') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Add') }}
      </AButton>
    </div>

    <ATable :dataSource="nodes" :columns="columns" :loading="loading" rowKey="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <ATag :color="record.status === 'online' ? 'green' : 'red'">
            <template #icon>
              <span :class="record.status === 'online' ? 'i-tabler:circle-check' : 'i-tabler:circle-x'" />
            </template>
            {{ record.status === 'online' ? t('Online') || 'Online' : t('Offline') || 'Offline' }}
          </ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="handleTest(record)" :loading="testing === record.id">
            {{ t('Test Connection') }}
          </AButton>
          <AButton type="link" size="small" @click="handleEdit(record)">{{ t('Edit') }}</AButton>
          <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(record)">
            <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
          </APopconfirm>
        </template>
      </template>
    </ATable>

    <AModal v-model:open="modalVisible" :title="(editing ? t('Edit') : t('Add')) + ' Node'" @ok="handleSubmit">
      <AForm layout="vertical">
        <AFormItem :label="t('Name') || 'Name'" required>
          <AInput v-model:value="formData.name" placeholder="prod-lb-01" />
        </AFormItem>
        <AFormItem :label="t('Host') || 'Host'" required>
          <AInput v-model:value="formData.host" placeholder="192.168.1.10" />
        </AFormItem>
        <AFormItem :label="t('Port') || 'Port'">
          <AInputNumber v-model:value="formData.port" :min="1" :max="65535" class="w-full" />
        </AFormItem>
        <AFormItem :label="t('Environment') || 'Environment'">
          <AInput v-model:value="formData.environment" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>