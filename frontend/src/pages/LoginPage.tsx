import { useState } from 'react'
import type { FormEvent } from 'react'
import { isAxiosError } from 'axios'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export function LoginPage() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault()
    try {
      await login(email, password)
      navigate('/dashboard')
    } catch (err) {
      if (isAxiosError(err)) {
        const message = (err.response?.data as { message?: string } | undefined)?.message
        setError(message ?? 'Login failed. Check your credentials.')
        return
      }
      setError('Login failed. Check your credentials.')
    }
  }

  return (
    <div className="auth-page">
      <form className="auth-card" onSubmit={onSubmit}>
        <h2>Welcome Back</h2>
        <p>Sign in to manage inventory operations.</p>
        <label>
          Email
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
        </label>
        <label>
          Password
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
        </label>
        {error && <div className="error-text">{error}</div>}
        <button type="submit">Login</button>
        <div className="credential-box">
          <strong>Evaluation Credentials</strong>
          <p>Admin: admin@inventory.local / Admin@12345</p>
          <p>Staff: staff@inventory.local / Staff@12345</p>
        </div>
        <small>
          No account? <Link to="/signup">Create one</Link>
        </small>
      </form>
    </div>
  )
}
