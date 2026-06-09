import { apiClient } from './api-client'

export interface HealthStatus {
  status: string
  version: string
  environment: string
  database: string
  uptime: string
  go_version: string
}

export async function getHealthStatus(): Promise<HealthStatus> {
  const response = await apiClient.get<HealthStatus>('/health')
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Failed to fetch health status')
  }
  return response.data
}
