import type { PortfolioDetail, PortfolioHolding } from '../../../types'

export const formatCurrency = (value: number) =>
  new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
  }).format(value)

export const formatCompactCurrency = (value: number) =>
  new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    notation: 'compact',
    maximumFractionDigits: 1,
  }).format(value)

export const formatNumber = (value: number, maximumFractionDigits = 2) =>
  new Intl.NumberFormat('en-US', { maximumFractionDigits }).format(value)

export const formatPercent = (value: number) =>
  `${value > 0 ? '+' : ''}${value.toFixed(2)}%`

export function recalculatePortfolio(
  detail: PortfolioDetail,
  holdings: PortfolioHolding[],
): PortfolioDetail {
  const totalValue = holdings.reduce((sum, holding) => sum + holding.currentValue, 0)
  const totalCost = holdings.reduce(
    (sum, holding) => sum + holding.shares * holding.averageCost,
    0,
  )
  const totalGainLoss = totalValue - totalCost

  return {
    ...detail,
    holdings: holdings.map((holding) => ({
      ...holding,
      allocationPercent: totalValue > 0 ? (holding.currentValue / totalValue) * 100 : 0,
    })),
    totalValue,
    totalCost,
    totalGainLoss,
    totalGainLossPercent: totalCost > 0 ? (totalGainLoss / totalCost) * 100 : 0,
  }
}
