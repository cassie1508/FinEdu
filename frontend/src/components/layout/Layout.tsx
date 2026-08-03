import type { ReactNode } from 'react'
import { useLocation } from 'react-router-dom'
import { Navbar } from './Navbar'

export function Layout({ children }: { children: ReactNode }) {
  const { pathname } = useLocation()
  const isHome = pathname === '/'

  if (isHome) {
    return (
      <div className="relative min-h-screen overflow-hidden bg-[linear-gradient(135deg,#9EC0FF_0%,#F8FAFF_45%,#FFF6E2_72%,#FFDF94_100%)]">
        <Navbar transparent />
        <main>{children}</main>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-[#F1F0F3] text-[#31302F]">
      <Navbar />
      <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
        {children}
      </main>
    </div>
  )
}
