import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useCommissionRules, useCommissionRuleById, useCreateCommissionRule, useUpdateCommissionRule, useDeleteCommissionRule, useStaffCommission, useMonthlyCommission, useCommissionStats } from './useCommissions'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useCommissions hooks', () => {
  it('useCommissionRules fetches rules', async () => {
    const { result } = renderHook(() => useCommissionRules(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.rules).toHaveLength(1)
  })

  it('useCommissionRuleById fetches rule', async () => {
    const { result } = renderHook(() => useCommissionRuleById('rule1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.rule_type).toBe('percentage')
  })

  it('useStaffCommission fetches commission', async () => {
    const { result } = renderHook(() => useStaffCommission('staff1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total).toBe(5000)
  })

  it('useMonthlyCommission fetches commission', async () => {
    const { result } = renderHook(() => useMonthlyCommission(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total).toBe(15000)
  })

  it('useCommissionStats fetches stats', async () => {
    const { result } = renderHook(() => useCommissionStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_commission).toBe(50000)
  })

  it('useCreateCommissionRule executes mutation', async () => {
    const { result } = renderHook(() => useCreateCommissionRule(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ rule_type: 'percentage', value: 10 } as any) })
  })

  it('useUpdateCommissionRule executes mutation', async () => {
    const { result } = renderHook(() => useUpdateCommissionRule(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'rule1', input: { value: 15 } as any }) })
  })

  it('useDeleteCommissionRule executes mutation', async () => {
    const { result } = renderHook(() => useDeleteCommissionRule(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('rule1') })
  })
})
