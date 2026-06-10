import { describe, it, expect } from 'vitest'
import { apiClient } from './api-client'

describe('API Client', () => {
  it('makes GET requests', async () => {
    const mockResponse = { success: true, data: { id: 1 } }
    globalThis.fetch = async () => new Response(JSON.stringify(mockResponse))

    const result = await apiClient.get('/test')
    expect(result.success).toBe(true)
  })

  it('makes POST requests', async () => {
    const mockResponse = { success: true, data: { id: 1 } }
    globalThis.fetch = async (url: any, opts: any) => {
      expect(opts.method).toBe('POST')
      expect(opts.headers['Content-Type']).toBe('application/json')
      return new Response(JSON.stringify(mockResponse))
    }

    const result = await apiClient.post('/test', { name: 'test' })
    expect(result.success).toBe(true)
  })

  it('makes POST requests without body', async () => {
    const mockResponse = { success: true }
    globalThis.fetch = async () => new Response(JSON.stringify(mockResponse))

    const result = await apiClient.post('/test')
    expect(result.success).toBe(true)
  })

  it('makes PUT requests', async () => {
    const mockResponse = { success: true, data: { id: 1 } }
    globalThis.fetch = async (url: any, opts: any) => {
      expect(opts.method).toBe('PUT')
      return new Response(JSON.stringify(mockResponse))
    }

    const result = await apiClient.put('/test/1', { name: 'updated' })
    expect(result.success).toBe(true)
  })

  it('makes PUT requests without body', async () => {
    const mockResponse = { success: true }
    globalThis.fetch = async () => new Response(JSON.stringify(mockResponse))

    const result = await apiClient.put('/test/1')
    expect(result.success).toBe(true)
  })

  it('makes DELETE requests', async () => {
    const mockResponse = { success: true }
    globalThis.fetch = async (url: any, opts: any) => {
      expect(opts.method).toBe('DELETE')
      return new Response(JSON.stringify(mockResponse))
    }

    const result = await apiClient.delete('/test/1')
    expect(result.success).toBe(true)
  })
})
