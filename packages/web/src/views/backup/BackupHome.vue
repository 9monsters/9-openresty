<script setup lang="ts">
import { listBackups, getAutoBackupSettings } from '@/api/backup'
import dayjs from 'dayjs'
const backups = ref<any[]>([])
const autoSettings = ref<any>({})
onMounted(async () => {
  try {
    const [b, a] = await Promise.all([listBackups(), getAutoBackupSettings()])
    backups.value = (b as any).backups
    autoSettings.value = (a as any).settings
  } catch {}
})
</script>
<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">备份管理</h3>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <ACard title="备份列表" size="small">
        <ATable :data-source="backups" row-key="id" :pagination="false" size="small">
          <ATableColumn title="名称" data-index="name" />
          <ATableColumn title="类型" data-index="type" width="80"><template #default="{ record }"><ATag :color="record.type === 'auto' ? 'blue' : 'purple'">{{ record.type }}</ATag></template></ATableColumn>
          <ATableColumn title="大小" width="80"><template #default="{ record }">{{ (record.size / 1024).toFixed(1) }}K</template></ATableColumn>
          <ATableColumn title="时间" width="180"><template #default="{ record }">{{ dayjs(record.timestamp).format('YYYY-MM-DD HH:mm') }}</template></ATableColumn>
          <ATableColumn title="操作" width="120"><template #default><AButton size="small" type="link">恢复</AButton></template></ATableColumn>
        </ATable>
      </ACard>
      <ACard title="自动备份设置" size="small">
        <div class="grid grid-cols-2 gap-2 text-sm">
          <div class="text-gray-500">状态</div><div><ATag :color="autoSettings.enabled ? 'green' : 'default'">{{ autoSettings.enabled ? '已启用' : '已禁用' }}</ATag></div>
          <div class="text-gray-500">频率</div><div class="font-medium">{{ autoSettings.interval }}</div>
          <div class="text-gray-500">保留天数</div><div class="font-medium">{{ autoSettings.retention }}</div>
          <div class="text-gray-500">时间</div><div class="font-medium">{{ autoSettings.time }}</div>
          <div class="text-gray-500">路径</div><div class="font-medium text-xs">{{ autoSettings.path }}</div>
        </div>
      </ACard>
    </div>
  </div>
</template>
