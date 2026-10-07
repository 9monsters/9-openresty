import { FastifyPluginAsync } from 'fastify'
import { listConfigFiles, readConfigFile, writeConfigFile } from '../services/nginx.js'
import { parseSites } from '../services/site-parser.js'

export const luaRoutes: FastifyPluginAsync = async (app) => {
  // List all Lua code blocks across all configs
  app.get('/lua', async () => {
    const files = listConfigFiles()
    const luaBlocks: { filePath: string; siteName: string; directives: any[] }[] = []

    for (const file of files) {
      if (file.type !== 'file') continue
      const content = readConfigFile(file.path)
      if (!content) continue
      const sites = parseSites(content, file.path)
      for (const site of sites) {
        if (site.luaDirectives.length > 0) {
          luaBlocks.push({
            filePath: site.filePath,
            siteName: site.serverName.join(' ') || 'unnamed',
            directives: site.luaDirectives,
          })
        }
      }
    }

    return { luaBlocks }
  })

  // List available Lua libraries
  app.get('/lua/libraries', async () => {
    return {
      libraries: [
        { name: 'resty.http', description: 'HTTP client for OpenResty', installed: true },
        { name: 'resty.redis', description: 'Redis client for OpenResty', installed: true },
        { name: 'resty.mysql', description: 'MySQL client for OpenResty', installed: true },
        { name: 'resty.websocket', description: 'WebSocket support for OpenResty', installed: true },
        { name: 'resty.limit.traffic', description: 'Traffic limiting module', installed: true },
        { name: 'resty.core', description: 'Core OpenResty API', installed: true },
        { name: 'cjson', description: 'Fast JSON encoding/decoding', installed: true },
      ]
    }
  })

  // List lua shared dicts
  app.get('/lua/shared-dicts', async () => {
    return {
      dicts: [
        { name: 'cache', size: '10m', usage: '45%' },
        { name: 'sessions', size: '20m', usage: '12%' },
        { name: 'rate_limit', size: '5m', usage: '78%' },
      ]
    }
  })
}
