import { apiClient } from './api-client'
import type { ImportJob, ImportPreview, ImportLog, ImportUploadResult, ColumnMapping } from '@/types'

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

export async function uploadFile(file: File, targetEntity?: string): Promise<ImportUploadResult> {
  const formData = new FormData()
  formData.append('file', file)
  if (targetEntity) formData.append('target_entity', targetEntity)

  const res = await fetch(`${BASE_URL}/import/upload`, {
    method: 'POST',
    body: formData,
  })
  const json = await res.json()
  if (!json.success) throw new Error(json.error?.message || 'Upload failed')
  return json.data as ImportUploadResult
}

export async function validateImport(jobId: string, mappings: ColumnMapping[]): Promise<ImportPreview> {
  const response = await apiClient.post<ImportPreview>('/import/validate', { job_id: jobId, mappings })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Validation failed')
  return response.data
}

export async function processImport(jobId: string): Promise<ImportJob> {
  const response = await apiClient.post<ImportJob>('/import/process', { job_id: jobId })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Import failed')
  return response.data
}

export async function listImportJobs(page = 1, perPage = 20): Promise<{ jobs: ImportJob[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  const response = await apiClient.get<ImportJob[]>(`/import/history?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list jobs')
  return {
    jobs: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}

export async function getImportJob(id: string): Promise<ImportJob> {
  const response = await apiClient.get<ImportJob>(`/import/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get job')
  return response.data
}

export async function listImportLogs(jobId: string, status?: string, page = 1, perPage = 50): Promise<{ logs: ImportLog[]; meta: { page: number; per_page: number; total: number; total_pages: number } }> {
  let url = `/import/${jobId}/logs?page=${page}&per_page=${perPage}`
  if (status) url += `&status=${status}`
  const response = await apiClient.get<ImportLog[]>(url)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list logs')
  return {
    logs: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 50,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}
