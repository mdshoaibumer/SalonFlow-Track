import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { EmptyState } from './EmptyState'
import { Users } from 'lucide-react'

describe('EmptyState', () => {
  it('renders title and description', () => {
    render(<EmptyState title="No results" description="Try a different search" />)
    expect(screen.getByText('No results')).toBeInTheDocument()
    expect(screen.getByText('Try a different search')).toBeInTheDocument()
  })

  it('renders action button when provided', () => {
    const onClick = vi.fn()
    render(
      <EmptyState
        title="Empty"
        description="Nothing here"
        action={{ label: 'Create', onClick }}
      />
    )
    expect(screen.getByRole('button', { name: 'Create' })).toBeInTheDocument()
  })

  it('calls action onClick when button clicked', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(
      <EmptyState
        title="Empty"
        description="Nothing here"
        action={{ label: 'Add Item', onClick }}
      />
    )
    await user.click(screen.getByRole('button', { name: 'Add Item' }))
    expect(onClick).toHaveBeenCalledTimes(1)
  })

  it('renders custom icon', () => {
    render(<EmptyState icon={Users} title="No staff" description="Add some" />)
    expect(screen.getByText('No staff')).toBeInTheDocument()
  })

  it('does not render button when action not provided', () => {
    render(<EmptyState title="Empty" description="Nothing" />)
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })
})
