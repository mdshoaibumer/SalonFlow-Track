import { describe, it, expect, vi } from 'vitest'
import { uploadFile, validateImport, processImport, listImportJobs, getImportJob, listImportLogs } from './import'

describe('Import Service', () => {
  it('uploads file', async () => {
    const file = new File(['data'], 'test.csv', { type: 'text/csv' })
    const result = await uploadFile(file)
    expect(result.job.id).toBe('imp1')
    expect(result.columns).toHaveLength(3)
  })

  it('uploads file with target entity', async () => {
    const file = new File(['data'], 'test.csv', { type: 'text/csv' })
    const result = await uploadFile(file, 'customers')
    expect(result.job.id).toBe('imp1')
  })

  it('validates import', async () => {
    const result = await validateImport('imp1', [{ source: 'Name', target: 'full_name' }] as any)
    expect(result.valid_rows).toBe(10)
  })

  it('processes import', async () => {
    const job = await processImport('imp1')
    expect(job.status).toBe('completed')
  })

  it('lists import jobs', async () => {
    const result = await listImportJobs()
    expect(result.jobs).toHaveLength(1)
    expect(result.meta.total).toBe(1)
  })

  it('lists import jobs with pagination', async () => {
    const result = await listImportJobs(2, 10)
    expect(result.jobs).toHaveLength(1)
  })

  it('gets import job', async () => {
    const job = await getImportJob('imp1')
    expect(job.status).toBe('completed')
  })

  it('lists import logs', async () => {
    const result = await listImportLogs('imp1')
    expect(result.logs).toHaveLength(1)
  })

  it('lists import logs with status and pagination', async () => {
    const result = await listImportLogs('imp1', 'success', 2, 10)
    expect(result.logs).toHaveLength(1)
  })

  it('listImportJobs returns empty when API returns undefined', async () => {
    vi.spyOn(window.go.main.ImportService, 'ListJobs').mockResolvedValueOnce([undefined as any, 0])
    const r = await listImportJobs()
    expect(r.jobs).toEqual([])
  })

  it('listImportLogs returns empty when API returns undefined', async () => {
    vi.spyOn(window.go.main.ImportService, 'ListLogs').mockResolvedValueOnce([undefined as any, 0])
    const r = await listImportLogs('imp1')
    expect(r.logs).toEqual([])
  })
})
