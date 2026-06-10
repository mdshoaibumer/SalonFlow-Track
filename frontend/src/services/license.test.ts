import { describe, it, expect, vi } from 'vitest'
import { getLicenseStatus, validateLicense, activateLicense, renewLicense, listLicenseEvents } from './license'

describe('License Service', () => {
  it('gets license status', async () => {
    const status = await getLicenseStatus()
    expect(status.status).toBe('active')
  })

  it('validates license', async () => {
    const result = await validateLicense()
    expect(result.valid).toBe(true)
  })

  it('activates license', async () => {
    const license = await activateLicense('XXXX-XXXX', 'John', 'My Salon')
    expect(license.id).toBe('lic1')
  })

  it('renews license', async () => {
    const license = await renewLicense('XXXX-XXXX')
    expect(license.status).toBe('active')
  })

  it('renews license without key', async () => {
    const license = await renewLicense()
    expect(license.status).toBe('active')
  })

  it('lists license events', async () => {
    const result = await listLicenseEvents()
    expect(result.events).toHaveLength(1)
    expect(result.meta.total).toBe(1)
  })

  it('lists license events with pagination', async () => {
    const result = await listLicenseEvents(2, 10)
    expect(result.events).toHaveLength(1)
  })

  it('listLicenseEvents returns empty when API returns undefined', async () => {
    vi.spyOn(window.go.main.LicenseService, 'ListEvents').mockResolvedValueOnce([undefined as any, 0])
    const r = await listLicenseEvents()
    expect(r.events).toEqual([])
  })
})
