import { describe, it, expect, vi } from 'vitest'
import { listInvoices, getInvoiceById, createInvoice, recordPayment, getInvoiceStats } from './invoices'

describe('Invoices Service', () => {
  it('lists invoices', async () => {
    const result = await listInvoices()
    expect(result.invoices).toHaveLength(1)
    expect(result.meta.total).toBe(1)
  })

  it('lists invoices with params', async () => {
    const result = await listInvoices({ customer_id: 'c1', staff_id: 's1', payment_status: 'paid', date_from: '2024-01-01', date_to: '2024-12-31', search: 'INV', page: 1, per_page: 10 })
    expect(result.invoices).toHaveLength(1)
  })

  it('gets invoice by id', async () => {
    const invoice = await getInvoiceById('inv1')
    expect(invoice.invoice_number).toBe('INV-001')
  })

  it('creates invoice', async () => {
    const invoice = await createInvoice({ customer_id: 'c1', items: [] } as any)
    expect(invoice.id).toBeDefined()
  })

  it('records payment', async () => {
    const payment = await recordPayment('inv1', { amount: 2242, method: 'upi' } as any)
    expect(payment.id).toBe('pay1')
  })

  it('gets invoice stats', async () => {
    const stats = await getInvoiceStats()
    expect(stats.today_revenue).toBe(12500)
  })

  it('listInvoices returns empty when API returns undefined invoices', async () => {
    vi.spyOn(window.go.main.InvoiceService, 'ListInvoices').mockResolvedValueOnce({ invoices: undefined, total: 0, total_amount: 0 } as any)
    const r = await listInvoices()
    expect(r.invoices).toEqual([])
  })
})
