import { useState } from 'react'
import { Search, Plus, Users } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { StaffTable } from '@/components/staff/StaffTable'
import { StaffFormDialog } from '@/components/staff/StaffFormDialog'
import { DeleteStaffDialog } from '@/components/staff/DeleteStaffDialog'
import { useStaffList, useCreateStaff, useUpdateStaff, useDeleteStaff } from '@/hooks/useStaff'
import type { Staff, CreateStaffInput, UpdateStaffInput } from '@/types'

export function StaffPage() {
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [formOpen, setFormOpen] = useState(false)
  const [editingStaff, setEditingStaff] = useState<Staff | null>(null)
  const [deletingStaff, setDeletingStaff] = useState<Staff | null>(null)

  const { data, isLoading, error } = useStaffList({
    page,
    per_page: 20,
    search: search || undefined,
    status: statusFilter || undefined,
  })

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
        { onSuccess: () => setFormOpen(false) }
      )
    } else {
      createMutation.mutate(values as unknown as CreateStaffInput, {
        onSuccess: () => setFormOpen(false),
      })
    }
  }

  const handleDeleteConfirm = () => {
    if (!deletingStaff) return
    deleteMutation.mutate(deletingStaff.id, {
      onSuccess: () => setDeletingStaff(null),
    })
  }

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setPage(1)
  }

  const totalPages = data?.meta?.total_pages || 1

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Staff</h1>
          <p className="text-muted-foreground">
            Manage your salon staff members
          </p>
        </div>
        <Button onClick={handleAdd}>
          <Plus className="mr-2 h-4 w-4" />
          Add Staff
        </Button>
      </div>

      {/* Filters */}
      <div className="flex items-center gap-4">
        <form onSubmit={handleSearchSubmit} className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Search by name, phone, or code..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-10"
          />
        </form>

        <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v === 'all' ? '' : v); setPage(1) }}>
          <SelectTrigger className="w-[140px]">
            <SelectValue placeholder="All Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="active">Active</SelectItem>
            <SelectItem value="inactive">Inactive</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Content */}
      {isLoading ? (
        <div className="flex items-center justify-center h-48">
          <div className="text-muted-foreground">Loading staff...</div>
        </div>
      ) : error ? (
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-6 text-center">
          <p className="text-destructive">
            Failed to load staff. Please ensure the backend is running.
          </p>
        </div>
      ) : (
        <>
          {data && data.staff.length === 0 && !search && !statusFilter ? (
            <div className="flex flex-col items-center justify-center rounded-lg border bg-card p-12 text-center">
              <Users className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-semibold">No staff members yet</h3>
              <p className="text-muted-foreground mt-1 mb-4">
                Get started by adding your first staff member.
              </p>
              <Button onClick={handleAdd}>
                <Plus className="mr-2 h-4 w-4" />
                Add Staff
              </Button>
            </div>
          ) : (
            <StaffTable
              staff={data?.staff || []}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between">
              <p className="text-sm text-muted-foreground">
                Showing page {page} of {totalPages} ({data?.meta?.total || 0} total)
              </p>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
                >
                  Previous
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}

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
