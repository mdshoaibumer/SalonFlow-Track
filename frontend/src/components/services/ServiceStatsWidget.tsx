import { useServiceStats } from '@/hooks/useServices'
import { Package, PackageCheck, PackageMinus } from 'lucide-react'

export function ServiceStatsWidget() {
  const { data: stats, isLoading } = useServiceStats()

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="rounded-lg border bg-card p-4 animate-pulse">
            <div className="h-4 w-20 bg-muted rounded mb-2" />
            <div className="h-8 w-12 bg-muted rounded" />
          </div>
        ))}
      </div>
    )
  }

  if (!stats) return null

  const items = [
    { label: 'Total Services', value: stats.total, icon: Package },
    { label: 'Active', value: stats.active, icon: PackageCheck },
    { label: 'Inactive', value: stats.inactive, icon: PackageMinus },
    { label: 'Avg Price', value: `₹${Math.round(stats.avg_price).toLocaleString('en-IN')}`, icon: Package },
  ]

  return (
    <div className="grid gap-4 md:grid-cols-4">
      {items.map((item) => (
        <div key={item.label} className="rounded-lg border bg-card p-4">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <item.icon className="h-4 w-4" />
            {item.label}
          </div>
          <p className="text-2xl font-bold mt-1">{item.value}</p>
        </div>
      ))}
    </div>
  )
}
