import { describe, it, expect } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useDailyPerformance, useWeeklyPerformance, useMonthlyPerformance, useTopPerformers, useRevenueTrend, useStaffPerformanceDetail, useStaffRevenueTrend, usePerformanceStats } from './usePerformance'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('usePerformance hooks', () => {
  it('useDailyPerformance fetches data', async () => {
    const { result } = renderHook(() => useDailyPerformance(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual([])
  })

  it('useWeeklyPerformance fetches data', async () => {
    const { result } = renderHook(() => useWeeklyPerformance(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual([])
  })

  it('useMonthlyPerformance fetches data', async () => {
    const { result } = renderHook(() => useMonthlyPerformance(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual([])
  })

  it('useTopPerformers fetches data', async () => {
    const { result } = renderHook(() => useTopPerformers(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(2)
  })

  it('useRevenueTrend fetches data', async () => {
    const { result } = renderHook(() => useRevenueTrend(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(3)
  })

  it('useStaffPerformanceDetail fetches data', async () => {
    const { result } = renderHook(() => useStaffPerformanceDetail('staff1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual([])
  })

  it('useStaffRevenueTrend fetches data', async () => {
    const { result } = renderHook(() => useStaffRevenueTrend('staff1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual([])
  })

  it('usePerformanceStats fetches stats', async () => {
    const { result } = renderHook(() => usePerformanceStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_revenue_today).toBe(12500)
  })
})
