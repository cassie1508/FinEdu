import { NavLink } from 'react-router-dom'

const links = [
  { to: '/', label: 'Home' },
  { to: '/dashboard', label: 'Market Dashboard' },
  { to: '/news', label: 'News & Charts' },
  { to: '/learn', label: 'Learning Center' },
  { to: '/portfolio', label: 'Portfolio' },
]

export function Navbar() {
  return (
    <nav className="border-b border-[#CACDDC] bg-[#F1F0F3] px-4 sm:px-6 lg:px-8">
      <div className="mx-auto flex max-w-7xl flex-col gap-3 py-4 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
        <span className="shrink-0 text-lg font-bold tracking-[-0.03em] text-[#31302F]">
          Fin<span className="text-[#6E6C6F]">Edu</span>
        </span>
        <ul className="-mx-1 flex min-w-0 gap-1 overflow-x-auto px-1 pb-1 sm:mx-0 sm:gap-2 sm:px-0 sm:pb-0">
          {links.map((link) => (
            <li key={link.to} className="shrink-0">
              <NavLink
                to={link.to}
                end={link.to === '/'}
                className={({ isActive }) =>
                  `block rounded-full px-3 py-2 text-sm font-semibold transition-colors ${
                    isActive
                      ? 'bg-[#31302F] text-[#F1F0F3]'
                      : 'text-[#6E6C6F] hover:bg-[#E3DEDE] hover:text-[#31302F]'
                  }`
                }
              >
                {link.label}
              </NavLink>
            </li>
          ))}
        </ul>
      </div>
    </nav>
  )
}
