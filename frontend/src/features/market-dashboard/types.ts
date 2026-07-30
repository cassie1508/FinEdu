// No "1D" — there's no free-tier intraday data source, so a true 1-day
// window would be a single daily bar, and a widened one just duplicates 1W.
export type ChartRange = '1W' | '1M' | '6M' | '1Y' | '5Y'

export const CHART_RANGES: ChartRange[] = ['1W', '1M', '6M', '1Y', '5Y']

export interface CandlePoint {
  t: number
  o: number
  h: number
  l: number
  c: number
  v: number
}

export interface PriceHistory {
  symbol: string
  range: ChartRange
  resolution: string
  candles: CandlePoint[]
  delayedMinutes: number
}

export interface NewsArticle {
  id: string
  headline: string
  summary: string
  source: string
  url: string
  publishedAt: string
  imageUrl: string
}

// Matches the backend's NewsDailySummary (models.go) — the cached AI-generated
// daily digest for a ticker, regenerated on a schedule rather than per request.
export interface TickerNewsSummary {
  ticker: string
  date: string
  dailySummary: string
  sentiment: 'bullish' | 'neutral' | 'bearish'
  sentimentScore: number
  potentialImpact: string
  sourceArticleIds: string[]
}
