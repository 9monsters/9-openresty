import { FastifyPluginAsync } from 'fastify'

export const userRoutes: FastifyPluginAsync = async (app) => {
  app.get('/users', async () => ({
    users: [
      { id: 1, name: 'admin', email: 'admin@openresty.local', role: 'Administrator', status: 'active', lastLogin: new Date().toISOString() },
      { id: 2, name: 'operator', email: 'ops@openresty.local', role: 'Operator', status: 'active', lastLogin: new Date(Date.now() - 86400000).toISOString() },
    ]
  }))

  app.get('/users/profile', async () => ({
    id: 1, name: 'admin', email: 'admin@openresty.local', role: 'Administrator',
    twoFactorEnabled: false, theme: 'auto', language: 'zh_CN',
  }))
}
