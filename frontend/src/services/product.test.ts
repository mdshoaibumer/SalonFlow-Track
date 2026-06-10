import { describe, it, expect, vi } from 'vitest'
import { listProducts, getProductById, createProduct, updateProduct, deleteProduct, adjustStock, listStockHistory, createPurchase, listPurchases, getInventoryStats, getLowStockProducts } from './product'

describe('Product Service', () => {
  it('lists products', async () => {
    const result = await listProducts()
    expect(result.products).toHaveLength(1)
    expect(result.products[0].name).toBe('Shampoo')
  })

  it('lists products with params', async () => {
    const result = await listProducts({ category: 'hair', status: 'active', search: 'shampoo', low_stock: true, page: 1, per_page: 10 })
    expect(result.products).toHaveLength(1)
  })

  it('gets product by id', async () => {
    const product = await getProductById('prod1')
    expect(product.name).toBe('Shampoo')
  })

  it('creates product', async () => {
    const product = await createProduct({ name: 'Shampoo', price: 500 } as any)
    expect(product.id).toBe('prod1')
  })

  it('updates product', async () => {
    const product = await updateProduct('prod1', { name: 'Shampoo Updated' } as any)
    expect(product.id).toBe('prod1')
  })

  it('deletes product', async () => {
    await expect(deleteProduct('prod1')).resolves.toBeUndefined()
  })

  it('adjusts stock', async () => {
    const txn = await adjustStock({ product_id: 'prod1', transaction_type: 'purchase', quantity: 5, notes: '' })
    expect(txn.id).toBe('txn1')
  })

  it('lists stock history', async () => {
    const result = await listStockHistory({})
    expect(result.transactions).toHaveLength(1)
  })

  it('lists stock history with params', async () => {
    const result = await listStockHistory({ product_id: 'prod1', transaction_type: 'purchase', date_from: '2024-01-01', date_to: '2024-12-31', page: 1, per_page: 10 })
    expect(result.transactions).toHaveLength(1)
  })

  it('creates purchase', async () => {
    const purchase = await createPurchase({ product_id: 'prod1', quantity: 50, total_cost: 5000 })
    expect(purchase.id).toBe('pur1')
  })

  it('lists purchases', async () => {
    const result = await listPurchases()
    expect(result.purchases).toHaveLength(1)
  })

  it('lists purchases with params', async () => {
    const result = await listPurchases({ date_from: '2024-01-01', date_to: '2024-12-31', page: 1, per_page: 10 })
    expect(result.purchases).toHaveLength(1)
  })

  it('gets inventory stats', async () => {
    const stats = await getInventoryStats()
    expect(stats.total_products).toBe(30)
  })

  it('gets low stock products', async () => {
    const items = await getLowStockProducts()
    expect(items).toHaveLength(1)
    expect(items[0].name).toBe('Conditioner')
  })

  it('listProducts returns empty when API returns undefined products', async () => {
    vi.spyOn(window.go.main.ProductService, 'ListProducts').mockResolvedValueOnce({ products: undefined, total: 0 } as any)
    const r = await listProducts()
    expect(r.products).toEqual([])
  })
})
