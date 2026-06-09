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
]
