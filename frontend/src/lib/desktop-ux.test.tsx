import { describe, it, expect, vi, afterEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { type ReactNode } from 'react'
import { MemoryRouter } from 'react-router-dom'
import {
  useKeyboardShortcut,
  useNavigationShortcuts,
  useWindowDrag,
  useEscapeKey,
  useFocusTrap,
  useDebouncedValue,
} from './desktop-ux'

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

const wrapper = ({ children }: { children: ReactNode }) => (
  <MemoryRouter>{children}</MemoryRouter>
)

describe('useKeyboardShortcut', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('calls handler when matching key combo is pressed', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', ctrlKey: true }))
    })

    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('does not call handler when key does not match', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'j', ctrlKey: true }))
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('does not call handler when ctrl is required but not pressed', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', ctrlKey: false }))
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('handles shift modifier', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'n', ctrl: true, shift: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'n', ctrlKey: true, shiftKey: true }))
    })

    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('does not fire when shift is pressed but not expected', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', ctrlKey: true, shiftKey: true }))
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('handles alt modifier', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'p', alt: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'p', altKey: true }))
    })

    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('does not fire when alt is pressed but not expected', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', ctrlKey: true, altKey: true }))
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('prevents default on matching combo', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    const event = new KeyboardEvent('keydown', { key: 'k', ctrlKey: true })
    const preventSpy = vi.spyOn(event, 'preventDefault')

    act(() => {
      window.dispatchEvent(event)
    })

    expect(preventSpy).toHaveBeenCalled()
  })

  it('cleans up listener on unmount', () => {
    const handler = vi.fn()
    const { unmount } = renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    unmount()

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', ctrlKey: true }))
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('works with metaKey (Cmd on Mac)', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut({ key: 'k', ctrl: true }, handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }))
    })

    expect(handler).toHaveBeenCalledTimes(1)
  })
})

describe('useNavigationShortcuts', () => {
  afterEach(() => {
    mockNavigate.mockClear()
  })

  it('navigates to dashboard on Ctrl+1', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '1', ctrlKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/')
  })

  it('navigates to billing on Ctrl+2', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '2', ctrlKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/billing')
  })

  it('navigates to appointments on Ctrl+3', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '3', ctrlKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/appointments')
  })

  it('navigates to staff on Ctrl+4', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '4', ctrlKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/staff')
  })

  it('navigates to customers on Ctrl+5', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '5', ctrlKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/customers')
  })

  it('navigates to invoices on Ctrl+6', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '6', ctrlKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/invoices')
  })

  it('navigates to settings on Ctrl+7', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '7', ctrlKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/settings')
  })

  it('does not navigate on non-number keys', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'a', ctrlKey: true }))
    })

    expect(mockNavigate).not.toHaveBeenCalled()
  })

  it('does not navigate without ctrl', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '1', ctrlKey: false }))
    })

    expect(mockNavigate).not.toHaveBeenCalled()
  })

  it('does not navigate on Ctrl+8 (out of range)', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '8', ctrlKey: true }))
    })

    expect(mockNavigate).not.toHaveBeenCalled()
  })

  it('navigates with metaKey (Cmd on Mac)', () => {
    renderHook(() => useNavigationShortcuts(), { wrapper })

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '2', metaKey: true }))
    })

    expect(mockNavigate).toHaveBeenCalledWith('/billing')
  })
})

describe('useWindowDrag', () => {
  it('sets wails drag attributes on element', () => {
    const el = document.createElement('div')
    const ref = { current: el }

    renderHook(() => useWindowDrag(ref))

    expect(el.style.getPropertyValue('--wails-draggable')).toBe('drag')
    expect(el.hasAttribute('data-wails-drag')).toBe(true)
  })

  it('removes attributes on unmount', () => {
    const el = document.createElement('div')
    const ref = { current: el }

    const { unmount } = renderHook(() => useWindowDrag(ref))
    unmount()

    expect(el.style.getPropertyValue('--wails-draggable')).toBe('')
    expect(el.hasAttribute('data-wails-drag')).toBe(false)
  })

  it('handles null ref gracefully', () => {
    const ref = { current: null }
    expect(() => renderHook(() => useWindowDrag(ref))).not.toThrow()
  })
})

describe('useEscapeKey', () => {
  it('calls handler on Escape key', () => {
    const handler = vi.fn()
    renderHook(() => useEscapeKey(handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    })

    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('does not call handler for other keys', () => {
    const handler = vi.fn()
    renderHook(() => useEscapeKey(handler))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter' }))
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('does not attach listener when disabled', () => {
    const handler = vi.fn()
    renderHook(() => useEscapeKey(handler, false))

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('removes listener on unmount', () => {
    const handler = vi.fn()
    const { unmount } = renderHook(() => useEscapeKey(handler))
    unmount()

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    })

    expect(handler).not.toHaveBeenCalled()
  })
})

describe('useFocusTrap', () => {
  it('traps focus within container and focuses first element', () => {
    const container = document.createElement('div')
    const btn1 = document.createElement('button')
    btn1.textContent = 'First'
    const btn2 = document.createElement('button')
    btn2.textContent = 'Last'
    container.appendChild(btn1)
    container.appendChild(btn2)
    document.body.appendChild(container)

    const ref = { current: container }
    renderHook(() => useFocusTrap(ref))

    // First focusable element should be focused
    expect(document.activeElement).toBe(btn1)

    document.body.removeChild(container)
  })

  it('wraps focus forward from last to first on Tab', () => {
    const container = document.createElement('div')
    const btn1 = document.createElement('button')
    const btn2 = document.createElement('button')
    container.appendChild(btn1)
    container.appendChild(btn2)
    document.body.appendChild(container)

    const ref = { current: container }
    renderHook(() => useFocusTrap(ref))

    // Focus last element
    btn2.focus()
    expect(document.activeElement).toBe(btn2)

    // Tab on last element should wrap to first
    const event = new KeyboardEvent('keydown', { key: 'Tab', bubbles: true })
    Object.defineProperty(event, 'shiftKey', { value: false })
    const preventSpy = vi.spyOn(event, 'preventDefault')

    act(() => {
      container.dispatchEvent(event)
    })

    expect(preventSpy).toHaveBeenCalled()
    expect(document.activeElement).toBe(btn1)

    document.body.removeChild(container)
  })

  it('wraps focus backward from first to last on Shift+Tab', () => {
    const container = document.createElement('div')
    const btn1 = document.createElement('button')
    const btn2 = document.createElement('button')
    container.appendChild(btn1)
    container.appendChild(btn2)
    document.body.appendChild(container)

    const ref = { current: container }
    renderHook(() => useFocusTrap(ref))

    // Focus should already be on btn1
    expect(document.activeElement).toBe(btn1)

    // Shift+Tab on first element should wrap to last
    const event = new KeyboardEvent('keydown', { key: 'Tab', shiftKey: true, bubbles: true })
    const preventSpy = vi.spyOn(event, 'preventDefault')

    act(() => {
      container.dispatchEvent(event)
    })

    expect(preventSpy).toHaveBeenCalled()
    expect(document.activeElement).toBe(btn2)

    document.body.removeChild(container)
  })

  it('does not wrap backward when Shift+Tab and not on first element', () => {
    const container = document.createElement('div')
    const btn1 = document.createElement('button')
    const btn2 = document.createElement('button')
    const btn3 = document.createElement('button')
    container.appendChild(btn1)
    container.appendChild(btn2)
    container.appendChild(btn3)
    document.body.appendChild(container)

    const ref = { current: container }
    renderHook(() => useFocusTrap(ref))

    // Focus middle element
    btn2.focus()
    expect(document.activeElement).toBe(btn2)

    // Shift+Tab on middle element — should NOT prevent default
    const event = new KeyboardEvent('keydown', { key: 'Tab', shiftKey: true, bubbles: true })
    const preventSpy = vi.spyOn(event, 'preventDefault')

    act(() => {
      container.dispatchEvent(event)
    })

    expect(preventSpy).not.toHaveBeenCalled()

    document.body.removeChild(container)
  })

  it('does not prevent default when not on boundary element', () => {
    const container = document.createElement('div')
    const btn1 = document.createElement('button')
    const btn2 = document.createElement('button')
    const btn3 = document.createElement('button')
    container.appendChild(btn1)
    container.appendChild(btn2)
    container.appendChild(btn3)
    document.body.appendChild(container)

    const ref = { current: container }
    renderHook(() => useFocusTrap(ref))

    // Focus middle element
    btn2.focus()

    // Tab should not prevent default (not on last)
    const event = new KeyboardEvent('keydown', { key: 'Tab', bubbles: true })
    const preventSpy = vi.spyOn(event, 'preventDefault')

    act(() => {
      container.dispatchEvent(event)
    })

    expect(preventSpy).not.toHaveBeenCalled()

    document.body.removeChild(container)
  })

  it('ignores non-Tab keys', () => {
    const container = document.createElement('div')
    const btn1 = document.createElement('button')
    container.appendChild(btn1)
    document.body.appendChild(container)

    const ref = { current: container }
    renderHook(() => useFocusTrap(ref))

    const event = new KeyboardEvent('keydown', { key: 'Enter', bubbles: true })
    const preventSpy = vi.spyOn(event, 'preventDefault')

    act(() => {
      container.dispatchEvent(event)
    })

    expect(preventSpy).not.toHaveBeenCalled()

    document.body.removeChild(container)
  })

  it('does not trap when inactive', () => {
    const container = document.createElement('div')
    const btn1 = document.createElement('button')
    container.appendChild(btn1)
    document.body.appendChild(container)

    const ref = { current: container }
    renderHook(() => useFocusTrap(ref, false))

    // Should not auto-focus
    expect(document.activeElement).not.toBe(btn1)

    document.body.removeChild(container)
  })

  it('handles null ref gracefully', () => {
    const ref = { current: null }
    expect(() => renderHook(() => useFocusTrap(ref))).not.toThrow()
  })

  it('handles container with no focusable elements', () => {
    const container = document.createElement('div')
    container.innerHTML = '<span>no focusable</span>'
    document.body.appendChild(container)

    const ref = { current: container }
    expect(() => renderHook(() => useFocusTrap(ref))).not.toThrow()

    document.body.removeChild(container)
  })
})

describe('useDebouncedValue', () => {
  it('returns initial value immediately', () => {
    vi.useFakeTimers()
    const { result } = renderHook(() => useDebouncedValue('hello'))
    expect(result.current).toBe('hello')
    vi.useRealTimers()
  })

  it('debounces value changes', () => {
    vi.useFakeTimers()
    const { result, rerender } = renderHook(({ value }) => useDebouncedValue(value, 200), {
      initialProps: { value: 'a' },
    })

    rerender({ value: 'b' })
    expect(result.current).toBe('a') // not yet updated

    act(() => {
      vi.advanceTimersByTime(200)
    })
    expect(result.current).toBe('b') // updated after delay

    vi.useRealTimers()
  })

  it('uses default delay of 200ms', () => {
    vi.useFakeTimers()
    const { result, rerender } = renderHook(({ value }) => useDebouncedValue(value), {
      initialProps: { value: 'initial' },
    })

    rerender({ value: 'changed' })
    expect(result.current).toBe('initial')

    act(() => {
      vi.advanceTimersByTime(199)
    })
    expect(result.current).toBe('initial')

    act(() => {
      vi.advanceTimersByTime(1)
    })
    expect(result.current).toBe('changed')

    vi.useRealTimers()
  })
})
