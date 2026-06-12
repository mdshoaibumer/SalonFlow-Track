import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import type { SessionInfo, LoginInput } from '@/types/auth'
import * as authService from '@/services/auth'

interface AuthContextType {
  user: SessionInfo | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (input: LoginInput) => Promise<void>
  logout: () => Promise<void>
  refreshSession: () => Promise<void>
  hasPermission: (permission: string) => boolean
  hasAnyPermission: (permissions: string[]) => boolean
  hasRole: (role: string) => boolean
}

const AuthContext = createContext<AuthContextType | null>(null)

const TOKEN_KEY = 'salonflow_session_token'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<SessionInfo | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Restore session on mount
  useEffect(() => {
    restoreSession()
  }, [])

  const restoreSession = async () => {
    try {
      // Try to get the current session from backend
      const savedToken = localStorage.getItem(TOKEN_KEY)
      if (savedToken) {
        await authService.setToken(savedToken)
      }
      const session = await authService.getCurrentSession()
      setUser(session)
    } catch {
      setUser(null)
      localStorage.removeItem(TOKEN_KEY)
    } finally {
      setIsLoading(false)
    }
  }

  const login = useCallback(async (input: LoginInput) => {
    const result = await authService.login(input)
    setUser(result.user)
    if (input.remember_me) {
      localStorage.setItem(TOKEN_KEY, result.token)
    } else {
      // Store in sessionStorage for the current session
      sessionStorage.setItem(TOKEN_KEY, result.token)
    }
  }, [])

  const logout = useCallback(async () => {
    try {
      await authService.logout()
    } finally {
      setUser(null)
      localStorage.removeItem(TOKEN_KEY)
      sessionStorage.removeItem(TOKEN_KEY)
    }
  }, [])

  const refreshSession = useCallback(async () => {
    try {
      const session = await authService.getCurrentSession()
      setUser(session)
    } catch {
      setUser(null)
    }
  }, [])

  const hasPermission = useCallback((permission: string): boolean => {
    if (!user) return false
    return user.permissions.includes(permission)
  }, [user])

  const hasAnyPermission = useCallback((permissions: string[]): boolean => {
    if (!user) return false
    return permissions.some(p => user.permissions.includes(p))
  }, [user])

  const hasRole = useCallback((role: string): boolean => {
    if (!user) return false
    return user.roles.includes(role)
  }, [user])

  return (
    <AuthContext.Provider value={{
      user,
      isAuthenticated: !!user,
      isLoading,
      login,
      logout,
      refreshSession,
      hasPermission,
      hasAnyPermission,
      hasRole,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
