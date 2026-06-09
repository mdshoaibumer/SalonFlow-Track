import { apiClient } from './api-client'
import type { LicenseRecord, LicenseStatus, LicenseValidation, LicenseEvent } from '@/types'

export async function getLicenseStatus(): Promise<LicenseStatus> {
  const response = await apiClient.get<LicenseStatus>('/license')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get license status')
  return response.data
}

export async function validateLicense(): Promise<LicenseValidation> {
  const response = await apiClient.post<LicenseValidation>('/license/validate', {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to validate license')
  return response.data
}

export async function activateLicense(licenseKey: string, customerName: string, salonName: string): Promise<LicenseRecord> {
  const response = await apiClient.post<LicenseRecord>('/license/activate', { license_key: licenseKey, customer_name: customerName, salon_name: salonName })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to activate license')
  return response.data
}

export async function renewLicense(licenseKey?: string): Promise<LicenseRecord> {
  const response = await apiClient.post<LicenseRecord>('/license/renew', { license_key: licenseKey || '' })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to renew license')
  return response.data
}

export async function listLicenseEvents(page = 1, perPage = 20): Promise<{ events: LicenseEvent[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const response = await apiClient.get<LicenseEvent[]>(`/license/events?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list events')
  return {
    events: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}
