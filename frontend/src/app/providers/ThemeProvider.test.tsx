import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/react'
import { ThemeProvider, useTheme } from './ThemeProvider'

beforeAll(() => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation((query: string) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })),
  })
})

function TestComponent() {
  const { theme, setTheme } = useTheme()
  return (
    <div>
      <span data-testid="theme">{theme}</span>
      <button onClick={() => setTheme('light')}>Set Light</button>
    </div>
  )
}

describe('ThemeProvider', () => {
  it('provides light theme by default', () => {
    const { getByTestId } = render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )
    expect(getByTestId('theme').textContent).toBe('light')
  })

  it('adds light class to document root', () => {
    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )
    expect(document.documentElement.classList.contains('light')).toBe(true)
  })

  it('useTheme works outside provider with default context', () => {
    function Orphan() {
      const ctx = useTheme()
      return <span data-testid="theme-val">{ctx.theme}</span>
    }
    const { getByTestId } = render(<Orphan />)
    expect(getByTestId('theme-val').textContent).toBe('light')
  })

  it('default context setTheme returns null', () => {
    let setThemeFn: (t: 'light') => void = () => {}
    function Capture() {
      const { setTheme } = useTheme()
      setThemeFn = setTheme
      return null
    }
    render(<Capture />)
    expect(setThemeFn('light')).toBeNull()
  })

  it('setTheme inside provider is a no-op', () => {
    const { getByText, getByTestId } = render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )
    getByText('Set Light').click()
    expect(getByTestId('theme').textContent).toBe('light')
  })

  it('accepts storageKey and defaultTheme props', () => {
    const { getByTestId } = render(
      <ThemeProvider defaultTheme="light" storageKey="custom-key">
        <TestComponent />
      </ThemeProvider>
    )
    expect(getByTestId('theme').textContent).toBe('light')
  })
})
