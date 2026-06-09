import { useQuery } from '@tanstack/react-query'
import {
  getDashboardStats,
  getKPIMetrics,
  getRevenueReport,
  getCustomerReport,
  getStaffReport,
  getServiceReport,
  getExpenseReport,
  getInventoryReport,
  getProfitLossReport,
  type DateRangeParams,
} from '@/services/analytics'

export function useDashboardStats() {
  return useQuery({
    queryKey: ['analytics', 'dashboard'],
    queryFn: getDashboardStats,
    refetchInterval: 60000,
  })
}

export function useKPIMetrics(params: DateRangeParams = {}) {
  return useQuery({
    queryKey: ['analytics', 'kpis', params],
    queryFn: () => getKPIMetrics(params),
  })
}

export function useRevenueReport(params: DateRangeParams & { group_by?: string } = {}) {
  return useQuery({
    queryKey: ['analytics', 'revenue', params],
    queryFn: () => getRevenueReport(params),
  })
}

export function useCustomerReport(params: DateRangeParams = {}) {
  return useQuery({
    queryKey: ['analytics', 'customers', params],
    queryFn: () => getCustomerReport(params),
  })
}

export function useStaffAnalyticsReport(params: DateRangeParams = {}) {
  return useQuery({
    queryKey: ['analytics', 'staff', params],
    queryFn: () => getStaffReport(params),
  })
}

export function useServiceReport(params: DateRangeParams = {}) {
  return useQuery({
    queryKey: ['analytics', 'services', params],
    queryFn: () => getServiceReport(params),
  })
}

export function useExpenseAnalyticsReport(params: DateRangeParams = {}) {
  return useQuery({
    queryKey: ['analytics', 'expenses', params],
    queryFn: () => getExpenseReport(params),
  })
}

export function useInventoryAnalyticsReport() {
  return useQuery({
    queryKey: ['analytics', 'inventory'],
    queryFn: getInventoryReport,
  })
}

export function useProfitLossReport(params: DateRangeParams & { group_by?: string } = {}) {
  return useQuery({
    queryKey: ['analytics', 'profit-loss', params],
    queryFn: () => getProfitLossReport(params),
  })
}
