import http from './http'
import type { Transaction, CreateTransactionRequest, PaginatedResponse } from '@/types/api'

export const createTransaction = (data: CreateTransactionRequest) => {
  return http.post<Transaction>('/transactions', data)
}

export const getTransactions = (limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/transactions', {
    params: { limit, offset },
  })
}

export const getTransactionById = (id: string) => {
  return http.get<Transaction>(`/transactions/${id}`)
}

export const getTransactionsByDateRange = (startDate: string, endDate: string, limit = 20, offset = 0) => {
  return http.get<PaginatedResponse>('/transactions/date-range', {
    params: { start_date: startDate, end_date: endDate, limit, offset },
  })
}
