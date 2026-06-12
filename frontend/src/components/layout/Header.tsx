import { Sun, Moon, Monitor, Search, Command } from 'lucide-react'
import { useLocation } from 'react-router-dom'
import { motion } from 'motion/react'
import { Button } from '@/components/ui/button'
import { useTheme } from '@/app/providers/ThemeProvider'

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
  '/reports/revenue': 'Revenue Reports',
  '/reports/staff': 'Staff Reports',
  '/reports/customers': 'Customer Reports',
  '/reports/services': 'Service Reports',
  '/reports/expenses': 'Expense Reports',
  '/reports/inventory': 'Inventory Reports',
  '/reports/profit-loss': 'P&L Reports',
  '/gst': 'GST & Tax',
  '/printer': 'Printer',
  '/whatsapp': 'WhatsApp',
  '/import': 'Import',
  '/license': 'License',
  '/updates': 'Updates',
  '/cloud-backup': 'Cloud Backup',
}

function getBreadcrumbs(pathname: string): { label: string; path: string }[] {
  const crumbs: { label: string; path: string }[] = []
  if (pathname === '/') return [{ label: 'Dashboard', path: '/' }]
  
  const parts = pathname.split('/').filter(Boolean)
  let currentPath = ''
  for (const part of parts) {
    currentPath += `/${part}`
    const title = routeTitles[currentPath] || part.replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
    crumbs.push({ label: title, path: currentPath })
  }
  return crumbs
}

export function Header() {
  const location = useLocation()
  const { theme, setTheme } = useTheme()

  const breadcrumbs = getBreadcrumbs(location.pathname)

  const cycleTheme = () => {
    const next = theme === 'light' ? 'dark' : theme === 'dark' ? 'system' : 'light'
    setTheme(next)
  }

  const ThemeIcon = theme === 'dark' ? Moon : theme === 'system' ? Monitor : Sun

  return (
    <header className="flex h-12 items-center justify-between border-b border-border/60 bg-background/80 backdrop-blur-sm px-6">
      <div className="flex items-center gap-2">
        {/* Breadcrumbs */}
        <nav className="flex items-center gap-1 text-[13px]">
          {breadcrumbs.map((crumb, idx) => (
            <motion.span
              key={crumb.path}
              initial={{ opacity: 0, x: -4 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.15, delay: idx * 0.05 }}
              className="flex items-center gap-1"
            >
              {idx > 0 && <span className="text-muted-foreground/40">/</span>}
              <span className={idx === breadcrumbs.length - 1 ? 'font-semibold text-foreground' : 'text-muted-foreground'}>
                {crumb.label}
              </span>
            </motion.span>
          ))}
        </nav>
      </div>

      <div className="flex items-center gap-2">
        {/* Command palette hint */}
        <div className="hidden md:flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-muted/50 border border-border/40 text-muted-foreground/60 text-[11px] cursor-pointer hover:bg-muted/80 hover:text-muted-foreground transition-colors">
          <Search className="h-3 w-3" />
          <span>Search</span>
          <kbd className="ml-2 inline-flex items-center gap-0.5 rounded border border-border/60 bg-background px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground/70">
            <Command className="h-2.5 w-2.5" />K
          </kbd>
        </div>

        {/* Theme toggle */}
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7 rounded-md"
          onClick={cycleTheme}
          aria-label={`Theme: ${theme}`}
        >
          <motion.div
            key={theme}
            initial={{ rotate: -90, opacity: 0 }}
            animate={{ rotate: 0, opacity: 1 }}
            transition={{ duration: 0.2, ease: [0.2, 0, 0, 1] }}
          >
            <ThemeIcon className="h-3.5 w-3.5 text-muted-foreground" />
          </motion.div>
        </Button>
      </div>
    </header>
  )
}
