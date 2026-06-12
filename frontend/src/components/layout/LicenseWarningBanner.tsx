import { useLicenseStatus } from '@/hooks/useLicense'
import { useNavigate } from 'react-router-dom'

export function LicenseWarningBanner() {
  const { data: status } = useLicenseStatus()
  const navigate = useNavigate()

  if (!status) return null

  // No warning for healthy licenses with > 7 days remaining
  if (!status.is_restricted && !status.needs_renewal) return null
  if (status.days_remaining > 7 && !status.is_restricted) return null

  const isRestricted = status.is_restricted
  const isGrace = status.license?.status === 'grace_period'
  const isExpiring = status.days_remaining <= 7 && status.days_remaining > 0

  let bgColor = 'bg-yellow-50 border-yellow-200'
  let textColor = 'text-yellow-800'
  let message = ''

  if (isRestricted) {
    bgColor = 'bg-red-50 border-red-200'
    textColor = 'text-red-800'
    message = 'License expired. Restricted mode active — new transactions are blocked.'
  } else if (isGrace) {
    bgColor = 'bg-orange-50 border-orange-200'
    textColor = 'text-orange-800'
    message = `Grace period: ${status.grace_days_remaining} days remaining before restricted mode.`
  } else if (isExpiring) {
    message = `License expires in ${status.days_remaining} day${status.days_remaining === 1 ? '' : 's'}. Please renew.`
  } else {
    message = `License expires in ${status.days_remaining} days. Consider renewing soon.`
  }

  return (
    <div className={`flex items-center justify-between px-4 py-2 border rounded-md ${bgColor}`}>
      <p className={`text-sm font-medium ${textColor}`}>{message}</p>
      <button
        onClick={() => navigate('/license')}
        className={`text-xs font-medium px-3 py-1 rounded border ${textColor} hover:opacity-80`}
      >
        Manage License
      </button>
    </div>
  )
}
