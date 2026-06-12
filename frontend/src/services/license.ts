import type {
  License,
  LicenseStatus,
  LicenseValidation,
  LicenseNotification,
} from '@/types/license'

const svc = () => window.go.main.LicenseService

export async function getLicenseStatus(): Promise<LicenseStatus> {
  return svc().GetStatus()
}

export async function validateLicense(): Promise<LicenseValidation> {
  return svc().Validate()
}

export async function activateLicense(key: string, customerName: string, salonName: string): Promise<License> {
  return svc().Activate(key, customerName, salonName)
}

export async function importLicenseFile(fileData: number[]): Promise<License> {
  return svc().ImportLicenseFile(fileData)
}

export async function exportLicenseFile(): Promise<number[]> {
  return svc().ExportLicenseFile()
}

export async function renewLicense(key?: string): Promise<License> {
  return svc().Renew(key || '')
}

export async function getDeviceID(): Promise<string> {
  return svc().GetDeviceID()
}

export async function getNotifications(unreadOnly = true): Promise<LicenseNotification[]> {
  return svc().GetNotifications(unreadOnly) || []
}

export async function markNotificationRead(id: string): Promise<void> {
  return svc().MarkNotificationRead(id)
}

export async function dismissNotification(id: string): Promise<void> {
  return svc().DismissNotification(id)
}

export async function isOperationAllowed(operation: string): Promise<void> {
  return svc().IsOperationAllowed(operation)
}

export async function listLicenseEvents(page = 1, perPage = 20) {
  const [events, total] = await svc().ListEvents(page, perPage)
  const total_pages = Math.ceil(total / perPage)
  return { events: events || [], meta: { page, per_page: perPage, total, total_pages } }
}
