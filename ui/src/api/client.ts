import axios, { type AxiosError } from 'axios'
import { useAuthStore } from '@/stores/auth'

const client = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
  withCredentials: false,
})

// Attach access token to every request
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Refresh token on 401, then retry once
let refreshing = false
let queue: Array<() => void> = []

client.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as any
    if (error.response?.status === 401 && !original._retry) {
      if (refreshing) {
        return new Promise((resolve) => {
          queue.push(() => resolve(client(original)))
        })
      }

      original._retry = true
      refreshing = true

      try {
        const auth = useAuthStore()
        await auth.refreshTokens()
        queue.forEach((fn) => fn())
        queue = []
        return client(original)
      } catch {
        const auth = useAuthStore()
        auth.clearSession()
        window.location.href = '/login'
        return Promise.reject(error)
      } finally {
        refreshing = false
      }
    }
    return Promise.reject(error)
  },
)

export default client
