import { api } from './api'
import type { ApiEnvelope, PaginatedResponse, Product, ProductFormInput } from '../types'

export const productService = {
  async list(params?: Record<string, string | number>) {
    const res = await api.get<ApiEnvelope<PaginatedResponse<Product>>>('/products', { params })
    return res.data.data
  },

  async create(payload: ProductFormInput) {
    const res = await api.post<ApiEnvelope<Product>>('/products', payload)
    return res.data.data
  },

  async update(id: string, payload: ProductFormInput) {
    const res = await api.put<ApiEnvelope<Product>>(`/products/${id}`, payload)
    return res.data.data
  },

  async remove(id: string) {
    await api.delete(`/products/${id}`)
  },
}
