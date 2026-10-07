import api from './client'
export function listCredentials() { return api.get('/dns/credentials') }
export function listDomains() { return api.get('/dns/domains') }
export function listGroups() { return api.get('/dns/groups') }
export function listDDNS() { return api.get('/dns/ddns') }