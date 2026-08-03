import { useEffect, useState, type FormEvent, type ReactNode } from 'react'
import type { PortfolioHolding } from '../../../types'

interface DialogShellProps {
  open: boolean
  title: string
  description: string
  children: ReactNode
  onClose: () => void
}

function DialogShell({ open, title, description, children, onClose }: DialogShellProps) {
  useEffect(() => {
    if (!open) return
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [open, onClose])

  if (!open) return null

  return (
    <div role="presentation" className="fixed inset-0 z-50 grid place-items-end bg-[#6E6C6F] p-0 sm:place-items-center sm:p-6" onMouseDown={onClose}>
      <section
        role="dialog"
        aria-modal="true"
        aria-labelledby="portfolio-dialog-title"
        className="max-h-[92vh] w-full overflow-y-auto rounded-t-3xl border border-[#CACDDC] bg-[#F1F0F3] p-5 sm:max-w-lg sm:rounded-3xl sm:p-7"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div className="flex items-start justify-between gap-5">
          <div>
            <h2 id="portfolio-dialog-title" className="text-2xl font-bold tracking-[-0.03em] text-[#31302F]">{title}</h2>
            <p className="mt-2 text-sm leading-6 text-[#6E6C6F]">{description}</p>
          </div>
          <button type="button" onClick={onClose} aria-label="Close dialog" className="grid h-10 w-10 shrink-0 place-items-center rounded-full border border-[#CACDDC] text-lg text-[#6E6C6F] hover:bg-[#E3DEDE]">×</button>
        </div>
        {children}
      </section>
    </div>
  )
}

const inputClass = 'mt-2 min-h-11 w-full rounded-xl border border-[#CACDDC] bg-[#F1F0F3] px-4 text-sm text-[#31302F] placeholder:text-[#A5A4A8]'

interface PortfolioFormDialogProps {
  open: boolean
  busy: boolean
  error: string
  onClose: () => void
  onSubmit: (name: string, description: string) => Promise<void>
}

export function PortfolioFormDialog({ open, busy, error, onClose, onSubmit }: PortfolioFormDialogProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

  useEffect(() => {
    if (open) {
      setName('')
      setDescription('')
    }
  }, [open])

  const submit = async (event: FormEvent) => {
    event.preventDefault()
    await onSubmit(name.trim(), description.trim())
  }

  return (
    <DialogShell open={open} title="Create a portfolio" description="Use a clear goal-based name so each portfolio is easy to recognize." onClose={onClose}>
      <form onSubmit={submit} className="mt-6 grid gap-5">
        <label className="text-sm font-bold text-[#31302F]">Portfolio name
          <input autoFocus required maxLength={100} value={name} onChange={(event) => setName(event.target.value)} placeholder="e.g. Long-term growth" className={inputClass} />
        </label>
        <label className="text-sm font-bold text-[#31302F]">Description <span className="font-normal text-[#6E6C6F]">(optional)</span>
          <textarea value={description} onChange={(event) => setDescription(event.target.value)} placeholder="What is this portfolio for?" rows={3} className={`${inputClass} py-3`} />
        </label>
        {error && <p role="alert" className="rounded-xl bg-[#E3DEDE] p-3 text-sm font-semibold text-[#31302F]">{error}</p>}
        <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <button type="button" onClick={onClose} className="min-h-11 rounded-xl px-5 text-sm font-bold text-[#6E6C6F] hover:bg-[#E3DEDE]">Cancel</button>
          <button disabled={busy || !name.trim()} className="min-h-11 rounded-xl bg-[#31302F] px-5 text-sm font-bold text-[#F1F0F3] disabled:cursor-not-allowed disabled:bg-[#A5A4A8]">{busy ? 'Creating…' : 'Create portfolio'}</button>
        </div>
      </form>
    </DialogShell>
  )
}

interface HoldingFormDialogProps {
  open: boolean
  busy: boolean
  error: string
  holding: PortfolioHolding | null
  onClose: () => void
  onSubmit: (symbol: string, shares: number, averageCost: number) => Promise<void>
}

export function HoldingFormDialog({ open, busy, error, holding, onClose, onSubmit }: HoldingFormDialogProps) {
  const [symbol, setSymbol] = useState('')
  const [shares, setShares] = useState('')
  const [averageCost, setAverageCost] = useState('')

  useEffect(() => {
    if (!open) return
    setSymbol(holding?.symbol ?? '')
    setShares(holding ? String(holding.shares) : '')
    setAverageCost(holding ? String(holding.averageCost) : '')
  }, [holding, open])

  const submit = async (event: FormEvent) => {
    event.preventDefault()
    await onSubmit(symbol.trim().toUpperCase(), Number(shares), Number(averageCost))
  }

  return (
    <DialogShell open={open} title={holding ? `Edit ${holding.symbol}` : 'Add a holding'} description="Enter your position and average purchase price. Market value is calculated from the latest quote." onClose={onClose}>
      <form onSubmit={submit} className="mt-6 grid gap-5">
        <label className="text-sm font-bold text-[#31302F]">Ticker symbol
          <input autoFocus={!holding} required maxLength={10} disabled={Boolean(holding)} value={symbol} onChange={(event) => setSymbol(event.target.value.replace(/[^a-zA-Z.-]/g, '').toUpperCase())} placeholder="AAPL" className={`${inputClass} uppercase disabled:bg-[#E3DEDE] disabled:text-[#6E6C6F]`} />
        </label>
        <div className="grid gap-4 sm:grid-cols-2">
          <label className="text-sm font-bold text-[#31302F]">Shares
            <input required min="0.0001" step="any" inputMode="decimal" type="number" value={shares} onChange={(event) => setShares(event.target.value)} placeholder="10" className={inputClass} />
          </label>
          <label className="text-sm font-bold text-[#31302F]">Average cost
            <span className="relative block"><span className="pointer-events-none absolute left-4 top-1/2 mt-1 -translate-y-1/2 text-sm text-[#6E6C6F]">$</span><input required min="0" step="any" inputMode="decimal" type="number" value={averageCost} onChange={(event) => setAverageCost(event.target.value)} placeholder="0.00" className={`${inputClass} pl-8`} /></span>
          </label>
        </div>
        {error && <p role="alert" className="rounded-xl bg-[#E3DEDE] p-3 text-sm font-semibold text-[#31302F]">{error}</p>}
        <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <button type="button" onClick={onClose} className="min-h-11 rounded-xl px-5 text-sm font-bold text-[#6E6C6F] hover:bg-[#E3DEDE]">Cancel</button>
          <button disabled={busy || !symbol || Number(shares) <= 0 || Number(averageCost) < 0} className="min-h-11 rounded-xl bg-[#31302F] px-5 text-sm font-bold text-[#F1F0F3] disabled:cursor-not-allowed disabled:bg-[#A5A4A8]">{busy ? 'Saving…' : holding ? 'Save changes' : 'Add holding'}</button>
        </div>
      </form>
    </DialogShell>
  )
}

interface ConfirmDialogProps {
  open: boolean
  busy: boolean
  title: string
  description: string
  confirmLabel: string
  onClose: () => void
  onConfirm: () => Promise<void>
}

export function ConfirmDialog({ open, busy, title, description, confirmLabel, onClose, onConfirm }: ConfirmDialogProps) {
  return (
    <DialogShell open={open} title={title} description={description} onClose={onClose}>
      <div className="mt-7 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
        <button type="button" onClick={onClose} className="min-h-11 rounded-xl px-5 text-sm font-bold text-[#6E6C6F] hover:bg-[#E3DEDE]">Cancel</button>
        <button type="button" disabled={busy} onClick={() => void onConfirm()} className="min-h-11 rounded-xl bg-[#31302F] px-5 text-sm font-bold text-[#F1F0F3] disabled:bg-[#A5A4A8]">{busy ? 'Removing…' : confirmLabel}</button>
      </div>
    </DialogShell>
  )
}
