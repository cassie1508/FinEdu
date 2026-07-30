import { useMemo, useState, type FormEvent } from 'react'
import {
  Area,
  CartesianGrid,
  ComposedChart,
  Customized,
  Line,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  type TooltipProps,
} from 'recharts'
import { usePriceHistory } from '../hooks/usePriceHistory'
import { CHART_RANGES, type CandlePoint, type ChartRange } from '../types'
import { formatAxisTick, formatChange, formatPrice, formatTooltipDate, formatVolume } from '../lib/format'

type ChartType = 'line' | 'candlestick'

// The SMA overlay needs enough bars to be meaningful — 1W/1M only have a
// handful of daily bars (SMA(10) barely has room to compute on 1W), so
// Indicators stays restricted to the weekly-resolution ranges. Candlestick
// view has no such constraint (every range already carries full OHLC) and
// is available everywhere.
const INDICATOR_RANGES: ChartRange[] = ['6M', '1Y', '5Y']

interface StockChartProps {
  symbol: string
  onSymbolChange: (symbol: string) => void
}

// Symbol is controlled by the parent (MarketDashboardPage) rather than
// owned internally — the news panel next to this chart needs to know which
// company is currently displayed so it can fetch news for the same ticker.
export function StockChart({ symbol, onSymbolChange }: StockChartProps) {
  const [symbolInput, setSymbolInput] = useState('')
  const [range, setRange] = useState<ChartRange>('1M')
  const [chartType, setChartType] = useState<ChartType>('line')
  const [showIndicators, setShowIndicators] = useState(false)

  const { data, isLoading, error } = usePriceHistory(symbol, range)

  const supportsIndicators = INDICATOR_RANGES.includes(range)
  const effectiveShowIndicators = supportsIndicators && showIndicators

  function handleSymbolSubmit(event: FormEvent) {
    event.preventDefault()
    const trimmed = symbolInput.trim().toUpperCase()
    if (trimmed) {
      onSymbolChange(trimmed)
      setSymbolInput('')
    }
  }

  const candles = data?.candles ?? []
  const latest = candles.at(-1)
  const previous = candles.at(-2)
  const change = latest && previous ? latest.c - previous.c : 0
  const changePercent = latest && previous && previous.c !== 0 ? (change / previous.c) * 100 : 0
  const isPositive = change >= 0

  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-800 dark:bg-gray-900">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <TrendIcon className="h-4 w-4 text-indigo-500" />
          <h2 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
            Interactive Stock Chart
          </h2>
        </div>

        <form onSubmit={handleSymbolSubmit} className="relative">
          <SearchIcon className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-gray-400" />
          <input
            value={symbolInput}
            onChange={(e) => setSymbolInput(e.target.value)}
            placeholder="Search symbol..."
            className="w-40 rounded-lg border border-gray-200 bg-gray-50 py-1.5 pr-3 pl-8 text-xs text-gray-900 placeholder:text-gray-400 focus:border-indigo-300 focus:bg-white focus:ring-2 focus:ring-indigo-100 focus:outline-none dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100"
          />
        </form>
      </div>

      <div className="mt-4 flex items-center justify-between">
        <div className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-gray-900 text-xs font-bold text-white dark:bg-gray-100 dark:text-gray-900">
            {symbol.charAt(0)}
          </div>
          <span className="text-base font-bold text-gray-900 dark:text-gray-100">{symbol}</span>
          <ChevronDownIcon className="h-3.5 w-3.5 text-gray-400" />
        </div>

        {latest && (
          <div className="flex items-baseline gap-2">
            <span className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {formatPrice(latest.c)}
            </span>
            <span
              className={`text-sm font-medium ${
                isPositive ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-500 dark:text-rose-400'
              }`}
            >
              {formatChange(change, changePercent)}
            </span>
          </div>
        )}
      </div>

      <div className="mt-4 flex items-center justify-between border-b border-gray-100 pb-3 dark:border-gray-800">
        <div className="flex items-center gap-1">
          {CHART_RANGES.map((r) => (
            <button
              key={r}
              type="button"
              onClick={() => setRange(r)}
              className={`rounded-md px-2 py-1 text-xs font-medium transition ${
                r === range
                  ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400'
                  : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
              }`}
            >
              {r}
            </button>
          ))}
        </div>

        <div className="flex items-center gap-1.5">
          <button
            type="button"
            onClick={() => setChartType('candlestick')}
            title="Candlestick view"
            className={`rounded p-1 transition ${
              chartType === 'candlestick'
                ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400'
                : 'text-gray-400 hover:text-gray-700 dark:text-gray-500 dark:hover:text-gray-300'
            }`}
          >
            <CandlestickIcon className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => setChartType('line')}
            title="Line view"
            className={`rounded p-1 transition ${
              chartType === 'line'
                ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400'
                : 'text-gray-400 hover:text-gray-700 dark:text-gray-500 dark:hover:text-gray-300'
            }`}
          >
            <LineChartIcon className="h-4 w-4" />
          </button>
          <span className="mx-0.5 h-4 w-px bg-gray-200 dark:bg-gray-700" />
          <button
            type="button"
            onClick={() => setShowIndicators((v) => !v)}
            disabled={!supportsIndicators}
            title={
              supportsIndicators
                ? 'Toggle SMA (10) overlay'
                : 'Indicators are available on 6M, 1Y, and 5Y'
            }
            className={`flex items-center gap-1 rounded px-1.5 py-1 text-[11px] font-medium transition ${
              effectiveShowIndicators
                ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400'
                : supportsIndicators
                  ? 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
                  : 'cursor-not-allowed text-gray-300 dark:text-gray-700'
            }`}
          >
            <IndicatorsIcon className="h-3.5 w-3.5" />
            Indicators
          </button>
        </div>
      </div>

      <div className="mt-3 h-80">
        {isLoading && (
          <div className="flex h-full items-center justify-center">
            <div className="h-full w-full animate-pulse rounded-xl bg-gray-100 dark:bg-gray-800" />
          </div>
        )}

        {!isLoading && error && (
          <div className="flex h-full flex-col items-center justify-center gap-1 text-center">
            <p className="text-sm font-medium text-red-600 dark:text-red-400">
              Couldn't load price history
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">{error}</p>
          </div>
        )}

        {!isLoading && !error && candles.length === 0 && (
          <div className="flex h-full items-center justify-center text-sm text-gray-500 dark:text-gray-400">
            No price data available for {symbol}.
          </div>
        )}

        {!isLoading && !error && candles.length > 0 && latest && (
          <PriceAreaChart
            candles={candles}
            range={range}
            latestClose={latest.c}
            chartType={chartType}
            showIndicators={effectiveShowIndicators}
          />
        )}
      </div>

      {latest && (
        <div className="mt-4 grid grid-cols-3 gap-y-3 border-t border-gray-100 pt-4 sm:grid-cols-6 dark:border-gray-800">
          <Stat label="Open" value={formatPrice(latest.o)} />
          <Stat label="High" value={formatPrice(latest.h)} />
          <Stat label="Low" value={formatPrice(latest.l)} />
          <Stat label="Volume" value={formatVolume(latest.v)} />
          <Stat label="Market Cap" value="—" />
          <Stat label="P/E Ratio" value="—" />
        </div>
      )}

      {data && (
        <p className="mt-3 text-[11px] text-gray-400 dark:text-gray-500">
          Data delayed ~{data.delayedMinutes} min
        </p>
      )}
    </div>
  )
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-xs text-gray-500 dark:text-gray-400">{label}</p>
      <p className="mt-0.5 text-sm font-semibold text-gray-900 dark:text-gray-100">{value}</p>
    </div>
  )
}

// FinEdu's brand blue/lavender (matches the primary button/accent color used
// on the login page). The chart line is a consistent brand color regardless
// of gain/loss — only the +/- change text next to the price switches
// between green and red.
const CHART_LINE_COLOR = '#8EB3E5'

// Period for the "Indicators" SMA overlay. Fixed rather than configurable —
// there's no UI for choosing a period, this is a single fixed indicator.
const SMA_PERIOD = 10
const SMA_COLOR = '#F59E0B'

function computeSma(candles: CandlePoint[], period: number): (number | undefined)[] {
  const result: (number | undefined)[] = []
  let sum = 0
  for (let i = 0; i < candles.length; i++) {
    sum += candles[i].c
    if (i >= period) sum -= candles[i - period].c
    result.push(i >= period - 1 ? sum / period : undefined)
  }
  return result
}

function PriceAreaChart({
  candles,
  range,
  latestClose,
  chartType,
  showIndicators,
}: {
  candles: CandlePoint[]
  range: ChartRange
  latestClose: number
  chartType: ChartType
  showIndicators: boolean
}) {
  const strokeColor = CHART_LINE_COLOR

  const chartData = useMemo(() => {
    const sma = computeSma(candles, SMA_PERIOD)
    return candles.map((candle, i) => ({ ...candle, sma: sma[i] }))
  }, [candles])

  return (
    <ResponsiveContainer width="100%" height="100%">
      <ComposedChart data={chartData} margin={{ top: 8, right: 70, bottom: 0, left: 0 }}>
        <defs>
          <linearGradient id="stockChartGradient" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor={strokeColor} stopOpacity={0.28} />
            <stop offset="100%" stopColor={strokeColor} stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid vertical={false} stroke="#F1F1F5" />
        <XAxis
          dataKey="t"
          tickFormatter={(t: number) => formatAxisTick(t, range)}
          tick={{ fontSize: 11, fill: '#A6A9C4' }}
          tickLine={false}
          axisLine={false}
          minTickGap={40}
        />
        <YAxis
          domain={['auto', 'auto']}
          tick={{ fontSize: 11, fill: '#A6A9C4' }}
          tickLine={false}
          axisLine={false}
          width={56}
          tickFormatter={(v: number) => formatPrice(v)}
        />
        <Tooltip content={<PriceTooltip />} />

        {chartType === 'line' ? (
          <Area
            type="monotone"
            dataKey="c"
            stroke={strokeColor}
            strokeWidth={1.75}
            fill="url(#stockChartGradient)"
            dot={(props: { cx: number; cy: number; index: number; key: string }) =>
              props.index === candles.length - 1 ? (
                <circle
                  key={props.key}
                  cx={props.cx}
                  cy={props.cy}
                  r={4}
                  fill={strokeColor}
                  stroke="#ffffff"
                  strokeWidth={2}
                />
              ) : (
                <g key={props.key} />
              )
            }
          />
        ) : (
          <>
            {/* Invisible series bound to the high/low fields purely so the
                YAxis 'auto' domain covers the full wick range — Customized
                below draws the actual candles and doesn't register a
                dataKey the axis can scale against on its own. */}
            <Line dataKey="h" stroke="transparent" dot={false} activeDot={false} isAnimationActive={false} />
            <Line dataKey="l" stroke="transparent" dot={false} activeDot={false} isAnimationActive={false} />
            <Customized component={CandlestickLayer} candles={candles} />
          </>
        )}

        {showIndicators && (
          <Line
            type="monotone"
            dataKey="sma"
            stroke={SMA_COLOR}
            strokeWidth={1.5}
            dot={false}
            activeDot={false}
            connectNulls
            isAnimationActive={false}
          />
        )}

        <ReferenceLine y={latestClose} stroke="#E5E7EB" strokeDasharray="4 4" strokeOpacity={0.8} />
        <ReferenceLine
          y={latestClose}
          stroke="transparent"
          label={<PriceBadge value={latestClose} />}
        />
      </ComposedChart>
    </ResponsiveContainer>
  )
}

interface ChartAxisScale {
  scale: (value: number) => number
}

// Recharts' parent chart wrapper clones this element with its internal
// state (xAxisMap/yAxisMap, among others) merged into its props — see
// generateCategoricalChart's renderCustomized — so `candles` (passed
// explicitly above) and the injected axis scales both end up here together.
function CandlestickLayer(props: {
  candles?: CandlePoint[]
  xAxisMap?: Record<string, ChartAxisScale>
  yAxisMap?: Record<string, ChartAxisScale>
}) {
  const { candles, xAxisMap, yAxisMap } = props
  const xAxis = xAxisMap ? Object.values(xAxisMap)[0] : undefined
  const yAxis = yAxisMap ? Object.values(yAxisMap)[0] : undefined

  if (!candles || candles.length === 0 || !xAxis || !yAxis) return null

  const xPositions = candles.map((candle) => xAxis.scale(candle.t))
  const step = xPositions.length > 1 ? Math.abs(xPositions[1] - xPositions[0]) : 8
  const bodyWidth = Math.max(Math.min(step * 0.6, 14), 2)

  return (
    <g>
      {candles.map((candle, i) => {
        const x = xPositions[i]
        const isUp = candle.c >= candle.o
        const color = isUp ? '#10B981' : '#F87171'
        const yOpen = yAxis.scale(candle.o)
        const yClose = yAxis.scale(candle.c)
        const bodyTop = Math.min(yOpen, yClose)
        const bodyHeight = Math.max(Math.abs(yClose - yOpen), 1)

        return (
          <g key={candle.t}>
            <line
              x1={x}
              x2={x}
              y1={yAxis.scale(candle.h)}
              y2={yAxis.scale(candle.l)}
              stroke={color}
              strokeWidth={1}
            />
            <rect x={x - bodyWidth / 2} y={bodyTop} width={bodyWidth} height={bodyHeight} fill={color} />
          </g>
        )
      })}
    </g>
  )
}

interface RechartsLabelViewBox {
  x: number
  y: number
  width: number
  height: number
}

const PRICE_BADGE_FILL = '#1F1F1F'
const PRICE_BADGE_TEXT = '#ffffff'
const PRICE_BADGE_GAP = 12

// Recharts clones this element and injects `viewBox` for the reference
// line's position, so it can draw a price badge as its own floating chip —
// offset from the plot's right edge (not touching the line/dot) — instead
// of Recharts' default plain-text label.
function PriceBadge({ viewBox, value }: { viewBox?: RechartsLabelViewBox; value: number }) {
  if (!viewBox) return null

  const text = formatPrice(value)
  const paddingX = 7
  const badgeWidth = text.length * 6.4 + paddingX * 2
  const badgeHeight = 20
  const x = viewBox.x + viewBox.width + PRICE_BADGE_GAP
  const y = viewBox.y - badgeHeight / 2

  return (
    <g>
      <rect x={x} y={y} width={badgeWidth} height={badgeHeight} rx={4} fill={PRICE_BADGE_FILL} />
      <text
        x={x + badgeWidth / 2}
        y={y + badgeHeight / 2 + 4}
        textAnchor="middle"
        fontSize={11}
        fontWeight={600}
        fill={PRICE_BADGE_TEXT}
      >
        {text}
      </text>
    </g>
  )
}

function PriceTooltip({ active, payload }: TooltipProps<number, string>) {
  if (!active || !payload?.length) return null
  const candle = payload[0].payload as CandlePoint

  return (
    <div className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs shadow-md dark:border-gray-700 dark:bg-gray-800">
      <p className="font-medium text-gray-900 dark:text-gray-100">{formatTooltipDate(candle.t)}</p>
      <p className="mt-1 text-gray-600 dark:text-gray-300">Close: {formatPrice(candle.c)}</p>
    </div>
  )
}

function TrendIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} xmlns="http://www.w3.org/2000/svg">
      <path
        d="M3 17L9 11L13 15L21 7"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path d="M15 7H21V13" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function SearchIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} xmlns="http://www.w3.org/2000/svg">
      <circle cx="11" cy="11" r="7" stroke="currentColor" strokeWidth="1.8" />
      <path d="M20 20L16.5 16.5" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  )
}

function ChevronDownIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} xmlns="http://www.w3.org/2000/svg">
      <path
        d="M6 9L12 15L18 9"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function CandlestickIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} xmlns="http://www.w3.org/2000/svg">
      <path d="M6 4V20" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <rect x="3.5" y="8" width="5" height="7" rx="1" stroke="currentColor" strokeWidth="1.5" />
      <path d="M14 2V22" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <rect x="11.5" y="10" width="5" height="6" rx="1" stroke="currentColor" strokeWidth="1.5" />
      <path d="M20 6V18" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <rect x="17.5" y="9" width="5" height="5" rx="1" stroke="currentColor" strokeWidth="1.5" />
    </svg>
  )
}

function LineChartIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} xmlns="http://www.w3.org/2000/svg">
      <path
        d="M3 15L9 9L13 13L21 5"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function IndicatorsIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} xmlns="http://www.w3.org/2000/svg">
      <path d="M4 6H20M4 12H14M4 18H17" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
      <circle cx="17" cy="6" r="1.5" fill="currentColor" />
      <circle cx="16" cy="12" r="1.5" fill="currentColor" />
      <circle cx="19" cy="18" r="1.5" fill="currentColor" />
    </svg>
  )
}
