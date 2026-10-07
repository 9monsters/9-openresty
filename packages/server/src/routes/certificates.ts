import { FastifyPluginAsync } from 'fastify'

export const certificateRoutes: FastifyPluginAsync = async (app) => {
  app.get('/certificates', async () => {
    return {
      certificates: [
        { name: 'example.com', domains: ['example.com', '*.example.com'], issuer: "Let's Encrypt",
          notBefore: '2026-01-15T00:00:00Z', notAfter: '2026-04-15T00:00:00Z',
          autoRenew: true, status: 'valid' },
        { name: 'api.example.com', domains: ['api.example.com'], issuer: "Let's Encrypt",
          notBefore: '2026-02-01T00:00:00Z', notAfter: '2026-05-01T00:00:00Z',
          autoRenew: true, status: 'valid' },
        { name: 'internal.local', domains: ['*.internal.local'], issuer: 'Self-signed',
          notBefore: '2025-12-01T00:00:00Z', notAfter: '2027-12-01T00:00:00Z',
          autoRenew: false, status: 'expiring' },
      ]
    }
  })
}
