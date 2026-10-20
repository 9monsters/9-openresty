import api from './client'

export function listAcmeUsers(): Promise<{ users: any[] }> {
  return api.get('/acme-users')
}

export function createAcmeUser(data: { email: string; provider?: string; server?: string }): Promise<{ ok: boolean; user: any }> {
  return api.post('/acme-users', data)
}

export function updateAcmeUser(id: number, data: any): Promise<{ ok: boolean; user: any }> {
  return api.put(`/acme-users/${id}`, data)
}

export function deleteAcmeUser(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/acme-users/${id}`)
}