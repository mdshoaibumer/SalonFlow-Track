import { Sun } from 'lucide-react'
import { useLocation } from 'react-router-dom'
import { Button } from '@/components/ui/button'

const routeTitles: Record<string, string> = {
  '/': 'Dashboard',
  '/staff': 'Staff Management',
  '/services': 'Services',
  '/customers': 'Customers',
  '/billing': 'Billing',
  '/invoices': 'Invoices',
  '/performance': 'Performance',
  '/commissions': 'Commissions',
  '/salary': 'Salary Management',
  '/advances': 'Advances',
  '/expenses': 'Expenses',
  '/profit-loss': 'Profit & Loss',
  '/products': 'Products',
  '/purchases': 'Purchases',
  '/inventory': 'Inventory',
  '/analytics': 'Analytics',
  '/settings': 'Settings',
  '/backups': 'Backups',
  '/appointments': 'Appointments',
  '/memberships': 'Memberships',
}

export function Header() {
  const location = useLocation()

  const title = routeTitles[location.pathname] || 'SalonFlow Track'

  return (
    <header className="flex h-14 items-center justify-between border-b bg-background px-6">
      <div>
        <h2 className="text-sm font-semibold">{title}</h2>
      </div>

      <div className="flex items-center gap-2">
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8"
          aria-label="Light theme"
        >
          <Sun className="h-4 w-4" />
        </Button>
      </div>
    </header>
  )
}
