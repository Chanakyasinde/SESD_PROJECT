import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export function ProtectedRoute() {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return <div className="centered">Loading session...</div>
  }

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}
