import { useState } from 'react'
import { useAdvanceList, useCreateAdvance, useApproveAdvance, useRejectAdvance } from '@/hooks/useSalary'
import { Plus, Check, X } from 'lucide-react'
import type { CreateAdvanceInput } from '@/types'

export function AdvancesPage() {
  const [showForm, setShowForm] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const { data, isLoading } = useAdvanceList({ status: statusFilter || undefined })
  const createAdv = useCreateAdvance()
  const approveAdv = useApproveAdvance()
  const rejectAdv = useRejectAdvance()

  const handleCreate = (input: CreateAdvanceInput) => {
    createAdv.mutate(input, { onSuccess: () => setShowForm(false) })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Advance Management</h1>
          <p className="text-muted-foreground">Manage staff salary advances and recovery</p>
        </div>
        <div className="flex gap-2">
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="rounded-md border px-3 py-2 text-sm"
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="recovering">Recovering</option>
            <option value="recovered">Recovered</option>
            <option value="rejected">Rejected</option>
          </select>
          <button
            onClick={() => setShowForm(!showForm)}
            className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            <Plus className="h-4 w-4" />
            New Advance
          </button>
        </div>
      </div>

      {showForm && <AdvanceForm onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={createAdv.isPending} />}

      <div className="rounded-lg border bg-card">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Staff</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Amount</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Date</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Reason</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Recovered</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Remaining</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Status</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Actions</th>
              </tr>
            </thead>
            <tbody>
              {isLoading && (
                <tr><td colSpan={8} className="px-6 py-8 text-center text-muted-foreground">Loading...</td></tr>
              )}
              {!isLoading && (!data?.advances || data.advances.length === 0) && (
                <tr><td colSpan={8} className="px-6 py-8 text-center text-muted-foreground">No advances found</td></tr>
              )}
              {data?.advances?.map((adv) => (
                <tr key={adv.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm font-medium">{adv.staff_name}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{adv.amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm">{adv.advance_date}</td>
                  <td className="px-6 py-4 text-sm">{adv.reason}</td>
                  <td className="px-6 py-4 text-sm text-right text-green-600">₹{adv.recovered_amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right text-orange-600">₹{adv.remaining_amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={adv.status} />
                  </td>
                  <td className="px-6 py-4 text-right">
                    {adv.status === 'pending' && (
                      <div className="flex justify-end gap-1">
                        <button
                          onClick={() => approveAdv.mutate(adv.id)}
                          className="rounded bg-green-600 p-1 text-white hover:bg-green-700"
                          title="Approve"
                        >
                          <Check className="h-3 w-3" />
                        </button>
                        <button
                          onClick={() => rejectAdv.mutate(adv.id)}
                          className="rounded bg-red-600 p-1 text-white hover:bg-red-700"
                          title="Reject"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-blue-100 text-blue-800',
    recovering: 'bg-orange-100 text-orange-800',
    recovered: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
  }
  return (
    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${colors[status] || 'bg-gray-100 text-gray-800'}`}>
      {status}
    </span>
  )
}

function AdvanceForm({ onSubmit, onCancel, isLoading }: { onSubmit: (input: CreateAdvanceInput) => void; onCancel: () => void; isLoading: boolean }) {
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
    <form onSubmit={handleSubmit} className="rounded-lg border bg-card p-6 space-y-4">
      <h3 className="text-lg font-semibold">New Advance Request</h3>
      <div className="grid gap-4 md:grid-cols-2">
        <div>
          <label className="text-sm font-medium">Staff ID</label>
          <input
            type="text"
            value={form.staff_id}
            onChange={(e) => setForm({ ...form, staff_id: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            placeholder="Staff UUID"
            required
          />
        </div>
        <div>
          <label className="text-sm font-medium">Amount (₹)</label>
          <input
            type="number"
            value={form.amount}
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
            value={form.advance_date}
            onChange={(e) => setForm({ ...form, advance_date: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            required
          />
        </div>
        <div>
          <label className="text-sm font-medium">Reason</label>
          <input
            type="text"
            value={form.reason}
            onChange={(e) => setForm({ ...form, reason: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            placeholder="Personal, Medical, etc."
          />
        </div>
      </div>
      <div className="flex gap-2">
        <button type="submit" disabled={isLoading} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
          {isLoading ? 'Creating...' : 'Create Advance'}
        </button>
        <button type="button" onClick={onCancel} className="rounded-lg border px-4 py-2 text-sm font-medium">Cancel</button>
      </div>
    </form>
  )
}
