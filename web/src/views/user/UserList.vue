<script setup lang="ts">
import { listUsers, createUser, updateUser, deleteUser } from '@/api/users'
import { t } from '@/i18n'
import dayjs from 'dayjs'

const loading = ref(true)
const users = ref<any[]>([])
const modalVisible = ref(false)
const editing = ref<any>(null)
const formData = reactive({ username: '', email: '', password: '', role: 'Operator' })

const columns = [
  { title: t('Username') || 'Username', dataIndex: 'username', key: 'username' },
  { title: 'Email', dataIndex: 'email', key: 'email' },
  { title: t('Role') || 'Role', dataIndex: 'role', key: 'role' },
  { title: t('Status'), key: 'status', width: 100 },
  { title: t('Last Login') || 'Last Login', dataIndex: 'lastLogin', key: 'lastLogin', width: 180 },
  { title: t('Actions'), key: 'actions', width: 180 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await listUsers()
    users.value = res.users
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  editing.value = null
  formData.username = ''
  formData.email = ''
  formData.password = ''
  formData.role = 'Operator'
  modalVisible.value = true
}

function handleEdit(record: any) {
  editing.value = record
  formData.username = record.username
  formData.email = record.email
  formData.password = ''
  formData.role = record.role
  modalVisible.value = true
}

async function handleDelete(record: any) {
  try {
    await deleteUser(record.id)
    fetchData()
  } catch { /**/ }
}

async function handleSubmit() {
  try {
    if (editing.value) {
      await updateUser(editing.value.id, { username: formData.username, email: formData.email, role: formData.role })
    } else {
      if (!formData.password) {
        message.error('Password is required')
        return
      }
      await createUser({ ...formData })
    }
    modalVisible.value = false
    fetchData()
  } catch { /**/ }
}

function formatDate(date: string) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Manage Users') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Add') }}
      </AButton>
    </div>

    <ATable :dataSource="users" :columns="columns" :loading="loading" rowKey="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <ATag :color="record.status === 'active' ? 'green' : 'red'">{{ record.status }}</ATag>
        </template>
        <template v-if="column.key === 'lastLogin'">
          {{ formatDate(record.lastLogin) }}
        </template>
        <template v-if="column.key === 'actions'">
          <AButton type="link" size="small" @click="handleEdit(record)">{{ t('Edit') }}</AButton>
          <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(record)">
            <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
          </APopconfirm>
        </template>
      </template>
    </ATable>

    <AModal v-model:open="modalVisible" :title="(editing ? t('Edit') : t('Add')) + ' ' + (t('User') || 'User')" @ok="handleSubmit">
      <AForm layout="vertical">
        <AFormItem :label="t('Username') || 'Username'" required>
          <AInput v-model:value="formData.username" placeholder="username" />
        </AFormItem>
        <AFormItem label="Email" required>
          <AInput v-model:value="formData.email" placeholder="user@example.com" />
        </AFormItem>
        <AFormItem label="Password" v-if="!editing" required>
          <AInputPassword v-model:value="formData.password" placeholder="password" />
        </AFormItem>
        <AFormItem :label="t('Role') || 'Role'">
          <ASelect v-model:value="formData.role">
            <ASelectOption value="Administrator">Administrator</ASelectOption>
            <ASelectOption value="Operator">Operator</ASelectOption>
            <ASelectOption value="Viewer">Viewer</ASelectOption>
          </ASelect>
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>