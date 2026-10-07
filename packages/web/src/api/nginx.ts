import api from './client'

export function getNginxStatus() {
  return api.get('/nginx/status') as Promise<{ status: string; version: string; isOpenResty: boolean }>
}

export function testConfig() {
  return api.post('/nginx/test') as Promise<{ ok: boolean; output: string }>
}

export function reloadNginx() {
  return api.post('/nginx/reload') as Promise<{ ok: boolean; output: string }>
}
