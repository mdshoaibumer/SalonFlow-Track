import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useSalaryList, useSalaryById, useSalaryCycles, useSalaryStats, useGenerateSalary, usePaySalary, useAdvanceList, useCreateAdvance, useApproveAdvance, useRejectAdvance } from './useSalary'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useSalary hooks', () => {
  it('useSalaryList fetches list', async () => {
    const { result } = renderHook(() => useSalaryList(12, 2024), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useSalaryById fetches record', async () => {
    const { result } = renderHook(() => useSalaryById('sal1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.staff_name).toBe('Priya Sharma')
  })

  it('useSalaryCycles fetches cycles', async () => {
    const { result } = renderHook(() => useSalaryCycles(2024), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useSalaryStats fetches stats', async () => {
    const { result } = renderHook(() => useSalaryStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_payroll).toBe(150000)
  })

  it('useAdvanceList fetches advances', async () => {
    const { result } = renderHook(() => useAdvanceList(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.advances).toHaveLength(1)
  })

  it('useGenerateSalary executes mutation', async () => {
    const { result } = renderHook(() => useGenerateSalary(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ month: 12, year: 2024 } as any) })
  })

  it('usePaySalary executes mutation', async () => {
    const { result } = renderHook(() => usePaySalary(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('sal1') })
  })

  it('useCreateAdvance executes mutation', async () => {
    const { result } = renderHook(() => useCreateAdvance(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ staff_id: 'staff1', amount: 5000 } as any) })
  })

  it('useApproveAdvance executes mutation', async () => {
    const { result } = renderHook(() => useApproveAdvance(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('adv1') })
  })

  it('useRejectAdvance executes mutation', async () => {
    const { result } = renderHook(() => useRejectAdvance(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('adv1') })
  })
})
