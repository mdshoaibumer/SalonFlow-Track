import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useGSTSettings, useSaveGSTSettings, useTaxRates, useCreateTaxRate, useUpdateTaxRate, useDeleteTaxRate, useGSTReport } from './useGST'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useGST hooks', () => {
  it('useGSTSettings fetches settings', async () => {
    const { result } = renderHook(() => useGSTSettings(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.gstin).toBe('27AABCU9603R1ZM')
  })

  it('useTaxRates fetches rates', async () => {
    const { result } = renderHook(() => useTaxRates(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useGSTReport fetches report', async () => {
    const { result } = renderHook(() => useGSTReport('2024-01-01', '2024-12-31'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_taxable).toBe(250000)
  })

  it('useSaveGSTSettings executes mutation', async () => {
    const { result } = renderHook(() => useSaveGSTSettings(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ gstin: '27AABCU9603R1ZM' } as any) })
  })

  it('useCreateTaxRate executes mutation', async () => {
    const { result } = renderHook(() => useCreateTaxRate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ name: 'GST 18%', rate: 18 } as any) })
  })

  it('useUpdateTaxRate executes mutation', async () => {
    const { result } = renderHook(() => useUpdateTaxRate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'rate1', rate: { name: 'GST 18%', rate: 18 } as any }) })
  })

  it('useDeleteTaxRate executes mutation', async () => {
    const { result } = renderHook(() => useDeleteTaxRate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('rate1') })
  })
})
