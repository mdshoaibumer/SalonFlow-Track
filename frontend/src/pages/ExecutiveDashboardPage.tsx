import { useState } from 'react'
import { useDashboardStats, useKPIMetrics, useRevenueReport, useProfitLossReport } from '@/hooks/useAnalytics'
import {
  IndianRupee, Users, FileText, TrendingUp, TrendingDown,
  Package, Wallet, Banknote, AlertTriangle, BarChart3,
} from 'lucide-react'
import {
  LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend,
} from 'recharts'

const PIE_COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4', '#ec4899', '#6b7280']

type DatePreset = 'today' | '7d' | '30d' | 'month' | 'last_month' | 'custom'

function getDateRange(preset: DatePreset): { date_from: string; date_to: string } {
  const now = new Date()
  const fmt = (d: Date) => d.toISOString().slice(0, 10)
  switch (preset) {
    case 'today':
      return { date_from: fmt(now), date_to: fmt(now) }
    case '7d': {
      const from = new Date(now); from.setDate(from.getDate() - 7)
      return { date_from: fmt(from), date_to: fmt(now) }
    }
    case '30d': {
      const from = new Date(now); from.setDate(from.getDate() - 30)
      return { date_from: fmt(from), date_to: fmt(now) }
    }
    case 'month':
      return { date_from: `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`, date_to: fmt(now) }
    case 'last_month': {
      const lm = new Date(now.getFullYear(), now.getMonth() - 1, 1)
      const lmEnd = new Date(now.getFullYear(), now.getMonth(), 0)
      return { date_from: fmt(lm), date_to: fmt(lmEnd) }
    }
    default:
      return { date_from: `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`, date_to: fmt(now) }
  }
}

export function ExecutiveDashboardPage() {
  const [preset, setPreset] = useState<DatePreset>('month')
  const range = getDateRange(preset)

  const { data: stats } = useDashboardStats()
  const { data: kpis } = useKPIMetrics(range)
  const { data: revenue } = useRevenueReport({ ...range, group_by: 'day' })
  const { data: pl } = useProfitLossReport({ ...range, group_by: 'month' })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Business Analytics</h1>
          <p className="text-muted-foreground">Executive overview of business performance</p>
        </div>
        <DateFilter value={preset} onChange={setPreset} />
      </div>

      {/* Owner Dashboard - Primary KPIs */}
      <div className="grid gap-4 md:grid-cols-5">
        <StatCard title="Today's Revenue" value={`₹${(stats?.today_revenue ?? 0).toLocaleString()}`} icon={<IndianRupee className="h-4 w-4 text-green-600" />} />
        <StatCard title="Today's Customers" value={String(stats?.today_customers ?? 0)} icon={<Users className="h-4 w-4 text-blue-600" />} />
        <StatCard title="Today's Invoices" value={String(stats?.today_invoices ?? 0)} icon={<FileText className="h-4 w-4" />} />
        <StatCard title="Monthly Revenue" value={`₹${Math.round(stats?.monthly_revenue ?? 0).toLocaleString()}`} icon={<TrendingUp className="h-4 w-4 text-green-600" />} />
        <StatCard title="Monthly Profit" value={`₹${Math.round(stats?.monthly_profit ?? 0).toLocaleString()}`} icon={<BarChart3 className="h-4 w-4 text-emerald-600" />} highlight={stats?.monthly_profit !== undefined && stats.monthly_profit > 0} />
      </div>

      {/* Secondary KPIs */}
      <div className="grid gap-4 md:grid-cols-5">
        <StatCard title="Monthly Expenses" value={`₹${Math.round(stats?.monthly_expenses ?? 0).toLocaleString()}`} icon={<TrendingDown className="h-4 w-4 text-red-500" />} />
        <StatCard title="Inventory Value" value={`₹${Math.round(stats?.inventory_value ?? 0).toLocaleString()}`} icon={<Package className="h-4 w-4" />} />
        <StatCard title="Outstanding Salary" value={`₹${Math.round(stats?.outstanding_salary ?? 0).toLocaleString()}`} icon={<Wallet className="h-4 w-4 text-orange-500" />} />
        <StatCard title="Outstanding Advances" value={`₹${Math.round(stats?.outstanding_advances ?? 0).toLocaleString()}`} icon={<Banknote className="h-4 w-4 text-purple-500" />} />
        <StatCard title="Low Stock Alerts" value={String(stats?.low_stock_count ?? 0)} icon={<AlertTriangle className="h-4 w-4 text-red-500" />} />
      </div>

      {/* KPI Metrics */}
      {kpis && (
        <div className="grid gap-4 md:grid-cols-6">
          <KPICard label="Revenue Growth" value={kpis.revenue_growth_pct} suffix="%" />
          <KPICard label="Customer Growth" value={kpis.customer_growth_pct} suffix="%" />
          <KPICard label="Profit Margin" value={kpis.profit_margin_pct} suffix="%" />
          <KPICard label="Avg Bill Value" value={kpis.average_bill_value} prefix="₹" />
          <KPICard label="Repeat Customers" value={kpis.repeat_customer_pct} suffix="%" />
          <KPICard label="Staff Productivity" value={kpis.staff_productivity_pct} suffix="%" />
        </div>
      )}

      {/* Charts Row */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Revenue Trend */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Revenue Trend</h3>
          {revenue?.trend && revenue.trend.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={revenue.trend}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} />
                <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Revenue']} />
                <Line type="monotone" dataKey="revenue" stroke="#3b82f6" strokeWidth={2} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          ) : <p className="text-center text-muted-foreground py-8">No revenue data</p>}
        </div>

        {/* Revenue By Service Pie */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Revenue by Service</h3>
          {revenue?.by_service && revenue.by_service.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie data={revenue.by_service} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={90}
                  label={(entry) => `${entry.name}`}>
                  {revenue.by_service.map((_, idx) => <Cell key={idx} fill={PIE_COLORS[idx % PIE_COLORS.length]} />)}
                </Pie>
                <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Revenue']} />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          ) : <p className="text-center text-muted-foreground py-8">No data</p>}
        </div>

        {/* Profit & Loss Trend */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Profit & Loss Trend</h3>
          {pl?.trend && pl.trend.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={pl.trend}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} />
                <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                <Bar dataKey="revenue" fill="#3b82f6" name="Revenue" />
                <Bar dataKey="expenses" fill="#ef4444" name="Expenses" />
                <Bar dataKey="profit" fill="#10b981" name="Profit" />
                <Legend />
              </BarChart>
            </ResponsiveContainer>
          ) : <p className="text-center text-muted-foreground py-8">No data</p>}
        </div>

        {/* Top Staff */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Top Staff (by Revenue)</h3>
          {revenue?.by_staff && revenue.by_staff.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={revenue.by_staff} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis type="number" tick={{ fontSize: 11 }} />
                <YAxis type="category" dataKey="name" tick={{ fontSize: 11 }} width={100} />
                <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`, 'Revenue']} />
                <Bar dataKey="value" fill="#8b5cf6" />
              </BarChart>
            </ResponsiveContainer>
          ) : <p className="text-center text-muted-foreground py-8">No data</p>}
        </div>
      </div>
    </div>
  )
}

function StatCard({ title, value, icon, highlight }: { title: string; value: string; icon: React.ReactNode; highlight?: boolean }) {
  return (
    <div className={`rounded-lg border bg-card p-4 ${highlight ? 'ring-2 ring-green-200' : ''}`}>
      <div className="flex items-center gap-2 text-muted-foreground mb-1">
        {icon}
        <span className="text-xs font-medium">{title}</span>
      </div>
      <p className="text-xl font-bold">{value}</p>
    </div>
  )
}

function KPICard({ label, value, prefix, suffix }: { label: string; value: number; prefix?: string; suffix?: string }) {
  const formatted = `${prefix || ''}${Math.round(value * 10) / 10}${suffix || ''}`
  const isPositive = value > 0
  return (
    <div className="rounded-lg border bg-card p-4 text-center">
      <p className="text-xs text-muted-foreground mb-1">{label}</p>
      <p className={`text-lg font-bold ${isPositive ? 'text-green-600' : value < 0 ? 'text-red-600' : ''}`}>{formatted}</p>
    </div>
  )
}

function DateFilter({ value, onChange }: { value: DatePreset; onChange: (v: DatePreset) => void }) {
  return (
    <select value={value} onChange={(e) => onChange(e.target.value as DatePreset)} className="rounded-md border px-3 py-2 text-sm">
      <option value="today">Today</option>
      <option value="7d">Last 7 Days</option>
      <option value="30d">Last 30 Days</option>
      <option value="month">This Month</option>
      <option value="last_month">Last Month</option>
    </select>
  )
}
