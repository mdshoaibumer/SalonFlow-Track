import { useNavigate } from 'react-router-dom'
import { usePerformanceStats, useRevenueTrend, useTopPerformers } from '@/hooks/usePerformance'
import { useStaffStats } from '@/hooks/useStaff'
import { useInvoiceStats } from '@/hooks/useInvoices'
import { PageHeader } from '@/components/shared/PageHeader'
import { KPICard } from '@/components/shared/KPICard'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
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
    const data = {
      date: new Date().toISOString().split('T')[0],
      revenue_today: perfStats?.total_revenue_today ?? 0,
      customers_today: perfStats?.total_customers_today ?? 0,
      avg_bill: perfStats?.avg_bill_today ?? 0,
    }
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `dashboard-${data.date}.json`
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <PageHeader
        title="Dashboard"
        description="Overview of your salon's performance"
        actions={
          <Button variant="outline" size="sm" onClick={handleExport}>
            <Download className="mr-2 h-4 w-4" />
            Export
          </Button>
        }
      />

      {/* KPI Cards */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        <KPICard
          title="Revenue Today"
          value={`₹${(perfStats?.total_revenue_today ?? 0).toLocaleString('en-IN')}`}
          icon={IndianRupee}
        />
        <KPICard
          title="Customers Today"
          value={perfStats?.total_customers_today ?? 0}
          icon={Users}
        />
        <KPICard
          title="Average Bill"
          value={`₹${Math.round(perfStats?.avg_bill_today ?? 0).toLocaleString('en-IN')}`}
          icon={TrendingUp}
        />
        <KPICard
          title="Top Performer"
          value={perfStats?.top_performer_month?.staff_name ?? '-'}
          icon={Trophy}
          description={perfStats?.top_performer_month ? `₹${perfStats.top_performer_month.revenue.toLocaleString('en-IN')} this month` : undefined}
        />
      </div>

      {/* Charts Row */}
      <div className="grid gap-4 lg:grid-cols-7">
        {/* Revenue Trend */}
        <Card className="lg:col-span-4">
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-semibold">Revenue Trend (14 Days)</CardTitle>
          </CardHeader>
          <CardContent>
            {revenueTrend && revenueTrend.length > 0 ? (
              <ResponsiveContainer width="100%" height={240}>
                <AreaChart data={revenueTrend} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                  <XAxis
                    dataKey="date"
                    tickFormatter={(v) => new Date(v).toLocaleDateString('en-IN', { day: 'numeric', month: 'short' })}
                    tick={{ fontSize: 11 }}
                  />
                  <YAxis
                    tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`}
                    tick={{ fontSize: 11 }}
                    width={50}
                  />
                  <Tooltip
                    formatter={(value) => [`₹${Number(value).toLocaleString('en-IN')}`, 'Revenue']}
                    labelFormatter={(label) => new Date(label).toLocaleDateString('en-IN', { weekday: 'short', day: 'numeric', month: 'short' })}
                  />
                  <Area
                    type="monotone"
                    dataKey="revenue"
                    stroke="hsl(var(--primary))"
                    fill="hsl(var(--primary))"
                    fillOpacity={0.1}
                    strokeWidth={2}
                  />
                </AreaChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-[240px] text-sm text-muted-foreground">
                No revenue data available yet
              </div>
            )}
          </CardContent>
        </Card>

        {/* Top Performers */}
        <Card className="lg:col-span-3">
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-semibold">Top Performers</CardTitle>
          </CardHeader>
          <CardContent>
            {topPerformers && topPerformers.length > 0 ? (
              <ResponsiveContainer width="100%" height={240}>
                <BarChart data={topPerformers.slice(0, 5)} layout="vertical" margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" className="stroke-muted" horizontal={false} />
                  <XAxis
                    type="number"
                    tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`}
                    tick={{ fontSize: 11 }}
                  />
                  <YAxis
                    dataKey="staff_name"
                    type="category"
                    width={80}
                    tick={{ fontSize: 11 }}
                  />
                  <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString('en-IN')}`, 'Revenue']} />
                  <Bar dataKey="revenue" fill="hsl(var(--primary))" radius={[0, 4, 4, 0]} />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-[240px] text-sm text-muted-foreground">
                No performance data available
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions + Stats Row */}
      <div className="grid gap-4 lg:grid-cols-3">
        {/* Quick Actions */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold">Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-2">
            <Button variant="outline" size="sm" className="h-auto py-3 flex-col gap-1" onClick={() => navigate('/billing')}>
              <Plus className="h-4 w-4" />
              <span className="text-xs">New Bill</span>
            </Button>
            <Button variant="outline" size="sm" className="h-auto py-3 flex-col gap-1" onClick={() => navigate('/customers')}>
              <UserPlus className="h-4 w-4" />
              <span className="text-xs">Add Customer</span>
            </Button>
            <Button variant="outline" size="sm" className="h-auto py-3 flex-col gap-1" onClick={() => navigate('/appointments')}>
              <Calendar className="h-4 w-4" />
              <span className="text-xs">Appointment</span>
            </Button>
            <Button variant="outline" size="sm" className="h-auto py-3 flex-col gap-1" onClick={() => navigate('/services')}>
              <Scissors className="h-4 w-4" />
              <span className="text-xs">Services</span>
            </Button>
          </CardContent>
        </Card>

        {/* Staff Summary */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold">Staff Summary</CardTitle>
          </CardHeader>
          <CardContent>
            {staffStats ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Total Staff</span>
                  <span className="text-sm font-medium">{staffStats.total}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Active</span>
                  <span className="text-sm font-medium text-green-600">{staffStats.active}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Inactive</span>
                  <span className="text-sm font-medium text-red-600">{staffStats.inactive}</span>
                </div>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Loading...</p>
            )}
          </CardContent>
        </Card>

        {/* Invoice Summary */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold">Today's Billing</CardTitle>
          </CardHeader>
          <CardContent>
            {invoiceStats ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Revenue</span>
                  <span className="text-sm font-medium">₹{invoiceStats.today_revenue.toLocaleString('en-IN')}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Invoices</span>
                  <span className="text-sm font-medium">{invoiceStats.today_invoices}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Avg Bill</span>
                  <span className="text-sm font-medium">₹{Math.round(invoiceStats.avg_bill_value).toLocaleString('en-IN')}</span>
                </div>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Loading...</p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
