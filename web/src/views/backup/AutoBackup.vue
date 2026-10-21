<script setup lang="ts">
import { t } from '@/i18n'
import { getAutoBackup, updateAutoBackup } from '@/api/backup'
import dayjs from 'dayjs'

const loading = ref(true)
const saving = ref(false)
const settings = reactive({ enabled: false, interval: 'daily', retention: 7, time: '03:00', path: '/var/backups/openresty' })
const schedule = ref<any>(null)

onMounted(async () => {
  try {
    const res = await getAutoBackup()
    Object.assign(settings, res.settings)
    schedule.value = res.schedule
  } catch { /**/ }
  loading.value = false
})

async function handleSave() {
  saving.value = true
  try {
    await updateAutoBackup({ ...settings })
    message.success(t('Operation successful'))
  } catch { message.error(t('Operation failed')) }
  saving.value = false
}
</script>

<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">{{ t('Auto Backup') }}</h3>
    <ASpin v-if="loading" class="flex-center w-full py-10" />
    <ACard v-else>
      <AForm layout="vertical" class="max-w-lg">
        <AFormItem :label="t('Enabled')">
          <ASwitch v-model:checked="settings.enabled" />
        </AFormItem>
        <AFormItem :label="t('Interval') || 'Interval'">
          <ASelect v-model:value="settings.interval">
            <ASelectOption value="hourly">Hourly</ASelectOption>
            <ASelectOption value="daily">Daily</ASelectOption>
            <ASelectOption value="weekly">Weekly</ASelectOption>
          </ASelect>
        </AFormItem>
        <AFormItem :label="t('Retention') || 'Retention (days)'">
          <AInputNumber v-model:value="settings.retention" :min="1" :max="90" class="w-full" />
        </AFormItem>
        <AFormItem :label="t('Time') || 'Time'">
          <AInput v-model:value="settings.time" placeholder="03:00" />
        </AFormItem>
        <AFormItem :label="t('Path') || 'Path'">
          <AInput v-model:value="settings.path" placeholder="/var/backups/openresty" />
        </AFormItem>
        <AButton type="primary" :loading="saving" @click="handleSave">{{ t('Save') }}</AButton>
      </AForm>
      <div v-if="schedule" class="mt-4 text-sm text-gray-500">
        <div>Next Run: {{ dayjs(schedule.nextRun).format('YYYY-MM-DD HH:mm:ss') }}</div>
        <div>Last Run: {{ dayjs(schedule.lastRun).format('YYYY-MM-DD HH:mm:ss') }}</div>
      </div>
    </ACard>
  </div>
</template>