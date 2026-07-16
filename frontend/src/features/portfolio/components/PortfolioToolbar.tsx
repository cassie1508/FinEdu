import type { PortfolioListItem } from '../../../types'

interface PortfolioToolbarProps {
  portfolios: PortfolioListItem[]
  activePortfolioId: string
  previewMode: boolean
  onPortfolioChange: (portfolioId: string) => void
  onCreatePortfolio: () => void
  onDeletePortfolio: () => void
}

export function PortfolioToolbar({
  portfolios,
  activePortfolioId,
  previewMode,
  onPortfolioChange,
  onCreatePortfolio,
  onDeletePortfolio,
}: PortfolioToolbarProps) {
  return (
    <header className="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <div className="mb-3 flex flex-wrap items-center gap-2">
          <span className="rounded-full bg-[#E3DEDE] px-3 py-1 text-xs font-bold uppercase tracking-[0.16em] text-[#6E6C6F]">
            Portfolio workspace
          </span>
          {previewMode && (
            <span className="rounded-full border border-[#A5A4A8] px-3 py-1 text-xs font-bold text-[#6E6C6F]">
              Preview data
            </span>
          )}
        </div>
        <h1 className="text-3xl font-bold tracking-[-0.04em] text-[#31302F] sm:text-4xl">
          Your investments, clearly understood.
        </h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-[#6E6C6F] sm:text-base">
          Track what you own, understand your allocation, and spot portfolio risks without the noise.
        </p>
      </div>

      <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
        {portfolios.length > 0 && (
          <label className="min-w-0 sm:min-w-56">
            <span className="sr-only">Choose portfolio</span>
            <select
              value={activePortfolioId}
              onChange={(event) => onPortfolioChange(event.target.value)}
              className="min-h-11 w-full rounded-xl border border-[#CACDDC] bg-[#F1F0F3] px-4 text-sm font-semibold text-[#31302F]"
            >
              {portfolios.map((portfolio) => (
                <option value={portfolio.id} key={portfolio.id}>
                  {portfolio.name}
                </option>
              ))}
            </select>
          </label>
        )}
        <button
          type="button"
          onClick={onCreatePortfolio}
          className="min-h-11 rounded-xl bg-[#31302F] px-5 text-sm font-bold text-[#F1F0F3] transition-colors hover:bg-[#6E6C6F]"
        >
          + New portfolio
        </button>
        {activePortfolioId && (
          <button
            type="button"
            onClick={onDeletePortfolio}
            className="min-h-11 rounded-xl border border-[#A5A4A8] px-4 text-sm font-bold text-[#6E6C6F] transition-colors hover:bg-[#E3DEDE] hover:text-[#31302F]"
          >
            Delete
          </button>
        )}
      </div>
    </header>
  )
}
