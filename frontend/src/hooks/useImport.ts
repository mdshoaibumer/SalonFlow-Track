import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as importService from '@/services/import'
import type { ColumnMapping } from '@/types'

const IMPORT_KEYS = {
  all: ['import'] as const,
  jobs: (page: number) => ['import', 'jobs', page] as const,
  job: (id: string) => ['import', 'job', id] as const,
  logs: (jobId: string, status: string, page: number) => ['import', 'logs', jobId, status, page] as const,
}

export function useImportJobs(page = 1) {
  return useQuery({
    queryKey: IMPORT_KEYS.jobs(page),
    queryFn: () => importService.listImportJobs(page),
  })
}

export function useImportJob(id: string) {
  return useQuery({
    queryKey: IMPORT_KEYS.job(id),
    queryFn: () => importService.getImportJob(id),
    enabled: !!id,
  })
}

export function useImportLogs(jobId: string, status = '', page = 1) {
  return useQuery({
    queryKey: IMPORT_KEYS.logs(jobId, status, page),
    queryFn: () => importService.listImportLogs(jobId, status || undefined, page),
    enabled: !!jobId,
  })
}

export function useUploadFile() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ file, targetEntity }: { file: File; targetEntity?: string }) =>
      importService.uploadFile(file, targetEntity),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: IMPORT_KEYS.all })
    },
  })
}

export function useValidateImport() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ jobId, mappings }: { jobId: string; mappings: ColumnMapping[] }) =>
      importService.validateImport(jobId, mappings),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: IMPORT_KEYS.all })
    },
  })
}

export function useProcessImport() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (jobId: string) => importService.processImport(jobId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: IMPORT_KEYS.all })
    },
  })
}
