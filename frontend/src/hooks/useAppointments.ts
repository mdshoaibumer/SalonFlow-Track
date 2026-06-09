import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as api from '@/services/appointments'
import type { Appointment, AppointmentFilter, AppointmentService } from '@/types'

export function useAppointments(filter?: AppointmentFilter) {
  return useQuery({ queryKey: ['appointments', filter], queryFn: () => api.listAppointments(filter) })
}

export function useAppointment(id: string) {
  return useQuery({ queryKey: ['appointments', id], queryFn: () => api.getAppointment(id), enabled: !!id })
}

export function useCreateAppointment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<Appointment> & { services?: Partial<AppointmentService>[] }) => api.createAppointment(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['appointments'] }),
  })
}

export function useUpdateAppointment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: Partial<Appointment> & { id: string }) => api.updateAppointment(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['appointments'] }),
  })
}

export function useUpdateAppointmentStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) => api.updateAppointmentStatus(id, status),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['appointments'] }),
  })
}

export function useDeleteAppointment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.deleteAppointment(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['appointments'] }),
  })
}

export function useAppointmentHistory(id: string) {
  return useQuery({ queryKey: ['appointments', id, 'history'], queryFn: () => api.getAppointmentHistory(id), enabled: !!id })
}
