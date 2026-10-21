<script setup lang="ts">
import {
  listDnsCredentials, createDnsCredential, updateDnsCredential, deleteDnsCredential,
  listDnsDomains, createDnsDomain, updateDnsDomain, deleteDnsDomain,
  listDnsGroups, createDnsGroup, updateDnsGroup, deleteDnsGroup,
  listDdnsEntries, createDdnsEntry, updateDdnsEntry, deleteDdnsEntry,
} from '@/api/dns'
import { t } from '@/i18n'

const activeTab = ref('credentials')
const loading = ref(true)

// Credentials
const credentials = ref<any[]>([])
const credVisible = ref(false)
const credEditing = ref<any>(null)
const credForm = reactive({ name: '', provider: 'cloudflare', email: '', accessKeyId: '' })

// Domains
const domains = ref<any[]>([])
const domainVisible = ref(false)
const domainEditing = ref<any>(null)
const domainForm = reactive({ name: '', provider: '' })

// Groups
const groups = ref<any[]>([])
const groupVisible = ref(false)
const groupEditing = ref<any>(null)
const groupForm = reactive({ name: '', strategy: 'round_robin' })

// DDNS
const ddnsEntries = ref<any[]>([])
const ddnsVisible = ref(false)
const ddnsEditing = ref<any>(null)
const ddnsForm = reactive({ domain: '', credential: '' })

async function fetchAll() {
  loading.value = true
  try {
    const [cr, dm, gr, dd] = await Promise.all([
      listDnsCredentials(), listDnsDomains(), listDnsGroups(), listDdnsEntries(),
    ])
    credentials.value = cr.credentials
    domains.value = dm.domains
    groups.value = gr.groups
    ddnsEntries.value = dd.ddns
  } catch { /**/ }
  loading.value = false
}

// --- Credentials CRUD ---
function addCredential() {
  credEditing.value = null
  credForm.name = ''; credForm.provider = 'cloudflare'; credForm.email = ''; credForm.accessKeyId = ''
  credVisible.value = true
}
function editCredential(r: any) {
  credEditing.value = r
  credForm.name = r.name; credForm.provider = r.provider; credForm.email = r.email || ''; credForm.accessKeyId = r.accessKeyId || ''
  credVisible.value = true
}
async function saveCredential() {
  try {
    if (credEditing.value) await updateDnsCredential(credEditing.value.id, { ...credForm })
    else await createDnsCredential({ ...credForm })
    credVisible.value = false; fetchAll()
  } catch { /**/ }
}
async function delCredential(r: any) { try { await deleteDnsCredential(r.id); fetchAll() } catch { /**/ } }

// --- Domains CRUD ---
function addDomain() { domainEditing.value = null; domainForm.name = ''; domainForm.provider = ''; domainVisible.value = true }
function editDomain(r: any) { domainEditing.value = r; domainForm.name = r.name; domainForm.provider = r.provider; domainVisible.value = true }
async function saveDomain() {
  try {
    if (domainEditing.value) await updateDnsDomain(domainEditing.value.id, { ...domainForm })
    else await createDnsDomain({ ...domainForm })
    domainVisible.value = false; fetchAll()
  } catch { /**/ }
}
async function delDomain(r: any) { try { await deleteDnsDomain(r.id); fetchAll() } catch { /**/ } }

// --- Groups CRUD ---
function addGroup() { groupEditing.value = null; groupForm.name = ''; groupForm.strategy = 'round_robin'; groupVisible.value = true }
function editGroup(r: any) { groupEditing.value = r; groupForm.name = r.name; groupForm.strategy = r.strategy; groupVisible.value = true }
async function saveGroup() {
  try {
    if (groupEditing.value) await updateDnsGroup(groupEditing.value.id, { ...groupForm })
    else await createDnsGroup({ ...groupForm })
    groupVisible.value = false; fetchAll()
  } catch { /**/ }
}
async function delGroup(r: any) { try { await deleteDnsGroup(r.id); fetchAll() } catch { /**/ } }

// --- DDNS CRUD ---
function addDdns() { ddnsEditing.value = null; ddnsForm.domain = ''; ddnsForm.credential = ''; ddnsVisible.value = true }
function editDdns(r: any) { ddnsEditing.value = r; ddnsForm.domain = r.domain; ddnsForm.credential = r.credential; ddnsVisible.value = true }
async function saveDdns() {
  try {
    if (ddnsEditing.value) await updateDdnsEntry(ddnsEditing.value.id, { ...ddnsForm })
    else await createDdnsEntry({ ...ddnsForm })
    ddnsVisible.value = false; fetchAll()
  } catch { /**/ }
}
async function delDdns(r: any) { try { await deleteDdnsEntry(r.id); fetchAll() } catch { /**/ } }

const credColumns = [
  { title: 'Name', dataIndex: 'name', key: 'name' },
  { title: 'Provider', dataIndex: 'provider', key: 'provider' },
  { title: 'Status', dataIndex: 'status', key: 'status', width: 100 },
  { title: '', key: 'actions', width: 150 },
]
const domainColumns = [
  { title: 'Domain', dataIndex: 'name', key: 'name' },
  { title: 'Provider', dataIndex: 'provider', key: 'provider' },
  { title: 'Records', dataIndex: 'recordCount', key: 'recordCount', width: 80 },
  { title: 'Status', dataIndex: 'status', key: 'status', width: 100 },
  { title: '', key: 'actions', width: 150 },
]
const groupColumns = [
  { title: 'Name', dataIndex: 'name', key: 'name' },
  { title: 'Strategy', dataIndex: 'strategy', key: 'strategy' },
  { title: 'Domains', key: 'domains' },
  { title: '', key: 'actions', width: 150 },
]
const ddnsColumns = [
  { title: 'Domain', dataIndex: 'domain', key: 'domain' },
  { title: 'IP', dataIndex: 'ip', key: 'ip' },
  { title: 'Credential', dataIndex: 'credential', key: 'credential' },
  { title: 'Enabled', dataIndex: 'enabled', key: 'enabled', width: 80 },
  { title: '', key: 'actions', width: 150 },
]

onMounted(fetchAll)
</script>

<template>
  <div>
    <div class="flex-between mb-4">
      <h3 class="text-xl font-semibold m-0">{{ t('DNS') }}</h3>
    </div>

    <ATabs v-model:activeKey="activeTab">
      <ATabPane key="credentials" tab="DNS Credentials">
        <div class="flex justify-end mb-3">
          <AButton type="primary" size="small" @click="addCredential"><span class="i-tabler:plus mr-1" /> Add</AButton>
        </div>
        <ATable :data-source="credentials" :loading="loading" :columns="credColumns" row-key="id" :pagination="false" size="small">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'provider'">
              <ATag>{{ record.provider }}</ATag>
            </template>
            <template v-if="column.key === 'status'">
              <ATag :color="record.status === 'active' ? 'green' : 'red'">{{ record.status }}</ATag>
            </template>
            <template v-if="column.key === 'actions'">
              <AButton type="link" size="small" @click="editCredential(record)">{{ t('Edit') }}</AButton>
              <AButton type="link" size="small" danger @click="delCredential(record)">{{ t('Delete') }}</AButton>
            </template>
          </template>
        </ATable>
      </ATabPane>

      <ATabPane key="domains" tab="DNS Domains">
        <div class="flex justify-end mb-3">
          <AButton type="primary" size="small" @click="addDomain"><span class="i-tabler:plus mr-1" /> Add</AButton>
        </div>
        <ATable :data-source="domains" :loading="loading" :columns="domainColumns" row-key="id" :pagination="false" size="small">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <ATag color="blue">{{ record.status }}</ATag>
            </template>
            <template v-if="column.key === 'actions'">
              <AButton type="link" size="small" @click="editDomain(record)">{{ t('Edit') }}</AButton>
              <AButton type="link" size="small" danger @click="delDomain(record)">{{ t('Delete') }}</AButton>
            </template>
          </template>
        </ATable>
      </ATabPane>

      <ATabPane key="groups" tab="DNS Groups">
        <div class="flex justify-end mb-3">
          <AButton type="primary" size="small" @click="addGroup"><span class="i-tabler:plus mr-1" /> Add</AButton>
        </div>
        <ATable :data-source="groups" :loading="loading" :columns="groupColumns" row-key="id" :pagination="false" size="small">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'strategy'">
              <ATag>{{ record.strategy }}</ATag>
            </template>
            <template v-if="column.key === 'domains'">
              <ATag v-for="d in record.domains" :key="d" class="mr-1">{{ d }}</ATag>
            </template>
            <template v-if="column.key === 'actions'">
              <AButton type="link" size="small" @click="editGroup(record)">{{ t('Edit') }}</AButton>
              <AButton type="link" size="small" danger @click="delGroup(record)">{{ t('Delete') }}</AButton>
            </template>
          </template>
        </ATable>
      </ATabPane>

      <ATabPane key="ddns" tab="DDNS">
        <div class="flex justify-end mb-3">
          <AButton type="primary" size="small" @click="addDdns"><span class="i-tabler:plus mr-1" /> Add</AButton>
        </div>
        <ATable :data-source="ddnsEntries" :loading="loading" :columns="ddnsColumns" row-key="id" :pagination="false" size="small">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'enabled'">
              <ATag :color="record.enabled ? 'green' : 'red'">{{ record.enabled ? t('Yes') || 'Yes' : t('No') || 'No' }}</ATag>
            </template>
            <template v-if="column.key === 'actions'">
              <AButton type="link" size="small" @click="editDdns(record)">{{ t('Edit') }}</AButton>
              <AButton type="link" size="small" danger @click="delDdns(record)">{{ t('Delete') }}</AButton>
            </template>
          </template>
        </ATable>
      </ATabPane>
    </ATabs>

    <!-- Credential Modal -->
    <AModal v-model:open="credVisible" title="DNS Credential" @ok="saveCredential">
      <AForm layout="vertical">
        <AFormItem label="Name" required><AInput v-model:value="credForm.name" /></AFormItem>
        <AFormItem label="Provider"><AInput v-model:value="credForm.provider" /></AFormItem>
        <AFormItem label="Email"><AInput v-model:value="credForm.email" /></AFormItem>
        <AFormItem label="API Key"><AInput v-model:value="credForm.accessKeyId" /></AFormItem>
      </AForm>
    </AModal>

    <!-- Domain Modal -->
    <AModal v-model:open="domainVisible" title="DNS Domain" @ok="saveDomain">
      <AForm layout="vertical">
        <AFormItem label="Domain" required><AInput v-model:value="domainForm.name" placeholder="example.com" /></AFormItem>
        <AFormItem label="Provider"><AInput v-model:value="domainForm.provider" placeholder="cloudflare" /></AFormItem>
      </AForm>
    </AModal>

    <!-- Group Modal -->
    <AModal v-model:open="groupVisible" title="DNS Group" @ok="saveGroup">
      <AForm layout="vertical">
        <AFormItem label="Name" required><AInput v-model:value="groupForm.name" /></AFormItem>
        <AFormItem label="Strategy"><ASelect v-model:value="groupForm.strategy">
          <ASelectOption value="round_robin">Round Robin</ASelectOption>
          <ASelectOption value="weighted">Weighted</ASelectOption>
          <ASelectOption value="geo">Geographic</ASelectOption>
        </ASelect></AFormItem>
      </AForm>
    </AModal>

    <!-- DDNS Modal -->
    <AModal v-model:open="ddnsVisible" title="DDNS Entry" @ok="saveDdns">
      <AForm layout="vertical">
        <AFormItem label="Domain" required><AInput v-model:value="ddnsForm.domain" placeholder="home.example.com" /></AFormItem>
        <AFormItem label="Credential"><AInput v-model:value="ddnsForm.credential" placeholder="Cloudflare" /></AFormItem>
      </AForm>
    </AModal>
  </div>
</template>