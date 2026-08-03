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

export interface NewsSummary {
  companySymbol: string
  headline: string
  summary: string
  sentiment: 'positive' | 'neutral' | 'negative'
  potentialImpact: string
  articleUrl: string
  publishedAt: string
}

// --- Portfolio Types ---

export interface Portfolio {
  id: string
  userId: string
  name: string
  description: string
  createdAt: string
  updatedAt: string
}

export interface PortfolioListItem {
  id: string
  name: string
  description: string
  holdingsCount: number
  createdAt: string
}

export interface PortfolioHolding {
  id: string
  portfolioId: string
  symbol: string
  shares: number
  averageCost: number
  currentPrice: number
  currentValue: number
  unrealizedGainLoss: number
  unrealizedGainLossPercent: number
  allocationPercent: number
}

export interface PortfolioDetail {
  id: string
  name: string
  description: string
  holdings: PortfolioHolding[]
  totalValue: number
  totalCost: number
  totalGainLoss: number
  totalGainLossPercent: number
}

export interface PortfolioSummary {
  totalValue: number
  totalCost: number
  totalGainLoss: number
  totalGainLossPercent: number
  allocations: AllocationEntry[]
}

export interface AllocationEntry {
  symbol: string
  percent: number
}

export interface PortfolioRisk {
  healthScore: number | null
  riskLevel: string | null
  sectorConcentration: SectorEntry[]
  diversificationScore: number | null
  recommendations: string[]
  message?: string
}

export interface SectorEntry {
  sector: string
  percent: number
}

// --- Request Types ---

export interface CreatePortfolioRequest {
  name: string
  description?: string
}

export interface AddHoldingRequest {
  symbol: string
  shares: number
  averageCost: number
}

export interface UpdateHoldingRequest {
  shares: number
  averageCost: number
}

// --- Response Wrappers ---

export interface PortfolioListResponse {
  portfolios: PortfolioListItem[]
}
