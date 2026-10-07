import api from './client'
export function listUsers() { return api.get('/users') }
export function getProfile() { return api.get('/users/profile') }