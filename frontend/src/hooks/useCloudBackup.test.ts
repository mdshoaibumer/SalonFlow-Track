import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useCloudConfig, useSaveCloudConfig, useTestCloudConnection, useBackupNow, useRestoreBackup, useCloudHistory, useCloudStats } from './useCloudBackup'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useCloudBackup hooks', () => {
  it('useCloudConfig fetches config', async () => {
    const { result } = renderHook(() => useCloudConfig(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.provider).toBe('google_drive')
  })

  it('useCloudHistory fetches history', async () => {
    const { result } = renderHook(() => useCloudHistory(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.data).toHaveLength(1)
  })

  it('useCloudStats fetches stats', async () => {
    const { result } = renderHook(() => useCloudStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_backups).toBe(10)
  })

  it('useSaveCloudConfig executes mutation', async () => {
    const { result } = renderHook(() => useSaveCloudConfig(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ provider: 'google_drive' } as any) })
  })

  it('useTestCloudConnection executes mutation', async () => {
    const { result } = renderHook(() => useTestCloudConnection(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync() })
  })

  it('useBackupNow executes mutation', async () => {
    const { result } = renderHook(() => useBackupNow(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync() })
  })

  it('useRestoreBackup executes mutation', async () => {
    const { result } = renderHook(() => useRestoreBackup(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('cloud-backup-1') })
  })
})
