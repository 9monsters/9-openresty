import { defineStore } from 'pinia'

interface User {
  id: number
  username: string
  email: string
  role: string
  status: string
  lastLogin?: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('openresty-token') || '')
  const info = ref<User | null>(JSON.parse(localStorage.getItem('openresty-user') || 'null'))

  const isLogin = computed(() => !!token.value)

  function login(t: string, user: User) {
    token.value = t
    info.value = user
    localStorage.setItem('openresty-token', t)
    localStorage.setItem('openresty-user', JSON.stringify(user))
  }

  function logout() {
    token.value = ''
    info.value = null
    localStorage.removeItem('openresty-token')
    localStorage.removeItem('openresty-user')
  }

  function updateInfo(user: Partial<User>) {
    if (info.value) {
      info.value = { ...info.value, ...user }
      localStorage.setItem('openresty-user', JSON.stringify(info.value))
    }
  }

  return { token, info, isLogin, login, logout, updateInfo }
})