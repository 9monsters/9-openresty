import { FastifyPluginAsync } from 'fastify'

export const nodeRoutes: FastifyPluginAsync = async (app) => {
  app.get('/nodes', async () => ({
    nodes: [
      { id: 1, name: 'lb-01', host: '192.168.1.10', port: 9001, status: 'online', version: 'OpenResty 1.25.3.1', uptime: '15d 4h' },
      { id: 2, name: 'lb-02', host: '192.168.1.11', port: 9001, status: 'online', version: 'OpenResty 1.25.3.1', uptime: '15d 3h' },
      { id: 3, name: 'cache-01', host: '10.0.0.5', port: 9001, status: 'offline', version: 'N/A', uptime: 'N/A' },
    ]
  }))
}
