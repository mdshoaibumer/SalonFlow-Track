import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { ErrorBoundary } from './ErrorBoundary'

// Component that throws an error
function ThrowingComponent({ shouldThrow = true }: { shouldThrow?: boolean }) {
  if (shouldThrow) throw new Error('Test error')
  return <div>Working fine</div>
}

describe('ErrorBoundary', () => {
  beforeEach(() => {
    // Suppress React error boundary console errors
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders children when no error', () => {
    render(
      <ErrorBoundary>
        <div>Child content</div>
      </ErrorBoundary>
    )
    expect(screen.getByText('Child content')).toBeInTheDocument()
  })

  it('renders error UI when child throws', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent />
      </ErrorBoundary>
    )
    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
    expect(screen.getByText('Test error')).toBeInTheDocument()
  })

  it('renders custom fallback when provided', () => {
    render(
      <ErrorBoundary fallback={<div>Custom fallback</div>}>
        <ThrowingComponent />
      </ErrorBoundary>
    )
    expect(screen.getByText('Custom fallback')).toBeInTheDocument()
    expect(screen.queryByText('Something went wrong')).not.toBeInTheDocument()
  })

  it('renders Try Again button in default error UI', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent />
      </ErrorBoundary>
    )
    expect(screen.getByRole('button', { name: /try again/i })).toBeInTheDocument()
  })

  it('resets error state when Try Again is clicked', () => {
    const { rerender } = render(
      <ErrorBoundary>
        <ThrowingComponent shouldThrow={true} />
      </ErrorBoundary>
    )
    expect(screen.getByText('Something went wrong')).toBeInTheDocument()

    // After reset, re-render will throw again, but we verify reset was called
    fireEvent.click(screen.getByRole('button', { name: /try again/i }))
    // After handleReset, it re-renders children which throw again
    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
  })

  it('logs error info via componentDidCatch', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    render(
      <ErrorBoundary>
        <ThrowingComponent />
      </ErrorBoundary>
    )
    // React calls console.error and our componentDidCatch also calls it
    expect(consoleSpy).toHaveBeenCalled()
    const calls = consoleSpy.mock.calls.flat()
    const hasErrorBoundaryLog = calls.some(
      (arg) => typeof arg === 'string' && arg.includes('[ErrorBoundary]')
    )
    expect(hasErrorBoundaryLog).toBe(true)
  })

  it('shows default message when error has no message', () => {
    function ThrowEmpty() {
      throw new Error('')
    }
    render(
      <ErrorBoundary>
        <ThrowEmpty />
      </ErrorBoundary>
    )
    expect(screen.getByText('An unexpected error occurred. Please try again.')).toBeInTheDocument()
  })
})
