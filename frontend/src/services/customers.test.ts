import { describe, it, expect, vi } from 'vitest'
import { listCustomers, getCustomerById, createCustomer, updateCustomer, deleteCustomer, getCustomerStats } from './customers'

describe('Customers Service', () => {
  it('lists customers', async () => {
    const result = await listCustomers()
    expect(result.customers).toHaveLength(1)
    expect(result.customers[0].full_name).toBe('Anjali Desai')
    expect(result.meta.total).toBe(1)
  })

  it('lists customers with params', async () => {
    const result = await listCustomers({ search: 'Anjali', status: 'active', page: 1, per_page: 10 })
    expect(result.customers).toHaveLength(1)
  })

  it('gets customer by id', async () => {
    const customer = await getCustomerById('01912345-6789-7abc-def0-123456789003')
    expect(customer.full_name).toBe('Anjali Desai')
  })

  it('creates customer', async () => {
    const customer = await createCustomer({ full_name: 'New', phone: '9876543212', gender: 'female' } as any)
    expect(customer.id).toBeDefined()
  })

  it('updates customer', async () => {
    const customer = await updateCustomer('cust1', { full_name: 'Updated' } as any)
    expect(customer.id).toBeDefined()
  })

  it('deletes customer', async () => {
    await expect(deleteCustomer('cust1')).resolves.toBeUndefined()
  })

  it('gets customer stats', async () => {
    const stats = await getCustomerStats()
    expect(stats.total).toBe(50)
    expect(stats.active).toBe(45)
  })

  it('listCustomers returns empty when API returns undefined customers', async () => {
    vi.spyOn(window.go.main.CustomerService, 'ListCustomers').mockResolvedValueOnce({ customers: undefined, total: 0 } as any)
    const r = await listCustomers()
    expect(r.customers).toEqual([])
  })
})
