import api from './client'

export function getSystemInfo(): Promise<{
  cpu: { count: number; model: string }
  memory: { total: number; free: number; used: number; usagePercent: number }
  uptime: number
  loadAvg: number[]
  platform: string
  arch: string
  hostname: string
  nodejs: string
}> {
  return api.get('/system/info')
}

export function getSelfCheck(): Promise<{ checks: Array<{ name: string; status: string; detail: string }> }> {
  return api.get('/system/self-check')
}

export function getUpgradeInfo(): Promise<{ currentVersion: string; latestVersion: string; updateAvailable: boolean; changelog: string }> {
  return api.get('/system/upgrade')
}

export function getAbout(): Promise<{
  name: string; version: string; buildTime: string; node: string
  platform: string; arch: string; openresty: string; goVersion: null; dbType: string
}> {
  return api.get('/system/about')
}