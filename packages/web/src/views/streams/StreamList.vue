<script setup lang="ts">
import { listStreams, createStream, updateStream, deleteStream } from '@/api/streams'
import { t } from '@/i18n'

const loading = ref(true)
const streams = ref<any[]>([])
const modalVisible = ref(false)
const editing = ref<any>(null)
const formData = reactive({ name: '', listen: '', proxyPass: '', protocol: 'tcp' })

const columns = [
  { title: t('Name') || 'Name', dataIndex: 'name', key: 'name' },
  { title: 'Listen', dataIndex: 'listen', key: 'listen', width: 100 },
  { title: 'Proxy Pass', dataIndex: 'proxyPass', key: 'proxyPass' },
  { title: t('Protocol') || 'Protocol', dataIndex: 'protocol', key: 'protocol', width: 80 },
  { title: t('Status'), key: 'status', width: 100 },
  { title: t('Actions'), key: 'actions', width: 200 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await listStreams()
    streams.value = res.streams
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  editing.value = null
  formData.name = ''
  formData.listen = ''
  formData.proxyPass = ''
  formData.protocol = 'tcp'
  modalVisible.value = true
}

function handleEdit(record: any) {
  editing.value = record
  formData.name = record.name
  formData.listen = record.listen
  formData.proxyPass = record.proxyPass
  formData.protocol = record.protocol || 'tcp'
  modalVisible.value = true
}

async function handleDelete(record: any) {
  try {
    await deleteStream(record.id)
    fetchData()
  } catch { /**/ }
}

async function handleSubmit() {
  try {
    if (editing.value) {
      await updateStream(editing.value.id, { ...formData })
    } else {
      await createStream({ ...formData })
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
      <h3 class="text-xl font-semibold m-0">{{ t('Manage Streams') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Add') }}
      </AButton>
    </div>

    <ATable :dataSource="streams" :columns="columns" :loading="loading" rowKey="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'protocol'">
          <ATag :color="record.protocol === 'tcp' ? 'blue' : 'purple'">{{ record.protocol.toUpperCase() }}</ATag>
        </template>
        <template v-if="column.key === 'status'">
          <ATag :color="record.status === 'active' ? 'green' : 'red'">{{ record.status }}</ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="handleEdit(record)">{{ t('Edit') }}</AButton>
          <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(record)">
            <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
          </APopconfirm>
        </template>
      </template>
    </ATable>

    <AModal v-model:open="modalVisible" :title="(editing ? t('Edit') : t('Add')) + ' Stream'" @ok="handleSubmit">
      <AForm layout="vertical">
        <AFormItem :label="t('Name') || 'Name'" required>
          <AInput v-model:value="formData.name" placeholder="mysql_proxy" />
        </AFormItem>
        <AFormItem label="Listen Port" required>
          <AInput v-model:value="formData.listen" placeholder="3306" />
        </AFormItem>
        <AFormItem label="Proxy Pass" required>
          <AInput v-model:value="formData.proxyPass" placeholder="10.0.0.1:3306" />
        </AFormItem>
        <AFormItem :label="t('Protocol') || 'Protocol'">
          <ASelect v-model:value="formData.protocol">
            <ASelectOption value="tcp">TCP</ASelectOption>
            <ASelectOption value="udp">UDP</ASelectOption>
          </ASelect>
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>