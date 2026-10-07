import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    children: [
      // --- Dashboard ---
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard/Dashboard.vue'), meta: { title: '仪表盘', icon: 'i-tabler:dashboard' } },
      { path: 'dashboard/server', name: 'ServerStatus', component: () => import('@/views/dashboard/ServerDash.vue'), meta: { title: '服务器', hidden: true } },
      { path: 'dashboard/nginx', name: 'NginxPerf', component: () => import('@/views/dashboard/NginxDash.vue'), meta: { title: 'Nginx 性能', hidden: true } },
      { path: 'dashboard/sites', name: 'SiteNav', component: () => import('@/views/dashboard/SiteDash.vue'), meta: { title: '站点导航', hidden: true } },

      // --- Sites ---
      { path: 'sites', name: 'Sites', component: () => import('@/views/sites/SiteList.vue'), meta: { title: '站点管理', icon: 'i-tabler:world-www' } },
      { path: 'sites/add', name: 'SiteAdd', component: () => import('@/views/sites/SiteEdit.vue'), meta: { title: '添加站点', hidden: true } },
      { path: 'sites/:name/edit', name: 'SiteEdit', component: () => import('@/views/sites/SiteEdit.vue'), meta: { title: '编辑站点', hidden: true } },
      { path: 'snippets', name: 'Snippets', component: () => import('@/views/snippet/SnippetList.vue'), meta: { title: '代码片段', icon: 'i-tabler:code-dots' } },

      // --- Streams ---
      { path: 'streams', name: 'Streams', component: () => import('@/views/streams/StreamList.vue'), meta: { title: 'Stream 管理', icon: 'i-tabler:arrows-exchange' } },

      // --- Upstream ---
      { path: 'upstream', name: 'Upstreams', component: () => import('@/views/upstreams/UpstreamList.vue'), meta: { title: 'Upstream', icon: 'i-tabler:server' } },

      // --- Config ---
      { path: 'config', name: 'Config', component: () => import('@/views/config/ConfigEditor.vue'), meta: { title: '配置文件', icon: 'i-tabler:file-code' } },

      // --- Certificates ---
      { path: 'certificates', name: 'Certificates', component: () => import('@/views/certificates/CertificateList.vue'), meta: { title: '证书管理', icon: 'i-tabler:certificate' } },
      { path: 'certificates/acme', name: 'ACMEUsers', component: () => import('@/views/certificate/AcmeUser.vue'), meta: { title: 'ACME 用户', icon: 'i-tabler:users' } },

      // --- Nginx Log ---
      { path: 'logs', name: 'Logs', component: () => import('@/views/logs/LogViewer.vue'), meta: { title: 'Nginx 日志', icon: 'i-tabler:file-text' } },

      // --- Terminal ---
      { path: 'terminal', name: 'Terminal', component: () => import('@/views/terminal/Terminal.vue'), meta: { title: '终端', icon: 'i-tabler:terminal-2' } },

      // --- DNS ---
      { path: 'dns', name: 'DNS', component: () => import('@/views/dns/DNSHome.vue'), meta: { title: 'DNS 管理', icon: 'i-tabler:cloud-network' } },

      // --- Access Lists ---
      { path: 'access-lists', name: 'AccessLists', component: () => import('@/views/access-list/AccessList.vue'), meta: { title: '访问控制', icon: 'i-tabler:shield-lock' } },

      // --- Namespaces ---
      { path: 'namespaces', name: 'Namespaces', component: () => import('@/views/namespace/Namespace.vue'), meta: { title: '命名空间', icon: 'i-tabler:stack-2' } },

      // --- Nodes ---
      { path: 'nodes', name: 'Nodes', component: () => import('@/views/node/NodeList.vue'), meta: { title: '节点管理', icon: 'i-tabler:database' } },

      // --- Notifications ---
      { path: 'notifications', name: 'Notifications', component: () => import('@/views/notification/NotificationList.vue'), meta: { title: '通知', icon: 'i-tabler:bell' } },

      // --- Lua Manager (OpenResty specific) ---
      { path: 'lua', name: 'Lua', component: () => import('@/views/lua/LuaManager.vue'), meta: { title: 'Lua 代码', icon: 'i-tabler:code' } },

      // --- Backup ---
      { path: 'backup', name: 'Backup', component: () => import('@/views/backup/BackupHome.vue'), meta: { title: '备份管理', icon: 'i-tabler:clock' } },

      // --- Preference (Settings) ---
      { path: 'preference', name: 'Preference', component: () => import('@/views/settings/SettingsPage.vue'), meta: { title: '系统设置', icon: 'i-tabler:settings' } },

      // --- System ---
      { path: 'system', name: 'System', component: () => import('@/views/system/SystemHome.vue'), meta: { title: '系统', icon: 'i-tabler:info-circle' } },

      // --- Users ---
      { path: 'users', name: 'Users', component: () => import('@/views/user/UserList.vue'), meta: { title: '用户管理', icon: 'i-tabler:user-group' } },
      { path: 'profile', name: 'Profile', component: () => import('@/views/user/Profile.vue'), meta: { title: '个人资料', icon: 'i-tabler:user' } },
    ],
  },
]

const router = createRouter({ history: createWebHashHistory(), routes })

router.beforeEach((to) => {
  document.title = (to.meta.title as string || 'OpenResty Admin') + ' | OpenResty Admin'
})

export default router