import api from './client'

export function listBackups(): Promise<{ backups: any[] }> {
  return api.get('/backup')
}

export function createBackup(name?: string): Promise<{ ok: boolean; backup: any }> {
  return api.post('/backup', { name })
}

export function restoreBackup(id: number): Promise<{ ok: boolean; message: string }> {
  return api.post(`/backup/${id}/restore`)
}

export function deleteBackup(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/backup/${id}`)
}

export function getAutoBackup(): Promise<{ settings: any; schedule: any }> {
  return api.get('/backup/auto-backup')
}

export function updateAutoBackup(data: any): Promise<{ ok: boolean; settings: any }> {
  return api.put('/backup/auto-backup', data)
}