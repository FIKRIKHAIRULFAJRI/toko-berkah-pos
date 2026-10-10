import http from './http'
import type { Product, CreateProductRequest, UpdateProductRequest, PaginatedResponse } from '@/types/api'

export const getProducts = (limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/products', {
    params: { limit, offset },
  })
}

export const getProductById = (id: string) => {
  return http.get<Product>(`/products/${id}`)
}

export const getProductBySKU = (sku: string) => {
  return http.get<Product>('/products/sku', {
    params: { sku },
  })
}

export const searchProducts = (query: string, limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/products/search', {
    params: { q: query, limit, offset },
  })
}

export const createProduct = (data: CreateProductRequest) => {
  return http.post<Product>('/products', data)
}

export const updateProduct = (id: string, data: UpdateProductRequest) => {
  return http.put<Product>(`/products/${id}`, data)
}
