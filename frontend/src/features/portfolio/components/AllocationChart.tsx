import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'
import type { PortfolioHolding } from '../../../types'
import { formatCurrency } from '../lib/portfolioFormatters'

const chartColors = ['#31302F', '#6E6C6F', '#A5A4A8', '#CACDDC', '#E3DEDE']

interface TooltipPayload {
  payload?: PortfolioHolding
}

function AllocationTooltip({ active, payload }: { active?: boolean; payload?: TooltipPayload[] }) {
  const holding = payload?.[0]?.payload
  if (!active || !holding) return null

  return (
    <div className="rounded-xl border border-[#CACDDC] bg-[#F1F0F3] px-3 py-2 text-xs text-[#31302F]">
      <p className="font-bold">{holding.symbol}</p>
      <p className="mt-1 text-[#6E6C6F]">
        {holding.allocationPercent.toFixed(1)}% · {formatCurrency(holding.currentValue)}
      </p>
    </div>
  )
}

export function AllocationChart({ holdings }: { holdings: PortfolioHolding[] }) {
  return (
    <section className="rounded-2xl border border-[#CACDDC] bg-[#F1F0F3] p-5 sm:p-6">
      <div>
        <p className="text-xs font-bold uppercase tracking-[0.14em] text-[#6E6C6F]">Allocation</p>
        <h2 className="mt-1 text-xl font-bold tracking-[-0.03em] text-[#31302F]">Where your money sits</h2>
      </div>

      <div className="mt-4 grid items-center gap-5 sm:grid-cols-[minmax(0,1fr)_minmax(150px,0.8fr)]">
        <div className="h-56 min-w-0" aria-label="Portfolio allocation donut chart">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={holdings}
                dataKey="currentValue"
                nameKey="symbol"
                innerRadius="61%"
                outerRadius="88%"
                paddingAngle={2}
                stroke="#F1F0F3"
                strokeWidth={3}
              >
                {holdings.map((holding, index) => (
                  <Cell key={holding.id} fill={chartColors[index % chartColors.length]} />
                ))}
              </Pie>
              <Tooltip cursor={false} content={<AllocationTooltip />} />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <ul className="grid gap-3" aria-label="Allocation legend">
          {holdings.map((holding, index) => (
            <li key={holding.id} className="flex items-center justify-between gap-3 text-sm">
              <span className="flex min-w-0 items-center gap-2 font-bold text-[#31302F]">
                <span
                  className="h-2.5 w-2.5 shrink-0 rounded-full"
                  style={{ backgroundColor: chartColors[index % chartColors.length] }}
                />
                <span className="truncate">{holding.symbol}</span>
              </span>
              <span className="font-semibold tabular-nums text-[#6E6C6F]">
                {holding.allocationPercent.toFixed(1)}%
              </span>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
