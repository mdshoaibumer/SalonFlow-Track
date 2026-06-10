import { describe, it, expect } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useLocalStorage } from './useLocalStorage'

describe('useLocalStorage', () => {
  beforeEach(() => {
    window.localStorage.clear()
  })

  it('returns initial value when key not in storage', () => {
    const { result } = renderHook(() => useLocalStorage('test-key', 'default'))
    expect(result.current[0]).toBe('default')
  })

  it('reads existing value from localStorage', () => {
    window.localStorage.setItem('existing-key', JSON.stringify('stored-value'))
    const { result } = renderHook(() => useLocalStorage('existing-key', 'default'))
    expect(result.current[0]).toBe('stored-value')
  })

  it('updates value and stores in localStorage', () => {
    const { result } = renderHook(() => useLocalStorage('test-key', 'initial'))
    act(() => { result.current[1]('updated') })
    expect(result.current[0]).toBe('updated')
    expect(JSON.parse(window.localStorage.getItem('test-key')!)).toBe('updated')
  })

  it('handles invalid JSON in localStorage gracefully', () => {
    window.localStorage.setItem('bad-key', 'not-json{{{')
    const { result } = renderHook(() => useLocalStorage('bad-key', 'fallback'))
    expect(result.current[0]).toBe('fallback')
  })
})
