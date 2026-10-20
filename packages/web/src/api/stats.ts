import api from './client'

export interface SystemHistoryPoint {
  time: string
  cpuCount: number
  memoryUsed: number
  memoryTotal: number
  loadAvg1: number
  loadAvg5: number
  loadAvg15: number
  uptimeSeconds: number
}

export interface NginxHistoryPoint {
  time: string
  activeConnections: number
  accepts: number
  handled: number
  requests: number
  reading: number
  writing: number
  waiting: number
}

export function getSystemHistory() {
  return api.get('/stats/system') as Promise<{ history: SystemHistoryPoint[] }>
}

export function getNginxHistory() {
  return api.get('/stats/nginx') as Promise<{ history: NginxHistoryPoint[] }>
}