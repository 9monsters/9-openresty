<script setup lang="ts">
import { getProfile, updateProfile, changePassword } from '@/api/users'
import { useSettingsStore } from '@/stores/settings'
import { t } from '@/i18n'
import type { Locale } from '@/i18n'

const settings = useSettingsStore()
const loading = ref(true)
const saving = ref(false)
const profile = reactive({
  username: '', email: '', role: '', language: 'zh-CN' as string,
})
const pwForm = reactive({ current: '', newPassword: '' })
const pwSaving = ref(false)

async function fetchProfile() {
  loading.value = true
  try {
    const data = await getProfile()
    profile.username = data.username
    profile.email = data.email
    profile.role = data.role
    profile.language = data.language || 'zh-CN'
  } catch { /**/ }
  loading.value = false
}

async function handleSave() {
  saving.value = true
  try {
    await updateProfile({ email: profile.email, language: profile.language })
    settings.setLanguage(profile.language as Locale)
    message.success(t('Operation successful'))
  } catch { message.error(t('Operation failed')) }
  saving.value = false
}

async function handleChangePassword() {
  if (!pwForm.current || !pwForm.newPassword) {
    message.error('Please fill in both fields')
    return
  }
  pwSaving.value = true
  try {
    await changePassword(pwForm.current, pwForm.newPassword)
    message.success(t('Operation successful'))
    pwForm.current = ''
    pwForm.newPassword = ''
  } catch { message.error(t('Operation failed')) }
  pwSaving.value = false
}

onMounted(fetchProfile)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('User Profile') }}</h3>
    </div>

    <ASpin v-if="loading" class="flex-center w-full py-10" />
    <div v-else class="max-w-lg">
      <ACard :title="t('Profile Info') || 'Profile Info'" class="mb-6">
        <AForm layout="vertical">
          <AFormItem :label="t('Username') || 'Username'">
            <AInput :value="profile.username" disabled />
          </AFormItem>
          <AFormItem label="Email">
            <AInput v-model:value="profile.email" />
          </AFormItem>
          <AFormItem :label="t('Role') || 'Role'">
            <AInput :value="profile.role" disabled />
          </AFormItem>
          <AFormItem :label="t('Language')">
            <ASelect v-model:value="profile.language">
              <ASelectOption value="zh-CN">中文</ASelectOption>
              <ASelectOption value="en-US">English</ASelectOption>
            </ASelect>
          </AFormItem>
          <AButton type="primary" :loading="saving" @click="handleSave">{{ t('Save') }}</AButton>
        </AForm>
      </ACard>

      <ACard :title="t('Change Password') || 'Change Password'">
        <AForm layout="vertical">
          <AFormItem label="Current Password">
            <AInputPassword v-model:value="pwForm.current" />
          </AFormItem>
          <AFormItem label="New Password">
            <AInputPassword v-model:value="pwForm.newPassword" />
          </AFormItem>
          <AButton :loading="pwSaving" @click="handleChangePassword">{{ t('Save') }}</AButton>
        </AForm>
      </ACard>
    </div>
  </div>
</template>