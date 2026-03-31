import type { ReactNode } from 'react'

interface Column<T> {
  key: keyof T
  label: string
  render?: (value: T[keyof T], row: T) => ReactNode
}

interface DataTableProps<T extends { id: string }> {
  columns: Array<Column<T>>
  rows: T[]
}

export function DataTable<T extends { id: string }>({ columns, rows }: DataTableProps<T>) {
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            {columns.map((column) => (
              <th key={String(column.key)}>{column.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.id}>
              {columns.map((column) => (
                <td key={String(column.key)}>
                  {column.render ? column.render(row[column.key], row) : String(row[column.key])}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
      {!rows.length && <p className="empty">No data available.</p>}
    </div>
  )
}
