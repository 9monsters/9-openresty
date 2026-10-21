import type { RouteRecordRaw } from 'vue-router'
import { h } from 'vue'

const HomeOutlined = 'HomeOutlined'
const CloudOutlined = 'CloudOutlined'
const ShareAltOutlined = 'ShareAltOutlined'
const ClusterOutlined = 'ClusterOutlined'
const SafetyOutlined = 'SafetyOutlined'
const FileOutlined = 'FileOutlined'
const SafetyCertificateOutlined = 'SafetyCertificateOutlined'
const FileTextOutlined = 'FileTextOutlined'
const CodeOutlined = 'CodeOutlined'
const StarOutlined = 'StarOutlined'
const CloudServerOutlined = 'CloudServerOutlined'
const DatabaseOutlined = 'DatabaseOutlined'
const BellOutlined = 'BellOutlined'
const SettingOutlined = 'SettingOutlined'
const ClockCircleOutlined = 'ClockCircleOutlined'
const InfoCircleOutlined = 'InfoCircleOutlined'
const UsergroupAddOutlined = 'UsergroupAddOutlined'
const UserOutlined = 'UserOutlined'

export const dashboardRoutes: RouteRecordRaw[] = [
  {
    path: 'dashboard',
    component: () => import('@/layouts/BaseRouterView.vue'),
    name: 'Dashboard',
    meta: { name: 'Dashboard', icon: HomeOutlined },
    children: [
      { path: '', name: 'Dashboard Home', component: () => import('@/views/dashboard/HomePage.vue'), meta: { name: 'Dashboard', hiddenInSidebar: true } },
      { path: 'server', name: 'Server', component: () => import('@/views/dashboard/ServerDash.vue'), meta: { name: 'Server' } },
      { path: 'nginx', name: 'NginxPerformance', component: () => import('@/views/dashboard/NginxDash.vue'), meta: { name: 'Nginx' } },
      { path: 'sites', name: 'SiteNavigation', component: () => import('@/views/dashboard/SiteDash.vue'), meta: { name: 'Sites' } },
    ],
  },
]

export const sitesRoutes: RouteRecordRaw[] = [
  {
    path: 'sites',
    name: 'Manage Sites',
    component: () => import('@/layouts/BaseRouterView.vue'),
    meta: { name: 'Manage Sites', icon: CloudOutlined },
    children: [
      { path: '', name: 'Sites Home', component: () => import('@/views/sites/HomePage.vue'), meta: { name: 'Manage Sites', hiddenInSidebar: true } },
      { path: 'list', name: 'Sites List', component: () => import('@/views/sites/SiteList.vue'), meta: { name: 'Sites List' } },
      { path: 'add', name: 'Add Site', component: () => import('@/views/sites/SiteEdit.vue'), meta: { name: 'Add Site', hiddenInSidebar: true, lastRouteName: 'Sites List' } },
      { path: 'snippets', name: 'Snippets', component: () => import('@/views/snippet/SnippetList.vue'), meta: { name: 'Snippets' } },
      { path: ':name', name: 'Edit Site', component: () => import('@/views/sites/SiteEdit.vue'), meta: { name: 'Edit Site', hiddenInSidebar: true, lastRouteName: 'Sites List' } },
    ],
  },
]

export const streamsRoutes: RouteRecordRaw[] = [
  { path: 'streams', name: 'Manage Streams', component: () => import('@/views/streams/StreamList.vue'), meta: { name: 'Manage Streams', icon: ShareAltOutlined } },
  { path: 'streams/:name', name: 'Edit Stream', component: () => import('@/views/streams/StreamEdit.vue'), meta: { name: 'Edit Stream', hiddenInSidebar: true, lastRouteName: 'Manage Streams' } },
]

export const upstreamRoutes: RouteRecordRaw[] = [
  {
    path: 'upstream',
    name: 'Upstream Management',
    component: () => import('@/layouts/BaseRouterView.vue'),
    meta: { name: 'Upstream', icon: ClusterOutlined },
    children: [
      { path: '', name: 'Upstream Home', component: () => import('@/views/upstreams/HomePage.vue'), meta: { name: 'Upstream', hiddenInSidebar: true } },
      { path: 'list', name: 'Upstream List', component: () => import('@/views/upstreams/UpstreamList.vue'), meta: { name: 'Upstream Groups' } },
    ],
  },
]

export const accessListRoutes: RouteRecordRaw[] = [
  { path: 'access-lists', name: 'Access Lists', component: () => import('@/views/access-list/AccessList.vue'), meta: { name: 'Access Lists', icon: SafetyOutlined } },
]

export const configRoutes: RouteRecordRaw[] = [
  { path: 'config', name: 'Manage Configs', component: () => import('@/views/config/ConfigEditor.vue'), meta: { name: 'Manage Configs', icon: FileOutlined, hideChildren: true } },
]

export const certificatesRoutes: RouteRecordRaw[] = [
  {
    path: 'certificates',
    name: 'Certificates',
    component: () => import('@/layouts/BaseRouterView.vue'),
    meta: { name: 'Certificates', icon: SafetyCertificateOutlined },
    children: [
      { path: '', name: 'Certificates Home', component: () => import('@/views/certificates/HomePage.vue'), meta: { name: 'Certificates', hiddenInSidebar: true } },
      { path: 'acme_users', name: 'ACME User', component: () => import('@/views/certificate/AcmeUser.vue'), meta: { name: 'ACME User' } },
      { path: 'list', name: 'Certificates List', component: () => import('@/views/certificates/CertificateList.vue'), meta: { name: 'Certificates List' } },
    ],
  },
]

export const logRoutes: RouteRecordRaw[] = [
  {
    path: 'nginx_log',
    name: 'Nginx Log',
    component: () => import('@/layouts/BaseRouterView.vue'),
    meta: { name: 'Nginx Log', icon: FileTextOutlined },
    children: [
      { path: '', name: 'Nginx Log Home', component: () => import('@/views/logs/HomePage.vue'), meta: { name: 'Nginx Log', hiddenInSidebar: true } },
      { path: 'access', name: 'Access Logs', component: () => import('@/views/logs/LogViewer.vue'), meta: { name: 'Access Logs' } },
      { path: 'error', name: 'Error Logs', component: () => import('@/views/logs/LogViewer.vue'), meta: { name: 'Error Logs' } },
      { path: 'list', name: 'Log List', component: () => import('@/views/logs/LogList.vue'), meta: { name: 'Log List' } },
    ],
  },
]

// OpenResty-specific
export const luaRoutes: RouteRecordRaw[] = [
  { path: 'lua', name: 'Lua Manager', component: () => import('@/views/lua/LuaManager.vue'), meta: { name: 'Lua Manager', icon: StarOutlined } },
]

export const dnsRoutes: RouteRecordRaw[] = [
  {
    path: 'dns',
    name: 'DNS',
    component: () => import('@/layouts/BaseRouterView.vue'),
    meta: { name: 'DNS', icon: CloudServerOutlined },
    children: [
      { path: '', name: 'DNS Home', component: () => import('@/views/dns/HomePage.vue'), meta: { name: 'DNS', hiddenInSidebar: true } },
      { path: 'credentials', name: 'DNS Credentials', component: () => import('@/views/dns/DNSCredential.vue'), meta: { name: 'Credentials' } },
      { path: 'domains', name: 'DNS Domains', component: () => import('@/views/dns/DNSDomainList.vue'), meta: { name: 'DNS Domains' } },
      { path: 'groups', name: 'DNS Groups', component: () => import('@/views/dns/DNSGroupList.vue'), meta: { name: 'Groups' } },
      { path: 'ddns', name: 'DNS DDNS', component: () => import('@/views/dns/DDNSManager.vue'), meta: { name: 'DDNS' } },
    ],
  },
]

export const namespaceRoutes: RouteRecordRaw[] = [
  { path: 'namespaces', name: 'Namespaces', component: () => import('@/views/namespace/Namespace.vue'), meta: { name: 'Namespaces', icon: StarOutlined } },
]

export const nodeRoutes: RouteRecordRaw[] = [
  { path: 'nodes', name: 'Nodes', component: () => import('@/views/node/NodeList.vue'), meta: { name: 'Nodes', icon: DatabaseOutlined } },
]

export const notificationRoutes: RouteRecordRaw[] = [
  { path: 'notifications', name: 'Notifications', component: () => import('@/views/notification/NotificationList.vue'), meta: { name: 'Notifications', icon: BellOutlined } },
]

export const backupRoutes: RouteRecordRaw[] = [
  {
    path: 'backup',
    name: 'Backup',
    component: () => import('@/layouts/BaseRouterView.vue'),
    meta: { name: 'Backup Manage', icon: ClockCircleOutlined },
    children: [
      { path: '', name: 'Backup Home', component: () => import('@/views/backup/HomePage.vue'), meta: { name: 'Backup Manage', hiddenInSidebar: true } },
      { path: 'backup-and-restore', name: 'BackupAndRestore', component: () => import('@/views/backup/BackupHome.vue'), meta: { name: 'Backup' } },
      { path: 'auto-backup', name: 'AutoBackup', component: () => import('@/views/backup/AutoBackup.vue'), meta: { name: 'Auto Backup' } },
    ],
  },
]

export const systemRoutes: RouteRecordRaw[] = [
  {
    path: 'system',
    component: () => import('@/layouts/BaseRouterView.vue'),
    name: 'System',
    meta: { name: 'System', icon: InfoCircleOutlined },
    children: [
      { path: '', name: 'System Home', component: () => import('@/views/system/HomePage.vue'), meta: { name: 'System', hiddenInSidebar: true } },
      { path: 'self_check', name: 'Self Check', component: () => import('@/views/system/SelfCheck.vue'), meta: { name: 'Self Check' } },
      { path: 'upgrade', name: 'Upgrade', component: () => import('@/views/system/Upgrade.vue'), meta: { name: 'Upgrade' } },
      { path: 'about', name: 'About', component: () => import('@/views/system/About.vue'), meta: { name: 'About' } },
    ],
  },
]

export const userRoutes: RouteRecordRaw[] = [
  { path: 'users', name: 'Manage Users', component: () => import('@/views/user/UserList.vue'), meta: { name: 'Manage Users', icon: UsergroupAddOutlined } },
  { path: 'profile', name: 'User Profile', component: () => import('@/views/user/Profile.vue'), meta: { name: 'User Profile', icon: UserOutlined } },
]

export const settingsRoutes: RouteRecordRaw[] = [
  { path: 'preference', name: 'Preference', component: () => import('@/views/settings/SettingsPage.vue'), meta: { name: 'Preference', icon: SettingOutlined } },
  { path: 'preference/nginx-host-setup', name: 'Nginx Host SSH Setup', component: () => import('@/views/settings/NginxHostSetup.vue'), meta: { name: 'Host SSH Setup', hiddenInSidebar: true, lastRouteName: 'Preference' } },
]