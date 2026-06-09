import { useState } from 'react'
import { useLicenseStatus, useLicenseEvents, useActivateLicense, useRenewLicense, useValidateLicense } from '@/hooks/useLicense'
import { Shield, ShieldCheck, ShieldAlert, ShieldX, Clock, RefreshCw, Key } from 'lucide-react'

function formatDate(dateStr: string): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('en-IN', { dateStyle: 'medium', timeStyle: 'short' })
}

function StatusIcon({ status }: { status: string }) {
  switch (status) {
    case 'active':
      return <ShieldCheck className="h-8 w-8 text-green-600" />
    case 'grace_period':
      return <ShieldAlert className="h-8 w-8 text-yellow-600" />
    case 'expired':
      return <ShieldX className="h-8 w-8 text-red-600" />
    case 'suspended':
      return <ShieldX className="h-8 w-8 text-red-800" />
    default:
      return <Shield className="h-8 w-8 text-gray-400" />
  }
}

function StatusBadge({ status }: { status: string }) {
  const styles: Record<string, string> = {
    active: 'bg-green-100 text-green-800',
    grace_period: 'bg-yellow-100 text-yellow-800',
    expired: 'bg-red-100 text-red-800',
    suspended: 'bg-red-100 text-red-800',
  }
  return (
    <span className={`inline-flex items-center rounded-full px-3 py-1 text-sm font-medium ${styles[status] || 'bg-gray-100 text-gray-800'}`}>
      {status.replace('_', ' ')}
    </span>
  )
}

export function LicensePage() {
  const [showActivate, setShowActivate] = useState(false)
  const [showRenew, setShowRenew] = useState(false)
  const [activateForm, setActivateForm] = useState({ licenseKey: '', customerName: '', salonName: '' })
  const [renewKey, setRenewKey] = useState('')
  const [eventsPage, setEventsPage] = useState(1)

  const { data: status, isLoading } = useLicenseStatus()
  const { data: eventsData } = useLicenseEvents(eventsPage)
  const activateMut = useActivateLicense()
  const renewMut = useRenewLicense()
  const validateMut = useValidateLicense()

  const handleActivate = (e: React.FormEvent) => {
    e.preventDefault()
    activateMut.mutate(activateForm, { onSuccess: () => setShowActivate(false) })
  }

  const handleRenew = (e: React.FormEvent) => {
    e.preventDefault()
    renewMut.mutate(renewKey || undefined, { onSuccess: () => setShowRenew(false) })
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">License & Subscription</h1>
            <p className="text-muted-foreground">Manage your SalonFlow Track license</p>
          </div>
        </div>
        <div className="flex items-center justify-center p-12 text-muted-foreground">Loading license info...</div>
      </div>
    )
  }

  const lic = status?.license
  const noLicense = !lic

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">License & Subscription</h1>
          <p className="text-muted-foreground">Manage your SalonFlow Track license</p>
        </div>
        <button
          onClick={() => validateMut.mutate()}
          disabled={validateMut.isPending}
          className="inline-flex items-center gap-2 rounded-lg border px-4 py-2 text-sm font-medium hover:bg-muted"
        >
          <RefreshCw className={`h-4 w-4 ${validateMut.isPending ? 'animate-spin' : ''}`} />
          Validate
        </button>
      </div>

      {/* Restricted Mode Banner */}
      {status?.is_restricted && (
        <div className="rounded-lg border border-red-200 bg-red-50 p-4">
          <div className="flex items-center gap-3">
            <ShieldX className="h-5 w-5 text-red-600" />
            <div>
              <p className="font-medium text-red-800">Restricted Mode Active</p>
              <p className="text-sm text-red-600">New transactions are disabled. Please renew your license to restore full access.</p>
            </div>
          </div>
        </div>
      )}

      {/* License Status Card */}
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-start gap-6">
          <StatusIcon status={lic?.status || 'expired'} />
          <div className="flex-1 space-y-4">
            {noLicense ? (
              <div>
                <h2 className="text-xl font-semibold">No License Active</h2>
                <p className="text-muted-foreground">Activate a license to get started.</p>
              </div>
            ) : (
              <>
                <div className="flex items-center gap-3">
                  <h2 className="text-xl font-semibold">{lic.salon_name}</h2>
                  <StatusBadge status={lic.status} />
                </div>
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                  <div>
                    <p className="text-sm text-muted-foreground">License Key</p>
                    <p className="font-mono text-sm font-medium">{lic.license_key}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Expiry Date</p>
                    <p className="font-medium">{lic.expiry_date}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Days Remaining</p>
                    <p className={`text-2xl font-bold ${(status?.days_remaining ?? 0) <= 7 ? 'text-red-600' : (status?.days_remaining ?? 0) <= 30 ? 'text-yellow-600' : 'text-green-600'}`}>
                      {status?.days_remaining ?? 0}
                    </p>
                  </div>
                  {lic.status === 'grace_period' && (
                    <div>
                      <p className="text-sm text-muted-foreground">Grace Days Left</p>
                      <p className="text-2xl font-bold text-yellow-600">{status?.grace_days_remaining ?? 0}</p>
                    </div>
                  )}
                </div>
                <div className="text-sm text-muted-foreground">
                  <span className="mr-4">Customer: {lic.customer_name}</span>
                  <span>Last Validated: {formatDate(lic.last_validation)}</span>
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="flex gap-3">
        {noLicense && (
          <button onClick={() => setShowActivate(true)} className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
            <Key className="h-4 w-4" /> Activate License
          </button>
        )}
        {!noLicense && (
          <button onClick={() => setShowRenew(true)} className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
            <RefreshCw className="h-4 w-4" /> Renew License
          </button>
        )}
      </div>

      {/* Activate Form */}
      {showActivate && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="mb-4 text-lg font-semibold">Activate License</h3>
          <form onSubmit={handleActivate} className="space-y-4">
            <div>
              <label className="block text-sm font-medium">License Key</label>
              <input type="text" value={activateForm.licenseKey} onChange={(e) => setActivateForm((f) => ({ ...f, licenseKey: e.target.value }))} placeholder="SALONFLOW-XXXX-XXXX-XXXX" className="mt-1 block w-full rounded-lg border bg-background px-3 py-2 text-sm" required />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <label className="block text-sm font-medium">Customer Name</label>
                <input type="text" value={activateForm.customerName} onChange={(e) => setActivateForm((f) => ({ ...f, customerName: e.target.value }))} className="mt-1 block w-full rounded-lg border bg-background px-3 py-2 text-sm" required />
              </div>
              <div>
                <label className="block text-sm font-medium">Salon Name</label>
                <input type="text" value={activateForm.salonName} onChange={(e) => setActivateForm((f) => ({ ...f, salonName: e.target.value }))} className="mt-1 block w-full rounded-lg border bg-background px-3 py-2 text-sm" required />
              </div>
            </div>
            <div className="flex gap-2">
              <button type="submit" disabled={activateMut.isPending} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
                {activateMut.isPending ? 'Activating...' : 'Activate'}
              </button>
              <button type="button" onClick={() => setShowActivate(false)} className="rounded-lg border px-4 py-2 text-sm font-medium hover:bg-muted">Cancel</button>
            </div>
          </form>
        </div>
      )}

      {/* Renew Form */}
      {showRenew && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="mb-4 text-lg font-semibold">Renew License</h3>
          <form onSubmit={handleRenew} className="space-y-4">
            <div>
              <label className="block text-sm font-medium">Renewal Key (optional)</label>
              <input type="text" value={renewKey} onChange={(e) => setRenewKey(e.target.value)} placeholder="Leave blank to use existing key" className="mt-1 block w-full rounded-lg border bg-background px-3 py-2 text-sm" />
            </div>
            <div className="flex gap-2">
              <button type="submit" disabled={renewMut.isPending} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
                {renewMut.isPending ? 'Renewing...' : 'Renew (₹1500/month)'}
              </button>
              <button type="button" onClick={() => setShowRenew(false)} className="rounded-lg border px-4 py-2 text-sm font-medium hover:bg-muted">Cancel</button>
            </div>
          </form>
        </div>
      )}

      {/* Audit Events */}
      <div>
        <h3 className="mb-3 text-lg font-semibold">License Events</h3>
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead className="border-b bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">Event</th>
                <th className="px-4 py-3 text-left font-medium">Date</th>
                <th className="px-4 py-3 text-left font-medium">Notes</th>
              </tr>
            </thead>
            <tbody>
              {!eventsData?.events?.length ? (
                <tr><td colSpan={3} className="px-4 py-8 text-center text-muted-foreground">No events yet.</td></tr>
              ) : (
                eventsData.events.map((ev) => (
                  <tr key={ev.id} className="border-b last:border-0">
                    <td className="px-4 py-3 capitalize">{ev.event_type.replace('_', ' ')}</td>
                    <td className="px-4 py-3">{formatDate(ev.event_date)}</td>
                    <td className="px-4 py-3 text-muted-foreground">{ev.notes || '—'}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
          {eventsData && eventsData.meta.total_pages > 1 && (
            <div className="flex items-center justify-between border-t px-4 py-3">
              <span className="text-sm text-muted-foreground">Page {eventsData.meta.page} of {eventsData.meta.total_pages}</span>
              <div className="flex gap-2">
                <button onClick={() => setEventsPage((p) => Math.max(1, p - 1))} disabled={eventsPage === 1} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Prev</button>
                <button onClick={() => setEventsPage((p) => p + 1)} disabled={eventsPage >= eventsData.meta.total_pages} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Next</button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Renewal Notification */}
      {status?.needs_renewal && !status?.is_restricted && (
        <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-4">
          <div className="flex items-center gap-3">
            <Clock className="h-5 w-5 text-yellow-600" />
            <div>
              <p className="font-medium text-yellow-800">Renewal Reminder</p>
              <p className="text-sm text-yellow-600">
                {(status?.days_remaining ?? 0) > 0
                  ? `Your license expires in ${status?.days_remaining} days. Please renew at ₹1500/month.`
                  : `Your license is in grace period. ${status?.grace_days_remaining} days remaining before restricted mode.`}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
