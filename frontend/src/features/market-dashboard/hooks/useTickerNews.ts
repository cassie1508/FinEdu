import { useEffect, useState } from 'react'
import { getTickerNews } from '../api'
import type { NewsArticle } from '../types'

interface UseTickerNewsResult {
  articles: NewsArticle[]
  isLoading: boolean
  error: string | null
}

export function useTickerNews(ticker: string): UseTickerNewsResult {
  const [articles, setArticles] = useState<NewsArticle[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    setIsLoading(true)
    setError(null)

    getTickerNews(ticker)
      .then((result) => {
        if (!cancelled) setArticles(result)
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(err instanceof Error ? err.message : 'Failed to load news')
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [ticker])

  return { articles, isLoading, error }
}
