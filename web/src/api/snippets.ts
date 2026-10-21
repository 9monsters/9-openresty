import api from './client'

export function listSnippets(): Promise<{ snippets: any[] }> {
  return api.get('/snippets')
}

export function getSnippet(id: number): Promise<{ snippet: any }> {
  return api.get(`/snippets/${id}`)
}

export function createSnippet(data: { name: string; description?: string; content?: string }): Promise<{ ok: boolean; snippet: any }> {
  return api.post('/snippets', data)
}

export function updateSnippet(id: number, data: any): Promise<{ ok: boolean; snippet: any }> {
  return api.put(`/snippets/${id}`, data)
}

export function deleteSnippet(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/snippets/${id}`)
}

export function listTemplates(): Promise<{ templates: any[] }> {
  return api.get('/snippets/templates')
}