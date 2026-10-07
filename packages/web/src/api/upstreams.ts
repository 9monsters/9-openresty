import api from './client'

export function listUpstreams() {
  return api.get('/upstreams')
}
