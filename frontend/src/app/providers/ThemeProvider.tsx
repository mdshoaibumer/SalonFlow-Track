import { createContext, useContext, useEffect, useState } from 'react'

type Theme = 'light'

interface ThemeProviderState {
  theme: Theme
  setTheme: (theme: Theme) => void
}

const ThemeProviderContext = createContext<ThemeProviderState>({
  theme: 'light',
  setTheme: () => null,
})

interface ThemeProviderProps {
  children: React.ReactNode
  defaultTheme?: Theme
  storageKey?: string
}

export function ThemeProvider({
  children,
  defaultTheme: _defaultTheme = 'light',
  storageKey: _storageKey = 'salonflow-theme',
}: ThemeProviderProps) {
  const [theme] = useState<Theme>('light')

  useEffect(() => {
    const root = window.document.documentElement
    root.classList.remove('light', 'dark')
    root.classList.add('light')
  }, [theme])

  const value = {
    theme,
    setTheme: (_theme: Theme) => {
      // Light theme only - no theme switching
    },
  }

  return (
    <ThemeProviderContext.Provider value={value}>
      {children}
    </ThemeProviderContext.Provider>
  )
}

export const useTheme = () => {
  return useContext(ThemeProviderContext)
}
