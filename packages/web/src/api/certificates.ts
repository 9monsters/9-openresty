import api from './client'

export function listCertificates(): Promise<{ certificates: any[] }> {
  return api.get('/certificates')
}

export function getCertificate(id: number): Promise<{ certificate: any }> {
  return api.get(`/certificates/${id}`)
}

export function createCertificate(data: { name: string; domains: string[]; certPEM?: string; keyPEM?: string; issuer?: string }): Promise<{ ok: boolean; certificate: any }> {
  return api.post('/certificates', data)
}

export function deleteCertificate(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/certificates/${id}`)
}