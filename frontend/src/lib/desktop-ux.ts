/**
 * SalonFlow Desktop UX Utilities
 * ================================
 * Patterns that make the app feel native on desktop:
 * - Keyboard navigation
 * - Instant tooltip delays
 * - Window drag region
 * - Context menus
 */

import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

// ─── Keyboard Shortcut System ────────────────────────────────────────────────

type KeyCombo = {
  key: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
}

type ShortcutHandler = () => void

/**
 * Register a global keyboard shortcut.
 * Handles Ctrl/Cmd abstraction for cross-platform.
 */
export function useKeyboardShortcut(
  combo: KeyCombo,
  handler: ShortcutHandler,
  deps: unknown[] = []
) {
  useEffect(() => {
    const listener = (e: KeyboardEvent) => {
      const ctrlMatch = combo.ctrl ? (e.ctrlKey || e.metaKey) : true
      const shiftMatch = combo.shift ? e.shiftKey : !e.shiftKey
      const altMatch = combo.alt ? e.altKey : !e.altKey

      if (ctrlMatch && shiftMatch && altMatch && e.key.toLowerCase() === combo.key.toLowerCase()) {
        e.preventDefault()
        handler()
      }
    }

    window.addEventListener('keydown', listener)
    return () => window.removeEventListener('keydown', listener)
  }, [combo.key, combo.ctrl, combo.shift, combo.alt, handler, ...deps])
}

/**
 * Quick navigation shortcuts (Ctrl+1..9 for sidebar sections).
 */
export function useNavigationShortcuts() {
  const navigate = useNavigate()

  const routes = [
    '/',           // Ctrl+1 → Dashboard
    '/billing',    // Ctrl+2 → Billing
    '/appointments', // Ctrl+3 → Appointments
    '/staff',      // Ctrl+4 → Staff
    '/customers',  // Ctrl+5 → Customers
    '/invoices',   // Ctrl+6 → Invoices
    '/settings',   // Ctrl+7 → Settings
  ]

  useEffect(() => {
    const listener = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key >= '1' && e.key <= '7') {
        e.preventDefault()
        const index = parseInt(e.key) - 1
        navigate(routes[index])
      }
    }

    window.addEventListener('keydown', listener)
    return () => window.removeEventListener('keydown', listener)
  }, [navigate])
}

// ─── Desktop Window Utilities ────────────────────────────────────────────────

/**
 * Prevents the default drag behavior on the window chrome
 * and enables wails-style window dragging on specified regions.
 */
export function useWindowDrag(elementRef: React.RefObject<HTMLElement | null>) {
  useEffect(() => {
    const el = elementRef.current
    if (!el) return

    el.style.setProperty('--wails-draggable', 'drag')
    el.setAttribute('data-wails-drag', '')

    return () => {
      el.style.removeProperty('--wails-draggable')
      el.removeAttribute('data-wails-drag')
    }
  }, [elementRef])
}

// ─── Escape Key Hook (for modals, popovers) ─────────────────────────────────

export function useEscapeKey(handler: () => void, enabled = true) {
  useEffect(() => {
    if (!enabled) return

    const listener = (e: KeyboardEvent) => {
      if (e.key === 'Escape') handler()
    }

    window.addEventListener('keydown', listener)
    return () => window.removeEventListener('keydown', listener)
  }, [handler, enabled])
}

// ─── Focus Trap (for modal dialogs) ──────────────────────────────────────────

export function useFocusTrap(containerRef: React.RefObject<HTMLElement | null>, active = true) {
  useEffect(() => {
    if (!active || !containerRef.current) return

    const container = containerRef.current
    const focusableElements = container.querySelectorAll<HTMLElement>(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    )

    if (focusableElements.length === 0) return

    const first = focusableElements[0]!
    const last = focusableElements[focusableElements.length - 1]!

    const trapHandler = (e: KeyboardEvent) => {
      if (e.key !== 'Tab') return

      if (e.shiftKey) {
        if (document.activeElement === first) {
          e.preventDefault()
          last.focus()
        }
      } else {
        if (document.activeElement === last) {
          e.preventDefault()
          first.focus()
        }
      }
    }

    container.addEventListener('keydown', trapHandler)
    first.focus()

    return () => container.removeEventListener('keydown', trapHandler)
  }, [containerRef, active])
}

// ─── Debounced Search Hook ───────────────────────────────────────────────────

export function useDebouncedValue<T>(value: T, delay = 200): T {
  const [debounced, setDebounced] = useState(value)

  useEffect(() => {
    const timer = setTimeout(() => setDebounced(value), delay)
    return () => clearTimeout(timer)
  }, [value, delay])

  return debounced
}
