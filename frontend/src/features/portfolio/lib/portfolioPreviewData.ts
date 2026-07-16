import type {
  PortfolioDetail,
  PortfolioListItem,
  PortfolioRisk,
} from '../../../types'

export const previewPortfolioList: PortfolioListItem[] = [
  {
    id: 'preview-core-growth',
    name: 'Core Growth',
    description: 'A long-term portfolio focused on durable businesses.',
    holdingsCount: 5,
    createdAt: '2026-07-02T09:00:00.000Z',
  },
  {
    id: 'preview-retirement',
    name: 'Retirement',
    description: 'A steady, diversified portfolio for the long horizon.',
    holdingsCount: 0,
    createdAt: '2026-06-18T09:00:00.000Z',
  },
]

export const previewPortfolioDetails: Record<string, PortfolioDetail> = {
  'preview-core-growth': {
    id: 'preview-core-growth',
    name: 'Core Growth',
    description: 'A long-term portfolio focused on durable businesses.',
    totalValue: 48736.48,
    totalCost: 43180,
    totalGainLoss: 5556.48,
    totalGainLossPercent: 12.87,
    holdings: [
      {
        id: 'holding-aapl', portfolioId: 'preview-core-growth', symbol: 'AAPL',
        shares: 60, averageCost: 184.2, currentPrice: 213.32, currentValue: 12799.2,
        unrealizedGainLoss: 1747.2, unrealizedGainLossPercent: 15.81, allocationPercent: 26.26,
      },
      {
        id: 'holding-msft', portfolioId: 'preview-core-growth', symbol: 'MSFT',
        shares: 24, averageCost: 408.5, currentPrice: 511.16, currentValue: 12267.84,
        unrealizedGainLoss: 2463.84, unrealizedGainLossPercent: 25.13, allocationPercent: 25.17,
      },
      {
        id: 'holding-vti', portfolioId: 'preview-core-growth', symbol: 'VTI',
        shares: 40, averageCost: 281.25, currentPrice: 301.44, currentValue: 12057.6,
        unrealizedGainLoss: 807.6, unrealizedGainLossPercent: 7.18, allocationPercent: 24.74,
      },
      {
        id: 'holding-jpm', portfolioId: 'preview-core-growth', symbol: 'JPM',
        shares: 25, averageCost: 227.4, currentPrice: 290.11, currentValue: 7252.75,
        unrealizedGainLoss: 1567.75, unrealizedGainLossPercent: 27.58, allocationPercent: 14.88,
      },
      {
        id: 'holding-nke', portfolioId: 'preview-core-growth', symbol: 'NKE',
        shares: 50, averageCost: 108.78, currentPrice: 87.18, currentValue: 4359,
        unrealizedGainLoss: -1080, unrealizedGainLossPercent: -19.86, allocationPercent: 8.95,
      },
    ],
  },
  'preview-retirement': {
    id: 'preview-retirement',
    name: 'Retirement',
    description: 'A steady, diversified portfolio for the long horizon.',
    holdings: [],
    totalValue: 0,
    totalCost: 0,
    totalGainLoss: 0,
    totalGainLossPercent: 0,
  },
}

export const previewPortfolioRisk: PortfolioRisk = {
  healthScore: 78,
  riskLevel: 'Moderate',
  diversificationScore: 72,
  sectorConcentration: [
    { sector: 'Technology', percent: 51 },
    { sector: 'Index fund', percent: 25 },
    { sector: 'Financials', percent: 15 },
    { sector: 'Consumer', percent: 9 },
  ],
  recommendations: [
    'Technology makes up more than half of this portfolio. Consider adding exposure to less-correlated sectors.',
    'Your position sizes are generally balanced, with no single holding above 30%.',
  ],
  message: 'Illustrative preview — live scores come from the portfolio risk service.',
}
