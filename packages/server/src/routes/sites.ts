import { FastifyPluginAsync } from 'fastify'
import { listConfigFiles, readConfigFile, writeConfigFile, deleteConfigFile } from '../services/nginx.js'
import { parseSites, type ParsedSite, type LuaDirective } from '../services/site-parser.js'
import { backupConfig } from '../services/backup.js'

export const sitesRoutes: FastifyPluginAsync = async (app) => {
  // List all sites from config files
  app.get('/sites', async () => {
    const files = listConfigFiles()
    const allSites: ParsedSite[] = []

    for (const file of files) {
      if (file.type !== 'file') continue
      const content = readConfigFile(file.path)
      if (!content) continue
      const sites = parseSites(content, file.path)
      allSites.push(...sites)
    }

    return { sites: allSites }
  })

  // Get a single site by server name
  app.get('/sites/:name', async (request, reply) => {
    const { name } = request.params as { name: string }
    const files = listConfigFiles()

    for (const file of files) {
      if (file.type !== 'file') continue
      const content = readConfigFile(file.path)
      if (!content) continue
      const sites = parseSites(content, file.path)
      const site = sites.find(s => s.serverName.includes(name) || s.name === name)
      if (site) return site
    }

    return reply.status(404).send({ error: 'Site not found' })
  })

  // Create a new site
  app.post('/sites', async (request, reply) => {
    const { serverName, listen, root, ssl, luaDirectives, locations, accessLog, errorLog } = request.body as {
      serverName: string[]
      listen: string[]
      root?: string
      ssl?: boolean
      luaDirectives?: LuaDirective[]
      locations?: { path: string; proxyPass?: string; root?: string; luaDirectives?: LuaDirective[] }[]
      accessLog?: string
      errorLog?: string
    }

    if (!serverName?.length) {
      return reply.status(400).send({ error: 'Server name is required' })
    }

    const listenStrs = listen?.length ? listen : ['80']
    const luaBlock = luaDirectives?.map(d => buildLuaBlock(d)).join('\n') || ''
    const locBlocks = locations?.map(loc => {
      const locLua = loc.luaDirectives?.map(d => buildLuaBlock(d)).join('\n') || ''
      return `  location ${loc.path} {
    ${loc.proxyPass ? 'proxy_pass ' + loc.proxyPass + ';' : ''}
    ${loc.root ? 'root ' + loc.root + ';' : ''}
    ${locLua}
  }`
    }).join('\n\n') || ''

    const config = `server {
  listen ${listenStrs.join('\n  listen ')};
  server_name ${serverName.join(' ')};
  ${ssl ? 'ssl on;\n  ssl_certificate certs/' + serverName[0] + '.pem;\n  ssl_certificate_key certs/' + serverName[0] + '.key;' : ''}
  ${root ? 'root ' + root + ';' : ''}
  ${accessLog ? 'access_log ' + accessLog + ';' : ''}
  ${errorLog ? 'error_log ' + errorLog + ';' : ''}
  ${luaBlock}
${locBlocks}
}`

    const safeName = serverName[0].replace(/[*.]/g, '_')
    const filePath = `sites-enabled/${safeName}.conf`

    const existing = readConfigFile(filePath)
    if (existing !== null) {
      backupConfig(filePath, existing, 'Auto backup before creating new site')
    }

    const result = writeConfigFile(filePath, config)
    if (!result.ok) {
      return reply.status(500).send({ error: result.error })
    }

    const parsed = parseSites(config, filePath)
    return { ok: true, site: parsed[0], filePath }
  })

  // Update an existing site
  app.put('/sites/:name', async (request, reply) => {
    const { name } = request.params as { name: string }
    const { content } = request.body as { content: string }

    // Find the file containing this site
    const files = listConfigFiles()
    let foundPath = ''

    for (const file of files) {
      if (file.type !== 'file') continue
      const fileContent = readConfigFile(file.path)
      if (!fileContent) continue
      const sites = parseSites(fileContent, file.path)
      if (sites.some(s => s.serverName.includes(name) || s.name === name)) {
        foundPath = file.path
        break
      }
    }

    if (!foundPath) {
      return reply.status(404).send({ error: 'Site not found' })
    }

    const existing = readConfigFile(foundPath)
    if (existing) {
      backupConfig(foundPath, existing, 'Auto backup before update')
    }

    const result = writeConfigFile(foundPath, content)
    if (!result.ok) {
      return reply.status(500).send({ error: result.error })
    }

    return { ok: true, path: foundPath }
  })

  // Delete a site
  app.delete('/sites/:name', async (request, reply) => {
    const { name } = request.params as { name: string }
    const files = listConfigFiles()
    let foundPath = ''

    for (const file of files) {
      if (file.type !== 'file') continue
      const content = readConfigFile(file.path)
      if (!content) continue
      const sites = parseSites(content, file.path)
      if (sites.some(s => s.serverName.includes(name) || s.name === name)) {
        foundPath = file.path
        break
      }
    }

    if (!foundPath) {
      return reply.status(404).send({ error: 'Site not found' })
    }

    const existing = readConfigFile(foundPath)
    if (existing) {
      backupConfig(foundPath, existing, 'Backup before delete')
    }

    const result = deleteConfigFile(foundPath)
    if (!result.ok) {
      return reply.status(500).send({ error: result.error })
    }

    return { ok: true, name }
  })
}

function buildLuaBlock(directive: LuaDirective): string {
  if (directive.type.endsWith('_block')) {
    return `  ${directive.type.replace('_block', '')}_block {
    ${directive.value}
  }`
  }
  return `  ${directive.type} ${directive.value};`
}