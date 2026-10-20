import api from './client'

export function listStreams(): Promise<{ streams: any[] }> {
  return api.get('/streams')
}

export function createStream(data: { name: string; listen: string; proxyPass: string; protocol?: string }): Promise<{ ok: boolean; stream: any }> {
  return api.post('/streams', data)
}

export function updateStream(id: number, data: any): Promise<{ ok: boolean; stream: any }> {
  return api.put(`/streams/${id}`, data)
}

export function deleteStream(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/streams/${id}`)
}