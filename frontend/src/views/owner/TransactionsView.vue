<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <h1 class="text-2xl font-bold text-gray-900">Daftar Transaksi</h1>
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
          Transaksi
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <!-- Filters -->
      <div class="bg-white rounded-lg shadow mb-6 p-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Dari Tanggal</label>
            <input
              v-model="startDate"
              type="date"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Sampai Tanggal</label>
            <input
              v-model="endDate"
              type="date"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tipe</label>
            <select
              v-model="filterType"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            >
              <option value="">Semua</option>
              <option value="retail">Retail</option>
              <option value="wholesale">Grosir</option>
            </select>
          </div>
          <div class="flex items-end gap-2">
            <button
              @click="loadTransactions"
              class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition"
            >
              Filter
            </button>
          </div>
        </div>
      </div>

      <!-- Statistics -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-600 text-sm">Total Transaksi</p>
          <p class="text-3xl font-bold text-gray-900">{{ transactions.length }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-600 text-sm">Total Penjualan</p>
          <p class="text-3xl font-bold text-blue-600">{{ formatCurrency(totalSales) }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-600 text-sm">Rata-rata Transaksi</p>
          <p class="text-3xl font-bold text-green-600">{{ formatCurrency(avgTransaction) }}</p>
        </div>
      </div>

      <!-- Transactions Table -->
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <div v-if="loading" class="p-8 text-center">
          <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto mb-2"></div>
          <p class="text-gray-600">Memuat transaksi...</p>
        </div>

        <div v-else-if="transactions.length === 0" class="p-8 text-center">
          <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p class="text-gray-600">Tidak ada transaksi ditemukan</p>
        </div>

        <table v-else class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">No. Ref</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Tanggal</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Kasir</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Tipe</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Item</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Total</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Metode</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="tx in transactions" :key="tx.id" class="hover:bg-gray-50 transition">
              <td class="px-6 py-4 text-sm font-medium text-blue-600">{{ tx.no_ref }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ formatDate(tx.transaction_date) }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">-</td>
              <td class="px-6 py-4 text-sm">
                <span
                  :class="[
                    'px-2 py-1 rounded-full text-xs font-medium',
                    tx.type === 'retail'
                      ? 'bg-blue-100 text-blue-800'
                      : 'bg-purple-100 text-purple-800',
                  ]"
                >
                  {{ tx.type }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ tx.total_items }}</td>
              <td class="px-6 py-4 text-sm font-medium text-gray-900">{{ formatCurrency(tx.total_amount) }}</td>
              <td class="px-6 py-4 text-sm">
                <span :class="[
                  'px-2 py-1 rounded text-xs font-medium',
                  tx.payment_method === 'cash' ? 'bg-green-100 text-green-800' : 'bg-orange-100 text-orange-800'
                ]">
                  {{ tx.payment_method === 'cash' ? 'Tunai' : 'Transfer' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <span
                  :class="[
                    'px-2 py-1 rounded-full text-xs font-medium',
                    tx.payment_status === 'completed'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-yellow-100 text-yellow-800',
                  ]"
                >
                  {{ tx.payment_status === 'completed' ? 'Selesai' : 'Pending' }}
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
import * as transactionApi from '@/api/transactions'
import type { Transaction } from '@/types/api'

const transactions = ref<Transaction[]>([])
const loading = ref(false)
const startDate = ref(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0])
const endDate = ref(new Date().toISOString().split('T')[0])
const filterType = ref('')

const totalSales = computed(() => {
  return transactions.value.reduce((sum, tx) => sum + tx.total_amount, 0)
})

const avgTransaction = computed(() => {
  return transactions.value.length > 0 ? totalSales.value / transactions.value.length : 0
})

onMounted(() => {
  loadTransactions()
})

const loadTransactions = async () => {
  loading.value = true
  try {
    const response = await transactionApi.getTransactionsByDateRange(startDate.value, endDate.value, 1000, 0)
    let data = response.data.data

    if (filterType.value) {
      data = data.filter((tx: Transaction) => tx.type === filterType.value)
    }

    transactions.value = data
  } catch (err) {
    console.error('Failed to load transactions:', err)
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

const formatDate = (dateString: string) => {
  return new Intl.DateTimeFormat('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(dateString))
}
</script>
