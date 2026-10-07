<script setup lang="ts">
import { getProfile } from '@/api/users'
const profile = ref<any>(null)
onMounted(async () => { try { profile.value = await getProfile() } catch {} })
</script>
<template>
  <div class="max-w-lg">
    <h3 class="text-xl font-semibold mb-4">个人资料</h3>
    <ACard v-if="profile" size="small">
      <AForm layout="vertical">
        <AFormItem label="用户名"><AInput :value="profile.name" disabled /></AFormItem>
        <AFormItem label="邮箱"><AInput :value="profile.email" /></AFormItem>
        <AFormItem label="角色"><AInput :value="profile.role" disabled /></AFormItem>
        <AFormItem label="双因素认证"><ASwitch :checked="profile.twoFactorEnabled" /></AFormItem>
        <AFormItem label="语言"><ASelect :value="profile.language"><ASelectOption value="zh_CN">简体中文</ASelectOption><ASelectOption value="en">English</ASelectOption></ASelect></AFormItem>
        <AButton type="primary">保存</AButton>
      </AForm>
    </ACard>
  </div>
</template>
