import { FastifyPluginAsync } from 'fastify'

export const namespaceRoutes: FastifyPluginAsync = async (app) => {
  app.get('/namespaces', async () => ({
    namespaces: [
      { name: 'default', description: 'Default configuration namespace', siteCount: 5, enabled: true },
      { name: 'staging', description: 'Staging environment', siteCount: 2, enabled: true },
    ]
  }))
}
