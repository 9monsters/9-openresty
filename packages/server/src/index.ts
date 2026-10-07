import Fastify from 'fastify'
import cors from '@fastify/cors'
import fastifyStatic from '@fastify/static'
import { join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'
import { dashboardRoutes } from './routes/dashboard.js'
import { configRoutes } from './routes/config.js'
import { sitesRoutes } from './routes/sites.js'
import { upstreamRoutes } from './routes/upstreams.js'
import { streamRoutes } from './routes/streams.js'
import { certificateRoutes } from './routes/certificates.js'
import { logRoutes } from './routes/logs.js'
import { nginxRoutes } from './routes/nginx.js'
import { luaRoutes } from './routes/lua.js'
import { settingsRoutes } from './routes/settings.js'
import { terminalRoutes } from './routes/terminal.js'
import { dnsRoutes } from './routes/dns.js'
import { accessListRoutes } from './routes/access-lists.js'
import { namespaceRoutes } from './routes/namespaces.js'
import { nodeRoutes } from './routes/nodes.js'
import { notificationRoutes } from './routes/notifications.js'
import { backupRoutes as apiBackupRoutes } from './routes/backup.js'
import { systemRoutes } from './routes/system.js'
import { userRoutes } from './routes/users.js'
import { snippetRoutes } from './routes/snippets.js'
import { acmeRoutes } from './routes/acme.js'

const __dirname = dirname(fileURLToPath(import.meta.url))

export async function createServer() {
  const app = Fastify({
    logger: {
      transport: {
        target: 'pino-pretty',
        options: { colorize: true }
      }
    }
  })

  await app.register(cors, { origin: true })

  // API routes
  await app.register(dashboardRoutes, { prefix: '/api' })
  await app.register(configRoutes, { prefix: '/api' })
  await app.register(sitesRoutes, { prefix: '/api' })
  await app.register(upstreamRoutes, { prefix: '/api' })
  await app.register(streamRoutes, { prefix: '/api' })
  await app.register(certificateRoutes, { prefix: '/api' })
  await app.register(logRoutes, { prefix: '/api' })
  await app.register(nginxRoutes, { prefix: '/api' })
  await app.register(luaRoutes, { prefix: '/api' })
    await app.register(terminalRoutes, { prefix: '/api' })
  await app.register(dnsRoutes, { prefix: '/api' })
  await app.register(accessListRoutes, { prefix: '/api' })
  await app.register(namespaceRoutes, { prefix: '/api' })
  await app.register(nodeRoutes, { prefix: '/api' })
  await app.register(notificationRoutes, { prefix: '/api' })
  await app.register(apiBackupRoutes, { prefix: '/api' })
  await app.register(systemRoutes, { prefix: '/api' })
  await app.register(userRoutes, { prefix: '/api' })
  await app.register(snippetRoutes, { prefix: '/api' })
  await app.register(acmeRoutes, { prefix: '/api' })
  await app.register(settingsRoutes, { prefix: '/api' })

  // Serve built frontend in production
  const webDist = join(__dirname, '../../web/dist')
  await app.register(fastifyStatic, {
    root: webDist,
    prefix: '/',
    wildcard: false
  })

  // SPA fallback
  app.setNotFoundHandler(async (request, reply) => {
    if (request.url.startsWith('/api/')) {
      return reply.status(404).send({ error: 'Not found' })
    }
    return reply.sendFile('index.html')
  })

  return app
}

// Start directly when run as main module
const port = parseInt(process.env.PORT || '3000')
const host = process.env.HOST || '0.0.0.0'

const app = await createServer()
await app.listen({ port, host })
console.log(`OpenResty Admin server running at http://${host}:${port}`)