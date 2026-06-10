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

;(globalThis as any).window = globalThis.window || {}
;(window as any).go = {
  main: {
    StaffService: {
      ListStaff: async () => ({ staff: [mockStaff], page: 1, per_page: 20, total: 1, total_pages: 1 }),
      GetStaff: async () => mockStaff,
      CreateStaff: async () => mockStaff,
      UpdateStaff: async () => mockStaff,
      DeleteStaff: async () => {},
      GetStaffStats: async () => ({ total: 5, active: 4, inactive: 1, avg_salary: 20000 }),
    },
  },
}

beforeAll(() => server.listen({ onUnhandledRequest: 'warn' }))
afterEach(() => {
  cleanup()
  server.resetHandlers()
})
afterAll(() => server.close())
