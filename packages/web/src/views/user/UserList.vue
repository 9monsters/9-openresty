<script setup lang="ts">
import { listUsers } from '@/api/users'
import dayjs from 'dayjs'
const users = ref<any[]>([])
onMounted(async () => { try { users.value = ((await listUsers()) as any).users } catch {} })
</script>
<template>
  <div>
    <div class="flex-between mb-4"><h3 class="text-xl font-semibold m-0">用户管理</h3><AButton type="primary" size="small"><span class="i-tabler:plus" /> 添加用户</AButton></div>
    <ATable :data-source="users" row-key="id">
      <ATableColumn title="用户名" data-index="name" />
      <ATableColumn title="邮箱" data-index="email" />
      <ATableColumn title="角色" data-index="role" />
      <ATableColumn title="状态" data-index="status" width="80"><template #default="{ record }"><ATag :color="record.status === 'active' ? 'green' : 'default'">{{ record.status }}</ATag></template></ATableColumn>
      <ATableColumn title="最后登录" width="180"><template #default="{ record }">{{ dayjs(record.lastLogin).format('YYYY-MM-DD HH:mm') }}</template></ATableColumn>
    </ATable>
  </div>
</template>
