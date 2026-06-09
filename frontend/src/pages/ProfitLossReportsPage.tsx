import { useState } from 'react'
import { useProfitLossReport } from '@/hooks/useAnalytics'
import { TrendingUp, TrendingDown, IndianRupee, Minus } from 'lucide-react'
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer, Legend, LineChart, Line,
} from 'recharts'

type GroupBy = 'day' | 'month' | 'year'

export function ProfitLossReportsPage() {
  const [dateFrom, setDateFrom] = useState(() => {
    const d = new Date(); d.setMonth(d.getMonth() - 6)
    return d.toISOString().slice(0, 10)
  })
  const [dateTo, setDateTo] = useState(() => new Date().toISOString().slice(0, 10))
  const [groupBy, setGroupBy] = useState<GroupBy>('month')

  const { data, isLoading } = useProfitLossReport({ date_from: dateFrom, date_to: dateTo, group_by: groupBy })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Profit & Loss Reports</h1>
        <p className="text-muted-foreground">Revenue - Expenses - Salary = Net Profit</p>
      </div>

      <div className="flex gap-3 flex-wrap items-end">
        <div>
          <label className="text-xs font-medium text-muted-foreground">From</label>
          <input type="date" value={dateFrom} onChange={(e) => setDateFrom(e.target.value)} className="block mt-1 rounded-md border px-3 py-2 text-sm" />
        </div>
        <div>
          <label className="text-xs font-medium text-muted-foreground">To</label>
          <input type="date" value={dateTo} onChange={(e) => setDateTo(e.target.value)} className="block mt-1 rounded-md border px-3 py-2 text-sm" />
        </div>
        <div>
          <label className="text-xs font-medium text-muted-foreground">Group By</label>
          <select value={groupBy} onChange={(e) => setGroupBy(e.target.value as GroupBy)} className="block mt-1 rounded-md border px-3 py-2 text-sm">
            <option value="day">Daily</option>
            <option value="month">Monthly</option>
            <option value="year">Yearly</option>
          </select>
        </div>
      </div>

      {isLoading && <p className="text-muted-foreground">Loading...</p>}

      {data && (
        <>
          {/* P&L Summary */}
          <div className="rounded-lg border bg-card p-6">
            <div className="grid gap-4 md:grid-cols-4">
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 text-green-600 mb-1">
                  <TrendingUp className="h-4 w-4" />
                  <span className="text-sm font-medium">Revenue</span>
                </div>
                <p className="text-2xl font-bold text-green-700">₹{Math.round(data.revenue).toLocaleString()}</p>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 text-red-600 mb-1">
                  <Minus className="h-4 w-4" />
                  <span className="text-sm font-medium">Expenses</span>
                </div>
                <p className="text-2xl font-bold text-red-700">₹{Math.round(data.expenses).toLocaleString()}</p>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 text-orange-600 mb-1">
                  <TrendingDown className="h-4 w-4" />
                  <span className="text-sm font-medium">Salary Cost</span>
                </div>
                <p className="text-2xl font-bold text-orange-700">₹{Math.round(data.salary_cost).toLocaleString()}</p>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 mb-1">
                  <IndianRupee className="h-4 w-4" />
                  <span className="text-sm font-medium">Net Profit</span>
                </div>
                <p className={`text-2xl font-bold ${data.net_profit >= 0 ? 'text-green-700' : 'text-red-700'}`}>
                  ₹{Math.round(data.net_profit).toLocaleString()}
                </p>
              </div>
            </div>
          </div>

          {/* Trend Charts */}
          <div className="grid gap-6 lg:grid-cols-2">
            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">P&L Breakdown</h3>
              {data.trend && data.trend.length > 0 ? (
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={data.trend}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                    <Bar dataKey="revenue" fill="#3b82f6" name="Revenue" />
                    <Bar dataKey="expenses" fill="#ef4444" name="Expenses" />
                    <Bar dataKey="salary" fill="#f59e0b" name="Salary" />
                    <Legend />
                  </BarChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Profit Trend</h3>
              {data.trend && data.trend.length > 0 ? (
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={data.trend}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Profit']} />
                    <Line type="monotone" dataKey="profit" stroke="#10b981" strokeWidth={2} name="Net Profit" />
                  </LineChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
