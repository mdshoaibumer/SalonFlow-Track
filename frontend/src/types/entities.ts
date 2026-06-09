// Base entity type matching backend
export interface BaseEntity {
  id: string
  created_at: string
  updated_at: string
}

export interface User extends BaseEntity {
  name: string
  email: string
  phone: string
  role: string
  is_active: boolean
}

export type StaffDesignation = 'stylist' | 'assistant' | 'receptionist' | 'manager'
export type StaffStatus = 'active' | 'inactive'
export type StaffGender = 'male' | 'female' | 'other'

export interface Staff extends BaseEntity {
  staff_code: string
  full_name: string
  phone: string
  email: string
  gender: StaffGender
  designation: StaffDesignation
  joining_date: string
  base_salary: number
  commission_percentage: number
  status: StaffStatus
}

export interface StaffStats {
  total: number
  active: number
  inactive: number
}

export interface CreateStaffInput {
  full_name: string
  phone: string
  email?: string
  gender: StaffGender
  designation: StaffDesignation
  joining_date: string
  base_salary: number
  commission_percentage: number
}

export interface UpdateStaffInput extends CreateStaffInput {
  status: StaffStatus
}

// --- Services ---

export type ServiceCategory = 'hair' | 'facial' | 'skin' | 'spa' | 'massage' | 'coloring' | 'treatment' | 'other'
export type ServiceStatus = 'active' | 'inactive'
export type CommissionType = 'fixed' | 'percentage'

export interface Service extends BaseEntity {
  service_code: string
  name: string
  category: ServiceCategory
  description: string
  duration_minutes: number
  price: number
  cost_price: number
  commission_type: CommissionType
  commission_value: number
  status: ServiceStatus
}

export interface ServiceStats {
  total: number
  active: number
  inactive: number
  avg_price: number
}

export interface CreateServiceInput {
  name: string
  category: ServiceCategory
  description?: string
  duration_minutes: number
  price: number
  cost_price?: number
  commission_type: CommissionType
  commission_value: number
}

export interface UpdateServiceInput extends CreateServiceInput {
  status: ServiceStatus
}

// --- Customers ---

export type CustomerStatus = 'active' | 'inactive'

export interface Customer extends BaseEntity {
  customer_code: string
  full_name: string
  phone: string
  email: string
  gender: StaffGender
  date_of_birth?: string
  anniversary_date?: string
  address: string
  notes: string
  total_visits: number
  total_spent: number
  last_visit_date?: string
  status: CustomerStatus
}

export interface CustomerStats {
  total: number
  active: number
  inactive: number
  new_this_month: number
  birthday_today: number
}

export interface CreateCustomerInput {
  full_name: string
  phone: string
  email?: string
  gender?: StaffGender
  date_of_birth?: string
  anniversary_date?: string
  address?: string
  notes?: string
}

export interface UpdateCustomerInput extends CreateCustomerInput {
  status?: CustomerStatus
}

// --- Invoices ---

export type PaymentStatus = 'pending' | 'paid' | 'partial'
export type PaymentMethod = 'cash' | 'upi' | 'card' | 'bank_transfer'

export interface InvoiceItem {
  id: string
  invoice_id: string
  service_id: string
  service_name_snapshot: string
  quantity: number
  unit_price: number
  discount: number
  line_total: number
}

export interface Invoice extends BaseEntity {
  invoice_number: string
  customer_id: string
  staff_id: string
  items: InvoiceItem[]
  subtotal: number
  discount: number
  tax: number
  grand_total: number
  payment_status: PaymentStatus
  payment_method: string
  notes: string
  invoice_date: string
}

export interface Payment {
  id: string
  invoice_id: string
  amount: number
  payment_method: PaymentMethod
  reference_number: string
  payment_date: string
}

export interface InvoiceStats {
  today_revenue: number
  today_invoices: number
  avg_bill_value: number
}

export interface CreateInvoiceItemInput {
  service_id: string
  quantity: number
  discount: number
}

export interface CreateInvoiceInput {
  customer_id: string
  staff_id: string
  items: CreateInvoiceItemInput[]
  discount: number
  tax: number
  payment_method?: PaymentMethod
  notes?: string
}

export interface RecordPaymentInput {
  amount: number
  payment_method: PaymentMethod
  reference_number?: string
}

export interface Setting extends BaseEntity {
  key: string
  value: string
  description: string
  category: string
}

export interface License extends BaseEntity {
  license_key: string
  expiry_date: string
  status: 'active' | 'expired' | 'revoked'
  issued_to: string
  issued_at: string
}

// --- Staff Performance ---

export interface StaffPerformanceSummary {
  staff_id: string
  staff_name: string
  revenue: number
  customer_count: number
  invoice_count: number
  service_count: number
  avg_bill: number
  commission: number
  rank: number
}

export interface StaffPerformanceDaily {
  id: string
  staff_id: string
  business_date: string
  invoice_count: number
  customer_count: number
  service_count: number
  revenue: number
  commission_amount: number
}

export interface RevenueTrendPoint {
  date: string
  revenue: number
}

export interface PerformanceStats {
  top_performer_today: StaffPerformanceSummary | null
  top_performer_month: StaffPerformanceSummary | null
  total_revenue_today: number
  total_customers_today: number
  avg_bill_today: number
}

// --- Commission ---

export type CommissionRuleType = 'revenue_based' | 'service_based' | 'fixed'
export type CommissionTargetType = 'global' | 'staff' | 'service'
export type CommissionCalcType = 'percentage' | 'fixed_amount' | 'tiered'
export type CommissionTxStatus = 'pending' | 'approved' | 'paid'

export interface CommissionRule extends BaseEntity {
  rule_name: string
  rule_type: CommissionRuleType
  target_type: CommissionTargetType
  target_id: string
  calculation_type: CommissionCalcType
  calculation_value: number
  minimum_target: number
  maximum_target: number
  is_active: boolean
}

export interface CommissionTransaction {
  id: string
  staff_id: string
  invoice_id: string
  rule_id: string
  revenue_amount: number
  commission_amount: number
  business_date: string
  status: CommissionTxStatus
}

export interface CommissionStaffSummary {
  staff_id: string
  staff_name: string
  revenue: number
  commission: number
}

export interface CommissionStats {
  total_commission_this_month: number
  top_earner: CommissionStaffSummary | null
  avg_commission: number
}

export interface CreateRuleInput {
  rule_name: string
  rule_type: CommissionRuleType
  target_type: CommissionTargetType
  target_id?: string
  calculation_type: CommissionCalcType
  calculation_value: number
  minimum_target?: number
  maximum_target?: number
}

export interface UpdateRuleInput extends CreateRuleInput {
  is_active: boolean
}

export interface StaffCommissionOutput {
  staff_id: string
  total_revenue: number
  commission: number
  transactions: CommissionTransaction[]
}

// --- Salary & Advances ---

export type SalaryCycleStatus = 'draft' | 'generated' | 'finalized'
export type SalaryPaymentStatus = 'pending' | 'partial' | 'paid'
export type AdvanceStatus = 'pending' | 'approved' | 'recovering' | 'recovered' | 'rejected'

export interface SalaryCycle {
  id: string
  month: number
  year: number
  status: SalaryCycleStatus
  generated_at: string
  generated_by: string
  created_at: string
  updated_at: string
}

export interface SalaryRecord {
  id: string
  salary_cycle_id: string
  staff_id: string
  staff_name: string
  base_salary: number
  commission_amount: number
  bonus_amount: number
  advance_amount: number
  deduction_amount: number
  gross_salary: number
  net_salary: number
  payment_status: SalaryPaymentStatus
  payment_date: string
  notes: string
  created_at: string
  updated_at: string
}

export interface Advance {
  id: string
  staff_id: string
  staff_name: string
  amount: number
  advance_date: string
  reason: string
  recovered_amount: number
  remaining_amount: number
  status: AdvanceStatus
  created_at: string
  updated_at: string
}

export interface SalaryStats {
  total_payroll: number
  pending_payments: number
  paid_salaries: number
  outstanding_advances: number
}

export interface CreateAdvanceInput {
  staff_id: string
  amount: number
  advance_date: string
  reason: string
}

export interface GenerateSalaryInput {
  month: number
  year: number
}

export interface GenerateSalaryOutput {
  cycle: SalaryCycle
  records: SalaryRecord[]
}

// --- Expense Management ---

export type ExpensePaymentMethod = 'cash' | 'upi' | 'bank_transfer' | 'card' | 'cheque'
export type ExpenseStatus = 'pending' | 'approved' | 'paid' | 'rejected'

export interface ExpenseCategory {
  id: string
  name: string
  description: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Expense {
  id: string
  expense_number: string
  category_id: string
  category_name: string
  amount: number
  expense_date: string
  payment_method: ExpensePaymentMethod
  vendor_name: string
  invoice_reference: string
  description: string
  attachment_path: string
  status: ExpenseStatus
  created_by: string
  created_at: string
  updated_at: string
}

export interface CreateExpenseInput {
  category_id: string
  amount: number
  expense_date: string
  payment_method: ExpensePaymentMethod
  vendor_name: string
  invoice_reference?: string
  description?: string
}

export interface UpdateExpenseInput extends CreateExpenseInput {
  status: ExpenseStatus
}

export interface ExpenseStats {
  today_expenses: number
  monthly_expenses: number
  today_revenue: number
  monthly_revenue: number
  monthly_profit: number
  profit_margin: number
}

export interface CategoryExpense {
  category_id: string
  category_name: string
  amount: number
  percentage: number
}

export interface ProfitLoss {
  period: string
  total_revenue: number
  total_expenses: number
  gross_profit: number
  profit_margin: number
  expenses_by_category: CategoryExpense[]
}

export interface MonthlyTrend {
  month: string
  revenue: number
  expenses: number
  profit: number
}

export interface ExpenseReport {
  date_from: string
  date_to: string
  total_expenses: number
  expenses_by_category: CategoryExpense[]
  expense_count: number
}

// --- Products & Inventory ---

export type ProductCategory = 'hair_care' | 'facial' | 'spa' | 'coloring' | 'treatment' | 'retail' | 'equipment' | 'other'
export type ProductStatus = 'active' | 'inactive' | 'discontinued'
export type StockTransactionType = 'purchase' | 'consumption' | 'sale' | 'adjustment' | 'return' | 'damage'

export interface Product {
  id: string
  product_code: string
  name: string
  category: ProductCategory
  brand: string
  unit: string
  sku: string
  purchase_price: number
  selling_price: number
  current_stock: number
  minimum_stock: number
  maximum_stock: number
  status: ProductStatus
  created_at: string
  updated_at: string
}

export interface CreateProductInput {
  name: string
  category: ProductCategory
  brand: string
  unit: string
  sku?: string
  purchase_price: number
  selling_price: number
  minimum_stock: number
  maximum_stock: number
}

export interface UpdateProductInput extends CreateProductInput {
  status: ProductStatus
}

export interface StockTransaction {
  id: string
  product_id: string
  product_name: string
  transaction_type: StockTransactionType
  quantity: number
  unit_cost: number
  reference_type: string
  reference_id: string
  notes: string
  transaction_date: string
  created_at: string
  updated_at: string
}

export interface PurchaseEntry {
  id: string
  purchase_number: string
  vendor_name: string
  invoice_number: string
  purchase_date: string
  total_amount: number
  notes: string
  items?: PurchaseItem[]
  created_at: string
  updated_at: string
}

export interface PurchaseItem {
  id: string
  purchase_entry_id: string
  product_id: string
  product_name: string
  quantity: number
  unit_price: number
  line_total: number
}

export interface CreatePurchaseInput {
  vendor_name: string
  invoice_number?: string
  purchase_date: string
  notes?: string
  items: { product_id: string; quantity: number; unit_price: number }[]
}

export interface StockAdjustInput {
  product_id: string
  transaction_type: StockTransactionType
  quantity: number
  notes: string
}

export interface InventoryStats {
  total_products: number
  active_products: number
  low_stock_count: number
  total_value: number
  total_purchases_this_month: number
}

export interface LowStockItem {
  product_id: string
  product_code: string
  product_name: string
  category: string
  current_stock: number
  minimum_stock: number
  deficit: number
}

// =============== Analytics Types ===============

export interface DashboardStats {
  today_revenue: number
  today_customers: number
  today_invoices: number
  monthly_revenue: number
  monthly_expenses: number
  monthly_profit: number
  inventory_value: number
  outstanding_salary: number
  outstanding_advances: number
  low_stock_count: number
}

export interface KPIMetrics {
  revenue_growth_pct: number
  customer_growth_pct: number
  profit_margin_pct: number
  average_bill_value: number
  repeat_customer_pct: number
  staff_productivity_pct: number
}

export interface NameValuePair {
  name: string
  value: number
}

export interface TrendPoint {
  period: string
  value: number
}

export interface DualTrendPoint {
  period: string
  value1: number
  value2: number
}

export interface RevenueTrendPoint {
  date: string
  revenue: number
}

export interface RevenueReport {
  trend: RevenueTrendPoint[]
  by_service: NameValuePair[]
  by_staff: NameValuePair[]
  by_customer: NameValuePair[]
  total_revenue: number
  invoice_count: number
}

export interface CustomerReport {
  total_customers: number
  new_customers: number
  repeat_customers: number
  birthday_today: number
  inactive_count: number
  top_customers: NameValuePair[]
  growth_trend: TrendPoint[]
}

export interface StaffAnalyticsReport {
  top_performers: NameValuePair[]
  revenue_by_staff: NameValuePair[]
  customers_by_staff: NameValuePair[]
  commission_earned: NameValuePair[]
  salary_cost: number
}

export interface ServiceReport {
  top_services: NameValuePair[]
  least_used: NameValuePair[]
  revenue_by_service: NameValuePair[]
  avg_service_value: number
  total_bookings: number
}

export interface ExpenseAnalyticsReport {
  total_expenses: number
  by_category: NameValuePair[]
  monthly_trend: TrendPoint[]
  revenue_vs_expense: DualTrendPoint[]
}

export interface InventoryAnalyticsReport {
  total_value: number
  low_stock_count: number
  fast_moving: NameValuePair[]
  slow_moving: NameValuePair[]
  purchase_trend: TrendPoint[]
  consumption_trend: TrendPoint[]
}

export interface PLTrendPoint {
  period: string
  revenue: number
  expenses: number
  salary: number
  profit: number
}

export interface ProfitLossReport {
  revenue: number
  expenses: number
  salary_cost: number
  net_profit: number
  trend: PLTrendPoint[]
}

// =============== Backup Types ===============

export interface BackupRecord {
  id: string
  backup_name: string
  backup_type: string
  backup_path: string
  file_size: number
  checksum: string
  status: string
  error_message: string
  created_at: string
}

export interface RestoreRecord {
  id: string
  backup_id: string
  backup_name: string
  restore_date: string
  status: string
  notes: string
  error_message: string
  created_at: string
}

export interface BackupStats {
  total_backups: number
  last_backup_name: string
  last_backup_date: string
  last_backup_size: number
  last_status: string
  total_restores: number
}

export interface BackupVerification {
  backup_id: string
  file_exists: boolean
  can_open: boolean
  integrity_ok: boolean
  checksum_ok: boolean
  status: string
  error_message: string
}

// =============== License Types ===============

export interface LicenseRecord {
  id: string
  license_key: string
  customer_name: string
  salon_name: string
  device_id: string
  issued_date: string
  expiry_date: string
  status: string
  signature: string
  last_validation: string
  created_at: string
  updated_at: string
}

export interface LicenseEvent {
  id: string
  license_id: string
  event_type: string
  event_date: string
  notes: string
  created_at: string
}

export interface LicenseStatus {
  license: LicenseRecord | null
  days_remaining: number
  grace_days_remaining: number
  is_restricted: boolean
  needs_renewal: boolean
}

export interface LicenseValidation {
  valid: boolean
  status: string
  days_remaining: number
  is_restricted: boolean
  message: string
}

// =============== Update Types ===============

export interface AppVersion {
  id: string
  version: string
  release_date: string
  release_notes: string
  installed_at: string
  status: string
  created_at: string
}

export interface UpdateRecord {
  id: string
  from_version: string
  to_version: string
  update_date: string
  status: string
  error_message: string
  created_at: string
}

export interface UpdateStatus {
  current_version: string
  latest_version: string
  update_available: boolean
  status: string
  release_notes: string
}

// =============== Import Types ===============

export interface ImportJob {
  id: string
  template_id: string
  file_name: string
  file_path: string
  target_entity: string
  status: string
  total_rows: number
  valid_rows: number
  invalid_rows: number
  imported_rows: number
  column_mapping: string
  error_message: string
  created_at: string
  updated_at: string
}

export interface ImportLog {
  id: string
  job_id: string
  row_number: number
  status: string
  message: string
  row_data: string
  created_at: string
}

export interface ImportPreview {
  job_id: string
  total_rows: number
  valid_rows: number
  invalid_rows: number
  warnings: number
  headers: string[]
  sample_rows: string[][]
  errors: ImportLogRow[]
}

export interface ImportLogRow {
  row_number: number
  status: string
  message: string
}

export interface ColumnMapping {
  source_column: string
  target_field: string
}

export interface ImportUploadResult {
  job: ImportJob
  headers: string[]
  mappings: ColumnMapping[]
}

// =============== GST Types ===============

export interface GSTSettings {
  id: string
  business_name: string
  gstin: string
  state: string
  address: string
  hsn_code: string
  cgst_rate: number
  sgst_rate: number
  igst_rate: number
  is_gst_enabled: boolean
  created_at: string
  updated_at: string
}

export interface TaxRate {
  id: string
  name: string
  hsn_code: string
  cgst_rate: number
  sgst_rate: number
  igst_rate: number
  category: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface InvoiceTaxLine {
  id: string
  invoice_id: string
  item_id: string
  taxable_amount: number
  cgst_rate: number
  cgst_amount: number
  sgst_rate: number
  sgst_amount: number
  igst_rate: number
  igst_amount: number
  total_tax: number
  is_interstate: boolean
  hsn_code: string
  created_at: string
}

export interface GSTInvoiceSummary {
  invoice_id: string
  taxable_amount: number
  total_cgst: number
  total_sgst: number
  total_igst: number
  total_tax: number
  grand_total: number
  tax_lines: InvoiceTaxLine[]
}

export interface GSTReport {
  period: string
  start_date: string
  end_date: string
  total_invoices: number
  taxable_amount: number
  total_cgst: number
  total_sgst: number
  total_igst: number
  total_tax: number
  grand_total: number
}

// =============== Printer Types ===============

export interface PrinterSettings {
  id: string
  default_printer: string
  paper_width: string
  margin_top: number
  margin_bottom: number
  margin_left: number
  margin_right: number
  header_text: string
  footer_text: string
  show_logo: boolean
  show_qr: boolean
  upi_id: string
  created_at: string
  updated_at: string
}

export interface PrintJob {
  id: string
  document_type: string
  document_id: string
  printer_name: string
  paper_width: string
  status: string
  copies: number
  created_at: string
}

export interface ReceiptData {
  salon_name: string
  gstin: string
  address: string
  invoice_number: string
  date: string
  customer_name: string
  customer_phone: string
  items: ReceiptItem[]
  subtotal: number
  cgst: number
  sgst: number
  igst: number
  discount: number
  grand_total: number
  payment_method: string
  footer_text: string
  upi_id: string
}

export interface ReceiptItem {
  name: string
  quantity: number
  price: number
  total: number
}

// =============== Appointment Types ===============

export type AppointmentStatus = 'booked' | 'confirmed' | 'in_progress' | 'completed' | 'cancelled' | 'no_show'

export interface Appointment {
  id: string
  customer_id: string
  customer_name?: string
  staff_id: string
  staff_name?: string
  start_time: string
  end_time: string
  status: AppointmentStatus
  notes: string
  total_amount: number
  services?: AppointmentService[]
  created_at: string
  updated_at: string
}

export interface AppointmentService {
  id: string
  appointment_id: string
  service_id: string
  service_name?: string
  staff_id: string
  price: number
  duration: number
}

export interface AppointmentHistory {
  id: string
  appointment_id: string
  action: string
  old_value: string
  new_value: string
  changed_by: string
  created_at: string
}

export interface AppointmentFilter {
  start_date?: string
  end_date?: string
  staff_id?: string
  customer_id?: string
  status?: AppointmentStatus
}

// =============== WhatsApp Types ===============

export interface WhatsAppTemplate {
  id: string
  name: string
  category: string
  body: string
  variables: string[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface WhatsAppMessage {
  id: string
  template_id: string
  template_name?: string
  customer_id: string
  customer_name?: string
  phone: string
  body: string
  status: string
  sent_at: string
  delivered_at: string
  error_message: string
  created_at: string
}

export interface AutomationRule {
  id: string
  name: string
  trigger: string
  template_id: string
  template_name?: string
  delay_minutes: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface WAMessageStats {
  total_sent: number
  total_delivered: number
  total_failed: number
  total_pending: number
}

// =============== Membership Types ===============

export type PlanType = 'package' | 'membership'
export type SubscriptionStatus = 'active' | 'expired' | 'cancelled' | 'paused'

export interface MembershipPlan {
  id: string
  name: string
  plan_type: PlanType
  price: number
  validity_days: number
  total_sessions: number
  description: string
  is_active: boolean
  services?: PackageService[]
  created_at: string
  updated_at: string
}

export interface PackageService {
  id: string
  plan_id: string
  service_id: string
  service_name?: string
  sessions_included: number
}

export interface MemberSubscription {
  id: string
  plan_id: string
  plan_name?: string
  customer_id: string
  customer_name?: string
  start_date: string
  end_date: string
  total_sessions: number
  used_sessions: number
  remaining_sessions: number
  status: SubscriptionStatus
  amount_paid: number
  created_at: string
}

export interface MembershipStats {
  active_subscriptions: number
  total_revenue: number
  expiring_soon: number
  top_plan: string
}

// =============== Cloud Backup Types ===============

export interface CloudBackupConfig {
  id: string
  provider: string
  bucket_name: string
  region: string
  access_key: string
  endpoint: string
  encrypt_backups: boolean
  auto_backup: boolean
  auto_backup_interval_hours: number
  max_versions: number
  created_at: string
  updated_at: string
}

export interface CloudBackupHistory {
  id: string
  provider: string
  file_name: string
  file_size: number
  remote_path: string
  status: string
  is_encrypted: boolean
  error_message: string
  started_at: string
  completed_at: string
  created_at: string
}

export interface CloudBackupStats {
  last_backup_at: string
  total_backups: number
  total_size_bytes: number
  provider: string
  auto_enabled: boolean
}
