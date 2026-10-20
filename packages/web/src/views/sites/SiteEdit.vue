<script setup lang="ts">
import { createSite, getSite, updateSite } from '@/api/sites'
import type { ParsedSite } from '@/types'

const router = useRouter()
const route = useRoute()
const isEdit = computed(() => !!route.params.name)
const loading = ref(false)

// Form
const serverName = ref('')
const listen80 = ref(true)
const listen443 = ref(false)
const root = ref('')
const locations = ref<{ path: string; proxyPass: string; root: string }[]>([])
const luaBlocks = ref<{ type: string; code: string }[]>([])

const LUA_BLOCK_TYPES = [
  'access_by_lua_block',
  'content_by_lua_block',
  'rewrite_by_lua_block',
  'log_by_lua_block',
  'header_filter_by_lua_block',
  'body_filter_by_lua_block',
]

function addLocation() {
  locations.value.push({ path: '/', proxyPass: '', root: '' })
}

function removeLocation(idx: number) {
  locations.value.splice(idx, 1)
}

function addLuaBlock() {
  luaBlocks.value.push({ type: 'content_by_lua_block', code: '' })
}

function removeLuaBlock(idx: number) {
  luaBlocks.value.splice(idx, 1)
}

async function handleSubmit() {
  loading.value = true
  try {
    const listen = [listen80.value ? '80' : '', listen443.value ? '443 ssl http2' : ''].filter(Boolean)
    
    if (isEdit.value) {
      // Build full config string for update
      const content = buildConfig()
      await updateSite(route.params.name as string, content)
    } else {
      await createSite({
        serverName: serverName.value.split(/[,\s]+/).filter(Boolean),
        listen,
        root: root.value || undefined,
        ssl: listen443.value,
        luaDirectives: luaBlocks.value.filter(l => l.code.trim()).map(l => ({
          type: l.type,
          value: l.code,
        })),
        locations: locations.value.filter(l => l.path.trim()).map(l => ({
          path: l.path,
          proxyPass: l.proxyPass || undefined,
          root: l.root || undefined,
        })),
      })
    }
    message.success(isEdit.value ? '站点已更新' : '站点已创建')
    router.push('/sites')
  } catch (e: any) {
    message.error(e.message || '操作失败')
  } finally {
    loading.value = false
  }
}

function buildConfig(): string {
  const lines: string[] = []
  lines.push('server {')
  const listen = [listen80.value ? '80' : '', listen443.value ? '443 ssl http2' : ''].filter(Boolean)
  for (const l of listen) lines.push(`  listen ${l};`)
  lines.push(`  server_name ${serverName.value};`)
  if (root.value) lines.push(`  root ${root.value};`)
  
  for (const loc of locations.value) {
    if (!loc.path.trim()) continue
    lines.push(`  location ${loc.path} {`)
    if (loc.proxyPass) lines.push(`    proxy_pass ${loc.proxyPass};`)
    if (loc.root) lines.push(`    root ${loc.root};`)
    lines.push('  }')
  }

  for (const lb of luaBlocks.value) {
    if (!lb.code.trim()) continue
    lines.push(`  ${lb.type} {`)
    lines.push(lb.code)
    lines.push('  }')
  }

  lines.push('}')
  return lines.join('\n')
}

// Load existing site data
onMounted(async () => {
  if (isEdit.value) {
    try {
      const site = await getSite(route.params.name as string)
      serverName.value = site.serverName.join(' ')
      listen80.value = site.listen.some(l => l.includes('80'))
      listen443.value = site.listen.some(l => l.includes('443') || l.includes('ssl') || site.ssl)
      root.value = site.root || ''
      locations.value = site.locations.map(l => ({
        path: l.path,
        proxyPass: l.proxyPass || '',
        root: l.root || '',
      }))
      if (site.luaDirectives.length > 0) {
        luaBlocks.value = site.luaDirectives.map(d => ({
          type: d.type,
          code: d.value,
        }))
      }
    } catch {} 
  }
})
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ isEdit ? '编辑站点' : '添加站点' }}</h3>
    </div>

    <AForm layout="vertical" class="max-w-3xl">
      <!-- Server Name -->
      <AFormItem label="域名 (Server Name)" required>
        <AInput v-model:value="serverName" placeholder="example.com www.example.com" />
      </AFormItem>

      <!-- Listen -->
      <AFormItem label="监听端口">
        <ASpace>
          <ACheckbox v-model:checked="listen80">HTTP (80)</ACheckbox>
          <ACheckbox v-model:checked="listen443">HTTPS (443)</ACheckbox>
        </ASpace>
      </AFormItem>

      <!-- Root -->
      <AFormItem label="网站根目录 (Root)">
        <AInput v-model:value="root" placeholder="/usr/local/openresty/nginx/html" />
      </AFormItem>

      <!-- Locations -->
      <div class="mb-4">
        <div class="flex-between mb-2">
          <div class="font-medium">Location 配置</div>
          <AButton size="small" @click="addLocation">
            <span class="i-tabler:plus" /> 添加 Location
          </AButton>
        </div>
        <ACard v-for="(loc, idx) in locations" :key="idx" size="small" class="mb-2">
          <div class="flex gap-2 items-start">
            <AInput v-model:value="loc.path" placeholder="路径 e.g. /api/" class="w-30" />
            <AInput v-model:value="loc.proxyPass" placeholder="proxy_pass" class="flex-1" />
            <AInput v-model:value="loc.root" placeholder="root" class="flex-1" />
            <AButton size="small" danger @click="removeLocation(idx)">
              <span class="i-tabler:trash" />
            </AButton>
          </div>
        </ACard>
      </div>

      <!-- Lua Blocks (OpenResty specific) -->
      <div class="mb-4">
        <div class="flex-between mb-2">
          <div class="font-medium">Lua 代码块 <ATag color="purple">OpenResty</ATag></div>
          <AButton size="small" @click="addLuaBlock">
            <span class="i-tabler:plus" /> 添加 Lua 块
          </AButton>
        </div>
        <ACard v-for="(lb, idx) in luaBlocks" :key="idx" size="small" class="mb-2">
          <div class="flex gap-2 mb-2 items-center">
            <ASelect v-model:value="lb.type" class="w-50">
              <ASelectOption v-for="t in LUA_BLOCK_TYPES" :key="t" :value="t">{{ t }}</ASelectOption>
            </ASelect>
            <AButton size="small" danger @click="removeLuaBlock(idx)">
              <span class="i-tabler:trash" />
            </AButton>
          </div>
          <ATextarea
            v-model:value="lb.code"
            :rows="6"
            placeholder="ngx.say('Hello OpenResty!')"
            class="font-mono text-sm"
          />
        </ACard>
      </div>

      <!-- Submit -->
      <AFormItem>
        <ASpace>
          <AButton type="primary" :loading="loading" @click="handleSubmit">
            {{ isEdit ? '更新站点' : '创建站点' }}
          </AButton>
          <AButton @click="router.back()">取消</AButton>
        </ASpace>
      </AFormItem>
    </AForm>
  </div>
</template>
