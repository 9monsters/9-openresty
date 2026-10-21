<script setup lang="ts">
import { listLuaBlocks, listLuaLibraries, listSharedDicts } from '@/api/lua'

interface LuaBlock {
  filePath: string
  siteName: string
  directives: { type: string; value: string }[]
}

interface Library {
  name: string
  description: string
  installed: boolean
}

interface SharedDict {
  name: string
  size: string
  usage: string
}

const luaBlocks = ref<LuaBlock[]>([])
const libraries = ref<Library[]>([])
const sharedDicts = ref<SharedDict[]>([])
const selectedBlock = ref<LuaBlock | null>(null)

onMounted(async () => {
  try {
    const [blocks, libs, dicts] = await Promise.all([
      listLuaBlocks(),
      listLuaLibraries(),
      listSharedDicts(),
    ])
    luaBlocks.value = blocks.luaBlocks
    libraries.value = libs.libraries
    sharedDicts.value = dicts.dicts
  } catch (e) { console.error(e) }
})
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Left: Lua Blocks -->
    <div class="lg:col-span-2">
      <ACard title="Lua 代码块" class="mb-4">
        <template #extra>
          <ATag color="purple">{{ luaBlocks.length }} 个代码块</ATag>
        </template>

        <div v-for="block in luaBlocks" :key="block.filePath" class="mb-3">
          <div class="flex-between mb-1">
            <div>
              <span class="font-medium">{{ block.siteName }}</span>
              <span class="text-gray-400 text-xs ml-2">{{ block.filePath }}</span>
            </div>
          </div>

          <ACard
            v-for="(dir, di) in block.directives"
            :key="di"
            size="small"
            class="mb-1 bg-gray-50 dark:bg-gray-700"
          >
            <div class="flex-between mb-1">
              <ATag :color="dir.type.includes('_block') ? 'blue' : 'green'" class="text-xs">
                {{ dir.type }}
              </ATag>
              <AButton size="small" type="link" @click="selectedBlock = block">展开</AButton>
            </div>
            <pre class="text-xs font-mono text-gray-600 dark:text-gray-300 overflow-auto max-h-32 m-0">{{ dir.value.slice(0, 300) }}{{ dir.value.length > 300 ? '...' : '' }}</pre>
          </ACard>
        </div>

        <AEmpty v-if="luaBlocks.length === 0" description="暂无 Lua 代码块" />
      </ACard>

      <!-- Selected block detail -->
      <ACard v-if="selectedBlock" title="代码详情">
        <template #extra>
          <AButton size="small" @click="selectedBlock = null">关闭</AButton>
        </template>
        <div v-for="(dir, di) in selectedBlock!.directives" :key="di" class="mb-3">
          <div class="font-semibold mb-1">{{ dir.type }}</div>
          <pre class="bg-gray-900 text-green-400 p-3 rounded text-xs overflow-auto max-h-80">{{ dir.value }}</pre>
        </div>
      </ACard>
    </div>

    <!-- Right: Libraries & Dicts -->
    <div>
      <!-- Libraries -->
      <ACard title="Lua 库" class="mb-4">
        <div v-for="lib in libraries" :key="lib.name"
          class="flex-between py-2 border-b border-gray-100 dark:border-gray-700 last:border-0"
        >
          <div>
            <div class="font-medium text-sm">{{ lib.name }}</div>
            <div class="text-xs text-gray-400">{{ lib.description }}</div>
          </div>
          <ATag :color="lib.installed ? 'green' : 'default'" size="small">
            {{ lib.installed ? '已安装' : '未安装' }}
          </ATag>
        </div>
      </ACard>

      <!-- Shared Dicts -->
      <ACard title="共享字典 (Shared Dict)">
        <div v-for="dict in sharedDicts" :key="dict.name"
          class="mb-3 last:mb-0"
        >
          <div class="flex-between text-sm">
            <span class="font-medium">{{ dict.name }}</span>
            <span class="text-gray-400">{{ dict.size }}</span>
          </div>
          <AProgress
            :percent="parseInt(dict.usage)"
            :stroke-color="parseInt(dict.usage) > 80 ? '#ff4d4f' : '#1890ff'"
            size="small"
          />
        </div>
        <AEmpty v-if="sharedDicts.length === 0" description="无共享字典" />
      </ACard>
    </div>
  </div>
</template>
