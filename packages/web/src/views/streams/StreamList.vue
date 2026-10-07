<script setup lang="ts">
import { listStreams, createStream } from '@/api/streams'

const streams = ref<any[]>([])
const loading = ref(true)
const showAdd = ref(false)

const newStream = reactive({ name: '', listen: '', proxyPass: '', protocol: 'tcp' })

const columns = [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: '监听', dataIndex: 'listen', key: 'listen' },
  { title: '代理目标', dataIndex: 'proxyPass', key: 'proxyPass' },
  { title: '协议', key: 'protocol', width: 80 },
  { title: '状态', key: 'status', width: 100 },
]

async function fetch() {
  loading.value = true
  try {
    const res = await listStreams()
    streams.value = res.streams
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  try {
    await createStream({ ...newStream })
    showAdd.value = false
    Object.assign(newStream, { name: '', listen: '', proxyPass: '', protocol: 'tcp' })
    AMessage.success('Stream 已添加')
    await fetch()
  } catch (e: any) {
    AMessage.error(e.message)
  }
}

onMounted(fetch)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">Stream 管理</h3>
      <AButton type="primary" @click="showAdd = !showAdd">
        <span class="i-tabler:plus mr-1" /> 添加 Stream
      </AButton>
    </div>

    <!-- Add form -->
    <ACard v-if="showAdd" title="添加 Stream" class="mb-4">
      <AForm layout="inline">
        <AFormItem label="名称"><AInput v-model:value="newStream.name" /></AFormItem>
        <AFormItem label="监听端口"><AInput v-model:value="newStream.listen" placeholder="3306" /></AFormItem>
        <AFormItem label="代理目标"><AInput v-model:value="newStream.proxyPass" placeholder="192.168.1.100:3306" /></AFormItem>
        <AFormItem>
          <ASelect v-model:value="newStream.protocol" style="width:80px">
            <ASelectOption value="tcp">TCP</ASelectOption>
            <ASelectOption value="udp">UDP</ASelectOption>
          </ASelect>
        </AFormItem>
        <AFormItem>
          <ASpace><AButton type="primary" @click="handleCreate">添加</AButton><AButton @click="showAdd = false">取消</AButton></ASpace>
        </AFormItem>
      </AForm>
    </ACard>

    <ATable :columns="columns" :data-source="streams" :loading="loading" row-key="name">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'protocol'">
          <ATag :color="record.protocol === 'tcp' ? 'blue' : 'purple'">{{ record.protocol.toUpperCase() }}</ATag>
        </template>
        <template v-else-if="column.key === 'status'">
          <ATag :color="record.status === 'active' ? 'green' : 'default'">{{ record.status }}</ATag>
        </template>
      </template>
    </ATable>
  </div>
</template>
