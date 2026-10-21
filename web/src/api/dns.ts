import api from './client'

export function listDnsCredentials(): Promise<{ credentials: any[] }> {
  return api.get('/dns/credentials')
}

export function createDnsCredential(data: any): Promise<{ ok: boolean; credential: any }> {
  return api.post('/dns/credentials', data)
}

export function updateDnsCredential(id: number, data: any): Promise<{ ok: boolean; credential: any }> {
  return api.put(`/dns/credentials/${id}`, data)
}

export function deleteDnsCredential(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/dns/credentials/${id}`)
}

export function listDnsDomains(): Promise<{ domains: any[] }> {
  return api.get('/dns/domains')
}

export function createDnsDomain(data: any): Promise<{ ok: boolean; domain: any }> {
  return api.post('/dns/domains', data)
}

export function updateDnsDomain(id: number, data: any): Promise<{ ok: boolean; domain: any }> {
  return api.put(`/dns/domains/${id}`, data)
}

export function deleteDnsDomain(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/dns/domains/${id}`)
}

export function listDnsGroups(): Promise<{ groups: any[] }> {
  return api.get('/dns/groups')
}

export function createDnsGroup(data: any): Promise<{ ok: boolean; group: any }> {
  return api.post('/dns/groups', data)
}

export function updateDnsGroup(id: number, data: any): Promise<{ ok: boolean; group: any }> {
  return api.put(`/dns/groups/${id}`, data)
}

export function deleteDnsGroup(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/dns/groups/${id}`)
}

export function listDdnsEntries(): Promise<{ ddns: any[] }> {
  return api.get('/dns/ddns')
}

export function createDdnsEntry(data: any): Promise<{ ok: boolean; ddns: any }> {
  return api.post('/dns/ddns', data)
}

export function updateDdnsEntry(id: number, data: any): Promise<{ ok: boolean; ddns: any }> {
  return api.put(`/dns/ddns/${id}`, data)
}

export function deleteDdnsEntry(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/dns/ddns/${id}`)
}