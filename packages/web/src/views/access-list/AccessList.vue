<script setup lang="ts">
import { listAccessLists, createAccessList, deleteAccessList, addRule, deleteRule } from '@/api/access-lists'
import { t } from '@/i18n'

const loading = ref(true)
const lists = ref<any[]>([])
const modalVisible = ref(false)
const ruleModalVisible = ref(false)
const currentListId = ref<number | null>(null)
const listName = ref('')
const ruleForm = reactive({ directive: 'allow', address: '' })

async function fetchData() {
  loading.value = true
  try {
    const res = await listAccessLists()
    lists.value = res.lists
  } catch { /**/ }
  loading.value = false
}

async function handleCreateList() {
  if (!listName.value) return
  try {
    await createAccessList({ name: listName.value })
    listName.value = ''
    modalVisible.value = false
    fetchData()
  } catch { /**/ }
}

async function handleDeleteList(list: any) {
  try {
    await deleteAccessList(list.id)
    fetchData()
  } catch { /**/ }
}

function handleAddRule(listId: number) {
  currentListId.value = listId
  ruleForm.directive = 'allow'
  ruleForm.address = ''
  ruleModalVisible.value = true
}

async function handleSubmitRule() {
  if (!currentListId.value || !ruleForm.address) return
  try {
    await addRule(currentListId.value, { directive: ruleForm.directive, address: ruleForm.address })
    ruleModalVisible.value = false
    fetchData()
  } catch { /**/ }
}

async function handleDeleteRule(listId: number, ruleId: number) {
  try {
    await deleteRule(listId, ruleId)
    fetchData()
  } catch { /**/ }
}

const ruleColumns = [
  { title: 'Directive', dataIndex: 'directive', key: 'directive', width: 120 },
  { title: 'Address', dataIndex: 'address', key: 'address' },
  { title: '', key: 'actions', width: 80 },
]

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('Access Lists') }}</h3>
      <AButton type="primary" @click="modalVisible = true">
        <span class="i-tabler:plus mr-1" /> {{ t('Add') }}
      </AButton>
    </div>

    <ASpin v-if="loading" class="flex-center w-full py-10" />

    <div v-else class="space-y-4">
      <ACard v-for="list in lists" :key="list.id" :title="list.name" size="small">
        <template #extra>
          <AButton type="link" size="small" @click="handleAddRule(list.id)">
            + {{ t('Add Rule') || 'Add Rule' }}
          </AButton>
          <APopconfirm :title="t('Are you sure to delete') + '?'" @confirm="handleDeleteList(list)">
            <AButton type="link" size="small" danger>{{ t('Delete') }}</AButton>
          </APopconfirm>
        </template>

        <ATable :data-source="list.rules" :pagination="false" size="small" :columns="ruleColumns" row-key="id">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'directive'">
              <ATag :color="record.directive === 'allow' ? 'green' : 'red'">{{ record.directive }}</ATag>
            </template>
            <template v-if="column.key === 'actions'">
              <AButton type="link" size="small" danger @click="handleDeleteRule(list.id, record.id)">{{ t('Delete') }}</AButton>
            </template>
          </template>
        </ATable>
      </ACard>

      <AEmpty v-if="lists.length === 0" :description="t('No Data')" />
    </div>

    <!-- Create list modal -->
    <AModal v-model:open="modalVisible" :title="t('Add') + ' Access List'" @ok="handleCreateList">
      <AForm layout="vertical">
        <AFormItem :label="t('Name') || 'Name'" required>
          <AInput v-model:value="listName" placeholder="allow_internal" />
        </AFormItem>
      </AForm>
    </AModal>

    <!-- Add rule modal -->
    <AModal v-model:open="ruleModalVisible" :title="t('Add Rule') || 'Add Rule'" @ok="handleSubmitRule">
      <AForm layout="vertical">
        <AFormItem :label="'Directive'" required>
          <ARadioGroup v-model:value="ruleForm.directive">
            <ARadioButton value="allow">allow</ARadioButton>
            <ARadioButton value="deny">deny</ARadioButton>
          </ARadioGroup>
        </AFormItem>
        <AFormItem label="Address" required>
          <AInput v-model:value="ruleForm.address" placeholder="10.0.0.0/8 or all" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>