<script setup lang="ts">
import { listConfig, readConfig, updateConfig, deleteConfig, createConfig } from '@/api/config'
import { testConfig } from '@/api/nginx'
import { t } from '@/i18n'
import type { ConfigFile } from '@/types'

const files = ref<ConfigFile[]>([])
const selectedFile = ref('')
const content = ref('')
const originalContent = ref('')
const loading = ref(true)
const saving = ref(false)
const showNewFile = ref(false)
const newFileName = ref('')

// File tree
const fileTree = computed(() => {
  const tree: any[] = []
  const dirMap = new Map<string, any>()

  for (const f of files.value) {
    const parts = f.path.split('/')
    if (f.type === 'directory') {
      dirMap.set(f.path, { name: f.name, path: f.path, children: [] })
    }
  }

  for (const f of files.value) {
    if (f.type === 'file') {
      const parentPath = f.path.includes('/') ? f.path.split('/').slice(0, -1).join('/') : ''
      const node = { name: f.name, path: f.path, isFile: true, size: f.size, mtime: f.mtime }
      if (parentPath && dirMap.has(parentPath)) {
        dirMap.get(parentPath)!.children.push(node)
      } else {
        tree.push(node)
      }
    }
  }

  for (const [, dir] of dirMap) {
    if (!tree.find(n => n.name === dir.name)) tree.push(dir)
  }

  return tree
})

async function fetchFiles() {
  loading.value = true
  try {
    const res = await listConfig()
    files.value = res.files
  } finally {
    loading.value = false
  }
}

async function selectFile(f: { path: string }) {
  try {
    const res = await readConfig(f.path)
    selectedFile.value = res.path
    content.value = res.content
    originalContent.value = res.content
  } catch (e) {
    message.error('读取文件失败')
  }
}

async function handleSave() {
  if (!selectedFile.value) return
  saving.value = true
  try {
    await updateConfig(selectedFile.value, content.value)
    originalContent.value = content.value
    message.success(t('Published and reloaded'))
  } catch (e: any) {
    message.error(e.message)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (!selectedFile.value) return
  try {
    await deleteConfig(selectedFile.value)
    selectedFile.value = ''
    content.value = ''
    message.success(t('Published and reloaded'))
    await fetchFiles()
  } catch (e: any) {
    message.error(e.message)
  }
}

async function handleCreate() {
  if (!newFileName.value) return
  try {
    const path = newFileName.value.endsWith('.conf') ? newFileName.value : newFileName.value + '.conf'
    await createConfig({ path, content: '# New Config\n', overwrite: false })
    showNewFile.value = false
    newFileName.value = ''
    message.success(t('Published and reloaded'))
    await fetchFiles()
  } catch (e: any) {
    message.error(e.message)
  }
}

async function handleTest() {
  try {
    const result = await testConfig()
    if (result.ok) {
      message.success('配置测试通过！')
    } else {
      message.warning(result.output)
    }
  } catch (e: any) {
    message.error(e.message)
  }
}


const isModified = computed(() => content.value !== originalContent.value)

onMounted(fetchFiles)
</script>

<template>
  <div class="flex h-[calc(100vh-12rem)] gap-4">
    <!-- File Tree Panel -->
    <div class="w-64 flex-shrink-0 border-r border-gray-200 dark:border-gray-700 pr-4 overflow-auto">
      <div class="flex-between mb-3">
        <h4 class="font-semibold m-0">配置文件</h4>
        <AButton size="small" @click="showNewFile = true">
          <span class="i-tabler:plus" />
        </AButton>
      </div>

      <!-- New file dialog -->
      <div v-if="showNewFile" class="mb-3 p-2 bg-gray-50 dark:bg-gray-700 rounded">
        <AInput v-model:value="newFileName" placeholder="文件名.conf" size="small" class="mb-1" />
        <AButton size="small" type="primary" block @click="handleCreate">创建</AButton>
      </div>

      <!-- File tree -->
      <div v-for="node in fileTree" :key="node.path" class="mb-1">
        <div
          v-if="!node.isFile"
          class="px-2 py-1 text-xs text-gray-500 font-semibold uppercase"
        >
          📁 {{ node.name }}
        </div>
        <div
          v-else
          class="px-2 py-1 cursor-pointer rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-sm"
          :class="{ 'bg-blue-50 dark:bg-blue-900 text-blue-600': selectedFile === node.path }"
          @click="selectFile(node)"
        >
          📄 {{ node.name }}
          <span class="text-xs text-gray-400 ml-1">{{ node.size ? (node.size / 1024).toFixed(1) + 'K' : '' }}</span>
        </div>

        <!-- Children (files in directory) -->
        <template v-if="node.children">
          <div
            v-for="child in node.children"
            :key="child.path"
            class="px-4 py-1 cursor-pointer rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-sm ml-2"
            :class="{ 'bg-blue-50 dark:bg-blue-900 text-blue-600': selectedFile === child.path }"
            @click="selectFile(child)"
          >
            📄 {{ child.name }}
          </div>
        </template>
      </div>
    </div>

    <!-- Editor Panel -->
    <div class="flex-1 flex flex-col min-w-0">
      <div v-if="!selectedFile" class="flex-center flex-1 text-gray-400">
        选择左侧文件进行编辑
      </div>

      <template v-else>
        <!-- Toolbar -->
        <div class="flex-between mb-2">
          <div class="text-sm text-gray-500">{{ selectedFile }}</div>
          <ASpace>
            <AButton size="small" @click="handleTest">
              <span class="i-tabler:test-pipe mr-1" /> 测试配置
            </AButton>
            <AButton size="small" :danger="isModified" type="primary" @click="handleSave" :loading="saving">
              <span class="i-tabler:device-floppy mr-1" /> {{ t('Save and publish') }}
            </AButton>
            <AButton size="small" danger @click="handleDelete">
              <span class="i-tabler:trash" />
            </AButton>
          </ASpace>
        </div>

        <!-- Editor -->
        <ATextarea
          v-model:value="content"
          class="flex-1 font-mono text-sm"
          :auto-size="false"
          style="height: 100%"
        />
      </template>
    </div>
  </div>
</template>
