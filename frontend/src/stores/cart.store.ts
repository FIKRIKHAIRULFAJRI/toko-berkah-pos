import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Product } from '@/types/api'

export interface CartItem {
  product: Product
  qty: number
  unitPrice: number
  subtotal: number
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])

  const totalItems = computed(() => {
    return items.value.reduce((sum, item) => sum + item.qty, 0)
  })

  const subtotal = computed(() => {
    return items.value.reduce((sum, item) => sum + item.subtotal, 0)
  })

  const totalAmount = computed(() => subtotal.value)

  const isEmpty = computed(() => items.value.length === 0)

  // Add item to cart
  const addItem = (product: Product, qty: number) => {
    const existingItem = items.value.find(item => item.product.id === product.id)

    if (existingItem) {
      existingItem.qty += qty
      existingItem.subtotal = existingItem.qty * existingItem.unitPrice
    } else {
      // Determine price based on quantity
      const unitPrice = qty >= 10 ? product.price_wholesale : product.price_retail

      items.value.push({
        product,
        qty,
        unitPrice,
        subtotal: qty * unitPrice,
      })
    }
  }

  // Update item quantity
  const updateQty = (productId: string, qty: number) => {
    const item = items.value.find(i => i.product.id === productId)
    if (item) {
      if (qty <= 0) {
        removeItem(productId)
      } else {
        item.qty = qty
        item.subtotal = qty * item.unitPrice
      }
    }
  }

  // Remove item from cart
  const removeItem = (productId: string) => {
    items.value = items.value.filter(item => item.product.id !== productId)
  }

  // Clear cart
  const clearCart = () => {
    items.value = []
  }

  return {
    items,
    totalItems,
    subtotal,
    totalAmount,
    isEmpty,
    addItem,
    updateQty,
    removeItem,
    clearCart,
  }
})
