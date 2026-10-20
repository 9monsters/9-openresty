import api from './client'

export function listAccessLists(): Promise<{ lists: any[] }> {
  return api.get('/access-lists')
}

export function createAccessList(data: { name: string; description?: string }): Promise<{ ok: boolean; list: any }> {
  return api.post('/access-lists', data)
}

export function updateAccessList(id: number, data: any): Promise<{ ok: boolean; list: any }> {
  return api.put(`/access-lists/${id}`, data)
}

export function deleteAccessList(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/access-lists/${id}`)
}

export function addRule(listId: number, rule: { directive: string; address: string }): Promise<{ ok: boolean; rules: any[] }> {
  return api.post(`/access-lists/${listId}/rules`, rule)
}

export function deleteRule(listId: number, ruleId: number): Promise<{ ok: boolean; rules: any[] }> {
  return api.delete(`/access-lists/${listId}/rules/${ruleId}`)
}