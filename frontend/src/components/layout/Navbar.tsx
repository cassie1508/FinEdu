import { useEffect, useRef, useState } from 'react'
import { Link, NavLink } from 'react-router-dom'
import { useSession } from '../../lib/useSession'
import { supabase } from '../../lib/supabaseClient'

const links = [
  { to: '/', label: 'Home' },
  { to: '/dashboard', label: 'Market Dashboard' },
  { to: '/news', label: 'News' },
  { to: '/learn', label: 'Learning Center' },
  { to: '/portfolio', label: 'Portfolio' },
]

export function Navbar({ transparent = false }: { transparent?: boolean }) {
  const session = useSession()
  const [menuOpen, setMenuOpen] = useState(false)
  const menuRef = useRef<HTMLLIElement>(null)

  useEffect(() => {
    if (!menuOpen) return
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setMenuOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [menuOpen])

  const handleLogout = async () => {
    setMenuOpen(false)
    await supabase.auth.signOut()
  }

  return (
    <nav
      className={
        transparent
          ? 'px-6 py-4'
          : 'border-b border-gray-200 bg-white px-6 py-4 dark:border-gray-800 dark:bg-gray-950'
      }
    >
      <div className="mx-auto flex max-w-6xl items-center justify-between">
        <Link to="/" className="flex items-center gap-2">
          <FinEduLogo />
          <span className="text-lg font-semibold text-gray-900 dark:text-gray-100">FinEdu</span>
        </Link>
        <ul className="flex items-center gap-6">
          {links.map((link) => (
            <li key={link.to}>
              <NavLink
                to={link.to}
                end={link.to === '/'}
                className={({ isActive }) =>
                  `text-sm font-medium ${
                    isActive
                      ? 'text-indigo-600 dark:text-indigo-400'
                      : 'text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100'
                  }`
                }
              >
                {link.label}
              </NavLink>
            </li>
          ))}
          {session && (
            <li className="relative" ref={menuRef}>
              <button
                type="button"
                aria-label="Profile"
                aria-expanded={menuOpen}
                onClick={() => setMenuOpen((open) => !open)}
                className="flex h-9 w-9 items-center justify-center rounded-full bg-gray-100 text-gray-600 transition hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
              >
                <ProfileIcon />
              </button>
              {menuOpen && (
                <div className="absolute right-0 top-11 z-10 w-40 rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-gray-800 dark:bg-gray-900">
                  <button
                    type="button"
                    onClick={handleLogout}
                    className="block w-full px-4 py-2 text-left text-sm font-medium text-gray-600 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:text-gray-100"
                  >
                    Log out
                  </button>
                </div>
              )}
            </li>
          )}
        </ul>
        {!session && (
          <div className="flex items-center gap-3">
            <Link
              to="/login"
              className="text-sm font-medium text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100"
            >
              Sign in
            </Link>
            <Link
              to="/register"
              className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-700"
            >
              Register
            </Link>
          </div>
        )}
      </div>
    </nav>
  )
}

function FinEduLogo() {
  return (
    <svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="28" height="28" rx="8" fill="#1F1F1F" />
      <path
        d="M8 18.5L11.8 14.2L14.6 16.8L20 10.5"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M16.2 10.5H20V14.3"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function ProfileIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="9" cy="6" r="3" stroke="currentColor" strokeWidth="1.5" />
      <path
        d="M3.5 15.25C4.36 12.64 6.47 11 9 11C11.53 11 13.64 12.64 14.5 15.25"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  )
}
