import { FastifyPluginAsync } from 'fastify'
import { readLogFile } from '../services/nginx.js'

export const logRoutes: FastifyPluginAsync = async (app) => {
  // Get access log
  app.get('/logs/access', async (request) => {
    const query = request.query as { lines?: string }
    const lines = parseInt(query.lines || '200')
    const content = readLogFile('access', lines)
    return {
      type: 'access',
      lines: content.split('\n').filter(Boolean).length,
      content,
    }
  })

  // Get error log
  app.get('/logs/error', async (request) => {
    const query = request.query as { lines?: string }
    const lines = parseInt(query.lines || '200')
    const content = readLogFile('error', lines)
    return {
      type: 'error',
      lines: content.split('\n').filter(Boolean).length,
      content,
    }
  })
}
