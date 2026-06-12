/**
 * SalonFlow Toast System
 * =======================
 * Lightweight wrapper over sonner for consistent mutation feedback.
 */

import { toast } from 'sonner'

export { toast }

/** Show success toast after mutation */
export function toastSuccess(message: string) {
  toast.success(message, {
    duration: 2500,
  })
}

/** Show error toast on mutation failure */
export function toastError(message: string, detail?: string) {
  toast.error(message, {
    description: detail,
    duration: 4000,
  })
}

/** Show promise-based toast (loading → success/error) */
export function toastPromise<T>(
  promise: Promise<T>,
  messages: { loading: string; success: string; error: string }
) {
  return toast.promise(promise, messages)
}

/**
 * Helper for react-query mutation callbacks.
 * Usage: mutate(data, mutationToast('Staff created', 'Failed to create staff'))
 */
export function mutationToast(successMsg: string, errorMsg?: string) {
  return {
    onSuccess: () => toastSuccess(successMsg),
    onError: (err: Error) => toastError(errorMsg ?? 'Operation failed', err.message),
  }
}
