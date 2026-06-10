import '@testing-library/jest-dom/vitest'
import { cleanup } from '@testing-library/react'
import { afterEach, beforeAll, afterAll } from 'vitest'
import { server } from './src/mocks/server'

// Mock Wails Go bindings for desktop service tests
const mockStaff = {
  id: '01912345-6789-7abc-def0-123456789001',
  staff_code: 'STF001',
  full_name: 'Priya Sharma',
  phone: '9876543210',
  email: 'priya@example.com',
  gender: 'female',
  designation: 'stylist',
  joining_date: '2024-01-15',
  base_salary: 25000,
  commission_percentage: 10,
  status: 'active',
  created_at: '2024-01-15T00:00:00Z',
  updated_at: '2024-01-15T00:00:00Z',
}

const mockService = {
  id: '01912345-6789-7abc-def0-123456789002',
  service_code: 'SVC001',
  name: 'Haircut - Ladies',
  category: 'hair',
  description: 'Standard ladies haircut',
  duration_minutes: 45,
  price: 500,
  cost_price: 100,
  commission_type: 'percentage',
  commission_value: 10,
  status: 'active',
  created_at: '2024-01-15T00:00:00Z',
  updated_at: '2024-01-15T00:00:00Z',
}

const mockCustomer = {
  id: '01912345-6789-7abc-def0-123456789003',
  customer_code: 'CUS001',
  full_name: 'Anjali Desai',
  phone: '9876543211',
  email: 'anjali@example.com',
  gender: 'female',
  total_visits: 5,
  total_spent: 2500,
  status: 'active',
  created_at: '2024-01-15T00:00:00Z',
  updated_at: '2024-01-15T00:00:00Z',
}

const mockInvoice = {
  id: 'inv1',
  invoice_number: 'INV-001',
  customer_id: 'cust1',
  staff_id: 'staff1',
  items: [],
  subtotal: 2000,
  discount: 100,
  tax: 342,
  grand_total: 2242,
  payment_status: 'paid',
  payment_method: 'upi',
  notes: '',
  invoice_date: '2024-12-19',
  created_at: '2024-12-19T00:00:00Z',
  updated_at: '2024-12-19T00:00:00Z',
}

const mockExpense = {
  id: 'exp1',
  expense_number: 'EXP-001',
  category_id: 'cat1',
  category_name: 'Rent',
  amount: 25000,
  expense_date: '2024-12-01',
  payment_method: 'bank_transfer',
  vendor_name: 'Landlord',
  invoice_reference: '',
  description: 'Monthly rent',
  attachment_path: '',
  status: 'paid',
  created_by: '',
  created_at: '2024-12-01T00:00:00Z',
  updated_at: '2024-12-01T00:00:00Z',
}

const mockAdvance = {
  id: 'adv1',
  staff_id: 'staff2',
  staff_name: 'Rahul Kumar',
  amount: 5000,
  advance_date: '2024-12-10',
  reason: 'Personal',
  recovered_amount: 2000,
  remaining_amount: 3000,
  status: 'pending',
  created_at: '2024-12-10T00:00:00Z',
  updated_at: '2024-12-10T00:00:00Z',
}

const mockSalaryRecord = {
  id: 'sal1',
  salary_cycle_id: 'cyc1',
  staff_id: 'staff1',
  staff_name: 'Priya Sharma',
  base_salary: 25000,
  commission_amount: 5000,
  bonus_amount: 0,
  advance_amount: 2000,
  deduction_amount: 0,
  gross_salary: 30000,
  net_salary: 28000,
  payment_status: 'pending',
  payment_date: '',
  notes: '',
  created_at: '2024-12-01T00:00:00Z',
  updated_at: '2024-12-01T00:00:00Z',
}

;(globalThis as any).window = globalThis.window || {}
;(window as any).go = {
  main: {
    App: {
      GetVersion: async () => '0.2.0',
      GetEnvironment: async () => 'production',
    },
    StaffService: {
      ListStaff: async () => ({ staff: [mockStaff], page: 1, per_page: 20, total: 1, total_pages: 1 }),
      GetStaff: async () => mockStaff,
      CreateStaff: async () => mockStaff,
      UpdateStaff: async () => mockStaff,
      DeleteStaff: async () => {},
      GetStaffStats: async () => ({ total: 5, active: 4, inactive: 1, avg_salary: 20000 }),
    },
    ServiceService: {
      ListServices: async () => ({ services: [mockService], page: 1, per_page: 20, total: 1, total_pages: 1 }),
      GetService: async () => mockService,
      CreateService: async () => mockService,
      UpdateService: async () => mockService,
      DeleteService: async () => {},
      GetServiceStats: async () => ({ total: 10, active: 8, inactive: 2, avg_price: 750 }),
    },
    CustomerService: {
      ListCustomers: async () => ({ customers: [mockCustomer], page: 1, per_page: 20, total: 1, total_pages: 1 }),
      GetCustomer: async () => mockCustomer,
      CreateCustomer: async () => mockCustomer,
      UpdateCustomer: async () => mockCustomer,
      DeleteCustomer: async () => {},
      GetCustomerStats: async () => ({ total: 50, active: 45, inactive: 5, total_revenue: 125000, avg_visits: 5 }),
    },
    InvoiceService: {
      ListInvoices: async () => ({ invoices: [mockInvoice], page: 1, per_page: 20, total: 1, total_pages: 1 }),
      GetInvoice: async () => mockInvoice,
      CreateInvoice: async () => mockInvoice,
      RecordPayment: async () => ({ id: 'pay1', amount: 2242, method: 'upi' }),
      GetInvoiceStats: async () => ({ today_revenue: 12500, today_invoices: 8, avg_bill_value: 1562 }),
    },
    ExpenseService: {
      ListExpenses: async () => ({ expenses: [mockExpense], page: 1, per_page: 20, total: 1, total_pages: 1 }),
      GetExpense: async () => mockExpense,
      CreateExpense: async () => mockExpense,
      UpdateExpense: async () => mockExpense,
      DeleteExpense: async () => {},
      ListCategories: async () => [
        { id: 'cat1', name: 'Rent', description: 'Monthly rent', is_active: true, created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
        { id: 'cat2', name: 'Utilities', description: 'Electricity and water', is_active: true, created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      ],
      CreateCategory: async () => ({ id: 'cat3', name: 'New', description: '' }),
      UpdateCategory: async () => ({ id: 'cat1', name: 'Rent', description: 'Monthly rent', is_active: true }),
      GetExpenseStats: async () => ({ today_expenses: 5000, monthly_expenses: 85000, today_revenue: 12500, monthly_revenue: 300000, monthly_profit: 215000, profit_margin: 71.7 }),
      GetExpenseReport: async () => ({ categories: [], total: 0 }),
      GetProfitLoss: async () => ({ revenue: 300000, expenses: 85000, salary_cost: 0, net_profit: 215000, margin: 71.7 }),
      GetMonthlyTrend: async () => [],
    },
    SalaryService: {
      ListSalaries: async () => [mockSalaryRecord],
      GetSalary: async () => mockSalaryRecord,
      GenerateSalary: async () => ({ message: 'Generated', records: [mockSalaryRecord] }),
      PaySalary: async () => {},
      ListCycles: async () => [{ id: 'cyc1', month: 6, year: 2026, status: 'generated', generated_at: '2024-12-01', created_at: '2024-12-01' }],
      CreateAdvance: async () => mockAdvance,
      ApproveAdvance: async () => ({ ...mockAdvance, status: 'approved' }),
      RejectAdvance: async () => ({ ...mockAdvance, status: 'rejected' }),
      ListAdvances: async () => ({ advances: [mockAdvance], page: 1, per_page: 20, total: 1, total_pages: 1 }),
      GetSalaryStats: async () => ({ total_payroll: 150000, pending_payments: 3, paid_salaries: 2, outstanding_advances: 10000 }),
    },
    PerformanceService: {
      GetPerformanceStats: async () => ({
        top_performer_today: { staff_id: '1', staff_name: 'Priya Sharma', revenue: 5000, customer_count: 3, invoice_count: 3, service_count: 5, avg_bill: 1667, commission: 500, rank: 1 },
        top_performer_month: { staff_id: '1', staff_name: 'Priya Sharma', revenue: 50000, customer_count: 30, invoice_count: 30, service_count: 50, avg_bill: 1667, commission: 5000, rank: 1 },
        total_revenue_today: 12500,
        total_customers_today: 8,
        avg_bill_today: 1562,
      }),
      GetRevenueTrend: async () => [
        { date: '2024-12-01', revenue: 10000 },
        { date: '2024-12-02', revenue: 12000 },
        { date: '2024-12-03', revenue: 8000 },
      ],
      GetTopPerformers: async () => [
        { staff_id: '1', staff_name: 'Priya Sharma', revenue: 50000, customer_count: 30, invoice_count: 30, service_count: 50, avg_bill: 1667, commission: 5000, rank: 1 },
        { staff_id: '2', staff_name: 'Rahul Kumar', revenue: 35000, customer_count: 20, invoice_count: 20, service_count: 35, avg_bill: 1750, commission: 3500, rank: 2 },
      ],
      GetDailyPerformance: async () => [],
      GetWeeklyPerformance: async () => [],
      GetMonthlyPerformance: async () => [],
      GetStaffPerformance: async () => [],
      GetStaffRevenueTrend: async () => [],
    },
    MembershipService: {
      ListPlans: async () => [{ id: 'plan1', name: 'Gold', plan_type: 'package', price: 5000, duration_days: 90, max_sessions: 12, is_active: true, services: [] }],
      GetPlan: async () => ({ id: 'plan1', name: 'Gold', plan_type: 'package', price: 5000, duration_days: 90, max_sessions: 12, is_active: true, services: [] }),
      CreatePlan: async () => {},
      UpdatePlan: async () => {},
      DeletePlan: async () => {},
      SellPlan: async () => ({ id: 'sub1', plan_id: 'plan1', customer_id: 'cust1', status: 'active' }),
      UseSession: async () => {},
      ListSubscriptions: async () => [[], 0],
      GetMembershipStats: async () => ({ total_plans: 3, active_subscriptions: 15, revenue: 75000, expiring_soon: 2, top_plan: 'Gold' }),
    },
  },
}

beforeAll(() => server.listen({ onUnhandledRequest: 'warn' }))
afterEach(() => {
  cleanup()
  server.resetHandlers()
})
afterAll(() => server.close())
