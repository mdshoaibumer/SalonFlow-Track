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
  const result = await window.go.main.StaffService.ListStaff({
    search: params.search || '',
    status: params.status || '',
    designation: params.designation || '',
    page: params.page || 1,
    per_page: params.per_page || 20,
  })
  return {
    staff: result.staff || [],
    meta: {
      page: result.page,
      per_page: result.per_page,
      total: result.total,
      total_pages: result.total_pages,
    },
  }
}

export async function getStaffById(id: string): Promise<Staff> {
  return window.go.main.StaffService.GetStaff(id)
}

export async function createStaff(input: CreateStaffInput): Promise<Staff> {
  return window.go.main.StaffService.CreateStaff(input)
}

export async function updateStaff(id: string, input: UpdateStaffInput): Promise<Staff> {
  return window.go.main.StaffService.UpdateStaff(id, input)
}

export async function deleteStaff(id: string): Promise<void> {
  await window.go.main.StaffService.DeleteStaff(id)
}

export async function getStaffStats(): Promise<StaffStats> {
  return window.go.main.StaffService.GetStaffStats()
}
