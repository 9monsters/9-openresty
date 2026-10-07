<script setup lang="ts">
import { listAcmeUsers } from '@/api/acme'
const users = ref<any[]>([])
onMounted(async () => { try { users.value = ((await listAcmeUsers()) as any).users } catch {} })
</script>
<template>
  <div>
    <div class="flex-between mb-4"><h3 class="text-xl font-semibold m-0">ACME 用户</h3><AButton type="primary" size="small"><span class="i-tabler:plus" /> 添加</AButton></div>
    <ATable :data-source="users" row-key="id" :pagination="false">
      <ATableColumn title="邮箱" data-index="name" />
      <ATableColumn title="提供商" data-index="provider" />
      <ATableColumn title="目录服务器" data-index="server" ellipsis />
      <ATableColumn title="证书数" data-index="certCount" width="80" />
      <ATableColumn title="状态" data-index="status" width="80"><template #default="{ record }"><ATag :color="record.status === 'active' ? 'green' : 'default'">{{ record.status }}</ATag></template></ATableColumn>
    </ATable>
  </div>
</template>
