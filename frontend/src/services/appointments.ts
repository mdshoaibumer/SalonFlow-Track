import type { Appointment, AppointmentHistory } from '@/types'

export interface AppointmentFilter {
  staff_id?: string
  customer_id?: string
  status?: string
  date_from?: string
  date_to?: string
  page?: number
  per_page?: number
}

export async function createAppointment(data: any): Promise<void> {
  const { services, ...appt } = data
  await window.go.main.AppointmentService.CreateAppointment(appt, services || [])
}

export async function updateAppointment(id: string, data: any): Promise<void> {
  const { services, ...appt } = data
  await window.go.main.AppointmentService.UpdateAppointment({ ...appt, id }, services || [])
}

export async function updateAppointmentStatus(id: string, status: string, note = ''): Promise<void> {
  await window.go.main.AppointmentService.UpdateAppointmentStatus(id, status, note)
}

export async function deleteAppointment(id: string): Promise<void> {
  await window.go.main.AppointmentService.DeleteAppointment(id)
}

export async function getAppointment(id: string): Promise<Appointment> {
  return window.go.main.AppointmentService.GetAppointment(id)
}

export async function listAppointments(filter: AppointmentFilter = {}): Promise<Appointment[]> {
  const [appointments] = await window.go.main.AppointmentService.ListAppointments(filter)
  return appointments || []
}

export async function getAppointmentHistory(id: string): Promise<AppointmentHistory[]> {
  return window.go.main.AppointmentService.GetAppointmentHistory(id)
}
