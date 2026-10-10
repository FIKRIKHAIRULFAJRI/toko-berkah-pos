import http from './http'
import type { Inventory, AdjustInventoryRequest, PaginatedResponse } from '@/types/api'

export const getInventory = (productID?: string) => {
  return http.get<Inventory | PaginatedResponse>('/inventory', {
    params: productID ? { productID } : undefined,
  })
}

export const getAllInventory = (limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/inventory/all', {
    params: { limit, offset },
  })
}

export const getLowStockProducts = (limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/inventory/low-stock', {
    params: { limit, offset },
  })
}

export const adjustInventory = (data: AdjustInventoryRequest) => {
  return http.post('/inventory/adjustment', data)
}

export const getInventoryLogs = (limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/inventory/logs', {
    params: { limit, offset },
  })
}

export const getInventoryLogsByProduct = (productID: string, limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/inventory/logs/product', {
    params: { productID, limit, offset },
  })
}
