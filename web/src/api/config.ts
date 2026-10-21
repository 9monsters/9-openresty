import api from './client'

export interface ConfigFile {
  name: string
  path: string
  type: 'file' | 'directory'
  size?: number
  mtime?: string
}

export function listConfig(path?: string) {
  return api.get('/config', { params: { path } }) as Promise<{ files: ConfigFile[] }>
}

export function readConfig(path: string) {
  return api.get('/config/read', { params: { path } }) as Promise<{ path: string; content: string }>
}

export function createConfig(data: { path: string; content: string; overwrite?: boolean }) {
  return api.post('/config', data)
}

export function updateConfig(path: string, content: string) {
  return api.put('/config', { path, content })
}

export function deleteConfig(path: string) {
  return api.delete('/config', { params: { path } })
}

export function getConfigHistory(path: string) {
  return api.get('/config/history', { params: { path } })
}
