import api from './client'

export function listLuaBlocks() {
  return api.get('/lua')
}

export function listLuaLibraries() {
  return api.get('/lua/libraries')
}

export function listSharedDicts() {
  return api.get('/lua/shared-dicts')
}
