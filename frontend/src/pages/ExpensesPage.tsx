import { useState } from 'react'
import { useExpenseList, useExpenseCategories, useCreateExpense, useDeleteExpense, useExpenseStats } from '@/hooks/useExpense'
import { Plus, CreditCard, TrendingDown, TrendingUp, Percent, MoreHorizontal, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { PageHeader } from '@/components/shared/PageHeader'
import { KPICard } from '@/components/shared/KPICard'
import { DataTable, type ColumnDef } from '@/components/shared/DataTable'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'
import type { Expense, CreateExpenseInput, ExpensePaymentMethod } from '@/types'
import { toastSuccess } from '@/lib/toast'

const PAYMENT_METHODS: { value: ExpensePaymentMethod; label: string }[] = [
  { value: 'cash', label: 'Cash' },
  { value: 'upi', label: 'UPI' },
  { value: 'bank_transfer', label: 'Bank Transfer' },
  { value: 'card', label: 'Card' },
  { value: 'cheque', label: 'Cheque' },
]

const statusColors: Record<string, 'default' | 'secondary' | 'outline' | 'destructive'> = {
  pending: 'outline',
  approved: 'secondary',
  paid: 'default',
  rejected: 'destructive',
}

export function ExpensesPage() {
  const [formOpen, setFormOpen] = useState(false)
  const [categoryFilter, setCategoryFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')

  const { data: stats } = useExpenseStats()
  const { data: categories } = useExpenseCategories()
  const { data, isLoading, error, refetch } = useExpenseList({
    category_id: categoryFilter || undefined,
    status: statusFilter || undefined,
    page,
    per_page: 20,
  })
  const createExp = useCreateExpense()
  const deleteExp = useDeleteExpense()

  const handleDelete = (id: string) => {
    deleteExp.mutate(id, { onSuccess: () => toastSuccess('Expense deleted') })
  }

  const handleExport = () => {
    if (!data?.expenses) return
    const csv = [
      ['Expense #', 'Category', 'Amount', 'Date', 'Payment', 'Vendor', 'Status'].join(','),
      ...data.expenses.map((e) =>
        [e.expense_number, e.category_name, e.amount, e.expense_date, e.payment_method, e.vendor_name, e.status].join(',')
      ),
    ].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `expenses-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const columns: ColumnDef<Expense, unknown>[] = [
    {
      accessorKey: 'expense_number',
      header: 'Expense #',
      cell: ({ row }) => <span className="font-mono text-xs">{row.original.expense_number}</span>,
    },
    {
      accessorKey: 'category_name',
      header: 'Category',
      cell: ({ row }) => <Badge variant="outline">{row.original.category_name}</Badge>,
    },
    {
      accessorKey: 'amount',
      header: 'Amount',
      cell: ({ row }) => <span className="font-medium">₹{row.original.amount.toLocaleString('en-IN')}</span>,
    },
    {
      accessorKey: 'expense_date',
      header: 'Date',
      cell: ({ row }) => new Date(row.original.expense_date).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' }),
    },
    {
      accessorKey: 'payment_method',
      header: 'Payment',
      cell: ({ row }) => <span className="capitalize">{row.original.payment_method.replace('_', ' ')}</span>,
    },
    {
      accessorKey: 'vendor_name',
      header: 'Vendor',
      cell: ({ row }) => row.original.vendor_name || '—',
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => (
        <Badge variant={statusColors[row.original.status] || 'outline'}>
          {row.original.status}
        </Badge>
      ),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => handleDelete(row.original.id)} className="text-destructive">
              <Trash2 className="mr-2 h-4 w-4" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ]

  if (isLoading) {
    return (
      <div className="space-y-6">
        <PageHeader title="Expenses" description="Track and manage all business expenses" />
        <LoadingState variant="page" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Expenses" description="Track and manage all business expenses" />
        <ErrorState
          title="Failed to load expenses"
          message="Please ensure the backend is running and try again."
          onRetry={() => refetch()}
        />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Expenses"
        description="Track and manage all business expenses"
        actions={
          <Button onClick={() => setFormOpen(true)} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Expense
          </Button>
        }
      />

      {/* KPIs */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        <KPICard title="Today's Expenses" value={`₹${Math.round(stats?.today_expenses ?? 0).toLocaleString('en-IN')}`} icon={TrendingDown} />
        <KPICard title="Monthly Expenses" value={`₹${Math.round(stats?.monthly_expenses ?? 0).toLocaleString('en-IN')}`} icon={CreditCard} />
        <KPICard title="Monthly Profit" value={`₹${Math.round(stats?.monthly_profit ?? 0).toLocaleString('en-IN')}`} icon={TrendingUp} />
        <KPICard title="Profit Margin" value={`${(stats?.profit_margin ?? 0).toFixed(1)}%`} icon={Percent} />
      </div>

      {/* Filters */}
      <div className="flex items-center gap-3">
        <Select value={categoryFilter || 'all'} onValueChange={(v) => { setCategoryFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[160px] h-9">
            <SelectValue placeholder="Category" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Categories</SelectItem>
            {categories?.map((c) => (
              <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select value={statusFilter || 'all'} onValueChange={(v) => { setStatusFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[130px] h-9">
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="pending">Pending</SelectItem>
            <SelectItem value="approved">Approved</SelectItem>
            <SelectItem value="paid">Paid</SelectItem>
            <SelectItem value="rejected">Rejected</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Data Table */}
      <DataTable
        columns={columns}
        data={data?.expenses || []}
        searchPlaceholder="Search expenses..."
        onSearchChange={(v) => { setSearch(v); setPage(1) }}
        searchValue={search}
        pageCount={data?.meta?.total_pages || 1}
        page={page}
        onPageChange={setPage}
        emptyTitle="No expenses recorded"
        emptyDescription="Start tracking your business expenses."
        emptyAction={{ label: 'Add Expense', onClick: () => setFormOpen(true) }}
        onExport={handleExport}
        exportLabel="Export CSV"
      />

      {/* Add Expense Dialog */}
      <AddExpenseDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        categories={categories || []}
        onSubmit={(input) => {
          createExp.mutate(input, { onSuccess: () => { setFormOpen(false); toastSuccess('Expense created') } })
        }}
        isLoading={createExp.isPending}
      />
    </div>
  )
}

function AddExpenseDialog({
  open,
  onOpenChange,
  categories,
  onSubmit,
  isLoading,
}: {
  open: boolean
  onOpenChange: (v: boolean) => void
  categories: { id: string; name: string }[]
  onSubmit: (input: CreateExpenseInput) => void
  isLoading: boolean
}) {
  const [form, setForm] = useState<CreateExpenseInput>({
    category_id: '',
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
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Add Expense</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Category *</label>
              <Select value={form.category_id} onValueChange={(v) => setForm({ ...form, category_id: v })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select" />
                </SelectTrigger>
                <SelectContent>
                  {categories.map((c) => (
                    <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Amount (₹) *</label>
              <Input
                type="number"
                value={form.amount || ''}
                onChange={(e) => setForm({ ...form, amount: Number(e.target.value) })}
                min={1}
                required
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Date *</label>
              <Input
                type="date"
                value={form.expense_date}
                onChange={(e) => setForm({ ...form, expense_date: e.target.value })}
                required
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Payment Method</label>
              <Select value={form.payment_method} onValueChange={(v) => setForm({ ...form, payment_method: v as ExpensePaymentMethod })}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {PAYMENT_METHODS.map((m) => (
                    <SelectItem key={m.value} value={m.value}>{m.label}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Vendor</label>
              <Input
                value={form.vendor_name}
                onChange={(e) => setForm({ ...form, vendor_name: e.target.value })}
                placeholder="Vendor name"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Invoice Ref</label>
              <Input
                value={form.invoice_reference || ''}
                onChange={(e) => setForm({ ...form, invoice_reference: e.target.value })}
                placeholder="Bill number"
              />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <Input
              value={form.description || ''}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              placeholder="What was this expense for?"
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={isLoading || !form.category_id || !form.amount}>
              {isLoading ? 'Creating...' : 'Add Expense'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
