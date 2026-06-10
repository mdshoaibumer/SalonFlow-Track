import { describe, it, expect } from 'vitest'
import { checkForUpdate, downloadUpdate, installUpdate, getUpdateStatus, listUpdateHistory } from './update'

describe('Update Service', () => {
  it('checks for update', async () => {
    const status = await checkForUpdate()
    expect(status.available).toBe(true)
    expect(status.version).toBe('0.3.0')
  })

  it('downloads update', async () => {
    const record = await downloadUpdate()
    expect(record.status).toBe('downloaded')
  })

  it('installs update', async () => {
    const record = await installUpdate()
    expect(record.status).toBe('installed')
  })

  it('gets update status', async () => {
    const status = await getUpdateStatus()
    expect(status.current_version).toBe('0.2.0')
  })

  it('lists update history', async () => {
    const result = await listUpdateHistory()
    expect(result.records).toHaveLength(1)
    expect(result.total).toBe(1)
  })

  it('lists update history with pagination', async () => {
    const result = await listUpdateHistory(2, 10)
    expect(result.records).toHaveLength(1)
  })
})
