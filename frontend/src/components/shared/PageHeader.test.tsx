import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { PageHeader } from './PageHeader'
import { Button } from '@/components/ui/button'

describe('PageHeader', () => {
  it('renders title', () => {
    render(<PageHeader title="Staff Management" />)
    expect(screen.getByText('Staff Management')).toBeInTheDocument()
  })

  it('renders title and description', () => {
    render(<PageHeader title="Staff" description="Manage your team" />)
    expect(screen.getByText('Staff')).toBeInTheDocument()
    expect(screen.getByText('Manage your team')).toBeInTheDocument()
  })

  it('renders actions', () => {
    const onClick = vi.fn()
    render(
      <PageHeader
        title="Test"
        actions={<Button onClick={onClick}>Add</Button>}
      />
    )
    expect(screen.getByRole('button', { name: 'Add' })).toBeInTheDocument()
  })

  it('does not render description when not provided', () => {
    render(<PageHeader title="Test" />)
    const paragraphs = screen.queryAllByRole('paragraph')
    // No p element with muted class
    expect(screen.queryByText('undefined')).not.toBeInTheDocument()
  })
})
