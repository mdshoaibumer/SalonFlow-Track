import { useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import {
  LayoutDashboard,
  Users,
  Scissors,
  UserCircle,
  Receipt,
  FileText,
  BarChart3,
  Coins,
  Wallet,
  Banknote,
  CreditCard,
  PieChart,
  Package,
  ShoppingCart,
  Warehouse,
  LineChart,
  HardDrive,
  Shield,
  ArrowUpCircle,
  FileUp,
  IndianRupee,
  Printer,
  Calendar,
  MessageSquare,
  Crown,
  Cloud,
  Settings,
  ChevronDown,
  type LucideIcon,
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface NavItem {
  name: string
  href: string
  icon: LucideIcon
}

interface NavGroup {
  label: string
  items: NavItem[]
}

const navGroups: NavGroup[] = [
  {
    label: 'Main',
    items: [
      { name: 'Dashboard', href: '/', icon: LayoutDashboard },
      { name: 'Billing', href: '/billing', icon: Receipt },
      { name: 'Appointments', href: '/appointments', icon: Calendar },
    ],
  },
  {
    label: 'Management',
    items: [
      { name: 'Staff', href: '/staff', icon: Users },
      { name: 'Services', href: '/services', icon: Scissors },
      { name: 'Customers', href: '/customers', icon: UserCircle },
      { name: 'Memberships', href: '/memberships', icon: Crown },
    ],
  },
  {
    label: 'Finance',
    items: [
      { name: 'Invoices', href: '/invoices', icon: FileText },
      { name: 'Commissions', href: '/commissions', icon: Coins },
      { name: 'Salary', href: '/salary', icon: Wallet },
      { name: 'Advances', href: '/advances', icon: Banknote },
      { name: 'Expenses', href: '/expenses', icon: CreditCard },
      { name: 'Profit & Loss', href: '/profit-loss', icon: PieChart },
      { name: 'GST & Tax', href: '/gst', icon: IndianRupee },
    ],
  },
  {
    label: 'Inventory',
    items: [
      { name: 'Products', href: '/products', icon: Package },
      { name: 'Purchases', href: '/purchases', icon: ShoppingCart },
      { name: 'Stock', href: '/inventory', icon: Warehouse },
    ],
  },
  {
    label: 'Reports',
    items: [
      { name: 'Analytics', href: '/analytics', icon: LineChart },
      { name: 'Performance', href: '/performance', icon: BarChart3 },
      { name: 'Revenue', href: '/reports/revenue', icon: IndianRupee },
      { name: 'Staff Reports', href: '/reports/staff', icon: Users },
      { name: 'Customer Reports', href: '/reports/customers', icon: UserCircle },
    ],
  },
  {
    label: 'System',
    items: [
      { name: 'Settings', href: '/settings', icon: Settings },
      { name: 'Backups', href: '/backups', icon: HardDrive },
      { name: 'Cloud Backup', href: '/cloud-backup', icon: Cloud },
      { name: 'Printer', href: '/printer', icon: Printer },
      { name: 'WhatsApp', href: '/whatsapp', icon: MessageSquare },
      { name: 'Import', href: '/import', icon: FileUp },
      { name: 'License', href: '/license', icon: Shield },
      { name: 'Updates', href: '/updates', icon: ArrowUpCircle },
    ],
  },
]

function NavGroupSection({ group }: { group: NavGroup }) {
  const location = useLocation()
  const isGroupActive = group.items.some((item) => {
    if (item.href === '/') return location.pathname === '/'
    return location.pathname.startsWith(item.href)
  })
  const [isOpen, setIsOpen] = useState(isGroupActive)

  return (
    <div>
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="flex w-full items-center justify-between px-3 py-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground hover:text-foreground"
      >
        {group.label}
        <ChevronDown className={cn('h-3 w-3 transition-transform', isOpen && 'rotate-180')} />
      </button>
      {isOpen && (
        <div className="mt-1 space-y-0.5">
          {group.items.map((item) => (
            <NavLink
              key={item.href}
              to={item.href}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                  isActive
                    ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                    : 'text-sidebar-foreground hover:bg-sidebar-accent/50 hover:text-sidebar-accent-foreground'
                )
              }
            >
              <item.icon className="h-4 w-4 shrink-0" />
              <span className="truncate">{item.name}</span>
            </NavLink>
          ))}
        </div>
      )}
    </div>
  )
}

export function Sidebar() {
  return (
    <aside className="flex w-60 flex-col border-r bg-sidebar">
      {/* Brand */}
      <div className="flex h-14 items-center gap-2 border-b px-5">
        <Scissors className="h-5 w-5 text-primary" />
        <span className="text-base font-bold tracking-tight">SalonFlow</span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto px-2 py-3 space-y-4">
        {navGroups.map((group) => (
          <NavGroupSection key={group.label} group={group} />
        ))}
      </nav>

      {/* Footer */}
      <div className="border-t px-4 py-3">
        <p className="text-[11px] text-muted-foreground">
          SalonFlow Track v0.1.0
        </p>
      </div>
    </aside>
  )
}
