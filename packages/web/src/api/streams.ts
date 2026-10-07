import api from './client'

export function listStreams() {
  return api.get('/streams')
}

export function createStream(data: { name: string; listen: string; proxyPass: string; protocol?: string }) {
  return api.post('/streams', data)
}
