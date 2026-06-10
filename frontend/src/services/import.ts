import type { ImportJob, ImportPreview, ColumnMapping } from '@/types'

export interface ImportUploadResult {
  job: ImportJob
  columns: string[]
  headers: string[]
  mappings: ColumnMapping[]
}

export async function uploadFile(file: File, targetEntity?: string): Promise<ImportUploadResult> {
  const fileName = file.name
  const filePath = (file as any).path || file.name
  const [job, columns, mappings] = await window.go.main.ImportService.Upload(fileName, filePath, targetEntity || '')
  return { job, columns, headers: columns, mappings }
}

export async function validateImport(jobId: string, mappings: ColumnMapping[]): Promise<ImportPreview> {
  return window.go.main.ImportService.Validate(jobId, mappings)
}

export async function processImport(jobId: string): Promise<ImportJob> {
  return window.go.main.ImportService.Process(jobId)
}

export async function listImportJobs(page = 1, perPage = 20) {
  const [jobs, total] = await window.go.main.ImportService.ListJobs(page, perPage)
  const total_pages = Math.ceil(total / perPage)
  return { jobs: jobs || [], meta: { page, per_page: perPage, total, total_pages } }
}

export async function getImportJob(id: string): Promise<ImportJob> {
  return window.go.main.ImportService.GetJob(id)
}

export async function listImportLogs(jobId: string, status?: string, page = 1, perPage = 20) {
  const [logs, total] = await window.go.main.ImportService.ListLogs(jobId, status || '', page, perPage)
  const total_pages = Math.ceil(total / perPage)
  return { logs: logs || [], meta: { page, per_page: perPage, total, total_pages } }
}
