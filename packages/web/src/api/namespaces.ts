import api from './client'
export function listNamespaces() { return api.get('/namespaces') }