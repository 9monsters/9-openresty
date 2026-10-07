import { FastifyPluginAsync } from 'fastify'
import { testConfig, reloadNginx, getNginxVersion, getNginxStatus } from '../services/nginx.js'

export const nginxRoutes: FastifyPluginAsync = async (app) => {
  // Get nginx status
  app.get('/nginx/status', async () => {
    const version = await getNginxVersion().catch(() => ({ version: 'Unknown', openresty: false }))
    return {
      status: getNginxStatus(),
      version: version.version,
      isOpenResty: version.openresty,
    }
  })

  // Test configuration
  app.post('/nginx/test', async () => {
    const result = await testConfig()
    return result
  })

  // Reload nginx
  app.post('/nginx/reload', async () => {
    const result = await reloadNginx()
    return result
  })

  // Restart nginx (mock implementation)
  app.post('/nginx/restart', async () => {
    return { ok: true, message: 'Restart signal sent' }
  })

  // Stop nginx (mock)
  app.post('/nginx/stop', async () => {
    return { ok: true, message: 'Stop signal sent' }
  })
}
