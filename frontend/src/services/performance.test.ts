import { describe, it, expect } from 'vitest'
import { getDailyPerformance, getWeeklyPerformance, getMonthlyPerformance, getTopPerformers, getRevenueTrend, getStaffPerformanceDetail, getStaffRevenueTrend, getPerformanceStats } from './performance'

describe('Performance Service', () => {
  it('gets daily performance', async () => {
    const result = await getDailyPerformance()
    expect(result).toEqual([])
  })

  it('gets daily performance with params', async () => {
    const result = await getDailyPerformance({ staff_id: 'staff1', date: '2024-12-19' })
    expect(result).toEqual([])
  })

  it('gets weekly performance', async () => {
    const result = await getWeeklyPerformance()
    expect(result).toEqual([])
  })

  it('gets weekly performance with params', async () => {
    const result = await getWeeklyPerformance({ staff_id: 's1', date_from: '2024-12-01', date_to: '2024-12-07' })
    expect(result).toEqual([])
  })

  it('gets monthly performance', async () => {
    const result = await getMonthlyPerformance()
    expect(result).toEqual([])
  })

  it('gets monthly performance with params', async () => {
    const result = await getMonthlyPerformance({ staff_id: 's1', date_from: '2024-12-01', date_to: '2024-12-31' })
    expect(result).toEqual([])
  })

  it('gets top performers', async () => {
    const result = await getTopPerformers()
    expect(result).toHaveLength(2)
  })

  it('gets top performers with params', async () => {
    const result = await getTopPerformers({ date_from: '2024-12-01', date_to: '2024-12-31', limit: 3 })
    expect(result).toHaveLength(2)
  })

  it('gets revenue trend', async () => {
    const result = await getRevenueTrend()
    expect(result).toHaveLength(3)
  })

  it('gets revenue trend with params', async () => {
    const result = await getRevenueTrend({ date_from: '2024-12-01', date_to: '2024-12-31' })
    expect(result).toHaveLength(3)
  })

  it('gets staff performance detail', async () => {
    const result = await getStaffPerformanceDetail('staff1')
    expect(result).toEqual([])
  })

  it('gets staff performance detail with params', async () => {
    const result = await getStaffPerformanceDetail('staff1', { date_from: '2024-12-01', date_to: '2024-12-31' })
    expect(result).toEqual([])
  })

  it('gets staff revenue trend', async () => {
    const result = await getStaffRevenueTrend('staff1')
    expect(result).toEqual([])
  })

  it('gets staff revenue trend with params', async () => {
    const result = await getStaffRevenueTrend('staff1', { date_from: '2024-12-01', date_to: '2024-12-31' })
    expect(result).toEqual([])
  })

  it('gets performance stats', async () => {
    const stats = await getPerformanceStats()
    expect(stats.total_revenue_today).toBe(12500)
  })
})
