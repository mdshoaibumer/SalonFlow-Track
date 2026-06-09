import { apiClient } from './api-client'
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
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.customer_id) query.set('customer_id', params.customer_id)
  if (params.staff_id) query.set('staff_id', params.staff_id)
  if (params.payment_status) query.set('payment_status', params.payment_status)
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  if (params.search) query.set('search', params.search)

  const qs = query.toString()
  const path = qs ? `/invoices?${qs}` : '/invoices'
  const response = await apiClient.get<Invoice[]>(path)

  if (!response.success) {
    throw new Error(response.error?.message || 'Failed to fetch invoices')
  }

  return {
    invoices: response.data || [],
    meta: response.meta as { page: number; per_page: number; total: number; total_pages: number } || { page: 1, per_page: 20, total: 0, total_pages: 0 },
  }
}

export async function getInvoiceById(id: string): Promise<Invoice> {
  const response = await apiClient.get<Invoice>(`/invoices/${id}`)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch invoice')
  }
  return response.data
}

export async function createInvoice(input: CreateInvoiceInput): Promise<Invoice> {
  const response = await apiClient.post<Invoice>('/invoices', input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to create invoice')
  }
  return response.data
}

export async function recordPayment(invoiceId: string, input: RecordPaymentInput): Promise<Payment> {
  const response = await apiClient.post<Payment>(`/invoices/${invoiceId}/payment`, input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to record payment')
  }
  return response.data
}

export async function getInvoiceStats(): Promise<InvoiceStats> {
  const response = await apiClient.get<InvoiceStats>('/invoices/stats')
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch invoice stats')
  }
  return response.data
}
