import { useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { isValidEmail } from '../lib/validation'
import heroIllustration from '../assets/hero-illustration.png'

const TRENDING = [
  { symbol: 'AAPL', change: '+2.31%', direction: 'up' as const },
  { symbol: 'NVDA', change: '+4.82%', direction: 'up' as const },
  { symbol: 'TSLA', change: '-1.20%', direction: 'down' as const },
  { symbol: 'MSFT', change: '+1.01%', direction: 'up' as const },
  { symbol: 'BTC', change: '+3.45%', direction: 'up' as const },
]

const FOOTER_COLUMNS: { title: string; links: { label: string; to?: string }[] }[] = [
  {
    title: 'Use cases',
    links: [
      { label: 'Portfolio Tracking', to: '/portfolio' },
      { label: 'Company Research', to: '/dashboard' },
      { label: 'Market Analysis', to: '/news' },
      { label: 'Financial Literacy', to: '/learn' },
      { label: 'Investment Simulator' },
      { label: 'Risk Assessment' },
    ],
  },
  {
    title: 'Explore',
    links: [
      { label: 'Market Dashboard', to: '/dashboard' },
      { label: 'News & Charts', to: '/news' },
      { label: 'Learning Center', to: '/learn' },
      { label: 'Portfolio', to: '/portfolio' },
      { label: 'Community' },
      { label: 'Blog' },
    ],
  },
  {
    title: 'Resources',
    links: [
      { label: 'Help Center' },
      { label: 'Glossary' },
      { label: 'FAQs' },
      { label: 'Support' },
      { label: 'Guides' },
      { label: 'Contact Us' },
    ],
  },
]

export function HomePage() {
  return (
    <div>
      <HeroSection />
      <NewsletterSection />
      <Footer />
    </div>
  )
}

function HeroSection() {
  return (
    <section className="relative overflow-hidden">
      <div
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(62% 62% at 74% 40%, rgba(158,192,255,0.42) 0%, rgba(199,207,255,0.28) 26%, rgba(232,228,255,0.18) 46%, rgba(255,246,226,0.14) 68%, rgba(255,255,255,0) 88%)',
        }}
      />

      <div className="relative z-10 mx-auto max-w-6xl px-6 pt-16 pb-24">
        <div className="grid items-center gap-12 lg:grid-cols-2">
          <div className="relative z-20">
            <h1 className="text-4xl leading-[1.1] font-bold tracking-tight text-[#1F1F1F] sm:text-5xl">
              Learn. Analyze.
              <br />
              Invest <span className="text-[#7c97d4]">Smarter</span>.
            </h1>
            <p className="mt-5 max-w-md text-base leading-relaxed text-[#6D7280]">
              AI-powered market insights, interactive charts, and bite-sized lessons to grow your
              wealth.
            </p>

            <div className="mt-8 flex flex-wrap items-center gap-4">
              <Link
                to="/dashboard"
                className="flex items-center gap-2 rounded-xl bg-[#1F1F1F] px-6 py-3 text-sm font-semibold text-white transition hover:bg-black"
              >
                Explore Market
                <ArrowRightIcon />
              </Link>
              <Link
                to="/learn"
                className="rounded-xl border border-[#EAEAEA] bg-white/70 px-6 py-3 text-sm font-semibold text-[#1F1F1F] transition hover:bg-white"
              >
                Start Learning
              </Link>
            </div>
          </div>

          <MarketIllustration />
        </div>

        <TrendingStrip />
      </div>
    </section>
  )
}

function TrendingStrip() {
  return (
    <div className="mt-16">
      <span className="text-sm font-medium text-[#6D7280]">Trending</span>
      <div className="mt-3 flex flex-wrap gap-3">
        {TRENDING.map((item) => (
          <div
            key={item.symbol}
            className="flex items-center gap-2 rounded-full border border-[#EAEAEA] bg-white/70 px-4 py-2 text-sm shadow-[0_2px_10px_-4px_rgba(31,31,31,0.08)]"
          >
            <span className="font-semibold text-[#1F1F1F]">{item.symbol}</span>
            <span className={item.direction === 'up' ? 'text-[#16A34A]' : 'text-[#DC2626]'}>
              {item.change}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}

function NewsletterSection() {
  const [email, setEmail] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [subscribed, setSubscribed] = useState(false)

  function handleSubmit(event: FormEvent) {
    event.preventDefault()
    if (!isValidEmail(email)) {
      setError('Enter a valid email address')
      setSubscribed(false)
      return
    }
    setError(null)
    setSubscribed(true)
    setEmail('')
  }

  return (
    <section className="relative z-10 border-t border-[#EAEAEA]/70 px-6 py-16 text-center">
      <h2 className="text-2xl font-bold text-[#1F1F1F]">Follow the latest trends</h2>
      <p className="mt-2 text-base text-[#6D7280]">With our daily newsletter</p>

      <form
        onSubmit={handleSubmit}
        className="mx-auto mt-6 flex w-full max-w-md flex-col gap-3 sm:flex-row"
      >
        <input
          type="email"
          required
          placeholder="you@example.com"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value)
            setSubscribed(false)
          }}
          aria-label="Email address"
          className="w-full rounded-xl border border-[#EAEAEA] bg-white px-4 py-3 text-sm text-[#1F1F1F] placeholder:text-[#9CA0A8] focus:border-[#6C7CFF] focus:ring-2 focus:ring-[#DDE7FF] focus:outline-none"
        />
        <button
          type="submit"
          className="shrink-0 rounded-xl bg-[#1F1F1F] px-6 py-3 text-sm font-semibold text-white transition hover:bg-black"
        >
          Submit
        </button>
      </form>

      {error && <p className="mt-3 text-sm text-red-600">{error}</p>}
      {subscribed && <p className="mt-3 text-sm text-[#16A34A]">You're subscribed — welcome aboard!</p>}
    </section>
  )
}

function Footer() {
  return (
    <footer className="relative z-10 border-t border-[#EAEAEA]/70 px-6 py-12">
      <div className="mx-auto max-w-6xl">
        <div className="grid grid-cols-2 gap-10 sm:grid-cols-4">
          <div className="col-span-2 sm:col-span-1">
            <div className="flex items-center gap-2">
              <FinEduLogo />
              <span className="text-lg font-semibold text-[#1F1F1F]">FinEdu</span>
            </div>
            <p className="mt-3 max-w-55 text-sm leading-relaxed text-[#6D7280]">
              AI-powered market insights and bite-sized lessons to help you invest smarter.
            </p>
          </div>

          {FOOTER_COLUMNS.map((column) => (
            <div key={column.title}>
              <h3 className="text-sm font-semibold text-[#1F1F1F]">{column.title}</h3>
              <ul className="mt-4 space-y-3">
                {column.links.map((link) => (
                  <li key={link.label}>
                    {link.to ? (
                      <Link
                        to={link.to}
                        className="text-sm text-[#6D7280] transition hover:text-[#1F1F1F]"
                      >
                        {link.label}
                      </Link>
                    ) : (
                      <a href="#" className="text-sm text-[#6D7280] transition hover:text-[#1F1F1F]">
                        {link.label}
                      </a>
                    )}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="mt-12 flex flex-col-reverse items-center justify-between gap-4 border-t border-[#EAEAEA]/70 pt-6 sm:flex-row">
          <span className="text-[13px] text-[#6D7280]">
            © {new Date().getFullYear()} FinEdu. All rights reserved.
          </span>
          <div className="flex items-center gap-4 text-[#6D7280]">
            <a href="#" aria-label="X" className="transition hover:text-[#1F1F1F]">
              <XIcon />
            </a>
            <a href="#" aria-label="Instagram" className="transition hover:text-[#1F1F1F]">
              <InstagramIcon />
            </a>
            <a href="#" aria-label="YouTube" className="transition hover:text-[#1F1F1F]">
              <YouTubeIcon />
            </a>
            <a href="#" aria-label="LinkedIn" className="transition hover:text-[#1F1F1F]">
              <LinkedInIcon />
            </a>
          </div>
        </div>
      </div>
    </footer>
  )
}

function MarketIllustration() {
  return (
    <div className="pointer-events-none relative z-0 mx-auto h-80 w-full max-w-md select-none overflow-visible">
      <div
        className="absolute inset-0 rounded-full blur-3xl"
        style={{
          background: 'radial-gradient(closest-side, rgba(140,170,255,0.42), rgba(140,170,255,0) 72%)',
        }}
      />
      <div className="absolute inset-0 flex items-center justify-center overflow-visible">
        <img
          src={heroIllustration}
          alt=""
          className="h-full w-full scale-150 object-contain"
          style={{
            maskImage: 'radial-gradient(ellipse closest-side at 50% 46%, black 40%, transparent 92%)',
            WebkitMaskImage: 'radial-gradient(ellipse closest-side at 50% 46%, black 40%, transparent 92%)',
          }}
        />
      </div>
    </div>
  )
}

function ArrowRightIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M3.5 8H12.5M12.5 8L8.5 4M12.5 8L8.5 12"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
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

function XIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M2 2L16 16M16 2L2 16"
        stroke="currentColor"
        strokeWidth="1.6"
        strokeLinecap="round"
      />
    </svg>
  )
}

function InstagramIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="2" y="2" width="14" height="14" rx="4" stroke="currentColor" strokeWidth="1.5" />
      <circle cx="9" cy="9" r="3.5" stroke="currentColor" strokeWidth="1.5" />
      <circle cx="13" cy="5" r="1" fill="currentColor" />
    </svg>
  )
}

function YouTubeIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="1.5" y="4" width="15" height="10" rx="3" stroke="currentColor" strokeWidth="1.5" />
      <path d="M7.5 6.8L11.5 9L7.5 11.2V6.8Z" fill="currentColor" />
    </svg>
  )
}

function LinkedInIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="2" y="2" width="14" height="14" rx="3" stroke="currentColor" strokeWidth="1.5" />
      <circle cx="5.5" cy="5.8" r="1" fill="currentColor" />
      <path d="M5.5 8.2V13" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <path
        d="M8.5 13V10.2C8.5 9.1 9.2 8.2 10.3 8.2C11.4 8.2 12 9.1 12 10.2V13"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  )
}
