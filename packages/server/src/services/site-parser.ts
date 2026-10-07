// Parses nginx config to extract site/server-block information
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

export interface ParsedLocation {
  path: string
  proxyPass?: string
  root?: string
  return?: string
  luaDirectives: LuaDirective[]
}

export interface LuaDirective {
  type: string  // e.g. 'content_by_lua_block', 'access_by_lua_file'
  value: string
}

// Parse a single server block from raw config
export function parseSites(configContent: string, filePath: string): ParsedSite[] {
  const sites: ParsedSite[] = []

  // Simple regex-based parser for server blocks
  const serverBlockRegex = /servers*{([^}]*(?:{[^}]*}[^}]*)*)}/gs
  let match

  while ((match = serverBlockRegex.exec(configContent)) !== null) {
    const block = match[1]
    const serverName = extractValues(block, /server_names+([^;]+);/)
    const listen = extractValues(block, /listens+([^;]+);/)
    const ssl = block.includes('ssl on') || listen.some(l => l.includes('ssl'))
    const http2 = listen.some(l => l.includes('http2'))
    const root = extractSingle(block, /roots+([^;]+);/)
    const accessLog = extractSingle(block, /access_logs+([^;]+);/)
    const errorLog = extractSingle(block, /error_logs+([^;]+);/)
    const proxyPass = extractSingle(block, /proxy_passs+([^;]+);/)

    // Extract lua directives
    const luaDirectives = extractLuaDirectives(block)

    // Extract locations
    const locations = parseLocations(block)

    const defaultName = serverName.length > 0 ? serverName[0] : 'unnamed'

    sites.push({
      name: defaultName,
      filePath,
      serverName,
      listen,
      ssl,
      http2,
      locations,
      accessLog,
      errorLog,
      root,
      proxyPass,
      luaDirectives,
      content: block,
    })
  }

  return sites
}

function extractValues(content: string, regex: RegExp): string[] {
  const match = content.match(regex)
  if (!match) return []
  return match[1].split(/s+/).filter(Boolean)
}

function extractSingle(content: string, regex: RegExp): string | undefined {
  const match = content.match(regex)
  return match?.[1]?.trim()
}

const LUA_DIRECTIVES = [
  'init_by_lua', 'init_by_lua_block', 'init_by_lua_file',
  'init_worker_by_lua', 'init_worker_by_lua_block', 'init_worker_by_lua_file',
  'content_by_lua', 'content_by_lua_block', 'content_by_lua_file',
  'access_by_lua', 'access_by_lua_block', 'access_by_lua_file',
  'rewrite_by_lua', 'rewrite_by_lua_block', 'rewrite_by_lua_file',
  'log_by_lua', 'log_by_lua_block', 'log_by_lua_file',
  'header_filter_by_lua', 'header_filter_by_lua_block', 'header_filter_by_lua_file',
  'body_filter_by_lua', 'body_filter_by_lua_block', 'body_filter_by_lua_file',
  'set_by_lua', 'set_by_lua_block', 'set_by_lua_file',
  'balancer_by_lua', 'balancer_by_lua_block', 'balancer_by_lua_file',
  'ssl_certificate_by_lua', 'ssl_certificate_by_lua_block', 'ssl_certificate_by_lua_file',
  'server_rewrite_by_lua',
]

function extractLuaDirectives(content: string): LuaDirective[] {
  const directives: LuaDirective[] = []
  for (const dir of LUA_DIRECTIVES) {
    const regex = new RegExp(`${dir}_block\\s*((?:[^}]|\\}[^}])*)\\}`, 's')
    const m = content.match(regex)
    if (m) {
      directives.push({ type: dir + '_block', value: m[1].trim() })
      continue
    }
    const fileRegex = new RegExp(`${dir}_file\\s+([^;]+);`)
    const fm = content.match(fileRegex)
    if (fm) {
      directives.push({ type: dir + '_file', value: fm[1].trim() })
    }
  }
  return directives
}

function parseLocations(serverBlock: string): ParsedLocation[] {
  const locations: ParsedLocation[] = []
  const locRegex = /locations+(S+)s*{([^}]*(?:{[^}]*}[^}]*)*)}/gs
  let match

  while ((match = locRegex.exec(serverBlock)) !== null) {
    const path = match[1].trim()
    const block = match[2]
    locations.push({
      path,
      proxyPass: extractSingle(block, /proxy_passs+([^;]+);/),
      root: extractSingle(block, /roots+([^;]+);/),
      return: extractSingle(block, /returns+([^;]+);/),
      luaDirectives: extractLuaDirectives(block),
    })
  }

  return locations
}
