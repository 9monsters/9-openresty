import api from './client'

export function listUpstreams(): Promise<{ upstreams: any[] }> {
  return api.get('/upstreams')
}

export function getUpstream(id: number): Promise<{ upstream: any }> {
  return api.get(`/upstreams/${id}`)
}

export function createUpstream(data: { name: string; strategy?: string; servers?: any[] }): Promise<{ ok: boolean; upstream: any }> {
  return api.post('/upstreams', data)
}

export function updateUpstream(id: number, data: any): Promise<{ ok: boolean; upstream: any }> {
  return api.put(`/upstreams/${id}`, data)
}

export function deleteUpstream(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/upstreams/${id}`)
}