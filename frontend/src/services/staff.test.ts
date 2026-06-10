import { describe, it, expect, vi } from 'vitest'
import { listStaff, getStaffById, createStaff, updateStaff, deleteStaff, getStaffStats } from '@/services/staff'

describe('Staff Service', () => {
  it('lists staff', async () => {
    const result = await listStaff()
    expect(result.staff).toHaveLength(1)
    expect(result.staff[0].full_name).toBe('Priya Sharma')
    expect(result.meta.total).toBe(1)
  })

  it('lists staff with params', async () => {
    const result = await listStaff({ search: 'Priya', status: 'active', designation: 'stylist', page: 1, per_page: 10 })
    expect(result.staff).toHaveLength(1)
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

  it('updates staff', async () => {
    const staff = await updateStaff('01912345-6789-7abc-def0-123456789001', {
      full_name: 'Priya Updated',
      phone: '9876543210',
      gender: 'female',
      designation: 'stylist',
      joining_date: '2024-01-15',
      base_salary: 30000,
      commission_percentage: 12,
      status: 'active',
    })
    expect(staff.id).toBeDefined()
  })

  it('deletes staff', async () => {
    await expect(deleteStaff('01912345-6789-7abc-def0-123456789001')).resolves.toBeUndefined()
  })

  it('gets staff stats', async () => {
    const stats = await getStaffStats()
    expect(stats.total).toBe(5)
    expect(stats.active).toBe(4)
  })

  it('listStaff returns empty when API returns undefined staff', async () => {
    vi.spyOn(window.go.main.StaffService, 'ListStaff').mockResolvedValueOnce({ staff: undefined, total: 0 } as any)
    const r = await listStaff()
    expect(r.staff).toEqual([])
  })
})
