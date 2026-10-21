import api from './client'

export function listNotifications(): Promise<{ notifications: any[]; total: number }> {
  return api.get('/notifications')
}

export function createNotification(data: { type: string; title: string; message: string }): Promise<{ ok: boolean; notification: any }> {
  return api.post('/notifications', data)
}

export function markRead(id: number): Promise<{ ok: boolean }> {
  return api.put(`/notifications/${id}/read`)
}

export function markAllRead(): Promise<{ ok: boolean }> {
  return api.put('/notifications/read-all')
}

export function deleteNotification(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/notifications/${id}`)
}