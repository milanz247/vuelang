import client from './client'

export interface LoginPayload { email: string; password: string }
export interface RegisterPayload { name: string; email: string; password: string }
export interface ForgotPayload { email: string }
export interface ResetPayload { token: string; password: string }

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_at: string
}

export interface AuthUser {
  id: number
  name: string
  email: string
  roles: string[]
  is_active: boolean
  email_verified_at: string | null
  created_at: string
}

export const authApi = {
  register: (data: RegisterPayload) =>
    client.post<{ data: { user: AuthUser; tokens: TokenPair } }>('/auth/register', data),

  login: (data: LoginPayload) =>
    client.post<{ data: { user: AuthUser; tokens: TokenPair } }>('/auth/login', data),

  logout: (refresh_token: string) =>
    client.post('/auth/logout', { refresh_token }),

  refresh: (refresh_token: string) =>
    client.post<{ data: TokenPair }>('/auth/refresh', { refresh_token }),

  me: () => client.get<{ data: AuthUser }>('/auth/me'),

  forgotPassword: (data: ForgotPayload) =>
    client.post('/auth/forgot-password', data),

  resetPassword: (data: ResetPayload) =>
    client.post('/auth/reset-password', data),
}
