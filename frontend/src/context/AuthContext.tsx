import { createContext, useEffect, useMemo, useState } from 'react'
import type { PropsWithChildren } from 'react'
import type { User } from '../types'
import { authService } from '../services/authService'

interface AuthContextShape {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  signup: (name: string, email: string, password: string) => Promise<void>
  logout: () => void
}

export const AuthContext = createContext<AuthContextShape | null>(null)

const TOKEN_KEY = 'ims_token'
const USER_KEY = 'ims_user'

export function AuthProvider({ children }: PropsWithChildren) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const savedToken = localStorage.getItem(TOKEN_KEY)
    const savedUser = localStorage.getItem(USER_KEY)

    if (savedToken && savedUser) {
      try {
        const parsedUser = JSON.parse(savedUser) as User
        if (parsedUser && typeof parsedUser === 'object' && parsedUser.id) {
          setToken(savedToken)
          setUser(parsedUser)
        } else {
          localStorage.removeItem(TOKEN_KEY)
          localStorage.removeItem(USER_KEY)
        }
      } catch {
        localStorage.removeItem(TOKEN_KEY)
        localStorage.removeItem(USER_KEY)
      }
    }
    setIsLoading(false)
  }, [])

  const login = async (email: string, password: string) => {
    const result = await authService.login({ email, password })
    localStorage.setItem(TOKEN_KEY, result.token)
    localStorage.setItem(USER_KEY, JSON.stringify(result.user))
    setToken(result.token)
    setUser(result.user)
  }

  const signup = async (name: string, email: string, password: string) => {
    const result = await authService.signup({ name, email, password, role: 'customer' })
    localStorage.setItem(TOKEN_KEY, result.token)
    localStorage.setItem(USER_KEY, JSON.stringify(result.user))
    setToken(result.token)
    setUser(result.user)
  }

  const logout = () => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    setToken(null)
    setUser(null)
  }

  const value = useMemo(
    () => ({
      user,
      token,
      isAuthenticated: Boolean(token && user),
      isLoading,
      login,
      signup,
      logout,
    }),
    [user, token, isLoading],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
