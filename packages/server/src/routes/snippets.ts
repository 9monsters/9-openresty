import { FastifyPluginAsync } from 'fastify'

export const snippetRoutes: FastifyPluginAsync = async (app) => {
  app.get('/snippets', async () => ({
    snippets: [
      { name: 'proxy.conf', description: 'Reverse proxy configuration snippet', usage: 'proxy settings' },
      { name: 'ssl.conf', description: 'SSL/TLS configuration snippet', usage: 'ssl settings' },
      { name: 'cache.conf', description: 'Proxy cache configuration', usage: 'cache settings' },
    ]
  }))

  app.get('/snippets/templates', async () => ({
    templates: [
      { name: 'wordpress', description: 'WordPress site template', category: 'CMS' },
      { name: 'proxy', description: 'Generic reverse proxy template', category: 'Proxy' },
      { name: 'static', description: 'Static file server template', category: 'Web' },
      { name: 'websocket', description: 'WebSocket proxy template', category: 'Proxy' },
    ]
  }))
}
