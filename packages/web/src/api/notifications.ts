import api from './client'
export function listNotifications() { return api.get('/notifications') }