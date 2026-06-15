import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type AuthUser, type TokenPair } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUser | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))

  const isAuthenticated = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.roles?.includes('admin') || user.value?.roles?.includes('super_admin'))

  function saveTokens(tokens: TokenPair) {
    accessToken.value = tokens.access_token
    refreshToken.value = tokens.refresh_token
    localStorage.setItem('access_token', tokens.access_token)
    localStorage.setItem('refresh_token', tokens.refresh_token)
  }

  function clearSession() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  async function register(name: string, email: string, password: string) {
    const res = await authApi.register({ name, email, password })
    saveTokens(res.data.data.tokens)
    user.value = res.data.data.user
    return res.data.data.user
  }

  async function login(email: string, password: string) {
    const res = await authApi.login({ email, password })
    saveTokens(res.data.data.tokens)
    user.value = res.data.data.user
    return res.data.data.user
  }

  async function logout() {
    if (refreshToken.value) {
      try { await authApi.logout(refreshToken.value) } catch { /* ignore */ }
    }
    clearSession()
  }

  async function refreshTokens() {
    if (!refreshToken.value) throw new Error('No refresh token')
    const res = await authApi.refresh(refreshToken.value)
    saveTokens(res.data.data)
  }

  async function fetchMe() {
    const res = await authApi.me()
    user.value = res.data.data
    return user.value
  }

  return {
    user,
    accessToken,
    isAuthenticated,
    isAdmin,
    register,
    login,
    logout,
    refreshTokens,
    fetchMe,
    clearSession,
    saveTokens,
  }
})
