/**
 * useCountUp Hook
 * ================
 * Animates a number from 0 to target value with easing.
 * Creates the "counting up" effect seen in premium dashboards.
 */

import { useEffect, useRef, useState } from 'react'

interface UseCountUpOptions {
  /** Duration of the count animation in ms */
  duration?: number
  /** Delay before starting in ms */
  delay?: number
  /** Easing function */
  easing?: (t: number) => number
  /** Number of decimal places */
  decimals?: number
  /** Only animate on first render */
  once?: boolean
}

// Smooth ease-out-cubic (exported for custom use)
export const easeOutCubic = (t: number) => 1 - Math.pow(1 - t, 3)

// Smooth ease-out-expo (more dramatic)
const easeOutExpo = (t: number) => (t === 1 ? 1 : 1 - Math.pow(2, -10 * t))

export function useCountUp(
  end: number,
  options: UseCountUpOptions = {}
): number {
  const {
    duration = 800,
    delay = 0,
    easing = easeOutExpo,
    decimals = 0,
    once = true,
  } = options

  const [count, setCount] = useState(0)
  const prevEnd = useRef(end)
  const hasAnimated = useRef(false)
  const rafRef = useRef<number>(0)

  useEffect(() => {
    if (once && hasAnimated.current && prevEnd.current === end) return
    
    const startValue = once ? 0 : prevEnd.current
    prevEnd.current = end

    if (end === 0) {
      setCount(0)
      return
    }

    // Skip animation if RAF is not available or in test environments
    if (typeof requestAnimationFrame === 'undefined') {
      setCount(end)
      return
    }

    let startTime: number | null = null
    hasAnimated.current = true

    const tick = (timestamp: number) => {
      if (startTime === null) startTime = timestamp
      const elapsed = timestamp - startTime

      if (elapsed < delay) {
        rafRef.current = requestAnimationFrame(tick)
        return
      }

      const progress = Math.min((elapsed - delay) / duration, 1)
      const easedProgress = easing(progress)
      const current = startValue + (end - startValue) * easedProgress
      
      setCount(Number(current.toFixed(decimals)))

      if (progress < 1) {
        rafRef.current = requestAnimationFrame(tick)
      }
    }

    rafRef.current = requestAnimationFrame(tick)
    
    return () => {
      cancelAnimationFrame(rafRef.current)
    }
  }, [end, duration, delay, easing, decimals, once])

  return count
}

/**
 * Format a number with Indian locale and optional prefix/suffix.
 * Used in combination with useCountUp for animated currency displays.
 */
export function formatAnimatedValue(
  value: number,
  options: { prefix?: string; suffix?: string; locale?: string } = {}
): string {
  const { prefix = '', suffix = '', locale = 'en-IN' } = options
  return `${prefix}${value.toLocaleString(locale)}${suffix}`
}
