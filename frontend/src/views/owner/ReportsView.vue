<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <h1 class="text-2xl font-bold text-gray-900">Laporan</h1>
      </div>
    </div>

    <!-- Navigation -->
    <div class="bg-white border-b border-gray-200 sticky top-16 z-30">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <button
          @click="$router.push('/owner')"
          class="px-4 py-4 font-medium text-sm border-b-2 border-transparent text-gray-600 hover:text-gray-900 transition"
        >
          Dashboard
        </button>
        <button
          class="px-4 py-4 font-medium text-sm border-b-2 border-blue-600 text-blue-600"
        >
          Laporan
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="bg-white rounded-lg shadow p-8 text-center">
        <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <h2 class="text-xl font-bold text-gray-900 mb-2">Halaman Laporan</h2>
        <p class="text-gray-600">Fitur laporan dan analitik sedang dalam pengembangan</p>
      </div>
    </div>
  </div>
</template>
