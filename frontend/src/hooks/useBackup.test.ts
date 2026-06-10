import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useBackups, useBackupStats, useRestores, useCreateBackup, useVerifyBackup, useRestoreBackup, useDeleteBackup } from './useBackup'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useBackup hooks', () => {
  it('useBackups fetches list', async () => {
    const { result } = renderHook(() => useBackups(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.backups).toHaveLength(1)
  })

  it('useBackupStats fetches stats', async () => {
    const { result } = renderHook(() => useBackupStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_backups).toBe(5)
  })

  it('useRestores fetches list', async () => {
    const { result } = renderHook(() => useRestores(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.restores).toHaveLength(1)
  })

  it('useCreateBackup executes mutation', async () => {
    const { result } = renderHook(() => useCreateBackup(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('full') })
  })

  it('useVerifyBackup executes mutation', async () => {
    const { result } = renderHook(() => useVerifyBackup(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('backup1') })
  })

  it('useRestoreBackup executes mutation', async () => {
    const { result } = renderHook(() => useRestoreBackup(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'backup1', notes: 'restoring' }) })
  })

  it('useDeleteBackup executes mutation', async () => {
    const { result } = renderHook(() => useDeleteBackup(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('backup1') })
  })
})
