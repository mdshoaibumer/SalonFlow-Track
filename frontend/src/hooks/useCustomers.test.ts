import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useCustomerList, useCustomerById, useCustomerStats, useCreateCustomer, useUpdateCustomer, useDeleteCustomer } from './useCustomers'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useCustomers hooks', () => {
  it('useCustomerList fetches list', async () => {
    const { result } = renderHook(() => useCustomerList(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.customers).toHaveLength(1)
  })

  it('useCustomerById fetches customer', async () => {
    const { result } = renderHook(() => useCustomerById('cust1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.full_name).toBe('Anjali Desai')
  })

  it('useCustomerStats fetches stats', async () => {
    const { result } = renderHook(() => useCustomerStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total).toBe(50)
  })

  it('useCreateCustomer executes mutation', async () => {
    const { result } = renderHook(() => useCreateCustomer(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync({ full_name: 'New', phone: '9876543212', gender: 'female' } as any)
    })
  })

  it('useUpdateCustomer executes mutation', async () => {
    const { result } = renderHook(() => useUpdateCustomer(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync({ id: 'cust1', input: { full_name: 'Updated' } as any })
    })
  })

  it('useDeleteCustomer executes mutation', async () => {
    const { result } = renderHook(() => useDeleteCustomer(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync('cust1')
    })
  })
})
