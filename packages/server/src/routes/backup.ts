import { FastifyPluginAsync } from 'fastify'

export const backupRoutes: FastifyPluginAsync = async (app) => {
  app.get('/backup', async () => ({
    backups: [
      { id: 1, name: 'daily-2026-10-07', type: 'auto', size: 245760, timestamp: new Date(Date.now() - 3600000).toISOString(), files: 12 },
      { id: 2, name: 'manual-2026-10-06', type: 'manual', size: 238912, timestamp: new Date(Date.now() - 86400000).toISOString(), files: 12 },
    ]
  }))

  app.get('/backup/auto-backup', async () => ({
    settings: { enabled: true, interval: 'daily', retention: 7, time: '03:00', path: '/var/backups/openresty' },
    schedule: { nextRun: '2026-10-08T03:00:00Z', lastRun: '2026-10-07T03:00:00Z', lastStatus: 'success' }
  }))
}
