import api from './client'

export function getSettings() {
  return api.get('/settings') as Promise<Record<string, any>>
}

export function updateSettings(data: Record<string, any>) {
  return api.put('/settings', data)
}

export function getServerName() {
  return api.get('/settings/server-name') as Promise<{ name: string }>
}
