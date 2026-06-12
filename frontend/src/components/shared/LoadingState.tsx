import { Skeleton } from '@/components/ui/skeleton'

interface LoadingStateProps {
  rows?: number
  variant?: 'table' | 'cards' | 'page' | 'chart'
}

export function LoadingState({ rows = 5, variant = 'table' }: LoadingStateProps) {
  if (variant === 'cards') {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="surface-base p-5 space-y-3" style={{ animationDelay: `${i * 50}ms` }}>
            <div className="flex items-center justify-between">
              <Skeleton className="h-3.5 w-20" />
              <Skeleton className="h-4 w-4 rounded" />
            </div>
            <Skeleton className="h-7 w-24" />
            <Skeleton className="h-3 w-14" />
          </div>
        ))}
      </div>
    )
  }

  if (variant === 'chart') {
    return (
      <div className="surface-base p-6 space-y-4">
        <div className="flex items-center justify-between">
          <Skeleton className="h-5 w-40" />
          <Skeleton className="h-8 w-24 rounded-md" />
        </div>
        <div className="flex items-end gap-2 h-[200px] pt-8">
          {Array.from({ length: 12 }).map((_, i) => (
            <Skeleton
              key={i}
              className="flex-1 rounded-t-sm"
              style={{ height: `${30 + Math.random() * 60}%`, animationDelay: `${i * 40}ms` }}
            />
          ))}
        </div>
      </div>
    )
  }

  if (variant === 'page') {
    return (
      <div className="space-y-6 animate-fade-in">
        {/* Header skeleton */}
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <Skeleton className="h-7 w-44" />
            <Skeleton className="h-4 w-64" />
          </div>
          <Skeleton className="h-9 w-28 rounded-lg" />
        </div>
        {/* KPI cards */}
        <div className="grid gap-4 md:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="surface-base p-5 space-y-3" style={{ animationDelay: `${i * 60}ms` }}>
              <div className="flex items-center justify-between">
                <Skeleton className="h-3.5 w-20" />
                <Skeleton className="h-4 w-4 rounded" />
              </div>
              <Skeleton className="h-7 w-24" />
            </div>
          ))}
        </div>
        {/* Table skeleton */}
        <div className="surface-base overflow-hidden">
          <div className="p-4 border-b border-border/60">
            <Skeleton className="h-9 w-72 rounded-lg" />
          </div>
          {Array.from({ length: rows }).map((_, i) => (
            <div
              key={i}
              className="flex items-center gap-4 px-4 py-3.5 border-b border-border/40 last:border-0"
              style={{ animationDelay: `${(i + 4) * 40}ms` }}
            >
              <Skeleton className="h-4 w-12" />
              <Skeleton className="h-4 flex-1 max-w-[200px]" />
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-4 w-20" />
              <Skeleton className="h-6 w-16 rounded-full" />
            </div>
          ))}
        </div>
      </div>
    )
  }

  // table variant
  return (
    <div className="surface-base overflow-hidden">
      <div className="p-4 border-b border-border/60">
        <Skeleton className="h-9 w-64 rounded-lg" />
      </div>
      {Array.from({ length: rows }).map((_, i) => (
        <div
          key={i}
          className="flex items-center gap-4 px-4 py-3.5 border-b border-border/40 last:border-0"
          style={{ animationDelay: `${i * 40}ms` }}
        >
          <Skeleton className="h-4 w-8" />
          <Skeleton className="h-4 flex-1 max-w-[180px]" />
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-4 w-20" />
          <Skeleton className="h-7 w-16 rounded-full" />
        </div>
      ))}
    </div>
  )
}
