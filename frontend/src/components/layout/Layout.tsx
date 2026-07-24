import type { ReactNode } from 'react'
import { useLocation } from 'react-router-dom'
import { Navbar } from './Navbar'

export function Layout({ children }: { children: ReactNode }) {
  const location = useLocation()
  const isFullHeightPage = location.pathname === '/learn'

  if (isFullHeightPage) {
    return children
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navbar />
      <main className="mx-auto max-w-6xl px-6 py-8">{children}</main>
    </div>
  )
}
