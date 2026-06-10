import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useInvoiceList, useInvoiceById, useInvoiceStats, useCreateInvoice, useRecordPayment } from './useInvoices'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useInvoices hooks', () => {
  it('useInvoiceList fetches list', async () => {
    const { result } = renderHook(() => useInvoiceList(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.invoices).toHaveLength(1)
  })

  it('useInvoiceById fetches invoice', async () => {
    const { result } = renderHook(() => useInvoiceById('inv1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.invoice_number).toBe('INV-001')
  })

  it('useInvoiceStats fetches stats', async () => {
    const { result } = renderHook(() => useInvoiceStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.today_revenue).toBe(12500)
  })

  it('useCreateInvoice executes mutation', async () => {
    const { result } = renderHook(() => useCreateInvoice(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ customer_id: 'c1', items: [] } as any) })
  })

  it('useRecordPayment executes mutation', async () => {
    const { result } = renderHook(() => useRecordPayment(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ invoiceId: 'inv1', input: { amount: 2242, method: 'upi' } as any }) })
  })
})
