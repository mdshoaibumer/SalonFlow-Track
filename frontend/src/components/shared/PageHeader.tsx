import type { ReactNode } from 'react'
import { Button } from '@/components/ui/button'

interface PageHeaderProps {
  title: string
  description?: string
  actions?: ReactNode
}

export function PageHeader({ title, description, actions }: PageHeaderProps) {
  return (
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">{title}</h1>
        {description && <p className="text-sm text-muted-foreground mt-1">{description}</p>}
      </div>
      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </div>
  )
}

interface PageHeaderActionProps {
  label: string
  icon?: ReactNode
  onClick: () => void
  variant?: 'default' | 'outline' | 'secondary' | 'ghost' | 'destructive'
}

export function PageHeaderAction({ label, icon, onClick, variant = 'default' }: PageHeaderActionProps) {
  return (
    <Button variant={variant} onClick={onClick} size="sm">
      {icon}
      {label}
    </Button>
  )
}
