import type { BackupRecord, BackupVerification, RestoreRecord, BackupStats } from '@/types'

export async function createBackup(backupType = 'manual'): Promise<BackupRecord> {
  return window.go.main.BackupService.CreateBackup(backupType)
}

export async function verifyBackup(id: string): Promise<BackupVerification> {
  return window.go.main.BackupService.VerifyBackup(id)
}

export async function restoreBackup(id: string, notes = ''): Promise<RestoreRecord> {
  return window.go.main.BackupService.RestoreBackup(id, notes)
}

export async function listBackups(page = 1, perPage = 20) {
  const [backups, total] = await window.go.main.BackupService.ListBackups(page, perPage)
  return { backups: backups || [], meta: { page, per_page: perPage, total, total_pages: Math.ceil(total / perPage) } }
}

export async function listRestores(page = 1, perPage = 20) {
  const [restores, total] = await window.go.main.BackupService.ListRestores(page, perPage)
  return { restores: restores || [], meta: { page, per_page: perPage, total, total_pages: Math.ceil(total / perPage) } }
}

export async function getBackupStats(): Promise<BackupStats> {
  return window.go.main.BackupService.GetBackupStats()
}

export async function deleteBackup(id: string): Promise<void> {
  await window.go.main.BackupService.DeleteBackup(id)
}
