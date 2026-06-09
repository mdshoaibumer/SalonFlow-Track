import { apiClient } from './api-client'
import type { GSTSettings, TaxRate, InvoiceTaxLine, GSTReport } from '@/types'

export async function getGSTSettings(): Promise<GSTSettings> {
  const response = await apiClient.get<GSTSettings>('/gst/settings')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get GST settings')
  return response.data
}

export async function saveGSTSettings(settings: Partial<GSTSettings>): Promise<GSTSettings> {
  const response = await apiClient.post<GSTSettings>('/gst/settings', settings)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to save GST settings')
  return response.data
}

export async function listTaxRates(category?: string): Promise<TaxRate[]> {
  const params = category ? `?category=${category}` : ''
  const response = await apiClient.get<TaxRate[]>(`/gst/tax-rates${params}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list tax rates')
  return response.data || []
}

export async function createTaxRate(rate: Partial<TaxRate>): Promise<TaxRate> {
  const response = await apiClient.post<TaxRate>('/gst/tax-rates', rate)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create tax rate')
  return response.data
}

export async function updateTaxRate(id: string, rate: Partial<TaxRate>): Promise<TaxRate> {
  const response = await apiClient.put<TaxRate>(`/gst/tax-rates/${id}`, rate)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update tax rate')
  return response.data
}

export async function deleteTaxRate(id: string): Promise<void> {
  const response = await apiClient.delete(`/gst/tax-rates/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete tax rate')
}

export async function getInvoiceTaxLines(invoiceId: string): Promise<InvoiceTaxLine[]> {
  const response = await apiClient.get<InvoiceTaxLine[]>(`/gst/invoice/${invoiceId}/tax`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to get tax lines')
  return response.data || []
}

export async function getGSTReport(startDate: string, endDate: string, period = 'daily'): Promise<GSTReport> {
  const response = await apiClient.get<GSTReport>(`/gst/reports?period=${period}&start_date=${startDate}&end_date=${endDate}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get GST report')
  return response.data
}
