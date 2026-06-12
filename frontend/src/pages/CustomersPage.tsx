import { useState } from 'react'
import { toastSuccess } from '@/lib/toast'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Plus, UserCircle, MoreHorizontal, Pencil, Trash2 } from 'lucide-react'
import { useCustomerList, useCreateCustomer, useUpdateCustomer, useDeleteCustomer } from '@/hooks/useCustomers'
import { CustomerFormDialog } from '@/components/customers/CustomerFormDialog'
import { DeleteCustomerDialog } from '@/components/customers/DeleteCustomerDialog'
import { PageHeader } from '@/components/shared/PageHeader'
import { KPICard } from '@/components/shared/KPICard'
import { DataTable, type ColumnDef } from '@/components/shared/DataTable'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { Customer, CreateCustomerInput, UpdateCustomerInput } from '@/types'

export function CustomersPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const [formOpen, setFormOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null)

  const { data, isLoading, error, refetch } = useCustomerList({
    page,
    per_page: 20,
    search: search || undefined,
    status: statusFilter || undefined,
  })

  const createMutation = useCreateCustomer()
  const updateMutation = useUpdateCustomer()
  const deleteMutation = useDeleteCustomer()

  const handleAdd = () => {
    setSelectedCustomer(null)
    setFormOpen(true)
  }

  const handleEdit = (customer: Customer) => {
    setSelectedCustomer(customer)
    setFormOpen(true)
  }

  const handleDelete = (customer: Customer) => {
    setSelectedCustomer(customer)
    setDeleteOpen(true)
  }

  const handleFormSubmit = (values: Record<string, unknown>) => {
    if (selectedCustomer) {
      updateMutation.mutate(
        { id: selectedCustomer.id, input: values as unknown as UpdateCustomerInput },
        { onSuccess: () => { setFormOpen(false); toastSuccess('Customer updated') } }
      )
    } else {
      createMutation.mutate(values as unknown as CreateCustomerInput, {
        onSuccess: () => { setFormOpen(false); toastSuccess('Customer added') },
      })
    }
  }

  const handleDeleteConfirm = () => {
    if (selectedCustomer) {
      deleteMutation.mutate(selectedCustomer.id, {
        onSuccess: () => { setDeleteOpen(false); toastSuccess('Customer deleted') },
      })
    }
  }

  const handleExport = () => {
    if (!data?.customers) return
    const csv = [
      ['Code', 'Name', 'Phone', 'Email', 'Visits', 'Total Spent', 'Status'].join(','),
      ...data.customers.map((c) =>
        [c.customer_code, c.full_name, c.phone, c.email, c.total_visits, c.total_spent, c.status].join(',')
      ),
    ].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `customers-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const columns: ColumnDef<Customer, unknown>[] = [
    {
      accessorKey: 'customer_code',
      header: 'Code',
      cell: ({ row }) => <span className="font-mono text-xs">{row.original.customer_code}</span>,
    },
    {
      accessorKey: 'full_name',
      header: 'Name',
      cell: ({ row }) => <span className="font-medium">{row.original.full_name}</span>,
    },
    {
      accessorKey: 'phone',
      header: 'Phone',
    },
    {
      accessorKey: 'total_visits',
      header: 'Visits',
    },
    {
      accessorKey: 'total_spent',
      header: 'Total Spent',
      cell: ({ row }) => `₹${row.original.total_spent.toLocaleString('en-IN')}`,
    },
    {
      accessorKey: 'last_visit_date',
      header: 'Last Visit',
      cell: ({ row }) =>
        row.original.last_visit_date
          ? new Date(row.original.last_visit_date).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' })
          : '-',
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => (
        <Badge variant={row.original.status === 'active' ? 'default' : 'secondary'}>
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
            <DropdownMenuItem onClick={() => handleEdit(row.original)}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleDelete(row.original)} className="text-destructive">
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
        <PageHeader title="Customers" description="Manage your customer database" />
        <LoadingState variant="page" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Customers" description="Manage your customer database" />
        <ErrorState
          title="Failed to load customers"
          message="Please ensure the backend is running and try again."
          onRetry={() => refetch()}
        />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Customers"
        description="Manage your customer database"
        actions={
          <Button onClick={handleAdd} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Customer
          </Button>
        }
      />

      {/* KPI Stats */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        <KPICard title="Total Customers" value={data?.meta?.total ?? 0} icon={UserCircle} />
        <KPICard
          title="Active"
          value={data?.customers.filter(c => c.status === 'active').length ?? 0}
          icon={UserCircle}
        />
        <KPICard
          title="Total Revenue"
          value={`₹${(data?.customers.reduce((sum, c) => sum + c.total_spent, 0) ?? 0).toLocaleString('en-IN')}`}
          icon={UserCircle}
        />
        <KPICard
          title="Avg Visits"
          value={data?.customers.length ? Math.round(data.customers.reduce((sum, c) => sum + c.total_visits, 0) / data.customers.length) : 0}
          icon={UserCircle}
        />
      </div>

      {/* Filter */}
      <div className="flex items-center gap-3">
        <Select value={statusFilter || 'all'} onValueChange={(v) => { setStatusFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[140px] h-9">
            <SelectValue placeholder="All Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="active">Active</SelectItem>
            <SelectItem value="inactive">Inactive</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Data Table */}
      <DataTable
        columns={columns}
        data={data?.customers || []}
        searchPlaceholder="Search by name, phone, or code..."
        onSearchChange={(v) => { setSearch(v); setPage(1) }}
        searchValue={search}
        pageCount={data?.meta?.total_pages || 1}
        page={page}
        onPageChange={setPage}
        emptyTitle="No customers yet"
        emptyDescription="Start building your customer database."
        emptyAction={{ label: 'Add Customer', onClick: handleAdd }}
        onExport={handleExport}
        exportLabel="Export CSV"
      />

      <CustomerFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        customer={selectedCustomer}
        onSubmit={handleFormSubmit}
        isLoading={createMutation.isPending || updateMutation.isPending}
      />

      <DeleteCustomerDialog
        open={deleteOpen}
        onOpenChange={setDeleteOpen}
        customer={selectedCustomer}
        onConfirm={handleDeleteConfirm}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
