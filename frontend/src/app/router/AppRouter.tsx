import { Routes, Route, Navigate } from 'react-router-dom'
import { MainLayout } from '../layouts/MainLayout'
import { useAuth } from '@/app/providers/AuthProvider'
import { LoginPage } from '@/pages/LoginPage'
import { DashboardPage } from '@/pages/DashboardPage'
import { StaffPage } from '@/pages/StaffPage'
import { ServicesPage } from '@/pages/ServicesPage'
import { CustomersPage } from '@/pages/CustomersPage'
import { BillingPage } from '@/pages/BillingPage'
import { InvoicesPage } from '@/pages/InvoicesPage'
import { PerformancePage } from '@/pages/PerformancePage'
import { CommissionsPage } from '@/pages/CommissionsPage'
import { SalaryPage } from '@/pages/SalaryPage'
import { AdvancesPage } from '@/pages/AdvancesPage'
import { ExpensesPage } from '@/pages/ExpensesPage'
import { ProfitLossPage } from '@/pages/ProfitLossPage'
import { ProductsPage } from '@/pages/ProductsPage'
import { PurchasesPage } from '@/pages/PurchasesPage'
import { InventoryPage } from '@/pages/InventoryPage'
import { ExecutiveDashboardPage } from '@/pages/ExecutiveDashboardPage'
import { RevenueReportsPage } from '@/pages/RevenueReportsPage'
import { CustomerReportsPage } from '@/pages/CustomerReportsPage'
import { StaffReportsPage } from '@/pages/StaffReportsPage'
import { ServiceReportsPage } from '@/pages/ServiceReportsPage'
import { ExpenseReportsPage } from '@/pages/ExpenseReportsPage'
import { InventoryReportsPage } from '@/pages/InventoryReportsPage'
import { ProfitLossReportsPage } from '@/pages/ProfitLossReportsPage'
import { BackupPage } from '@/pages/BackupPage'
import { LicensePage } from '@/pages/LicensePage'
import { UpdatePage } from '@/pages/UpdatePage'
import { ImportPage } from '@/pages/ImportPage'
import { GSTPage } from '@/pages/GSTPage'
import { PrinterPage } from '@/pages/PrinterPage'
import { AppointmentsPage } from '@/pages/AppointmentsPage'
import { WhatsAppPage } from '@/pages/WhatsAppPage'
import { MembershipPage } from '@/pages/MembershipPage'
import { CloudBackupPage } from '@/pages/CloudBackupPage'
import { SettingsPage } from '@/pages/SettingsPage'
import { UserManagementPage } from '@/pages/UserManagementPage'
import { AuditLogPage } from '@/pages/AuditLogPage'
import { DiagnosticsPage } from '@/pages/DiagnosticsPage'
import { ChangePasswordPage } from '@/pages/ChangePasswordPage'

export function AppRouter() {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-background">
        <div className="animate-pulse text-sm text-muted-foreground">Loading...</div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return (
      <Routes>
        <Route path="*" element={<LoginPage />} />
      </Routes>
    )
  }

  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/staff" element={<StaffPage />} />
        <Route path="/services" element={<ServicesPage />} />
        <Route path="/customers" element={<CustomersPage />} />
        <Route path="/billing" element={<BillingPage />} />
        <Route path="/invoices" element={<InvoicesPage />} />
        <Route path="/performance" element={<PerformancePage />} />
        <Route path="/commissions" element={<CommissionsPage />} />
        <Route path="/salary" element={<SalaryPage />} />
        <Route path="/advances" element={<AdvancesPage />} />
        <Route path="/expenses" element={<ExpensesPage />} />
        <Route path="/profit-loss" element={<ProfitLossPage />} />
        <Route path="/products" element={<ProductsPage />} />
        <Route path="/purchases" element={<PurchasesPage />} />
        <Route path="/inventory" element={<InventoryPage />} />
        <Route path="/analytics" element={<ExecutiveDashboardPage />} />
        <Route path="/reports/revenue" element={<RevenueReportsPage />} />
        <Route path="/reports/customers" element={<CustomerReportsPage />} />
        <Route path="/reports/staff" element={<StaffReportsPage />} />
        <Route path="/reports/services" element={<ServiceReportsPage />} />
        <Route path="/reports/expenses" element={<ExpenseReportsPage />} />
        <Route path="/reports/inventory" element={<InventoryReportsPage />} />
        <Route path="/reports/profit-loss" element={<ProfitLossReportsPage />} />
        <Route path="/backups" element={<BackupPage />} />
        <Route path="/license" element={<LicensePage />} />
        <Route path="/updates" element={<UpdatePage />} />
        <Route path="/import" element={<ImportPage />} />
        <Route path="/gst" element={<GSTPage />} />
        <Route path="/printer" element={<PrinterPage />} />
        <Route path="/appointments" element={<AppointmentsPage />} />
        <Route path="/whatsapp" element={<WhatsAppPage />} />
        <Route path="/memberships" element={<MembershipPage />} />
        <Route path="/cloud-backup" element={<CloudBackupPage />} />
        <Route path="/settings" element={<SettingsPage />} />
        <Route path="/users" element={<UserManagementPage />} />
        <Route path="/audit-logs" element={<AuditLogPage />} />
        <Route path="/diagnostics" element={<DiagnosticsPage />} />
        <Route path="/change-password" element={<ChangePasswordPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
