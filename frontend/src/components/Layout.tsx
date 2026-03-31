import { Link, Outlet, useLocation } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import type { Role } from '../types'

const links = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/products', label: 'Products' },
  { to: '/orders', label: 'Orders' },
]

const navigationByRole: Record<Role, string[]> = {
  admin: ['/dashboard', '/products', '/orders'],
  staff: ['/dashboard', '/products', '/orders'],
  customer: ['/dashboard', '/products', '/orders'],
}

export function Layout() {
  const location = useLocation()
  const { user, logout } = useAuth()
  const role = user?.role ?? 'customer'
  const visibleLinks = links.filter((item) => navigationByRole[role].includes(item.to))

  return (
    <div className="shell">
      <aside className="sidebar">
        <h1>Inventory OS</h1>
        <p>Modern inventory and order workflow hub.</p>
        <p className="role-note">Role Access: {role.toUpperCase()}</p>
        <nav>
          {visibleLinks.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              className={location.pathname.startsWith(item.to) ? 'active' : ''}
            >
              {item.label}
            </Link>
          ))}
        </nav>
      </aside>

      <main className="content">
        <header className="topbar">
          <div>
            <strong>{user?.name}</strong>
            <span>{user?.role}</span>
          </div>
          <button onClick={logout}>Logout</button>
        </header>
        <Outlet />
      </main>
    </div>
  )
}
