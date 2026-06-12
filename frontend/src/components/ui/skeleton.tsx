import { cn } from '@/lib/utils'

function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'rounded-md bg-gradient-to-r from-muted via-muted/60 via-50% to-muted bg-[length:200%_100%] animate-shimmer',
        className
      )}
      {...props}
    />
  )
}

/** Circular skeleton for avatars */
function SkeletonCircle({ className, size = 40, ...props }: React.HTMLAttributes<HTMLDivElement> & { size?: number }) {
  return (
    <Skeleton
      className={cn('rounded-full', className)}
      style={{ width: size, height: size }}
      {...props}
    />
  )
}

/** Text-line skeleton with natural widths */
function SkeletonText({ className, lines = 3, ...props }: React.HTMLAttributes<HTMLDivElement> & { lines?: number }) {
  const widths = ['100%', '92%', '78%', '85%', '65%']
  return (
    <div className={cn('space-y-2', className)} {...props}>
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton
          key={i}
          className="h-3.5"
          style={{ width: widths[i % widths.length] }}
        />
      ))}
    </div>
  )
}

export { Skeleton, SkeletonCircle, SkeletonText }
