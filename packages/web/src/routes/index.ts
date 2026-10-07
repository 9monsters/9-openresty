import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        meta: { title: '仪表盘', icon: 'i-tabler:dashboard' },
        component: () => import('@/views/dashboard/Dashboard.vue'),
      },
      {
        path: 'sites',
        name: 'Sites',
        meta: { title: '站点管理', icon: 'i-tabler:world-www' },
        component: () => import('@/views/sites/SiteList.vue'),
      },
      {
        path: 'sites/add',
        name: 'SiteAdd',
        meta: { title: '添加站点', hidden: true },
        component: () => import('@/views/sites/SiteEdit.vue'),
      },
      {
        path: 'sites/:name/edit',
        name: 'SiteEdit',
        meta: { title: '编辑站点', hidden: true },
        component: () => import('@/views/sites/SiteEdit.vue'),
      },
      {
        path: 'streams',
        name: 'Streams',
        meta: { title: 'Stream 管理', icon: 'i-tabler:arrows-exchange' },
        component: () => import('@/views/streams/StreamList.vue'),
      },
      {
        path: 'upstreams',
        name: 'Upstreams',
        meta: { title: 'Upstream 管理', icon: 'i-tabler:server' },
        component: () => import('@/views/upstreams/UpstreamList.vue'),
      },
      {
        path: 'config',
        name: 'Config',
        meta: { title: '配置文件', icon: 'i-tabler:file-code' },
        component: () => import('@/views/config/ConfigEditor.vue'),
      },
      {
        path: 'certificates',
        name: 'Certificates',
        meta: { title: '证书管理', icon: 'i-tabler:certificate' },
        component: () => import('@/views/certificates/CertificateList.vue'),
      },
      {
        path: 'lua',
        name: 'Lua',
        meta: { title: 'Lua 代码', icon: 'i-tabler:code' },
        component: () => import('@/views/lua/LuaManager.vue'),
      },
      {
        path: 'logs',
        name: 'Logs',
        meta: { title: '日志查看', icon: 'i-tabler:file-text' },
        component: () => import('@/views/logs/LogViewer.vue'),
      },
      {
        path: 'settings',
        name: 'Settings',
        meta: { title: '系统设置', icon: 'i-tabler:settings' },
        component: () => import('@/views/settings/SettingsPage.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.beforeEach((to) => {
  document.title = `${to.meta.title || 'OpenResty Admin'} | OpenResty Admin`
})

export default router
