import { FastifyPluginAsync } from 'fastify'

export const notificationRoutes: FastifyPluginAsync = async (app) => {
  app.get('/notifications', async () => ({
    notifications: [
      { id: 1, type: 'success', title: 'Config Reloaded', message: 'Nginx configuration reloaded successfully', time: new Date(Date.now() - 3600000).toISOString() },
      { id: 2, type: 'warning', title: 'Certificate Expiring', message: 'SSL certificate for example.com expires in 15 days', time: new Date(Date.now() - 7200000).toISOString() },
      { id: 3, type: 'error', title: 'Upstream Down', message: 'Backend server 192.168.1.12:8080 is unreachable', time: new Date(Date.now() - 10800000).toISOString() },
      { id: 4, type: 'info', title: 'Backup Completed', message: 'Daily config backup completed successfully', time: new Date(Date.now() - 86400000).toISOString() },
    ]
  }))
}
