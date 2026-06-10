import type { Product, InventoryStats, CreateProductInput, UpdateProductInput, StockTransaction, PurchaseEntry, LowStockItem } from '@/types'

export interface ListProductsParams {
  page?: number
  per_page?: number
  category?: string
  status?: string
  search?: string
  low_stock?: boolean
}

export interface ListProductsResponse {
  products: Product[]
  meta: { page: number; per_page: number; total: number; total_pages: number }
}

export async function listProducts(params: ListProductsParams = {}): Promise<ListProductsResponse> {
  const result = await window.go.main.ProductService.ListProducts({
    Category: params.category || '',
    Status: params.status || '',
    Search: params.search || '',
    LowStock: params.low_stock || false,
    Page: params.page || 1,
    PerPage: params.per_page || 20,
  })
  return {
    products: result.products || [],
    meta: { page: result.page, per_page: result.per_page, total: result.total, total_pages: result.total_pages },
  }
}

export async function getProductById(id: string): Promise<Product> {
  return window.go.main.ProductService.GetProduct(id)
}

export async function createProduct(input: CreateProductInput): Promise<Product> {
  return window.go.main.ProductService.CreateProduct(input)
}

export async function updateProduct(id: string, input: UpdateProductInput): Promise<Product> {
  return window.go.main.ProductService.UpdateProduct(id, input)
}

export async function deleteProduct(id: string): Promise<void> {
  await window.go.main.ProductService.DeleteProduct(id)
}

export async function adjustStock(input: { product_id: string; transaction_type: string; quantity: number; notes: string }): Promise<StockTransaction> {
  return window.go.main.ProductService.AdjustStock(input)
}

export async function listStockHistory(input: { product_id?: string; transaction_type?: string; date_from?: string; date_to?: string; page?: number; per_page?: number }) {
  return window.go.main.ProductService.ListStockHistory({
    ProductID: input.product_id || '',
    TransactionType: input.transaction_type || '',
    DateFrom: input.date_from || '',
    DateTo: input.date_to || '',
    Page: input.page || 1,
    PerPage: input.per_page || 20,
  })
}

export async function createPurchase(input: any): Promise<PurchaseEntry> {
  return window.go.main.ProductService.CreatePurchase(input)
}

export async function listPurchases(input: { date_from?: string; date_to?: string; page?: number; per_page?: number } = {}) {
  return window.go.main.ProductService.ListPurchases({
    DateFrom: input.date_from || '',
    DateTo: input.date_to || '',
    Page: input.page || 1,
    PerPage: input.per_page || 20,
  })
}

export async function getInventoryStats(): Promise<InventoryStats> {
  return window.go.main.ProductService.GetInventoryStats()
}

export async function getLowStockProducts(): Promise<LowStockItem[]> {
  return window.go.main.ProductService.GetLowStockProducts()
}
