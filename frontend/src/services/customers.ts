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
  const result = await window.go.main.CustomerService.ListCustomers({
    search: params.search || '',
    status: params.status || '',
    page: params.page || 1,
    per_page: params.per_page || 20,
  })
  return {
    customers: result.customers || [],
    meta: {
      page: result.page,
      per_page: result.per_page,
      total: result.total,
      total_pages: result.total_pages,
    },
  }
}

export async function getCustomerById(id: string): Promise<Customer> {
  return window.go.main.CustomerService.GetCustomer(id)
}

export async function createCustomer(input: CreateCustomerInput): Promise<Customer> {
  return window.go.main.CustomerService.CreateCustomer(input)
}

export async function updateCustomer(id: string, input: UpdateCustomerInput): Promise<Customer> {
  return window.go.main.CustomerService.UpdateCustomer(id, input)
}

export async function deleteCustomer(id: string): Promise<void> {
  await window.go.main.CustomerService.DeleteCustomer(id)
}

export async function getCustomerStats(): Promise<CustomerStats> {
  return window.go.main.CustomerService.GetCustomerStats()
}
