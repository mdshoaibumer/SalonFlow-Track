import { describe, it, expect, vi } from 'vitest'
import { getCloudConfig, saveCloudConfig, testCloudConnection, backupNow, restoreBackup, listCloudHistory, getCloudStats } from './cloud-backup'

describe('Cloud Backup Service', () => {
  it('gets cloud config', async () => {
    const config = await getCloudConfig()
    expect(config.provider).toBe('google_drive')
  })

  it('saves cloud config', async () => {
    await expect(saveCloudConfig({ provider: 'google_drive', auto_backup: true })).resolves.toBeUndefined()
  })

  it('tests cloud connection', async () => {
    await expect(testCloudConnection()).resolves.toBeUndefined()
  })

  it('backs up now', async () => {
    const history = await backupNow()
    expect(history.id).toBe('cbk1')
  })

  it('restores backup', async () => {
    await expect(restoreBackup('cbk1')).resolves.toBeUndefined()
  })

  it('lists cloud history', async () => {
    const result = await listCloudHistory()
    expect(result.data).toHaveLength(1)
    expect(result.total).toBe(1)
  })

  it('lists cloud history with params', async () => {
    const result = await listCloudHistory(10, 5)
    expect(result.data).toHaveLength(1)
  })

  it('gets cloud stats', async () => {
    const stats = await getCloudStats()
    expect(stats.total_backups).toBe(10)
  })

  it('listCloudHistory returns empty when API returns undefined', async () => {
    vi.spyOn(window.go.main.CloudBackupService, 'ListHistory').mockResolvedValueOnce([undefined as any, 0])
    const r = await listCloudHistory()
    expect(r.data).toEqual([])
  })
})
