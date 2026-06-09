import { apiClient } from './api-client'
import type {
  Product,
  CreateProductInput,
  UpdateProductInput,
  StockTransaction,
  PurchaseEntry,
  CreatePurchaseInput,
  StockAdjustInput,
  InventoryStats,
  LowStockItem,
} from '@/types'

export interface ListProductsParams {
  page?: number
  per_page?: number
  category?: string
  status?: string
  search?: string
  low_stock?: boolean
}

// --- Products ---

export async function createProduct(input: CreateProductInput): Promise<Product> {
  const response = await apiClient.post<Product>('/products', input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create product')
  return response.data
}

export async function listProducts(params: ListProductsParams = {}): Promise<{ products: Product[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.category) query.set('category', params.category)
  if (params.status) query.set('status', params.status)
  if (params.search) query.set('search', params.search)
  if (params.low_stock) query.set('low_stock', 'true')

  const response = await apiClient.get<Product[]>(`/products?${query.toString()}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch products')
  return {
    products: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}

export async function getProductById(id: string): Promise<Product> {
  const response = await apiClient.get<Product>(`/products/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch product')
  return response.data
}

export async function updateProduct(id: string, input: UpdateProductInput): Promise<Product> {
  const response = await apiClient.put<Product>(`/products/${id}`, input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update product')
  return response.data
}

export async function deleteProduct(id: string): Promise<void> {
  const response = await apiClient.delete(`/products/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete product')
}

// --- Stock ---

export async function adjustStock(input: StockAdjustInput): Promise<StockTransaction> {
  const response = await apiClient.post<StockTransaction>('/products/stock/adjust', input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to adjust stock')
  return response.data
}

export async function listStockHistory(params: { product_id?: string; transaction_type?: string; date_from?: string; date_to?: string; page?: number; per_page?: number } = {}): Promise<{ transactions: StockTransaction[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const query = new URLSearchParams()
  if (params.product_id) query.set('product_id', params.product_id)
  if (params.transaction_type) query.set('transaction_type', params.transaction_type)
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))

  const response = await apiClient.get<StockTransaction[]>(`/products/stock/history?${query.toString()}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch stock history')
  return {
    transactions: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}

// --- Purchases ---

export async function createPurchase(input: CreatePurchaseInput): Promise<PurchaseEntry> {
  const response = await apiClient.post<PurchaseEntry>('/products/purchases', input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create purchase')
  return response.data
}

export async function listPurchases(params: { date_from?: string; date_to?: string; page?: number; per_page?: number } = {}): Promise<{ purchases: PurchaseEntry[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const query = new URLSearchParams()
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))

  const response = await apiClient.get<PurchaseEntry[]>(`/products/purchases?${query.toString()}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch purchases')
  return {
    purchases: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}

export async function getPurchaseById(id: string): Promise<PurchaseEntry> {
  const response = await apiClient.get<PurchaseEntry>(`/products/purchases/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch purchase')
  return response.data
}

// --- Reporting ---

export async function getInventoryStats(): Promise<InventoryStats> {
  const response = await apiClient.get<InventoryStats>('/products/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch inventory stats')
  return response.data
}

export async function getLowStockProducts(): Promise<LowStockItem[]> {
  const response = await apiClient.get<LowStockItem[]>('/products/low-stock')
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch low stock')
  return response.data || []
}
