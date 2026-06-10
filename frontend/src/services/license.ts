import type { LicenseStatus, LicenseValidation, License } from '@/types'

export async function getLicenseStatus(): Promise<LicenseStatus> {
  return window.go.main.LicenseService.GetStatus()
}

export async function validateLicense(): Promise<LicenseValidation> {
  return window.go.main.LicenseService.Validate()
}

export async function activateLicense(key: string, customerName: string, salonName: string): Promise<License> {
  return window.go.main.LicenseService.Activate(key, customerName, salonName)
}

export async function renewLicense(key?: string): Promise<License> {
  return window.go.main.LicenseService.Renew(key || '')
}

export async function listLicenseEvents(page = 1, perPage = 20) {
  const [events, total] = await window.go.main.LicenseService.ListEvents(page, perPage)
  const total_pages = Math.ceil(total / perPage)
  return { events: events || [], meta: { page, per_page: perPage, total, total_pages } }
}
