<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Kasir</h1>
          <p class="text-sm text-gray-600">{{ authStore.user?.username }}</p>
        </div>
        <button
          @click="logout"
          class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition"
        >
          Logout
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Products Section -->
      <div class="flex-1 overflow-auto">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <!-- Search Bar -->
          <div class="mb-6">
            <div class="flex gap-2">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Cari produk (nama/SKU) atau scan barcode..."
                class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                @keydown.enter="handleSearch"
              />
              <button
                @click="handleSearch"
                class="px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition"
              >
                Cari
              </button>
            </div>
          </div>

          <!-- Products Grid -->
          <div v-if="!loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            <div
              v-for="product in displayedProducts"
              :key="product.id"
              class="bg-white rounded-lg shadow hover:shadow-lg transition cursor-pointer overflow-hidden"
              @click="selectProduct(product)"
            >
              <!-- Product Image Placeholder -->
              <div class="h-32 bg-gradient-to-br from-gray-200 to-gray-300 flex items-center justify-center">
                <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
              </div>

              <!-- Product Info -->
              <div class="p-3">
                <h3 class="font-medium text-gray-900 truncate text-sm">{{ product.name }}</h3>
                <p class="text-xs text-gray-600 mb-2">{{ product.sku }}</p>
                <div class="flex justify-between items-end">
                  <div>
                    <p class="text-xs text-gray-600">Retail</p>
                    <p class="font-bold text-gray-900">{{ formatCurrency(product.price_retail) }}</p>
                  </div>
                  <button
                    @click.stop="addToCart(product)"
                    class="px-2 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-xs font-medium transition"
                  >
                    +
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Loading State -->
          <div v-else class="flex items-center justify-center h-64">
            <div class="text-center">
              <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto mb-2"></div>
              <p class="text-gray-600">Memuat produk...</p>
            </div>
          </div>

          <!-- No Products -->
          <div v-if="displayedProducts.length === 0 && !loading" class="text-center py-12">
            <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
            <p class="text-gray-600">Tidak ada produk ditemukan</p>
          </div>
        </div>
      </div>

      <!-- Cart Sidebar -->
      <div class="w-80 bg-white border-l border-gray-200 flex flex-col overflow-hidden">
        <!-- Cart Header -->
        <div class="px-4 py-4 border-b border-gray-200">
          <h2 class="text-lg font-bold text-gray-900">Keranjang</h2>
          <p class="text-sm text-gray-600">{{ cartStore.totalItems }} item</p>
        </div>

        <!-- Cart Items -->
        <div class="flex-1 overflow-auto">
          <div v-if="!cartStore.isEmpty" class="divide-y">
            <div
              v-for="item in cartStore.items"
              :key="item.product.id"
              class="p-4 hover:bg-gray-50 transition"
            >
              <div class="mb-2">
                <h3 class="font-medium text-gray-900 text-sm">{{ item.product.name }}</h3>
                <p class="text-xs text-gray-600">{{ item.product.sku }}</p>
              </div>

              <div class="flex items-center gap-2 mb-2">
                <button
                  @click="decrementQty(item.product.id)"
                  class="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded text-xs font-medium transition"
                >
                  −
                </button>
                <input
                  :value="item.qty"
                  type="number"
                  min="1"
                  class="w-12 px-2 py-1 text-center border border-gray-300 rounded text-xs"
                  @change="(e) => updateQty(item.product.id, parseInt((e.target as HTMLInputElement).value))"
                />
                <button
                  @click="incrementQty(item.product.id)"
                  class="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded text-xs font-medium transition"
                >
                  +
                </button>
                <span class="ml-auto text-xs text-gray-600">{{ item.product.unit }}</span>
              </div>

              <div class="flex justify-between items-center mb-2">
                <span class="text-xs text-gray-600">{{ formatCurrency(item.unitPrice) }} x {{ item.qty }}</span>
                <span class="font-medium text-gray-900">{{ formatCurrency(item.subtotal) }}</span>
              </div>

              <button
                @click="removeFromCart(item.product.id)"
                class="w-full px-2 py-1 text-xs text-red-600 hover:bg-red-50 rounded transition"
              >
                Hapus
              </button>
            </div>
          </div>

          <div v-else class="flex items-center justify-center h-full">
            <div class="text-center text-gray-500">
              <svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
              </svg>
              <p class="text-sm">Keranjang kosong</p>
            </div>
          </div>
        </div>

        <!-- Cart Summary -->
        <div class="border-t border-gray-200 p-4 space-y-3 bg-gray-50">
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Subtotal</span>
            <span class="font-medium">{{ formatCurrency(cartStore.subtotal) }}</span>
          </div>

          <div class="flex justify-between text-lg font-bold border-t pt-3">
            <span>Total</span>
            <span class="text-blue-600">{{ formatCurrency(cartStore.totalAmount) }}</span>
          </div>

          <!-- Checkout Section -->
          <div class="space-y-2 pt-2">
            <div>
              <label class="block text-xs font-medium text-gray-700 mb-1">Metode Pembayaran</label>
              <select
                v-model="paymentMethod"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              >
                <option value="cash">Tunai</option>
                <option value="transfer">Transfer</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-medium text-gray-700 mb-1">Jumlah Bayar</label>
              <input
                v-model.number="amountPaid"
                type="number"
                placeholder="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
            </div>

            <div v-if="amountPaid > 0" class="bg-blue-50 p-2 rounded">
              <div class="flex justify-between text-xs mb-1">
                <span class="text-gray-600">Total</span>
                <span class="font-medium">{{ formatCurrency(cartStore.totalAmount) }}</span>
              </div>
              <div class="flex justify-between text-xs">
                <span class="text-gray-600">Kembalian</span>
                <span :class="['font-medium', amountPaid >= cartStore.totalAmount ? 'text-green-600' : 'text-red-600']">
                  {{ formatCurrency(amountPaid - cartStore.totalAmount) }}
                </span>
              </div>
            </div>

            <button
              @click="checkout"
              :disabled="cartStore.isEmpty || !amountPaid || amountPaid < cartStore.totalAmount || checkoutLoading"
              class="w-full mt-3 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-300 text-white font-medium rounded-lg transition"
            >
              <span v-if="checkoutLoading" class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Proses...
              </span>
              <span v-else>Checkout</span>
            </button>

            <button
              @click="clearCart"
              :disabled="cartStore.isEmpty"
              class="w-full py-2 bg-gray-200 hover:bg-gray-300 disabled:bg-gray-100 text-gray-700 font-medium rounded-lg transition text-sm"
            >
              Bersihkan
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Product Selection Modal -->
    <div
      v-if="selectedProduct"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
      @click="selectedProduct = null"
    >
      <div class="bg-white rounded-lg p-6 max-w-sm w-full" @click.stop>
        <h2 class="text-xl font-bold mb-4">{{ selectedProduct.name }}</h2>
        <div class="space-y-4 mb-6">
          <div>
            <p class="text-sm text-gray-600">SKU</p>
            <p class="font-medium">{{ selectedProduct.sku }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-600">Unit</p>
            <p class="font-medium">{{ selectedProduct.unit }}</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-600">Harga Retail</p>
              <p class="font-medium">{{ formatCurrency(selectedProduct.price_retail) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Harga Grosir</p>
              <p class="font-medium">{{ formatCurrency(selectedProduct.price_wholesale) }}</p>
            </div>
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">Jumlah</label>
          <input
            v-model.number="addQty"
            type="number"
            min="1"
            placeholder="1"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
          />
        </div>

        <div class="flex gap-2">
          <button
            @click="confirmAddToCart"
            class="flex-1 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
          >
            Tambah
          </button>
          <button
            @click="selectedProduct = null"
            class="flex-1 py-2 bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium rounded-lg transition"
          >
            Batal
          </button>
        </div>
      </div>
    </div>

    <!-- Success Modal -->
    <div
      v-if="successModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-sm w-full text-center">
        <div class="mb-4">
          <svg class="w-16 h-16 text-green-600 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">Transaksi Berhasil!</h2>
        <p class="text-gray-600 mb-2">No. Ref: {{ successTransaction?.no_ref }}</p>
        <p class="text-sm text-gray-600 mb-6">Total: {{ formatCurrency(successTransaction?.total_amount || 0) }}</p>
        <button
          @click="closeSuccessModal"
          class="w-full py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
        >
          OK
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import { useCartStore } from '@/stores/cart.store'
import * as productApi from '@/api/products'
import * as transactionApi from '@/api/transactions'
import type { Product, Transaction } from '@/types/api'

const router = useRouter()
const authStore = useAuthStore()
const cartStore = useCartStore()

const products = ref<Product[]>([])
const loading = ref(false)
const searchQuery = ref('')
const selectedProduct = ref<Product | null>(null)
const addQty = ref(1)
const paymentMethod = ref<'cash' | 'transfer'>('cash')
const amountPaid = ref(0)
const checkoutLoading = ref(false)
const successModal = ref(false)
const successTransaction = ref<Transaction | null>(null)

const displayedProducts = computed(() => {
  if (!searchQuery.value) {
    return products.value
  }
  const query = searchQuery.value.toLowerCase()
  return products.value.filter(p =>
    p.name.toLowerCase().includes(query) || p.sku.toLowerCase().includes(query)
  )
})

onMounted(() => {
  loadProducts()
})

const loadProducts = async () => {
  loading.value = true
  try {
    const response = await productApi.getProducts(100, 0)
    products.value = response.data.data
  } catch (err) {
    console.error('Failed to load products:', err)
  } finally {
    loading.value = false
  }
}

const handleSearch = async () => {
  if (!searchQuery.value.trim()) {
    loadProducts()
    return
  }

  loading.value = true
  try {
    const response = await productApi.searchProducts(searchQuery.value, 100, 0)
    products.value = response.data.data
  } catch (err) {
    console.error('Failed to search products:', err)
  } finally {
    loading.value = false
  }
}

const selectProduct = (product: Product) => {
  selectedProduct.value = product
  addQty.value = 1
}

const confirmAddToCart = () => {
  if (selectedProduct.value) {
    cartStore.addItem(selectedProduct.value, addQty.value)
    selectedProduct.value = null
  }
}

const addToCart = (product: Product) => {
  cartStore.addItem(product, 1)
}

const incrementQty = (productId: string) => {
  const item = cartStore.items.find(i => i.product.id === productId)
  if (item) {
    cartStore.updateQty(productId, item.qty + 1)
  }
}

const decrementQty = (productId: string) => {
  const item = cartStore.items.find(i => i.product.id === productId)
  if (item && item.qty > 1) {
    cartStore.updateQty(productId, item.qty - 1)
  }
}

const updateQty = (productId: string, qty: number) => {
  cartStore.updateQty(productId, qty)
}

const removeFromCart = (productId: string) => {
  cartStore.removeItem(productId)
}

const clearCart = () => {
  cartStore.clearCart()
}

const checkout = async () => {
  if (cartStore.isEmpty || !amountPaid.value || amountPaid.value < cartStore.totalAmount) {
    return
  }

  checkoutLoading.value = true
  try {
    const items = cartStore.items.map(item => ({
      product_id: item.product.id,
      qty: item.qty,
    }))

    const response = await transactionApi.createTransaction({
      items,
      payment_method: paymentMethod.value,
      amount_paid: amountPaid.value,
    })

    successTransaction.value = response.data.data
    successModal.value = true
    cartStore.clearCart()
    amountPaid.value = 0
  } catch (err) {
    console.error('Checkout failed:', err)
  } finally {
    checkoutLoading.value = false
  }
}

const closeSuccessModal = () => {
  successModal.value = false
  successTransaction.value = null
}

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(value)
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
