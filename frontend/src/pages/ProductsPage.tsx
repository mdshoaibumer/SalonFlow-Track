import { useState } from 'react'
import { useProductList, useCreateProduct, useDeleteProduct, useInventoryStats } from '@/hooks/useProduct'
import { Plus, Trash2, Package, AlertTriangle, IndianRupee, ShoppingCart } from 'lucide-react'
import type { CreateProductInput, ProductCategory } from '@/types'

const CATEGORIES: { value: ProductCategory; label: string }[] = [
  { value: 'hair_care', label: 'Hair Care' },
  { value: 'facial', label: 'Facial' },
  { value: 'spa', label: 'Spa' },
  { value: 'coloring', label: 'Coloring' },
  { value: 'treatment', label: 'Treatment' },
  { value: 'retail', label: 'Retail' },
  { value: 'equipment', label: 'Equipment' },
  { value: 'other', label: 'Other' },
]

export function ProductsPage() {
  const [showForm, setShowForm] = useState(false)
  const [filters, setFilters] = useState({ category: '', status: '', search: '', page: 1 })

  const { data: stats } = useInventoryStats()
  const { data, isLoading } = useProductList({ ...filters, per_page: 20 })
  const createProd = useCreateProduct()
  const deleteProd = useDeleteProduct()

  const handleCreate = (input: CreateProductInput) => {
    createProd.mutate(input, { onSuccess: () => setShowForm(false) })
  }

  const handleDelete = (id: string, name: string) => {
    if (confirm(`Delete product "${name}"?`)) {
      deleteProd.mutate(id)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Products & Inventory</h1>
          <p className="text-muted-foreground">Manage products, stock levels, and purchases</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="h-4 w-4" />
          Add Product
        </button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <StatCard title="Total Products" value={String(stats?.active_products ?? 0)} icon={<Package className="h-4 w-4" />} />
        <StatCard title="Low Stock Alerts" value={String(stats?.low_stock_count ?? 0)} icon={<AlertTriangle className="h-4 w-4 text-orange-500" />} />
        <StatCard title="Inventory Value" value={`₹${Math.round(stats?.total_value ?? 0).toLocaleString()}`} icon={<IndianRupee className="h-4 w-4 text-green-500" />} />
        <StatCard title="Purchases (Month)" value={`₹${Math.round(stats?.total_purchases_this_month ?? 0).toLocaleString()}`} icon={<ShoppingCart className="h-4 w-4" />} />
      </div>

      {showForm && <ProductForm onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={createProd.isPending} />}

      {/* Filters */}
      <div className="flex gap-2 flex-wrap">
        <input
          type="text"
          placeholder="Search products..."
          value={filters.search}
          onChange={(e) => setFilters({ ...filters, search: e.target.value, page: 1 })}
          className="rounded-md border px-3 py-2 text-sm w-64"
        />
        <select
          value={filters.category}
          onChange={(e) => setFilters({ ...filters, category: e.target.value, page: 1 })}
          className="rounded-md border px-3 py-2 text-sm"
        >
          <option value="">All Categories</option>
          {CATEGORIES.map((c) => (
            <option key={c.value} value={c.value}>{c.label}</option>
          ))}
        </select>
        <select
          value={filters.status}
          onChange={(e) => setFilters({ ...filters, status: e.target.value, page: 1 })}
          className="rounded-md border px-3 py-2 text-sm"
        >
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
          <option value="discontinued">Discontinued</option>
        </select>
      </div>

      {/* Product Table */}
      <div className="rounded-lg border bg-card">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Code</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Product</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Category</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Stock</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Min</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Purchase ₹</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Selling ₹</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Status</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Actions</th>
              </tr>
            </thead>
            <tbody>
              {isLoading && (
                <tr><td colSpan={9} className="px-6 py-8 text-center text-muted-foreground">Loading...</td></tr>
              )}
              {!isLoading && (!data?.products || data.products.length === 0) && (
                <tr><td colSpan={9} className="px-6 py-8 text-center text-muted-foreground">No products found</td></tr>
              )}
              {data?.products?.map((p) => (
                <tr key={p.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm font-mono">{p.product_code}</td>
                  <td className="px-6 py-4 text-sm font-medium">
                    {p.name}
                    {p.brand && <span className="ml-2 text-xs text-muted-foreground">({p.brand})</span>}
                  </td>
                  <td className="px-6 py-4 text-sm capitalize">{p.category.replace('_', ' ')}</td>
                  <td className={`px-6 py-4 text-sm text-right font-medium ${p.current_stock < p.minimum_stock ? 'text-red-600' : ''}`}>
                    {p.current_stock} {p.unit}
                  </td>
                  <td className="px-6 py-4 text-sm text-right text-muted-foreground">{p.minimum_stock}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{p.purchase_price.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{p.selling_price.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={p.status} lowStock={p.current_stock < p.minimum_stock} />
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button onClick={() => handleDelete(p.id, p.name)} className="rounded p-1 text-red-600 hover:bg-red-50" title="Delete">
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {data && data.meta.total_pages > 1 && (
          <div className="flex items-center justify-between border-t px-6 py-3">
            <span className="text-sm text-muted-foreground">Page {data.meta.page} of {data.meta.total_pages} ({data.meta.total} total)</span>
            <div className="flex gap-2">
              <button disabled={filters.page <= 1} onClick={() => setFilters({ ...filters, page: filters.page - 1 })} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Previous</button>
              <button disabled={filters.page >= data.meta.total_pages} onClick={() => setFilters({ ...filters, page: filters.page + 1 })} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Next</button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

function StatusBadge({ status, lowStock }: { status: string; lowStock: boolean }) {
  if (lowStock) return <span className="inline-flex rounded-full px-2 py-1 text-xs font-medium bg-red-100 text-red-800">Low Stock</span>
  const colors: Record<string, string> = {
    active: 'bg-green-100 text-green-800',
    inactive: 'bg-gray-100 text-gray-800',
    discontinued: 'bg-red-100 text-red-800',
  }
  return <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${colors[status] || 'bg-gray-100 text-gray-800'}`}>{status}</span>
}

function StatCard({ title, value, icon }: { title: string; value: string; icon: React.ReactNode }) {
  return (
    <div className="rounded-lg border bg-card p-6">
      <div className="flex items-center gap-2 text-muted-foreground mb-2">
        {icon}
        <span className="text-sm font-medium">{title}</span>
      </div>
      <p className="text-2xl font-bold">{value}</p>
    </div>
  )
}

function ProductForm({ onSubmit, onCancel, isLoading }: { onSubmit: (input: CreateProductInput) => void; onCancel: () => void; isLoading: boolean }) {
  const [form, setForm] = useState<CreateProductInput>({
    name: '',
    category: 'hair_care',
    brand: '',
    unit: 'pcs',
    sku: '',
    purchase_price: 0,
    selling_price: 0,
    minimum_stock: 5,
    maximum_stock: 50,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="rounded-lg border bg-card p-6 space-y-4">
      <h3 className="text-lg font-semibold">New Product</h3>
      <div className="grid gap-4 md:grid-cols-3">
        <div>
          <label className="text-sm font-medium">Product Name</label>
          <input type="text" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" required />
        </div>
        <div>
          <label className="text-sm font-medium">Category</label>
          <select value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value as ProductCategory })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm">
            {CATEGORIES.map((c) => <option key={c.value} value={c.value}>{c.label}</option>)}
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Brand</label>
          <input type="text" value={form.brand} onChange={(e) => setForm({ ...form, brand: e.target.value })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" />
        </div>
        <div>
          <label className="text-sm font-medium">Unit</label>
          <select value={form.unit} onChange={(e) => setForm({ ...form, unit: e.target.value })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm">
            <option value="pcs">Pieces</option>
            <option value="ml">ml</option>
            <option value="gm">grams</option>
            <option value="ltr">Litres</option>
            <option value="kg">kg</option>
            <option value="box">Box</option>
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Purchase Price (₹)</label>
          <input type="number" value={form.purchase_price || ''} onChange={(e) => setForm({ ...form, purchase_price: Number(e.target.value) })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" min={0} required />
        </div>
        <div>
          <label className="text-sm font-medium">Selling Price (₹)</label>
          <input type="number" value={form.selling_price || ''} onChange={(e) => setForm({ ...form, selling_price: Number(e.target.value) })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" min={0} required />
        </div>
        <div>
          <label className="text-sm font-medium">Min Stock</label>
          <input type="number" value={form.minimum_stock || ''} onChange={(e) => setForm({ ...form, minimum_stock: Number(e.target.value) })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" min={0} />
        </div>
        <div>
          <label className="text-sm font-medium">Max Stock</label>
          <input type="number" value={form.maximum_stock || ''} onChange={(e) => setForm({ ...form, maximum_stock: Number(e.target.value) })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" min={0} />
        </div>
        <div>
          <label className="text-sm font-medium">SKU (Optional)</label>
          <input type="text" value={form.sku || ''} onChange={(e) => setForm({ ...form, sku: e.target.value })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" />
        </div>
      </div>
      <div className="flex gap-2">
        <button type="submit" disabled={isLoading} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
          {isLoading ? 'Creating...' : 'Add Product'}
        </button>
        <button type="button" onClick={onCancel} className="rounded-lg border px-4 py-2 text-sm font-medium">Cancel</button>
      </div>
    </form>
  )
}
