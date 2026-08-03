import type { PortfolioDetail } from '../../../types'
import { formatCompactCurrency, formatCurrency, formatPercent } from '../lib/portfolioFormatters'

export function SummaryCards({ portfolio }: { portfolio: PortfolioDetail }) {
  const performanceDirection = portfolio.totalGainLoss >= 0 ? 'Up' : 'Down'

  const metrics = [
    {
      label: 'Current value',
      value: formatCurrency(portfolio.totalValue),
      note: `${portfolio.holdings.length} holding${portfolio.holdings.length === 1 ? '' : 's'}`,
    },
    {
      label: 'Total return',
      value: formatCurrency(portfolio.totalGainLoss),
      note: `${performanceDirection} ${formatPercent(Math.abs(portfolio.totalGainLossPercent))}`,
      directional: true,
    },
    {
      label: 'Amount invested',
      value: formatCompactCurrency(portfolio.totalCost),
      note: 'Based on average cost',
    },
  ]

  return (
    <section aria-label="Portfolio summary" className="grid gap-3 sm:grid-cols-3">
      {metrics.map((metric, index) => (
        <article
          key={metric.label}
          className={`rounded-2xl border p-5 ${
            index === 0
              ? 'border-[#31302F] bg-[#31302F] text-[#F1F0F3]'
              : 'border-[#CACDDC] bg-[#F1F0F3] text-[#31302F]'
          }`}
        >
          <p className={`text-xs font-bold uppercase tracking-[0.14em] ${index === 0 ? 'text-[#CACDDC]' : 'text-[#6E6C6F]'}`}>
            {metric.label}
          </p>
          <p className="mt-3 break-words text-2xl font-bold tracking-[-0.04em] sm:text-[1.7rem] xl:text-3xl">
            {metric.value}
          </p>
          <p className={`mt-2 text-sm font-semibold ${
            index === 0 ? 'text-[#E3DEDE]' : metric.directional ? 'text-[#31302F]' : 'text-[#6E6C6F]'
          }`}>
            {metric.directional && (portfolio.totalGainLoss >= 0 ? '↗ ' : '↘ ')}
            {metric.note}
          </p>
        </article>
      ))}
    </section>
  )
}
