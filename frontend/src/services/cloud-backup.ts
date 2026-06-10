import type { CloudBackupConfig, CloudBackupHistory, CloudBackupStats } from '@/types'

export async function getCloudConfig(): Promise<CloudBackupConfig> {
  return window.go.main.CloudBackupService.GetConfig()
}

export async function saveCloudConfig(cfg: Partial<CloudBackupConfig>): Promise<void> {
  await window.go.main.CloudBackupService.SaveConfig(cfg as CloudBackupConfig)
}

export async function testCloudConnection(): Promise<void> {
  await window.go.main.CloudBackupService.TestConnection()
}

export async function backupNow(): Promise<CloudBackupHistory> {
  return window.go.main.CloudBackupService.BackupNow()
}

export async function restoreBackup(historyId: string): Promise<void> {
  await window.go.main.CloudBackupService.Restore(historyId)
}

export async function listCloudHistory(limit = 20, offset = 0) {
  const [history, total] = await window.go.main.CloudBackupService.ListHistory(limit, offset)
  return { data: history || [], total }
}

export async function getCloudStats(): Promise<CloudBackupStats> {
  return window.go.main.CloudBackupService.GetCloudBackupStats()
}
