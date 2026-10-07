export interface LuaDirective {
  type: string
  value: string
}

export interface ParsedLocation {
  path: string
  proxyPass?: string
  root?: string
  return?: string
  luaDirectives: LuaDirective[]
}

export interface ParsedSite {
  name: string
  filePath: string
  serverName: string[]
  listen: string[]
  ssl: boolean
  http2: boolean
  locations: ParsedLocation[]
  accessLog?: string
  errorLog?: string
  root?: string
  proxyPass?: string
  luaDirectives: LuaDirective[]
  content: string
}

export interface UpstreamServer {
  address: string
  weight: number
  maxFails?: number
  failTimeout?: string
  status: 'up' | 'down' | 'checking'
}

export interface Upstream {
  name: string
  strategy: string
  hashKey?: string
  servers: UpstreamServer[]
}

export interface Certificate {
  name: string
  domains: string[]
  issuer: string
  notBefore: string
  notAfter: string
  autoRenew: boolean
  status: 'valid' | 'expiring' | 'expired'
}

export interface ConfigFile {
  name: string
  path: string
  type: 'file' | 'directory'
  size?: number
  mtime?: string
}
