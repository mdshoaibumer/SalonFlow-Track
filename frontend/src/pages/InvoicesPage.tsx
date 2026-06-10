import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { useNavigate } from 'react-router-dom'
import { IndianRupee, Receipt, TrendingUp, Plus } from 'lucide-react'
import { useInvoiceList, useInvoiceStats } from '@/hooks/useInvoices'
import { PageHeader } from '@/components/shared/PageHeader'
import { KPICard } from '@/components/shared/KPICard'
import { DataTable, type ColumnDef } from '@/components/shared/DataTable'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'
import type { Invoice } from '@/types'

export function InvoicesPage() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)

  const { data, isLoading, error, refetch } = useInvoiceList({
    page,
    per_page: 20,
    search: search || undefined,
    payment_status: statusFilter || undefined,
  })

  const { data: stats } = useInvoiceStats()

  const handleExport = () => {
    if (!data?.invoices) return
    const csv = [
      ['Invoice #', 'Date', 'Subtotal', 'Discount', 'Grand Total', 'Payment Method', 'Status'].join(','),
      ...data.invoices.map((inv) =>
        [inv.invoice_number, inv.invoice_date, inv.subtotal, inv.discount, inv.grand_total, inv.payment_method, inv.payment_status].join(',')
      ),
    ].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `invoices-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const columns: ColumnDef<Invoice, unknown>[] = [
    {
      accessorKey: 'invoice_number',
      header: 'Invoice #',
      cell: ({ row }) => <span className="font-mono text-xs font-medium">{row.original.invoice_number}</span>,
    },
    {
      accessorKey: 'invoice_date',
      header: 'Date',
      cell: ({ row }) => new Date(row.original.invoice_date).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' }),
    },
    {
      accessorKey: 'subtotal',
      header: 'Subtotal',
      cell: ({ row }) => `₹${row.original.subtotal.toLocaleString('en-IN')}`,
    },
    {
      accessorKey: 'discount',
      header: 'Discount',
      cell: ({ row }) => row.original.discount > 0 ? `₹${row.original.discount.toLocaleString('en-IN')}` : '-',
    },
    {
      accessorKey: 'grand_total',
      header: 'Total',
      cell: ({ row }) => <span className="font-medium">₹{row.original.grand_total.toLocaleString('en-IN')}</span>,
    },
    {
      accessorKey: 'payment_method',
      header: 'Payment',
      cell: ({ row }) => <span className="capitalize">{row.original.payment_method || '—'}</span>,
    },
    {
      accessorKey: 'payment_status',
      header: 'Status',
      cell: ({ row }) => {
        const s = row.original.payment_status
        return (
          <Badge variant={s === 'paid' ? 'default' : s === 'partial' ? 'secondary' : 'outline'}>
            {s}
          </Badge>
        )
      },
    },
  ]

  if (isLoading) {
    return (
      <div className="space-y-6">
        <PageHeader title="Invoices" description="View and manage billing invoices" />
        <LoadingState variant="page" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Invoices" description="View and manage billing invoices" />
        <ErrorState
          title="Failed to load invoices"
          message="Please ensure the backend is running and try again."
          onRetry={() => refetch()}
        />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Invoices"
        description="View and manage billing invoices"
        actions={
          <Button onClick={() => navigate('/billing')} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            New Invoice
          </Button>
        }
      />

      {/* KPIs */}
      {stats && (
        <div className="grid gap-4 grid-cols-3">
          <KPICard title="Today's Revenue" value={`₹${stats.today_revenue.toLocaleString('en-IN')}`} icon={IndianRupee} />
          <KPICard title="Today's Invoices" value={stats.today_invoices} icon={Receipt} />
          <KPICard title="Avg Bill Value" value={`₹${Math.round(stats.avg_bill_value).toLocaleString('en-IN')}`} icon={TrendingUp} />
        </div>
      )}

      {/* Filters */}
      <div className="flex items-center gap-3">
        <Select value={statusFilter || 'all'} onValueChange={(v) => { setStatusFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[140px] h-9">
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="paid">Paid</SelectItem>
            <SelectItem value="pending">Pending</SelectItem>
            <SelectItem value="partial">Partial</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Data Table */}
      <DataTable
        columns={columns}
        data={data?.invoices || []}
        searchPlaceholder="Search by invoice number..."
        onSearchChange={(v) => { setSearch(v); setPage(1) }}
        searchValue={search}
        pageCount={data?.meta?.total_pages || 1}
        page={page}
        onPageChange={setPage}
        emptyTitle="No invoices yet"
        emptyDescription="Create invoices from the Billing page."
        emptyAction={{ label: 'Create Invoice', onClick: () => navigate('/billing') }}
        onExport={handleExport}
        exportLabel="Export CSV"
      />
    </div>
  )
}
