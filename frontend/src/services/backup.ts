import { apiClient } from './api-client'
import type { BackupRecord, RestoreRecord, BackupStats, BackupVerification } from '@/types'

export async function createBackup(backupType?: string): Promise<BackupRecord> {
  const response = await apiClient.post<BackupRecord>('/backups', { backup_type: backupType || 'manual' })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create backup')
  return response.data
}

export async function listBackups(page = 1, perPage = 20): Promise<{ backups: BackupRecord[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const response = await apiClient.get<BackupRecord[]>(`/backups?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list backups')
  return {
    backups: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}

export async function getBackupStats(): Promise<BackupStats> {
  const response = await apiClient.get<BackupStats>('/backups/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get backup stats')
  return response.data
}

export async function verifyBackup(id: string): Promise<BackupVerification> {
  const response = await apiClient.post<BackupVerification>(`/backups/${id}/verify`, {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to verify backup')
  return response.data
}

export async function restoreBackup(id: string, notes?: string): Promise<RestoreRecord> {
  const response = await apiClient.post<RestoreRecord>(`/backups/${id}/restore`, { notes: notes || '' })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to restore backup')
  return response.data
}

export async function deleteBackup(id: string): Promise<void> {
  const response = await apiClient.delete(`/backups/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete backup')
}

export async function listRestores(page = 1, perPage = 20): Promise<{ restores: RestoreRecord[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const response = await apiClient.get<RestoreRecord[]>(`/backups/restores?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list restores')
  return {
    restores: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}
