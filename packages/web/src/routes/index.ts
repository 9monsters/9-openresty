import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import {
  dashboardRoutes, sitesRoutes, streamsRoutes, upstreamRoutes,
  configRoutes, certificatesRoutes, logRoutes,
  luaRoutes, dnsRoutes, accessListRoutes, namespaceRoutes,
  nodeRoutes, notificationRoutes, backupRoutes, systemRoutes,
  userRoutes, settingsRoutes,
} from './modules/all'

export const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginPage.vue'),
    meta: { hiddenInSidebar: true, public: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: '',
        name: 'Dashboard Home',
        component: () => import('@/views/dashboard/HomePage.vue'),
        meta: { name: 'Dashboard', hiddenInSidebar: true },
      },
      ...dashboardRoutes,
      ...sitesRoutes,
      ...streamsRoutes,
      ...upstreamRoutes,
      ...accessListRoutes,
      ...configRoutes,
      ...certificatesRoutes,
      ...logRoutes,
      ...luaRoutes,
      ...dnsRoutes,
      ...namespaceRoutes,
      ...nodeRoutes,
      ...notificationRoutes,
      ...backupRoutes,
      ...systemRoutes,
      ...userRoutes,
      ...settingsRoutes,
    ],
  },
]

const router = createRouter({ history: createWebHashHistory(), routes })

router.beforeEach((to) => {
  const token = localStorage.getItem('openresty-token')
  const isPublic = to.meta?.public === true

  if (!token && !isPublic) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (token && to.path === '/login') {
    return { path: '/' }
  }

  const name = typeof to.meta?.name === 'function' ? '' : (to.meta?.name as string || '')
  document.title = name ? `${name} | OpenResty Admin` : 'OpenResty Admin'
})

export default router