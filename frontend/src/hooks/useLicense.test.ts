import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useLicenseStatus, useLicenseEvents, useValidateLicense, useActivateLicense, useRenewLicense } from './useLicense'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useLicense hooks', () => {
  it('useLicenseStatus fetches status', async () => {
    const { result } = renderHook(() => useLicenseStatus(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.status).toBe('active')
  })

  it('useLicenseEvents fetches events', async () => {
    const { result } = renderHook(() => useLicenseEvents(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.events).toHaveLength(1)
  })

  it('useValidateLicense executes mutation', async () => {
    const { result } = renderHook(() => useValidateLicense(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync() })
  })

  it('useActivateLicense executes mutation', async () => {
    const { result } = renderHook(() => useActivateLicense(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ licenseKey: 'ABC-123', customerName: 'Test', salonName: 'Salon' }) })
  })

  it('useRenewLicense executes mutation', async () => {
    const { result } = renderHook(() => useRenewLicense(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('ABC-123') })
  })
})
