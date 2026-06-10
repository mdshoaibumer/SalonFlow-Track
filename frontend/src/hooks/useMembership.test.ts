import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useMembershipPlans, useMembershipPlan, useCreatePlan, useUpdatePlan, useDeletePlan, useSellPlan, useUseSession, useSubscriptions, useMembershipStats } from '@/hooks/useMembership'

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useMembership hooks', () => {
  it('useMembershipPlans fetches plans', async () => {
    const { result } = renderHook(() => useMembershipPlans(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useMembershipPlan fetches single plan', async () => {
    const { result } = renderHook(() => useMembershipPlan('plan1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.name).toBe('Gold')
  })

  it('useSubscriptions fetches subscriptions', async () => {
    const { result } = renderHook(() => useSubscriptions(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.data).toEqual([])
  })

  it('useMembershipStats fetches stats', async () => {
    const { result } = renderHook(() => useMembershipStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_plans).toBe(3)
  })

  it('useCreatePlan executes mutation', async () => {
    const { result } = renderHook(() => useCreatePlan(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync({ name: 'Silver', price: 3000 } as any)
    })
  })

  it('useUpdatePlan executes mutation', async () => {
    const { result } = renderHook(() => useUpdatePlan(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync({ id: 'plan1', name: 'Gold Pro' } as any)
    })
  })

  it('useDeletePlan executes mutation', async () => {
    const { result } = renderHook(() => useDeletePlan(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync('plan1')
    })
  })

  it('useSellPlan executes mutation', async () => {
    const { result } = renderHook(() => useSellPlan(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync({ plan_id: 'plan1', customer_id: 'cust1', amount_paid: 5000 })
    })
  })

  it('useUseSession executes mutation', async () => {
    const { result } = renderHook(() => useUseSession(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync('sub1')
    })
  })
})
