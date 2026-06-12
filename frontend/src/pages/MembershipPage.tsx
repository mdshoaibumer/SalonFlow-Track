import { useState } from 'react'
import { toastSuccess } from '@/lib/toast'
import { useMembershipPlans, useCreatePlan, useDeletePlan, useSellPlan, useUseSession, useSubscriptions, useMembershipStats } from '@/hooks/useMembership'
import type { MembershipPlan, MemberSubscription } from '@/types'

export function MembershipPage() {
  const [tab, setTab] = useState<'plans' | 'subscriptions' | 'sell'>('plans')

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">Memberships & Packages</h1>
        <p className="text-muted-foreground">Create plans, sell memberships, and track subscriptions</p>
      </div>

      <StatsCards />

      <div className="flex gap-2 border-b pb-2">
        <button onClick={() => setTab('plans')} className={`px-4 py-2 rounded-t ${tab === 'plans' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Plans</button>
        <button onClick={() => setTab('subscriptions')} className={`px-4 py-2 rounded-t ${tab === 'subscriptions' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Subscriptions</button>
        <button onClick={() => setTab('sell')} className={`px-4 py-2 rounded-t ${tab === 'sell' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Sell Plan</button>
      </div>

      {tab === 'plans' && <PlansTab />}
      {tab === 'subscriptions' && <SubscriptionsTab />}
      {tab === 'sell' && <SellPlanTab />}
    </div>
  )
}

function StatsCards() {
  const { data: stats } = useMembershipStats()
  if (!stats) return null
  return (
    <div className="grid grid-cols-4 gap-4">
      <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">{stats.active_subscriptions}</p><p className="text-xs text-muted-foreground">Active</p></div>
      <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">₹{stats.total_revenue}</p><p className="text-xs text-muted-foreground">Revenue</p></div>
      <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">{stats.expiring_soon}</p><p className="text-xs text-muted-foreground">Expiring Soon</p></div>
      <div className="border rounded p-3 text-center"><p className="text-lg font-bold">{stats.top_plan || '-'}</p><p className="text-xs text-muted-foreground">Top Plan</p></div>
    </div>
  )
}

function PlansTab() {
  const { data: plans, isLoading } = useMembershipPlans()
  const createMutation = useCreatePlan()
  const deleteMutation = useDeletePlan()
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ name: '', plan_type: 'package' as 'package' | 'membership', price: 0, validity_days: 30, total_sessions: 10, description: '' })

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate({ ...form, is_active: true }, { onSuccess: () => { setShowForm(false); toastSuccess('Membership created') } })
  }

  return (
    <div className="space-y-4">
      <button onClick={() => setShowForm(!showForm)} className="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">
        {showForm ? 'Cancel' : 'New Plan'}
      </button>

      {showForm && (
        <form onSubmit={handleCreate} className="border rounded p-4 space-y-3 max-w-lg">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Plan Name</label>
              <input className="w-full border rounded px-3 py-2 mt-1" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required />
            </div>
            <div>
              <label className="text-sm font-medium">Type</label>
              <select className="w-full border rounded px-3 py-2 mt-1" value={form.plan_type} onChange={e => setForm({ ...form, plan_type: e.target.value as 'package' | 'membership' })}>
                <option value="package">Package</option>
                <option value="membership">Membership</option>
              </select>
            </div>
            <div>
              <label className="text-sm font-medium">Price (₹)</label>
              <input type="number" className="w-full border rounded px-3 py-2 mt-1" value={form.price} onChange={e => setForm({ ...form, price: parseFloat(e.target.value) || 0 })} required />
            </div>
            <div>
              <label className="text-sm font-medium">Validity (days)</label>
              <input type="number" className="w-full border rounded px-3 py-2 mt-1" value={form.validity_days} onChange={e => setForm({ ...form, validity_days: parseInt(e.target.value) || 0 })} required />
            </div>
            <div>
              <label className="text-sm font-medium">Total Sessions</label>
              <input type="number" className="w-full border rounded px-3 py-2 mt-1" value={form.total_sessions} onChange={e => setForm({ ...form, total_sessions: parseInt(e.target.value) || 0 })} />
            </div>
          </div>
          <div>
            <label className="text-sm font-medium">Description</label>
            <textarea className="w-full border rounded px-3 py-2 mt-1" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} rows={2} />
          </div>
          <button type="submit" disabled={createMutation.isPending} className="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">Create Plan</button>
        </form>
      )}

      {isLoading && <p>Loading plans...</p>}

      <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
        {(plans || []).map((p: MembershipPlan) => (
          <div key={p.id} className="border rounded p-4">
            <div className="flex justify-between items-start">
              <div>
                <h3 className="font-medium">{p.name}</h3>
                <span className={`text-xs px-2 py-0.5 rounded ${p.plan_type === 'membership' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800'}`}>{p.plan_type}</span>
              </div>
              <button onClick={() => deleteMutation.mutate(p.id, { onSuccess: () => toastSuccess('Membership deleted') })} className="text-xs px-2 py-1 bg-red-100 text-red-800 rounded">Delete</button>
            </div>
            <p className="text-xl font-bold mt-2">₹{p.price}</p>
            <p className="text-xs text-muted-foreground">{p.validity_days} days &middot; {p.total_sessions} sessions</p>
            {p.description && <p className="text-sm mt-2 text-muted-foreground">{p.description}</p>}
          </div>
        ))}
      </div>
    </div>
  )
}

function SubscriptionsTab() {
  const { data: subs, isLoading } = useSubscriptions()
  const useSessionMutation = useUseSession()

  if (isLoading) return <p>Loading...</p>

  return (
    <div className="border rounded-lg overflow-hidden">
      <table className="w-full text-sm">
        <thead className="bg-muted">
          <tr>
            <th className="px-4 py-2 text-left">Customer</th>
            <th className="px-4 py-2 text-left">Plan</th>
            <th className="px-4 py-2 text-left">Sessions</th>
            <th className="px-4 py-2 text-left">Expires</th>
            <th className="px-4 py-2 text-left">Status</th>
            <th className="px-4 py-2 text-left">Actions</th>
          </tr>
        </thead>
        <tbody>
          {(subs?.data || []).map((s: MemberSubscription) => (
            <tr key={s.id} className="border-t">
              <td className="px-4 py-2">{s.customer_name || s.customer_id}</td>
              <td className="px-4 py-2">{s.plan_name || s.plan_id}</td>
              <td className="px-4 py-2">{s.used_sessions}/{s.total_sessions}</td>
              <td className="px-4 py-2">{new Date(s.end_date).toLocaleDateString()}</td>
              <td className="px-4 py-2"><span className={`px-2 py-0.5 rounded text-xs ${s.status === 'active' ? 'bg-green-100' : 'bg-gray-100'}`}>{s.status}</span></td>
              <td className="px-4 py-2">
                {s.status === 'active' && s.remaining_sessions > 0 && (
                  <button onClick={() => useSessionMutation.mutate(s.id)} className="text-xs px-2 py-1 bg-blue-100 rounded">Use Session</button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function SellPlanTab() {
  const { data: plans } = useMembershipPlans()
  const sellMutation = useSellPlan()
  const [form, setForm] = useState({ plan_id: '', customer_id: '', amount_paid: 0 })

  const handleSell = (e: React.FormEvent) => {
    e.preventDefault()
    sellMutation.mutate(form)
  }

  const selectedPlan = (plans || []).find(p => p.id === form.plan_id)

  return (
    <form onSubmit={handleSell} className="space-y-4 max-w-lg">
      <div>
        <label className="text-sm font-medium">Select Plan</label>
        <select className="w-full border rounded px-3 py-2 mt-1" value={form.plan_id} onChange={e => { const p = (plans || []).find(x => x.id === e.target.value); setForm({ ...form, plan_id: e.target.value, amount_paid: p?.price || 0 }) }} required>
          <option value="">Choose plan...</option>
          {(plans || []).filter(p => p.is_active).map(p => <option key={p.id} value={p.id}>{p.name} - ₹{p.price}</option>)}
        </select>
      </div>
      {selectedPlan && <p className="text-sm text-muted-foreground">{selectedPlan.plan_type} &middot; {selectedPlan.validity_days} days &middot; {selectedPlan.total_sessions} sessions</p>}
      <div>
        <label className="text-sm font-medium">Customer ID</label>
        <input className="w-full border rounded px-3 py-2 mt-1" value={form.customer_id} onChange={e => setForm({ ...form, customer_id: e.target.value })} required />
      </div>
      <div>
        <label className="text-sm font-medium">Amount Paid (₹)</label>
        <input type="number" className="w-full border rounded px-3 py-2 mt-1" value={form.amount_paid} onChange={e => setForm({ ...form, amount_paid: parseFloat(e.target.value) || 0 })} required />
      </div>
      <button type="submit" disabled={sellMutation.isPending} className="px-4 py-2 bg-green-600 text-white rounded">
        {sellMutation.isPending ? 'Processing...' : 'Sell Plan'}
      </button>
      {sellMutation.isSuccess && <p className="text-green-600 text-sm">Subscription created!</p>}
    </form>
  )
}
