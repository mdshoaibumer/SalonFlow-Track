import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useStaffList, useStaffById, useStaffStats, useCreateStaff, useUpdateStaff, useDeleteStaff } from '@/hooks/useStaff'

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
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

  it('useStaffById fetches staff', async () => {
    const { result } = renderHook(() => useStaffById('staff1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.full_name).toBe('Priya Sharma')
  })

  it('useStaffStats fetches stats', async () => {
    const { result } = renderHook(() => useStaffStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total).toBe(5)
    expect(result.current.data?.active).toBe(4)
  })

  it('useCreateStaff executes mutation', async () => {
    const { result } = renderHook(() => useCreateStaff(), { wrapper: createWrapper() })
    await act(async () => {
      const staff = await result.current.mutateAsync({ full_name: 'New', phone: '9876543210', gender: 'male', designation: 'stylist', joining_date: '2024-01-01', base_salary: 20000, commission_percentage: 10 } as any)
      expect(staff.full_name).toBe('Priya Sharma')
    })
  })

  it('useUpdateStaff executes mutation', async () => {
    const { result } = renderHook(() => useUpdateStaff(), { wrapper: createWrapper() })
    await act(async () => {
      const staff = await result.current.mutateAsync({ id: 'staff1', input: { full_name: 'Updated' } as any })
      expect(staff.full_name).toBe('Priya Sharma')
    })
  })

  it('useDeleteStaff executes mutation', async () => {
    const { result } = renderHook(() => useDeleteStaff(), { wrapper: createWrapper() })
    await act(async () => {
      await result.current.mutateAsync('staff1')
    })
  })
})
