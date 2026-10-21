import api from './client'

export interface User {
  id: number
  username: string
  email: string
  role: string
  status: string
  twoFactorEnabled?: boolean
  lastLogin?: string
  language?: string
}

export function listUsers(): Promise<{ users: User[] }> {
  return api.get('/users')
}

export function createUser(data: Partial<User> & { password: string }): Promise<{ ok: boolean; user: User }> {
  return api.post('/users', data)
}

export function updateUser(id: number, data: Partial<User>): Promise<{ ok: boolean; user: User }> {
  return api.put(`/users/${id}`, data)
}

export function deleteUser(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/users/${id}`)
}

export function getProfile(): Promise<User> {
  return api.get('/users/profile')
}

export function updateProfile(data: Partial<User>): Promise<{ ok: boolean; user: User }> {
  return api.put('/users/profile', data)
}

export function changePassword(currentPassword: string, newPassword: string): Promise<{ ok: boolean }> {
  return api.put('/users/profile/password', { currentPassword, newPassword })
}