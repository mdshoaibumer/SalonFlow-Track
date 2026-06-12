import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { KPICard } from './KPICard'
import { Users } from 'lucide-react'

describe('KPICard', () => {
  it('renders title and value', () => {
    render(<KPICard title="Total Staff" value={42} icon={Users} />)
    expect(screen.getByText('Total Staff')).toBeInTheDocument()
    expect(screen.getByText('42')).toBeInTheDocument()
  })

  it('renders string value', () => {
    render(<KPICard title="Revenue" value="₹12,500" icon={Users} />)
    expect(screen.getByText('₹12,500')).toBeInTheDocument()
  })

  it('renders description when provided', () => {
    render(<KPICard title="Performer" value="Priya" icon={Users} description="₹50,000 this month" />)
    expect(screen.getByText('₹50,000 this month')).toBeInTheDocument()
  })

  it('renders positive trend', () => {
    render(<KPICard title="Revenue" value="₹12,500" icon={Users} trend={{ value: 15, label: 'vs last week' }} />)
    expect(screen.getByText(/15/)).toBeInTheDocument()
    expect(screen.getByText(/vs last week/)).toBeInTheDocument()
  })

  it('renders negative trend', () => {
    render(<KPICard title="Revenue" value="₹8,000" icon={Users} trend={{ value: -10, label: 'vs yesterday' }} />)
    expect(screen.getByText(/10/)).toBeInTheDocument()
    expect(screen.getByText(/vs yesterday/)).toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = render(<KPICard title="Test" value={0} icon={Users} className="border-red-500" />)
    expect(container.firstChild).toHaveClass('border-red-500')
  })
})
