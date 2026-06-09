import { describe, it, expect } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useStaffList, useStaffStats } from '@/hooks/useStaff'

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useStaff hooks', () => {
  it('useStaffList fetches staff list', async () => {
    const { result } = renderHook(() => useStaffList(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.staff).toHaveLength(1)
    expect(result.current.data?.staff[0].full_name).toBe('Priya Sharma')
  })

  it('useStaffStats fetches stats', async () => {
    const { result } = renderHook(() => useStaffStats(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.total).toBe(5)
    expect(result.current.data?.active).toBe(4)
  })
})
