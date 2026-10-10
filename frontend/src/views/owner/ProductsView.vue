<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900">Manajemen Produk</h1>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition"
        >
          + Tambah Produk
        </button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="bg-white border-b border-gray-200 sticky top-16 z-30">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex gap-0">
          <button
            @click="$router.push('/owner')"
            class="px-4 py-4 font-medium text-sm border-b-2 border-transparent text-gray-600 hover:text-gray-900 transition"
          >
            Dashboard
          </button>
          <button
            class="px-4 py-4 font-medium text-sm border-b-2 border-blue-600 text-blue-600"
          >
            Produk
          </button>
          <button
            @click="$router.push('/owner/inventory')"
            class="px-4 py-4 font-medium text-sm border-b-2 border-transparent text-gray-600 hover:text-gray-900 transition"
          >
            Inventori
          </button>
          <button
            @click="$router.push('/owner/transactions')"
            class="px-4 py-4 font-medium text-sm border-b-2 border-transparent text-gray-600 hover:text-gray-900 transition"
          >
            Transaksi
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="bg-white rounded-lg shadow p-8 text-center">
        <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
        <h2 class="text-xl font-bold text-gray-900 mb-2">Halaman Produk</h2>
        <p class="text-gray-600">Fitur manajemen produk sedang dalam pengembangan</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const openCreateModal = () => {
  // TODO: Implement product creation modal
}
</script>
