import api from './client'
import type { ParsedSite } from '@/types'

export function listSites() {
  return api.get('/sites') as Promise<{ sites: ParsedSite[] }>
}

export function getSite(name: string) {
  return api.get(`/sites/${encodeURIComponent(name)}`) as Promise<ParsedSite>
}

export function createSite(data: {
  serverName: string[]
  listen: string[]
  root?: string
  ssl?: boolean
  luaDirectives?: any[]
  locations?: any[]
}) {
  return api.post('/sites', data)
}

export function updateSite(name: string, content: string) {
  return api.put(`/sites/${encodeURIComponent(name)}`, { content })
}

export function deleteSite(name: string) {
  return api.delete(`/sites/${encodeURIComponent(name)}`)
}
