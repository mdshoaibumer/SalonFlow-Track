import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createElement } from 'react'
import { useWhatsAppTemplates, useCreateTemplate, useUpdateTemplate, useDeleteTemplate, useSendMessage, useWhatsAppMessages, useWAMessageStats, useAutomationRules, useCreateRule, useUpdateRule, useDeleteRule } from './useWhatsApp'

function createWrapper() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('useWhatsApp hooks', () => {
  it('useWhatsAppTemplates fetches templates', async () => {
    const { result } = renderHook(() => useWhatsAppTemplates(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useWhatsAppMessages fetches messages', async () => {
    const { result } = renderHook(() => useWhatsAppMessages(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.data).toHaveLength(1)
  })

  it('useWAMessageStats fetches stats', async () => {
    const { result } = renderHook(() => useWAMessageStats(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data?.total_sent).toBe(100)
  })

  it('useAutomationRules fetches rules', async () => {
    const { result } = renderHook(() => useAutomationRules(), { wrapper: createWrapper() })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toHaveLength(1)
  })

  it('useCreateTemplate executes mutation', async () => {
    const { result } = renderHook(() => useCreateTemplate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ name: 'Welcome', body: 'Hi {{name}}' } as any) })
  })

  it('useUpdateTemplate executes mutation', async () => {
    const { result } = renderHook(() => useUpdateTemplate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'tmpl1', name: 'Updated' } as any) })
  })

  it('useDeleteTemplate executes mutation', async () => {
    const { result } = renderHook(() => useDeleteTemplate(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('tmpl1') })
  })

  it('useSendMessage executes mutation', async () => {
    const { result } = renderHook(() => useSendMessage(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ template_id: 'tmpl1', customer_id: 'c1', phone: '9876543210' }) })
  })

  it('useCreateRule executes mutation', async () => {
    const { result } = renderHook(() => useCreateRule(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ name: 'Birthday', trigger: 'birthday' } as any) })
  })

  it('useUpdateRule executes mutation', async () => {
    const { result } = renderHook(() => useUpdateRule(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync({ id: 'rule1', name: 'Updated Rule' } as any) })
  })

  it('useDeleteRule executes mutation', async () => {
    const { result } = renderHook(() => useDeleteRule(), { wrapper: createWrapper() })
    await act(async () => { await result.current.mutateAsync('rule1') })
  })
})
