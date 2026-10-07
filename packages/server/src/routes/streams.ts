import { FastifyPluginAsync } from 'fastify'

export const streamRoutes: FastifyPluginAsync = async (app) => {
  app.get('/streams', async () => {
    return {
      streams: [
        { name: 'mysql_proxy', listen: '3306', proxyPass: '192.168.1.100:3306', protocol: 'tcp', status: 'active' },
        { name: 'redis_proxy', listen: '6379', proxyPass: '10.0.1.10:6379', protocol: 'tcp', status: 'active' },
        { name: 'dns_proxy', listen: '53', proxyPass: '8.8.8.8:53', protocol: 'udp', status: 'disabled' },
      ]
    }
  })

  app.post('/streams', async (request) => {
    const { name, listen, proxyPass, protocol } = request.body as {
      name: string; listen: string; proxyPass: string; protocol?: string
    }
    return { ok: true, name, listen, proxyPass, protocol: protocol || 'tcp' }
  })
}
