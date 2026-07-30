import { api } from '../../lib/api'
import type { ChartRange, NewsArticle, PriceHistory, TickerNewsSummary } from './types'

export function getPriceHistory(symbol: string, range: ChartRange): Promise<PriceHistory> {
  return api.get<PriceHistory>(
    `/api/v1/companies/${encodeURIComponent(symbol)}/prices?range=${range}`,
  )
}

export function getTickerNews(ticker: string, limit = 20): Promise<NewsArticle[]> {
  return api.get<NewsArticle[]>(`/api/v1/news/${encodeURIComponent(ticker)}?limit=${limit}`)
}

export function getTickerNewsSummary(ticker: string): Promise<TickerNewsSummary> {
  return api.get<TickerNewsSummary>(`/api/v1/news/${encodeURIComponent(ticker)}/summary`)
}
