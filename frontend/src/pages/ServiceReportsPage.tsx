import { useState } from 'react'
import { useServiceReport } from '@/hooks/useAnalytics'
import { Scissors, IndianRupee, Hash } from 'lucide-react'
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend,
} from 'recharts'

const PIE_COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4', '#ec4899', '#6b7280']

export function ServiceReportsPage() {
  const [dateFrom, setDateFrom] = useState(() => {
    const d = new Date(); d.setDate(d.getDate() - 30)
    return d.toISOString().slice(0, 10)
  })
  const [dateTo, setDateTo] = useState(() => new Date().toISOString().slice(0, 10))

  const { data, isLoading } = useServiceReport({ date_from: dateFrom, date_to: dateTo })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Service Reports</h1>
        <p className="text-muted-foreground">Service popularity and revenue analytics</p>
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
          <div className="grid gap-4 md:grid-cols-3">
            <div className="rounded-lg border bg-card p-4">
              <div className="flex items-center gap-2 text-muted-foreground mb-1"><Hash className="h-4 w-4" /><span className="text-xs font-medium">Total Bookings</span></div>
              <p className="text-xl font-bold">{data.total_bookings}</p>
            </div>
            <div className="rounded-lg border bg-card p-4">
              <div className="flex items-center gap-2 text-muted-foreground mb-1"><IndianRupee className="h-4 w-4" /><span className="text-xs font-medium">Avg Service Value</span></div>
              <p className="text-xl font-bold">₹{Math.round(data.avg_service_value).toLocaleString()}</p>
            </div>
            <div className="rounded-lg border bg-card p-4">
              <div className="flex items-center gap-2 text-muted-foreground mb-1"><Scissors className="h-4 w-4" /><span className="text-xs font-medium">Top Service</span></div>
              <p className="text-xl font-bold">{data.top_services?.[0]?.name || 'N/A'}</p>
            </div>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Top Services (by Bookings)</h3>
              {data.top_services && data.top_services.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <BarChart data={data.top_services} layout="vertical">
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis type="number" tick={{ fontSize: 10 }} />
                    <YAxis type="category" dataKey="name" tick={{ fontSize: 10 }} width={120} />
                    <Tooltip />
                    <Bar dataKey="value" fill="#3b82f6" name="Bookings" />
                  </BarChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Revenue by Service</h3>
              {data.revenue_by_service && data.revenue_by_service.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <PieChart>
                    <Pie data={data.revenue_by_service} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80}>
                      {data.revenue_by_service.map((_, idx) => <Cell key={idx} fill={PIE_COLORS[idx % PIE_COLORS.length]} />)}
                    </Pie>
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Least Used Services</h3>
              {data.least_used && data.least_used.length > 0 ? (
                <div className="space-y-2">
                  {data.least_used.map((s) => (
                    <div key={s.name} className="flex items-center justify-between border-b pb-2">
                      <span className="text-sm">{s.name}</span>
                      <span className="text-sm font-medium text-muted-foreground">{s.value} bookings</span>
                    </div>
                  ))}
                </div>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
