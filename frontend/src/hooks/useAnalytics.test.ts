import { describe, it, expect } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useDashboardStats, useKPIMetrics, useRevenueReport, useCustomerReport, useStaffAnalyticsReport, useServiceReport, useExpenseAnalyticsReport, useInventoryAnalyticsReport, useProfitLossReport } from './useAnalytics'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useAnalytics hooks', () => {
  it('useDashboardStats fetches data', async () => {
    const { result } = renderHook(() => useDashboardStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.today_revenue).toBe(12500)
  })

  it('useKPIMetrics fetches data', async () => {
    const { result } = renderHook(() => useKPIMetrics(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.revenue).toBe(300000)
  })

  it('useRevenueReport fetches data', async () => {
    const { result } = renderHook(() => useRevenueReport(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total).toBe(300000)
  })

  it('useCustomerReport fetches data', async () => {
    const { result } = renderHook(() => useCustomerReport(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_customers).toBe(50)
  })

  it('useStaffAnalyticsReport fetches data', async () => {
    const { result } = renderHook(() => useStaffAnalyticsReport(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_revenue).toBe(300000)
  })

  it('useServiceReport fetches data', async () => {
    const { result } = renderHook(() => useServiceReport(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_revenue).toBe(300000)
  })

  it('useExpenseAnalyticsReport fetches data', async () => {
    const { result } = renderHook(() => useExpenseAnalyticsReport(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total).toBe(85000)
  })

  it('useInventoryAnalyticsReport fetches data', async () => {
    const { result } = renderHook(() => useInventoryAnalyticsReport(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_products).toBe(30)
  })

  it('useProfitLossReport fetches data', async () => {
    const { result } = renderHook(() => useProfitLossReport(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.net_profit).toBe(215000)
  })
})
