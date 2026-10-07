import api from './client'

export function getAccessLog(lines?: number) {
  return api.get('/logs/access', { params: { lines } }) as Promise<{ content: string; lines: number }>
}

export function getErrorLog(lines?: number) {
  return api.get('/logs/error', { params: { lines } }) as Promise<{ content: string; lines: number }>
}
