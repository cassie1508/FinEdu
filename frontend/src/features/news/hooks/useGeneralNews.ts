import { useEffect, useState } from 'react'
import { getGeneralNews } from '../api'
import type { NewsArticle, NewsCategory } from '../types'

interface UseGeneralNewsResult {
  articles: NewsArticle[]
  isLoading: boolean
  error: string | null
}

export function useGeneralNews(category: NewsCategory): UseGeneralNewsResult {
  const [articles, setArticles] = useState<NewsArticle[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    setIsLoading(true)
    setError(null)

    getGeneralNews(category)
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
  }, [category])

  return { articles, isLoading, error }
}
