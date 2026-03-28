import { api } from './api'
import type { ApiEnvelope, AuthResponse } from '../types'

export const authService = {
  async login(payload: { email: string; password: string }) {
    const res = await api.post<ApiEnvelope<AuthResponse>>('/auth/login', payload)
    return res.data.data
  },

  async signup(payload: { name: string; email: string; password: string; role?: string }) {
    const res = await api.post<ApiEnvelope<AuthResponse>>('/auth/register', payload)
    return res.data.data
  },
}
