import { useState } from 'react'
import { Plus, Users } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { PageHeader } from '@/components/shared/PageHeader'
import { KPICard } from '@/components/shared/KPICard'
import { DataTable, type ColumnDef } from '@/components/shared/DataTable'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'
import { StaffFormDialog } from '@/components/staff/StaffFormDialog'
import { DeleteStaffDialog } from '@/components/staff/DeleteStaffDialog'
import { useStaffList, useStaffStats, useCreateStaff, useUpdateStaff, useDeleteStaff } from '@/hooks/useStaff'
import type { Staff, CreateStaffInput, UpdateStaffInput } from '@/types'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreHorizontal, Pencil, Trash2 } from 'lucide-react'
import { toastSuccess } from '@/lib/toast'

const designationLabels: Record<string, string> = {
  stylist: 'Stylist',
  assistant: 'Assistant',
  receptionist: 'Receptionist',
  manager: 'Manager',
}

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(amount)

export function StaffPage() {
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [formOpen, setFormOpen] = useState(false)
  const [editingStaff, setEditingStaff] = useState<Staff | null>(null)
  const [deletingStaff, setDeletingStaff] = useState<Staff | null>(null)

  const { data, isLoading, error, refetch } = useStaffList({
    page,
    per_page: 20,
    search: search || undefined,
    status: statusFilter || undefined,
  })
  const { data: stats } = useStaffStats()

  const createMutation = useCreateStaff()
  const updateMutation = useUpdateStaff()
  const deleteMutation = useDeleteStaff()

  const handleAdd = () => {
    setEditingStaff(null)
    setFormOpen(true)
  }

  const handleEdit = (staff: Staff) => {
    setEditingStaff(staff)
    setFormOpen(true)
  }

  const handleDelete = (staff: Staff) => {
    setDeletingStaff(staff)
  }

  const handleFormSubmit = (values: Record<string, unknown>) => {
    if (editingStaff) {
      updateMutation.mutate(
        { id: editingStaff.id, input: values as unknown as UpdateStaffInput },
        { onSuccess: () => { setFormOpen(false); toastSuccess('Staff member updated') } }
      )
    } else {
      createMutation.mutate(values as unknown as CreateStaffInput, {
        onSuccess: () => { setFormOpen(false); toastSuccess('Staff member created') },
      })
    }
  }

  const handleDeleteConfirm = () => {
    if (!deletingStaff) return
    deleteMutation.mutate(deletingStaff.id, {
      onSuccess: () => { setDeletingStaff(null); toastSuccess('Staff member deleted') },
    })
  }

  const handleExport = () => {
    if (!data?.staff) return
    const csv = [
      ['Code', 'Name', 'Phone', 'Designation', 'Base Salary', 'Commission %', 'Status'].join(','),
      ...data.staff.map((s) =>
        [s.staff_code, s.full_name, s.phone, s.designation, s.base_salary, s.commission_percentage, s.status].join(',')
      ),
    ].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `staff-export-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const columns: ColumnDef<Staff, unknown>[] = [
    {
      accessorKey: 'staff_code',
      header: 'Code',
      cell: ({ row }) => <span className="font-mono text-xs">{row.original.staff_code}</span>,
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
      accessorKey: 'designation',
      header: 'Designation',
      cell: ({ row }) => designationLabels[row.original.designation] || row.original.designation,
    },
    {
      accessorKey: 'base_salary',
      header: 'Salary',
      cell: ({ row }) => formatCurrency(row.original.base_salary),
    },
    {
      accessorKey: 'commission_percentage',
      header: 'Comm %',
      cell: ({ row }) => `${row.original.commission_percentage}%`,
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
        <PageHeader title="Staff" description="Manage your salon staff members" />
        <LoadingState variant="page" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Staff" description="Manage your salon staff members" />
        <ErrorState
          title="Failed to load staff"
          message="Please ensure the backend is running and try again."
          onRetry={() => refetch()}
        />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <PageHeader
        title="Staff"
        description="Manage your salon staff members"
        actions={
          <Button onClick={handleAdd} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Staff
          </Button>
        }
      />

      {/* Stats */}
      {stats && (
        <div className="grid gap-4 grid-cols-3">
          <KPICard title="Total Staff" value={stats.total} icon={Users} />
          <KPICard title="Active" value={stats.active} icon={Users} className="border-green-200 dark:border-green-800" />
          <KPICard title="Inactive" value={stats.inactive} icon={Users} className="border-red-200 dark:border-red-800" />
        </div>
      )}

      {/* Filters */}
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
        data={data?.staff || []}
        searchPlaceholder="Search by name, phone, or code..."
        onSearchChange={(v) => { setSearch(v); setPage(1) }}
        searchValue={search}
        pageCount={data?.meta?.total_pages || 1}
        page={page}
        onPageChange={setPage}
        emptyTitle="No staff members yet"
        emptyDescription="Get started by adding your first staff member."
        emptyAction={{ label: 'Add Staff', onClick: handleAdd }}
        onExport={handleExport}
        exportLabel="Export CSV"
      />

      {/* Dialogs */}
      <StaffFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        staff={editingStaff}
        onSubmit={handleFormSubmit}
        isLoading={createMutation.isPending || updateMutation.isPending}
      />

      <DeleteStaffDialog
        open={!!deletingStaff}
        onOpenChange={(open) => { if (!open) setDeletingStaff(null) }}
        staff={deletingStaff}
        onConfirm={handleDeleteConfirm}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
