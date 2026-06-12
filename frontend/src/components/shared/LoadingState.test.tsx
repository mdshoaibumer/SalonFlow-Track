import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { LoadingState } from './LoadingState'

describe('LoadingState', () => {
  it('renders table variant by default', () => {
    const { container } = render(<LoadingState />)
    // Should have skeleton elements (shimmer divs)
    const skeletons = container.querySelectorAll('[class*="animate-"]')
    expect(skeletons.length).toBeGreaterThan(0)
  })

  it('renders correct number of rows for table', () => {
    const { container } = render(<LoadingState rows={3} variant="table" />)
    // Each row has multiple skeletons in a flex container
    const rows = container.querySelectorAll('.flex.items-center.gap-4')
    expect(rows.length).toBe(3)
  })

  it('renders cards variant with 4 cards', () => {
    const { container } = render(<LoadingState variant="cards" />)
    const cards = container.querySelectorAll('.surface-base')
    expect(cards.length).toBe(4)
  })

  it('renders page variant with header + cards + table', () => {
    const { container } = render(<LoadingState variant="page" rows={3} />)
    // Page variant has all sections with skeleton elements
    const skeletons = container.querySelectorAll('[class*="animate-"]')
    expect(skeletons.length).toBeGreaterThan(10)
  })

  it('defaults to 5 rows', () => {
    const { container } = render(<LoadingState variant="table" />)
    const rows = container.querySelectorAll('.flex.items-center.gap-4')
    expect(rows.length).toBe(5)
  })

  it('renders chart variant with bar skeletons', () => {
    const { container } = render(<LoadingState variant="chart" />)
    const surface = container.querySelector('.surface-base')
    expect(surface).toBeInTheDocument()
    // Should have 12 bar skeletons
    const bars = container.querySelectorAll('.flex-1.rounded-t-sm')
    expect(bars.length).toBe(12)
  })
})
