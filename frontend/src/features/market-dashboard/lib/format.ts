import type { ChartRange } from '../types'

export function formatPrice(value: number): string {
  return value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

export function formatChange(change: number, changePercent: number): string {
  const sign = change >= 0 ? '+' : ''
  return `${sign}${formatPrice(change)} (${sign}${changePercent.toFixed(2)}%)`
}

export function formatVolume(value: number): string {
  if (value >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(2)}B`
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(2)}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(2)}K`
  return value.toFixed(0)
}

// Daily/weekly resolution only (no intraday data source available on the
// free tier), so ticks are always a date, never a time of day.
//
// The backend parses Alpha Vantage's date-only bars (e.g. "2026-07-23") as
// UTC midnight, since AV gives no timezone info. Formatting must pin
// timeZone: 'UTC' here to match — otherwise a viewer west of UTC (i.e.
// anywhere in the US) sees every date rolled back by a day, since midnight
// UTC is still the previous evening in their local time.
export function formatAxisTick(unixSeconds: number, range: ChartRange): string {
  const date = new Date(unixSeconds * 1000)
  if (range === '5Y' || range === '1Y') {
    return date.toLocaleDateString('en-US', { month: 'short', year: '2-digit', timeZone: 'UTC' })
  }
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', timeZone: 'UTC' })
}

// News timestamps are ISO strings with real times (unlike the date-only
// chart bars above), so this needs no UTC pinning — Date parses the offset
// that's already in the string.
export function formatRelativeTime(publishedAt: string): string {
  const diffMs = Date.now() - new Date(publishedAt).getTime()
  const minutes = Math.round(diffMs / 60_000)

  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`

  const hours = Math.round(minutes / 60)
  if (hours < 24) return `${hours}h ago`

  const days = Math.round(hours / 24)
  if (days < 7) return `${days}d ago`

  return new Date(publishedAt).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

export function formatTooltipDate(unixSeconds: number): string {
  return new Date(unixSeconds * 1000).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  })
}
