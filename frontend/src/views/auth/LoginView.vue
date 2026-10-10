<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-600 to-blue-800 flex items-center justify-center p-4">
    <!-- Login Card -->
    <div class="w-full max-w-md">
      <!-- Logo Section -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-white mb-2">Toko Berkah</h1>
        <p class="text-blue-100">Point of Sale System</p>
      </div>

      <!-- Login Form -->
      <div class="bg-white rounded-lg shadow-2xl p-8">
        <!-- Title -->
        <h2 class="text-2xl font-bold text-gray-800 mb-6">Login</h2>

        <!-- Error Alert -->
        <div v-if="error" class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-red-600 text-sm">{{ error }}</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
          <!-- Username -->
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 mb-1">
              Username
            </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              placeholder="Masukkan username"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
              :disabled="loading"
            />
          </div>

          <!-- Password -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              placeholder="Masukkan password"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
              :disabled="loading"
            />
          </div>

          <!-- Login Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full mt-6 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-medium py-2 px-4 rounded-lg transition duration-200"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Sedang login...
            </span>
            <span v-else>Login</span>
          </button>
        </form>

        <!-- Demo Credentials -->
        <div class="mt-6 pt-6 border-t border-gray-200">
          <p class="text-xs text-gray-600 text-center mb-3">Demo Credentials:</p>
          <div class="grid grid-cols-2 gap-2 text-xs">
            <button
              type="button"
              @click="useDemoOwner"
              class="p-2 bg-gray-100 hover:bg-gray-200 rounded text-gray-700 transition"
            >
              <span class="font-medium">Owner</span>
              <br />
              owner / password123
            </button>
            <button
              type="button"
              @click="useDemoKasir"
              class="p-2 bg-gray-100 hover:bg-gray-200 rounded text-gray-700 transition"
            >
              <span class="font-medium">Kasir</span>
              <br />
              kasir / password123
            </button>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <p class="text-center text-blue-100 text-sm mt-6">
        © 2026 Toko Berkah. All rights reserved.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import * as authApi from '@/api/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    error.value = 'Username dan password harus diisi'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await authApi.login(form.value.username, form.value.password)
    const data = response.data.data

    // Store tokens and user
    authStore.setToken(data.access_token)
    authStore.setRefreshToken(data.refresh_token)
    authStore.setUser(data.user)

    // Redirect based on role
    if (data.user.role === 'owner') {
      router.push('/owner')
    } else {
      router.push('/kasir')
    }
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Login gagal. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}

const useDemoOwner = () => {
  form.value.username = 'owner'
  form.value.password = 'password123'
  handleLogin()
}

const useDemoKasir = () => {
  form.value.username = 'kasir'
  form.value.password = 'password123'
  handleLogin()
}
</script>
