import { apiClient } from './api-client'
import type { Service, ServiceStats, CreateServiceInput, UpdateServiceInput } from '@/types'

export interface ListServiceParams {
  page?: number
  per_page?: number
  search?: string
  status?: string
  category?: string
}

export interface ListServiceResponse {
  services: Service[]
  meta: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
}

export async function listServices(params: ListServiceParams = {}): Promise<ListServiceResponse> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.search) query.set('search', params.search)
  if (params.status) query.set('status', params.status)
  if (params.category) query.set('category', params.category)

  const qs = query.toString()
  const path = qs ? `/services?${qs}` : '/services'
  const response = await apiClient.get<Service[]>(path)

  if (!response.success) {
    throw new Error(response.error?.message || 'Failed to fetch services')
  }

  return {
    services: response.data || [],
    meta: response.meta as { page: number; per_page: number; total: number; total_pages: number } || { page: 1, per_page: 20, total: 0, total_pages: 0 },
  }
}

export async function getServiceById(id: string): Promise<Service> {
  const response = await apiClient.get<Service>(`/services/${id}`)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch service')
  }
  return response.data
}

export async function createService(input: CreateServiceInput): Promise<Service> {
  const response = await apiClient.post<Service>('/services', input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to create service')
  }
  return response.data
}

export async function updateService(id: string, input: UpdateServiceInput): Promise<Service> {
  const response = await apiClient.put<Service>(`/services/${id}`, input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to update service')
  }
  return response.data
}

export async function deleteService(id: string): Promise<void> {
  const response = await apiClient.delete<{ message: string }>(`/services/${id}`)
  if (!response.success) {
    throw new Error(response.error?.message || 'Failed to delete service')
  }
}

export async function getServiceStats(): Promise<ServiceStats> {
  const response = await apiClient.get<ServiceStats>('/services/stats')
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch service stats')
  }
  return response.data
}
