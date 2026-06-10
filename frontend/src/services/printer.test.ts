import { describe, it, expect, vi } from 'vitest'
import { getPrinterSettings, savePrinterSettings, printInvoice, printReceipt, printTest, listPrintJobs, getPrintJob } from './printer'

describe('Printer Service', () => {
  it('gets printer settings', async () => {
    const settings = await getPrinterSettings()
    expect(settings.printer_name).toBe('POS-80')
  })

  it('saves printer settings', async () => {
    await expect(savePrinterSettings({ printer_name: 'POS-80', auto_print: true })).resolves.toBeUndefined()
  })

  it('prints invoice', async () => {
    const result = await printInvoice({} as any)
    expect(result.job.id).toBe('job1')
    expect(result.html).toContain('Invoice')
  })

  it('prints receipt', async () => {
    const result = await printReceipt({} as any)
    expect(result.job.id).toBe('job2')
  })

  it('prints test', async () => {
    const result = await printTest()
    expect(result.job.id).toBe('job3')
  })

  it('lists print jobs', async () => {
    const result = await listPrintJobs()
    expect(result.jobs).toHaveLength(1)
    expect(result.meta.total).toBe(1)
  })

  it('lists print jobs with pagination', async () => {
    const result = await listPrintJobs(2, 10)
    expect(result.jobs).toHaveLength(1)
  })

  it('gets print job', async () => {
    const job = await getPrintJob('job1')
    expect(job.id).toBe('job1')
  })

  it('listPrintJobs returns empty when API returns undefined', async () => {
    vi.spyOn(window.go.main.PrinterService, 'ListPrintJobs').mockResolvedValueOnce([undefined as any, 0])
    const r = await listPrintJobs()
    expect(r.jobs).toEqual([])
  })
})
