import { useEffect, useState } from 'react'
import { getPriceHistory } from '../api'
import type { ChartRange, PriceHistory } from '../types'

interface UsePriceHistoryResult {
  data: PriceHistory | null
  isLoading: boolean
  error: string | null
}

export function usePriceHistory(symbol: string, range: ChartRange): UsePriceHistoryResult {
  const [data, setData] = useState<PriceHistory | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    setIsLoading(true)
    setError(null)

    getPriceHistory(symbol, range)
      .then((result) => {
        if (!cancelled) setData(result)
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(err instanceof Error ? err.message : 'Failed to load price history')
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [symbol, range])

  return { data, isLoading, error }
}
