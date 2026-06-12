import { motion } from 'motion/react'
import type { LucideIcon } from 'lucide-react'
import { cn } from '@/lib/utils'
import { kpiVariants } from '@/lib/motion'
import { useCountUp } from '@/hooks/useCountUp'

interface KPICardProps {
  title: string
  value: string | number
  icon: LucideIcon
  description?: string
  trend?: { value: number; label: string }
  className?: string
  /** Accent color for the icon background */
  accent?: 'violet' | 'emerald' | 'amber' | 'blue' | 'rose'
}

const accentStyles = {
  violet: 'bg-violet-100 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400',
  emerald: 'bg-emerald-100 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400',
  amber: 'bg-amber-100 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400',
  blue: 'bg-blue-100 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400',
  rose: 'bg-rose-100 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400',
}

function AnimatedValue({ value }: { value: string | number }) {
  // Try to extract numeric value for animation
  const strValue = String(value)
  const numericMatch = strValue.match(/[\d,]+/)
  const numericPart = numericMatch ? parseInt(numericMatch[0].replace(/,/g, ''), 10) : null

  const animatedNum = useCountUp(numericPart ?? 0, { duration: 900, delay: 100 })

  if (numericPart !== null && numericPart > 0) {
    const prefix = strValue.substring(0, strValue.indexOf(numericMatch![0]))
    const suffix = strValue.substring(strValue.indexOf(numericMatch![0]) + numericMatch![0].length)
    return (
      <span className="tabular-nums">
        {prefix}{animatedNum.toLocaleString('en-IN')}{suffix}
        {/* Hidden span for test accessibility */}
        <span className="sr-only">{strValue}</span>
      </span>
    )
  }

  return <span className="tabular-nums">{strValue}</span>
}

export function KPICard({ title, value, icon: Icon, description, trend, className, accent = 'violet' }: KPICardProps) {
  return (
    <motion.div
      variants={kpiVariants}
      className={cn(
        'relative overflow-hidden surface-base p-5 group cursor-default',
        'hover:shadow-elevation-2 hover:-translate-y-0.5 transition-all duration-200',
        className
      )}
    >
      {/* Subtle gradient overlay on hover */}
      <div className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-300 bg-gradient-to-br from-violet-500/[0.02] to-indigo-500/[0.02] pointer-events-none" />
      
      <div className="relative flex items-center justify-between">
        <p className="text-[12.5px] font-medium text-muted-foreground">{title}</p>
        <div className={cn(
          'flex h-8 w-8 items-center justify-center rounded-lg transition-transform duration-200 group-hover:scale-110',
          accentStyles[accent]
        )}>
          <Icon className="h-4 w-4" />
        </div>
      </div>
      <div className="relative mt-2">
        <p className="text-2xl font-bold tracking-tight">
          <AnimatedValue value={value} />
        </p>
        {trend && (
          <motion.div
            initial={{ opacity: 0, y: 4 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.5, duration: 0.3 }}
          >
            <p className={cn(
              'text-[11px] font-semibold mt-1.5 inline-flex items-center gap-1 px-1.5 py-0.5 rounded-md',
              trend.value >= 0
                ? 'text-emerald-700 dark:text-emerald-400 bg-emerald-100/60 dark:bg-emerald-500/10'
                : 'text-red-700 dark:text-red-400 bg-red-100/60 dark:bg-red-500/10'
            )}>
              {trend.value >= 0 ? '↑' : '↓'} {Math.abs(trend.value)}% {trend.label}
            </p>
          </motion.div>
        )}
        {description && <p className="text-[11px] text-muted-foreground mt-1">{description}</p>}
      </div>
    </motion.div>
  )
}
