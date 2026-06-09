import { apiClient } from './api-client'
import type { CloudBackupConfig, CloudBackupHistory, CloudBackupStats } from '@/types'

export async function getCloudConfig(): Promise<CloudBackupConfig> {
  const response = await apiClient.get<CloudBackupConfig>('/cloud-backup/config')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get config')
  return response.data
}

export async function saveCloudConfig(data: Partial<CloudBackupConfig>): Promise<CloudBackupConfig> {
  const response = await apiClient.post<CloudBackupConfig>('/cloud-backup/config', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to save config')
  return response.data
}

export async function testCloudConnection(): Promise<void> {
  const response = await apiClient.post<{ message: string }>('/cloud-backup/test', {})
  if (!response.success) throw new Error(response.error?.message || 'Connection test failed')
}

export async function backupNow(): Promise<CloudBackupHistory> {
  const response = await apiClient.post<CloudBackupHistory>('/cloud-backup/backup', {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to backup')
  return response.data
}

export async function restoreBackup(id: string): Promise<void> {
  const response = await apiClient.post<{ message: string }>(`/cloud-backup/restore/${id}`, {})
  if (!response.success) throw new Error(response.error?.message || 'Failed to restore')
}

export async function listCloudHistory(page = 1, perPage = 20): Promise<{ data: CloudBackupHistory[]; total: number }> {
  const response = await apiClient.get<CloudBackupHistory[]>(`/cloud-backup/history?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list history')
  return { data: response.data || [], total: response.meta?.total || 0 }
}

export async function getCloudStats(): Promise<CloudBackupStats> {
  const response = await apiClient.get<CloudBackupStats>('/cloud-backup/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get stats')
  return response.data
}
