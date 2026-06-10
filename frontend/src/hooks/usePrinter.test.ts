import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { usePrinterSettings, useSavePrinterSettings, usePrintInvoice, usePrintReceipt, usePrintTest, usePrintJobs } from './usePrinter'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('usePrinter hooks', () => {
  it('usePrinterSettings fetches settings', async () => {
    const { result } = renderHook(() => usePrinterSettings(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.printer_name).toBe('POS-80')
  })

  it('usePrintJobs fetches jobs', async () => {
    const { result } = renderHook(() => usePrintJobs(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.jobs).toHaveLength(1)
  })

  it('useSavePrinterSettings executes mutation', async () => {
    const { result } = renderHook(() => useSavePrinterSettings(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ printer_name: 'POS-80' } as any) })
  })

  it('usePrintInvoice executes mutation', async () => {
    const { result } = renderHook(() => usePrintInvoice(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ invoice_id: 'inv1' } as any) })
  })

  it('usePrintReceipt executes mutation', async () => {
    const { result } = renderHook(() => usePrintReceipt(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ invoice_id: 'inv1' } as any) })
  })

  it('usePrintTest executes mutation', async () => {
    const { result } = renderHook(() => usePrintTest(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync() })
  })
})
