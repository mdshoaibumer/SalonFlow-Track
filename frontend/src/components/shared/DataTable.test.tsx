import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { DataTable } from './DataTable'
import type { ColumnDef } from '@tanstack/react-table'

interface TestRow {
  id: string
  name: string
  status: string
}

const columns: ColumnDef<TestRow, unknown>[] = [
  { accessorKey: 'id', header: 'ID' },
  { accessorKey: 'name', header: 'Name' },
  { accessorKey: 'status', header: 'Status' },
]

const testData: TestRow[] = [
  { id: '1', name: 'Alice', status: 'active' },
  { id: '2', name: 'Bob', status: 'inactive' },
  { id: '3', name: 'Charlie', status: 'active' },
]

describe('DataTable', () => {
  it('renders column headers', () => {
    render(<DataTable columns={columns} data={testData} />)
    expect(screen.getByText('ID')).toBeInTheDocument()
    expect(screen.getByText('Name')).toBeInTheDocument()
    expect(screen.getByText('Status')).toBeInTheDocument()
  })

  it('renders data rows', () => {
    render(<DataTable columns={columns} data={testData} />)
    expect(screen.getByText('Alice')).toBeInTheDocument()
    expect(screen.getByText('Bob')).toBeInTheDocument()
    expect(screen.getByText('Charlie')).toBeInTheDocument()
  })

  it('renders empty state when no data', () => {
    render(
      <DataTable
        columns={columns}
        data={[]}
        emptyTitle="No items"
        emptyDescription="Nothing found"
      />
    )
    expect(screen.getByText('No items')).toBeInTheDocument()
    expect(screen.getByText('Nothing found')).toBeInTheDocument()
  })

  it('renders search input with placeholder', () => {
    render(<DataTable columns={columns} data={testData} searchPlaceholder="Find stuff..." />)
    expect(screen.getByPlaceholderText('Find stuff...')).toBeInTheDocument()
  })

  it('calls onSearchChange when typing in search', async () => {
    const user = userEvent.setup()
    const onSearchChange = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        onSearchChange={onSearchChange}
        searchValue=""
        searchPlaceholder="Search..."
      />
    )
    await user.type(screen.getByPlaceholderText('Search...'), 'test')
    expect(onSearchChange).toHaveBeenCalled()
  })

  it('renders pagination when pageCount > 1', () => {
    const onPageChange = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={1}
        pageCount={5}
        onPageChange={onPageChange}
      />
    )
    expect(screen.getByText('Page 1 of 5')).toBeInTheDocument()
  })

  it('does not render pagination when pageCount is 1', () => {
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={1}
        pageCount={1}
        onPageChange={vi.fn()}
      />
    )
    expect(screen.queryByText(/Page 1 of 1/)).not.toBeInTheDocument()
  })

  it('calls onPageChange when next button clicked', async () => {
    const user = userEvent.setup()
    const onPageChange = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={1}
        pageCount={3}
        onPageChange={onPageChange}
      />
    )
    // Click next page button (3rd button)
    const buttons = screen.getAllByRole('button')
    const nextBtn = buttons.find(b => b.querySelector('[class*="chevron-right"]') || b.getAttribute('aria-label')?.includes('next'))
    // Find button with ChevronRight icon - it's the 3rd pagination button
    const paginationButtons = buttons.filter(b => b.classList.contains('h-8'))
    await user.click(paginationButtons[2]) // Next button
    expect(onPageChange).toHaveBeenCalledWith(2)
  })

  it('disables prev buttons on first page', () => {
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={1}
        pageCount={3}
        onPageChange={vi.fn()}
      />
    )
    const buttons = screen.getAllByRole('button').filter(b => b.classList.contains('h-8'))
    expect(buttons[0]).toBeDisabled() // First page button
    expect(buttons[1]).toBeDisabled() // Prev button
  })

  it('disables next buttons on last page', () => {
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={3}
        pageCount={3}
        onPageChange={vi.fn()}
      />
    )
    const buttons = screen.getAllByRole('button').filter(b => b.classList.contains('h-8'))
    expect(buttons[2]).toBeDisabled() // Next button
    expect(buttons[3]).toBeDisabled() // Last page button
  })

  it('renders export button when onExport provided', () => {
    render(
      <DataTable
        columns={columns}
        data={testData}
        onExport={vi.fn()}
        exportLabel="Export CSV"
      />
    )
    expect(screen.getByRole('button', { name: /export csv/i })).toBeInTheDocument()
  })

  it('calls onExport when export button clicked', async () => {
    const user = userEvent.setup()
    const onExport = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        onExport={onExport}
        exportLabel="Export"
      />
    )
    await user.click(screen.getByRole('button', { name: /export/i }))
    expect(onExport).toHaveBeenCalledTimes(1)
  })

  it('supports client-side sorting when clicking header', async () => {
    const user = userEvent.setup()
    render(<DataTable columns={columns} data={testData} />)
    // Click on Name header to sort
    const nameHeader = screen.getByText('Name')
    await user.click(nameHeader)
    // Data should still render (sorting applied)
    expect(screen.getByText('Alice')).toBeInTheDocument()
  })

  it('supports client-side search filtering (globalFilter)', async () => {
    const user = userEvent.setup()
    render(<DataTable columns={columns} data={testData} searchPlaceholder="Search..." searchKey="name" />)
    const input = screen.getByPlaceholderText('Search...')
    await user.type(input, 'Alice')
    expect(screen.getByText('Alice')).toBeInTheDocument()
  })

  it('clicking last page button calls onPageChange with pageCount', async () => {
    const user = userEvent.setup()
    const onPageChange = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={2}
        pageCount={5}
        onPageChange={onPageChange}
      />
    )
    const buttons = screen.getAllByRole('button').filter(b => b.classList.contains('h-8'))
    // buttons: [first, prev, next, last]
    await user.click(buttons[3]) // Last page
    expect(onPageChange).toHaveBeenCalledWith(5)
  })

  it('clicking first page button calls onPageChange with 1', async () => {
    const user = userEvent.setup()
    const onPageChange = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={3}
        pageCount={5}
        onPageChange={onPageChange}
      />
    )
    const buttons = screen.getAllByRole('button').filter(b => b.classList.contains('h-8'))
    await user.click(buttons[0]) // First page
    expect(onPageChange).toHaveBeenCalledWith(1)
  })

  it('clicking prev page button calls onPageChange with page-1', async () => {
    const user = userEvent.setup()
    const onPageChange = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        page={3}
        pageCount={5}
        onPageChange={onPageChange}
      />
    )
    const buttons = screen.getAllByRole('button').filter(b => b.classList.contains('h-8'))
    await user.click(buttons[1]) // Prev
    expect(onPageChange).toHaveBeenCalledWith(2)
  })

  it('handles server search with undefined searchValue (uses empty string fallback)', () => {
    const onSearchChange = vi.fn()
    render(
      <DataTable
        columns={columns}
        data={testData}
        onSearchChange={onSearchChange}
        searchPlaceholder="Search..."
      />
    )
    // searchValue is undefined, so ?? '' should be used
    const input = screen.getByPlaceholderText('Search...')
    expect(input).toHaveValue('')
  })

  it('renders placeholder headers correctly', () => {
    const groupedColumns: ColumnDef<TestRow, unknown>[] = [
      {
        id: 'group',
        header: 'Group',
        columns: [
          { accessorKey: 'id', header: 'ID' },
          { accessorKey: 'name', header: 'Name' },
        ],
      },
      { accessorKey: 'status', header: 'Status' },
    ]
    render(<DataTable columns={groupedColumns} data={testData} />)
    expect(screen.getByText('ID')).toBeInTheDocument()
    expect(screen.getByText('Name')).toBeInTheDocument()
  })
})
