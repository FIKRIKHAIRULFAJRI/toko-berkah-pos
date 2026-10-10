<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-600 to-blue-800 flex items-center justify-center p-4">
    <div class="text-center">
      <h1 class="text-9xl font-bold text-white mb-4">404</h1>
      <h2 class="text-4xl font-bold text-white mb-2">Halaman Tidak Ditemukan</h2>
      <p class="text-blue-100 mb-8">Maaf, halaman yang Anda cari tidak ada.</p>
      <router-link
        to="/"
        class="inline-block px-6 py-3 bg-white text-blue-600 font-medium rounded-lg hover:bg-blue-50 transition"
      >
        Kembali ke Beranda
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
</script>
