import { describe, it, expect } from 'vitest'
import { getDashboardStats, getKPIMetrics, getRevenueReport, getCustomerReport, getStaffReport, getServiceReport, getExpenseReport, getInventoryReport, getProfitLossReport } from './analytics'

describe('Analytics Service', () => {
  it('gets dashboard stats', async () => {
    const stats = await getDashboardStats()
    expect(stats.today_revenue).toBe(12500)
  })

  it('gets KPI metrics', async () => {
    const kpis = await getKPIMetrics()
    expect(kpis.revenue).toBe(300000)
  })

  it('gets KPI metrics with params', async () => {
    const kpis = await getKPIMetrics({ date_from: '2024-01-01', date_to: '2024-12-31' })
    expect(kpis.revenue).toBe(300000)
  })

  it('gets revenue report', async () => {
    const report = await getRevenueReport()
    expect(report.total).toBe(300000)
  })

  it('gets revenue report with params', async () => {
    const report = await getRevenueReport({ date_from: '2024-01-01', date_to: '2024-12-31', group_by: 'month' })
    expect(report.total).toBe(300000)
  })

  it('gets customer report', async () => {
    const report = await getCustomerReport()
    expect(report.total_customers).toBe(50)
  })

  it('gets customer report with params', async () => {
    const report = await getCustomerReport({ date_from: '2024-01-01', date_to: '2024-12-31' })
    expect(report.total_customers).toBe(50)
  })

  it('gets staff report', async () => {
    const report = await getStaffReport()
    expect(report.total_revenue).toBe(300000)
  })

  it('gets staff report with params', async () => {
    const report = await getStaffReport({ date_from: '2024-01-01', date_to: '2024-12-31' })
    expect(report.total_revenue).toBe(300000)
  })

  it('gets service report', async () => {
    const report = await getServiceReport()
    expect(report.total_revenue).toBe(300000)
  })

  it('gets service report with params', async () => {
    const report = await getServiceReport({ date_from: '2024-01-01', date_to: '2024-12-31' })
    expect(report.total_revenue).toBe(300000)
  })

  it('gets expense report', async () => {
    const report = await getExpenseReport()
    expect(report.total).toBe(85000)
  })

  it('gets expense report with params', async () => {
    const report = await getExpenseReport({ date_from: '2024-01-01', date_to: '2024-12-31' })
    expect(report.total).toBe(85000)
  })

  it('gets inventory report', async () => {
    const report = await getInventoryReport()
    expect(report.total_products).toBe(30)
  })

  it('gets profit loss report', async () => {
    const report = await getProfitLossReport()
    expect(report.net_profit).toBe(215000)
  })

  it('gets profit loss report with params', async () => {
    const report = await getProfitLossReport({ date_from: '2024-01-01', date_to: '2024-12-31', group_by: 'month' })
    expect(report.net_profit).toBe(215000)
  })
})
