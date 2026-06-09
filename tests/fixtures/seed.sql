-- SalonFlow Track - Seed Data for Development/Testing
-- Run: sqlite3 salonflow.db < tests/fixtures/seed.sql

-- Staff Members
INSERT OR IGNORE INTO staff (id, staff_code, full_name, phone, email, gender, designation, joining_date, base_salary, commission_percentage, status, created_at, updated_at) VALUES
('01912345-0001-7abc-def0-100000000001', 'STF001', 'Priya Sharma', '9876543001', 'priya@salon.com', 'female', 'senior_stylist', '2023-01-15T00:00:00Z', 30000, 15, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0001-7abc-def0-100000000002', 'STF002', 'Rahul Verma', '9876543002', 'rahul@salon.com', 'male', 'stylist', '2023-03-20T00:00:00Z', 22000, 10, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0001-7abc-def0-100000000003', 'STF003', 'Sneha Patel', '9876543003', 'sneha@salon.com', 'female', 'stylist', '2023-06-01T00:00:00Z', 20000, 10, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0001-7abc-def0-100000000004', 'STF004', 'Amit Kumar', '9876543004', 'amit@salon.com', 'male', 'assistant', '2024-01-10T00:00:00Z', 12000, 5, 'active', '2024-01-10T00:00:00Z', '2024-01-10T00:00:00Z'),
('01912345-0001-7abc-def0-100000000005', 'STF005', 'Kavita Singh', '9876543005', '', 'female', 'receptionist', '2023-09-01T00:00:00Z', 18000, 0, 'inactive', '2024-01-01T00:00:00Z', '2024-06-01T00:00:00Z');

-- Services
INSERT OR IGNORE INTO services (id, service_code, name, category, description, duration_minutes, price, cost_price, commission_type, commission_value, status, created_at, updated_at) VALUES
('01912345-0002-7abc-def0-200000000001', 'SVC001', 'Haircut - Ladies', 'hair', 'Standard ladies haircut with styling', 45, 500, 100, 'percentage', 10, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0002-7abc-def0-200000000002', 'SVC002', 'Haircut - Gents', 'hair', 'Standard gents haircut', 30, 300, 50, 'percentage', 10, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0002-7abc-def0-200000000003', 'SVC003', 'Hair Color', 'hair', 'Full hair coloring with premium products', 90, 2500, 800, 'percentage', 12, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0002-7abc-def0-200000000004', 'SVC004', 'Facial - Gold', 'skin', 'Premium gold facial treatment', 60, 1500, 400, 'percentage', 10, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0002-7abc-def0-200000000005', 'SVC005', 'Manicure', 'nails', 'Complete manicure with polish', 40, 600, 150, 'flat', 50, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0002-7abc-def0-200000000006', 'SVC006', 'Pedicure', 'nails', 'Complete pedicure treatment', 50, 800, 200, 'flat', 60, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0002-7abc-def0-200000000007', 'SVC007', 'Head Massage', 'spa', 'Relaxing head and neck massage', 30, 400, 50, 'percentage', 8, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0002-7abc-def0-200000000008', 'SVC008', 'Bridal Package', 'bridal', 'Complete bridal makeup and styling', 180, 15000, 5000, 'percentage', 15, 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z');

-- Customers
INSERT OR IGNORE INTO customers (id, customer_code, full_name, phone, email, gender, address, notes, total_visits, total_spent, status, created_at, updated_at) VALUES
('01912345-0003-7abc-def0-300000000001', 'CUS001', 'Anjali Desai', '9876500001', 'anjali@gmail.com', 'female', '123 MG Road, Pune', 'VIP customer', 15, 22500, 'active', '2024-01-15T00:00:00Z', '2024-12-01T00:00:00Z'),
('01912345-0003-7abc-def0-300000000002', 'CUS002', 'Meera Kapoor', '9876500002', 'meera@gmail.com', 'female', '45 Park Street, Mumbai', 'Prefers Priya for haircuts', 8, 12000, 'active', '2024-02-10T00:00:00Z', '2024-11-15T00:00:00Z'),
('01912345-0003-7abc-def0-300000000003', 'CUS003', 'Rohan Mehta', '9876500003', '', 'male', '', '', 3, 900, 'active', '2024-06-01T00:00:00Z', '2024-10-01T00:00:00Z'),
('01912345-0003-7abc-def0-300000000004', 'CUS004', 'Sunita Agarwal', '9876500004', 'sunita@outlook.com', 'female', '78 Civil Lines, Delhi', 'Birthday: 15 March', 20, 45000, 'active', '2023-12-01T00:00:00Z', '2024-12-15T00:00:00Z'),
('01912345-0003-7abc-def0-300000000005', 'CUS005', 'Vikram Shah', '9876500005', '', 'male', '', 'Walk-in regular', 6, 1800, 'active', '2024-04-01T00:00:00Z', '2024-11-20T00:00:00Z');

-- Membership Plans
INSERT OR IGNORE INTO membership_plans (id, name, description, plan_type, price, duration_days, max_sessions, discount_percentage, priority_booking, is_active, created_at, updated_at) VALUES
('01912345-0004-7abc-def0-400000000001', 'Gold Package', '12 sessions of any hair service', 'package', 5000, 90, 12, 10, 1, 1, '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0004-7abc-def0-400000000002', 'Silver Membership', 'Monthly membership with 15% discount', 'membership', 3000, 30, 0, 15, 0, 1, '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0004-7abc-def0-400000000003', 'Bridal Package', 'Complete bridal preparation package', 'package', 25000, 60, 6, 20, 1, 1, '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z');

-- WhatsApp Templates
INSERT OR IGNORE INTO whatsapp_templates (id, name, category, body, variables, is_active, created_at, updated_at) VALUES
('01912345-0005-7abc-def0-500000000001', 'Appointment Confirmation', 'appointment', 'Hi {{name}}, your appointment on {{date}} at {{time}} is confirmed. See you at SalonFlow!', '["name","date","time"]', 1, '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0005-7abc-def0-500000000002', 'Appointment Reminder', 'appointment', 'Hi {{name}}, friendly reminder about your appointment tomorrow at {{time}}. Reply CANCEL to cancel.', '["name","time"]', 1, '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0005-7abc-def0-500000000003', 'Birthday Wish', 'marketing', 'Happy Birthday {{name}}! 🎂 Enjoy 20% off on all services today. Visit us at SalonFlow!', '["name"]', 1, '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'),
('01912345-0005-7abc-def0-500000000004', 'Payment Receipt', 'billing', 'Hi {{name}}, payment of ₹{{amount}} received. Invoice: {{invoice}}. Thank you for choosing SalonFlow!', '["name","amount","invoice"]', 1, '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z');

-- Appointments (today and upcoming)
INSERT OR IGNORE INTO appointments (id, customer_id, staff_id, appointment_date, start_time, end_time, status, notes, is_walkin, total_amount, created_at, updated_at) VALUES
('01912345-0006-7abc-def0-600000000001', '01912345-0003-7abc-def0-300000000001', '01912345-0001-7abc-def0-100000000001', '2024-12-20', '10:00', '10:45', 'confirmed', '', 0, 500, '2024-12-19T00:00:00Z', '2024-12-19T00:00:00Z'),
('01912345-0006-7abc-def0-600000000002', '01912345-0003-7abc-def0-300000000002', '01912345-0001-7abc-def0-100000000002', '2024-12-20', '11:00', '12:30', 'booked', 'Wants highlights', 0, 2500, '2024-12-19T00:00:00Z', '2024-12-19T00:00:00Z'),
('01912345-0006-7abc-def0-600000000003', '01912345-0003-7abc-def0-300000000003', '01912345-0001-7abc-def0-100000000002', '2024-12-20', '14:00', '14:30', 'booked', '', 1, 300, '2024-12-20T00:00:00Z', '2024-12-20T00:00:00Z');

-- Invoices
INSERT OR IGNORE INTO invoices (id, invoice_number, customer_id, customer_name, staff_id, staff_name, subtotal, discount, tax, grand_total, payment_method, status, notes, created_at, updated_at) VALUES
('01912345-0007-7abc-def0-700000000001', 'INV-2024-001', '01912345-0003-7abc-def0-300000000001', 'Anjali Desai', '01912345-0001-7abc-def0-100000000001', 'Priya Sharma', 500, 0, 90, 590, 'upi', 'completed', '', '2024-12-18T00:00:00Z', '2024-12-18T00:00:00Z'),
('01912345-0007-7abc-def0-700000000002', 'INV-2024-002', '01912345-0003-7abc-def0-300000000004', 'Sunita Agarwal', '01912345-0001-7abc-def0-100000000001', 'Priya Sharma', 2500, 250, 405, 2655, 'card', 'completed', 'Membership discount applied', '2024-12-19T00:00:00Z', '2024-12-19T00:00:00Z');

-- Expenses
INSERT OR IGNORE INTO expenses (id, category, amount, description, expense_date, payment_method, vendor, receipt_number, is_recurring, created_at, updated_at) VALUES
('01912345-0008-7abc-def0-800000000001', 'rent', 25000, 'Monthly shop rent', '2024-12-01', 'bank_transfer', 'Landlord', 'RNT-DEC24', 1, '2024-12-01T00:00:00Z', '2024-12-01T00:00:00Z'),
('01912345-0008-7abc-def0-800000000002', 'supplies', 8500, 'Hair color products restock', '2024-12-10', 'cash', 'Beauty Wholesale', 'BW-4521', 0, '2024-12-10T00:00:00Z', '2024-12-10T00:00:00Z'),
('01912345-0008-7abc-def0-800000000003', 'utilities', 3200, 'Electricity bill December', '2024-12-15', 'upi', 'MSEDCL', 'ELEC-DEC24', 1, '2024-12-15T00:00:00Z', '2024-12-15T00:00:00Z');
