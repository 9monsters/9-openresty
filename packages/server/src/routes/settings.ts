import { FastifyPluginAsync } from 'fastify'
import { getNginxPaths, updateNginxPaths } from '../services/nginx.js'

let savedSettings: Record<string, any> = {
  nginxBin: '/usr/local/openresty/nginx/sbin/nginx',
  nginxConfDir: '/usr/local/openresty/nginx/conf',
  nginxConfPath: '/usr/local/openresty/nginx/conf/nginx.conf',
  nginxLogDir: '/usr/local/openresty/nginx/logs',
  nginxPIDPath: '/usr/local/openresty/nginx/logs/nginx.pid',
  serverName: 'OpenResty Admin',
  language: 'zh_CN',
  theme: 'auto',
}

export const settingsRoutes: FastifyPluginAsync = async (app) => {
  // Get current settings
  app.get('/settings', async () => {
    return {
      ...savedSettings,
      paths: getNginxPaths(),
    }
  })

  // Update settings
  app.put('/settings', async (request) => {
    const body = request.body as Record<string, any>
    Object.assign(savedSettings, body)
    
    if (body.nginxBin || body.nginxConfDir || body.nginxConfPath || body.nginxLogDir || body.nginxPIDPath) {
      updateNginxPaths({
        nginxBin: body.nginxBin,
        nginxConfDir: body.nginxConfDir,
        nginxConfPath: body.nginxConfPath,
        nginxLogDir: body.nginxLogDir,
        nginxPIDPath: body.nginxPIDPath,
      })
    }

    return { ok: true, settings: savedSettings }
  })

  // Get server name
  app.get('/settings/server-name', async () => {
    return { name: savedSettings.serverName }
  })
}
