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
  const result = await window.go.main.ServiceService.ListServices({
    search: params.search || '',
    status: params.status || '',
    category: params.category || '',
    page: params.page || 1,
    per_page: params.per_page || 20,
  })
  return {
    services: result.services || [],
    meta: {
      page: result.page,
      per_page: result.per_page,
      total: result.total,
      total_pages: result.total_pages,
    },
  }
}

export async function getServiceById(id: string): Promise<Service> {
  return window.go.main.ServiceService.GetService(id)
}

export async function createService(input: CreateServiceInput): Promise<Service> {
  return window.go.main.ServiceService.CreateService(input)
}

export async function updateService(id: string, input: UpdateServiceInput): Promise<Service> {
  return window.go.main.ServiceService.UpdateService(id, input)
}

export async function deleteService(id: string): Promise<void> {
  await window.go.main.ServiceService.DeleteService(id)
}

export async function getServiceStats(): Promise<ServiceStats> {
  return window.go.main.ServiceService.GetServiceStats()
}
