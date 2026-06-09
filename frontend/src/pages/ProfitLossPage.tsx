import { useState } from 'react'
import { useProfitLoss, useMonthlyTrend, useExpenseStats } from '@/hooks/useExpense'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, LineChart, Line, Legend } from 'recharts'
import { TrendingUp, TrendingDown, IndianRupee, Percent } from 'lucide-react'

const COLORS = ['#ef4444', '#f97316', '#eab308', '#22c55e', '#06b6d4', '#3b82f6', '#8b5cf6', '#ec4899', '#6b7280', '#14b8a6']

export function ProfitLossPage() {
  const now = new Date()
  const [dateFrom, setDateFrom] = useState(() => {
    const d = new Date(now.getFullYear(), now.getMonth(), 1)
    return d.toISOString().split('T')[0] ?? ''
  })
  const [dateTo, setDateTo] = useState(() => now.toISOString().split('T')[0] ?? '')

  const { data: stats } = useExpenseStats()
  const { data: pl } = useProfitLoss(dateFrom, dateTo)
  const { data: trends } = useMonthlyTrend(6)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Profit & Loss</h1>
          <p className="text-muted-foreground">Business revenue, expenses, and profit analysis</p>
        </div>
        <div className="flex items-center gap-2">
          <input
            type="date"
            value={dateFrom}
            onChange={(e) => setDateFrom(e.target.value)}
            className="rounded-md border px-3 py-2 text-sm"
          />
          <span className="text-sm text-muted-foreground">to</span>
          <input
            type="date"
            value={dateTo}
            onChange={(e) => setDateTo(e.target.value)}
            className="rounded-md border px-3 py-2 text-sm"
          />
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <StatCard title="Monthly Revenue" value={`₹${Math.round(stats?.monthly_revenue ?? 0).toLocaleString()}`} icon={<IndianRupee className="h-4 w-4 text-green-500" />} />
        <StatCard title="Monthly Expenses" value={`₹${Math.round(stats?.monthly_expenses ?? 0).toLocaleString()}`} icon={<TrendingDown className="h-4 w-4 text-red-500" />} />
        <StatCard title="Net Profit" value={`₹${Math.round(stats?.monthly_profit ?? 0).toLocaleString()}`} icon={<TrendingUp className="h-4 w-4 text-blue-500" />} />
        <StatCard title="Profit Margin" value={`${(stats?.profit_margin ?? 0).toFixed(1)}%`} icon={<Percent className="h-4 w-4" />} />
      </div>

      {/* P&L Summary */}
      {pl && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">P&L Summary — {pl.period}</h3>
          <div className="grid gap-4 md:grid-cols-3">
            <div className="text-center p-4 rounded-lg bg-green-50">
              <p className="text-sm text-muted-foreground">Total Revenue</p>
              <p className="text-2xl font-bold text-green-700">₹{Math.round(pl.total_revenue).toLocaleString()}</p>
            </div>
            <div className="text-center p-4 rounded-lg bg-red-50">
              <p className="text-sm text-muted-foreground">Total Expenses</p>
              <p className="text-2xl font-bold text-red-700">₹{Math.round(pl.total_expenses).toLocaleString()}</p>
            </div>
            <div className="text-center p-4 rounded-lg bg-blue-50">
              <p className="text-sm text-muted-foreground">Gross Profit</p>
              <p className={`text-2xl font-bold ${pl.gross_profit >= 0 ? 'text-blue-700' : 'text-red-700'}`}>
                ₹{Math.round(pl.gross_profit).toLocaleString()}
              </p>
              <p className="text-xs text-muted-foreground mt-1">Margin: {pl.profit_margin.toFixed(1)}%</p>
            </div>
          </div>
        </div>
      )}

      {/* Charts Row */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Revenue vs Expense Trend */}
        {trends && trends.length > 0 && (
          <div className="rounded-lg border bg-card p-6">
            <h3 className="text-lg font-semibold mb-4">Revenue vs Expenses (6 Months)</h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={trends}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="month" />
                <YAxis />
                <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, '']} />
                <Legend />
                <Bar dataKey="revenue" fill="#22c55e" name="Revenue" />
                <Bar dataKey="expenses" fill="#ef4444" name="Expenses" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        )}

        {/* Expense Category Distribution */}
        {pl && pl.expenses_by_category && pl.expenses_by_category.length > 0 && (
          <div className="rounded-lg border bg-card p-6">
            <h3 className="text-lg font-semibold mb-4">Expense Distribution</h3>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={pl.expenses_by_category}
                  dataKey="amount"
                  nameKey="category_name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label={(entry) => `${entry.name} (${((entry.percent ?? 0) * 100).toFixed(0)}%)`}
                >
                  {pl.expenses_by_category.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, '']} />
              </PieChart>
            </ResponsiveContainer>
          </div>
        )}
      </div>

      {/* Profit Trend */}
      {trends && trends.length > 0 && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Monthly Profit Trend</h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={trends}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, '']} />
              <Legend />
              <Line type="monotone" dataKey="profit" stroke="#3b82f6" strokeWidth={2} name="Profit" />
              <Line type="monotone" dataKey="revenue" stroke="#22c55e" strokeWidth={1} strokeDasharray="5 5" name="Revenue" />
            </LineChart>
          </ResponsiveContainer>
        </div>
      )}

      {/* Expense Breakdown Table */}
      {pl && pl.expenses_by_category && pl.expenses_by_category.length > 0 && (
        <div className="rounded-lg border bg-card">
          <div className="p-6">
            <h3 className="text-lg font-semibold">Expense Breakdown by Category</h3>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-t bg-muted/50">
                  <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Category</th>
                  <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Amount</th>
                  <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">% of Total</th>
                </tr>
              </thead>
              <tbody>
                {pl.expenses_by_category.map((cat) => (
                  <tr key={cat.category_id} className="border-t hover:bg-muted/50">
                    <td className="px-6 py-4 text-sm font-medium">{cat.category_name}</td>
                    <td className="px-6 py-4 text-sm text-right">₹{Math.round(cat.amount).toLocaleString()}</td>
                    <td className="px-6 py-4 text-sm text-right">{cat.percentage.toFixed(1)}%</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}

function StatCard({ title, value, icon }: { title: string; value: string; icon: React.ReactNode }) {
  return (
    <div className="rounded-lg border bg-card p-6">
      <div className="flex items-center gap-2 text-muted-foreground mb-2">
        {icon}
        <span className="text-sm font-medium">{title}</span>
      </div>
      <p className="text-2xl font-bold">{value}</p>
    </div>
  )
}
