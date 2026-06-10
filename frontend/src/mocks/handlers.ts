import { http, HttpResponse } from 'msw'

const BASE = 'http://localhost:8080/api/v1'

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

const mockAppointment = {
  id: '01912345-6789-7abc-def0-123456789004',
  customer_id: mockCustomer.id,
  staff_id: mockStaff.id,
  appointment_date: '2024-12-20',
  start_time: '10:00',
  end_time: '11:00',
  status: 'booked',
  notes: '',
  is_walkin: false,
  total_amount: 500,
  services: [],
  created_at: '2024-12-19T00:00:00Z',
  updated_at: '2024-12-19T00:00:00Z',
}

export const handlers = [
  // Staff
  http.get(`${BASE}/staff`, () => {
    return HttpResponse.json({
      success: true,
      data: [mockStaff],
      meta: { page: 1, per_page: 20, total: 1, total_pages: 1 },
    })
  }),
  http.post(`${BASE}/staff`, () => {
    return HttpResponse.json({ success: true, data: mockStaff }, { status: 201 })
  }),
  http.get(`${BASE}/staff/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { total: 5, active: 4, inactive: 1, on_leave: 0 },
    })
  }),
  http.get(`${BASE}/staff/:id`, () => {
    return HttpResponse.json({ success: true, data: mockStaff })
  }),
  http.put(`${BASE}/staff/:id`, () => {
    return HttpResponse.json({ success: true, data: mockStaff })
  }),
  http.delete(`${BASE}/staff/:id`, () => {
    return HttpResponse.json({ success: true, data: null })
  }),

  // Services
  http.get(`${BASE}/services`, () => {
    return HttpResponse.json({
      success: true,
      data: [mockService],
      meta: { page: 1, per_page: 20, total: 1, total_pages: 1 },
    })
  }),
  http.post(`${BASE}/services`, () => {
    return HttpResponse.json({ success: true, data: mockService }, { status: 201 })
  }),
  http.get(`${BASE}/services/:id`, () => {
    return HttpResponse.json({ success: true, data: mockService })
  }),

  // Customers
  http.get(`${BASE}/customers`, () => {
    return HttpResponse.json({
      success: true,
      data: [mockCustomer],
      meta: { page: 1, per_page: 20, total: 1, total_pages: 1 },
    })
  }),
  http.post(`${BASE}/customers`, () => {
    return HttpResponse.json({ success: true, data: mockCustomer }, { status: 201 })
  }),
  http.get(`${BASE}/customers/:id`, () => {
    return HttpResponse.json({ success: true, data: mockCustomer })
  }),

  // Appointments
  http.get(`${BASE}/appointments`, () => {
    return HttpResponse.json({
      success: true,
      data: [mockAppointment],
      meta: { page: 1, per_page: 50, total: 1, total_pages: 1 },
    })
  }),
  http.post(`${BASE}/appointments`, () => {
    return HttpResponse.json({ success: true, data: mockAppointment }, { status: 201 })
  }),
  http.get(`${BASE}/appointments/:id`, () => {
    return HttpResponse.json({ success: true, data: mockAppointment })
  }),

  // Memberships
  http.get(`${BASE}/memberships/plans`, () => {
    return HttpResponse.json({
      success: true,
      data: [{ id: 'plan1', name: 'Gold', plan_type: 'package', price: 5000, duration_days: 90, max_sessions: 12, is_active: true }],
      meta: { page: 1, per_page: 20, total: 1, total_pages: 1 },
    })
  }),
  http.post(`${BASE}/memberships/plans`, () => {
    return HttpResponse.json({ success: true, data: { id: 'plan1', name: 'Gold' } }, { status: 201 })
  }),
  http.get(`${BASE}/memberships/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { total_plans: 3, active_subscriptions: 15, revenue: 75000 },
    })
  }),

  // WhatsApp
  http.get(`${BASE}/whatsapp/templates`, () => {
    return HttpResponse.json({
      success: true,
      data: [{ id: 'tpl1', name: 'Welcome', category: 'general', body: 'Hello!', is_active: true }],
    })
  }),
  http.post(`${BASE}/whatsapp/templates`, () => {
    return HttpResponse.json({ success: true, data: { id: 'tpl1', name: 'Welcome' } }, { status: 201 })
  }),
  http.get(`${BASE}/whatsapp/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { total_sent: 100, delivered: 95, read: 80, failed: 3, queued: 2 },
    })
  }),

  // Cloud Backup
  http.get(`${BASE}/cloud-backup/config`, () => {
    return HttpResponse.json({
      success: true,
      data: { provider: 'none', auto_backup: false },
    })
  }),
  http.get(`${BASE}/cloud-backup/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { total_backups: 5, last_backup: '2024-12-19T00:00:00Z', total_size: 1024000 },
    })
  }),

  // Health
  http.get(`${BASE}/health`, () => {
    return HttpResponse.json({
      success: true,
      data: { status: 'healthy', version: 'v0.1.0', database: 'sqlite', uptime: '2h30m' },
    })
  }),

  // Performance Stats
  http.get(`${BASE}/performance/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: {
        top_performer_today: { staff_id: '1', staff_name: 'Priya Sharma', revenue: 5000, customer_count: 3, invoice_count: 3, service_count: 5, avg_bill: 1667, commission: 500, rank: 1 },
        top_performer_month: { staff_id: '1', staff_name: 'Priya Sharma', revenue: 50000, customer_count: 30, invoice_count: 30, service_count: 50, avg_bill: 1667, commission: 5000, rank: 1 },
        total_revenue_today: 12500,
        total_customers_today: 8,
        avg_bill_today: 1562,
      },
    })
  }),
  http.get(`${BASE}/performance/revenue-trend`, () => {
    return HttpResponse.json({
      success: true,
      data: [
        { date: '2024-12-01', revenue: 10000 },
        { date: '2024-12-02', revenue: 12000 },
        { date: '2024-12-03', revenue: 8000 },
      ],
    })
  }),
  http.get(`${BASE}/performance/top-performers`, () => {
    return HttpResponse.json({
      success: true,
      data: [
        { staff_id: '1', staff_name: 'Priya Sharma', revenue: 50000, customer_count: 30, invoice_count: 30, service_count: 50, avg_bill: 1667, commission: 5000, rank: 1 },
        { staff_id: '2', staff_name: 'Rahul Kumar', revenue: 35000, customer_count: 20, invoice_count: 20, service_count: 35, avg_bill: 1750, commission: 3500, rank: 2 },
      ],
    })
  }),

  // Invoices
  http.get(`${BASE}/invoices`, () => {
    return HttpResponse.json({
      success: true,
      data: [{
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
      }],
      meta: { page: 1, per_page: 20, total: 1, total_pages: 1 },
    })
  }),
  http.get(`${BASE}/invoices/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { today_revenue: 12500, today_invoices: 8, avg_bill_value: 1562 },
    })
  }),

  // Expenses
  http.get(`${BASE}/expenses`, () => {
    return HttpResponse.json({
      success: true,
      data: [{
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
      }],
      meta: { page: 1, per_page: 20, total: 1, total_pages: 1 },
    })
  }),
  http.get(`${BASE}/expenses/categories`, () => {
    return HttpResponse.json({
      success: true,
      data: [
        { id: 'cat1', name: 'Rent', description: 'Monthly rent', is_active: true, created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
        { id: 'cat2', name: 'Utilities', description: 'Electricity and water', is_active: true, created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      ],
    })
  }),
  http.post(`${BASE}/expenses`, () => {
    return HttpResponse.json({ success: true, data: { id: 'exp2', expense_number: 'EXP-002', amount: 500, status: 'pending' } }, { status: 201 })
  }),
  http.delete(`${BASE}/expenses/:id`, () => {
    return HttpResponse.json({ success: true, data: null })
  }),
  http.get(`${BASE}/expenses/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { today_expenses: 5000, monthly_expenses: 85000, today_revenue: 12500, monthly_revenue: 300000, monthly_profit: 215000, profit_margin: 71.7 },
    })
  }),

  // Salary
  http.get(`${BASE}/salary`, () => {
    return HttpResponse.json({
      success: true,
      data: [{
        id: 'sal1',
        staff_name: 'Priya Sharma',
        base_salary: 25000,
        commission_amount: 5000,
        bonus_amount: 0,
        advance_amount: 2000,
        deduction_amount: 0,
        net_salary: 28000,
        payment_status: 'pending',
      }],
    })
  }),
  http.get(`${BASE}/salary/cycles`, () => {
    return HttpResponse.json({
      success: true,
      data: [{ month: 6, year: 2026, status: 'generated' }],
    })
  }),
  http.get(`${BASE}/salary/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { total_payroll: 150000, pending_payments: 3, paid_salaries: 2, outstanding_advances: 10000 },
    })
  }),
  http.post(`${BASE}/salary/generate`, () => {
    return HttpResponse.json({ success: true, data: { message: 'Salary generated' } }, { status: 201 })
  }),
  http.post(`${BASE}/salary/:id/pay`, () => {
    return HttpResponse.json({ success: true, data: { message: 'Paid' } })
  }),

  // Advances
  http.get(`${BASE}/salary/advances`, () => {
    return HttpResponse.json({
      success: true,
      data: [{
        id: 'adv1',
        staff_name: 'Rahul Kumar',
        amount: 5000,
        advance_date: '2024-12-10',
        reason: 'Personal',
        recovered_amount: 2000,
        remaining_amount: 3000,
        status: 'pending',
      }],
      meta: { page: 1, per_page: 20, total: 1, total_pages: 1 },
    })
  }),
  http.post(`${BASE}/salary/advances`, () => {
    return HttpResponse.json({ success: true, data: { id: 'adv2' } }, { status: 201 })
  }),
  http.put(`${BASE}/salary/advances/:id/approve`, () => {
    return HttpResponse.json({ success: true, data: { message: 'Approved' } })
  }),
  http.put(`${BASE}/salary/advances/:id/reject`, () => {
    return HttpResponse.json({ success: true, data: { message: 'Rejected' } })
  }),

  // Settings
  http.get(`${BASE}/settings`, () => {
    return HttpResponse.json({
      success: true,
      data: [
        { id: 's1', key: 'salon_name', value: 'SalonFlow Studio', description: 'Salon name', category: 'general', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
        { id: 's2', key: 'salon_phone', value: '9876543210', description: 'Phone', category: 'general', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
        { id: 's3', key: 'invoice_prefix', value: 'INV', description: 'Invoice prefix', category: 'billing', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
        { id: 's4', key: 'tax_rate', value: '18', description: 'Tax rate', category: 'billing', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      ],
    })
  }),

  // Services stats / delete / update
  http.get(`${BASE}/services/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { total: 8, active: 6, inactive: 2, avg_price: 750 },
    })
  }),
  http.put(`${BASE}/services/:id`, () => {
    return HttpResponse.json({ success: true, data: mockService })
  }),
  http.delete(`${BASE}/services/:id`, () => {
    return HttpResponse.json({ success: true, data: null })
  }),

  // Customers stats / delete / update
  http.get(`${BASE}/customers/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: { total: 50, active: 45, inactive: 5, new_this_month: 3, birthday_today: 1 },
    })
  }),
  http.put(`${BASE}/customers/:id`, () => {
    return HttpResponse.json({ success: true, data: mockCustomer })
  }),
  http.delete(`${BASE}/customers/:id`, () => {
    return HttpResponse.json({ success: true, data: null })
  }),
]
