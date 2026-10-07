import { FastifyPluginAsync } from 'fastify'
import os from 'node:os'

export const systemRoutes: FastifyPluginAsync = async (app) => {
  app.get('/system/self-check', async () => ({
    checks: [
      { name: 'Nginx Binary', status: 'ok', detail: '/usr/local/openresty/nginx/sbin/nginx' },
      { name: 'Config Syntax', status: 'ok', detail: 'nginx.conf syntax is ok' },
      { name: 'Write Permission', status: 'ok', detail: 'Config directory is writable' },
      { name: 'Port Availability', status: 'ok', detail: 'Port 80/443 available' },
      { name: 'SSL Certificates', status: 'warning', detail: '1 certificate expiring in 15 days' },
    ]
  }))

  app.get('/system/upgrade', async () => ({
    currentVersion: '2.5.0',
    latestVersion: '2.5.1',
    updateAvailable: true,
    changelog: 'Bug fixes and performance improvements',
  }))

  app.get('/system/about', async () => ({
    name: 'OpenResty Admin',
    version: '0.1.0',
    buildTime: '2026-10-07T09:00:00Z',
    node: process.version,
    platform: os.platform(),
    arch: os.arch(),
    openresty: '1.25.3.1',
    goVersion: null,
    dbType: 'SQLite',
  }))
}
