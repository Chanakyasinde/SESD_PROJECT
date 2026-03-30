import { api } from './api'
import type { ApiEnvelope, PaginatedResponse, User } from '../types'

export const userService = {
  async listCustomers(params?: Record<string, string | number>) {
    const res = await api.get<ApiEnvelope<PaginatedResponse<User>>>('/users/customers', { params })
    return res.data.data
  },
}
