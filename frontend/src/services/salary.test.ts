import { describe, it, expect } from 'vitest'
import { generateSalary, getSalaryById, listSalaries, paySalary, listSalaryCycles, createAdvance, approveAdvance, rejectAdvance, listAdvances, getSalaryStats } from './salary'

describe('Salary Service', () => {
  it('generates salary', async () => {
    const result = await generateSalary({ month: 12, year: 2024 } as any)
    expect(result.message).toBe('Generated')
  })

  it('gets salary by id', async () => {
    const salary = await getSalaryById('sal1')
    expect(salary.staff_name).toBe('Priya Sharma')
  })

  it('lists salaries', async () => {
    const salaries = await listSalaries(12, 2024)
    expect(salaries).toHaveLength(1)
  })

  it('pays salary', async () => {
    await expect(paySalary('sal1')).resolves.toBeUndefined()
  })

  it('lists salary cycles', async () => {
    const cycles = await listSalaryCycles(2024)
    expect(cycles).toHaveLength(1)
  })

  it('lists salary cycles defaults to current year', async () => {
    const cycles = await listSalaryCycles()
    expect(cycles).toHaveLength(1)
  })

  it('creates advance', async () => {
    const advance = await createAdvance({ staff_id: 'staff1', amount: 5000 } as any)
    expect(advance.id).toBe('adv1')
  })

  it('approves advance', async () => {
    const advance = await approveAdvance('adv1')
    expect(advance.status).toBe('approved')
  })

  it('rejects advance', async () => {
    const advance = await rejectAdvance('adv1')
    expect(advance.status).toBe('rejected')
  })

  it('lists advances', async () => {
    const result = await listAdvances()
    expect(result.advances).toHaveLength(1)
  })

  it('lists advances with params', async () => {
    const result = await listAdvances({ staff_id: 'staff1', status: 'pending', page: 1, per_page: 10 })
    expect(result.advances).toHaveLength(1)
  })

  it('gets salary stats', async () => {
    const stats = await getSalaryStats()
    expect(stats.total_payroll).toBe(150000)
  })
})
