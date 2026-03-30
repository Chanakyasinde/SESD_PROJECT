import { api } from './api'
import type { ApiEnvelope, CreateOrderInput, Order, PaginatedResponse } from '../types'

export const orderService = {
  async list(params?: Record<string, string | number>) {
    const res = await api.get<ApiEnvelope<PaginatedResponse<Order>>>('/orders', { params })
    return res.data.data
  },

  async create(payload: CreateOrderInput) {
    const res = await api.post<ApiEnvelope<Order>>('/orders', payload)
    return res.data.data
  },

  async updateStatus(id: string, status: Order['status']) {
    await api.patch(`/orders/${id}/status`, { status })
  },
}
