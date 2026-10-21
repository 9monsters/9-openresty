<script setup lang="ts">
import { listSnippets, listTemplates, createSnippet, updateSnippet, deleteSnippet } from '@/api/snippets'
import { t } from '@/i18n'

const loading = ref(true)
const snippets = ref<any[]>([])
const templates = ref<any[]>([])
const modalVisible = ref(false)
const editing = ref<any>(null)
const formData = reactive({ name: '', description: '', content: '' })

const columns = [
  { title: t('Name') || 'Name', dataIndex: 'name', key: 'name' },
  { title: t('Description') || 'Description', dataIndex: 'description', key: 'description' },
  { title: t('Usage') || 'Usage', dataIndex: 'usage', key: 'usage' },
  { title: t('Actions'), key: 'actions', width: 180 },
]

async function fetchData() {
  loading.value = true
  try {
    const [snipRes, tmplRes] = await Promise.all([listSnippets(), listTemplates()])
    snippets.value = snipRes.snippets
    templates.value = tmplRes.templates
  } catch { /**/ }
  loading.value = false
}

function handleAdd() {
  editing.value = null
  formData.name = ''
  formData.description = ''
  formData.content = ''
  modalVisible.value = true
}

function handleEdit(record: any) {
  editing.value = record
  formData.name = record.name
  formData.description = record.description || ''
  formData.content = record.content || ''
  modalVisible.value = true
}

async function handleDelete(record: any) {
  try {
    await deleteSnippet(record.id)
    fetchData()
  } catch { /**/ }
}

async function handleSubmit() {
  try {
    if (editing.value) {
      await updateSnippet(editing.value.id, { ...formData })
    } else {
      await createSnippet({ ...formData })
    }
    modalVisible.value = false
    fetchData()
  } catch { /**/ }
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Snippets') }}</h3>
      <AButton type="primary" @click="handleAdd">
        <span class="i-tabler:plus mr-1" /> {{ t('Create') }} {{ t('Snippet') || 'Snippet' }}
      </AButton>
    </div>

    <ASpin v-if="loading" class="flex-center w-full py-10" />

    <div v-else class="grid grid-cols-2 gap-6">
      <!-- User snippets -->
      <ACard :title="t('Snippets')">
        <ATable :dataSource="snippets" :columns="columns" :loading="loading" rowKey="id" :pagination="false" size="small">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'actions'">
              <AButton type="link" size="small" @click="handleEdit(record)">{{ t('Edit') }}</AButton>
              <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDelete(record)">
                <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
              </APopconfirm>
            </template>
          </template>
        </ATable>
        <div v-if="snippets.length === 0" class="text-center text-gray-400 py-8">
          <AEmpty :description="t('No Data')" />
        </div>
      </ACard>

      <!-- Built-in templates -->
      <ACard :title="t('Built-in Templates') || 'Built-in Templates'">
        <div class="space-y-3">
          <div v-for="tmpl in templates" :key="tmpl.name" class="flex-between p-3 bg-gray-50 dark:bg-gray-800 rounded">
            <div>
              <div class="font-semibold">{{ tmpl.name }}</div>
              <div class="text-sm text-gray-500">{{ tmpl.description }}</div>
            </div>
            <ATag>{{ tmpl.category }}</ATag>
          </div>
        </div>
        <div v-if="templates.length === 0" class="text-center text-gray-400 py-8">
          <AEmpty :description="t('No Data')" />
        </div>
      </ACard>
    </div>

    <AModal v-model:open="modalVisible" :title="(editing ? t('Edit') : t('Create')) + ' Snippet'" @ok="handleSubmit" width="600px">
      <AForm layout="vertical">
        <AFormItem :label="t('Name') || 'Name'" required>
          <AInput v-model:value="formData.name" placeholder="my-config.conf" />
        </AFormItem>
        <AFormItem :label="t('Description') || 'Description'">
          <AInput v-model:value="formData.description" />
        </AFormItem>
        <AFormItem :label="t('Content') || 'Content'">
          <ATextarea v-model:value="formData.content" :rows="10" class="font-mono text-sm" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>