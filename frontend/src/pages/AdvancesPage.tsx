import { useState } from 'react'
import { useAdvanceList, useCreateAdvance, useApproveAdvance, useRejectAdvance } from '@/hooks/useSalary'
import { Plus, Check, X } from 'lucide-react'
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
import { PageHeader } from '@/components/shared/PageHeader'
import { DataTable, type ColumnDef } from '@/components/shared/DataTable'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'
import type { CreateAdvanceInput } from '@/types'

interface Advance {
  id: string
  staff_name: string
  amount: number
  advance_date: string
  reason: string
  recovered_amount: number
  remaining_amount: number
  status: string
}

const statusVariants: Record<string, 'default' | 'secondary' | 'outline' | 'destructive'> = {
  pending: 'outline',
  approved: 'secondary',
  recovering: 'secondary',
  recovered: 'default',
  rejected: 'destructive',
}

export function AdvancesPage() {
  const [formOpen, setFormOpen] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const { data, isLoading, error, refetch } = useAdvanceList({ status: statusFilter || undefined })
  const createAdv = useCreateAdvance()
  const approveAdv = useApproveAdvance()
  const rejectAdv = useRejectAdvance()

  const handleCreate = (input: CreateAdvanceInput) => {
    createAdv.mutate(input, { onSuccess: () => setFormOpen(false) })
  }

  const columns: ColumnDef<Advance, unknown>[] = [
    {
      accessorKey: 'staff_name',
      header: 'Staff',
      cell: ({ row }) => <span className="font-medium">{row.original.staff_name}</span>,
    },
    {
      accessorKey: 'amount',
      header: 'Amount',
      cell: ({ row }) => `₹${row.original.amount.toLocaleString('en-IN')}`,
    },
    {
      accessorKey: 'advance_date',
      header: 'Date',
      cell: ({ row }) => new Date(row.original.advance_date).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' }),
    },
    {
      accessorKey: 'reason',
      header: 'Reason',
      cell: ({ row }) => row.original.reason || '—',
    },
    {
      accessorKey: 'recovered_amount',
      header: 'Recovered',
      cell: ({ row }) => <span className="text-green-600">₹{row.original.recovered_amount.toLocaleString('en-IN')}</span>,
    },
    {
      accessorKey: 'remaining_amount',
      header: 'Remaining',
      cell: ({ row }) => <span className="text-orange-600">₹{row.original.remaining_amount.toLocaleString('en-IN')}</span>,
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => (
        <Badge variant={statusVariants[row.original.status] || 'outline'}>
          {row.original.status}
        </Badge>
      ),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) =>
        row.original.status === 'pending' ? (
          <div className="flex justify-end gap-1">
            <Button size="icon" variant="ghost" className="h-7 w-7 text-green-600" onClick={() => approveAdv.mutate(row.original.id)}>
              <Check className="h-3.5 w-3.5" />
            </Button>
            <Button size="icon" variant="ghost" className="h-7 w-7 text-red-600" onClick={() => rejectAdv.mutate(row.original.id)}>
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
        ) : null,
    },
  ]

  if (error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Advance Management" description="Manage staff salary advances and recovery" />
        <ErrorState title="Failed to load advances" onRetry={() => refetch()} />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Advance Management"
        description="Manage staff salary advances and recovery"
        actions={
          <Button onClick={() => setFormOpen(true)} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            New Advance
          </Button>
        }
      />

      {/* Filter */}
      <div className="flex items-center gap-3">
        <Select value={statusFilter || 'all'} onValueChange={(v) => setStatusFilter(v === 'all' ? '' : v)}>
          <SelectTrigger className="w-[140px] h-9">
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="pending">Pending</SelectItem>
            <SelectItem value="approved">Approved</SelectItem>
            <SelectItem value="recovering">Recovering</SelectItem>
            <SelectItem value="recovered">Recovered</SelectItem>
            <SelectItem value="rejected">Rejected</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Data Table */}
      {isLoading ? (
        <LoadingState variant="table" />
      ) : (
        <DataTable
          columns={columns}
          data={data?.advances || []}
          searchPlaceholder="Search staff..."
          emptyTitle="No advances found"
          emptyDescription="Create a new advance request to get started."
          emptyAction={{ label: 'New Advance', onClick: () => setFormOpen(true) }}
        />
      )}

      {/* Add Advance Dialog */}
      <AddAdvanceDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        onSubmit={handleCreate}
        isLoading={createAdv.isPending}
      />
    </div>
  )
}

function AddAdvanceDialog({
  open,
  onOpenChange,
  onSubmit,
  isLoading,
}: {
  open: boolean
  onOpenChange: (v: boolean) => void
  onSubmit: (input: CreateAdvanceInput) => void
  isLoading: boolean
}) {
  const [form, setForm] = useState<CreateAdvanceInput>({
    staff_id: '',
    amount: 0,
    advance_date: new Date().toISOString().split('T')[0] ?? '',
    reason: '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(form)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle>New Advance Request</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Staff ID *</label>
              <Input
                value={form.staff_id}
                onChange={(e) => setForm({ ...form, staff_id: e.target.value })}
                placeholder="Staff UUID"
                required
              />
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
                value={form.advance_date}
                onChange={(e) => setForm({ ...form, advance_date: e.target.value })}
                required
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Reason</label>
              <Input
                value={form.reason}
                onChange={(e) => setForm({ ...form, reason: e.target.value })}
                placeholder="Personal, Medical..."
              />
            </div>
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={isLoading || !form.staff_id || !form.amount}>
              {isLoading ? 'Creating...' : 'Create Advance'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
