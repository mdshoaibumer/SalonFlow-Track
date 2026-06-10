import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useImportJobs, useImportJob, useImportLogs, useUploadFile, useValidateImport, useProcessImport } from './useImport'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useImport hooks', () => {
  it('useImportJobs fetches jobs', async () => {
    const { result } = renderHook(() => useImportJobs(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.jobs).toHaveLength(1)
  })

  it('useImportJob fetches job', async () => {
    const { result } = renderHook(() => useImportJob('imp1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.status).toBe('completed')
  })

  it('useImportLogs fetches logs', async () => {
    const { result } = renderHook(() => useImportLogs('imp1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.logs).toHaveLength(1)
  })

  it('useUploadFile executes mutation', async () => {
    const { result } = renderHook(() => useUploadFile(), { wrapper: createWrapper() })
    const file = new File(['test'], 'test.csv', { type: 'text/csv' })
    await act(async () => { await result.current.mutateAsync({ file, targetEntity: 'customers' }) })
  })

  it('useValidateImport executes mutation', async () => {
    const { result } = renderHook(() => useValidateImport(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ jobId: 'imp1', mappings: [] }) })
  })

  it('useProcessImport executes mutation', async () => {
    const { result } = renderHook(() => useProcessImport(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('imp1') })
  })
})
