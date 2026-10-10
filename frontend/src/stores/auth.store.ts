import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isOwner = computed(() => user.value?.role === 'owner')
  const isKasir = computed(() => user.value?.role === 'kasir')

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setRefreshToken = (newRefreshToken: string) => {
    refreshToken.value = newRefreshToken
    localStorage.setItem('refreshToken', newRefreshToken)
  }

  const setUser = (newUser: User) => {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const logout = () => {
    token.value = null
    refreshToken.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

  // Load user from localStorage on init
  const initializeAuth = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      user.value = JSON.parse(storedUser)
    }
  }

  return {
    token,
    refreshToken,
    user,
    isAuthenticated,
    isOwner,
    isKasir,
    setToken,
    setRefreshToken,
    setUser,
    logout,
    initializeAuth,
  }
})
