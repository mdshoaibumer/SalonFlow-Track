import type { StaffPerformanceSummary, StaffPerformanceDaily, RevenueTrendPoint, PerformanceStats } from '@/types'

export interface PerformanceParams {
  staff_id?: string
  date_from?: string
  date_to?: string
  date?: string
  limit?: number
}

export async function getDailyPerformance(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  return window.go.main.PerformanceService.GetDailyPerformance({
    staff_id: params.staff_id || '',
    date: params.date || '',
  })
}

export async function getWeeklyPerformance(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  return window.go.main.PerformanceService.GetWeeklyPerformance({
    staff_id: params.staff_id || '',
    date_from: params.date_from || '',
    date_to: params.date_to || '',
  })
}

export async function getMonthlyPerformance(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  return window.go.main.PerformanceService.GetMonthlyPerformance({
    staff_id: params.staff_id || '',
    date_from: params.date_from || '',
    date_to: params.date_to || '',
  })
}

export async function getTopPerformers(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  return window.go.main.PerformanceService.GetTopPerformers({
    date_from: params.date_from || '',
    date_to: params.date_to || '',
    limit: params.limit || 5,
  })
}

export async function getRevenueTrend(params: PerformanceParams = {}): Promise<RevenueTrendPoint[]> {
  return window.go.main.PerformanceService.GetRevenueTrend(params.date_from || '', params.date_to || '')
}

export async function getStaffPerformanceDetail(staffId: string, params: PerformanceParams = {}): Promise<StaffPerformanceDaily[]> {
  return window.go.main.PerformanceService.GetStaffPerformance(staffId, params.date_from || '', params.date_to || '')
}

export async function getStaffRevenueTrend(staffId: string, params: PerformanceParams = {}): Promise<RevenueTrendPoint[]> {
  return window.go.main.PerformanceService.GetStaffRevenueTrend(staffId, params.date_from || '', params.date_to || '')
}

export async function getPerformanceStats(): Promise<PerformanceStats> {
  return window.go.main.PerformanceService.GetPerformanceStats()
}
