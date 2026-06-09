import { apiClient } from './api-client'
import type { SalaryCycle, SalaryRecord, Advance, SalaryStats, CreateAdvanceInput, GenerateSalaryInput, GenerateSalaryOutput } from '@/types'

export interface ListAdvancesParams {
  page?: number
  per_page?: number
  staff_id?: string
  status?: string
}

export interface ListAdvancesResponse {
  advances: Advance[]
  meta: { page: number; per_page: number; total: number; total_pages: number }
}

export async function generateSalary(input: GenerateSalaryInput): Promise<GenerateSalaryOutput> {
  const response = await apiClient.post<GenerateSalaryOutput>('/salary/generate', input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to generate salary')
  return response.data
}

export async function listSalaries(month: number, year: number): Promise<SalaryRecord[]> {
  const response = await apiClient.get<SalaryRecord[]>(`/salary?month=${month}&year=${year}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch salaries')
  return response.data || []
}

export async function getSalaryById(id: string): Promise<SalaryRecord> {
  const response = await apiClient.get<SalaryRecord>(`/salary/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch salary')
  return response.data
}

export async function paySalary(id: string): Promise<void> {
  const response = await apiClient.post(`/salary/${id}/pay`, {})
  if (!response.success) throw new Error(response.error?.message || 'Failed to pay salary')
}

export async function listSalaryCycles(year?: number): Promise<SalaryCycle[]> {
  const query = year ? `?year=${year}` : ''
  const response = await apiClient.get<SalaryCycle[]>(`/salary/cycles${query}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch cycles')
  return response.data || []
}

export async function getSalaryStats(): Promise<SalaryStats> {
  const response = await apiClient.get<SalaryStats>('/salary/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch salary stats')
  return response.data
}

export async function createAdvance(input: CreateAdvanceInput): Promise<Advance> {
  const response = await apiClient.post<Advance>('/salary/advances', input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create advance')
  return response.data
}

export async function listAdvances(params: ListAdvancesParams = {}): Promise<ListAdvancesResponse> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.staff_id) query.set('staff_id', params.staff_id)
  if (params.status) query.set('status', params.status)
  const qs = query.toString()
  const path = qs ? `/salary/advances?${qs}` : '/salary/advances'
  const response = await apiClient.get<Advance[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch advances')
  return {
    advances: response.data || [],
    meta: response.meta as { page: number; per_page: number; total: number; total_pages: number } || { page: 1, per_page: 20, total: 0, total_pages: 0 },
  }
}

export async function approveAdvance(id: string): Promise<Advance> {
  const response = await apiClient.put<Advance>(`/salary/advances/${id}/approve`, {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to approve advance')
  return response.data
}

export async function rejectAdvance(id: string): Promise<Advance> {
  const response = await apiClient.put<Advance>(`/salary/advances/${id}/reject`, {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to reject advance')
  return response.data
}
