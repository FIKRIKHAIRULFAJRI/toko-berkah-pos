<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Dashboard Owner</h1>
          <p class="text-sm text-gray-600">Selamat datang, {{ authStore.user?.username }}</p>
        </div>
        <div class="flex gap-2">
          <button
            @click="logout"
            class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition"
          >
            Logout
          </button>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <div class="bg-white border-b border-gray-200 sticky top-16 z-30">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex gap-0 overflow-x-auto">
          <router-link
            to="/owner"
            :class="[
              'px-4 py-4 font-medium text-sm border-b-2 transition whitespace-nowrap',
              route.name === 'Dashboard'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900',
            ]"
          >
            Dashboard
          </router-link>
          <router-link
            to="/owner/products"
            :class="[
              'px-4 py-4 font-medium text-sm border-b-2 transition whitespace-nowrap',
              route.name === 'Products'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900',
            ]"
          >
            Produk
          </router-link>
          <router-link
            to="/owner/inventory"
            :class="[
              'px-4 py-4 font-medium text-sm border-b-2 transition whitespace-nowrap',
              route.name === 'Inventory'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900',
            ]"
          >
            Inventori
          </router-link>
          <router-link
            to="/owner/transactions"
            :class="[
              'px-4 py-4 font-medium text-sm border-b-2 transition whitespace-nowrap',
              route.name === 'Transactions'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900',
            ]"
          >
            Transaksi
          </router-link>
          <router-link
            to="/owner/reports"
            :class="[
              'px-4 py-4 font-medium text-sm border-b-2 transition whitespace-nowrap',
              route.name === 'Reports'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900',
            ]"
          >
            Laporan
          </router-link>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <!-- Date Range Filter -->
      <div class="bg-white rounded-lg shadow mb-6 p-4">
        <div class="flex gap-4 items-end">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Dari Tanggal</label>
            <input
              v-model="startDate"
              type="date"
              class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Sampai Tanggal</label>
            <input
              v-model="endDate"
              type="date"
              class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            />
          </div>
          <button
            @click="loadMetrics"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition"
          >
            Filter
          </button>
        </div>
      </div>

      <!-- KPI Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <!-- Total Penjualan -->
        <div class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg shadow p-6 text-white">
          <div class="flex justify-between items-start mb-4">
            <div>
              <p class="text-blue-100 text-sm font-medium">Total Penjualan</p>
              <p class="text-3xl font-bold mt-2">{{ formatCurrency(totalSales) }}</p>
            </div>
            <svg class="w-12 h-12 opacity-20" fill="currentColor" viewBox="0 0 20 20">
              <path d="M3 1a1 1 0 000 2h1.22l.305 1.222a.997.997 0 00.01.042L5.960 9H9a1 1 0 000 2H5.897l.017.085A2.008 2.008 0 005 13a2 2 0 104 0 .999.999 0 00-.898-.684.984.984 0 00.017-.085h4.014l.305-1.222a.997.997 0 00.01-.042L17.78 3H19a1 1 0 000-2H3z" />
            </svg>
          </div>
          <p class="text-blue-100 text-xs">Periode: {{ formatDate(startDate) }} - {{ formatDate(endDate) }}</p>
        </div>

        <!-- Total Transaksi -->
        <div class="bg-gradient-to-br from-green-500 to-green-600 rounded-lg shadow p-6 text-white">
          <div class="flex justify-between items-start mb-4">
            <div>
              <p class="text-green-100 text-sm font-medium">Total Transaksi</p>
              <p class="text-3xl font-bold mt-2">{{ totalTransactions }}</p>
            </div>
            <svg class="w-12 h-12 opacity-20" fill="currentColor" viewBox="0 0 20 20">
              <path d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" />
            </svg>
          </div>
          <p class="text-green-100 text-xs">Rata-rata: {{ formatCurrency(avgTransaction) }}/transaksi</p>
        </div>

        <!-- Total Produk -->
        <div class="bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg shadow p-6 text-white">
          <div class="flex justify-between items-start mb-4">
            <div>
              <p class="text-purple-100 text-sm font-medium">Total Produk</p>
              <p class="text-3xl font-bold mt-2">{{ totalProducts }}</p>
            </div>
            <svg class="w-12 h-12 opacity-20" fill="currentColor" viewBox="0 0 20 20">
              <path d="M3 1a1 1 0 000 2h1.22l.305 1.222a.997.997 0 00.01.042L5.960 9H9a1 1 0 000 2H5.897l.017.085A2.008 2.008 0 005 13a2 2 0 104 0 .999.999 0 00-.898-.684.984.984 0 00.017-.085h4.014l.305-1.222a.997.997 0 00.01-.042L17.78 3H19a1 1 0 000-2H3z" />
            </svg>
          </div>
          <p class="text-purple-100 text-xs">Produk aktif</p>
        </div>

        <!-- Low Stock Alert -->
        <div class="bg-gradient-to-br from-orange-500 to-orange-600 rounded-lg shadow p-6 text-white">
          <div class="flex justify-between items-start mb-4">
            <div>
              <p class="text-orange-100 text-sm font-medium">Stok Rendah</p>
              <p class="text-3xl font-bold mt-2">{{ lowStockCount }}</p>
            </div>
            <svg class="w-12 h-12 opacity-20" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
          </div>
          <p class="text-orange-100 text-xs">Butuh restocking segera</p>
        </div>
      </div>

      <!-- Charts Row -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <!-- Sales Trend -->
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-bold text-gray-900 mb-4">Tren Penjualan</h3>
          <div class="h-64 flex items-end justify-around gap-2 bg-gray-50 rounded p-4">
            <div
              v-for="(value, idx) in salesTrend"
              :key="idx"
              :style="{ height: value.percentage + '%' }"
              class="flex-1 bg-blue-500 rounded-t hover:bg-blue-600 transition relative group"
            >
              <div class="opacity-0 group-hover:opacity-100 absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap transition">
                {{ formatCurrency(value.amount) }}
              </div>
            </div>
          </div>
          <p class="text-xs text-gray-600 mt-4 text-center">Hover untuk melihat detail</p>
        </div>

        <!-- Retail vs Wholesale -->
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-bold text-gray-900 mb-4">Retail vs Grosir</h3>
          <div class="space-y-4">
            <div>
              <div class="flex justify-between items-center mb-2">
                <span class="text-sm font-medium text-gray-700">Retail</span>
                <span class="text-sm font-bold text-blue-600">{{ formatCurrency(retailTotal) }}</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div
                  class="bg-blue-500 h-2 rounded-full transition"
                  :style="{ width: retailPercentage + '%' }"
                ></div>
              </div>
            </div>

            <div>
              <div class="flex justify-between items-center mb-2">
                <span class="text-sm font-medium text-gray-700">Grosir</span>
                <span class="text-sm font-bold text-purple-600">{{ formatCurrency(wholesaleTotal) }}</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div
                  class="bg-purple-500 h-2 rounded-full transition"
                  :style="{ width: wholesalePercentage + '%' }"
                ></div>
              </div>
            </div>
          </div>

          <div class="mt-6 p-4 bg-gray-50 rounded-lg">
            <p class="text-xs text-gray-600 mb-2">Komposisi Penjualan</p>
            <div class="flex gap-4">
              <div>
                <p class="text-xs text-gray-600">Retail</p>
                <p class="text-lg font-bold text-blue-600">{{ retailPercentage }}%</p>
              </div>
              <div>
                <p class="text-xs text-gray-600">Grosir</p>
                <p class="text-lg font-bold text-purple-600">{{ wholesalePercentage }}%</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-bold text-gray-900">Transaksi Terbaru</h3>
        </div>
        <div v-if="recentTransactions.length === 0" class="p-8 text-center text-gray-500">
          Tidak ada transaksi
        </div>
        <table v-else class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">No. Ref</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Waktu</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Item</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Total</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Tipe</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="tx in recentTransactions.slice(0, 5)" :key="tx.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 text-sm font-medium text-blue-600">{{ tx.no_ref }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ formatTime(tx.transaction_date) }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ tx.total_items }} item</td>
              <td class="px-6 py-4 text-sm font-bold text-gray-900">{{ formatCurrency(tx.total_amount) }}</td>
              <td class="px-6 py-4 text-sm">
                <span
                  :class="[
                    'px-2 py-1 rounded text-xs font-medium',
                    tx.type === 'retail' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800',
                  ]"
                >
                  {{ tx.type }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import * as transactionApi from '@/api/transactions'
import * as productApi from '@/api/products'
import * as inventoryApi from '@/api/inventory'
import type { Transaction, Product, Inventory } from '@/types/api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const startDate = ref(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0])
const endDate = ref(new Date().toISOString().split('T')[0])
const transactions = ref<Transaction[]>([])
const products = ref<Product[]>([])
const inventory = ref<Inventory[]>([])
const loading = ref(false)

const totalSales = computed(() => {
  return transactions.value.reduce((sum, tx) => sum + tx.total_amount, 0)
})
const totalTransactions = computed(() => transactions.value.length)
const avgTransaction = computed(() => {
  return transactions.value.length > 0 ? totalSales.value / transactions.value.length : 0
})
const totalProducts = computed(() => products.value.length)
const lowStockCount = computed(() => {
  return inventory.value.filter(inv => inv.status === 'low').length
})
const recentTransactions = computed(() => {
  return [...transactions.value].sort((a, b) => {
    return new Date(b.transaction_date).getTime() - new Date(a.transaction_date).getTime()
  })
})
const retailTotal = computed(() => {
  return transactions.value
    .filter(tx => tx.type === 'retail')
    .reduce((sum, tx) => sum + tx.total_amount, 0)
})

const wholesaleTotal = computed(() => {
  return transactions.value
    .filter(tx => tx.type === 'wholesale')
    .reduce((sum, tx) => sum + tx.total_amount, 0)
})

const retailPercentage = computed(() => {
  const total = retailTotal.value + wholesaleTotal.value
  return total > 0 ? Math.round((retailTotal.value / total) * 100) : 0
})

const wholesalePercentage = computed(() => {
  return 100 - retailPercentage.value
})

const salesTrend = computed(() => {
  const days = 7
  const trend = []
  const maxSales = Math.max(...Array(days).fill(0).map((_, i) => {
    const date = new Date(Date.now() - i * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
    return transactions.value
      .filter(tx => tx.transaction_date.split('T')[0] === date)
      .reduce((sum, tx) => sum + tx.total_amount, 0)
  }))

  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(Date.now() - i * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
    const amount = transactions.value
      .filter(tx => tx.transaction_date.split('T')[0] === date)
      .reduce((sum, tx) => sum + tx.total_amount, 0)
    trend.push({
      amount,
      percentage: maxSales > 0 ? (amount / maxSales) * 100 : 0,
    })
  }
  return trend
})

onMounted(() => {
  loadMetrics()
})

const loadMetrics = async () => {
  loading.value = true
  try {
    const [txRes, prodRes, invRes] = await Promise.all([
      transactionApi.getTransactionsByDateRange(startDate.value, endDate.value, 1000, 0),
      productApi.getProducts(1000, 0),
      inventoryApi.getAllInventory(1000, 0),
    ])

    transactions.value = txRes.data.data
    products.value = prodRes.data.data
    inventory.value = Array.isArray(invRes.data.data) ? invRes.data.data : []
  } catch (err) {
    console.error('Failed to load metrics:', err)
  } finally {
    loading.value = false
  }
}

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(value)
}

const formatDate = (dateStr: string) => {
  return new Intl.DateTimeFormat('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(new Date(dateStr))
}

const formatTime = (dateStr: string) => {
  return new Intl.DateTimeFormat('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(dateStr))
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
