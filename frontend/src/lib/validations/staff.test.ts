import { describe, it, expect } from 'vitest'
import { staffFormSchema, updateStaffFormSchema } from './staff'

describe('Staff Validations', () => {
  it('validates valid staff form data', () => {
    const result = staffFormSchema.safeParse({
      full_name: 'Priya Sharma',
      phone: '9876543210',
      email: 'priya@example.com',
      gender: 'female',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 10,
    })
    expect(result.success).toBe(true)
  })

  it('rejects empty name', () => {
    const result = staffFormSchema.safeParse({
      full_name: '',
      phone: '9876543210',
      gender: 'female',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 10,
    })
    expect(result.success).toBe(false)
  })

  it('rejects invalid phone number', () => {
    const result = staffFormSchema.safeParse({
      full_name: 'Test',
      phone: '1234567890',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 10,
    })
    expect(result.success).toBe(false)
  })

  it('rejects commission over 100', () => {
    const result = staffFormSchema.safeParse({
      full_name: 'Test',
      phone: '9876543210',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 101,
    })
    expect(result.success).toBe(false)
  })

  it('accepts empty email', () => {
    const result = staffFormSchema.safeParse({
      full_name: 'Test',
      phone: '9876543210',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 10,
      email: '',
    })
    expect(result.success).toBe(true)
  })

  it('rejects name over 100 chars', () => {
    const result = staffFormSchema.safeParse({
      full_name: 'A'.repeat(101),
      phone: '9876543210',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 10,
    })
    expect(result.success).toBe(false)
  })

  it('rejects negative salary', () => {
    const result = staffFormSchema.safeParse({
      full_name: 'Test',
      phone: '9876543210',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: -1,
      commission_percentage: 10,
    })
    expect(result.success).toBe(false)
  })

  it('updateStaffFormSchema requires status', () => {
    const result = updateStaffFormSchema.safeParse({
      full_name: 'Test',
      phone: '9876543210',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 10,
      status: 'active',
    })
    expect(result.success).toBe(true)
  })

  it('updateStaffFormSchema rejects invalid status', () => {
    const result = updateStaffFormSchema.safeParse({
      full_name: 'Test',
      phone: '9876543210',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 25000,
      commission_percentage: 10,
      status: 'invalid',
    })
    expect(result.success).toBe(false)
  })
})
