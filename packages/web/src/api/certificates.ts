import api from './client'

export function listCertificates() {
  return api.get('/certificates')
}
