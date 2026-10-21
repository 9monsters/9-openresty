<script setup lang="ts">
import { listBackups, createBackup, restoreBackup, deleteBackup, getAutoBackup, updateAutoBackup } from '@/api/backup'
import { t } from '@/i18n'
import dayjs from 'dayjs'

const loading = ref(true)
const backupLoading = ref(false)
const backups = ref<any[]>([])
const autoSettings = ref<any>({ enabled: false, interval: 'daily', retention: 7, time: '03:00', path: '/var/backups/openresty' })
const schedule = ref<any>(null)
const savingSettings = ref(false)

const columns = [
  { title: t('Name') || 'Name', dataIndex: 'name', key: 'name' },
  { title: t('Type') || 'Type', dataIndex: 'type', key: 'type', width: 80 },
  { title: t('Size') || 'Size', key: 'size', width: 100 },
  { title: t('Files') || 'Files', dataIndex: 'files', key: 'files', width: 80 },
  { title: t('Time') || 'Time', dataIndex: 'timestamp', key: 'timestamp' },
  { title: t('Actions'), key: 'actions', width: 200 },
]

async function fetchData() {
  loading.value = true
  try {
    const [bkpRes, autoRes] = await Promise.all([listBackups(), getAutoBackup()])
    backups.value = bkpRes.backups
    autoSettings.value = autoRes.settings
    schedule.value = autoRes.schedule
  } catch { /**/ }
  loading.value = false
}

async function handleCreateBackup() {
  backupLoading.value = true
  try {
    await createBackup()
    message.success(t('Operation successful'))
    fetchData()
  } catch { message.error(t('Operation failed')) }
  backupLoading.value = false
}

async function handleRestore(backup: any) {
  try {
    await restoreBackup(backup.id)
    message.success(`Restored from ${backup.name}`)
  } catch { message.error(t('Operation failed')) }
}

async function handleDelete(backup: any) {
  try {
    await deleteBackup(backup.id)
    fetchData()
  } catch { /**/ }
}

async function handleSaveAutoBackup() {
  savingSettings.value = true
  try {
    await updateAutoBackup(autoSettings.value)
    message.success(t('Operation successful'))
  } catch { message.error(t('Operation failed')) }
  savingSettings.value = false
}

function formatSize(bytes: number) {
  return `${(bytes / 1024).toFixed(0)} KB`
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Backup') }}</h3>
      <AButton type="primary" :loading="backupLoading" @click="handleCreateBackup">
        <span class="i-tabler:upload mr-1" /> {{ t('Create Backup') || 'Create Backup' }}
      </AButton>
    </div>

    <ASpin v-if="loading" class="flex-center w-full py-10" />

    <div v-else class="grid grid-cols-2 gap-6">
      <ACard :title="t('Backups') || 'Backups'" class="mb-4">
        <ATable :dataSource="backups" :columns="columns" :pagination="false" size="small" rowKey="id">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'type'">
              <ATag :color="record.type === 'auto' ? 'blue' : 'green'">{{ record.type }}</ATag>
            </template>
            <template v-if="column.key === 'size'">{{ formatSize(record.size) }}</template>
            <template v-if="column.key === 'timestamp'">{{ formatTime(record.timestamp) }}</template>
            <template v-if="column.key === 'actions'">
              <AButton type="link" size="small" @click="handleRestore(record)">{{ t('Restore') || 'Restore' }}</AButton>
              <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(record)">
                <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
              </APopconfirm>
            </template>
          </template>
        </ATable>
      </ACard>

      <!-- Auto backup settings -->
      <ACard :title="t('Auto Backup')" class="mb-4">
        <AForm layout="vertical">
          <AFormItem :label="t('Enabled')">
            <ASwitch v-model:checked="autoSettings.enabled" />
          </AFormItem>
          <AFormItem :label="t('Interval') || 'Interval'">
            <ASelect v-model:value="autoSettings.interval">
              <ASelectOption value="hourly">Hourly</ASelectOption>
              <ASelectOption value="daily">Daily</ASelectOption>
              <ASelectOption value="weekly">Weekly</ASelectOption>
            </ASelect>
          </AFormItem>
          <AFormItem :label="t('Retention') || 'Retention'">
            <AInputNumber v-model:value="autoSettings.retention" :min="1" :max="90" class="w-full" />
          </AFormItem>
          <AFormItem :label="t('Time') || 'Time'">
            <AInput v-model:value="autoSettings.time" />
          </AFormItem>
          <AFormItem :label="t('Path') || 'Path'">
            <AInput v-model:value="autoSettings.path" />
          </AFormItem>
          <AButton type="primary" :loading="savingSettings" @click="handleSaveAutoBackup">{{ t('Save') }}</AButton>
        </AForm>
        <div v-if="schedule" class="mt-4 text-sm text-gray-500">
          <div>{{ t('Next Run') || 'Next Run' }}: {{ formatTime(schedule.nextRun) }}</div>
          <div>{{ t('Last Run') || 'Last Run' }}: {{ formatTime(schedule.lastRun) }}</div>
          <div>{{ t('Last Status') || 'Last Status' }}:
            <ATag :color="schedule.lastStatus === 'success' ? 'green' : 'red'">{{ schedule.lastStatus }}</ATag>
          </div>
        </div>
      </ACard>
    </div>
  </div>
</template>