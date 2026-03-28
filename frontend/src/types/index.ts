export type Role = 'admin' | 'staff' | 'customer'

export interface User {
  id: string
  name: string
  email: string
  role: Role
  createdAt: string
  updatedAt: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface Product {
  id: string
  name: string
  description: string
  category: string
  price: number
  stockQuantity: number
  lowStockThreshold: number
  version: number
  createdAt: string
  updatedAt: string
}

export interface OrderItem {
  productId: string
  productName: string
  quantity: number
  unitPrice: number
  subTotal: number
}

export interface Order {
  id: string
  customerId: string
  createdBy: string
  status: 'pending' | 'confirmed' | 'shipped' | 'delivered' | 'cancelled'
  items: OrderItem[]
  totalAmount: number
  createdAt: string
  updatedAt: string
}

export type OrderStatus = Order['status']

export const ORDER_TRANSITIONS: Record<OrderStatus, OrderStatus[]> = {
  pending: ['confirmed', 'cancelled'],
  confirmed: ['shipped', 'cancelled'],
  shipped: ['delivered'],
  delivered: [],
  cancelled: [],
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface ApiEnvelope<T> {
  success: boolean
  message: string
  data: T
}

export interface ProductFormInput {
  name: string
  description: string
  category: string
  price: number
  stockQuantity: number
  lowStockThreshold: number
}

export interface CreateOrderInput {
  customerId: string
  items: Array<{
    productId: string
    quantity: number
  }>
}
