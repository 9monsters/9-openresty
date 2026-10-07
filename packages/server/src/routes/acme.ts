import { FastifyPluginAsync } from 'fastify'

export const acmeRoutes: FastifyPluginAsync = async (app) => {
  app.get('/acme-users', async () => ({
    users: [
      { id: 1, name: 'admin@example.com', provider: "Let's Encrypt", server: 'https://acme-v02.api.letsencrypt.org/directory', status: 'active', certCount: 3 },
      { id: 2, name: 'dev@example.com', provider: 'ZeroSSL', server: 'https://acme.zerossl.com/v2/DV90', status: 'active', certCount: 1 },
    ]
  }))
}
