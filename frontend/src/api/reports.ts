import http from './http'
import type { DashboardMetrics, SalesBreakdown, BestsellerItem, PaginatedResponse } from '@/types/api'

export const getDashboard = (date?: string) => {
  return http.get<DashboardMetrics>('/reports/dashboard', {
    params: { date },
  })
}

export const getSalesBreakdown = (dateFrom: string, dateTo: string) => {
  return http.get<SalesBreakdown>('/reports/sales-breakdown', {
    params: { date_from: dateFrom, date_to: dateTo },
  })
}

export const getBestseller = (dateFrom: string, dateTo: string, limit: number = 10) => {
  return http.get<PaginatedResponse>('/reports/bestseller', {
    params: { date_from: dateFrom, date_to: dateTo, limit },
  })
}

export const getStockMovement = (dateFrom: string, dateTo: string, category?: string) => {
  return http.get<PaginatedResponse>('/reports/stock-movement', {
    params: { date_from: dateFrom, date_to: dateTo, category },
  })
}

export const exportReport = (reportType: string, dateFrom: string, dateTo: string, format: string = 'xlsx') => {
  return http.post('/reports/export', {
    report_type: reportType,
    date_from: dateFrom,
    date_to: dateTo,
    format,
  }, {
    responseType: 'blob',
  })
}
