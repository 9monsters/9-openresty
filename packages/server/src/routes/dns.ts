import { FastifyPluginAsync } from 'fastify'

export const dnsRoutes: FastifyPluginAsync = async (app) => {
  app.get('/dns/credentials', async () => ({
    credentials: [
      { id: 1, name: 'Cloudflare', provider: 'cloudflare', email: 'admin@example.com', status: 'active' },
      { id: 2, name: 'Aliyun DNS', provider: 'alidns', accessKeyId: 'ak-***', status: 'active' },
    ]
  }))

  app.get('/dns/domains', async () => ({
    domains: [
      { id: 1, name: 'example.com', provider: 'Cloudflare', recordCount: 12, status: 'managed' },
      { id: 2, name: 'api.example.com', provider: 'Cloudflare', recordCount: 3, status: 'managed' },
    ]
  }))

  app.get('/dns/groups', async () => ({
    groups: [
      { id: 1, name: 'Production', domains: ['example.com', 'app.example.com'], strategy: 'round_robin' },
      { id: 2, name: 'Staging', domains: ['staging.example.com'], strategy: 'weighted' },
    ]
  }))

  app.get('/dns/ddns', async () => ({
    ddns: [
      { id: 1, domain: 'home.example.com', credential: 'Cloudflare', ip: '203.0.113.1', lastUpdate: new Date().toISOString(), enabled: true },
    ]
  }))
}
