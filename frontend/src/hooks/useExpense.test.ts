import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useExpenseCategories, useCreateCategory, useUpdateCategory, useExpenseList, useExpenseById, useCreateExpense, useUpdateExpense, useDeleteExpense, useExpenseStats, useProfitLoss, useMonthlyTrend, useExpenseReport } from './useExpense'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useExpense hooks', () => {
  it('useExpenseCategories fetches categories', async () => {
    const { result } = renderHook(() => useExpenseCategories(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(2)
  })

  it('useExpenseList fetches list', async () => {
    const { result } = renderHook(() => useExpenseList(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.expenses).toHaveLength(1)
  })

  it('useExpenseById fetches expense', async () => {
    const { result } = renderHook(() => useExpenseById('exp1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.amount).toBe(25000)
  })

  it('useExpenseStats fetches stats', async () => {
    const { result } = renderHook(() => useExpenseStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.monthly_expenses).toBe(85000)
  })

  it('useProfitLoss fetches data', async () => {
    const { result } = renderHook(() => useProfitLoss('2024-01-01', '2024-12-31'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.revenue).toBe(300000)
  })

  it('useMonthlyTrend fetches data', async () => {
    const { result } = renderHook(() => useMonthlyTrend(6), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual([])
  })

  it('useExpenseReport fetches data', async () => {
    const { result } = renderHook(() => useExpenseReport('2024-01-01', '2024-12-31'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
  })

  it('useCreateCategory executes mutation', async () => {
    const { result } = renderHook(() => useCreateCategory(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ name: 'Rent', description: 'Monthly rent' }) })
  })

  it('useUpdateCategory executes mutation', async () => {
    const { result } = renderHook(() => useUpdateCategory(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'cat1', name: 'Rent', description: 'Office rent', is_active: true }) })
  })

  it('useCreateExpense executes mutation', async () => {
    const { result } = renderHook(() => useCreateExpense(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ category_id: 'cat1', amount: 5000, description: 'Test' } as any) })
  })

  it('useUpdateExpense executes mutation', async () => {
    const { result } = renderHook(() => useUpdateExpense(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'exp1', input: { amount: 6000 } as any }) })
  })

  it('useDeleteExpense executes mutation', async () => {
    const { result } = renderHook(() => useDeleteExpense(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('exp1') })
  })
})
