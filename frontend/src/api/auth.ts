import http from './http'
import type { LoginResponse, RefreshTokenResponse } from '@/types/api'

export const login = (username: string, password: string) => {
  return http.post<LoginResponse>('/auth/login', {
    username,
    password,
  })
}

export const logout = () => {
  return http.post('/auth/logout')
}

export const refreshToken = (refreshToken: string) => {
  return http.post<RefreshTokenResponse>('/auth/refresh', {
    refresh_token: refreshToken,
  })
}
