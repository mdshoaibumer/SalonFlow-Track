import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  createProduct,
  listProducts,
  getProductById,
  updateProduct,
  deleteProduct,
  adjustStock,
  listStockHistory,
  createPurchase,
  listPurchases,
  getInventoryStats,
  getLowStockProducts,
  type ListProductsParams,
} from '@/services/product'
import type { CreateProductInput, UpdateProductInput, StockAdjustInput, CreatePurchaseInput } from '@/types'

export function useProductList(params: ListProductsParams = {}) {
  return useQuery({
    queryKey: ['products', params],
    queryFn: () => listProducts(params),
  })
}

export function useProductById(id: string) {
  return useQuery({
    queryKey: ['products', id],
    queryFn: () => getProductById(id),
    enabled: !!id,
  })
}

export function useCreateProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateProductInput) => createProduct(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['inventory-stats'] })
    },
  })
}

export function useUpdateProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateProductInput }) => updateProduct(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['inventory-stats'] })
    },
  })
}

export function useDeleteProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['inventory-stats'] })
    },
  })
}

export function useAdjustStock() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: StockAdjustInput) => adjustStock(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['stock-history'] })
      queryClient.invalidateQueries({ queryKey: ['inventory-stats'] })
      queryClient.invalidateQueries({ queryKey: ['low-stock'] })
    },
  })
}

export function useStockHistory(params: { product_id?: string; transaction_type?: string; page?: number } = {}) {
  return useQuery({
    queryKey: ['stock-history', params],
    queryFn: () => listStockHistory(params),
  })
}

export function useCreatePurchase() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreatePurchaseInput) => createPurchase(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['purchases'] })
      queryClient.invalidateQueries({ queryKey: ['stock-history'] })
      queryClient.invalidateQueries({ queryKey: ['inventory-stats'] })
    },
  })
}

export function usePurchaseList(params: { date_from?: string; date_to?: string; page?: number } = {}) {
  return useQuery({
    queryKey: ['purchases', params],
    queryFn: () => listPurchases(params),
  })
}

export function useInventoryStats() {
  return useQuery({
    queryKey: ['inventory-stats'],
    queryFn: getInventoryStats,
  })
}

export function useLowStockProducts() {
  return useQuery({
    queryKey: ['low-stock'],
    queryFn: getLowStockProducts,
  })
}
