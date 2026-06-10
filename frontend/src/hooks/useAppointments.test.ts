import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useAppointments, useAppointment, useCreateAppointment, useUpdateAppointment, useUpdateAppointmentStatus, useDeleteAppointment, useAppointmentHistory } from './useAppointments'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useAppointments hooks', () => {
  it('useAppointments fetches list', async () => {
    const { result } = renderHook(() => useAppointments(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useAppointment fetches single', async () => {
    const { result } = renderHook(() => useAppointment('apt1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.status).toBe('booked')
  })

  it('useAppointmentHistory fetches history', async () => {
    const { result } = renderHook(() => useAppointmentHistory('apt1'), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useCreateAppointment executes mutation', async () => {
    const { result } = renderHook(() => useCreateAppointment(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ customer_id: 'c1', date: '2024-01-01' } as any) })
  })

  it('useUpdateAppointment executes mutation', async () => {
    const { result } = renderHook(() => useUpdateAppointment(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'apt1', notes: 'updated' } as any) })
  })

  it('useUpdateAppointmentStatus executes mutation', async () => {
    const { result } = renderHook(() => useUpdateAppointmentStatus(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'apt1', status: 'completed' }) })
  })

  it('useDeleteAppointment executes mutation', async () => {
    const { result } = renderHook(() => useDeleteAppointment(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('apt1') })
  })
})
