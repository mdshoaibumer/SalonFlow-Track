import { useState } from 'react'
import { useCreatePurchase, usePurchaseList, useProductList } from '@/hooks/useProduct'
import { Plus, Trash2 } from 'lucide-react'
import type { CreatePurchaseInput } from '@/types'

export function PurchasesPage() {
  const [showForm, setShowForm] = useState(false)
  const [page, setPage] = useState(1)

  const { data, isLoading } = usePurchaseList({ page })
  const createPurch = useCreatePurchase()

  const handleCreate = (input: CreatePurchaseInput) => {
    createPurch.mutate(input, { onSuccess: () => setShowForm(false) })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Purchases</h1>
          <p className="text-muted-foreground">Record product purchases from vendors</p>
        </div>
        <button onClick={() => setShowForm(!showForm)} className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
          <Plus className="h-4 w-4" />
          New Purchase
        </button>
      </div>

      {showForm && <PurchaseForm onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={createPurch.isPending} />}

      <div className="rounded-lg border bg-card">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Purchase #</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Vendor</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Invoice</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Date</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Amount</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Notes</th>
              </tr>
            </thead>
            <tbody>
              {isLoading && <tr><td colSpan={6} className="px-6 py-8 text-center text-muted-foreground">Loading...</td></tr>}
              {!isLoading && (!data?.purchases || data.purchases.length === 0) && (
                <tr><td colSpan={6} className="px-6 py-8 text-center text-muted-foreground">No purchases found</td></tr>
              )}
              {data?.purchases?.map((p: any) => (
                <tr key={p.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm font-mono">{p.purchase_number}</td>
                  <td className="px-6 py-4 text-sm font-medium">{p.vendor_name}</td>
                  <td className="px-6 py-4 text-sm">{p.invoice_number || '-'}</td>
                  <td className="px-6 py-4 text-sm">{new Date(p.purchase_date).toLocaleDateString('en-IN')}</td>
                  <td className="px-6 py-4 text-sm text-right font-medium">₹{p.total_amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-muted-foreground">{p.notes || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {data && data.meta.total_pages > 1 && (
          <div className="flex items-center justify-between border-t px-6 py-3">
            <span className="text-sm text-muted-foreground">Page {data.meta.page} of {data.meta.total_pages}</span>
            <div className="flex gap-2">
              <button disabled={page <= 1} onClick={() => setPage(page - 1)} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Previous</button>
              <button disabled={page >= data.meta.total_pages} onClick={() => setPage(page + 1)} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Next</button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

function PurchaseForm({ onSubmit, onCancel, isLoading }: { onSubmit: (input: CreatePurchaseInput) => void; onCancel: () => void; isLoading: boolean }) {
  const { data: productData } = useProductList({ per_page: 200, status: 'active' })
  const [vendorName, setVendorName] = useState('')
  const [invoiceNumber, setInvoiceNumber] = useState('')
  const [purchaseDate, setPurchaseDate] = useState(new Date().toISOString().slice(0, 10))
  const [notes, setNotes] = useState('')
  const [items, setItems] = useState<{ product_id: string; quantity: number; unit_price: number }[]>([
    { product_id: '', quantity: 1, unit_price: 0 },
  ])

  const addItem = () => setItems([...items, { product_id: '', quantity: 1, unit_price: 0 }])
  const removeItem = (idx: number) => setItems(items.filter((_, i) => i !== idx))
  const updateItem = (idx: number, field: string, value: string | number) => {
    setItems(items.map((it, i) => (i === idx ? { ...it, [field]: value } : it)))
  }

  const total = items.reduce((sum, it) => sum + it.quantity * it.unit_price, 0)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({
      vendor_name: vendorName,
      invoice_number: invoiceNumber,
      purchase_date: purchaseDate,
      notes,
      items: items.filter((it) => it.product_id),
    })
  }

  return (
    <form onSubmit={handleSubmit} className="rounded-lg border bg-card p-6 space-y-4">
      <h3 className="text-lg font-semibold">New Purchase Entry</h3>
      <div className="grid gap-4 md:grid-cols-4">
        <div>
          <label className="text-sm font-medium">Vendor Name</label>
          <input type="text" value={vendorName} onChange={(e) => setVendorName(e.target.value)} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" required />
        </div>
        <div>
          <label className="text-sm font-medium">Invoice Number</label>
          <input type="text" value={invoiceNumber} onChange={(e) => setInvoiceNumber(e.target.value)} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" />
        </div>
        <div>
          <label className="text-sm font-medium">Date</label>
          <input type="date" value={purchaseDate} onChange={(e) => setPurchaseDate(e.target.value)} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" required />
        </div>
        <div>
          <label className="text-sm font-medium">Notes</label>
          <input type="text" value={notes} onChange={(e) => setNotes(e.target.value)} className="mt-1 w-full rounded-md border px-3 py-2 text-sm" />
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <label className="text-sm font-semibold">Items</label>
          <button type="button" onClick={addItem} className="text-sm text-primary hover:underline">+ Add Item</button>
        </div>
        {items.map((it, idx) => (
          <div key={idx} className="grid gap-2 md:grid-cols-4 items-end">
            <div>
              <select value={it.product_id} onChange={(e) => updateItem(idx, 'product_id', e.target.value)} className="w-full rounded-md border px-3 py-2 text-sm" required>
                <option value="">Select Product</option>
                {productData?.products?.map((p) => <option key={p.id} value={p.id}>{p.name} ({p.product_code})</option>)}
              </select>
            </div>
            <div>
              <input type="number" placeholder="Qty" value={it.quantity || ''} onChange={(e) => updateItem(idx, 'quantity', Number(e.target.value))} className="w-full rounded-md border px-3 py-2 text-sm" min={1} required />
            </div>
            <div>
              <input type="number" placeholder="Unit Price" value={it.unit_price || ''} onChange={(e) => updateItem(idx, 'unit_price', Number(e.target.value))} className="w-full rounded-md border px-3 py-2 text-sm" min={0} required />
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium">₹{(it.quantity * it.unit_price).toLocaleString()}</span>
              {items.length > 1 && (
                <button type="button" onClick={() => removeItem(idx)} className="rounded p-1 text-red-600 hover:bg-red-50"><Trash2 className="h-4 w-4" /></button>
              )}
            </div>
          </div>
        ))}
        <div className="text-right text-sm font-semibold">Total: ₹{total.toLocaleString()}</div>
      </div>

      <div className="flex gap-2">
        <button type="submit" disabled={isLoading} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
          {isLoading ? 'Creating...' : 'Create Purchase'}
        </button>
        <button type="button" onClick={onCancel} className="rounded-lg border px-4 py-2 text-sm font-medium">Cancel</button>
      </div>
    </form>
  )
}
