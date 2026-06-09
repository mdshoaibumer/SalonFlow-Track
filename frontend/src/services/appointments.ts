import { apiClient } from './api-client'
import type { Appointment, AppointmentService, AppointmentHistory, AppointmentFilter } from '@/types'

export async function listAppointments(filter?: AppointmentFilter): Promise<Appointment[]> {
  const params = new URLSearchParams()
  if (filter?.start_date) params.set('start_date', filter.start_date)
  if (filter?.end_date) params.set('end_date', filter.end_date)
  if (filter?.staff_id) params.set('staff_id', filter.staff_id)
  if (filter?.customer_id) params.set('customer_id', filter.customer_id)
  if (filter?.status) params.set('status', filter.status)
  const q = params.toString()
  const response = await apiClient.get<Appointment[]>(`/appointments${q ? '?' + q : ''}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list appointments')
  return response.data || []
}

export async function getAppointment(id: string): Promise<Appointment> {
  const response = await apiClient.get<Appointment>(`/appointments/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get appointment')
  return response.data
}

export async function createAppointment(data: Partial<Appointment> & { services?: Partial<AppointmentService>[] }): Promise<Appointment> {
  const response = await apiClient.post<Appointment>('/appointments', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create appointment')
  return response.data
}

export async function updateAppointment(id: string, data: Partial<Appointment>): Promise<Appointment> {
  const response = await apiClient.put<Appointment>(`/appointments/${id}`, data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update appointment')
  return response.data
}

export async function updateAppointmentStatus(id: string, status: string): Promise<Appointment> {
  const response = await apiClient.put<Appointment>(`/appointments/${id}/status`, { status })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update status')
  return response.data
}

export async function deleteAppointment(id: string): Promise<void> {
  const response = await apiClient.delete(`/appointments/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete appointment')
}

export async function getAppointmentHistory(id: string): Promise<AppointmentHistory[]> {
  const response = await apiClient.get<AppointmentHistory[]>(`/appointments/${id}/history`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to get history')
  return response.data || []
}
