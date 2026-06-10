import { useState } from 'react'
import { useInventoryStats, useLowStockProducts, useAdjustStock, useProductList, useStockHistory } from '@/hooks/useProduct'
import { Package, AlertTriangle, IndianRupee, TrendingDown } from 'lucide-react'
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend } from 'recharts'
import type { StockAdjustInput, StockTransactionType } from '@/types'

const PIE_COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4', '#ec4899', '#6b7280']

export function InventoryPage() {
  const { data: stats } = useInventoryStats()
  const { data: lowStock } = useLowStockProducts()
  const { data: products } = useProductList({ per_page: 200, status: 'active' })
  const { data: stockHistory } = useStockHistory({ page: 1 })
  const adjustStock = useAdjustStock()
  const [showAdjust, setShowAdjust] = useState(false)

  // Build category distribution for pie chart
  const catMap = new Map<string, number>()
  products?.products?.forEach((p) => {
    catMap.set(p.category, (catMap.get(p.category) || 0) + p.current_stock * p.purchase_price)
  })
  const pieData = Array.from(catMap.entries()).map(([name, value]) => ({ name: name.replace('_', ' '), value: Math.round(value) }))

  const handleAdjust = (input: StockAdjustInput) => {
    adjustStock.mutate(input, { onSuccess: () => setShowAdjust(false) })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Inventory Dashboard</h1>
          <p className="text-muted-foreground">Stock levels, adjustments, and inventory valuation</p>
        </div>
        <button onClick={() => setShowAdjust(!showAdjust)} className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
          Stock Adjustment
        </button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <StatCard title="Active Products" value={String(stats?.active_products ?? 0)} icon={<Package className="h-4 w-4" />} />
        <StatCard title="Low Stock Alerts" value={String(stats?.low_stock_count ?? 0)} icon={<AlertTriangle className="h-4 w-4 text-orange-500" />} />
        <StatCard title="Inventory Value" value={`₹${Math.round(stats?.total_value ?? 0).toLocaleString()}`} icon={<IndianRupee className="h-4 w-4 text-green-500" />} />
        <StatCard title="Monthly Purchases" value={`₹${Math.round(stats?.total_purchases_this_month ?? 0).toLocaleString()}`} icon={<TrendingDown className="h-4 w-4 text-blue-500" />} />
      </div>

      {showAdjust && <AdjustForm products={products?.products || []} onSubmit={handleAdjust} onCancel={() => setShowAdjust(false)} isLoading={adjustStock.isPending} />}

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Pie Chart */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Inventory Value by Category</h3>
          {pieData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie data={pieData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={100} label={(entry) => `${entry.name} ₹${(entry.value ?? 0).toLocaleString()}`}>
                  {pieData.map((_, idx) => (
                    <Cell key={idx} fill={PIE_COLORS[idx % PIE_COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Value']} />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-center text-muted-foreground py-8">No data</p>
          )}
        </div>

        {/* Low Stock */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Low Stock Alerts</h3>
          {!lowStock || lowStock.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">All stock levels are healthy</p>
          ) : (
            <div className="space-y-3 max-h-[300px] overflow-y-auto">
              {lowStock.map((item) => (
                <div key={item.product_id} className="flex items-center justify-between rounded-md border px-4 py-3">
                  <div>
                    <p className="text-sm font-medium">{item.product_name}</p>
                    <p className="text-xs text-muted-foreground">{item.product_code}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium text-red-600">{item.current_stock} / {item.minimum_stock}</p>
                    <p className="text-xs text-muted-foreground">Need {item.minimum_stock - item.current_stock} more</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Recent Stock Transactions */}
      <div className="rounded-lg border bg-card">
        <div className="p-6 pb-0">
          <h3 className="text-lg font-semibold">Recent Stock Transactions</h3>
        </div>
        <div className="overflow-x-auto p-6 pt-4">
          <table className="w-full">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-muted-foreground">Date</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-muted-foreground">Product</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-muted-foreground">Type</th>
                <th className="px-4 py-2 text-right text-xs font-medium uppercase text-muted-foreground">Qty</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-muted-foreground">Notes</th>
              </tr>
            </thead>
            <tbody>
              {(!stockHistory?.transactions || stockHistory.transactions.length === 0) && (
                <tr><td colSpan={5} className="px-4 py-6 text-center text-muted-foreground">No transactions yet</td></tr>
              )}
              {stockHistory?.transactions?.map((t: any) => (
                <tr key={t.id} className="border-t">
                  <td className="px-4 py-2 text-sm">{new Date(t.transaction_date).toLocaleDateString('en-IN')}</td>
                  <td className="px-4 py-2 text-sm font-medium">{t.product_id.slice(0, 8)}...</td>
                  <td className="px-4 py-2 text-sm capitalize">
                    <TypeBadge type={t.transaction_type} />
                  </td>
                  <td className={`px-4 py-2 text-sm text-right font-medium ${['consumption', 'sale', 'damage'].includes(t.transaction_type) ? 'text-red-600' : 'text-green-600'}`}>
                    {['consumption', 'sale', 'damage'].includes(t.transaction_type) ? '-' : '+'}{t.quantity}
                  </td>
                  <td className="px-4 py-2 text-sm text-muted-foreground">{t.notes || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

function TypeBadge({ type }: { type: string }) {
  const colors: Record<string, string> = {
    purchase: 'bg-green-100 text-green-800',
    consumption: 'bg-orange-100 text-orange-800',
    sale: 'bg-blue-100 text-blue-800',
    adjustment: 'bg-purple-100 text-purple-800',
    return: 'bg-yellow-100 text-yellow-800',
    damage: 'bg-red-100 text-red-800',
  }
  return <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${colors[type] || 'bg-gray-100 text-gray-800'}`}>{type}</span>
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

function AdjustForm({ products, onSubmit, onCancel, isLoading }: { products: { id: string; name: string; product_code: string }[]; onSubmit: (i: StockAdjustInput) => void; onCancel: () => void; isLoading: boolean }) {
  const [form, setForm] = useState<StockAdjustInput>({
    product_id: '',
    transaction_type: 'consumption',
    quantity: 1,
    notes: '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="rounded-lg border bg-card p-6 space-y-4">
      <h3 className="text-lg font-semibold">Stock Adjustment</h3>
      <div className="grid gap-4 md:grid-cols-5">
        <div>
          <label className="text-sm font-medium">Product</label>
          <select value={form.product_id} onChange={(e) => setForm({ ...form, product_id: e.target.value })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" required>
            <option value="">Select Product</option>
            {products.map((p) => <option key={p.id} value={p.id}>{p.name} ({p.product_code})</option>)}
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Type</label>
          <select value={form.transaction_type} onChange={(e) => setForm({ ...form, transaction_type: e.target.value as StockTransactionType })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm">
            <option value="consumption">Consumption</option>
            <option value="sale">Sale</option>
            <option value="damage">Damage</option>
            <option value="adjustment">Adjustment</option>
            <option value="return">Return</option>
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Quantity</label>
          <input type="number" value={form.quantity || ''} onChange={(e) => setForm({ ...form, quantity: Number(e.target.value) })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" min={1} required />
        </div>

        <div>
          <label className="text-sm font-medium">Notes</label>
          <input type="text" value={form.notes || ''} onChange={(e) => setForm({ ...form, notes: e.target.value })} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" />
        </div>
      </div>
      <div className="flex gap-2">
        <button type="submit" disabled={isLoading} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
          {isLoading ? 'Adjusting...' : 'Submit'}
        </button>
        <button type="button" onClick={onCancel} className="rounded-lg border px-4 py-2 text-sm font-medium">Cancel</button>
      </div>
    </form>
  )
}
