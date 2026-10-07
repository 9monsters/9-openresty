<script setup lang="ts">
import { getSettings, updateSettings, getServerName } from '@/api/settings'
import { getNginxStatus, reloadNginx, testConfig } from '@/api/nginx'

const settings = reactive({
  nginxBin: '',
  nginxConfDir: '',
  nginxConfPath: '',
  nginxLogDir: '',
  nginxPIDPath: '',
  serverName: '',
  theme: 'light',
})

const loading = ref(true)
const saving = ref(false)
const nginxInfo = ref<{ status: string; version: string; isOpenResty: boolean } | null>(null)

onMounted(async () => {
  try {
    const [data, nameData, nginx] = await Promise.all([
      getSettings(),
      getServerName(),
      getNginxStatus(),
    ])
    Object.assign(settings, data)
    settings.serverName = nameData.name
    nginxInfo.value = nginx
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})

async function handleSave() {
  saving.value = true
  try {
    await updateSettings(settings)
    AMessage.success('设置已保存')
  } catch (e: any) {
    AMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

async function handleReload() {
  try {
    const r = await reloadNginx()
    AMessage.success(r.ok ? '重载成功' : r.output)
  } catch (e: any) {
    AMessage.error(e.message)
  }
}

async function handleTest() {
  try {
    const r = await testConfig()
    AMessage.info(r.ok ? '✅ 配置测试通过' : '❌ ' + r.output)
  } catch (e: any) {
    AMessage.error(e.message)
  }
}
</script>

<template>
  <div v-if="loading" class="flex-center py-20"><ASpin size="large" /></div>

  <div v-else class="max-w-3xl">
    <h3 class="text-xl font-semibold mb-6">系统设置</h3>

    <ACard title="OpenResty/Nginx 路径" class="mb-4">
      <AForm layout="vertical">
        <AFormItem label="Nginx 二进制路径">
          <AInput v-model:value="settings.nginxBin" placeholder="/usr/local/openresty/nginx/sbin/nginx" />
        </AFormItem>
        <AFormItem label="配置文件目录">
          <AInput v-model:value="settings.nginxConfDir" placeholder="/usr/local/openresty/nginx/conf" />
        </AFormItem>
        <AFormItem label="主配置文件">
          <AInput v-model:value="settings.nginxConfPath" placeholder="/usr/local/openresty/nginx/conf/nginx.conf" />
        </AFormItem>
        <AFormItem label="日志目录">
          <AInput v-model:value="settings.nginxLogDir" placeholder="/usr/local/openresty/nginx/logs" />
        </AFormItem>
        <AFormItem label="PID 文件路径">
          <AInput v-model:value="settings.nginxPIDPath" placeholder="/usr/local/openresty/nginx/logs/nginx.pid" />
        </AFormItem>
      </AForm>
    </ACard>

    <ACard title="服务器信息" class="mb-4">
      <div class="grid grid-cols-2 gap-4">
        <div class="text-gray-500">服务名称</div>
        <div class="font-medium">{{ settings.serverName }}</div>
        <div class="text-gray-500">Nginx 版本</div>
        <div class="font-medium text-sm">{{ nginxInfo?.version }}</div>
        <div class="text-gray-500">OpenResty</div>
        <div>
          <ATag :color="nginxInfo?.isOpenResty ? 'green' : 'default'">
            {{ nginxInfo?.isOpenResty ? '是' : '否' }}
          </ATag>
        </div>
        <div class="text-gray-500">当前状态</div>
        <div>
          <ATag :color="nginxInfo?.status === 'running' ? 'green' : 'red'">
            {{ nginxInfo?.status === 'running' ? '运行中' : '已停止' }}
          </ATag>
        </div>
      </div>
    </ACard>

    <ACard title="Nginx 控制" class="mb-4">
      <ASpace>
        <AButton @click="handleTest">
          <span class="i-tabler:test-pipe mr-1" /> 测试配置 (nginx -t)
        </AButton>
        <AButton type="primary" @click="handleReload">
          <span class="i-tabler:refresh mr-1" /> 重载 (nginx -s reload)
        </AButton>
        <AButton danger>
          <span class="i-tabler:power mr-1" /> 重启
        </AButton>
      </ASpace>
    </ACard>

    <!-- Save button -->
    <div class="flex justify-end">
      <AButton type="primary" size="large" :loading="saving" @click="handleSave">
        <span class="i-tabler:device-floppy mr-1" /> 保存设置
      </AButton>
    </div>
  </div>
</template>
