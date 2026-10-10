<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Riwayat Transaksi</h1>
        </div>
        <div class="flex gap-2">
          <router-link
            to="/kasir"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition"
          >
            Kembali ke Kasir
          </router-link>
          <button
            @click="logout"
            class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition"
          >
            Logout
          </button>
        </div>
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
            <button
              @click="resetFilter"
              class="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded-lg font-medium transition"
            >
              Reset
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
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Tipe</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Item</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Total</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Metode</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="tx in transactions" :key="tx.id" class="hover:bg-gray-50 transition">
              <td class="px-6 py-4 text-sm font-medium text-blue-600">{{ tx.no_ref }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ formatDate(tx.transaction_date) }}</td>
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
              <td class="px-6 py-4 text-sm text-gray-700">{{ tx.total_items }} item</td>
              <td class="px-6 py-4 text-sm font-medium text-gray-900">{{ formatCurrency(tx.total_amount) }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">
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
                      : tx.payment_status === 'pending'
                      ? 'bg-yellow-100 text-yellow-800'
                      : 'bg-red-100 text-red-800',
                  ]"
                >
                  {{ formatStatus(tx.payment_status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <button
                  @click="viewDetail(tx)"
                  class="text-blue-600 hover:text-blue-800 font-medium"
                >
                  Lihat
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Detail Modal -->
    <div
      v-if="selectedTransaction"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
      @click="selectedTransaction = null"
    >
      <div class="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-auto" @click.stop>
        <div class="sticky top-0 bg-gradient-to-r from-blue-600 to-blue-700 px-6 py-4 flex justify-between items-center">
          <h2 class="text-xl font-bold text-white">Detail Transaksi</h2>
          <button @click="selectedTransaction = null" class="text-white hover:opacity-80">
            ✕
          </button>
        </div>

        <div class="p-6">
          <!-- Header Info -->
          <div class="grid grid-cols-2 gap-4 mb-6 pb-6 border-b">
            <div>
              <p class="text-sm text-gray-600">No. Referensi</p>
              <p class="font-bold text-gray-900">{{ selectedTransaction.no_ref }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Tanggal</p>
              <p class="font-bold text-gray-900">{{ formatDate(selectedTransaction.transaction_date) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Tipe</p>
              <p class="font-bold text-gray-900">{{ selectedTransaction.type }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Status</p>
              <p class="font-bold">
                <span
                  :class="[
                    'px-2 py-1 rounded-full text-xs font-medium',
                    selectedTransaction.payment_status === 'completed'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-yellow-100 text-yellow-800',
                  ]"
                >
                  {{ formatStatus(selectedTransaction.payment_status) }}
                </span>
              </p>
            </div>
          </div>

          <!-- Items -->
          <div class="mb-6">
            <h3 class="font-bold text-gray-900 mb-3">Item Transaksi</h3>
            <div class="space-y-2">
              <div
                v-for="item in selectedTransaction.items"
                :key="item.product_id"
                class="flex justify-between items-center p-3 bg-gray-50 rounded"
              >
                <div>
                  <p class="font-medium text-gray-900">{{ item.product_name }}</p>
                  <p class="text-xs text-gray-600">{{ item.qty }} × {{ formatCurrency(item.unit_price) }}</p>
                </div>
                <p class="font-bold text-gray-900">{{ formatCurrency(item.subtotal) }}</p>
              </div>
            </div>
          </div>

          <!-- Summary -->
          <div class="bg-gray-50 rounded-lg p-4 mb-6">
            <div class="space-y-2">
              <div class="flex justify-between text-sm">
                <span class="text-gray-600">Subtotal</span>
                <span class="font-medium">{{ formatCurrency(selectedTransaction.subtotal) }}</span>
              </div>
              <div v-if="selectedTransaction.discount_amount > 0" class="flex justify-between text-sm">
                <span class="text-gray-600">Diskon</span>
                <span class="font-medium">-{{ formatCurrency(selectedTransaction.discount_amount) }}</span>
              </div>
              <div class="flex justify-between text-lg font-bold border-t pt-2">
                <span>Total</span>
                <span class="text-blue-600">{{ formatCurrency(selectedTransaction.total_amount) }}</span>
              </div>
            </div>
          </div>

          <!-- Payment Info -->
          <div class="mb-6">
            <h3 class="font-bold text-gray-900 mb-3">Informasi Pembayaran</h3>
            <div class="bg-blue-50 p-4 rounded-lg">
              <p class="text-sm text-gray-600">Metode Pembayaran</p>
              <p class="font-bold text-gray-900 mb-4">{{ selectedTransaction.payment_method === 'cash' ? 'Tunai' : 'Transfer' }}</p>
              <p class="text-sm text-gray-600">Total Dibayar</p>
              <p class="font-bold text-lg text-blue-600">{{ formatCurrency(selectedTransaction.total_amount) }}</p>
            </div>
          </div>

          <!-- Close Button -->
          <button
            @click="selectedTransaction = null"
            class="w-full py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import * as transactionApi from '@/api/transactions'
import type { Transaction } from '@/types/api'

const router = useRouter()
const authStore = useAuthStore()

const transactions = ref<Transaction[]>([])
const loading = ref(false)
const selectedTransaction = ref<Transaction | null>(null)
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
    const response = await transactionApi.getTransactionsByDateRange(startDate.value, endDate.value, 100, 0)
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

const resetFilter = () => {
  startDate.value = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
  endDate.value = new Date().toISOString().split('T')[0]
  filterType.value = ''
  loadTransactions()
}

const viewDetail = (tx: Transaction) => {
  selectedTransaction.value = tx
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

const formatStatus = (status: string) => {
  return {
    completed: 'Selesai',
    pending: 'Menunggu',
    failed: 'Gagal',
  }[status] || status
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
