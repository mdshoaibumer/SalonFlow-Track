import { describe, it, expect, vi } from 'vitest'
import { createBackup, verifyBackup, restoreBackup, listBackups, listRestores, getBackupStats, deleteBackup } from './backup'

describe('Backup Service', () => {
  it('creates backup', async () => {
    const backup = await createBackup()
    expect(backup.id).toBe('bk1')
    expect(backup.backup_type).toBe('manual')
  })

  it('creates backup with type', async () => {
    const backup = await createBackup('scheduled')
    expect(backup.id).toBe('bk1')
  })

  it('verifies backup', async () => {
    const result = await verifyBackup('bk1')
    expect(result.valid).toBe(true)
  })

  it('restores backup', async () => {
    const result = await restoreBackup('bk1', 'test restore')
    expect(result.id).toBe('rst1')
  })

  it('restores backup without notes', async () => {
    const result = await restoreBackup('bk1')
    expect(result.id).toBe('rst1')
  })

  it('lists backups', async () => {
    const result = await listBackups()
    expect(result.backups).toHaveLength(1)
    expect(result.meta.total).toBe(1)
  })

  it('lists backups with pagination', async () => {
    const result = await listBackups(2, 10)
    expect(result.backups).toHaveLength(1)
  })

  it('lists restores', async () => {
    const result = await listRestores()
    expect(result.restores).toHaveLength(1)
  })

  it('lists restores with pagination', async () => {
    const result = await listRestores(2, 10)
    expect(result.restores).toHaveLength(1)
  })

  it('gets backup stats', async () => {
    const stats = await getBackupStats()
    expect(stats.total_backups).toBe(5)
  })

  it('deletes backup', async () => {
    await expect(deleteBackup('bk1')).resolves.toBeUndefined()
  })

  it('listBackups returns empty array when API returns undefined', async () => {
    vi.spyOn(window.go.main.BackupService, 'ListBackups').mockResolvedValueOnce([undefined as any, 0])
    const result = await listBackups()
    expect(result.backups).toEqual([])
  })

  it('listRestores returns empty array when API returns undefined', async () => {
    vi.spyOn(window.go.main.BackupService, 'ListRestores').mockResolvedValueOnce([undefined as any, 0])
    const result = await listRestores()
    expect(result.restores).toEqual([])
  })
})
