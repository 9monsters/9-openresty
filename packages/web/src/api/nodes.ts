import api from './client'
export function listNodes() { return api.get('/nodes') }