import { useEffect, useState } from 'react'
import { getTickerNewsSummary } from '../api'
import type { TickerNewsSummary } from '../types'

interface UseTickerNewsSummaryResult {
  summary: TickerNewsSummary | null
  isLoading: boolean
  error: string | null
}

export function useTickerNewsSummary(ticker: string): UseTickerNewsSummaryResult {
  const [summary, setSummary] = useState<TickerNewsSummary | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    setIsLoading(true)
    setSummary(null)
    setError(null)

    getTickerNewsSummary(ticker)
      .then((result) => {
        if (!cancelled) setSummary(result)
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(err instanceof Error ? err.message : 'AI summary unavailable')
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [ticker])

  return { summary, isLoading, error }
}
