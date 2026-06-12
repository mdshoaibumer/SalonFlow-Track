import { useState, useCallback } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import { motion, AnimatePresence } from 'motion/react'
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
  ChevronRight,
  Sparkles,
  ScrollText,
  Activity,
  KeyRound,
  UserCog,
  LogOut,
  type LucideIcon,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { spring } from '@/lib/motion'
import { useAuth } from '@/app/providers/AuthProvider'

interface NavItem {
  name: string
  href: string
  icon: LucideIcon
  permission?: string
}

interface NavGroup {
  label: string
  id: string
  items: NavItem[]
  permission?: string
}

const navGroups: NavGroup[] = [
  {
    label: 'Main',
    id: 'main',
    items: [
      { name: 'Dashboard', href: '/', icon: LayoutDashboard, permission: 'dashboard.view' },
      { name: 'Billing', href: '/billing', icon: Receipt, permission: 'billing.view' },
      { name: 'Appointments', href: '/appointments', icon: Calendar, permission: 'appointments.view' },
    ],
  },
  {
    label: 'Management',
    id: 'management',
    items: [
      { name: 'Staff', href: '/staff', icon: Users, permission: 'staff.read' },
      { name: 'Services', href: '/services', icon: Scissors, permission: 'services.read' },
      { name: 'Customers', href: '/customers', icon: UserCircle, permission: 'customers.read' },
      { name: 'Memberships', href: '/memberships', icon: Crown, permission: 'memberships.view' },
    ],
  },
  {
    label: 'Finance',
    id: 'finance',
    permission: 'billing.view',
    items: [
      { name: 'Invoices', href: '/invoices', icon: FileText, permission: 'billing.view' },
      { name: 'Commissions', href: '/commissions', icon: Coins, permission: 'commissions.view' },
      { name: 'Salary', href: '/salary', icon: Wallet, permission: 'salary.view' },
      { name: 'Advances', href: '/advances', icon: Banknote, permission: 'advances.view' },
      { name: 'Expenses', href: '/expenses', icon: CreditCard, permission: 'expenses.view' },
      { name: 'Profit & Loss', href: '/profit-loss', icon: PieChart, permission: 'profitloss.view' },
      { name: 'GST & Tax', href: '/gst', icon: IndianRupee, permission: 'gst.view' },
    ],
  },
  {
    label: 'Inventory',
    id: 'inventory',
    permission: 'inventory.view',
    items: [
      { name: 'Products', href: '/products', icon: Package, permission: 'inventory.view' },
      { name: 'Purchases', href: '/purchases', icon: ShoppingCart, permission: 'inventory.purchase' },
      { name: 'Stock', href: '/inventory', icon: Warehouse, permission: 'inventory.view' },
    ],
  },
  {
    label: 'Reports',
    id: 'reports',
    permission: 'reports.view',
    items: [
      { name: 'Analytics', href: '/analytics', icon: LineChart, permission: 'analytics.view' },
      { name: 'Performance', href: '/performance', icon: BarChart3, permission: 'performance.view' },
      { name: 'Revenue', href: '/reports/revenue', icon: IndianRupee, permission: 'reports.view' },
      { name: 'Staff Reports', href: '/reports/staff', icon: Users, permission: 'reports.view' },
      { name: 'Customer Reports', href: '/reports/customers', icon: UserCircle, permission: 'reports.view' },
    ],
  },
  {
    label: 'Administration',
    id: 'admin',
    permission: 'users.view',
    items: [
      { name: 'Users', href: '/users', icon: UserCog, permission: 'users.view' },
      { name: 'Audit Logs', href: '/audit-logs', icon: ScrollText, permission: 'audit.view' },
      { name: 'Diagnostics', href: '/diagnostics', icon: Activity, permission: 'diagnostics.view' },
      { name: 'Change Password', href: '/change-password', icon: KeyRound },
    ],
  },
  {
    label: 'System',
    id: 'system',
    items: [
      { name: 'Settings', href: '/settings', icon: Settings, permission: 'settings.view' },
      { name: 'Backups', href: '/backups', icon: HardDrive, permission: 'backup.create' },
      { name: 'Cloud Backup', href: '/cloud-backup', icon: Cloud, permission: 'backup.create' },
      { name: 'Printer', href: '/printer', icon: Printer, permission: 'printer.use' },
      { name: 'WhatsApp', href: '/whatsapp', icon: MessageSquare, permission: 'whatsapp.view' },
      { name: 'Import', href: '/import', icon: FileUp, permission: 'import.execute' },
      { name: 'License', href: '/license', icon: Shield, permission: 'license.view' },
      { name: 'Updates', href: '/updates', icon: ArrowUpCircle, permission: 'updates.view' },
    ],
  },
]

function NavGroupSection({ group }: { group: NavGroup }) {
  const location = useLocation()
  const { hasPermission } = useAuth()

  // Filter items by permission
  const visibleItems = group.items.filter(item =>
    !item.permission || hasPermission(item.permission)
  )

  // Hide group if no visible items or group permission not met
  if (visibleItems.length === 0) return null
  if (group.permission && !hasPermission(group.permission)) return null

  const isGroupActive = visibleItems.some((item) => {
    if (item.href === '/') return location.pathname === '/'
    return location.pathname.startsWith(item.href)
  })
  const [isOpen, setIsOpen] = useState(isGroupActive)

  const toggleGroup = useCallback(() => setIsOpen(prev => !prev), [])

  return (
    <div className="space-y-0.5">
      <button
        type="button"
        onClick={toggleGroup}
        className="flex w-full items-center justify-between px-3 py-1.5 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground/70 hover:text-muted-foreground transition-colors duration-fast"
      >
        {group.label}
        <motion.div
          animate={{ rotate: isOpen ? 90 : 0 }}
          transition={{ duration: 0.15, ease: [0.2, 0, 0, 1] }}
        >
          <ChevronRight className="h-3 w-3" />
        </motion.div>
      </button>
      <AnimatePresence initial={false}>
        {isOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2, ease: [0.2, 0, 0, 1] }}
            style={{ overflow: 'hidden' }}
          >
            <div className="space-y-0.5 pb-1">
              {visibleItems.map((item) => (
                <NavLink
                  key={item.href}
                  to={item.href}
                  className={({ isActive }) =>
                    cn(
                      'group relative flex items-center gap-2.5 rounded-lg px-3 py-[7px] text-[13px] font-medium transition-all duration-fast',
                      isActive
                        ? 'bg-gradient-to-r from-violet-500/10 to-indigo-500/10 text-violet-700 dark:text-violet-300 shadow-sm shadow-violet-500/5'
                        : 'text-sidebar-foreground/80 hover:bg-sidebar-accent hover:text-sidebar-foreground'
                    )
                  }
                >
                  {({ isActive }) => (
                    <>
                      {/* Active indicator pill */}
                      {isActive && (
                        <motion.div
                          layoutId="sidebar-active-indicator"
                          className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-4 rounded-full bg-gradient-to-b from-violet-500 to-indigo-600"
                          transition={spring.snappy}
                        />
                      )}
                      <item.icon className={cn(
                        'h-4 w-4 shrink-0 transition-colors duration-fast',
                        isActive ? 'text-violet-600 dark:text-violet-400' : 'text-muted-foreground group-hover:text-foreground'
                      )} />
                      <span className="truncate">{item.name}</span>
                    </>
                  )}
                </NavLink>
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

export function Sidebar() {
  const { user, logout } = useAuth()

  return (
    <aside className="flex w-[240px] flex-col border-r border-sidebar-border bg-sidebar">
      {/* Brand */}
      <div className="flex h-14 items-center gap-2.5 px-5 border-b border-sidebar-border">
        <div className="relative flex h-8 w-8 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500 to-indigo-600 shadow-md shadow-violet-500/20">
          <Scissors className="h-4 w-4 text-white" />
          <div className="absolute -top-0.5 -right-0.5 h-2.5 w-2.5 rounded-full bg-emerald-400 border-2 border-sidebar animate-pulse" />
        </div>
        <div className="flex flex-col">
          <span className="text-[14px] font-bold tracking-tight bg-gradient-to-r from-violet-600 to-indigo-600 dark:from-violet-400 dark:to-indigo-400 bg-clip-text text-transparent">SalonFlow</span>
          <span className="text-[10px] text-muted-foreground/60 -mt-0.5 font-medium">Business Suite</span>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto scrollbar-hidden px-2.5 py-3 space-y-3">
        {navGroups.map((group) => (
          <NavGroupSection key={group.id} group={group} />
        ))}
      </nav>

      {/* Footer - User & Logout */}
      <div className="border-t border-sidebar-border px-4 py-3 space-y-2">
        {user && (
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 min-w-0">
              <div className="flex h-6 w-6 items-center justify-center rounded-full bg-violet-100 dark:bg-violet-900/30 text-[10px] font-bold text-violet-600">
                {user.display_name.charAt(0).toUpperCase()}
              </div>
              <div className="min-w-0">
                <p className="text-[11px] font-medium truncate">{user.display_name}</p>
                <p className="text-[9px] text-muted-foreground/60 truncate">{user.roles.join(', ')}</p>
              </div>
            </div>
            <button
              onClick={logout}
              className="rounded p-1.5 text-muted-foreground/60 hover:bg-red-100 hover:text-red-600 dark:hover:bg-red-900/20 transition-colors"
              title="Logout"
            >
              <LogOut className="h-3.5 w-3.5" />
            </button>
          </div>
        )}
        <div className="flex items-center gap-2">
          <Sparkles className="h-3 w-3 text-violet-500/60" />
          <p className="text-[10.5px] text-muted-foreground/50 font-medium">
            v0.2.0 — Desktop
          </p>
        </div>
      </div>
    </aside>
  )
}
