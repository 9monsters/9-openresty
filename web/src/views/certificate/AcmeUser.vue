<script setup lang="ts">
import { listAcmeUsers, createAcmeUser, updateAcmeUser, deleteAcmeUser } from '@/api/acme'
import { t } from '@/i18n'

const loading = ref(true)
const users = ref<any[]>([])
const modalVisible = ref(false)
const editing = ref<any>(null)
const formData = reactive({ email: '', provider: "Let's Encrypt", server: 'https://acme-v02.api.letsencrypt.org/directory' })

const columns = [
  { title: 'Email', dataIndex: 'email', key: 'email' },
  { title: t('Provider') || 'Provider', dataIndex: 'provider', key: 'provider' },
  { title: 'Server', dataIndex: 'server', key: 'server', ellipsis: true },
  { title: t('Certificates'), dataIndex: 'certCount', key: 'certCount', width: 100 },
  { title: t('Status'), dataIndex: 'status', key: 'status', width: 100 },
  { title: t('Actions'), key: 'actions', width: 180 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await listAcmeUsers()
    users.value = res.users
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  editing.value = null
  formData.email = ''
  formData.provider = "Let's Encrypt"
  formData.server = 'https://acme-v02.api.letsencrypt.org/directory'
  modalVisible.value = true
}

function handleEdit(user: any) {
  editing.value = user
  formData.email = user.email
  formData.provider = user.provider
  formData.server = user.server
  modalVisible.value = true
}

async function handleDelete(user: any) {
  try {
    await deleteAcmeUser(user.id)
    fetchData()
  } catch { /**/ }
}

async function handleSubmit() {
  try {
    if (editing.value) {
      await updateAcmeUser(editing.value.id, { ...formData })
    } else {
      await createAcmeUser({ ...formData })
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
      <h3 class="text-xl font-semibold m-0">{{ t('ACME Users') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Add') }}
      </AButton>
    </div>

    <ATable :dataSource="users" :columns="columns" :loading="loading" rowKey="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <ATag :color="record.status === 'active' ? 'green' : 'red'">{{ record.status }}</ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="handleEdit(record)">{{ t('Edit') }}</AButton>
          <AButton type="link" size="small" danger @click="handleDelete(record)">{{ t('Delete') }}</AButton>
        </template>
      </template>
    </ATable>

    <AModal v-model:open="modalVisible" :title="editing ? t('Edit') : t('Add') + ' ACME ' + (t('User') || 'User')" @ok="handleSubmit">
      <AForm layout="vertical">
        <AFormItem label="Email" required>
          <AInput v-model:value="formData.email" placeholder="admin@example.com" />
        </AFormItem>
        <AFormItem :label="t('Provider') || 'Provider'">
          <AInput v-model:value="formData.provider" />
        </AFormItem>
        <AFormItem label="Server URL">
          <AInput v-model:value="formData.server" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>