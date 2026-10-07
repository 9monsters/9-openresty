<script setup lang="ts">
import { listCredentials, listDomains, listGroups, listDDNS } from '@/api/dns'
const [credentials, domains, groups, ddns] = [ref<any[]>([]), ref<any[]>([]), ref<any[]>([]), ref<any[]>([])]
onMounted(async () => {
  try {
    const [c, d, g, dd] = await Promise.all([listCredentials(), listDomains(), listGroups(), listDDNS()])
    credentials.value = (c as any).credentials; domains.value = (d as any).domains
    groups.value = (g as any).groups; ddns.value = (dd as any).ddns
  } catch {}
})
</script>
<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">DNS 管理</h3>
    <div class="grid grid-cols-2 gap-4">
      <ACard title="凭据" size="small"><div v-for="c in credentials" :key="c.id" class="flex-between py-1 border-b last:border-0 text-sm"><span>{{ c.name }}</span><ATag color="blue">{{ c.provider }}</ATag></div></ACard>
      <ACard title="域名" size="small"><div v-for="d in domains" :key="d.id" class="flex-between py-1 border-b last:border-0 text-sm"><span>{{ d.name }}</span><span class="text-gray-400">{{ d.recordCount }} records</span></div></ACard>
      <ACard title="分组" size="small"><div v-for="g in groups" :key="g.id" class="flex-between py-1 border-b last:border-0 text-sm"><span>{{ g.name }}</span><ATag>{{ g.strategy }}</ATag></div></ACard>
      <ACard title="DDNS" size="small"><div v-for="d in ddns" :key="d.id" class="flex-between py-1 border-b last:border-0 text-sm"><span>{{ d.domain }}</span><ATag :color="d.enabled ? 'green' : 'default'">{{ d.ip }}</ATag></div></ACard>
    </div>
  </div>
</template>
