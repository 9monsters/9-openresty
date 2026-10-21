import api from './client'

export function getAccessLog(lines?: number) {
  return api.get('/logs/access', { params: { lines } }) as Promise<{ content: string; lines: number }>
}

export function getErrorLog(lines?: number) {
  return api.get('/logs/error', { params: { lines } }) as Promise<{ content: string; lines: number }>
}

export function getLogFiles() {
  return api.get('/logs/files') as Promise<{ files: Array<{ name: string; path: string; size: number; mtime: string }> }>
}

export function readLogByPath(path: string, lines?: number) {
  return api.get('/logs/file', { params: { path, lines } }) as Promise<{ content: string }>
}
