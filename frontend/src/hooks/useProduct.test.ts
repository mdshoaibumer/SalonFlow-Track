import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useProductList, useProductById, useCreateProduct, useUpdateProduct, useDeleteProduct, useAdjustStock, useStockHistory, useCreatePurchase, usePurchaseList, useInventoryStats, useLowStockProducts } from './useProduct'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useProduct hooks', () => {
  it('useProductList fetches list', async () => {
    const { result } = renderHook(() => useProductList(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.products).toHaveLength(1)
  })

  it('useProductById fetches product', async () => {
    const { result } = renderHook(() => useProductById('prod1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.name).toBe('Shampoo')
  })

  it('useStockHistory fetches history', async () => {
    const { result } = renderHook(() => useStockHistory(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.transactions).toHaveLength(1)
  })

  it('usePurchaseList fetches purchases', async () => {
    const { result } = renderHook(() => usePurchaseList(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.purchases).toHaveLength(1)
  })

  it('useInventoryStats fetches stats', async () => {
    const { result } = renderHook(() => useInventoryStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_products).toBe(30)
  })

  it('useLowStockProducts fetches products', async () => {
    const { result } = renderHook(() => useLowStockProducts(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useCreateProduct executes mutation', async () => {
    const { result } = renderHook(() => useCreateProduct(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ name: 'Conditioner', sku: 'COND-01' } as any) })
  })

  it('useUpdateProduct executes mutation', async () => {
    const { result } = renderHook(() => useUpdateProduct(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'prod1', input: { name: 'Updated' } as any }) })
  })

  it('useDeleteProduct executes mutation', async () => {
    const { result } = renderHook(() => useDeleteProduct(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('prod1') })
  })

  it('useAdjustStock executes mutation', async () => {
    const { result } = renderHook(() => useAdjustStock(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ product_id: 'prod1', quantity: 5, type: 'add' } as any) })
  })

  it('useCreatePurchase executes mutation', async () => {
    const { result } = renderHook(() => useCreatePurchase(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ supplier: 'Supplier A', items: [] } as any) })
  })
})
