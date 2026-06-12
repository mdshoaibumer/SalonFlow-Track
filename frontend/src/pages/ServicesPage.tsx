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
import { Plus, Scissors, MoreHorizontal, Pencil, Trash2 } from 'lucide-react'
import { useServiceList, useCreateService, useUpdateService, useDeleteService } from '@/hooks/useServices'
import { ServiceFormDialog } from '@/components/services/ServiceFormDialog'
import { DeleteServiceDialog } from '@/components/services/DeleteServiceDialog'
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
import type { Service, CreateServiceInput, UpdateServiceInput } from '@/types'

const categoryLabels: Record<string, string> = {
  hair: 'Hair',
  facial: 'Facial',
  skin: 'Skin',
  spa: 'Spa',
  massage: 'Massage',
  coloring: 'Coloring',
  treatment: 'Treatment',
  other: 'Other',
}

export function ServicesPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [page, setPage] = useState(1)
  const [formOpen, setFormOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [selectedService, setSelectedService] = useState<Service | null>(null)

  const { data, isLoading, error, refetch } = useServiceList({
    page,
    per_page: 20,
    search: search || undefined,
    status: statusFilter || undefined,
    category: categoryFilter || undefined,
  })

  const createMutation = useCreateService()
  const updateMutation = useUpdateService()
  const deleteMutation = useDeleteService()

  const handleAdd = () => {
    setSelectedService(null)
    setFormOpen(true)
  }

  const handleEdit = (service: Service) => {
    setSelectedService(service)
    setFormOpen(true)
  }

  const handleDelete = (service: Service) => {
    setSelectedService(service)
    setDeleteOpen(true)
  }

  const handleFormSubmit = (values: Record<string, unknown>) => {
    if (selectedService) {
      updateMutation.mutate(
        { id: selectedService.id, input: values as unknown as UpdateServiceInput },
        { onSuccess: () => { setFormOpen(false); toastSuccess('Service updated') } }
      )
    } else {
      createMutation.mutate(values as unknown as CreateServiceInput, {
        onSuccess: () => { setFormOpen(false); toastSuccess('Service created') },
      })
    }
  }

  const handleDeleteConfirm = () => {
    if (selectedService) {
      deleteMutation.mutate(selectedService.id, {
        onSuccess: () => { setDeleteOpen(false); toastSuccess('Service deleted') },
      })
    }
  }

  const handleExport = () => {
    if (!data?.services) return
    const csv = [
      ['Code', 'Name', 'Category', 'Duration', 'Price', 'Commission', 'Status'].join(','),
      ...data.services.map((s) =>
        [s.service_code, s.name, s.category, `${s.duration_minutes}min`, s.price, `${s.commission_value}${s.commission_type === 'percentage' ? '%' : ''}`, s.status].join(',')
      ),
    ].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `services-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const columns: ColumnDef<Service, unknown>[] = [
    {
      accessorKey: 'service_code',
      header: 'Code',
      cell: ({ row }) => <span className="font-mono text-xs">{row.original.service_code}</span>,
    },
    {
      accessorKey: 'name',
      header: 'Service',
      cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
    },
    {
      accessorKey: 'category',
      header: 'Category',
      cell: ({ row }) => (
        <Badge variant="outline">{categoryLabels[row.original.category] || row.original.category}</Badge>
      ),
    },
    {
      accessorKey: 'duration_minutes',
      header: 'Duration',
      cell: ({ row }) => `${row.original.duration_minutes} min`,
    },
    {
      accessorKey: 'price',
      header: 'Price',
      cell: ({ row }) => `₹${row.original.price.toLocaleString('en-IN')}`,
    },
    {
      accessorKey: 'commission_value',
      header: 'Commission',
      cell: ({ row }) => {
        const s = row.original
        return s.commission_type === 'percentage' ? `${s.commission_value}%` : `₹${s.commission_value}`
      },
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
        <PageHeader title="Services" description="Manage your salon service offerings" />
        <LoadingState variant="page" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Services" description="Manage your salon service offerings" />
        <ErrorState
          title="Failed to load services"
          message="Please ensure the backend is running and try again."
          onRetry={() => refetch()}
        />
      </div>
    )
  }

  const services = data?.services || []
  const activeCount = services.filter((s) => s.status === 'active').length
  const avgPrice = services.length ? Math.round(services.reduce((sum, s) => sum + s.price, 0) / services.length) : 0

  return (
    <div className="space-y-6">
      <PageHeader
        title="Services"
        description="Manage your salon service offerings"
        actions={
          <Button onClick={handleAdd} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Service
          </Button>
        }
      />

      {/* KPIs */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        <KPICard title="Total Services" value={data?.meta?.total ?? services.length} icon={Scissors} />
        <KPICard title="Active" value={activeCount} icon={Scissors} />
        <KPICard title="Categories" value={new Set(services.map((s) => s.category)).size} icon={Scissors} />
        <KPICard title="Avg Price" value={`₹${avgPrice.toLocaleString('en-IN')}`} icon={Scissors} />
      </div>

      {/* Filters */}
      <div className="flex items-center gap-3">
        <Select value={categoryFilter || 'all'} onValueChange={(v) => { setCategoryFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[150px] h-9">
            <SelectValue placeholder="Category" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Categories</SelectItem>
            <SelectItem value="hair">Hair</SelectItem>
            <SelectItem value="facial">Facial</SelectItem>
            <SelectItem value="skin">Skin</SelectItem>
            <SelectItem value="spa">Spa</SelectItem>
            <SelectItem value="massage">Massage</SelectItem>
            <SelectItem value="coloring">Coloring</SelectItem>
            <SelectItem value="treatment">Treatment</SelectItem>
            <SelectItem value="other">Other</SelectItem>
          </SelectContent>
        </Select>
        <Select value={statusFilter || 'all'} onValueChange={(v) => { setStatusFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[130px] h-9">
            <SelectValue placeholder="Status" />
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
        data={services}
        searchPlaceholder="Search services..."
        onSearchChange={(v) => { setSearch(v); setPage(1) }}
        searchValue={search}
        pageCount={data?.meta?.total_pages || 1}
        page={page}
        onPageChange={setPage}
        emptyTitle="No services yet"
        emptyDescription="Add your salon service offerings to get started."
        emptyAction={{ label: 'Add Service', onClick: handleAdd }}
        onExport={handleExport}
        exportLabel="Export CSV"
      />

      <ServiceFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        service={selectedService}
        onSubmit={handleFormSubmit}
        isLoading={createMutation.isPending || updateMutation.isPending}
      />

      <DeleteServiceDialog
        open={deleteOpen}
        onOpenChange={setDeleteOpen}
        service={selectedService}
        onConfirm={handleDeleteConfirm}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
