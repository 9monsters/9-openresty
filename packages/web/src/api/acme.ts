import api from './client'
export function listAcmeUsers() { return api.get('/acme-users') }