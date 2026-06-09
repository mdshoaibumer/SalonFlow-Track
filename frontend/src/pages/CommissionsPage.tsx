import { useState } from 'react'
import { useCommissionRules, useMonthlyCommission, useCommissionStats, useCreateCommissionRule, useDeleteCommissionRule } from '@/hooks/useCommissions'
import { XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar } from 'recharts'
import { Plus, Trash2, IndianRupee, Award, TrendingUp } from 'lucide-react'
import type { CreateRuleInput, CommissionRuleType, CommissionTargetType, CommissionCalcType } from '@/types'

export function CommissionsPage() {
  const [showForm, setShowForm] = useState(false)
  const { data: rulesData, isLoading: rulesLoading } = useCommissionRules()
  const { data: monthly } = useMonthlyCommission()
  const { data: stats } = useCommissionStats()
  const createRule = useCreateCommissionRule()
  const deleteRule = useDeleteCommissionRule()

  const handleCreateRule = (input: CreateRuleInput) => {
    createRule.mutate(input, { onSuccess: () => setShowForm(false) })
  }

  const handleDelete = (id: string) => {
    if (confirm('Delete this commission rule?')) {
      deleteRule.mutate(id)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Commissions</h1>
          <p className="text-muted-foreground">Commission rules and incentive tracking</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="h-4 w-4" />
          Add Rule
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <StatCard
          title="Commission This Month"
          value={`₹${Math.round(stats?.total_commission_this_month ?? 0).toLocaleString()}`}
          icon={<IndianRupee className="h-4 w-4" />}
        />
        <StatCard
          title="Top Earner"
          value={stats?.top_earner?.staff_name ?? '-'}
          subtitle={stats?.top_earner ? `₹${Math.round(stats.top_earner.commission).toLocaleString()}` : undefined}
          icon={<Award className="h-4 w-4" />}
        />
        <StatCard
          title="Avg Commission"
          value={`₹${Math.round(stats?.avg_commission ?? 0).toLocaleString()}`}
          icon={<TrendingUp className="h-4 w-4" />}
        />
      </div>

      {/* Create Rule Form */}
      {showForm && <RuleForm onSubmit={handleCreateRule} onCancel={() => setShowForm(false)} isLoading={createRule.isPending} />}

      {/* Monthly Commission Chart */}
      {monthly && monthly.length > 0 && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Monthly Commission by Staff</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={monthly} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis type="number" />
              <YAxis dataKey="staff_name" type="category" width={120} />
              <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Commission']} />
              <Bar dataKey="commission" fill="hsl(var(--primary))" radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}

      {/* Rules Table */}
      <div className="rounded-lg border bg-card">
        <div className="p-6">
          <h3 className="text-lg font-semibold">Commission Rules</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-t bg-muted/50">
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Rule Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Target</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Calculation</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Status</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Actions</th>
              </tr>
            </thead>
            <tbody>
              {rulesLoading && (
                <tr><td colSpan={6} className="px-6 py-8 text-center text-muted-foreground">Loading...</td></tr>
              )}
              {!rulesLoading && (!rulesData?.rules || rulesData.rules.length === 0) && (
                <tr><td colSpan={6} className="px-6 py-8 text-center text-muted-foreground">No commission rules configured</td></tr>
              )}
              {rulesData?.rules?.map((rule) => (
                <tr key={rule.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm font-medium">{rule.rule_name}</td>
                  <td className="px-6 py-4 text-sm capitalize">{rule.rule_type.replace('_', ' ')}</td>
                  <td className="px-6 py-4 text-sm capitalize">{rule.target_type}</td>
                  <td className="px-6 py-4 text-sm">
                    {rule.calculation_type === 'percentage' ? `${rule.calculation_value}%` : `₹${rule.calculation_value}`}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${rule.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                      {rule.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button onClick={() => handleDelete(rule.id)} className="text-destructive hover:text-destructive/80">
                      <Trash2 className="h-4 w-4" />
                    </button>
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

function RuleForm({ onSubmit, onCancel, isLoading }: { onSubmit: (input: CreateRuleInput) => void; onCancel: () => void; isLoading: boolean }) {
  const [form, setForm] = useState<CreateRuleInput>({
    rule_name: '',
    rule_type: 'revenue_based',
    target_type: 'global',
    target_id: '',
    calculation_type: 'percentage',
    calculation_value: 0,
    minimum_target: 0,
    maximum_target: 0,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="rounded-lg border bg-card p-6 space-y-4">
      <h3 className="text-lg font-semibold">New Commission Rule</h3>
      <div className="grid gap-4 md:grid-cols-2">
        <div>
          <label className="text-sm font-medium">Rule Name</label>
          <input
            type="text"
            value={form.rule_name}
            onChange={(e) => setForm({ ...form, rule_name: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            required
          />
        </div>
        <div>
          <label className="text-sm font-medium">Rule Type</label>
          <select
            value={form.rule_type}
            onChange={(e) => setForm({ ...form, rule_type: e.target.value as CommissionRuleType })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
          >
            <option value="revenue_based">Revenue Based</option>
            <option value="service_based">Service Based</option>
            <option value="fixed">Fixed</option>
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Target Type</label>
          <select
            value={form.target_type}
            onChange={(e) => setForm({ ...form, target_type: e.target.value as CommissionTargetType })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
          >
            <option value="global">Global</option>
            <option value="staff">Staff</option>
            <option value="service">Service</option>
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Target ID (optional)</label>
          <input
            type="text"
            value={form.target_id}
            onChange={(e) => setForm({ ...form, target_id: e.target.value })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            placeholder="Staff or Service ID"
          />
        </div>
        <div>
          <label className="text-sm font-medium">Calculation Type</label>
          <select
            value={form.calculation_type}
            onChange={(e) => setForm({ ...form, calculation_type: e.target.value as CommissionCalcType })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
          >
            <option value="percentage">Percentage</option>
            <option value="fixed_amount">Fixed Amount</option>
            <option value="tiered">Tiered</option>
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Value</label>
          <input
            type="number"
            value={form.calculation_value}
            onChange={(e) => setForm({ ...form, calculation_value: Number(e.target.value) })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            min={0}
            step="0.01"
          />
        </div>
        <div>
          <label className="text-sm font-medium">Minimum Target (₹)</label>
          <input
            type="number"
            value={form.minimum_target}
            onChange={(e) => setForm({ ...form, minimum_target: Number(e.target.value) })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            min={0}
          />
        </div>
        <div>
          <label className="text-sm font-medium">Maximum Target (₹)</label>
          <input
            type="number"
            value={form.maximum_target}
            onChange={(e) => setForm({ ...form, maximum_target: Number(e.target.value) })}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            min={0}
          />
        </div>
      </div>
      <div className="flex gap-2">
        <button
          type="submit"
          disabled={isLoading}
          className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
        >
          {isLoading ? 'Creating...' : 'Create Rule'}
        </button>
        <button type="button" onClick={onCancel} className="rounded-lg border px-4 py-2 text-sm font-medium">
          Cancel
        </button>
      </div>
    </form>
  )
}

function StatCard({ title, value, subtitle, icon }: { title: string; value: string; subtitle?: string; icon: React.ReactNode }) {
  return (
    <div className="rounded-lg border bg-card p-6">
      <div className="flex items-center gap-2 text-muted-foreground mb-2">
        {icon}
        <span className="text-sm font-medium">{title}</span>
      </div>
      <p className="text-2xl font-bold">{value}</p>
      {subtitle && <p className="text-sm text-muted-foreground mt-1">{subtitle}</p>}
    </div>
  )
}
