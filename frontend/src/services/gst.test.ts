import { describe, it, expect } from 'vitest'
import { getGSTSettings, saveGSTSettings, listTaxRates, createTaxRate, updateTaxRate, deleteTaxRate, getGSTReport } from './gst'

describe('GST Service', () => {
  it('gets GST settings', async () => {
    const settings = await getGSTSettings()
    expect(settings.gstin).toBe('27AABCU9603R1ZM')
  })

  it('saves GST settings', async () => {
    await expect(saveGSTSettings({ is_enabled: true })).resolves.toBeUndefined()
  })

  it('lists tax rates', async () => {
    const rates = await listTaxRates()
    expect(rates).toHaveLength(1)
    expect(rates[0].rate).toBe(18)
  })

  it('lists tax rates by category', async () => {
    const rates = await listTaxRates('service')
    expect(rates).toHaveLength(1)
  })

  it('creates tax rate', async () => {
    await expect(createTaxRate({ name: 'GST 5%', rate: 5 })).resolves.toBeUndefined()
  })

  it('updates tax rate', async () => {
    await expect(updateTaxRate('tax1', { name: 'GST 18%', rate: 18 })).resolves.toBeUndefined()
  })

  it('deletes tax rate', async () => {
    await expect(deleteTaxRate('tax1')).resolves.toBeUndefined()
  })

  it('gets GST report', async () => {
    const report = await getGSTReport('2024-01-01', '2024-12-31')
    expect(report.total_taxable).toBe(250000)
  })

  it('gets GST report with period', async () => {
    const report = await getGSTReport('2024-01-01', '2024-12-31', 'monthly')
    expect(report.total_taxable).toBe(250000)
  })
})
