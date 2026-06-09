import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Plus, Search } from 'lucide-react'
import { useServiceList, useCreateService, useUpdateService, useDeleteService } from '@/hooks/useServices'
import { ServiceTable } from '@/components/services/ServiceTable'
import { ServiceFormDialog } from '@/components/services/ServiceFormDialog'
import { DeleteServiceDialog } from '@/components/services/DeleteServiceDialog'
import { ServiceStatsWidget } from '@/components/services/ServiceStatsWidget'
import type { Service, CreateServiceInput, UpdateServiceInput } from '@/types'

export function ServicesPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [page, setPage] = useState(1)
  const [formOpen, setFormOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [selectedService, setSelectedService] = useState<Service | null>(null)

  const { data, isLoading, error } = useServiceList({
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
        { onSuccess: () => setFormOpen(false) }
      )
    } else {
      createMutation.mutate(values as unknown as CreateServiceInput, {
        onSuccess: () => setFormOpen(false),
      })
    }
  }

  const handleDeleteConfirm = () => {
    if (selectedService) {
      deleteMutation.mutate(selectedService.id, {
        onSuccess: () => setDeleteOpen(false),
      })
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Services</h1>
          <p className="text-muted-foreground">Manage your salon service offerings</p>
        </div>
        <Button onClick={handleAdd}>
          <Plus className="mr-2 h-4 w-4" /> Add Service
        </Button>
      </div>

      <ServiceStatsWidget />

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search services..."
            className="pl-9"
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1) }}
          />
        </div>
        <Select value={categoryFilter} onValueChange={(v) => { setCategoryFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[150px]">
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
        <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[130px]">
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="active">Active</SelectItem>
            <SelectItem value="inactive">Inactive</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {error && (
        <div className="rounded-lg border border-destructive bg-destructive/10 p-4">
          <p className="text-sm text-destructive">Failed to load services. Is the backend running?</p>
        </div>
      )}

      {isLoading ? (
        <div className="rounded-lg border bg-card p-12 text-center">
          <p className="text-muted-foreground">Loading services...</p>
        </div>
      ) : (
        data && <ServiceTable services={data.services} onEdit={handleEdit} onDelete={handleDelete} />
      )}

      {data && data.meta.total_pages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Page {data.meta.page} of {data.meta.total_pages} ({data.meta.total} total)
          </p>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>
              Previous
            </Button>
            <Button variant="outline" size="sm" disabled={page >= data.meta.total_pages} onClick={() => setPage(page + 1)}>
              Next
            </Button>
          </div>
        </div>
      )}

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
