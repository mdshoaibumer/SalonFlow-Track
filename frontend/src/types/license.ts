export interface License {
  id: string
  license_key: string
  customer_name: string
  salon_name: string
  device_id: string
  issued_date: string
  expiry_date: string
  grace_until: string
  status: LicenseStatusType
  signature: string
  last_validation: string
  last_verified_at: string
  created_at: string
  updated_at: string
}

export type LicenseStatusType = 'active' | 'grace_period' | 'expired' | 'suspended'

export interface LicenseStatus {
  license: License | null
  days_remaining: number
  grace_days_remaining: number
  is_restricted: boolean
  needs_renewal: boolean
}

export interface LicenseValidation {
  valid: boolean
  status: LicenseStatusType
  days_remaining: number
  is_restricted: boolean
  message: string
}

export interface LicenseEvent {
  id: string
  license_id: string
  event_type: string
  event_date: string
  notes: string
  created_at: string
}

export interface LicenseNotification {
  id: string
  license_id: string
  notification_type: string
  title: string
  message: string
  is_read: boolean
  is_dismissed: boolean
  created_at: string
}

export type NotificationType =
  | '7_days_remaining'
  | '3_days_remaining'
  | '1_day_remaining'
  | 'expired'
  | 'grace_period_remaining'
