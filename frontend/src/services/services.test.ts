import { describe, it, expect, vi } from 'vitest'
import { listServices, getServiceById, createService, updateService, deleteService } from './services'

describe('Services Service', () => {
  it('lists services', async () => {
    const result = await listServices()
    expect(result.services).toHaveLength(1)
    expect(result.services[0].name).toBe('Haircut - Ladies')
    expect(result.meta.total).toBe(1)
  })

  it('lists services with params', async () => {
    const result = await listServices({ search: 'hair', status: 'active', category: 'hair', page: 1, per_page: 10 })
    expect(result.services).toHaveLength(1)
  })

  it('gets service by id', async () => {
    const service = await getServiceById('svc1')
    expect(service.name).toBe('Haircut - Ladies')
  })

  it('creates service', async () => {
    const service = await createService({ name: 'New', category: 'hair', price: 500, duration_minutes: 30 } as any)
    expect(service.id).toBeDefined()
  })

  it('updates service', async () => {
    const service = await updateService('svc1', { name: 'Updated' } as any)
    expect(service.id).toBeDefined()
  })

  it('deletes service', async () => {
    await expect(deleteService('svc1')).resolves.toBeUndefined()
  })

  it('listServices returns empty when API returns undefined services', async () => {
    vi.spyOn(window.go.main.ServiceService, 'ListServices').mockResolvedValueOnce({ services: undefined, total: 0 } as any)
    const r = await listServices()
    expect(r.services).toEqual([])
  })
})
