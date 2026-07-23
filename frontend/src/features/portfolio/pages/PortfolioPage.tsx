import { useCallback, useEffect, useMemo, useState } from 'react'
import { api } from '../../../lib/api'
import { isSupabaseConfigured } from '../../../lib/supabaseClient'
import type {
  Portfolio,
  PortfolioDetail,
  PortfolioHolding,
  PortfolioListItem,
  PortfolioListResponse,
  PortfolioRisk,
} from '../../../types'
import { AllocationChart } from '../components/AllocationChart'
import { HoldingsTable } from '../components/HoldingsTable'
import {
  ConfirmDialog,
  HoldingFormDialog,
  PortfolioFormDialog,
} from '../components/PortfolioDialogs'
import { PortfolioToolbar } from '../components/PortfolioToolbar'
import { RiskPanel } from '../components/RiskPanel'
import { SummaryCards } from '../components/SummaryCards'
import {
  previewPortfolioDetails,
  previewPortfolioList,
  previewPortfolioRisk,
} from '../lib/portfolioPreviewData'
import { recalculatePortfolio } from '../lib/portfolioFormatters'

type DialogState =
  | { type: 'none' }
  | { type: 'create-portfolio' }
  | { type: 'add-holding' }
  | { type: 'edit-holding'; holding: PortfolioHolding }
  | { type: 'delete-holding'; holding: PortfolioHolding }
  | { type: 'delete-portfolio' }

const copyPreviewDetails = () =>
  Object.fromEntries(
    Object.entries(previewPortfolioDetails).map(([id, detail]) => [
      id,
      { ...detail, holdings: detail.holdings.map((holding) => ({ ...holding })) },
    ]),
  )

function getMessage(error: unknown) {
  return error instanceof Error ? error.message : 'Something went wrong. Please try again.'
}

export function PortfolioPage() {
  const previewMode = !isSupabaseConfigured
  const [portfolios, setPortfolios] = useState<PortfolioListItem[]>(
    previewMode ? previewPortfolioList : [],
  )
  const [previewDetails, setPreviewDetails] = useState<Record<string, PortfolioDetail>>(
    copyPreviewDetails,
  )
  const [activePortfolioId, setActivePortfolioId] = useState(
    previewMode ? previewPortfolioList[0].id : '',
  )
  const [liveDetail, setLiveDetail] = useState<PortfolioDetail | null>(null)
  const [risk, setRisk] = useState<PortfolioRisk | null>(
    previewMode ? previewPortfolioRisk : null,
  )
  const [loading, setLoading] = useState(!previewMode)
  const [busy, setBusy] = useState(false)
  const [pageError, setPageError] = useState('')
  const [formError, setFormError] = useState('')
  const [dialog, setDialog] = useState<DialogState>({ type: 'none' })

  const activeDetail = previewMode
    ? previewDetails[activePortfolioId] ?? null
    : liveDetail

  const closeDialog = useCallback(() => {
    if (busy) return
    setDialog({ type: 'none' })
    setFormError('')
  }, [busy])

  const loadPortfolioList = useCallback(async () => {
    if (previewMode) return
    setLoading(true)
    setPageError('')
    try {
      const response = await api.get<PortfolioListResponse>('/api/v1/portfolio/portfolios')
      setPortfolios(response.portfolios)
      setActivePortfolioId((current) => {
        if (current && response.portfolios.some((portfolio) => portfolio.id === current)) return current
        return response.portfolios[0]?.id ?? ''
      })
    } catch (error) {
      setPageError(getMessage(error))
    } finally {
      setLoading(false)
    }
  }, [previewMode])

  const loadActivePortfolio = useCallback(async () => {
    if (previewMode || !activePortfolioId) return
    setLoading(true)
    setPageError('')
    try {
      const [detailResult, riskResult] = await Promise.allSettled([
        api.get<PortfolioDetail>(`/api/v1/portfolio/portfolios/${activePortfolioId}`),
        api.get<PortfolioRisk>(`/api/v1/portfolio/portfolios/${activePortfolioId}/risk`),
      ])

      if (detailResult.status === 'rejected') throw detailResult.reason
      setLiveDetail(detailResult.value)
      setRisk(riskResult.status === 'fulfilled' ? riskResult.value : null)
    } catch (error) {
      setPageError(getMessage(error))
      setLiveDetail(null)
    } finally {
      setLoading(false)
    }
  }, [activePortfolioId, previewMode])

  useEffect(() => {
    void loadPortfolioList()
  }, [loadPortfolioList])

  useEffect(() => {
    if (previewMode) {
      setRisk(activePortfolioId === 'preview-core-growth' ? previewPortfolioRisk : null)
      return
    }
    void loadActivePortfolio()
  }, [activePortfolioId, loadActivePortfolio, previewMode])

  const updatePreviewDetail = (portfolioId: string, updater: (detail: PortfolioDetail) => PortfolioDetail) => {
    const detail = previewDetails[portfolioId]
    if (!detail) return
    const nextDetail = updater(detail)
    setPreviewDetails((current) => ({ ...current, [portfolioId]: nextDetail }))
    setPortfolios((items) => items.map((item) => (
      item.id === portfolioId
        ? { ...item, holdingsCount: nextDetail.holdings.length }
        : item
    )))
  }

  const createPortfolio = async (name: string, description: string) => {
    setBusy(true)
    setFormError('')
    try {
      if (previewMode) {
        const id = `preview-${Date.now()}`
        const createdAt = new Date().toISOString()
        setPortfolios((current) => [
          { id, name, description, holdingsCount: 0, createdAt },
          ...current,
        ])
        setPreviewDetails((current) => ({
          ...current,
          [id]: {
            id,
            name,
            description,
            holdings: [],
            totalValue: 0,
            totalCost: 0,
            totalGainLoss: 0,
            totalGainLossPercent: 0,
          },
        }))
        setActivePortfolioId(id)
      } else {
        const portfolio = await api.post<Portfolio>('/api/v1/portfolio/portfolios', {
          name,
          description,
        })
        setPortfolios((current) => [
          {
            id: portfolio.id,
            name: portfolio.name,
            description: portfolio.description,
            holdingsCount: 0,
            createdAt: portfolio.createdAt,
          },
          ...current,
        ])
        setActivePortfolioId(portfolio.id)
      }
      setDialog({ type: 'none' })
    } catch (error) {
      setFormError(getMessage(error))
    } finally {
      setBusy(false)
    }
  }

  const saveHolding = async (symbol: string, shares: number, averageCost: number) => {
    if (!activePortfolioId) return
    setBusy(true)
    setFormError('')
    const editedHolding = dialog.type === 'edit-holding' ? dialog.holding : null

    try {
      if (previewMode) {
        updatePreviewDetail(activePortfolioId, (detail) => {
          if (editedHolding) {
            const holdings = detail.holdings.map((holding) => {
              if (holding.id !== editedHolding.id) return holding
              const currentValue = shares * holding.currentPrice
              const totalCost = shares * averageCost
              const unrealizedGainLoss = currentValue - totalCost
              return {
                ...holding,
                shares,
                averageCost,
                currentValue,
                unrealizedGainLoss,
                unrealizedGainLossPercent: totalCost > 0 ? (unrealizedGainLoss / totalCost) * 100 : 0,
              }
            })
            return recalculatePortfolio(detail, holdings)
          }

          if (detail.holdings.some((holding) => holding.symbol === symbol)) {
            throw new Error(`${symbol} is already in this portfolio.`)
          }
          const holding: PortfolioHolding = {
            id: `preview-holding-${Date.now()}`,
            portfolioId: activePortfolioId,
            symbol,
            shares,
            averageCost,
            currentPrice: averageCost,
            currentValue: shares * averageCost,
            unrealizedGainLoss: 0,
            unrealizedGainLossPercent: 0,
            allocationPercent: 0,
          }
          return recalculatePortfolio(detail, [...detail.holdings, holding])
        })
      } else if (editedHolding) {
        await api.put(`/api/v1/portfolio/portfolios/${activePortfolioId}/holdings/${editedHolding.id}`, {
          shares,
          averageCost,
        })
        await loadActivePortfolio()
      } else {
        await api.post(`/api/v1/portfolio/portfolios/${activePortfolioId}/holdings`, {
          symbol,
          shares,
          averageCost,
        })
        await loadActivePortfolio()
      }
      setDialog({ type: 'none' })
    } catch (error) {
      setFormError(getMessage(error))
    } finally {
      setBusy(false)
    }
  }

  const deleteHolding = async () => {
    if (dialog.type !== 'delete-holding' || !activePortfolioId) return
    setBusy(true)
    try {
      if (previewMode) {
        updatePreviewDetail(activePortfolioId, (detail) =>
          recalculatePortfolio(
            detail,
            detail.holdings.filter((holding) => holding.id !== dialog.holding.id),
          ),
        )
      } else {
        await api.delete(`/api/v1/portfolio/portfolios/${activePortfolioId}/holdings/${dialog.holding.id}`)
        await loadActivePortfolio()
      }
      setDialog({ type: 'none' })
    } catch (error) {
      setPageError(getMessage(error))
    } finally {
      setBusy(false)
    }
  }

  const deletePortfolio = async () => {
    if (!activePortfolioId) return
    setBusy(true)
    try {
      if (!previewMode) {
        await api.delete(`/api/v1/portfolio/portfolios/${activePortfolioId}`)
      }
      const nextPortfolios = portfolios.filter((portfolio) => portfolio.id !== activePortfolioId)
      setPortfolios(nextPortfolios)
      setPreviewDetails((current) => {
        const next = { ...current }
        delete next[activePortfolioId]
        return next
      })
      setLiveDetail(null)
      setActivePortfolioId(nextPortfolios[0]?.id ?? '')
      setDialog({ type: 'none' })
    } catch (error) {
      setPageError(getMessage(error))
    } finally {
      setBusy(false)
    }
  }

  const selectedPortfolio = useMemo(
    () => portfolios.find((portfolio) => portfolio.id === activePortfolioId),
    [activePortfolioId, portfolios],
  )

  return (
    <div className="pb-10">
      <PortfolioToolbar
        portfolios={portfolios}
        activePortfolioId={activePortfolioId}
        previewMode={previewMode}
        onPortfolioChange={setActivePortfolioId}
        onCreatePortfolio={() => setDialog({ type: 'create-portfolio' })}
        onDeletePortfolio={() => setDialog({ type: 'delete-portfolio' })}
      />

      {previewMode && (
        <p className="mt-5 rounded-xl border border-[#CACDDC] bg-[#E3DEDE] px-4 py-3 text-sm leading-6 text-[#6E6C6F]">
          You’re viewing an interactive preview because Supabase is not configured. Add your project keys to use authenticated portfolio data.
        </p>
      )}

      {pageError && (
        <div role="alert" className="mt-6 flex flex-col gap-3 rounded-2xl border border-[#A5A4A8] bg-[#E3DEDE] p-5 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="font-bold text-[#31302F]">We couldn’t load your portfolio</p>
            <p className="mt-1 text-sm text-[#6E6C6F]">{pageError}</p>
          </div>
          <button type="button" onClick={() => void loadPortfolioList()} className="min-h-11 rounded-xl border border-[#31302F] px-4 text-sm font-bold text-[#31302F]">Try again</button>
        </div>
      )}

      {loading ? (
        <div aria-label="Loading portfolio" className="mt-8 animate-pulse space-y-4">
          <div className="h-24 rounded-2xl bg-[#E3DEDE]" />
          <div className="grid gap-4 lg:grid-cols-2"><div className="h-80 rounded-2xl bg-[#E3DEDE]" /><div className="h-80 rounded-2xl bg-[#CACDDC]" /></div>
        </div>
      ) : !activeDetail && portfolios.length === 0 && !pageError ? (
        <section className="mt-10 rounded-3xl border border-[#CACDDC] bg-[#E3DEDE] px-6 py-16 text-center">
          <p className="text-xs font-bold uppercase tracking-[0.14em] text-[#6E6C6F]">Start here</p>
          <h2 className="mt-3 text-2xl font-bold tracking-[-0.03em] text-[#31302F]">Create your first portfolio</h2>
          <p className="mx-auto mt-3 max-w-md text-sm leading-6 text-[#6E6C6F]">Organize holdings around a goal, then FinEdu will calculate value, allocation, returns, and risk.</p>
          <button type="button" onClick={() => setDialog({ type: 'create-portfolio' })} className="mt-6 min-h-11 rounded-xl bg-[#31302F] px-5 text-sm font-bold text-[#F1F0F3]">Create portfolio</button>
        </section>
      ) : activeDetail ? (
        <div className="mt-8 grid gap-6">
          <div className="flex flex-col gap-2 border-b border-[#CACDDC] pb-5 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <h2 className="text-2xl font-bold tracking-[-0.03em] text-[#31302F]">{activeDetail.name}</h2>
              {activeDetail.description && <p className="mt-1 text-sm text-[#6E6C6F]">{activeDetail.description}</p>}
            </div>
            <p className="text-xs font-semibold text-[#6E6C6F]">Values shown in USD</p>
          </div>

          <SummaryCards portfolio={activeDetail} />

          <div className={`grid gap-6 ${activeDetail.holdings.length > 0 ? 'lg:grid-cols-[minmax(0,1.08fr)_minmax(330px,0.92fr)]' : ''}`}>
            {activeDetail.holdings.length > 0 && <AllocationChart holdings={activeDetail.holdings} />}
            <RiskPanel risk={risk} />
          </div>

          <HoldingsTable
            holdings={activeDetail.holdings}
            busy={busy}
            onAdd={() => setDialog({ type: 'add-holding' })}
            onEdit={(holding) => setDialog({ type: 'edit-holding', holding })}
            onDelete={(holding) => setDialog({ type: 'delete-holding', holding })}
          />
        </div>
      ) : null}

      <PortfolioFormDialog
        open={dialog.type === 'create-portfolio'}
        busy={busy}
        error={formError}
        onClose={closeDialog}
        onSubmit={createPortfolio}
      />
      <HoldingFormDialog
        open={dialog.type === 'add-holding' || dialog.type === 'edit-holding'}
        busy={busy}
        error={formError}
        holding={dialog.type === 'edit-holding' ? dialog.holding : null}
        onClose={closeDialog}
        onSubmit={saveHolding}
      />
      <ConfirmDialog
        open={dialog.type === 'delete-holding'}
        busy={busy}
        title={dialog.type === 'delete-holding' ? `Remove ${dialog.holding.symbol}?` : 'Remove holding?'}
        description="This removes the holding from the portfolio. It does not place a market order."
        confirmLabel="Remove holding"
        onClose={closeDialog}
        onConfirm={deleteHolding}
      />
      <ConfirmDialog
        open={dialog.type === 'delete-portfolio'}
        busy={busy}
        title={`Delete ${selectedPortfolio?.name ?? 'this portfolio'}?`}
        description="This permanently removes the portfolio and all of its holdings. This action cannot be undone."
        confirmLabel="Delete portfolio"
        onClose={closeDialog}
        onConfirm={deletePortfolio}
      />
    </div>
  )
}
