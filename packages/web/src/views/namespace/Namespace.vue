<script setup lang="ts">
import { listNamespaces, createNamespace, updateNamespace, deleteNamespace } from '@/api/namespaces'
import { t } from '@/i18n'

const loading = ref(true)
const namespaces = ref<any[]>([])
const modalVisible = ref(false)
const editing = ref<any>(null)
const formData = reactive({ name: '', description: '' })

const columns = [
  { title: t('Name') || 'Name', dataIndex: 'name', key: 'name' },
  { title: t('Description') || 'Description', dataIndex: 'description', key: 'description' },
  { title: t('Site Count') || 'Site Count', dataIndex: 'siteCount', key: 'siteCount', width: 120 },
  { title: t('Status'), key: 'enabled', width: 100 },
  { title: t('Actions'), key: 'actions', width: 180 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await listNamespaces()
    namespaces.value = res.namespaces
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  editing.value = null
  formData.name = ''
  formData.description = ''
  modalVisible.value = true
}

function handleEdit(record: any) {
  editing.value = record
  formData.name = record.name
  formData.description = record.description || ''
  modalVisible.value = true
}

async function handleDelete(record: any) {
  try {
    await deleteNamespace(record.id)
    fetchData()
  } catch { /**/ }
}

async function handleSubmit() {
  try {
    if (editing.value) {
      await updateNamespace(editing.value.id, { ...formData })
    } else {
      await createNamespace({ ...formData })
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
      <h3 class="text-xl font-semibold m-0">{{ t('Namespaces') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Add') }}
      </AButton>
    </div>

    <ATable :dataSource="namespaces" :columns="columns" :loading="loading" rowKey="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'enabled'">
          <ATag :color="record.enabled ? 'green' : 'red'">{{ record.enabled ? t('Enabled') : t('Disabled') }}</ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="handleEdit(record)">{{ t('Edit') }}</AButton>
          <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(record)">
            <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
          </APopconfirm>
        </template>
      </template>
    </ATable>

    <AModal v-model:open="modalVisible" :title="(editing ? t('Edit') : t('Add')) + ' Namespace'" @ok="handleSubmit">
      <AForm layout="vertical">
        <AFormItem :label="t('Name') || 'Name'" required>
          <AInput v-model:value="formData.name" placeholder="production" />
        </AFormItem>
        <AFormItem :label="t('Description') || 'Description'">
          <ATextarea v-model:value="formData.description" :rows="3" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>