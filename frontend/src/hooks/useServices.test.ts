import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useServiceList, useServiceById, useServiceStats, useCreateService, useUpdateService, useDeleteService } from './useServices'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useServices hooks', () => {
  it('useServiceList fetches list', async () => {
    const { result } = renderHook(() => useServiceList(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.services).toHaveLength(1)
  })

  it('useServiceById fetches service', async () => {
    const { result } = renderHook(() => useServiceById('svc1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.name).toBe('Haircut - Ladies')
  })

  it('useServiceStats fetches stats', async () => {
    const { result } = renderHook(() => useServiceStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total).toBe(10)
  })

  it('useCreateService executes mutation', async () => {
    const { result } = renderHook(() => useCreateService(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ name: 'New', price: 500 } as any) })
  })

  it('useUpdateService executes mutation', async () => {
    const { result } = renderHook(() => useUpdateService(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'svc1', input: { name: 'Updated' } as any }) })
  })

  it('useDeleteService executes mutation', async () => {
    const { result } = renderHook(() => useDeleteService(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('svc1') })
  })
})
