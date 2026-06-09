import { useState } from 'react'
import { useExpenseAnalyticsReport } from '@/hooks/useAnalytics'
import { CreditCard } from 'lucide-react'
import {
  LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend,
} from 'recharts'

const PIE_COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4', '#ec4899', '#6b7280']

export function ExpenseReportsPage() {
  const [dateFrom, setDateFrom] = useState(() => {
    const d = new Date(); d.setDate(d.getDate() - 30)
    return d.toISOString().slice(0, 10)
  })
  const [dateTo, setDateTo] = useState(() => new Date().toISOString().slice(0, 10))

  const { data, isLoading } = useExpenseAnalyticsReport({ date_from: dateFrom, date_to: dateTo })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Expense Reports</h1>
        <p className="text-muted-foreground">Expense breakdown and trends</p>
      </div>

      <div className="flex gap-3">
        <div>
          <label className="text-xs font-medium text-muted-foreground">From</label>
          <input type="date" value={dateFrom} onChange={(e) => setDateFrom(e.target.value)} className="block mt-1 rounded-md border px-3 py-2 text-sm" />
        </div>
        <div>
          <label className="text-xs font-medium text-muted-foreground">To</label>
          <input type="date" value={dateTo} onChange={(e) => setDateTo(e.target.value)} className="block mt-1 rounded-md border px-3 py-2 text-sm" />
        </div>
      </div>

      {isLoading && <p className="text-muted-foreground">Loading...</p>}

      {data && (
        <>
          <div className="rounded-lg border bg-card p-6 flex items-center gap-4">
            <CreditCard className="h-8 w-8 text-red-500" />
            <div>
              <p className="text-sm text-muted-foreground">Total Expenses</p>
              <p className="text-2xl font-bold">₹{data.total_expenses.toLocaleString()}</p>
            </div>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Expenses by Category</h3>
              {data.by_category && data.by_category.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <PieChart>
                    <Pie data={data.by_category} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80}>
                      {data.by_category.map((_, idx) => <Cell key={idx} fill={PIE_COLORS[idx % PIE_COLORS.length]} />)}
                    </Pie>
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Monthly Expense Trend</h3>
              {data.monthly_trend && data.monthly_trend.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <LineChart data={data.monthly_trend}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                    <Line type="monotone" dataKey="value" stroke="#ef4444" strokeWidth={2} name="Expenses" />
                  </LineChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6 lg:col-span-2">
              <h3 className="text-lg font-semibold mb-4">Revenue vs Expenses</h3>
              {data.revenue_vs_expense && data.revenue_vs_expense.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <BarChart data={data.revenue_vs_expense}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                    <Bar dataKey="value1" fill="#3b82f6" name="Revenue" />
                    <Bar dataKey="value2" fill="#ef4444" name="Expenses" />
                    <Legend />
                  </BarChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
