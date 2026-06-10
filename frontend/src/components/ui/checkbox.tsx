import * as React from 'react'
import { cn } from '@/lib/utils'
import { Check } from 'lucide-react'

export interface CheckboxProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type'> {
  onCheckedChange?: (checked: boolean) => void
}

const Checkbox = React.forwardRef<HTMLInputElement, CheckboxProps>(
  ({ className, onCheckedChange, onChange, checked, ...props }, ref) => {
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      onChange?.(e)
      onCheckedChange?.(e.target.checked)
    }

    return (
      <label className={cn('relative inline-flex h-4 w-4 shrink-0 cursor-pointer', className)}>
        <input
          type="checkbox"
          ref={ref}
          checked={checked}
          onChange={handleChange}
          className="peer sr-only"
          {...props}
        />
        <span className="flex h-4 w-4 items-center justify-center rounded-sm border border-primary peer-focus-visible:ring-2 peer-focus-visible:ring-ring peer-focus-visible:ring-offset-2 peer-disabled:cursor-not-allowed peer-disabled:opacity-50 peer-checked:bg-primary peer-checked:text-primary-foreground">
          <Check className="h-3 w-3 opacity-0 peer-checked:opacity-100" />
        </span>
      </label>
    )
  }
)
Checkbox.displayName = 'Checkbox'

export { Checkbox }
