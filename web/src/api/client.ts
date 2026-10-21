import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('openresty-token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('openresty-token')
      localStorage.removeItem('openresty-user')
      window.location.hash = '#/login'
    }
    const msg = err.response?.data?.error || err.message || 'Request failed'
    return Promise.reject(new Error(msg))
  },
)

export default api