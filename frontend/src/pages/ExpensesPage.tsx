import { useState } from 'react'
import { useExpenseList, useExpenseCategories, useCreateExpense, useDeleteExpense, useExpenseStats } from '@/hooks/useExpense'
import { Plus, Trash2, IndianRupee, TrendingDown, TrendingUp, Percent } from 'lucide-react'
import type { CreateExpenseInput, ExpensePaymentMethod } from '@/types'

const PAYMENT_METHODS: { value: ExpensePaymentMethod; label: string }[] = [
  { value: 'cash', label: 'Cash' },
  { value: 'upi', label: 'UPI' },
  { value: 'bank_transfer', label: 'Bank Transfer' },
  { value: 'card', label: 'Card' },
  { value: 'cheque', label: 'Cheque' },
]

export function ExpensesPage() {
  const [showForm, setShowForm] = useState(false)
  const [filters, setFilters] = useState({ category_id: '', status: '', payment_method: '', page: 1 })

  const { data: stats } = useExpenseStats()
  const { data: categories } = useExpenseCategories()
  const { data, isLoading } = useExpenseList({ ...filters, per_page: 20 })
  const createExp = useCreateExpense()
  const deleteExp = useDeleteExpense()

  const handleCreate = (input: CreateExpenseInput) => {
    createExp.mutate(input, { onSuccess: () => setShowForm(false) })
  }

  const handleDelete = (id: string, num: string) => {
    if (confirm(`Delete expense ${num}?`)) {
      deleteExp.mutate(id)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Expenses</h1>
          <p className="text-muted-foreground">Track and manage all business expenses</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="h-4 w-4" />
          Add Expense
        </button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <StatCard title="Today's Expenses" value={`₹${Math.round(stats?.today_expenses ?? 0).toLocaleString()}`} icon={<TrendingDown className="h-4 w-4 text-red-500" />} />
        <StatCard title="Monthly Expenses" value={`₹${Math.round(stats?.monthly_expenses ?? 0).toLocaleString()}`} icon={<IndianRupee className="h-4 w-4" />} />
        <StatCard title="Monthly Profit" value={`₹${Math.round(stats?.monthly_profit ?? 0).toLocaleString()}`} icon={<TrendingUp className="h-4 w-4 text-green-500" />} />
        <StatCard title="Profit Margin" value={`${(stats?.profit_margin ?? 0).toFixed(1)}%`} icon={<Percent className="h-4 w-4" />} />
      </div>

      {/* Add Expense Form */}
      {showForm && (
        <ExpenseForm
          categories={categories || []}
          onSubmit={handleCreate}
          onCancel={() => setShowForm(false)}
          isLoading={createExp.isPending}
        />
      )}

      {/* Filters */}
      <div className="flex gap-2 flex-wrap">
        <select
          value={filters.category_id}
          onChange={(e) => setFilters({ ...filters, category_id: e.target.value, page: 1 })}
          className="rounded-md border px-3 py-2 text-sm"
        >
          <option value="">All Categories</option>
          {categories?.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>
        <select
          value={filters.status}
          onChange={(e) => setFilters({ ...filters, status: e.target.value, page: 1 })}
          className="rounded-md border px-3 py-2 text-sm"
        >
          <option value="">All Status</option>
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="paid">Paid</option>
          <option value="rejected">Rejected</option>
        </select>
        <select
          value={filters.payment_method}
          onChange={(e) => setFilters({ ...filters, payment_method: e.target.value, page: 1 })}
          className="rounded-md border px-3 py-2 text-sm"
        >
          <option value="">All Methods</option>
          {PAYMENT_METHODS.map((m) => (
            <option key={m.value} value={m.value}>{m.label}</option>
          ))}
        </select>
      </div>

      {/* Expense Table */}
      <div className="rounded-lg border bg-card">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Expense #</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Category</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Amount</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Date</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Payment</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Vendor</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Status</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Actions</th>
              </tr>
            </thead>
            <tbody>
              {isLoading && (
                <tr><td colSpan={8} className="px-6 py-8 text-center text-muted-foreground">Loading...</td></tr>
              )}
              {!isLoading && (!data?.expenses || data.expenses.length === 0) && (
                <tr><td colSpan={8} className="px-6 py-8 text-center text-muted-foreground">No expenses found</td></tr>
              )}
              {data?.expenses?.map((exp) => (
                <tr key={exp.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm font-mono">{exp.expense_number}</td>
                  <td className="px-6 py-4 text-sm">{exp.category_name}</td>
                  <td className="px-6 py-4 text-sm text-right font-medium">₹{exp.amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm">{exp.expense_date}</td>
                  <td className="px-6 py-4 text-sm capitalize">{exp.payment_method.replace('_', ' ')}</td>
                  <td className="px-6 py-4 text-sm">{exp.vendor_name || '—'}</td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={exp.status} />
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button
                      onClick={() => handleDelete(exp.id, exp.expense_number)}
                      className="rounded p-1 text-red-600 hover:bg-red-50"
                      title="Delete"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {/* Pagination */}
        {data && data.meta.total_pages > 1 && (
          <div className="flex items-center justify-between border-t px-6 py-3">
            <span className="text-sm text-muted-foreground">
              Page {data.meta.page} of {data.meta.total_pages} ({data.meta.total} total)
            </span>
            <div className="flex gap-2">
              <button
                disabled={filters.page <= 1}
                onClick={() => setFilters({ ...filters, page: filters.page - 1 })}
                className="rounded border px-3 py-1 text-sm disabled:opacity-50"
              >
                Previous
              </button>
              <button
                disabled={filters.page >= data.meta.total_pages}
                onClick={() => setFilters({ ...filters, page: filters.page + 1 })}
                className="rounded border px-3 py-1 text-sm disabled:opacity-50"
              >
                Next
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-blue-100 text-blue-800',
    paid: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
  }
  return (
    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${colors[status] || 'bg-gray-100 text-gray-800'}`}>
      {status}
    </span>
  )
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

function ExpenseForm({
  categories,
  onSubmit,
  onCancel,
  isLoading,
}: {
  categories: { id: string; name: string }[]
  onSubmit: (input: CreateExpenseInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [form, setForm] = useState<CreateExpenseInput>({
    category_id: categories[0]?.id || '',
    amount: 0,
    expense_date: new Date().toISOString().split('T')[0] ?? '',
    payment_method: 'cash',
    vendor_name: '',
    invoice_reference: '',
    description: '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="rounded-lg border bg-card p-6 space-y-4">
      <h3 className="text-lg font-semibold">New Expense</h3>
      <div className="grid gap-4 md:grid-cols-3">
        <div>
          <label className="text-sm font-medium">Category</label>
          <select
            value={form.category_id}
            onChange={(e) => setForm({ ...form, category_id: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            required
          >
            <option value="">Select category</option>
            {categories.map((c) => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Amount (₹)</label>
          <input
            type="number"
            value={form.amount || ''}
            onChange={(e) => setForm({ ...form, amount: Number(e.target.value) })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            min={1}
            required
          />
        </div>
        <div>
          <label className="text-sm font-medium">Date</label>
          <input
            type="date"
            value={form.expense_date}
            onChange={(e) => setForm({ ...form, expense_date: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            required
          />
        </div>
        <div>
          <label className="text-sm font-medium">Payment Method</label>
          <select
            value={form.payment_method}
            onChange={(e) => setForm({ ...form, payment_method: e.target.value as ExpensePaymentMethod })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
          >
            {PAYMENT_METHODS.map((m) => (
              <option key={m.value} value={m.value}>{m.label}</option>
            ))}
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Vendor</label>
          <input
            type="text"
            value={form.vendor_name}
            onChange={(e) => setForm({ ...form, vendor_name: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            placeholder="Vendor name"
          />
        </div>
        <div>
          <label className="text-sm font-medium">Invoice Ref</label>
          <input
            type="text"
            value={form.invoice_reference || ''}
            onChange={(e) => setForm({ ...form, invoice_reference: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            placeholder="Bill/Invoice number"
          />
        </div>
        <div className="md:col-span-3">
          <label className="text-sm font-medium">Description</label>
          <input
            type="text"
            value={form.description || ''}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            placeholder="What was this expense for?"
          />
        </div>
      </div>
      <div className="flex gap-2">
        <button type="submit" disabled={isLoading} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
          {isLoading ? 'Creating...' : 'Add Expense'}
        </button>
        <button type="button" onClick={onCancel} className="rounded-lg border px-4 py-2 text-sm font-medium">Cancel</button>
      </div>
    </form>
  )
}
