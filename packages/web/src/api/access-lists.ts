import api from './client'
export function listAccessLists() { return api.get('/access-lists') }