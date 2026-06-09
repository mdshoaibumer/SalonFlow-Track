-- Seed data: initial services
INSERT INTO services (id, service_code, name, category, description, duration_minutes, price, cost_price, commission_type, commission_value, status, created_at, updated_at)
VALUES
    ('019726a0-0001-7000-8000-000000000001', 'SVC-SEED0001', 'Hair Cut', 'hair', 'Standard hair cut for men and women', 30, 300, 100, 'fixed', 50, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726a0-0002-7000-8000-000000000002', 'SVC-SEED0002', 'Hair Spa', 'hair', 'Deep conditioning hair spa treatment', 60, 1200, 300, 'percentage', 10, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726a0-0003-7000-8000-000000000003', 'SVC-SEED0003', 'Facial', 'facial', 'Classic facial with cleansing and moisturizing', 45, 800, 200, 'percentage', 10, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726a0-0004-7000-8000-000000000004', 'SVC-SEED0004', 'Hair Coloring', 'coloring', 'Professional hair coloring service', 90, 2500, 800, 'percentage', 8, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726a0-0005-7000-8000-000000000005', 'SVC-SEED0005', 'Beard Trim', 'hair', 'Precision beard trimming and shaping', 15, 150, 30, 'fixed', 30, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726a0-0006-7000-8000-000000000006', 'SVC-SEED0006', 'Head Massage', 'massage', 'Relaxing head and scalp massage', 20, 200, 50, 'fixed', 40, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726a0-0007-7000-8000-000000000007', 'SVC-SEED0007', 'Keratin Treatment', 'treatment', 'Keratin smoothening for frizz-free hair', 120, 5000, 1500, 'percentage', 12, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726a0-0008-7000-8000-000000000008', 'SVC-SEED0008', 'Smoothening', 'treatment', 'Hair smoothening and straightening', 150, 6000, 2000, 'percentage', 12, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z');

-- Seed data: sample customers
INSERT INTO customers (id, customer_code, full_name, phone, email, gender, date_of_birth, address, notes, total_visits, total_spent, status, created_at, updated_at)
VALUES
    ('019726b0-0001-7000-8000-000000000001', 'CUS-SEED0001', 'Rahul Sharma', '9876543210', 'rahul@example.com', 'male', '1990-05-15', 'Mumbai', 'Regular customer', 12, 15600, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726b0-0002-7000-8000-000000000002', 'CUS-SEED0002', 'Priya Patel', '9876543211', 'priya@example.com', 'female', '1988-11-22', 'Delhi', 'Prefers hair spa', 8, 12800, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726b0-0003-7000-8000-000000000003', 'CUS-SEED0003', 'Amit Kumar', '9876543212', '', 'male', NULL, 'Pune', '', 3, 2100, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726b0-0004-7000-8000-000000000004', 'CUS-SEED0004', 'Sneha Reddy', '9876543213', 'sneha@example.com', 'female', '1995-03-08', 'Hyderabad', 'Birthday offer applied', 5, 8500, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
    ('019726b0-0005-7000-8000-000000000005', 'CUS-SEED0005', 'Vikas Singh', '9876543214', '', 'male', NULL, '', 'Walk-in customer', 1, 300, 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z');
