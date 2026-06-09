import { useState } from 'react'
import { useStaffAnalyticsReport } from '@/hooks/useAnalytics'
import { Trophy, IndianRupee, Users, Coins } from 'lucide-react'
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer,
} from 'recharts'

export function StaffReportsPage() {
  const [dateFrom, setDateFrom] = useState(() => {
    const d = new Date(); d.setDate(d.getDate() - 30)
    return d.toISOString().slice(0, 10)
  })
  const [dateTo, setDateTo] = useState(() => new Date().toISOString().slice(0, 10))

  const { data, isLoading } = useStaffAnalyticsReport({ date_from: dateFrom, date_to: dateTo })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Staff Reports</h1>
        <p className="text-muted-foreground">Staff performance and productivity analysis</p>
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
          <div className="grid gap-4 md:grid-cols-2">
            <div className="rounded-lg border bg-card p-6 flex items-center gap-4">
              <IndianRupee className="h-8 w-8 text-orange-500" />
              <div>
                <p className="text-sm text-muted-foreground">Total Salary Cost</p>
                <p className="text-2xl font-bold">₹{data.salary_cost.toLocaleString()}</p>
              </div>
            </div>
            <div className="rounded-lg border bg-card p-6 flex items-center gap-4">
              <Trophy className="h-8 w-8 text-yellow-500" />
              <div>
                <p className="text-sm text-muted-foreground">Top Performer</p>
                <p className="text-2xl font-bold">{data.top_performers?.[0]?.name || 'N/A'}</p>
              </div>
            </div>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <ChartCard title="Revenue by Staff" icon={<IndianRupee className="h-4 w-4" />} data={data.revenue_by_staff} color="#3b82f6" />
            <ChartCard title="Customers by Staff" icon={<Users className="h-4 w-4" />} data={data.customers_by_staff} color="#10b981" />
            <ChartCard title="Commission Earned" icon={<Coins className="h-4 w-4" />} data={data.commission_earned} color="#f59e0b" />
          </div>
        </>
      )}
    </div>
  )
}

function ChartCard({ title, icon, data, color }: { title: string; icon: React.ReactNode; data?: { name: string; value: number }[]; color: string }) {
  return (
    <div className="rounded-lg border bg-card p-6">
      <div className="flex items-center gap-2 mb-4">
        {icon}
        <h3 className="text-lg font-semibold">{title}</h3>
      </div>
      {data && data.length > 0 ? (
        <ResponsiveContainer width="100%" height={250}>
          <BarChart data={data} layout="vertical">
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis type="number" tick={{ fontSize: 10 }} />
            <YAxis type="category" dataKey="name" tick={{ fontSize: 10 }} width={100} />
            <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
            <Bar dataKey="value" fill={color} />
          </BarChart>
        </ResponsiveContainer>
      ) : <p className="text-center text-muted-foreground py-8">No data</p>}
    </div>
  )
}
