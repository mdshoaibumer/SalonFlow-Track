import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as api from '@/services/cloud-backup'
import type { CloudBackupConfig } from '@/types'

export function useCloudConfig() {
  return useQuery({ queryKey: ['cloud-config'], queryFn: api.getCloudConfig })
}

export function useSaveCloudConfig() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<CloudBackupConfig>) => api.saveCloudConfig(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['cloud-config'] }),
  })
}

export function useTestCloudConnection() {
  return useMutation({ mutationFn: () => api.testCloudConnection() })
}

export function useBackupNow() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.backupNow(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['cloud-history'] }),
  })
}

export function useRestoreBackup() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.restoreBackup(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['cloud-history'] }),
  })
}

export function useCloudHistory(page = 1) {
  return useQuery({ queryKey: ['cloud-history', page], queryFn: () => api.listCloudHistory(page) })
}

export function useCloudStats() {
  return useQuery({ queryKey: ['cloud-stats'], queryFn: api.getCloudStats })
}
