<script setup lang="ts">
import { login as loginApi } from '@/api/auth'
import { useUserStore } from '@/stores/user'
import { useSettingsStore } from '@/stores/settings'
import { t } from '@/i18n'
import { UserOutlined, LockOutlined } from '@antdv-next/icons'

const userStore = useUserStore()
const settings = useSettingsStore()
const router = useRouter()

const formRef = ref()
const form = reactive({ username: 'admin', password: '' })
const loading = ref(false)
const error = ref('')

// Server name is dynamic once logged in, but the login page can't call
// authenticated APIs. Use fallback; the settings store will update it later.
const serverName = computed(() => settings.serverName)

async function handleLogin() {
  if (!form.username || !form.password) {
    error.value = t('Username and password are required', '请输入用户名和密码')
    return
  }
  loading.value = true
  error.value = ''
  try {
    const res = await loginApi({ username: form.username, password: form.password })
    userStore.login(res.token, { id: res.user.id, username: res.user.username, email: '', role: res.user.role, status: 'active' })
    message.success(t('Login successful', '登录成功'))
    router.push('/')
  } catch (e: any) {
    error.value = e.message || t('Login failed', '登录失败')
  } finally {
    loading.value = false
  }
}

function switchLang() {
  settings.setLanguage(settings.language === 'zh-CN' ? 'en-US' : 'zh-CN')
}
</script>

<template>
  <div class="login-page">
    <div class="login-form">
      <!-- Title -->
      <div class="project-title">
        <h1>{{ serverName }}</h1>
        <p class="ios-footnote">{{ t('Sign in to continue', '登录以继续') }}</p>
      </div>

      <!-- Form -->
      <AForm ref="formRef" :model="form" layout="vertical" @keyup.enter="handleLogin">
        <AFormItem
          name="username"
          :rules="[{ required: true, message: t('Please input your username!', '请输入用户名！') }]"
        >
          <AInput
            v-model:value="form.username"
            size="large"
            :placeholder="t('Username', '用户名')"
            @pressEnter="handleLogin"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </AInput>
        </AFormItem>

        <AFormItem
          name="password"
          :rules="[{ required: true, message: t('Please input your password!', '请输入密码！') }]"
        >
          <AInput
            v-model:value="form.password"
            type="password"
            size="large"
            :placeholder="t('Password', '密码')"
            @pressEnter="handleLogin"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </AInput>
        </AFormItem>
      </AForm>

      <!-- Error message -->
      <p v-if="error" class="login-error ios-footnote" data-state="bad">{{ error }}</p>

      <!-- Submit -->
      <ASpin :spinning="loading">
        <AButton
          type="primary"
          size="large"
          block
          :loading="loading"
          :disabled="loading"
          @click="handleLogin"
        >
          {{ t('Sign In', '登录') }}
        </AButton>
      </ASpin>

      <!-- Footer -->
      <div class="footer">
        <button class="lang-btn ios-caption" @click="switchLang">
          {{ settings.language === 'zh-CN' ? 'English' : '中文' }}
        </button>
        <p class="ios-caption">{{ serverName }} · © {{ new Date().getFullYear() }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;
}

.login-form {
  width: 80%;
  max-width: 420px;
}

.project-title {
  margin: 48px 0;
  text-align: center;
}

.project-title h1 {
  font-size: 50px;
  font-weight: 100;
  line-height: 1.1;
  letter-spacing: -0.03em;
  color: var(--label);
  margin: 0 0 8px;
}

.login-error[data-state='bad'] {
  color: var(--bad);
  text-align: center;
  margin: 0 0 12px;
}

.footer {
  text-align: center;
  padding: 32px 20px;
}

.footer p {
  margin: 8px 0 0;
}

.lang-btn {
  background: none;
  border: 0;
  padding: 4px 12px;
  border-radius: var(--r-pill);
  font-size: 12px;
  font-weight: 600;
  color: var(--label-2);
  cursor: pointer;
  transition: background var(--dur) var(--ease-ios);
}

.lang-btn:hover {
  background: var(--fill);
}
</style>