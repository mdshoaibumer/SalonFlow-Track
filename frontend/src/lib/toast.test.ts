import { describe, it, expect, vi } from 'vitest'
import { toast } from 'sonner'
import { toastSuccess, toastError, toastPromise, mutationToast } from './toast'

vi.mock('sonner', () => ({
  toast: Object.assign(vi.fn(), {
    success: vi.fn(),
    error: vi.fn(),
    promise: vi.fn(),
  }),
}))

describe('toast utilities', () => {
  it('toastSuccess shows success toast with duration', () => {
    toastSuccess('Item created')
    expect(toast.success).toHaveBeenCalledWith('Item created', { duration: 2500 })
  })

  it('toastError shows error toast with message and detail', () => {
    toastError('Failed', 'Something went wrong')
    expect(toast.error).toHaveBeenCalledWith('Failed', {
      description: 'Something went wrong',
      duration: 4000,
    })
  })

  it('toastError shows error toast without detail', () => {
    toastError('Failed')
    expect(toast.error).toHaveBeenCalledWith('Failed', {
      description: undefined,
      duration: 4000,
    })
  })

  it('toastPromise wraps a promise with loading/success/error messages', () => {
    const p = Promise.resolve('data')
    const msgs = { loading: 'Loading...', success: 'Done!', error: 'Failed!' }
    toastPromise(p, msgs)
    expect(toast.promise).toHaveBeenCalledWith(p, msgs)
  })

  it('mutationToast returns onSuccess and onError callbacks', () => {
    const result = mutationToast('Created!', 'Create failed')
    expect(result).toHaveProperty('onSuccess')
    expect(result).toHaveProperty('onError')
  })

  it('mutationToast onSuccess calls toastSuccess', () => {
    const result = mutationToast('Created!')
    result.onSuccess()
    expect(toast.success).toHaveBeenCalledWith('Created!', { duration: 2500 })
  })

  it('mutationToast onError calls toastError with error message', () => {
    const result = mutationToast('Created!', 'Create failed')
    result.onError(new Error('Network error'))
    expect(toast.error).toHaveBeenCalledWith('Create failed', {
      description: 'Network error',
      duration: 4000,
    })
  })

  it('mutationToast onError uses default message when no errorMsg provided', () => {
    const result = mutationToast('Created!')
    result.onError(new Error('Oops'))
    expect(toast.error).toHaveBeenCalledWith('Operation failed', {
      description: 'Oops',
      duration: 4000,
    })
  })
})
