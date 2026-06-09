import { useState } from 'react'
import { useCustomerReport } from '@/hooks/useAnalytics'
import { Users, UserPlus, Repeat, Cake, UserX } from 'lucide-react'
import {
  LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer,
} from 'recharts'

export function CustomerReportsPage() {
  const [dateFrom, setDateFrom] = useState(() => {
    const d = new Date(); d.setDate(d.getDate() - 30)
    return d.toISOString().slice(0, 10)
  })
  const [dateTo, setDateTo] = useState(() => new Date().toISOString().slice(0, 10))

  const { data, isLoading } = useCustomerReport({ date_from: dateFrom, date_to: dateTo })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Customer Reports</h1>
        <p className="text-muted-foreground">Customer analytics and insights</p>
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
          <div className="grid gap-4 md:grid-cols-5">
            <Card title="Total Customers" value={String(data.total_customers)} icon={<Users className="h-4 w-4" />} />
            <Card title="New Customers" value={String(data.new_customers)} icon={<UserPlus className="h-4 w-4 text-green-600" />} />
            <Card title="Repeat Customers" value={String(data.repeat_customers)} icon={<Repeat className="h-4 w-4 text-blue-600" />} />
            <Card title="Birthday Today" value={String(data.birthday_today)} icon={<Cake className="h-4 w-4 text-pink-500" />} />
            <Card title="Inactive" value={String(data.inactive_count)} icon={<UserX className="h-4 w-4 text-red-500" />} />
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Customer Growth Trend</h3>
              {data.growth_trend && data.growth_trend.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <LineChart data={data.growth_trend}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip />
                    <Line type="monotone" dataKey="value" stroke="#3b82f6" strokeWidth={2} name="New Customers" />
                  </LineChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Top Customers (by Revenue)</h3>
              {data.top_customers && data.top_customers.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <BarChart data={data.top_customers} layout="vertical">
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis type="number" tick={{ fontSize: 10 }} />
                    <YAxis type="category" dataKey="name" tick={{ fontSize: 10 }} width={100} />
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                    <Bar dataKey="value" fill="#10b981" />
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

function Card({ title, value, icon }: { title: string; value: string; icon: React.ReactNode }) {
  return (
    <div className="rounded-lg border bg-card p-4">
      <div className="flex items-center gap-2 text-muted-foreground mb-1">{icon}<span className="text-xs font-medium">{title}</span></div>
      <p className="text-xl font-bold">{value}</p>
    </div>
  )
}
