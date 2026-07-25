import { useState } from 'react'
import { api } from '../../../lib/api'
import type { Company } from '../../../types'

const WATCHLIST = [
  { symbol: 'AAPL', name: 'Apple Inc.' },
  { symbol: 'MSFT', name: 'Microsoft Corp.' },
  { symbol: 'GOOGL', name: 'Alphabet Inc.' },
  { symbol: 'AMZN', name: 'Amazon.com Inc.' },
  { symbol: 'TSLA', name: 'Tesla, Inc.' },
]

function formatBig(n: number): string {
  if (n >= 1e12) return (n / 1e12).toFixed(2) + 'T'
  if (n >= 1e9) return (n / 1e9).toFixed(2) + 'B'
  if (n >= 1e6) return (n / 1e6).toFixed(2) + 'M'
  return n.toLocaleString()
}

export function CompanyDashboardPage() {
  const [query, setQuery] = useState('')
  const [company, setCompany] = useState<Company | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function loadCompany(symbol: string) {
    if (!symbol.trim()) return
    setLoading(true)
    setError('')
    try {
      const data = await api.get<Company>(`/api/v1/companies/${symbol.toUpperCase()}`)
      setCompany(data)
    } catch {
      setError(`Không tìm thấy công ty "${symbol.toUpperCase()}"`)
      setCompany(null)
    } finally {
      setLoading(false)
    }
  }

  const metrics = company
    ? [
        { label: 'Market Cap', value: formatBig(company.marketCap) },
        { label: 'P/E Ratio', value: company.peRatio.toFixed(2) },
        { label: 'Revenue', value: formatBig(company.revenue) },
        { label: 'EPS', value: company.eps.toFixed(2) },
        { label: 'Dividend Yield', value: company.dividendYield.toFixed(2) + '%' },
        { label: '52W High', value: company.weekHigh52.toFixed(2) },
        { label: '52W Low', value: company.weekLow52.toFixed(2) },
        { label: 'Industry', value: company.industry },
      ]
    : []

  return (
    <div className="-mx-6 -my-8 min-h-screen bg-[#F8FAFC] px-6 py-10">
      <div className="mx-auto max-w-[1200px] space-y-8">

        {/* Heading */}
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900">
            Company Details
          </h1>
          <p className="mt-1.5 text-[15px] text-slate-500">
            Search a company to view its overview and key metrics.
          </p>
        </div>

        {/* Search bar */}
        <div className="flex gap-3">
          <div className="relative flex-1">
            <svg
              className="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400"
              fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-4.35-4.35M17 11a6 6 0 11-12 0 6 6 0 0112 0z" />
            </svg>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && loadCompany(query)}
              placeholder="Search by name or ticker (e.g. AAPL)…"
              className="w-full rounded-full border border-slate-200 bg-white py-3.5 pl-12 pr-4 text-[15px] text-slate-900 shadow-sm transition placeholder:text-slate-400 focus:border-[#5B7CFA] focus:outline-none focus:ring-4 focus:ring-[#5B7CFA]/10"
            />
          </div>
          <button
            onClick={() => loadCompany(query)}
            className="rounded-full bg-[#5B7CFA] px-8 text-[15px] font-semibold text-white shadow-sm transition hover:bg-[#4a6bf0] hover:shadow-md active:scale-[0.98]"
          >
            Search
          </button>
        </div>

        {loading && <p className="text-sm text-slate-500">Loading…</p>}
        {error && (
          <p className="rounded-xl bg-red-50 px-4 py-3 text-sm text-red-600">{error}</p>
        )}

        {/* Company header */}
        {company && (
          <>
            <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-sm">
              <div className="flex items-center gap-4">
                <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-[#5B7CFA] to-[#7c93fb] text-lg font-bold text-white shadow-sm">
                  {company.symbol.slice(0, 2)}
                </div>
                <div>
                  <div className="text-2xl font-bold tracking-tight text-slate-900">
                    {company.symbol}
                  </div>
                  <div className="text-sm text-slate-500">
                    {company.name} · {company.sector}
                  </div>
                </div>
              </div>
            </div>

            {/* Key metrics */}
            <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-sm">
              <h2 className="mb-5 text-base font-semibold text-slate-900">
                Key Metrics
              </h2>
              <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
                {metrics.map((m) => (
                  <div
                    key={m.label}
                    className="rounded-xl border border-slate-100 bg-[#F8FAFC] p-4 transition hover:-translate-y-0.5 hover:shadow-md"
                  >
                    <div className="text-xs font-medium uppercase tracking-wide text-slate-400">
                      {m.label}
                    </div>
                    <div className="mt-2 text-xl font-bold text-slate-900">
                      {m.value}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </>
        )}

        {/* Placeholder chừa chỗ cho Nhi */}
        <div className="flex min-h-[220px] items-center justify-center rounded-2xl border border-dashed border-slate-200 bg-white/60 text-center text-sm text-slate-400">
          Interactive Stock Chart · News · AI Summary (owner: Nhi)
        </div>

        {/* Watchlist */}
        <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-sm">
          <h2 className="mb-5 text-base font-semibold text-slate-900">Watchlist</h2>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-5">
            {WATCHLIST.map((w) => (
              <button
                key={w.symbol}
                onClick={() => {
                  setQuery(w.symbol)
                  loadCompany(w.symbol)
                }}
                className="rounded-xl border border-slate-100 bg-[#F8FAFC] p-4 text-left transition hover:-translate-y-0.5 hover:border-[#5B7CFA]/40 hover:shadow-md"
              >
                <div className="text-base font-bold text-slate-900">{w.symbol}</div>
                <div className="mt-0.5 text-xs text-slate-500">{w.name}</div>
              </button>
            ))}
          </div>
        </div>

      </div>
    </div>
  )
}