import { FastifyPluginAsync } from 'fastify'
import { getNginxVersion, getNginxStatus, getNginxStats } from '../services/nginx.js'

export const dashboardRoutes: FastifyPluginAsync = async (app) => {
  app.get('/dashboard', async () => {
    const [version, stats] = await Promise.all([
      getNginxVersion().catch(() => ({ version: 'Unknown', openresty: false })),
      getNginxStats().catch(() => null),
    ])
    const status = getNginxStatus()

    return {
      status,
      version: version.version,
      isOpenResty: version.openresty,
      stats,
      serverTime: new Date().toISOString(),
    }
  })

  // Node monitoring data
  app.get('/dashboard/nodes', async () => {
    return {
      nodes: [{
        id: 1,
        name: 'local',
        host: 'localhost',
        status: getNginxStatus(),
        version: (await getNginxVersion().catch(() => ({ version: 'Unknown' }))).version,
        uptime: process.uptime(),
        memory: process.memoryUsage(),
      }]
    }
  })
}
