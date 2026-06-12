import { useNavigate } from 'react-router-dom'
import { motion } from 'motion/react'
import { usePerformanceStats, useRevenueTrend, useTopPerformers } from '@/hooks/usePerformance'
import { useStaffStats } from '@/hooks/useStaff'
import { useInvoiceStats } from '@/hooks/useInvoices'
import { PageHeader } from '@/components/shared/PageHeader'
import { KPICard } from '@/components/shared/KPICard'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { staggerContainer, staggerItem } from '@/lib/motion'
import {
  IndianRupee,
  Users,
  TrendingUp,
  Trophy,
  Plus,
  UserPlus,
  Scissors,
  Download,
  Calendar,
  ArrowRight,
} from 'lucide-react'
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
} from 'recharts'

export function DashboardPage() {
  const navigate = useNavigate()
  const { data: perfStats, isLoading: perfLoading, error: perfError, refetch: perfRefetch } = usePerformanceStats()
  const { data: staffStats } = useStaffStats()
  const { data: invoiceStats } = useInvoiceStats()
  const { data: revenueTrend } = useRevenueTrend({ limit: 14 })
  const { data: topPerformers } = useTopPerformers({ limit: 5 })

  if (perfLoading) {
    return (
      <div className="space-y-6">
        <PageHeader title="Dashboard" description="Overview of your salon's performance" />
        <LoadingState variant="page" />
      </div>
    )
  }

  if (perfError) {
    return (
      <div className="space-y-6">
        <PageHeader title="Dashboard" description="Overview of your salon's performance" />
        <ErrorState
          title="Unable to load dashboard"
          message="Failed to connect to the backend. Please ensure the server is running."
          onRetry={() => perfRefetch()}
        />
      </div>
    )
  }

  const handleExport = () => {
    const date = new Date().toISOString().split('T')[0]
    const rows = [
      ['Metric', 'Value'],
      ['Date', date],
      ['Revenue Today', `₹${(perfStats?.total_revenue_today ?? 0).toLocaleString('en-IN')}`],
      ['Customers Today', String(perfStats?.total_customers_today ?? 0)],
      ['Average Bill', `₹${Math.round(perfStats?.avg_bill_today ?? 0).toLocaleString('en-IN')}`],
      ['Top Performer', perfStats?.top_performer_month?.staff_name ?? '-'],
    ]
    if (revenueTrend && revenueTrend.length > 0) {
      rows.push(['', ''])
      rows.push(['Date', 'Revenue'])
      for (const entry of revenueTrend) {
        rows.push([entry.date, `₹${Number(entry.revenue).toLocaleString('en-IN')}`])
      }
    }
    const csv = rows.map(r => r.join(',')).join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `dashboard-${date}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <motion.div
      className="space-y-6"
      variants={staggerContainer}
      initial="hidden"
      animate="visible"
    >
      {/* Header */}
      <motion.div variants={staggerItem}>
        <PageHeader
          title="Dashboard"
          description="Overview of your salon's performance"
          actions={
            <Button variant="outline" size="sm" className="rounded-lg" onClick={handleExport}>
              <Download className="mr-2 h-3.5 w-3.5" />
              Export
            </Button>
          }
        />
      </motion.div>

      {/* KPI Cards */}
      <motion.div
        className="grid gap-4 grid-cols-2 lg:grid-cols-4"
        variants={staggerContainer}
        initial="hidden"
        animate="visible"
      >
        <KPICard
          title="Revenue Today"
          value={`₹${(perfStats?.total_revenue_today ?? 0).toLocaleString('en-IN')}`}
          icon={IndianRupee}
          accent="emerald"
        />
        <KPICard
          title="Customers Today"
          value={perfStats?.total_customers_today ?? 0}
          icon={Users}
          accent="blue"
        />
        <KPICard
          title="Average Bill"
          value={`₹${Math.round(perfStats?.avg_bill_today ?? 0).toLocaleString('en-IN')}`}
          icon={TrendingUp}
          accent="amber"
        />
        <KPICard
          title="Top Performer"
          value={perfStats?.top_performer_month?.staff_name ?? '-'}
          icon={Trophy}
          accent="rose"
          description={perfStats?.top_performer_month ? `₹${perfStats.top_performer_month.revenue.toLocaleString('en-IN')} this month` : undefined}
        />
      </motion.div>

      {/* Charts Row */}
      <motion.div variants={staggerItem} className="grid gap-4 lg:grid-cols-7">
        {/* Revenue Trend */}
        <Card className="lg:col-span-4 overflow-hidden">
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-semibold">Revenue Trend (14 Days)</CardTitle>
          </CardHeader>
          <CardContent>
            {revenueTrend && revenueTrend.length > 0 ? (
              <motion.div
                initial={{ opacity: 0, scale: 0.97 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.4, delay: 0.3, ease: [0.43, 0.13, 0.23, 0.96] }}
              >
                <ResponsiveContainer width="100%" height={240}>
                  <AreaChart data={revenueTrend} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
                    <defs>
                      <linearGradient id="revenueGradient" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity={0.2} />
                        <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" className="stroke-border/40" vertical={false} />
                    <XAxis
                      dataKey="date"
                      tickFormatter={(v) => new Date(v).toLocaleDateString('en-IN', { day: 'numeric', month: 'short' })}
                      tick={{ fontSize: 11 }}
                      axisLine={false}
                      tickLine={false}
                    />
                    <YAxis
                      tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`}
                      tick={{ fontSize: 11 }}
                      width={50}
                      axisLine={false}
                      tickLine={false}
                    />
                    <Tooltip
                      formatter={(value) => [`₹${Number(value).toLocaleString('en-IN')}`, 'Revenue']}
                      labelFormatter={(label) => new Date(label).toLocaleDateString('en-IN', { weekday: 'short', day: 'numeric', month: 'short' })}
                      contentStyle={{ borderRadius: '10px', border: '1px solid hsl(var(--border))', boxShadow: 'var(--shadow-md)' }}
                    />
                    <Area
                      type="monotone"
                      dataKey="revenue"
                      stroke="hsl(var(--primary))"
                      fill="url(#revenueGradient)"
                      strokeWidth={2.5}
                      animationDuration={1200}
                      animationEasing="ease-out"
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </motion.div>
            ) : (
              <div className="flex items-center justify-center h-[240px] text-sm text-muted-foreground">
                No revenue data available yet
              </div>
            )}
          </CardContent>
        </Card>

        {/* Top Performers */}
        <Card className="lg:col-span-3 overflow-hidden">
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-semibold">Top Performers</CardTitle>
          </CardHeader>
          <CardContent>
            {topPerformers && topPerformers.length > 0 ? (
              <motion.div
                initial={{ opacity: 0, x: 12 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.4, delay: 0.4, ease: [0.43, 0.13, 0.23, 0.96] }}
              >
                <ResponsiveContainer width="100%" height={240}>
                  <BarChart data={topPerformers.slice(0, 5)} layout="vertical" margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
                    <defs>
                      <linearGradient id="barGradient" x1="0" y1="0" x2="1" y2="0">
                        <stop offset="0%" stopColor="#8b5cf6" stopOpacity={0.9} />
                        <stop offset="100%" stopColor="#6366f1" stopOpacity={0.9} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" className="stroke-border/40" horizontal={false} />
                    <XAxis
                      type="number"
                      tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`}
                      tick={{ fontSize: 11 }}
                      axisLine={false}
                      tickLine={false}
                    />
                    <YAxis
                      dataKey="staff_name"
                      type="category"
                      width={80}
                      tick={{ fontSize: 11 }}
                      axisLine={false}
                      tickLine={false}
                    />
                    <Tooltip
                      formatter={(value) => [`₹${Number(value).toLocaleString('en-IN')}`, 'Revenue']}
                      contentStyle={{ borderRadius: '10px', border: '1px solid hsl(var(--border))', boxShadow: 'var(--shadow-md)' }}
                    />
                    <Bar
                      dataKey="revenue"
                      fill="url(#barGradient)"
                      radius={[0, 6, 6, 0]}
                      animationDuration={1000}
                      animationEasing="ease-out"
                    />
                  </BarChart>
                </ResponsiveContainer>
              </motion.div>
            ) : (
              <div className="flex items-center justify-center h-[240px] text-sm text-muted-foreground">
                No performance data available
              </div>
            )}
          </CardContent>
        </Card>
      </motion.div>

      {/* Quick Actions + Stats Row */}
      <motion.div variants={staggerItem} className="grid gap-4 lg:grid-cols-3">
        {/* Quick Actions */}
        <Card className="overflow-hidden">
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold">Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-2">
            {[
              { label: 'New Bill', icon: Plus, route: '/billing', color: 'group-hover:text-violet-600' },
              { label: 'Add Customer', icon: UserPlus, route: '/customers', color: 'group-hover:text-blue-600' },
              { label: 'Appointment', icon: Calendar, route: '/appointments', color: 'group-hover:text-emerald-600' },
              { label: 'Services', icon: Scissors, route: '/services', color: 'group-hover:text-amber-600' },
            ].map((action) => (
              <Button
                key={action.route}
                variant="outline"
                size="sm"
                className="group h-auto py-3 flex-col gap-1.5 rounded-xl border-border/60 hover:border-violet-200 dark:hover:border-violet-800 hover:bg-violet-50/50 dark:hover:bg-violet-500/5 transition-all duration-200"
                onClick={() => navigate(action.route)}
              >
                <action.icon className={`h-4 w-4 text-muted-foreground transition-colors ${action.color}`} />
                <span className="text-xs font-medium">{action.label}</span>
              </Button>
            ))}
          </CardContent>
        </Card>

        {/* Staff Summary */}
        <Card>
          <CardHeader className="pb-3 flex flex-row items-center justify-between">
            <CardTitle className="text-base font-semibold">Staff Summary</CardTitle>
            <Button variant="ghost" size="sm" className="h-7 px-2 text-xs text-muted-foreground" onClick={() => navigate('/staff')}>
              View <ArrowRight className="ml-1 h-3 w-3" />
            </Button>
          </CardHeader>
          <CardContent>
            {staffStats ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Total Staff</span>
                  <span className="text-sm font-semibold">{staffStats.total}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Active</span>
                  <span className="text-sm font-semibold text-emerald-600 dark:text-emerald-400">{staffStats.active}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Inactive</span>
                  <span className="text-sm font-semibold text-red-600 dark:text-red-400">{staffStats.inactive}</span>
                </div>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Loading...</p>
            )}
          </CardContent>
        </Card>

        {/* Invoice Summary */}
        <Card>
          <CardHeader className="pb-3 flex flex-row items-center justify-between">
            <CardTitle className="text-base font-semibold">Today's Billing</CardTitle>
            <Button variant="ghost" size="sm" className="h-7 px-2 text-xs text-muted-foreground" onClick={() => navigate('/invoices')}>
              View <ArrowRight className="ml-1 h-3 w-3" />
            </Button>
          </CardHeader>
          <CardContent>
            {invoiceStats ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Revenue</span>
                  <span className="text-sm font-semibold">₹{invoiceStats.today_revenue.toLocaleString('en-IN')}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Invoices</span>
                  <span className="text-sm font-semibold">{invoiceStats.today_invoices}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Avg Bill</span>
                  <span className="text-sm font-semibold">₹{Math.round(invoiceStats.avg_bill_value).toLocaleString('en-IN')}</span>
                </div>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Loading...</p>
            )}
          </CardContent>
        </Card>
      </motion.div>
    </motion.div>
  )
}
