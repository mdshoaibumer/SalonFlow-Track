import { useState } from 'react'
import { useMonthlyPerformance, useDailyPerformance, useTopPerformers, useRevenueTrend, usePerformanceStats } from '@/hooks/usePerformance'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar } from 'recharts'
import { Trophy, TrendingUp, Users, IndianRupee } from 'lucide-react'

type ViewMode = 'daily' | 'monthly'

export function PerformancePage() {
  const [viewMode, setViewMode] = useState<ViewMode>('daily')
  const { data: stats } = usePerformanceStats()
  const { data: dailyData, isLoading: dailyLoading } = useDailyPerformance()
  const { data: monthlyData, isLoading: monthlyLoading } = useMonthlyPerformance()
  const { data: topPerformers } = useTopPerformers({ limit: 5 })
  const { data: revenueTrend } = useRevenueTrend()

  const performanceData = viewMode === 'daily' ? dailyData : monthlyData
  const isLoading = viewMode === 'daily' ? dailyLoading : monthlyLoading

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Staff Performance</h1>
          <p className="text-muted-foreground">Track revenue, customers, and commissions</p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setViewMode('daily')}
            className={`rounded-lg px-4 py-2 text-sm font-medium ${viewMode === 'daily' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'}`}
          >
            Daily
          </button>
          <button
            onClick={() => setViewMode('monthly')}
            className={`rounded-lg px-4 py-2 text-sm font-medium ${viewMode === 'monthly' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'}`}
          >
            Monthly
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <StatCard
          title="Revenue Today"
          value={`₹${(stats?.total_revenue_today ?? 0).toLocaleString()}`}
          icon={<IndianRupee className="h-4 w-4" />}
        />
        <StatCard
          title="Customers Today"
          value={String(stats?.total_customers_today ?? 0)}
          icon={<Users className="h-4 w-4" />}
        />
        <StatCard
          title="Avg Bill Today"
          value={`₹${Math.round(stats?.avg_bill_today ?? 0).toLocaleString()}`}
          icon={<TrendingUp className="h-4 w-4" />}
        />
        <StatCard
          title="Top Performer (Month)"
          value={stats?.top_performer_month?.staff_name ?? '-'}
          subtitle={stats?.top_performer_month ? `₹${stats.top_performer_month.revenue.toLocaleString()}` : undefined}
          icon={<Trophy className="h-4 w-4" />}
        />
      </div>

      {/* Revenue Trend Chart */}
      {revenueTrend && revenueTrend.length > 0 && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Revenue Trend (Last 30 Days)</h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={revenueTrend}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" tickFormatter={(v: string) => v.slice(5)} />
              <YAxis />
              <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Revenue']} />
              <Line type="monotone" dataKey="revenue" stroke="hsl(var(--primary))" strokeWidth={2} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      )}

      {/* Performance Table */}
      <div className="rounded-lg border bg-card">
        <div className="p-6">
          <h3 className="text-lg font-semibold">
            {viewMode === 'daily' ? 'Today\'s Performance' : 'Monthly Performance'}
          </h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-t bg-muted/50">
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Rank</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Staff Name</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Revenue</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Customers</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Invoices</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Avg Bill</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Commission</th>
              </tr>
            </thead>
            <tbody>
              {isLoading && (
                <tr><td colSpan={7} className="px-6 py-8 text-center text-muted-foreground">Loading...</td></tr>
              )}
              {!isLoading && (!performanceData || performanceData.length === 0) && (
                <tr><td colSpan={7} className="px-6 py-8 text-center text-muted-foreground">No performance data</td></tr>
              )}
              {performanceData?.map((row) => (
                <tr key={row.staff_id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm font-medium">{row.rank}</td>
                  <td className="px-6 py-4 text-sm font-medium">{row.staff_name}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{row.revenue.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right">{row.customer_count}</td>
                  <td className="px-6 py-4 text-sm text-right">{row.invoice_count}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{Math.round(row.avg_bill).toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{Math.round(row.commission).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Leaderboard / Top Performers Bar Chart */}
      {topPerformers && topPerformers.length > 0 && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Top Performers This Month</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={topPerformers} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis type="number" />
              <YAxis dataKey="staff_name" type="category" width={120} />
              <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Revenue']} />
              <Bar dataKey="revenue" fill="hsl(var(--primary))" radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
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
