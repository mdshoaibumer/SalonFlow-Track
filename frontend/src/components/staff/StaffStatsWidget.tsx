import { Users, UserCheck, UserX } from 'lucide-react'
import { useStaffStats } from '@/hooks/useStaff'

export function StaffStatsWidget() {
  const { data: stats, isLoading } = useStaffStats()

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-3">
        {[1, 2, 3].map((i) => (
          <div key={i} className="rounded-lg border bg-card p-6 animate-pulse">
            <div className="h-4 w-20 bg-muted rounded mb-2" />
            <div className="h-8 w-10 bg-muted rounded" />
          </div>
        ))}
      </div>
    )
  }

  if (!stats) return null

  const items = [
    { label: 'Total Staff', value: stats.total, icon: Users, color: 'text-blue-600' },
    { label: 'Active', value: stats.active, icon: UserCheck, color: 'text-green-600' },
    { label: 'Inactive', value: stats.inactive, icon: UserX, color: 'text-gray-500' },
  ]

  return (
    <div className="grid gap-4 md:grid-cols-3">
      {items.map((item) => (
        <div key={item.label} className="rounded-lg border bg-card p-6">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-muted-foreground">{item.label}</p>
            <item.icon className={`h-5 w-5 ${item.color}`} />
          </div>
          <p className="mt-2 text-3xl font-bold">{item.value}</p>
        </div>
      ))}
    </div>
  )
}
