import { apiClient } from './api-client'
import type { Staff, StaffStats, CreateStaffInput, UpdateStaffInput } from '@/types'

export interface ListStaffParams {
  page?: number
  per_page?: number
  search?: string
  status?: string
  designation?: string
}

export interface ListStaffResponse {
  staff: Staff[]
  meta: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
}

export async function listStaff(params: ListStaffParams = {}): Promise<ListStaffResponse> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.search) query.set('search', params.search)
  if (params.status) query.set('status', params.status)
  if (params.designation) query.set('designation', params.designation)

  const qs = query.toString()
  const path = qs ? `/staff?${qs}` : '/staff'
  const response = await apiClient.get<Staff[]>(path)

  if (!response.success) {
    throw new Error(response.error?.message || 'Failed to fetch staff')
  }

  return {
    staff: response.data || [],
    meta: response.meta as { page: number; per_page: number; total: number; total_pages: number } || { page: 1, per_page: 20, total: 0, total_pages: 0 },
  }
}

export async function getStaffById(id: string): Promise<Staff> {
  const response = await apiClient.get<Staff>(`/staff/${id}`)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch staff')
  }
  return response.data
}

export async function createStaff(input: CreateStaffInput): Promise<Staff> {
  const response = await apiClient.post<Staff>('/staff', input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to create staff')
  }
  return response.data
}

export async function updateStaff(id: string, input: UpdateStaffInput): Promise<Staff> {
  const response = await apiClient.put<Staff>(`/staff/${id}`, input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to update staff')
  }
  return response.data
}

export async function deleteStaff(id: string): Promise<void> {
  const response = await apiClient.delete<{ message: string }>(`/staff/${id}`)
  if (!response.success) {
    throw new Error(response.error?.message || 'Failed to delete staff')
  }
}

export async function getStaffStats(): Promise<StaffStats> {
  const response = await apiClient.get<StaffStats>('/staff/stats')
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch staff stats')
  }
  return response.data
}
