<script setup lang="ts">
const requestsPerSec = ref(1523)
const avgLatency = ref('24ms')
const errorRate = ref('0.3%')
const bandwidth = ref('8.2 Mbps')
</script>
<template>
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
    <ACard size="small" class="text-center"><div class="text-3xl font-bold text-blue-600">{{ requestsPerSec }}</div><div class="text-gray-500 text-sm mt-1">请求/秒</div></ACard>
    <ACard size="small" class="text-center"><div class="text-3xl font-bold text-green-600">{{ avgLatency }}</div><div class="text-gray-500 text-sm mt-1">平均延迟</div></ACard>
    <ACard size="small" class="text-center"><div class="text-3xl font-bold" :class="parseFloat(errorRate) > 1 ? 'text-red-600' : 'text-orange-600'">{{ errorRate }}</div><div class="text-gray-500 text-sm mt-1">错误率</div></ACard>
    <ACard size="small" class="text-center"><div class="text-3xl font-bold text-purple-600">{{ bandwidth }}</div><div class="text-gray-500 text-sm mt-1">带宽</div></ACard>
  </div>
</template>
