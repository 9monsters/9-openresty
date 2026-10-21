import api from './client'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  expiresAt: string
  user: { id: number; username: string; role: string }
}

export async function login(params: LoginParams): Promise<LoginResult> {
  return api.post('/auth/login', params)
}

export async function refresh(): Promise<LoginResult> {
  return api.post('/auth/refresh')
}