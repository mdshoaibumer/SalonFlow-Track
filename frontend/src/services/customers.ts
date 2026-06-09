import { apiClient } from './api-client'
import type { Customer, CustomerStats, CreateCustomerInput, UpdateCustomerInput } from '@/types'

export interface ListCustomerParams {
  page?: number
  per_page?: number
  search?: string
  status?: string
}

export interface ListCustomerResponse {
  customers: Customer[]
  meta: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
}

export async function listCustomers(params: ListCustomerParams = {}): Promise<ListCustomerResponse> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.search) query.set('search', params.search)
  if (params.status) query.set('status', params.status)

  const qs = query.toString()
  const path = qs ? `/customers?${qs}` : '/customers'
  const response = await apiClient.get<Customer[]>(path)

  if (!response.success) {
    throw new Error(response.error?.message || 'Failed to fetch customers')
  }

  return {
    customers: response.data || [],
    meta: response.meta as { page: number; per_page: number; total: number; total_pages: number } || { page: 1, per_page: 20, total: 0, total_pages: 0 },
  }
}

export async function getCustomerById(id: string): Promise<Customer> {
  const response = await apiClient.get<Customer>(`/customers/${id}`)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch customer')
  }
  return response.data
}

export async function createCustomer(input: CreateCustomerInput): Promise<Customer> {
  const response = await apiClient.post<Customer>('/customers', input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to create customer')
  }
  return response.data
}

export async function updateCustomer(id: string, input: UpdateCustomerInput): Promise<Customer> {
  const response = await apiClient.put<Customer>(`/customers/${id}`, input)
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to update customer')
  }
  return response.data
}

export async function deleteCustomer(id: string): Promise<void> {
  const response = await apiClient.delete<{ message: string }>(`/customers/${id}`)
  if (!response.success) {
    throw new Error(response.error?.message || 'Failed to delete customer')
  }
}

export async function getCustomerStats(): Promise<CustomerStats> {
  const response = await apiClient.get<CustomerStats>('/customers/stats')
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch customer stats')
  }
  return response.data
}
