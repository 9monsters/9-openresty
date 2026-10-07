import { join } from 'node:path'
import { existsSync, mkdirSync, writeFileSync, readFileSync, readdirSync } from 'node:fs'
import { getNginxPaths } from './nginx.js'

const BACKUP_DIR = join(getNginxPaths().nginxConfDir, '.backups')

export interface BackupEntry {
  id: string
  filepath: string
  content: string
  timestamp: string
  message: string
}

// Save a backup of a config file before modification
export function backupConfig(relativePath: string, content: string, message = 'Manual backup'): string {
  if (!existsSync(BACKUP_DIR)) mkdirSync(BACKUP_DIR, { recursive: true })

  const id = Date.now().toString(36) + Math.random().toString(36).slice(2, 8)
  const safeName = relativePath.replaceAll('/', '_').replaceAll('\\', '_')
  const backupPath = join(BACKUP_DIR, `${id}_${safeName}`)

  writeFileSync(backupPath, content, 'utf-8')
  writeFileSync(backupPath + '.meta', JSON.stringify({
    id, filepath: relativePath, timestamp: new Date().toISOString(), message
  }), 'utf-8')

  return id
}

// List backups for a file
export function listBackups(relativePath: string): BackupEntry[] {
  if (!existsSync(BACKUP_DIR)) return []

  const safeName = relativePath.replaceAll('/', '_').replaceAll('\\', '_')
  const entries: BackupEntry[] = []

  for (const file of readdirSync(BACKUP_DIR)) {
    if (!file.endsWith(safeName)) continue
    const metaFile = file + '.meta'
    try {
      const meta = JSON.parse(readFileSync(join(BACKUP_DIR, metaFile), 'utf-8'))
      entries.push({
        ...meta,
        content: readFileSync(join(BACKUP_DIR, file), 'utf-8')
      })
    } catch { /* skip corrupt backups */ }
  }

  return entries.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
}