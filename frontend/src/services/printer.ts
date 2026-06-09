import { apiClient } from './api-client'
import type { PrinterSettings, PrintJob, ReceiptData } from '@/types'

export async function getPrinterSettings(): Promise<PrinterSettings> {
  const response = await apiClient.get<PrinterSettings>('/print/settings')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get printer settings')
  return response.data
}

export async function savePrinterSettings(settings: Partial<PrinterSettings>): Promise<PrinterSettings> {
  const response = await apiClient.post<PrinterSettings>('/print/settings', settings)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to save printer settings')
  return response.data
}

export async function printInvoice(data: ReceiptData): Promise<{ job: PrintJob; receipt: string }> {
  const response = await apiClient.post<{ job: PrintJob; receipt: string }>('/print/invoice', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to print invoice')
  return response.data
}

export async function printReceipt(data: ReceiptData): Promise<{ job: PrintJob; escpos: string }> {
  const response = await apiClient.post<{ job: PrintJob; escpos: string }>('/print/receipt', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to print receipt')
  return response.data
}

export async function printTest(): Promise<{ job: PrintJob; escpos: string }> {
  const response = await apiClient.post<{ job: PrintJob; escpos: string }>('/print/test', {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to print test page')
  return response.data
}

export async function listPrintJobs(page = 1, perPage = 20): Promise<{ jobs: PrintJob[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const response = await apiClient.get<PrintJob[]>(`/print/history?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list print jobs')
  return {
    jobs: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}
