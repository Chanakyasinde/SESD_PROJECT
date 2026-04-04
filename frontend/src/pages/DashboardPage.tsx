import { useEffect, useState } from 'react'
import { orderService } from '../services/orderService'
import { productService } from '../services/productService'

export function DashboardPage() {
  const [stats, setStats] = useState({ totalProducts: 0, lowStockCount: 0, totalOrders: 0 })

  useEffect(() => {
    Promise.all([
      productService.list({ page: 1, pageSize: 1 }),
      productService.list({ page: 1, pageSize: 100, lowStock: 'true' }),
      orderService.list({ page: 1, pageSize: 1 }),
    ])
      .then(([products, lowStock, orders]) => {
        setStats({
          totalProducts: products.total,
          lowStockCount: lowStock.total,
          totalOrders: orders.total,
        })
      })
      .catch(() => {
        setStats({ totalProducts: 0, lowStockCount: 0, totalOrders: 0 })
      })
  }, [])

  return (
    <section>
      <h2>Operational Dashboard</h2>
      <div className="stat-grid">
        <article>
          <h3>{stats.totalProducts}</h3>
          <span>Total Products</span>
        </article>
        <article>
          <h3>{stats.lowStockCount}</h3>
          <span>Low Stock Alerts</span>
        </article>
        <article>
          <h3>{stats.totalOrders}</h3>
          <span>Total Orders</span>
        </article>
      </div>
    </section>
  )
}
