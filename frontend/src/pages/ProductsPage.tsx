import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { isAxiosError } from 'axios'
import { DataTable } from '../components/DataTable'
import { useAuth } from '../hooks/useAuth'
import { productService } from '../services/productService'
import type { Product, ProductFormInput } from '../types'
import { formatINR } from '../utils/money'

const initialForm: ProductFormInput = {
  name: '',
  description: '',
  category: '',
  price: 0,
  stockQuantity: 0,
  lowStockThreshold: 0,
}

export function ProductsPage() {
  const { user } = useAuth()
  const canManageProducts = user?.role === 'admin'

  const [products, setProducts] = useState<Product[]>([])
  const [form, setForm] = useState<ProductFormInput>(initialForm)
  const [editProductId, setEditProductId] = useState<string | null>(null)
  const [editForm, setEditForm] = useState<ProductFormInput>(initialForm)
  const [error, setError] = useState('')

  const loadProducts = async () => {
    const result = await productService.list({ page: 1, pageSize: 50, sortBy: 'updatedAt', sortDir: -1 })
    setProducts(result.items)
  }

  useEffect(() => {
    loadProducts().catch(() => setProducts([]))
  }, [])

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault()
    if (!canManageProducts) {
      return
    }

    setError('')
    try {
      const created = await productService.create(form)
      setProducts((prev) => [created, ...prev])
      setForm(initialForm)
    } catch (err) {
      if (isAxiosError(err)) {
        setError((err.response?.data as { message?: string } | undefined)?.message ?? 'Could not create product.')
        return
      }
      setError('Could not create product.')
    }
  }

  const beginEdit = (product: Product) => {
    setEditProductId(product.id)
    setEditForm({
      name: product.name,
      description: product.description,
      category: product.category,
      price: product.price,
      stockQuantity: product.stockQuantity,
      lowStockThreshold: product.lowStockThreshold,
    })
    setError('')
  }

  const cancelEdit = () => {
    setEditProductId(null)
    setEditForm(initialForm)
  }

  const saveEdit = async (event: FormEvent) => {
    event.preventDefault()
    if (!editProductId) {
      return
    }

    const current = products.find((item) => item.id === editProductId)
    if (!current) {
      return
    }

    const optimistic: Product = {
      ...current,
      ...editForm,
      updatedAt: new Date().toISOString(),
    }
    setProducts((prev) => prev.map((item) => (item.id === editProductId ? optimistic : item)))
    setError('')

    try {
      const updated = await productService.update(editProductId, editForm)
      setProducts((prev) => prev.map((item) => (item.id === editProductId ? updated : item)))
      cancelEdit()
    } catch (err) {
      setProducts((prev) => prev.map((item) => (item.id === current.id ? current : item)))
      if (isAxiosError(err)) {
        setError((err.response?.data as { message?: string } | undefined)?.message ?? 'Could not update product.')
        return
      }
      setError('Could not update product.')
    }
  }

  const removeProduct = async (productId: string) => {
    const snapshot = [...products]
    setProducts((prev) => prev.filter((item) => item.id !== productId))
    setError('')

    try {
      await productService.remove(productId)
    } catch (err) {
      setProducts(snapshot)
      if (isAxiosError(err)) {
        setError((err.response?.data as { message?: string } | undefined)?.message ?? 'Could not delete product.')
        return
      }
      setError('Could not delete product.')
    }
  }

  return (
    <section>
      <h2>Inventory</h2>
      {!canManageProducts && <p className="info-text">You have read-only access to product inventory.</p>}

      {canManageProducts && (
        <form className="inline-form" onSubmit={onSubmit}>
          <input placeholder="Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          <input placeholder="Category" value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })} required />
          <input placeholder="Price" type="number" min="0" step="0.01" value={form.price} onChange={(e) => setForm({ ...form, price: Number(e.target.value) })} required />
          <input placeholder="Stock" type="number" min="0" value={form.stockQuantity} onChange={(e) => setForm({ ...form, stockQuantity: Number(e.target.value) })} required />
          <input placeholder="Low Stock Threshold" type="number" min="0" value={form.lowStockThreshold} onChange={(e) => setForm({ ...form, lowStockThreshold: Number(e.target.value) })} required />
          <input placeholder="Description" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
          <button type="submit">Add Product</button>
        </form>
      )}

      {canManageProducts && editProductId && (
        <form className="inline-form" onSubmit={saveEdit}>
          <input placeholder="Edit Name" value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} required />
          <input placeholder="Edit Category" value={editForm.category} onChange={(e) => setEditForm({ ...editForm, category: e.target.value })} required />
          <input placeholder="Edit Price" type="number" min="0" step="0.01" value={editForm.price} onChange={(e) => setEditForm({ ...editForm, price: Number(e.target.value) })} required />
          <input placeholder="Edit Stock" type="number" min="0" value={editForm.stockQuantity} onChange={(e) => setEditForm({ ...editForm, stockQuantity: Number(e.target.value) })} required />
          <input placeholder="Edit Low Stock Threshold" type="number" min="0" value={editForm.lowStockThreshold} onChange={(e) => setEditForm({ ...editForm, lowStockThreshold: Number(e.target.value) })} required />
          <input placeholder="Edit Description" value={editForm.description} onChange={(e) => setEditForm({ ...editForm, description: e.target.value })} />
          <button type="submit">Save</button>
          <button type="button" className="secondary-btn" onClick={cancelEdit}>Cancel</button>
        </form>
      )}

      {error && <p className="error-text">{error}</p>}

      <DataTable
        rows={products}
        columns={[
          { key: 'name', label: 'Name' },
          { key: 'category', label: 'Category' },
          { key: 'price', label: 'Price', render: (v) => formatINR(v) },
          { key: 'stockQuantity', label: 'Stock' },
          { key: 'lowStockThreshold', label: 'Threshold' },
          {
            key: 'id',
            label: 'Actions',
            render: (_, row) =>
              canManageProducts ? (
                <div className="actions-row">
                  <button type="button" className="secondary-btn" onClick={() => beginEdit(row)}>
                    Edit
                  </button>
                  <button type="button" className="danger-btn" onClick={() => removeProduct(row.id)}>
                    Delete
                  </button>
                </div>
              ) : (
                <span className="muted-text">No Actions</span>
              ),
          },
        ]}
      />
    </section>
  )
}
