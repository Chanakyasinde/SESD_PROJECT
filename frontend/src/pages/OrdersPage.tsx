import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { isAxiosError } from 'axios'
import { DataTable } from '../components/DataTable'
import { useAuth } from '../hooks/useAuth'
import { orderService } from '../services/orderService'
import { productService } from '../services/productService'
import { userService } from '../services/userService'
import type { Order, OrderStatus, Product, User } from '../types'
import { ORDER_TRANSITIONS } from '../types'
import { formatINR } from '../utils/money'

export function OrdersPage() {
  const { user } = useAuth()
  const canManageOrders = user?.role === 'admin' || user?.role === 'staff'
  const canCreateOrders = user?.role === 'admin' || user?.role === 'staff' || user?.role === 'customer'
  const isAdminOrStaff = user?.role === 'admin' || user?.role === 'staff'

  const [orders, setOrders] = useState<Order[]>([])
  const [customers, setCustomers] = useState<User[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [customerId, setCustomerId] = useState('')
  const [productId, setProductId] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [error, setError] = useState('')

  const loadOrders = async () => {
    const result = await orderService.list({ page: 1, pageSize: 50, sortBy: 'updatedAt', sortDir: -1 })
    setOrders(result.items)
  }

  useEffect(() => {
    loadOrders().catch(() => setOrders([]))
  }, [])

  useEffect(() => {
    if (!canCreateOrders) {
      return
    }

    if (isAdminOrStaff) {
      userService
        .listCustomers({ page: 1, pageSize: 200, sortBy: 'createdAt', sortDir: -1 })
        .then((result) => setCustomers(result.items))
        .catch(() => setCustomers([]))
    }

    productService
      .list({ page: 1, pageSize: 200, sortBy: 'updatedAt', sortDir: -1 })
      .then((result) => setProducts(result.items))
      .catch(() => setProducts([]))
  }, [canCreateOrders, isAdminOrStaff])

  const onCreate = async (event: FormEvent) => {
    event.preventDefault()
    if (!canCreateOrders) {
      return
    }

    const resolvedCustomerId = isAdminOrStaff ? customerId : user?.id ?? ''
    if (!resolvedCustomerId) {
      setError('Customer is required.')
      return
    }

    setError('')
    try {
      const created = await orderService.create({
        customerId: resolvedCustomerId,
        items: [{ productId, quantity }],
      })
      setOrders((prev) => [created, ...prev])
      if (isAdminOrStaff) {
        setCustomerId('')
      }
      setProductId('')
      setQuantity(1)
    } catch (err) {
      if (isAxiosError(err)) {
        setError((err.response?.data as { message?: string } | undefined)?.message ?? 'Could not create order.')
        return
      }
      setError('Could not create order.')
    }
  }

  const updateOrderStatus = async (orderId: string, nextStatus: OrderStatus) => {
    const previous = [...orders]
    setOrders((prev) =>
      prev.map((order) =>
        order.id === orderId
          ? {
              ...order,
              status: nextStatus,
              updatedAt: new Date().toISOString(),
            }
          : order,
      ),
    )
    setError('')

    try {
      await orderService.updateStatus(orderId, nextStatus)
    } catch (err) {
      setOrders(previous)
      if (isAxiosError(err)) {
        setError((err.response?.data as { message?: string } | undefined)?.message ?? 'Could not update order status.')
        return
      }
      setError('Could not update order status.')
    }
  }

  return (
    <section>
      <h2>Orders</h2>
      {user?.role === 'customer' && <p className="info-text">You can create orders for your own account and view only your orders.</p>}
      {!canCreateOrders && <p className="info-text">You can view orders, but you cannot create them.</p>}

      {canCreateOrders && (
        <form className="inline-form" onSubmit={onCreate}>
          {isAdminOrStaff && (
            <select value={customerId} onChange={(e) => setCustomerId(e.target.value)} required>
              <option value="" disabled>
                Select Customer...
              </option>
              {customers.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name} — {c.email}
                </option>
              ))}
            </select>
          )}

          <select value={productId} onChange={(e) => setProductId(e.target.value)} required>
            <option value="" disabled>
              Select Product...
            </option>
            {products.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </select>
          <input placeholder="Quantity" type="number" min="1" value={quantity} onChange={(e) => setQuantity(Number(e.target.value))} required />
          <button type="submit">Create Order</button>
        </form>
      )}

      {error && <p className="error-text">{error}</p>}

      <DataTable
        rows={orders}
        columns={[
          { key: 'id', label: 'Order ID' },
          {
            key: 'status',
            label: 'Status',
            render: (value, row) => {
              const current = value as OrderStatus
              const nextStates = ORDER_TRANSITIONS[current]

                if (!canManageOrders || nextStates.length === 0) {
                return <span className="status-pill">{current}</span>
              }

              return (
                <div className="status-controls">
                  <span className="status-pill">{current}</span>
                  <select
                    defaultValue=""
                    onChange={(e) => {
                      const selected = e.target.value as OrderStatus
                      if (!selected) {
                        return
                      }
                      void updateOrderStatus(row.id, selected)
                      e.currentTarget.value = ''
                    }}
                  >
                    <option value="" disabled>
                      Transition...
                    </option>
                    {nextStates.map((status) => (
                      <option key={status} value={status}>
                        {status}
                      </option>
                    ))}
                  </select>
                </div>
              )
            },
          },
          { key: 'totalAmount', label: 'Total', render: (v) => formatINR(v) },
          { key: 'customerId', label: 'Customer ID' },
        ]}
      />
    </section>
  )
}
