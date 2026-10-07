import api from './client'
export function selfCheck() { return api.get('/system/self-check') }
export function checkUpgrade() { return api.get('/system/upgrade') }
export function getAbout() { return api.get('/system/about') }