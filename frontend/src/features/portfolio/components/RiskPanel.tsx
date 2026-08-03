import type { PortfolioRisk } from '../../../types'

export function RiskPanel({ risk }: { risk: PortfolioRisk | null }) {
  const hasScore = risk?.healthScore != null

  return (
    <section className="rounded-2xl border border-[#CACDDC] bg-[#E3DEDE] p-5 sm:p-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-xs font-bold uppercase tracking-[0.14em] text-[#6E6C6F]">Risk snapshot</p>
          <h2 className="mt-1 text-xl font-bold tracking-[-0.03em] text-[#31302F]">Portfolio health</h2>
        </div>
        {hasScore && (
          <div className="grid h-20 w-20 shrink-0 place-items-center rounded-full border-[7px] border-[#A5A4A8] bg-[#F1F0F3] text-center">
            <span>
              <strong className="block text-2xl leading-none text-[#31302F]">{Math.round(risk.healthScore!)}</strong>
              <span className="text-[10px] font-bold uppercase text-[#6E6C6F]">of 100</span>
            </span>
          </div>
        )}
      </div>

      {!hasScore ? (
        <div className="mt-8 rounded-xl border border-[#A5A4A8] bg-[#F1F0F3] p-5">
          <p className="font-bold text-[#31302F]">Analysis is getting ready</p>
          <p className="mt-2 text-sm leading-6 text-[#6E6C6F]">
            {risk?.message ?? 'Add holdings to see concentration, diversification, and practical risk guidance.'}
          </p>
        </div>
      ) : (
        <>
          <div className="mt-5 grid grid-cols-2 gap-3">
            <div className="rounded-xl bg-[#F1F0F3] p-4">
              <p className="text-xs font-semibold text-[#6E6C6F]">Risk level</p>
              <p className="mt-1 font-bold text-[#31302F]">{risk.riskLevel}</p>
            </div>
            <div className="rounded-xl bg-[#F1F0F3] p-4">
              <p className="text-xs font-semibold text-[#6E6C6F]">Diversification</p>
              <p className="mt-1 font-bold text-[#31302F]">{Math.round(risk.diversificationScore ?? 0)}/100</p>
            </div>
          </div>

          {risk.sectorConcentration.length > 0 && (
            <div className="mt-6">
              <h3 className="text-sm font-bold text-[#31302F]">Sector concentration</h3>
              <ul className="mt-3 grid gap-3">
                {risk.sectorConcentration.map((entry) => (
                  <li key={entry.sector}>
                    <div className="mb-1.5 flex justify-between text-xs font-semibold text-[#6E6C6F]">
                      <span>{entry.sector}</span>
                      <span>{entry.percent.toFixed(0)}%</span>
                    </div>
                    <div className="h-2 overflow-hidden rounded-full bg-[#CACDDC]">
                      <div className="h-full rounded-full bg-[#31302F]" style={{ width: `${Math.min(entry.percent, 100)}%` }} />
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {risk.recommendations.length > 0 && (
            <div className="mt-6 border-t border-[#CACDDC] pt-5">
              <h3 className="text-sm font-bold text-[#31302F]">What to consider</h3>
              <ul className="mt-3 grid gap-3">
                {risk.recommendations.map((recommendation) => (
                  <li key={recommendation} className="flex gap-3 text-sm leading-5 text-[#6E6C6F]">
                    <span className="font-bold text-[#31302F]">→</span>
                    <span>{recommendation}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {risk.message && <p className="mt-5 text-xs leading-5 text-[#6E6C6F]">{risk.message}</p>}
        </>
      )}
    </section>
  )
}
