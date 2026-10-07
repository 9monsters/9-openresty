import { FastifyPluginAsync } from 'fastify'
import { listConfigFiles, readConfigFile, writeConfigFile, deleteConfigFile } from '../services/nginx.js'
import { backupConfig, listBackups } from '../services/backup.js'

export const configRoutes: FastifyPluginAsync = async (app) => {
  // List config files
  app.get('/config', async (request) => {
    const query = request.query as { path?: string }
    return { files: listConfigFiles(query.path || '') }
  })

  // Read a config file (path as query param)
  app.get('/config/read', async (request) => {
    const { path } = request.query as { path?: string }
    if (!path) return { error: 'path required' }
    const content = readConfigFile(path)
    if (content === null) return { error: 'File not found', path }
    return { path, content }
  })

  // Create a new config file
  app.post('/config', async (request, reply) => {
    const { path, content, overwrite } = request.body as { path: string; content: string; overwrite?: boolean }
    const existing = readConfigFile(path)
    if (existing !== null && !overwrite) {
      return reply.status(409).send({ error: 'File already exists' })
    }
    if (existing !== null) {
      backupConfig(path, existing, 'Auto backup before overwrite')
    }
    const result = writeConfigFile(path, content)
    if (!result.ok) return reply.status(500).send({ error: result.error })
    return { ok: true, path }
  })

  // Update a config file
  app.put('/config', async (request, reply) => {
    const { path, content } = request.body as { path: string; content: string }
    const existing = readConfigFile(path)
    if (existing !== null) {
      backupConfig(path, existing, 'Auto backup before edit')
    }
    const result = writeConfigFile(path, content)
    if (!result.ok) return reply.status(500).send({ error: result.error })
    return { ok: true, path }
  })

  // Delete a config file
  app.delete('/config', async (request, reply) => {
    const { path } = request.query as { path?: string }
    if (!path) return reply.status(400).send({ error: 'path required' })
    const existing = readConfigFile(path)
    if (existing !== null) {
      backupConfig(path, existing, 'Backup before delete')
    }
    const result = deleteConfigFile(path)
    if (!result.ok) return reply.status(500).send({ error: result.error })
    return { ok: true, path }
  })

  // List backup history
  app.get('/config/history', async (request) => {
    const { path } = request.query as { path?: string }
    if (!path) return { error: 'path required' }
    return {
      backups: listBackups(path).map(b => ({
        id: b.id, filepath: b.filepath, timestamp: b.timestamp, message: b.message,
      }))
    }
  })

  // Get a history entry
  app.get('/config/history/:id', async (request) => {
    const { id } = request.params as { id: string }
    const { path } = request.query as { path?: string }
    if (!path) return { error: 'path required' }
    const backups = listBackups(path)
    const backup = backups.find(b => b.id === id)
    if (!backup) return { error: 'Backup not found' }
    return { path, content: backup.content, timestamp: backup.timestamp, message: backup.message }
  })
}
