<script setup lang="ts">
import { listCertificates, createCertificate, deleteCertificate } from '@/api/certificates'
import { t } from '@/i18n'
import dayjs from 'dayjs'

const loading = ref(true)
const certs = ref<any[]>([])
const modalVisible = ref(false)
const formData = reactive({ name: '', domains: '', certPEM: '', keyPEM: '', issuer: 'Manual' })

const columns = [
  { title: t('Name') || 'Name', dataIndex: 'name', key: 'name' },
  { title: t('Domains') || 'Domains', key: 'domains' },
  { title: t('Issuer') || 'Issuer', dataIndex: 'issuer', key: 'issuer' },
  { title: t('Expiry') || 'Expiry', dataIndex: 'notAfter', key: 'notAfter', width: 140 },
  { title: t('Auto Renew') || 'Auto Renew', dataIndex: 'autoRenew', key: 'autoRenew', width: 100 },
  { title: t('Status'), key: 'status', width: 100 },
  { title: t('Actions'), key: 'actions', width: 120 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await listCertificates()
    certs.value = res.certificates
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  formData.name = ''
  formData.domains = ''
  formData.certPEM = ''
  formData.keyPEM = ''
  formData.issuer = 'Manual'
  modalVisible.value = true
}

async function handleDelete(cert: any) {
  try {
    await deleteCertificate(cert.id)
    fetchData()
  } catch { /**/ }
}

async function handleSubmit() {
  try {
    await createCertificate({
      name: formData.name,
      domains: formData.domains.split(',').map((d: string) => d.trim()).filter(Boolean),
      certPEM: formData.certPEM,
      keyPEM: formData.keyPEM,
      issuer: formData.issuer,
    })
    modalVisible.value = false
    fetchData()
  } catch { /**/ }
}

function getStatusTag(status: string) {
  switch (status) {
    case 'valid': return { color: 'green', label: t('Valid') || 'Valid' }
    case 'expiring': return { color: 'orange', label: t('Expiring') || 'Expiring' }
    case 'expired': return { color: 'red', label: t('Expired') || 'Expired' }
    default: return { color: 'default', label: status }
  }
}

function formatDate(date: string) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Certificates') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Import') || 'Import' }}
      </AButton>
    </div>

    <ATable :dataSource="certs" :columns="columns" :loading="loading" rowKey="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'domains'">
          <ATag v-for="d in record.domains" :key="d" class="mr-1">{{ d }}</ATag>
        </template>
        <template v-if="column.key === 'notAfter'">
          {{ formatDate(record.notAfter) }}
        </template>
        <template v-if="column.key === 'autoRenew'">
          <span :class="record.autoRenew ? 'i-tabler:check text-green-500' : 'i-tabler:x text-gray-400'" />
        </template>
        <template v-if="column.key === 'status'">
          <ATag :color="getStatusTag(record.status).color">{{ getStatusTag(record.status).label }}</ATag>
        </template>
        <template v-if="column.key === 'actions'">
          <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(record)">
            <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
          </APopconfirm>
        </template>
      </template>
    </ATable>

    <AModal v-model:open="modalVisible" :title="t('Import Certificate') || 'Import Certificate'" @ok="handleSubmit" width="600px">
      <AForm layout="vertical">
        <AFormItem :label="t('Name') || 'Name'" required>
          <AInput v-model:value="formData.name" placeholder="example.com" />
        </AFormItem>
        <AFormItem :label="t('Domains') || 'Domains'" required>
          <AInput v-model:value="formData.domains" placeholder="example.com, *.example.com" />
        </AFormItem>
        <AFormItem label="PEM Certificate">
          <ATextarea v-model:value="formData.certPEM" :rows="4" style="font-family:monospace" placeholder="-----BEGIN CERTIFICATE-----" />
        </AFormItem>
        <AFormItem label="Private Key">
          <ATextarea v-model:value="formData.keyPEM" :rows="4" style="font-family:monospace" placeholder="-----BEGIN PRIVATE KEY-----" />
        </AFormItem>
        <AFormItem :label="t('Issuer') || 'Issuer'">
          <AInput v-model:value="formData.issuer" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>