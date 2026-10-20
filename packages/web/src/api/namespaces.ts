import api from './client'

export function listNamespaces(): Promise<{ namespaces: any[] }> {
  return api.get('/namespaces')
}

export function createNamespace(data: { name: string; description?: string }): Promise<{ ok: boolean; namespace: any }> {
  return api.post('/namespaces', data)
}

export function updateNamespace(id: number, data: any): Promise<{ ok: boolean; namespace: any }> {
  return api.put(`/namespaces/${id}`, data)
}

export function deleteNamespace(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/namespaces/${id}`)
}