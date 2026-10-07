import { FastifyPluginAsync } from 'fastify'

export const upstreamRoutes: FastifyPluginAsync = async (app) => {
  // List upstreams with mock data
  app.get('/upstreams', async () => {
    return {
      upstreams: [
        { name: 'backend_api', strategy: 'round_robin', servers: [
          { address: '192.168.1.10:8080', weight: 1, maxFails: 3, failTimeout: '30s', status: 'up' },
          { address: '192.168.1.11:8080', weight: 2, maxFails: 3, failTimeout: '30s', status: 'up' },
          { address: '192.168.1.12:8080', weight: 1, maxFails: 3, failTimeout: '30s', status: 'down' },
        ]},
        { name: 'cache_servers', strategy: 'hash', hashKey: '$request_uri', servers: [
          { address: '10.0.0.1:6379', weight: 1, status: 'up' },
          { address: '10.0.0.2:6379', weight: 2, status: 'up' },
        ]},
        { name: 'websocket_backend', strategy: 'ip_hash', servers: [
          { address: '192.168.2.10:9090', weight: 1, status: 'up' },
          { address: '192.168.2.11:9090', weight: 1, status: 'up' },
        ]},
      ]
    }
  })

  // Add a new upstream
  app.post('/upstreams', async (request) => {
    const { name, strategy, servers } = request.body as {
      name: string
      strategy?: string
      servers?: { address: string; weight?: number; maxFails?: number; failTimeout?: string }[]
    }
    return { ok: true, name, strategy: strategy || 'round_robin', servers }
  })
}
