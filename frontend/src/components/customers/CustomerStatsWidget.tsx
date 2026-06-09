import { useCustomerStats } from '@/hooks/useCustomers'
import { Users, UserCheck, UserPlus, Cake } from 'lucide-react'

export function CustomerStatsWidget() {
  const { data: stats, isLoading } = useCustomerStats()

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
    { label: 'Total Customers', value: stats.total, icon: Users },
    { label: 'Active', value: stats.active, icon: UserCheck },
    { label: 'New This Month', value: stats.new_this_month, icon: UserPlus },
    { label: 'Birthday Today', value: stats.birthday_today, icon: Cake },
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
