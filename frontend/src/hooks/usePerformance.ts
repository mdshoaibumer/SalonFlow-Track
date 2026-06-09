import { useQuery } from '@tanstack/react-query'
import {
  getDailyPerformance,
  getWeeklyPerformance,
  getMonthlyPerformance,
  getTopPerformers,
  getRevenueTrend,
  getStaffPerformanceDetail,
  getStaffRevenueTrend,
  getPerformanceStats,
  type PerformanceParams,
} from '@/services/performance'

export function useDailyPerformance(params: PerformanceParams = {}) {
  return useQuery({
    queryKey: ['performance', 'daily', params],
    queryFn: () => getDailyPerformance(params),
  })
}

export function useWeeklyPerformance(params: PerformanceParams = {}) {
  return useQuery({
    queryKey: ['performance', 'weekly', params],
    queryFn: () => getWeeklyPerformance(params),
  })
}

export function useMonthlyPerformance(params: PerformanceParams = {}) {
  return useQuery({
    queryKey: ['performance', 'monthly', params],
    queryFn: () => getMonthlyPerformance(params),
  })
}

export function useTopPerformers(params: PerformanceParams = {}) {
  return useQuery({
    queryKey: ['performance', 'top-performers', params],
    queryFn: () => getTopPerformers(params),
  })
}

export function useRevenueTrend(params: PerformanceParams = {}) {
  return useQuery({
    queryKey: ['performance', 'revenue-trend', params],
    queryFn: () => getRevenueTrend(params),
  })
}

export function useStaffPerformanceDetail(staffId: string, params: PerformanceParams = {}) {
  return useQuery({
    queryKey: ['performance', 'staff', staffId, params],
    queryFn: () => getStaffPerformanceDetail(staffId, params),
    enabled: !!staffId,
  })
}

export function useStaffRevenueTrend(staffId: string, params: PerformanceParams = {}) {
  return useQuery({
    queryKey: ['performance', 'staff-trend', staffId, params],
    queryFn: () => getStaffRevenueTrend(staffId, params),
    enabled: !!staffId,
  })
}

export function usePerformanceStats() {
  return useQuery({
    queryKey: ['performance', 'stats'],
    queryFn: getPerformanceStats,
  })
}
