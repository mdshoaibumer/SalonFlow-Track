import { describe, it, expect, vi } from 'vitest'
import { createAppointment, updateAppointment, updateAppointmentStatus, deleteAppointment, getAppointment, listAppointments, getAppointmentHistory } from './appointments'

describe('Appointments Service', () => {
  it('creates appointment', async () => {
    await expect(createAppointment({ customer_id: 'c1', staff_id: 's1', date: '2024-12-20', services: ['svc1'] })).resolves.toBeUndefined()
  })

  it('creates appointment without services', async () => {
    await expect(createAppointment({ customer_id: 'c1', staff_id: 's1', date: '2024-12-20' })).resolves.toBeUndefined()
  })

  it('updates appointment', async () => {
    await expect(updateAppointment('apt1', { customer_id: 'c1', services: ['svc1'] })).resolves.toBeUndefined()
  })

  it('updates appointment without services', async () => {
    await expect(updateAppointment('apt1', { customer_id: 'c1' })).resolves.toBeUndefined()
  })

  it('updates appointment status', async () => {
    await expect(updateAppointmentStatus('apt1', 'confirmed')).resolves.toBeUndefined()
  })

  it('updates appointment status with note', async () => {
    await expect(updateAppointmentStatus('apt1', 'cancelled', 'Customer no-show')).resolves.toBeUndefined()
  })

  it('deletes appointment', async () => {
    await expect(deleteAppointment('apt1')).resolves.toBeUndefined()
  })

  it('gets appointment', async () => {
    const apt = await getAppointment('apt1')
    expect(apt.id).toBe('apt1')
    expect(apt.status).toBe('booked')
  })

  it('lists appointments', async () => {
    const result = await listAppointments()
    expect(result).toHaveLength(1)
  })

  it('lists appointments with filter', async () => {
    const result = await listAppointments({ staff_id: 's1', customer_id: 'c1', status: 'booked', date_from: '2024-12-01', date_to: '2024-12-31' })
    expect(result).toHaveLength(1)
  })

  it('gets appointment history', async () => {
    const history = await getAppointmentHistory('apt1')
    expect(history).toHaveLength(1)
  })

  it('listAppointments returns empty array when API returns undefined', async () => {
    vi.spyOn(window.go.main.AppointmentService, 'ListAppointments').mockResolvedValueOnce([undefined as any, 0])
    const result = await listAppointments()
    expect(result).toEqual([])
  })
})
