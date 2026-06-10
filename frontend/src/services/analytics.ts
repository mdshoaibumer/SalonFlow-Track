import type { DashboardStats, KPIMetrics, RevenueReport, CustomerReport, StaffAnalyticsReport, ServiceReport, ExpenseAnalyticsReport, InventoryAnalyticsReport, ProfitLossReport } from '@/types'

export interface DateRangeParams {
  date_from?: string
  date_to?: string
  group_by?: string
}

export async function getDashboardStats(): Promise<DashboardStats> {
  return window.go.main.AnalyticsService.GetDashboard()
}

export async function getKPIMetrics(params: DateRangeParams = {}): Promise<KPIMetrics> {
  return window.go.main.AnalyticsService.GetKPIs(params.date_from || '', params.date_to || '')
}

export async function getRevenueReport(params: DateRangeParams = {}): Promise<RevenueReport> {
  return window.go.main.AnalyticsService.GetRevenueReport(params.date_from || '', params.date_to || '', params.group_by || 'day')
}

export async function getCustomerReport(params: DateRangeParams = {}): Promise<CustomerReport> {
  return window.go.main.AnalyticsService.GetCustomerReport(params.date_from || '', params.date_to || '')
}

export async function getStaffReport(params: DateRangeParams = {}): Promise<StaffAnalyticsReport> {
  return window.go.main.AnalyticsService.GetStaffReport(params.date_from || '', params.date_to || '')
}

export async function getServiceReport(params: DateRangeParams = {}): Promise<ServiceReport> {
  return window.go.main.AnalyticsService.GetServiceReport(params.date_from || '', params.date_to || '')
}

export async function getExpenseReport(params: DateRangeParams = {}): Promise<ExpenseAnalyticsReport> {
  return window.go.main.AnalyticsService.GetExpenseAnalytics(params.date_from || '', params.date_to || '')
}

export async function getInventoryReport(): Promise<InventoryAnalyticsReport> {
  return window.go.main.AnalyticsService.GetInventoryReport()
}

export async function getProfitLossReport(params: DateRangeParams = {}): Promise<ProfitLossReport> {
  return window.go.main.AnalyticsService.GetProfitLossReport(params.date_from || '', params.date_to || '', params.group_by || 'month')
}
