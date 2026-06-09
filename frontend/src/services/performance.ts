import { apiClient } from './api-client'
import type { StaffPerformanceSummary, StaffPerformanceDaily, RevenueTrendPoint, PerformanceStats } from '@/types'

export interface PerformanceParams {
  staff_id?: string
  date?: string
  date_from?: string
  date_to?: string
  limit?: number
}

export async function getDailyPerformance(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  const query = new URLSearchParams()
  if (params.staff_id) query.set('staff_id', params.staff_id)
  if (params.date) query.set('date', params.date)
  const qs = query.toString()
  const path = qs ? `/performance/daily?${qs}` : '/performance/daily'
  const response = await apiClient.get<StaffPerformanceSummary[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch daily performance')
  return response.data || []
}

export async function getWeeklyPerformance(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  const query = new URLSearchParams()
  if (params.staff_id) query.set('staff_id', params.staff_id)
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  const qs = query.toString()
  const path = qs ? `/performance/weekly?${qs}` : '/performance/weekly'
  const response = await apiClient.get<StaffPerformanceSummary[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch weekly performance')
  return response.data || []
}

export async function getMonthlyPerformance(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  const query = new URLSearchParams()
  if (params.staff_id) query.set('staff_id', params.staff_id)
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  const qs = query.toString()
  const path = qs ? `/performance/monthly?${qs}` : '/performance/monthly'
  const response = await apiClient.get<StaffPerformanceSummary[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch monthly performance')
  return response.data || []
}

export async function getTopPerformers(params: PerformanceParams = {}): Promise<StaffPerformanceSummary[]> {
  const query = new URLSearchParams()
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  if (params.limit) query.set('limit', String(params.limit))
  const qs = query.toString()
  const path = qs ? `/performance/top-performers?${qs}` : '/performance/top-performers'
  const response = await apiClient.get<StaffPerformanceSummary[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch top performers')
  return response.data || []
}

export async function getRevenueTrend(params: PerformanceParams = {}): Promise<RevenueTrendPoint[]> {
  const query = new URLSearchParams()
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  const qs = query.toString()
  const path = qs ? `/performance/revenue-trend?${qs}` : '/performance/revenue-trend'
  const response = await apiClient.get<RevenueTrendPoint[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch revenue trend')
  return response.data || []
}

export async function getStaffPerformanceDetail(staffId: string, params: PerformanceParams = {}): Promise<StaffPerformanceDaily[]> {
  const query = new URLSearchParams()
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  const qs = query.toString()
  const path = qs ? `/performance/staff/${staffId}?${qs}` : `/performance/staff/${staffId}`
  const response = await apiClient.get<StaffPerformanceDaily[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch staff performance')
  return response.data || []
}

export async function getStaffRevenueTrend(staffId: string, params: PerformanceParams = {}): Promise<RevenueTrendPoint[]> {
  const query = new URLSearchParams()
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  const qs = query.toString()
  const path = qs ? `/performance/staff/${staffId}/trend?${qs}` : `/performance/staff/${staffId}/trend`
  const response = await apiClient.get<RevenueTrendPoint[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch staff revenue trend')
  return response.data || []
}

export async function getPerformanceStats(): Promise<PerformanceStats> {
  const response = await apiClient.get<PerformanceStats>('/performance/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch performance stats')
  return response.data
}
