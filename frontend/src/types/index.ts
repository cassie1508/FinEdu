export interface Company {
  symbol: string
  name: string
  sector: string
  industry: string
  marketCap: number
  revenue: number
  eps: number
  peRatio: number
  dividendYield: number
  weekHigh52: number
  weekLow52: number
}

export interface PricePoint {
  timestamp: string
  price: number
}

export interface PortfolioHolding {
  id: string
  symbol: string
  shares: number
  averageCost: number
  currentValue: number
  unrealizedGainLoss: number
}

export interface NewsSummary {
  companySymbol: string
  headline: string
  summary: string
  sentiment: 'positive' | 'neutral' | 'negative'
  potentialImpact: string
  articleUrl: string
  publishedAt: string
}
