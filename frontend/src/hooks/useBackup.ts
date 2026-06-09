import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as backupService from '@/services/backup'

const BACKUP_KEYS = {
  all: ['backups'] as const,
  list: (page: number) => ['backups', 'list', page] as const,
  stats: ['backups', 'stats'] as const,
  restores: (page: number) => ['backups', 'restores', page] as const,
}

export function useBackups(page = 1) {
  return useQuery({
    queryKey: BACKUP_KEYS.list(page),
    queryFn: () => backupService.listBackups(page),
  })
}

export function useBackupStats() {
  return useQuery({
    queryKey: BACKUP_KEYS.stats,
    queryFn: () => backupService.getBackupStats(),
  })
}

export function useRestores(page = 1) {
  return useQuery({
    queryKey: BACKUP_KEYS.restores(page),
    queryFn: () => backupService.listRestores(page),
  })
}

export function useCreateBackup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (backupType?: string) => backupService.createBackup(backupType),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: BACKUP_KEYS.all })
    },
  })
}

export function useVerifyBackup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => backupService.verifyBackup(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: BACKUP_KEYS.all })
    },
  })
}

export function useRestoreBackup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, notes }: { id: string; notes?: string }) => backupService.restoreBackup(id, notes),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: BACKUP_KEYS.all })
    },
  })
}

export function useDeleteBackup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => backupService.deleteBackup(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: BACKUP_KEYS.all })
    },
  })
}
