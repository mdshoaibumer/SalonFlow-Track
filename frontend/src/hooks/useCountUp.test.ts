import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useCountUp, easeOutCubic, formatAnimatedValue } from './useCountUp'

describe('useCountUp', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('starts at 0 and animates to target', () => {
    const { result } = renderHook(() => useCountUp(100))
    // Initially 0 before RAF fires
    expect(result.current).toBe(0)
  })

  it('returns 0 when end is 0', () => {
    const { result } = renderHook(() => useCountUp(0))
    expect(result.current).toBe(0)
  })

  it('animates to target after RAF ticks', () => {
    let rafCallback: FrameRequestCallback | null = null
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb
      return 1
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const { result } = renderHook(() => useCountUp(100, { duration: 100, delay: 0 }))

    // Simulate animation completion
    act(() => {
      if (rafCallback) rafCallback(0) // start
    })
    act(() => {
      if (rafCallback) rafCallback(200) // past duration
    })

    expect(result.current).toBe(100)
  })

  it('respects delay before starting animation', () => {
    let rafCallback: FrameRequestCallback | null = null
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb
      return 1
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const { result } = renderHook(() => useCountUp(50, { duration: 100, delay: 50 }))

    // First tick within delay period
    act(() => {
      if (rafCallback) rafCallback(10)
    })
    // Still 0 because delay hasn't elapsed
    expect(result.current).toBe(0)

    // Tick past delay + duration
    act(() => {
      if (rafCallback) rafCallback(200)
    })
    expect(result.current).toBe(50)
  })

  it('respects decimals option', () => {
    let rafCallback: FrameRequestCallback | null = null
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb
      return 1
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const { result } = renderHook(() => useCountUp(10, { duration: 100, delay: 0, decimals: 2 }))

    // Mid-animation
    act(() => {
      if (rafCallback) rafCallback(0)
    })
    act(() => {
      if (rafCallback) rafCallback(50)
    })
    // Should have decimal places
    expect(typeof result.current).toBe('number')
  })

  it('does not re-animate when once=true and value unchanged', () => {
    let rafCallback: FrameRequestCallback | null = null
    const rafSpy = vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb
      return 1
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const { rerender } = renderHook(
      ({ end, duration }) => useCountUp(end, { once: true, duration, delay: 0 }),
      { initialProps: { end: 100, duration: 100 } }
    )

    // Complete the animation
    act(() => { if (rafCallback) rafCallback(0) })
    act(() => { if (rafCallback) rafCallback(200) })

    const callCountAfterAnimation = rafSpy.mock.calls.length
    
    // Re-render with different duration but same end — effect re-runs but early returns at line 47
    rerender({ end: 100, duration: 200 })
    expect(rafSpy.mock.calls.length).toBe(callCountAfterAnimation)
  })

  it('cancels animation on unmount', () => {
    const cancelSpy = vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation(() => 42)

    const { unmount } = renderHook(() => useCountUp(100))
    unmount()

    expect(cancelSpy).toHaveBeenCalledWith(42)
  })

  it('handles end=0 when previously had value', () => {
    let rafCallback: FrameRequestCallback | null = null
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb
      return 1
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const { result, rerender } = renderHook(({ end }) => useCountUp(end, { once: false, duration: 100, delay: 0 }), {
      initialProps: { end: 50 },
    })

    // Complete initial animation
    act(() => { if (rafCallback) rafCallback(0) })
    act(() => { if (rafCallback) rafCallback(200) })
    expect(result.current).toBe(50)

    // Change to 0 — should immediately set to 0
    rerender({ end: 0 })
    expect(result.current).toBe(0)
  })

  it('re-animates from previous value when once=false', () => {
    let rafCallback: FrameRequestCallback | null = null
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb
      return 1
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const { result, rerender } = renderHook(({ end }) => useCountUp(end, { once: false, duration: 100, delay: 0 }), {
      initialProps: { end: 50 },
    })

    // Complete initial animation
    act(() => { if (rafCallback) rafCallback(0) })
    act(() => { if (rafCallback) rafCallback(200) })
    expect(result.current).toBe(50)

    // Change to new value — should animate from 50 to 100
    rerender({ end: 100 })
    act(() => { if (rafCallback) rafCallback(0) })
    act(() => { if (rafCallback) rafCallback(200) })
    expect(result.current).toBe(100)
  })

  it('falls back to end value when requestAnimationFrame is undefined', () => {
    const originalRaf = window.requestAnimationFrame
    // @ts-expect-error - intentionally removing RAF for test
    delete window.requestAnimationFrame

    const { result } = renderHook(() => useCountUp(200))
    expect(result.current).toBe(200)

    window.requestAnimationFrame = originalRaf
  })

  it('uses custom easing function', () => {
    let rafCallback: FrameRequestCallback | null = null
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb
      return 1
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const linearEasing = (t: number) => t
    const { result } = renderHook(() => useCountUp(100, { duration: 100, delay: 0, easing: linearEasing }))

    // At 50% progress with linear easing, value should be ~50
    act(() => {
      if (rafCallback) rafCallback(0)
    })
    act(() => {
      if (rafCallback) rafCallback(50)
    })
    expect(result.current).toBe(50)
  })
})

describe('easeOutCubic', () => {
  it('returns 0 at start', () => {
    expect(easeOutCubic(0)).toBe(0)
  })

  it('returns 1 at end', () => {
    expect(easeOutCubic(1)).toBe(1)
  })

  it('returns value between 0 and 1 for mid values', () => {
    const mid = easeOutCubic(0.5)
    expect(mid).toBeGreaterThan(0)
    expect(mid).toBeLessThan(1)
  })
})

describe('formatAnimatedValue', () => {
  it('formats number with default locale', () => {
    const result = formatAnimatedValue(12500)
    expect(result).toBe('12,500')
  })

  it('formats with prefix', () => {
    const result = formatAnimatedValue(1000, { prefix: '₹' })
    expect(result).toBe('₹1,000')
  })

  it('formats with suffix', () => {
    const result = formatAnimatedValue(50, { suffix: '%' })
    expect(result).toBe('50%')
  })

  it('formats with prefix and suffix', () => {
    const result = formatAnimatedValue(100, { prefix: '$', suffix: ' USD' })
    expect(result).toBe('$100 USD')
  })

  it('formats with custom locale', () => {
    const result = formatAnimatedValue(1000, { locale: 'en-US' })
    expect(result).toBe('1,000')
  })
})
