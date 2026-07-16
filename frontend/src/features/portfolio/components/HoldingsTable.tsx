import type { PortfolioHolding } from '../../../types'
import { formatCurrency, formatNumber, formatPercent } from '../lib/portfolioFormatters'

interface HoldingsTableProps {
  holdings: PortfolioHolding[]
  busy: boolean
  onAdd: () => void
  onEdit: (holding: PortfolioHolding) => void
  onDelete: (holding: PortfolioHolding) => void
}

function ReturnValue({ holding }: { holding: PortfolioHolding }) {
  const direction = holding.unrealizedGainLoss >= 0 ? 'Gain' : 'Loss'

  return (
    <span className="font-bold text-[#31302F]">
      <span className="block">{formatCurrency(holding.unrealizedGainLoss)}</span>
      <span className="mt-0.5 block text-xs font-semibold text-[#6E6C6F]">
        {direction} · {formatPercent(holding.unrealizedGainLossPercent)}
      </span>
    </span>
  )
}

export function HoldingsTable({ holdings, busy, onAdd, onEdit, onDelete }: HoldingsTableProps) {
  return (
    <section className="overflow-hidden rounded-2xl border border-[#CACDDC] bg-[#F1F0F3]">
      <div className="flex flex-col gap-3 border-b border-[#CACDDC] p-5 sm:flex-row sm:items-center sm:justify-between sm:p-6">
        <div>
          <p className="text-xs font-bold uppercase tracking-[0.14em] text-[#6E6C6F]">Holdings</p>
          <h2 className="mt-1 text-xl font-bold tracking-[-0.03em] text-[#31302F]">What you own</h2>
        </div>
        <button
          type="button"
          onClick={onAdd}
          disabled={busy}
          className="min-h-11 rounded-xl border border-[#31302F] px-4 text-sm font-bold text-[#31302F] transition-colors hover:bg-[#31302F] hover:text-[#F1F0F3] disabled:cursor-not-allowed disabled:border-[#A5A4A8] disabled:text-[#A5A4A8]"
        >
          + Add holding
        </button>
      </div>

      {holdings.length === 0 ? (
        <div className="px-5 py-14 text-center sm:px-6">
          <p className="text-lg font-bold text-[#31302F]">Build your first allocation</p>
          <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-[#6E6C6F]">
            Add a stock or fund with the number of shares and your average purchase price.
          </p>
          <button
            type="button"
            onClick={onAdd}
            className="mt-5 min-h-11 rounded-xl bg-[#31302F] px-5 text-sm font-bold text-[#F1F0F3]"
          >
            Add first holding
          </button>
        </div>
      ) : (
        <>
          <div className="hidden overflow-x-auto md:block">
            <table className="w-full min-w-[760px] border-collapse text-left">
              <thead className="bg-[#E3DEDE] text-xs font-bold uppercase tracking-[0.08em] text-[#6E6C6F]">
                <tr>
                  <th className="px-6 py-3.5">Asset</th>
                  <th className="px-4 py-3.5">Shares</th>
                  <th className="px-4 py-3.5">Avg. cost</th>
                  <th className="px-4 py-3.5">Price</th>
                  <th className="px-4 py-3.5">Market value</th>
                  <th className="px-4 py-3.5">Return</th>
                  <th className="px-6 py-3.5 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#CACDDC]">
                {holdings.map((holding) => (
                  <tr key={holding.id} className="transition-colors hover:bg-[#E3DEDE]">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <span className="grid h-10 w-10 place-items-center rounded-xl bg-[#31302F] text-xs font-bold text-[#F1F0F3]">
                          {holding.symbol.slice(0, 2)}
                        </span>
                        <span>
                          <strong className="block text-sm text-[#31302F]">{holding.symbol}</strong>
                          <span className="text-xs text-[#6E6C6F]">{holding.allocationPercent.toFixed(1)}% of portfolio</span>
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-4 text-sm font-semibold tabular-nums text-[#31302F]">{formatNumber(holding.shares, 4)}</td>
                    <td className="px-4 py-4 text-sm font-semibold tabular-nums text-[#6E6C6F]">{formatCurrency(holding.averageCost)}</td>
                    <td className="px-4 py-4 text-sm font-semibold tabular-nums text-[#6E6C6F]">{formatCurrency(holding.currentPrice)}</td>
                    <td className="px-4 py-4 text-sm font-bold tabular-nums text-[#31302F]">{formatCurrency(holding.currentValue)}</td>
                    <td className="px-4 py-4 text-sm tabular-nums"><ReturnValue holding={holding} /></td>
                    <td className="px-6 py-4">
                      <div className="flex justify-end gap-2">
                        <button type="button" onClick={() => onEdit(holding)} className="min-h-10 rounded-lg px-3 text-xs font-bold text-[#6E6C6F] hover:bg-[#CACDDC] hover:text-[#31302F]">Edit</button>
                        <button type="button" onClick={() => onDelete(holding)} className="min-h-10 rounded-lg px-3 text-xs font-bold text-[#6E6C6F] hover:bg-[#CACDDC] hover:text-[#31302F]">Remove</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <ul className="divide-y divide-[#CACDDC] md:hidden">
            {holdings.map((holding) => (
              <li key={holding.id} className="p-5">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex items-center gap-3">
                    <span className="grid h-10 w-10 place-items-center rounded-xl bg-[#31302F] text-xs font-bold text-[#F1F0F3]">
                      {holding.symbol.slice(0, 2)}
                    </span>
                    <span>
                      <strong className="block text-sm text-[#31302F]">{holding.symbol}</strong>
                      <span className="text-xs text-[#6E6C6F]">{formatNumber(holding.shares, 4)} shares</span>
                    </span>
                  </div>
                  <div className="text-right">
                    <strong className="block text-sm text-[#31302F]">{formatCurrency(holding.currentValue)}</strong>
                    <span className="text-xs font-semibold text-[#6E6C6F]">{holding.allocationPercent.toFixed(1)}% allocation</span>
                  </div>
                </div>
                <dl className="mt-4 grid grid-cols-2 gap-3 rounded-xl bg-[#E3DEDE] p-3 text-xs">
                  <div><dt className="text-[#6E6C6F]">Average cost</dt><dd className="mt-1 font-bold text-[#31302F]">{formatCurrency(holding.averageCost)}</dd></div>
                  <div><dt className="text-[#6E6C6F]">Unrealized return</dt><dd className="mt-1"><ReturnValue holding={holding} /></dd></div>
                </dl>
                <div className="mt-3 flex gap-2">
                  <button type="button" onClick={() => onEdit(holding)} className="min-h-10 flex-1 rounded-lg border border-[#CACDDC] text-xs font-bold text-[#31302F]">Edit</button>
                  <button type="button" onClick={() => onDelete(holding)} className="min-h-10 flex-1 rounded-lg border border-[#CACDDC] text-xs font-bold text-[#6E6C6F]">Remove</button>
                </div>
              </li>
            ))}
          </ul>
        </>
      )}
    </section>
  )
}
