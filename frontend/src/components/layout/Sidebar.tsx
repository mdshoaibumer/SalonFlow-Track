import { NavLink } from 'react-router-dom'
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
} from 'lucide-react'
import { cn } from '@/lib/utils'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Staff', href: '/staff', icon: Users },
  { name: 'Services', href: '/services', icon: Scissors },
  { name: 'Customers', href: '/customers', icon: UserCircle },
  { name: 'Billing', href: '/billing', icon: Receipt },
  { name: 'Invoices', href: '/invoices', icon: FileText },
  { name: 'Performance', href: '/performance', icon: BarChart3 },
  { name: 'Commissions', href: '/commissions', icon: Coins },
  { name: 'Salary', href: '/salary', icon: Wallet },
  { name: 'Advances', href: '/advances', icon: Banknote },
  { name: 'Expenses', href: '/expenses', icon: CreditCard },
  { name: 'Profit & Loss', href: '/profit-loss', icon: PieChart },
  { name: 'Products', href: '/products', icon: Package },
  { name: 'Purchases', href: '/purchases', icon: ShoppingCart },
  { name: 'Inventory', href: '/inventory', icon: Warehouse },
  { name: 'Analytics', href: '/analytics', icon: LineChart },
  { name: 'Backups', href: '/backups', icon: HardDrive },
  { name: 'License', href: '/license', icon: Shield },
  { name: 'Updates', href: '/updates', icon: ArrowUpCircle },
  { name: 'Import', href: '/import', icon: FileUp },
  { name: 'GST & Tax', href: '/gst', icon: IndianRupee },
  { name: 'Printer', href: '/printer', icon: Printer },
  { name: 'Appointments', href: '/appointments', icon: Calendar },
  { name: 'WhatsApp', href: '/whatsapp', icon: MessageSquare },
  { name: 'Memberships', href: '/memberships', icon: Crown },
  { name: 'Cloud Backup', href: '/cloud-backup', icon: Cloud },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export function Sidebar() {
  return (
    <aside className="flex w-64 flex-col border-r bg-sidebar">
      {/* Brand */}
      <div className="flex h-16 items-center gap-2 border-b px-6">
        <Scissors className="h-6 w-6 text-primary" />
        <span className="text-lg font-semibold">SalonFlow</span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 px-3 py-4">
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                  : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground'
              )
            }
          >
            <item.icon className="h-4 w-4" />
            {item.name}
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div className="border-t p-4">
        <p className="text-xs text-muted-foreground">
          SalonFlow Track v0.1.0
        </p>
      </div>
    </aside>
  )
}
