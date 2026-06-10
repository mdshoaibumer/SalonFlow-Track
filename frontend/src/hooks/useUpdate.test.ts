import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useUpdateStatus, useUpdateHistory, useCheckForUpdate, useDownloadUpdate, useInstallUpdate } from './useUpdate'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useUpdate hooks', () => {
  it('useUpdateStatus fetches status', async () => {
    const { result } = renderHook(() => useUpdateStatus(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.current_version).toBe('0.2.0')
  })

  it('useUpdateHistory fetches history', async () => {
    const { result } = renderHook(() => useUpdateHistory(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.records).toHaveLength(1)
  })

  it('useCheckForUpdate executes mutation', async () => {
    const { result } = renderHook(() => useCheckForUpdate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync() })
  })

  it('useDownloadUpdate executes mutation', async () => {
    const { result } = renderHook(() => useDownloadUpdate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync() })
  })

  it('useInstallUpdate executes mutation', async () => {
    const { result } = renderHook(() => useInstallUpdate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync() })
  })
})
