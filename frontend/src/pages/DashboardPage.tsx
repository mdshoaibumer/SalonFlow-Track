import { useQuery } from '@tanstack/react-query'
import { getHealthStatus } from '@/services/health'
import { StaffStatsWidget } from '@/components/staff/StaffStatsWidget'
import { ServiceStatsWidget } from '@/components/services/ServiceStatsWidget'
import { CustomerStatsWidget } from '@/components/customers/CustomerStatsWidget'
import { usePerformanceStats } from '@/hooks/usePerformance'
import { Trophy, TrendingUp, Users, IndianRupee } from 'lucide-react'

export function DashboardPage() {
  const { data: health, isLoading, error } = useQuery({
    queryKey: ['health'],
    queryFn: getHealthStatus,
    refetchInterval: 30000,
    retry: 5,
    retryDelay: 1000,
  })
  const { data: perfStats } = usePerformanceStats()

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome to SalonFlow Track
        </p>
      </div>

      {/* Performance Summary */}
      {perfStats && (
        <div className="grid gap-4 md:grid-cols-4">
          <div className="rounded-lg border bg-card p-4">
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <IndianRupee className="h-4 w-4" />
              <span className="text-xs font-medium">Revenue Today</span>
            </div>
            <p className="text-xl font-bold">₹{(perfStats.total_revenue_today ?? 0).toLocaleString()}</p>
          </div>
          <div className="rounded-lg border bg-card p-4">
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <Users className="h-4 w-4" />
              <span className="text-xs font-medium">Customers Today</span>
            </div>
            <p className="text-xl font-bold">{perfStats.total_customers_today ?? 0}</p>
          </div>
          <div className="rounded-lg border bg-card p-4">
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <TrendingUp className="h-4 w-4" />
              <span className="text-xs font-medium">Avg Bill</span>
            </div>
            <p className="text-xl font-bold">₹{Math.round(perfStats.avg_bill_today ?? 0).toLocaleString()}</p>
          </div>
          <div className="rounded-lg border bg-card p-4">
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <Trophy className="h-4 w-4" />
              <span className="text-xs font-medium">Top Staff (Month)</span>
            </div>
            <p className="text-xl font-bold">{perfStats.top_performer_month?.staff_name ?? '-'}</p>
          </div>
        </div>
      )}

      {/* Staff Stats */}
      <StaffStatsWidget />

      {/* Service Stats */}
      <ServiceStatsWidget />

      {/* Customer Stats */}
      <CustomerStatsWidget />

      {/* System Status Card */}
      <div className="rounded-lg border bg-card p-6">
        <h3 className="text-lg font-semibold mb-4">System Status</h3>
        {isLoading && (
          <p className="text-muted-foreground">Checking system status...</p>
        )}
        {error && (
          <p className="text-destructive">
            Unable to connect to backend. Please ensure the server is running.
          </p>
        )}
        {health && (
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            <StatusItem label="Status" value={health.status} />
            <StatusItem label="Version" value={health.version} />
            <StatusItem label="Database" value={health.database} />
            <StatusItem label="Uptime" value={health.uptime} />
          </div>
        )}
      </div>
    </div>
  )
}

function StatusItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-sm text-muted-foreground">{label}</p>
      <p className="text-sm font-medium capitalize">{value}</p>
    </div>
  )
}
