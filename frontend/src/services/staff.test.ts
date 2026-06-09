import { describe, it, expect } from 'vitest'
import { listStaff, getStaffById, createStaff, getStaffStats } from '@/services/staff'

describe('Staff Service', () => {
  it('lists staff', async () => {
    const result = await listStaff()
    expect(result.staff).toHaveLength(1)
    expect(result.staff[0].full_name).toBe('Priya Sharma')
    expect(result.meta.total).toBe(1)
  })

  it('gets staff by id', async () => {
    const staff = await getStaffById('01912345-6789-7abc-def0-123456789001')
    expect(staff.full_name).toBe('Priya Sharma')
    expect(staff.designation).toBe('stylist')
  })

  it('creates staff', async () => {
    const staff = await createStaff({
      full_name: 'New Staff',
      phone: '9876543212',
      gender: 'male',
      designation: 'stylist',
      joining_date: '2024-06-01',
      base_salary: 15000,
      commission_percentage: 5,
    })
    expect(staff.id).toBeDefined()
    expect(staff.full_name).toBe('Priya Sharma') // mock returns fixed data
  })

  it('gets staff stats', async () => {
    const stats = await getStaffStats()
    expect(stats.total).toBe(5)
    expect(stats.active).toBe(4)
  })
})
