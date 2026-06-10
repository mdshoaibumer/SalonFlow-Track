import type { UpdateStatus, UpdateRecord } from '@/types'

export async function checkForUpdate(): Promise<UpdateStatus> {
  return window.go.main.UpdateService.CheckForUpdate()
}

export async function downloadUpdate(): Promise<UpdateRecord> {
  return window.go.main.UpdateService.DownloadUpdate()
}

export async function installUpdate(): Promise<UpdateRecord> {
  return window.go.main.UpdateService.InstallUpdate()
}

export async function getUpdateStatus(): Promise<UpdateStatus> {
  return window.go.main.UpdateService.GetUpdateStatus()
}

export async function listUpdateHistory(page = 1, perPage = 20): Promise<{ records: UpdateRecord[]; total: number }> {
  const [records, total] = await window.go.main.UpdateService.ListUpdateHistory(page, perPage)
  return { records, total }
}
