import api from './client'
export function listBackups() { return api.get('/backup') }
export function getAutoBackupSettings() { return api.get('/backup/auto-backup') }