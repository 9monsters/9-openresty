import api from './client'
export function listSnippets() { return api.get('/snippets') }
export function listTemplates() { return api.get('/snippets/templates') }