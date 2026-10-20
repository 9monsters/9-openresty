import api from './client'

export interface LuaBlock {
  filePath: string
  siteName: string
  directives: Array<{ type: string; value: string }>
}

export interface LuaLibrary {
  name: string
  description: string
  installed: boolean
}

export interface LuaSharedDict {
  name: string
  size: string
  usage: string
}

export function listLuaBlocks(): Promise<{ luaBlocks: LuaBlock[] }> {
  return api.get('/lua')
}

export function listLuaLibraries(): Promise<{ libraries: LuaLibrary[] }> {
  return api.get('/lua/libraries')
}

export function listSharedDicts(): Promise<{ dicts: LuaSharedDict[] }> {
  return api.get('/lua/shared-dicts')
}
