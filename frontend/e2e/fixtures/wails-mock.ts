/**
 * Wails Runtime Mock for Playwright E2E Tests
 *
 * Since Wails IPC (window.go) only works in WebView2, not in regular Chromium,
 * we inject this mock via addInitScript to provide realistic demo data for testing.
 */
import { Page } from '@playwright/test'

export async function injectWailsMock(page: Page) {
  await page.addInitScript(() => {
    // Demo data
    const staffMembers = [
      { id: 'staff-1', staff_code: 'STF001', full_name: 'Priya Sharma', phone: '9876543210', email: 'priya@salon.com', gender: 'female', designation: 'stylist', joining_date: '2024-01-15', base_salary: 25000, commission_percentage: 10, status: 'active', created_at: '2024-01-15T00:00:00Z', updated_at: '2024-01-15T00:00:00Z' },
      { id: 'staff-2', staff_code: 'STF002', full_name: 'Rahul Verma', phone: '9876543211', email: 'rahul@salon.com', gender: 'male', designation: 'stylist', joining_date: '2024-02-01', base_salary: 22000, commission_percentage: 12, status: 'active', created_at: '2024-02-01T00:00:00Z', updated_at: '2024-02-01T00:00:00Z' },
      { id: 'staff-3', staff_code: 'STF003', full_name: 'Anjali Patel', phone: '9876543212', email: 'anjali@salon.com', gender: 'female', designation: 'assistant', joining_date: '2024-03-01', base_salary: 18000, commission_percentage: 8, status: 'active', created_at: '2024-03-01T00:00:00Z', updated_at: '2024-03-01T00:00:00Z' },
      { id: 'staff-4', staff_code: 'STF004', full_name: 'Vikram Singh', phone: '9876543213', email: 'vikram@salon.com', gender: 'male', designation: 'manager', joining_date: '2024-01-01', base_salary: 35000, commission_percentage: 5, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'staff-5', staff_code: 'STF005', full_name: 'Neha Gupta', phone: '9876543214', email: 'neha@salon.com', gender: 'female', designation: 'receptionist', joining_date: '2024-04-01', base_salary: 15000, commission_percentage: 0, status: 'inactive', created_at: '2024-04-01T00:00:00Z', updated_at: '2024-04-01T00:00:00Z' },
    ]

    const services = [
      { id: 'svc-1', service_code: 'SVC001', name: 'Haircut', category: 'Hair', description: 'Basic haircut', duration_minutes: 30, price: 500, cost_price: 100, commission_type: 'percentage', commission_value: 10, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'svc-2', service_code: 'SVC002', name: 'Hair Color', category: 'Coloring', description: 'Full hair coloring', duration_minutes: 90, price: 2500, cost_price: 800, commission_type: 'percentage', commission_value: 12, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'svc-3', service_code: 'SVC003', name: 'Facial', category: 'Facial', description: 'Deep cleansing facial', duration_minutes: 60, price: 1500, cost_price: 400, commission_type: 'percentage', commission_value: 10, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'svc-4', service_code: 'SVC004', name: 'Manicure', category: 'Spa', description: 'Nail care and polish', duration_minutes: 45, price: 800, cost_price: 200, commission_type: 'fixed', commission_value: 50, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'svc-5', service_code: 'SVC005', name: 'Head Massage', category: 'Massage', description: 'Relaxing head massage', duration_minutes: 30, price: 600, cost_price: 100, commission_type: 'percentage', commission_value: 8, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'svc-6', service_code: 'SVC006', name: 'Pedicure', category: 'Spa', description: 'Foot care and polish', duration_minutes: 45, price: 900, cost_price: 200, commission_type: 'percentage', commission_value: 10, status: 'inactive', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
    ]

    const customers = [
      { id: 'cust-1', customer_code: 'CUS001', full_name: 'Anjali Mehta', phone: '9988776655', email: 'anjali.m@email.com', gender: 'female', date_of_birth: '1990-05-15', anniversary_date: null, address: 'Mumbai', notes: '', total_visits: 12, total_spent: 15000, last_visit_date: '2024-05-01', status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-05-01T00:00:00Z' },
      { id: 'cust-2', customer_code: 'CUS002', full_name: 'Ravi Kumar', phone: '9988776656', email: 'ravi@email.com', gender: 'male', date_of_birth: '1985-08-20', anniversary_date: null, address: 'Delhi', notes: 'VIP customer', total_visits: 8, total_spent: 12000, last_visit_date: '2024-04-28', status: 'active', created_at: '2024-01-15T00:00:00Z', updated_at: '2024-04-28T00:00:00Z' },
      { id: 'cust-3', customer_code: 'CUS003', full_name: 'Sneha Reddy', phone: '9988776657', email: 'sneha@email.com', gender: 'female', date_of_birth: '1992-11-10', anniversary_date: '2018-02-14', address: 'Hyderabad', notes: '', total_visits: 5, total_spent: 8000, last_visit_date: '2024-04-15', status: 'active', created_at: '2024-02-01T00:00:00Z', updated_at: '2024-04-15T00:00:00Z' },
      { id: 'cust-4', customer_code: 'CUS004', full_name: 'Deepak Joshi', phone: '9988776658', email: 'deepak@email.com', gender: 'male', date_of_birth: '1988-03-25', anniversary_date: null, address: 'Pune', notes: '', total_visits: 3, total_spent: 4500, last_visit_date: '2024-03-20', status: 'active', created_at: '2024-02-15T00:00:00Z', updated_at: '2024-03-20T00:00:00Z' },
      { id: 'cust-5', customer_code: 'CUS005', full_name: 'Meera Nair', phone: '9988776659', email: 'meera@email.com', gender: 'female', date_of_birth: '1995-07-30', anniversary_date: null, address: 'Chennai', notes: '', total_visits: 1, total_spent: 1500, last_visit_date: '2024-01-10', status: 'inactive', created_at: '2024-01-10T00:00:00Z', updated_at: '2024-01-10T00:00:00Z' },
    ]

    const products = [
      { id: 'prod-1', product_code: 'PRD001', name: 'Shampoo Premium', category: 'Hair Care', brand: 'L\'Oreal', unit: 'bottle', sku: 'SHP-001', purchase_price: 400, selling_price: 650, current_stock: 25, minimum_stock: 5, maximum_stock: 50, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'prod-2', product_code: 'PRD002', name: 'Hair Color Kit', category: 'Coloring', brand: 'Garnier', unit: 'kit', sku: 'HCK-001', purchase_price: 300, selling_price: 500, current_stock: 15, minimum_stock: 3, maximum_stock: 30, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'prod-3', product_code: 'PRD003', name: 'Face Cream', category: 'Skin Care', brand: 'Neutrogena', unit: 'tube', sku: 'FCR-001', purchase_price: 250, selling_price: 450, current_stock: 2, minimum_stock: 5, maximum_stock: 20, status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
    ]

    const expenses = [
      { id: 'exp-1', category_id: 'cat-1', category_name: 'Rent', amount: 50000, description: 'Monthly rent', payment_method: 'bank_transfer', expense_date: '2024-05-01', status: 'paid', approved_by: 'staff-4', created_at: '2024-05-01T00:00:00Z', updated_at: '2024-05-01T00:00:00Z' },
      { id: 'exp-2', category_id: 'cat-2', category_name: 'Utilities', amount: 8000, description: 'Electricity bill', payment_method: 'cash', expense_date: '2024-05-05', status: 'paid', approved_by: 'staff-4', created_at: '2024-05-05T00:00:00Z', updated_at: '2024-05-05T00:00:00Z' },
      { id: 'exp-3', category_id: 'cat-3', category_name: 'Supplies', amount: 12000, description: 'Salon supplies', payment_method: 'cash', expense_date: '2024-05-10', status: 'pending', approved_by: '', created_at: '2024-05-10T00:00:00Z', updated_at: '2024-05-10T00:00:00Z' },
    ]

    const expenseCategories = [
      { id: 'cat-1', name: 'Rent', description: 'Monthly rent payments', is_active: true },
      { id: 'cat-2', name: 'Utilities', description: 'Electricity, water, internet', is_active: true },
      { id: 'cat-3', name: 'Supplies', description: 'Salon consumables', is_active: true },
      { id: 'cat-4', name: 'Marketing', description: 'Advertising and promotion', is_active: true },
      { id: 'cat-5', name: 'Maintenance', description: 'Equipment maintenance', is_active: true },
    ]

    const invoices = [
      { id: 'inv-1', invoice_number: 'INV-001', customer_id: 'cust-1', staff_id: 'staff-1', items: [{ id: 'item-1', invoice_id: 'inv-1', service_id: 'svc-1', service_name_snapshot: 'Haircut', quantity: 1, unit_price: 500, discount: 0, line_total: 500 }], subtotal: 500, discount: 0, tax: 90, grand_total: 590, payment_status: 'paid', payment_method: 'cash', notes: '', invoice_date: '2024-05-01', created_at: '2024-05-01T00:00:00Z', updated_at: '2024-05-01T00:00:00Z' },
      { id: 'inv-2', invoice_number: 'INV-002', customer_id: 'cust-2', staff_id: 'staff-2', items: [{ id: 'item-2', invoice_id: 'inv-2', service_id: 'svc-2', service_name_snapshot: 'Hair Color', quantity: 1, unit_price: 2500, discount: 200, line_total: 2300 }], subtotal: 2500, discount: 200, tax: 414, grand_total: 2714, payment_status: 'paid', payment_method: 'upi', notes: '', invoice_date: '2024-05-02', created_at: '2024-05-02T00:00:00Z', updated_at: '2024-05-02T00:00:00Z' },
      { id: 'inv-3', invoice_number: 'INV-003', customer_id: 'cust-3', staff_id: 'staff-1', items: [{ id: 'item-3', invoice_id: 'inv-3', service_id: 'svc-3', service_name_snapshot: 'Facial', quantity: 1, unit_price: 1500, discount: 0, line_total: 1500 }], subtotal: 1500, discount: 0, tax: 270, grand_total: 1770, payment_status: 'pending', payment_method: '', notes: '', invoice_date: '2024-05-03', created_at: '2024-05-03T00:00:00Z', updated_at: '2024-05-03T00:00:00Z' },
    ]

    const membershipPlans = [
      { id: 'plan-1', name: 'Gold Plan', description: 'Premium membership', plan_type: 'session_based', price: 5000, validity_days: 90, total_sessions: 10, discount_percentage: 15, services: ['svc-1', 'svc-3'], status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
      { id: 'plan-2', name: 'Silver Plan', description: 'Standard membership', plan_type: 'session_based', price: 3000, validity_days: 60, total_sessions: 5, discount_percentage: 10, services: ['svc-1'], status: 'active', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
    ]

    const appointments = [
      { id: 'appt-1', customer_id: 'cust-1', customer_name: 'Anjali Mehta', staff_id: 'staff-1', staff_name: 'Priya Sharma', appointment_date: '2024-05-15', start_time: '10:00', end_time: '10:30', status: 'confirmed', notes: '', services: [{ service_id: 'svc-1', service_name: 'Haircut', price: 500 }], created_at: '2024-05-14T00:00:00Z', updated_at: '2024-05-14T00:00:00Z' },
      { id: 'appt-2', customer_id: 'cust-2', customer_name: 'Ravi Kumar', staff_id: 'staff-2', staff_name: 'Rahul Verma', appointment_date: '2024-05-15', start_time: '11:00', end_time: '12:30', status: 'pending', notes: '', services: [{ service_id: 'svc-2', service_name: 'Hair Color', price: 2500 }], created_at: '2024-05-14T00:00:00Z', updated_at: '2024-05-14T00:00:00Z' },
      { id: 'appt-3', customer_id: 'cust-3', customer_name: 'Sneha Reddy', staff_id: 'staff-1', staff_name: 'Priya Sharma', appointment_date: '2024-05-15', start_time: '14:00', end_time: '15:00', status: 'completed', notes: '', services: [{ service_id: 'svc-3', service_name: 'Facial', price: 1500 }], created_at: '2024-05-14T00:00:00Z', updated_at: '2024-05-14T00:00:00Z' },
    ]

    const today = new Date().toISOString().split('T')[0]

    // Helper: filter/paginate list
    function paginate<T>(items: T[], page = 1, perPage = 20) {
      const total = items.length
      const totalPages = Math.ceil(total / perPage)
      const start = (page - 1) * perPage
      return { items: items.slice(start, start + perPage), page, per_page: perPage, total, total_pages: totalPages }
    }

    // Mock window.go
    ;(window as any).go = {
      main: {
        App: {
          GetVersion: async () => '0.2.0',
          GetEnvironment: async () => 'development',
        },
        StaffService: {
          ListStaff: async (params: any) => {
            let filtered = [...staffMembers]
            if (params?.search) {
              const s = params.search.toLowerCase()
              filtered = filtered.filter(st => st.full_name.toLowerCase().includes(s) || st.phone.includes(s) || st.staff_code.toLowerCase().includes(s))
            }
            if (params?.status) filtered = filtered.filter(st => st.status === params.status)
            if (params?.designation) filtered = filtered.filter(st => st.designation === params.designation)
            const p = paginate(filtered, params?.page || 1, params?.per_page || 20)
            return { staff: p.items, page: p.page, per_page: p.per_page, total: p.total, total_pages: p.total_pages }
          },
          GetStaff: async (id: string) => staffMembers.find(s => s.id === id) || staffMembers[0],
          CreateStaff: async (input: any) => ({ ...input, id: 'staff-new-' + Date.now(), staff_code: 'STF0' + (staffMembers.length + 1), status: 'active', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }),
          UpdateStaff: async (id: string, input: any) => ({ ...staffMembers.find(s => s.id === id), ...input, updated_at: new Date().toISOString() }),
          DeleteStaff: async () => {},
          GetStaffStats: async () => ({ total: 5, active: 4, inactive: 1 }),
        },
        ServiceService: {
          ListServices: async (params: any) => {
            let filtered = [...services]
            if (params?.search) {
              const s = params.search.toLowerCase()
              filtered = filtered.filter(svc => svc.name.toLowerCase().includes(s))
            }
            if (params?.status) filtered = filtered.filter(svc => svc.status === params.status)
            if (params?.category) filtered = filtered.filter(svc => svc.category === params.category)
            const p = paginate(filtered, params?.page || 1, params?.per_page || 20)
            return { services: p.items, page: p.page, per_page: p.per_page, total: p.total, total_pages: p.total_pages }
          },
          GetService: async (id: string) => services.find(s => s.id === id) || services[0],
          CreateService: async (input: any) => ({ ...input, id: 'svc-new-' + Date.now(), service_code: 'SVC0' + (services.length + 1), status: 'active', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }),
          UpdateService: async (id: string, input: any) => ({ ...services.find(s => s.id === id), ...input, updated_at: new Date().toISOString() }),
          DeleteService: async () => {},
          GetServiceStats: async () => ({ total: 6, active: 5, inactive: 1, avg_price: 1133 }),
        },
        CustomerService: {
          ListCustomers: async (params: any) => {
            let filtered = [...customers]
            if (params?.search) {
              const s = params.search.toLowerCase()
              filtered = filtered.filter(c => c.full_name.toLowerCase().includes(s) || c.phone.includes(s) || c.customer_code.toLowerCase().includes(s))
            }
            if (params?.status) filtered = filtered.filter(c => c.status === params.status)
            const p = paginate(filtered, params?.page || 1, params?.per_page || 20)
            return { customers: p.items, page: p.page, per_page: p.per_page, total: p.total, total_pages: p.total_pages }
          },
          GetCustomer: async (id: string) => customers.find(c => c.id === id) || customers[0],
          CreateCustomer: async (input: any) => ({ ...input, id: 'cust-new-' + Date.now(), customer_code: 'CUS0' + (customers.length + 1), status: 'active', total_visits: 0, total_spent: 0, created_at: new Date().toISOString(), updated_at: new Date().toISOString() }),
          UpdateCustomer: async (id: string, input: any) => ({ ...customers.find(c => c.id === id), ...input, updated_at: new Date().toISOString() }),
          DeleteCustomer: async () => {},
          GetCustomerStats: async () => ({ total: 5, active: 4, inactive: 1, new_this_month: 1, birthday_today: 0 }),
        },
        ExpenseService: {
          ListCategories: async () => expenseCategories,
          CreateCategory: async (name: string, description: string) => ({ id: 'cat-new-' + Date.now(), name, description, is_active: true }),
          UpdateCategory: async (id: string, name: string, description: string, isActive: boolean) => ({ id, name, description, is_active: isActive }),
          ListExpenses: async (params: any) => {
            let filtered = [...expenses]
            if (params?.Search) {
              const s = params.Search.toLowerCase()
              filtered = filtered.filter(e => e.description.toLowerCase().includes(s) || e.category_name.toLowerCase().includes(s))
            }
            if (params?.Status) filtered = filtered.filter(e => e.status === params.Status)
            if (params?.CategoryID) filtered = filtered.filter(e => e.category_id === params.CategoryID)
            const p = paginate(filtered, params?.Page || 1, params?.PerPage || 20)
            return { expenses: p.items, page: p.page, per_page: p.per_page, total: p.total, total_pages: p.total_pages }
          },
          GetExpense: async (id: string) => expenses.find(e => e.id === id) || expenses[0],
          CreateExpense: async (input: any) => ({ ...input, id: 'exp-new-' + Date.now(), status: 'pending', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }),
          UpdateExpense: async (id: string, input: any) => ({ ...expenses.find(e => e.id === id), ...input, updated_at: new Date().toISOString() }),
          DeleteExpense: async () => {},
          GetExpenseStats: async () => ({ today_expenses: 5000, monthly_expenses: 70000, today_revenue: 15000, monthly_revenue: 350000, monthly_profit: 280000, profit_margin: 80 }),
          GetExpenseReport: async () => ({ date_from: '2024-05-01', date_to: '2024-05-31', total_expenses: 70000, expenses_by_category: [{ category_id: 'cat-1', category_name: 'Rent', amount: 50000, percentage: 71 }], expense_count: 3 }),
          GetProfitLoss: async () => ({ period: '2024-05', total_revenue: 350000, total_expenses: 70000, gross_profit: 280000, profit_margin: 80, expenses_by_category: [] }),
          GetMonthlyTrend: async () => [{ month: '2024-03', revenue: 300000, expenses: 60000, profit: 240000 }, { month: '2024-04', revenue: 320000, expenses: 65000, profit: 255000 }, { month: '2024-05', revenue: 350000, expenses: 70000, profit: 280000 }],
        },
        InvoiceService: {
          ListInvoices: async (params: any) => {
            let filtered = [...invoices]
            if (params?.search) {
              const s = params.search.toLowerCase()
              filtered = filtered.filter(inv => inv.invoice_number.toLowerCase().includes(s))
            }
            if (params?.payment_status) filtered = filtered.filter(inv => inv.payment_status === params.payment_status)
            const p = paginate(filtered, params?.page || 1, params?.per_page || 20)
            return { invoices: p.items, page: p.page, per_page: p.per_page, total: p.total, total_pages: p.total_pages }
          },
          GetInvoice: async (id: string) => invoices.find(i => i.id === id) || invoices[0],
          CreateInvoice: async (input: any) => ({ ...input, id: 'inv-new-' + Date.now(), invoice_number: 'INV-' + (invoices.length + 1).toString().padStart(3, '0'), payment_status: 'pending', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }),
          RecordPayment: async () => ({ id: 'pay-1', amount: 1000, method: 'cash', date: today }),
          GetInvoiceStats: async () => ({ today_revenue: 5074, today_invoices: 3, avg_bill_value: 1691 }),
        },
        ProductService: {
          ListProducts: async (params: any) => {
            let filtered = [...products]
            if (params?.Search) {
              const s = params.Search.toLowerCase()
              filtered = filtered.filter(p => p.name.toLowerCase().includes(s))
            }
            if (params?.Status) filtered = filtered.filter(p => p.status === params.Status)
            if (params?.Category) filtered = filtered.filter(p => p.category === params.Category)
            const pg = paginate(filtered, params?.Page || 1, params?.PerPage || 20)
            return { products: pg.items, page: pg.page, per_page: pg.per_page, total: pg.total, total_pages: pg.total_pages }
          },
          GetProduct: async (id: string) => products.find(p => p.id === id) || products[0],
          CreateProduct: async (input: any) => ({ ...input, id: 'prod-new-' + Date.now(), product_code: 'PRD0' + (products.length + 1), status: 'active', current_stock: 0, created_at: new Date().toISOString(), updated_at: new Date().toISOString() }),
          UpdateProduct: async (id: string, input: any) => ({ ...products.find(p => p.id === id), ...input, updated_at: new Date().toISOString() }),
          DeleteProduct: async () => {},
          AdjustStock: async (input: any) => ({ id: 'stx-1', ...input, created_at: new Date().toISOString() }),
          ListStockHistory: async (params: any) => ({ transactions: [], meta: { page: params?.Page || 1, per_page: params?.PerPage || 20, total: 0, total_pages: 0 } }),
          CreatePurchase: async (input: any) => ({ ...input, id: 'pur-1', created_at: new Date().toISOString() }),
          ListPurchases: async (params: any) => ({ purchases: [], meta: { page: params?.Page || 1, per_page: params?.PerPage || 20, total: 0, total_pages: 0 } }),
          GetInventoryStats: async () => ({ total_products: 3, active_products: 3, low_stock_count: 1, total_value: 12750, total_purchases_this_month: 5000 }),
          GetLowStockProducts: async () => [{ product_id: 'prod-3', product_code: 'PRD003', product_name: 'Face Cream', category: 'Skin Care', current_stock: 2, minimum_stock: 5, deficit: 3 }],
        },
        MembershipService: {
          CreatePlan: async () => {},
          UpdatePlan: async () => {},
          DeletePlan: async () => {},
          GetPlan: async (id: string) => membershipPlans.find(p => p.id === id) || membershipPlans[0],
          ListPlans: async () => membershipPlans,
          SellPlan: async () => ({ id: 'sub-1', customer_id: 'cust-1', plan_id: 'plan-1', status: 'active', sessions_remaining: 10, start_date: today, end_date: '2024-08-01' }),
          UseSession: async () => {},
          ListSubscriptions: async () => [[], 0],
          GetMembershipStats: async () => ({ total_plans: 2, active_subscriptions: 3, expired_subscriptions: 1, revenue_this_month: 8000, total_sessions_used: 15 }),
        },
        AppointmentService: {
          CreateAppointment: async () => {},
          UpdateAppointment: async () => {},
          UpdateAppointmentStatus: async () => {},
          DeleteAppointment: async () => {},
          GetAppointment: async (id: string) => appointments.find(a => a.id === id) || appointments[0],
          ListAppointments: async () => [appointments, appointments.length],
          GetAppointmentHistory: async () => [],
        },
        PerformanceService: {
          GetPerformanceStats: async () => ({ top_performer_today: { staff_id: 'staff-1', staff_name: 'Priya Sharma', revenue: 5000, customer_count: 4, invoice_count: 4, service_count: 5, avg_bill: 1250, commission: 500, rank: 1 }, top_performer_month: { staff_id: 'staff-1', staff_name: 'Priya Sharma', revenue: 120000, customer_count: 45, invoice_count: 50, service_count: 60, avg_bill: 2400, commission: 12000, rank: 1 }, total_revenue_today: 15000, total_customers_today: 8, avg_bill_today: 1875 }),
          GetDailyPerformance: async () => [],
          GetWeeklyPerformance: async () => [],
          GetMonthlyPerformance: async () => [],
          GetTopPerformers: async () => [{ staff_id: 'staff-1', staff_name: 'Priya Sharma', revenue: 120000, customer_count: 45, invoice_count: 50, service_count: 60, avg_bill: 2400, commission: 12000, rank: 1 }, { staff_id: 'staff-2', staff_name: 'Rahul Verma', revenue: 95000, customer_count: 38, invoice_count: 40, service_count: 48, avg_bill: 2375, commission: 9500, rank: 2 }],
          GetRevenueTrend: async () => Array.from({ length: 14 }, (_, i) => ({ date: new Date(Date.now() - (13 - i) * 86400000).toISOString().split('T')[0], revenue: 10000 + Math.floor(Math.random() * 5000) })),
          GetStaffPerformance: async () => [],
          GetStaffRevenueTrend: async () => [],
        },
        AnalyticsService: {
          GetDashboard: async () => ({ today_revenue: 15000, today_customers: 8, today_invoices: 6, monthly_revenue: 350000, monthly_expenses: 70000, monthly_profit: 280000, inventory_value: 12750, outstanding_salary: 0, outstanding_advances: 5000, low_stock_count: 1 }),
          GetKPIs: async () => ({ revenue_growth_pct: 12.5, customer_growth_pct: 8.3, profit_margin_pct: 80, average_bill_value: 1875, repeat_customer_pct: 65, staff_productivity_pct: 78 }),
          GetRevenueReport: async () => ({ trend: [], by_service: [], by_staff: [], by_customer: [], total_revenue: 350000, invoice_count: 180 }),
          GetCustomerReport: async () => ({ total_customers: 5, new_customers: 1, repeat_customers: 4, birthday_today: 0, inactive_count: 1, top_customers: [], growth_trend: [] }),
          GetStaffReport: async () => ({ top_performers: [], revenue_by_staff: [], customers_by_staff: [], commission_earned: [], salary_cost: 115000 }),
          GetServiceReport: async () => ({ top_services: [], least_used: [], revenue_by_service: [], avg_service_value: 1133, total_bookings: 180 }),
          GetExpenseAnalytics: async () => ({ total_expenses: 70000, by_category: [], monthly_trend: [], revenue_vs_expense: [] }),
          GetInventoryReport: async () => ({ total_value: 12750, low_stock_count: 1, fast_moving: [], slow_moving: [], purchase_trend: [], consumption_trend: [] }),
          GetProfitLossReport: async () => ({ revenue: 350000, expenses: 70000, salary_cost: 115000, net_profit: 165000, trend: [] }),
        },
        CommissionService: {
          CreateRule: async (input: any) => ({ ...input, id: 'rule-new', created_at: new Date().toISOString() }),
          GetRule: async () => ({ id: 'rule-1', rule_type: 'percentage', target_type: 'service', target_id: '', value: 10, is_active: true }),
          ListRules: async () => ({ rules: [{ id: 'rule-1', rule_type: 'percentage', target_type: 'service', target_id: '', value: 10, is_active: true }], page: 1, per_page: 20, total: 1, total_pages: 1 }),
          UpdateRule: async (id: string, input: any) => ({ id, ...input }),
          DeleteRule: async () => {},
          GetStaffCommission: async () => ({ staff_id: 'staff-1', staff_name: 'Priya Sharma', total_commission: 12000, transactions: [] }),
          GetMonthlyCommission: async () => ({ month: '2024-05', staff_commissions: [] }),
          GetCommissionStats: async () => ({ total_commission_this_month: 45000, top_earner: 'Priya Sharma', avg_commission: 9000 }),
        },
        SalaryService: {
          GenerateSalary: async () => ({ cycle: { id: 'cycle-1', month: 5, year: 2024, status: 'generated' }, records: [] }),
          GetSalary: async () => ({ id: 'sal-1', staff_id: 'staff-1', staff_name: 'Priya Sharma', month: 5, year: 2024, base_salary: 25000, commission: 12000, deductions: 0, advances: 0, net_salary: 37000, status: 'pending' }),
          ListSalaries: async () => staffMembers.filter(s => s.status === 'active').map((s, i) => ({ id: `sal-${i + 1}`, staff_id: s.id, staff_name: s.full_name, month: 5, year: 2024, base_salary: s.base_salary, commission: 5000, deductions: 0, advances: 0, net_salary: s.base_salary + 5000, status: 'pending' })),
          PaySalary: async () => {},
          ListCycles: async () => [{ id: 'cycle-1', month: 5, year: 2024, status: 'generated', total_salary: 150000 }],
          CreateAdvance: async (input: any) => ({ ...input, id: 'adv-new', status: 'pending', created_at: new Date().toISOString() }),
          ApproveAdvance: async (id: string) => ({ id, status: 'approved' }),
          RejectAdvance: async (id: string) => ({ id, status: 'rejected' }),
          ListAdvances: async () => ({ advances: [{ id: 'adv-1', staff_id: 'staff-3', staff_name: 'Anjali Patel', amount: 5000, advance_date: '2024-05-01', reason: 'Medical emergency', recovered_amount: 2000, remaining_amount: 3000, status: 'approved', created_at: '2024-05-01T00:00:00Z' }], page: 1, per_page: 20, total: 1, total_pages: 1 }),
          GetSalaryStats: async () => ({ total_payroll: 150000, pending_payments: 4, paid_salaries: 0, outstanding_advances: 5000 }),
        },
        BackupService: {
          CreateBackup: async (type: string) => ({ id: 'bak-1', backup_type: type, file_name: 'backup-2024-05-15.db', file_size: 1024000, status: 'completed', created_at: new Date().toISOString() }),
          VerifyBackup: async () => ({ id: 'bak-1', is_valid: true, checked_at: new Date().toISOString() }),
          RestoreBackup: async () => ({ id: 'res-1', backup_id: 'bak-1', status: 'completed', restored_at: new Date().toISOString() }),
          ListBackups: async () => [[{ id: 'bak-1', backup_type: 'manual', file_name: 'backup-2024-05-15.db', file_size: 1024000, status: 'completed', created_at: '2024-05-15T00:00:00Z' }], 1],
          ListRestores: async () => [[], 0],
          GetBackupStats: async () => ({ total_backups: 5, last_backup_name: 'backup-2024-05-15.db', last_backup_date: '2024-05-15T00:00:00Z', last_backup_size: 1024000, last_status: 'completed', total_restores: 0 }),
          DeleteBackup: async () => {},
        },
        WhatsAppService: {
          CreateTemplate: async () => {},
          UpdateTemplate: async () => {},
          DeleteTemplate: async () => {},
          ListTemplates: async () => [{ id: 'tpl-1', name: 'Welcome', content: 'Welcome {{name}}!', category: 'greeting', variables: ['name'], is_active: true, created_at: '2024-01-01T00:00:00Z' }],
          SendMessage: async () => ({ id: 'msg-1', phone: '9876543210', template_id: 'tpl-1', status: 'sent', sent_at: new Date().toISOString() }),
          ListMessages: async () => [[], 0],
          GetWhatsAppStats: async () => ({ total_sent: 50, total_delivered: 45, total_failed: 5, sent_today: 3 }),
          CreateRule: async () => {},
          UpdateRule: async () => {},
          DeleteRule: async () => {},
          ListRules: async () => [{ id: 'auto-1', name: 'Birthday Greeting', trigger: 'birthday', template_id: 'tpl-1', is_active: true }],
        },
        GSTService: {
          GetSettings: async () => ({ gst_enabled: true, gst_number: '29AABCU9603R1ZM', business_name: 'SalonFlow Salon', address: 'Mumbai, MH', hsn_code: '9985', default_rate: 18, inclusive: false }),
          SaveSettings: async () => {},
          ListTaxRates: async () => [{ id: 'tax-1', category: 'service', name: 'Service Tax', rate: 18, hsn_code: '9985', is_active: true }],
          CreateTaxRate: async () => {},
          UpdateTaxRate: async () => {},
          DeleteTaxRate: async () => {},
          GetReport: async () => ({ period: '2024-05', total_taxable: 350000, total_cgst: 31500, total_sgst: 31500, total_igst: 0, total_tax: 63000, invoices: [] }),
        },
        PrinterService: {
          GetSettings: async () => ({ printer_name: 'Default', paper_size: '80mm', auto_print: false, header_text: 'SalonFlow Salon', footer_text: 'Thank you!' }),
          SaveSettings: async () => {},
          TestPrint: async () => ({ success: true, message: 'Test print sent' }),
          PrintReceipt: async () => ({ success: true }),
        },
        ImportService: {
          PreviewImport: async () => ({ total_rows: 0, valid_rows: 0, errors: [], preview: [] }),
          ExecuteImport: async () => ({ id: 'imp-1', status: 'completed', imported: 0, failed: 0 }),
          ListImports: async () => [],
          GetImportStatus: async () => ({ id: 'imp-1', status: 'completed', imported: 0, failed: 0, created_at: new Date().toISOString() }),
        },
        LicenseService: {
          GetStatus: async () => ({ is_valid: true, license_type: 'professional', expires_at: '2025-12-31T00:00:00Z', features: ['all'], max_staff: 50, max_customers: 10000 }),
          Validate: async () => ({ is_valid: true }),
          Activate: async () => {},
          Deactivate: async () => {},
        },
        UpdateService: {
          CheckForUpdate: async () => ({ available: false, current_version: '0.2.0', latest_version: '0.2.0' }),
          DownloadUpdate: async () => {},
          InstallUpdate: async () => {},
          GetUpdateStatus: async () => ({ status: 'up_to_date', current_version: '0.2.0' }),
        },
        CloudBackupService: {
          GetConfig: async () => ({ provider: 'google_drive', enabled: false, auto_backup: false, frequency: 'daily', retention_days: 30 }),
          SaveConfig: async () => {},
          TestConnection: async () => ({ success: true, message: 'Connection successful' }),
          RunBackup: async () => ({ id: 'cbak-1', status: 'completed', provider: 'google_drive', file_size: 1024000, created_at: new Date().toISOString() }),
          ListHistory: async () => ({ backups: [], total: 0 }),
          GetStats: async () => ({ total_backups: 0, last_backup_date: null, total_size: 0, provider: 'google_drive' }),
        },
      },
    }

    // Mock window.runtime (Wails runtime functions)
    ;(window as any).runtime = {
      LogPrint: () => {},
      LogTrace: () => {},
      LogDebug: () => {},
      LogInfo: () => {},
      LogWarning: () => {},
      LogError: () => {},
      LogFatal: () => {},
      EventsOn: () => {},
      EventsOff: () => {},
      EventsOnce: () => {},
      EventsOnMultiple: () => {},
      EventsEmit: () => {},
      WindowReload: () => {},
      WindowReloadApp: () => {},
      WindowSetAlwaysOnTop: () => {},
      WindowSetSystemDefaultTheme: () => {},
      WindowSetLightTheme: () => {},
      WindowSetDarkTheme: () => {},
      WindowCenter: () => {},
      WindowSetTitle: () => {},
      WindowFullscreen: () => {},
      WindowUnfullscreen: () => {},
      WindowIsFullscreen: () => false,
      WindowSetSize: () => {},
      WindowGetSize: () => ({ w: 1400, h: 900 }),
      WindowSetMaxSize: () => {},
      WindowSetMinSize: () => {},
      WindowSetPosition: () => {},
      WindowGetPosition: () => ({ x: 0, y: 0 }),
      WindowHide: () => {},
      WindowShow: () => {},
      WindowMaximise: () => {},
      WindowToggleMaximise: () => {},
      WindowUnmaximise: () => {},
      WindowIsMaximised: () => false,
      WindowMinimise: () => {},
      WindowUnminimise: () => {},
      WindowIsMinimised: () => false,
      WindowIsNormal: () => true,
      WindowSetBackgroundColour: () => {},
      BrowserOpenURL: () => {},
      Environment: () => ({ buildType: 'dev', platform: 'windows', arch: 'amd64' }),
      Quit: () => {},
      Hide: () => {},
      Show: () => {},
      ClipboardGetText: () => '',
      ClipboardSetText: () => {},
      OnFileDrop: () => {},
      OnFileDropOff: () => {},
    }
  })
}
