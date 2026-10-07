import api from './client'

export interface DashboardData {
  status: string
  version: string
  isOpenResty: boolean
  stats: {
    activeConnections: number
    accepts: number
    handled: number
    requests: number
    reading: number
    writing: number
    waiting: number
  } | null
  serverTime: string
}

export function getDashboard() {
  return api.get('/dashboard') as Promise<DashboardData>
}
