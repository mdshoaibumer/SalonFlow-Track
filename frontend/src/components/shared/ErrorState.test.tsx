import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ErrorState } from './ErrorState'

describe('ErrorState', () => {
  it('renders default title and message', () => {
    render(<ErrorState />)
    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
    expect(screen.getByText('Failed to load data. Please check your connection and try again.')).toBeInTheDocument()
  })

  it('renders custom title and message', () => {
    render(<ErrorState title="Network Error" message="Check your wifi" />)
    expect(screen.getByText('Network Error')).toBeInTheDocument()
    expect(screen.getByText('Check your wifi')).toBeInTheDocument()
  })

  it('renders retry button when onRetry provided', () => {
    const onRetry = vi.fn()
    render(<ErrorState onRetry={onRetry} />)
    expect(screen.getByRole('button', { name: /try again/i })).toBeInTheDocument()
  })

  it('does not render retry button without onRetry', () => {
    render(<ErrorState />)
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })

  it('calls onRetry when retry button clicked', async () => {
    const user = userEvent.setup()
    const onRetry = vi.fn()
    render(<ErrorState onRetry={onRetry} />)
    await user.click(screen.getByRole('button', { name: /try again/i }))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })
})
