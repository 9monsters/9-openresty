<script setup lang="ts">
const cpuUsage = ref(42)
const memUsage = ref(68)
const diskUsage = ref(35)
const uptime = ref('15d 4h 32m')
</script>
<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <ACard title="CPU 使用率" size="small"><AProgress :percent="cpuUsage" :stroke-color="{ '0%': '#108ee9', '100%': '#87d068' }" /></ACard>
    <ACard title="内存使用率" size="small"><AProgress :percent="memUsage" :stroke-color="{ '0%': '#ff4d4f', '100%': '#ffa940' }" /></ACard>
    <ACard title="磁盘使用率" size="small"><AProgress :percent="diskUsage" /></ACard>
    <ACard title="系统运行时间" size="small" class="flex-center"><span class="text-2xl font-bold text-blue-600">{{ uptime }}</span></ACard>
  </div>
</template>
