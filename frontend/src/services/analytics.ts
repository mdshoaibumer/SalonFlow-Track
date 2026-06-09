import { apiClient } from './api-client'
import type {
  DashboardStats,
  KPIMetrics,
  RevenueReport,
  CustomerReport,
  StaffAnalyticsReport,
  ServiceReport,
  ExpenseAnalyticsReport,
  InventoryAnalyticsReport,
  ProfitLossReport,
} from '@/types'

export interface DateRangeParams {
  date_from?: string
  date_to?: string
}

function buildQuery(params: object): string {
  const query = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== null && v !== '') query.set(k, String(v))
  }
  const qs = query.toString()
  return qs ? `?${qs}` : ''
}

export async function getDashboardStats(): Promise<DashboardStats> {
  const response = await apiClient.get<DashboardStats>('/reports/dashboard')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch dashboard')
  return response.data
}

export async function getKPIMetrics(params: DateRangeParams = {}): Promise<KPIMetrics> {
  const response = await apiClient.get<KPIMetrics>(`/reports/kpis${buildQuery(params)}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch KPIs')
  return response.data
}

export async function getRevenueReport(params: DateRangeParams & { group_by?: string } = {}): Promise<RevenueReport> {
  const response = await apiClient.get<RevenueReport>(`/reports/revenue${buildQuery(params)}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch revenue report')
  return response.data
}

export async function getCustomerReport(params: DateRangeParams = {}): Promise<CustomerReport> {
  const response = await apiClient.get<CustomerReport>(`/reports/customers${buildQuery(params)}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch customer report')
  return response.data
}

export async function getStaffReport(params: DateRangeParams = {}): Promise<StaffAnalyticsReport> {
  const response = await apiClient.get<StaffAnalyticsReport>(`/reports/staff${buildQuery(params)}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch staff report')
  return response.data
}

export async function getServiceReport(params: DateRangeParams = {}): Promise<ServiceReport> {
  const response = await apiClient.get<ServiceReport>(`/reports/services${buildQuery(params)}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch service report')
  return response.data
}

export async function getExpenseReport(params: DateRangeParams = {}): Promise<ExpenseAnalyticsReport> {
  const response = await apiClient.get<ExpenseAnalyticsReport>(`/reports/expenses${buildQuery(params)}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch expense report')
  return response.data
}

export async function getInventoryReport(): Promise<InventoryAnalyticsReport> {
  const response = await apiClient.get<InventoryAnalyticsReport>('/reports/inventory')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch inventory report')
  return response.data
}

export async function getProfitLossReport(params: DateRangeParams & { group_by?: string } = {}): Promise<ProfitLossReport> {
  const response = await apiClient.get<ProfitLossReport>(`/reports/profit-loss${buildQuery(params)}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch P&L report')
  return response.data
}
