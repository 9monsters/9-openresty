import { execSync, exec } from 'node:child_process'
import { readFileSync, existsSync, readdirSync, statSync, writeFileSync, mkdirSync } from 'node:fs'
import { join, dirname, basename, relative } from 'node:path'
import { promisify } from 'node:util'

const execAsync = promisify(exec)

// Configurable paths - these would come from settings in production
let nginxBin = '/usr/local/openresty/nginx/sbin/nginx'
let nginxConfDir = '/usr/local/openresty/nginx/conf'
let nginxConfPath = '/usr/local/openresty/nginx/conf/nginx.conf'
let nginxLogDir = '/usr/local/openresty/nginx/logs'
let nginxPIDPath = '/usr/local/openresty/nginx/logs/nginx.pid'

export function updateNginxPaths(paths: {
  nginxBin?: string
  nginxConfDir?: string
  nginxConfPath?: string
  nginxLogDir?: string
  nginxPIDPath?: string
}) {
  if (paths.nginxBin) nginxBin = paths.nginxBin
  if (paths.nginxConfDir) nginxConfDir = paths.nginxConfDir
  if (paths.nginxConfPath) nginxConfPath = paths.nginxConfPath
  if (paths.nginxLogDir) nginxLogDir = paths.nginxLogDir
  if (paths.nginxPIDPath) nginxPIDPath = paths.nginxPIDPath
}

export function getNginxPaths() {
  return { nginxBin, nginxConfDir, nginxConfPath, nginxLogDir, nginxPIDPath }
}

// Execute a nginx command
async function runNginx(args: string): Promise<{ stdout: string; stderr: string }> {
  try {
    const { stdout, stderr } = await execAsync(`${nginxBin} ${args}`, { timeout: 30000 })
    return { stdout, stderr }
  } catch (err: any) {
    return { stdout: err.stdout || '', stderr: err.stderr || err.message }
  }
}

// Get nginx version info
export async function getNginxVersion(): Promise<{ version: string; openresty: boolean }> {
  const { stdout } = await runNginx('-v 2>&1')
  const isOpenResty = stdout.toLowerCase().includes('openresty')
  return {
    version: stdout.trim(),
    openresty: isOpenResty
  }
}

// Test configuration syntax
export async function testConfig(): Promise<{ ok: boolean; output: string }> {
  const { stdout, stderr } = await runNginx('-t 2>&1')
  const ok = stderr.includes('successful') || stdout.includes('successful')
  return { ok, output: stderr || stdout }
}

// Reload nginx
export async function reloadNginx(): Promise<{ ok: boolean; output: string }> {
  const { stdout, stderr } = await runNginx('-s reload 2>&1')
  return { ok: !stderr || stderr.includes('signal'), output: stderr || stdout }
}

// Get nginx status
export function getNginxStatus(): 'running' | 'stopped' | 'unknown' {
  try {
    if (!existsSync(nginxPIDPath)) return 'stopped'
    const pid = readFileSync(nginxPIDPath, 'utf-8').trim()
    try {
      process.kill(parseInt(pid), 0)
      return 'running'
    } catch {
      return 'stopped'
    }
  } catch {
    return 'unknown'
  }
}

// List config files recursively
export function listConfigFiles(basePath?: string): { name: string; path: string; type: 'file' | 'directory'; size?: number; mtime?: string }[] {
  const root = basePath || nginxConfDir
  if (!existsSync(root)) return []

  const items: any[] = []
  try {
    const entries = readdirSync(root, { withFileTypes: true })
    for (const entry of entries) {
      // Skip backup files and hidden files
      if (entry.name.startsWith('.') || entry.name.endsWith('~') || entry.name.endsWith('.bak')) continue
      const fullPath = join(root, entry.name)
      const relPath = relative(nginxConfDir, fullPath)
      if (entry.isDirectory()) {
        items.push({ name: entry.name, path: relPath, type: 'directory' })
      } else if (entry.isFile() && entry.name.endsWith('.conf')) {
        const stat = statSync(fullPath)
        items.push({ name: entry.name, path: relPath, type: 'file', size: stat.size, mtime: stat.mtime.toISOString() })
      }
    }
  } catch (e) {
    // Directory doesn't exist or no permission
  }
  return items
}

// Read a config file
export function readConfigFile(relativePath: string): string | null {
  const fullPath = join(nginxConfDir, relativePath)
  // Prevent directory traversal
  if (!fullPath.startsWith(nginxConfDir)) return null
  try {
    return readFileSync(fullPath, 'utf-8')
  } catch {
    return null
  }
}

// Write a config file
export function writeConfigFile(relativePath: string, content: string): { ok: boolean; error?: string } {
  const fullPath = join(nginxConfDir, relativePath)
  if (!fullPath.startsWith(nginxConfDir)) return { ok: false, error: 'Invalid path' }
  try {
    const dir = dirname(fullPath)
    if (!existsSync(dir)) mkdirSync(dir, { recursive: true })
    writeFileSync(fullPath, content, 'utf-8')
    return { ok: true }
  } catch (err: any) {
    return { ok: false, error: err.message }
  }
}

// Delete a config file
export function deleteConfigFile(relativePath: string): { ok: boolean; error?: string } {
  const fullPath = join(nginxConfDir, relativePath)
  if (!fullPath.startsWith(nginxConfDir)) return { ok: false, error: 'Invalid path' }
  try {
    const { unlinkSync, rmdirSync } = require('node:fs')
    const stat = statSync(fullPath)
    if (stat.isDirectory()) {
      rmdirSync(fullPath)
    } else {
      unlinkSync(fullPath)
    }
    return { ok: true }
  } catch (err: any) {
    return { ok: false, error: err.message }
  }
}

// Read log file (tail)
export function readLogFile(logType: 'access' | 'error', lines = 200): string {
  const logPath = logType === 'access'
    ? join(nginxLogDir, 'access.log')
    : join(nginxLogDir, 'error.log')
  try {
    const content = readFileSync(logPath, 'utf-8')
    const allLines = content.split('\n').filter(l => l)
    return allLines.slice(-lines).join('\n')
  } catch {
    return ''
  }
}

// Get simple nginx stats via stub_status
export async function getNginxStats(): Promise<{
  activeConnections: number
  accepts: number
  handled: number
  requests: number
  reading: number
  writing: number
  waiting: number
} | null> {
  // This would normally parse stub_status output
  // For now return mock data that the frontend can work with
  return {
    activeConnections: Math.floor(Math.random() * 100) + 10,
    accepts: Math.floor(Math.random() * 10000) + 1000,
    handled: Math.floor(Math.random() * 10000) + 1000,
    requests: Math.floor(Math.random() * 50000) + 5000,
    reading: Math.floor(Math.random() * 5),
    writing: Math.floor(Math.random() * 10) + 1,
    waiting: Math.floor(Math.random() * 80) + 5,
  }
}
