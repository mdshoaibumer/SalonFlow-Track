import type { LucideIcon } from 'lucide-react'
import { cn } from '@/lib/utils'

interface KPICardProps {
  title: string
  value: string | number
  icon: LucideIcon
  description?: string
  trend?: { value: number; label: string }
  className?: string
}

export function KPICard({ title, value, icon: Icon, description, trend, className }: KPICardProps) {
  return (
    <div className={cn('rounded-lg border bg-card p-6', className)}>
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium text-muted-foreground">{title}</p>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </div>
      <div className="mt-2">
        <p className="text-2xl font-bold">{value}</p>
        {trend && (
          <p className={cn('text-xs mt-1', trend.value >= 0 ? 'text-green-600' : 'text-red-600')}>
            {trend.value >= 0 ? '+' : ''}{trend.value}% {trend.label}
          </p>
        )}
        {description && <p className="text-xs text-muted-foreground mt-1">{description}</p>}
      </div>
    </div>
  )
}
