import type { Invoice, InvoiceStats, CreateInvoiceInput, Payment, RecordPaymentInput } from '@/types'

export interface ListInvoiceParams {
  page?: number
  per_page?: number
  customer_id?: string
  staff_id?: string
  payment_status?: string
  date_from?: string
  date_to?: string
  search?: string
}

export interface ListInvoiceResponse {
  invoices: Invoice[]
  meta: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
}

export async function listInvoices(params: ListInvoiceParams = {}): Promise<ListInvoiceResponse> {
  const result = await window.go.main.InvoiceService.ListInvoices({
    customer_id: params.customer_id || '',
    staff_id: params.staff_id || '',
    payment_status: params.payment_status || '',
    date_from: params.date_from || '',
    date_to: params.date_to || '',
    search: params.search || '',
    page: params.page || 1,
    per_page: params.per_page || 20,
  })
  return {
    invoices: result.invoices || [],
    meta: {
      page: result.page,
      per_page: result.per_page,
      total: result.total,
      total_pages: result.total_pages,
    },
  }
}

export async function getInvoiceById(id: string): Promise<Invoice> {
  return window.go.main.InvoiceService.GetInvoice(id)
}

export async function createInvoice(input: CreateInvoiceInput): Promise<Invoice> {
  return window.go.main.InvoiceService.CreateInvoice(input)
}

export async function recordPayment(invoiceId: string, input: RecordPaymentInput): Promise<Payment> {
  return window.go.main.InvoiceService.RecordPayment(invoiceId, input)
}

export async function getInvoiceStats(): Promise<InvoiceStats> {
  return window.go.main.InvoiceService.GetInvoiceStats()
}
