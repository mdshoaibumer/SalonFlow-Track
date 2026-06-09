import { apiClient } from './api-client'
import type { UpdateStatus, UpdateRecord } from '@/types'

export async function checkForUpdate(): Promise<UpdateStatus> {
  const response = await apiClient.get<UpdateStatus>('/update/check')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to check for update')
  return response.data
}

export async function downloadUpdate(): Promise<UpdateRecord> {
  const response = await apiClient.post<UpdateRecord>('/update/download', {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to download update')
  return response.data
}

export async function installUpdate(): Promise<UpdateRecord> {
  const response = await apiClient.post<UpdateRecord>('/update/install', {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to install update')
  return response.data
}

export async function getUpdateStatus(): Promise<UpdateStatus> {
  const response = await apiClient.get<UpdateStatus>('/update/status')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get update status')
  return response.data
}

export async function listUpdateHistory(page = 1, perPage = 20): Promise<{ records: UpdateRecord[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const response = await apiClient.get<UpdateRecord[]>(`/update/history?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list update history')
  return {
    records: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}
