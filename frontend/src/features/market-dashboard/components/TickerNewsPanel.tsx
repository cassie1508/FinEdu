import { useEffect, useState } from 'react'
import { useTickerNews } from '../hooks/useTickerNews'
import { useTickerNewsSummary } from '../hooks/useTickerNewsSummary'
import type { NewsArticle, TickerNewsSummary } from '../types'
import { formatRelativeTime } from '../lib/format'

interface TickerNewsPanelProps {
  symbol: string
}

const ARTICLES_PER_PAGE = 5

// Companion to StockChart on the Market Dashboard — unlike the old
// News & Charts general feed, this is scoped to whichever ticker is
// currently displayed in the chart to the left.
export function TickerNewsPanel({ symbol }: TickerNewsPanelProps) {
  const { summary, isLoading: isSummaryLoading, error: summaryError } = useTickerNewsSummary(symbol)
  const { articles, isLoading, error } = useTickerNews(symbol)

  const [page, setPage] = useState(1)
  // Each symbol change re-fetches a fresh article list, so the page index
  // from the previous ticker shouldn't carry over.
  useEffect(() => setPage(1), [symbol])

  const totalPages = Math.max(1, Math.ceil(articles.length / ARTICLES_PER_PAGE))
  const currentPage = Math.min(page, totalPages)
  const pageArticles = articles.slice(
    (currentPage - 1) * ARTICLES_PER_PAGE,
    currentPage * ARTICLES_PER_PAGE,
  )

  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-800 dark:bg-gray-900">
      <div className="flex items-center gap-2">
        <NewsIcon className="h-4 w-4 text-indigo-500" />
        <h2 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
          {symbol} News
        </h2>
      </div>

      {isSummaryLoading && (
        <div className="mt-4 space-y-2">
          <div className="h-3 w-24 animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
          <div className="h-3.5 w-full animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
          <div className="h-3.5 w-4/5 animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
        </div>
      )}

      {!isSummaryLoading && summary && <DigestCallout summary={summary} />}
      {!isSummaryLoading && !summary && <DigestUnavailable error={summaryError} />}

      <div className="mt-2 flex flex-col divide-y divide-gray-100 border-t border-gray-100 dark:divide-gray-800 dark:border-gray-800">
        {isLoading &&
          Array.from({ length: 4 }).map((_, i) => <ArticleSkeleton key={i} />)}

        {!isLoading && error && (
          <div className="flex flex-col items-center gap-1 py-10 text-center">
            <p className="text-sm font-medium text-red-600 dark:text-red-400">
              Couldn't load news
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">{error}</p>
          </div>
        )}

        {!isLoading && !error && articles.length === 0 && (
          <div className="py-10 text-center text-sm text-gray-500 dark:text-gray-400">
            No news available for {symbol} right now.
          </div>
        )}

        {!isLoading &&
          !error &&
          pageArticles.map((article) => <ArticleCard key={article.id} article={article} />)}
      </div>

      {!isLoading && !error && articles.length > ARTICLES_PER_PAGE && (
        <Pagination currentPage={currentPage} totalPages={totalPages} onPageChange={setPage} />
      )}
    </div>
  )
}

function Pagination({
  currentPage,
  totalPages,
  onPageChange,
}: {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
}) {
  return (
    <div className="mt-3 flex items-center justify-center gap-1 border-t border-gray-100 pt-3 dark:border-gray-800">
      <button
        type="button"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="rounded-md px-2 py-1 text-xs font-medium text-gray-500 transition hover:text-gray-800 disabled:cursor-not-allowed disabled:text-gray-300 dark:text-gray-400 dark:hover:text-gray-200 dark:disabled:text-gray-700"
      >
        Prev
      </button>
      {Array.from({ length: totalPages }, (_, i) => i + 1).map((pageNumber) => (
        <button
          key={pageNumber}
          type="button"
          onClick={() => onPageChange(pageNumber)}
          className={`h-6 w-6 rounded-md text-xs font-medium transition ${
            pageNumber === currentPage
              ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400'
              : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
          }`}
        >
          {pageNumber}
        </button>
      ))}
      <button
        type="button"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="rounded-md px-2 py-1 text-xs font-medium text-gray-500 transition hover:text-gray-800 disabled:cursor-not-allowed disabled:text-gray-300 dark:text-gray-400 dark:hover:text-gray-200 dark:disabled:text-gray-700"
      >
        Next
      </button>
    </div>
  )
}

const SENTIMENT_STYLES: Record<TickerNewsSummary['sentiment'], string> = {
  bullish: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400',
  bearish: 'bg-rose-50 text-rose-700 dark:bg-rose-900/30 dark:text-rose-400',
  neutral: 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-300',
}

function DigestCallout({ summary }: { summary: TickerNewsSummary }) {
  return (
    <div className="mt-3 rounded-xl bg-indigo-50/60 p-3 dark:bg-indigo-900/10">
      <div className="flex items-center gap-2">
        <span className="text-xs font-semibold text-gray-900 dark:text-gray-100">
          AI Daily Digest
        </span>
        <span
          className={`rounded-full px-2 py-0.5 text-[10px] font-medium capitalize ${SENTIMENT_STYLES[summary.sentiment]}`}
        >
          {summary.sentiment}
        </span>
      </div>
      <p className="mt-1.5 text-xs text-gray-600 dark:text-gray-300">{summary.dailySummary}</p>
      {summary.potentialImpact && (
        <p className="mt-1.5 text-[11px] text-gray-500 dark:text-gray-400">
          <span className="font-medium">Potential impact:</span> {summary.potentialImpact}
        </p>
      )}
    </div>
  )
}

// Shown whenever the digest fetch fails or turns up nothing (no cached
// articles yet to summarize, or the AI provider errored) — kept visually
// distinct (dashed, muted) from DigestCallout so it reads as "nothing here
// right now" rather than as actual AI-generated content.
function DigestUnavailable({ error }: { error: string | null }) {
  return (
    <div className="mt-3 rounded-xl border border-dashed border-gray-200 p-3 dark:border-gray-700">
      <span className="text-xs font-semibold text-gray-400 dark:text-gray-500">AI Daily Digest</span>
      <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">
        {error ?? 'No AI summary available for this ticker yet.'}
      </p>
    </div>
  )
}

function ArticleCard({ article }: { article: NewsArticle }) {
  return (
    <a
      href={article.url}
      target="_blank"
      rel="noreferrer"
      className="flex gap-3 rounded-lg px-1 py-3 transition first:pt-2 hover:bg-gray-50 dark:hover:bg-gray-800/60"
    >
      <div className="h-14 w-14 shrink-0 overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800">
        {article.imageUrl ? (
          <img
            src={article.imageUrl}
            alt=""
            className="h-full w-full object-cover"
            onError={(e) => {
              e.currentTarget.style.display = 'none'
            }}
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center text-gray-300 dark:text-gray-600">
            <NewsIcon className="h-5 w-5" />
          </div>
        )}
      </div>
      <div className="min-w-0 flex-1">
        <p className="text-[11px] text-gray-400 dark:text-gray-500">
          {formatRelativeTime(article.publishedAt)}
        </p>
        <p className="mt-0.5 text-sm font-medium text-gray-900 dark:text-gray-100">
          {article.headline}
        </p>
        {article.summary && (
          <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">{article.summary}</p>
        )}
        <p className="mt-0.5 text-xs text-gray-400 dark:text-gray-500">{article.source}</p>
      </div>
    </a>
  )
}

function ArticleSkeleton() {
  return (
    <div className="flex gap-3 px-1 py-3">
      <div className="h-14 w-14 shrink-0 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-800" />
      <div className="flex-1 space-y-2 py-1">
        <div className="h-2.5 w-16 animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
        <div className="h-3.5 w-full animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
        <div className="h-3.5 w-2/3 animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
      </div>
    </div>
  )
}

function NewsIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} xmlns="http://www.w3.org/2000/svg">
      <rect x="3" y="4" width="18" height="16" rx="2" stroke="currentColor" strokeWidth="1.8" />
      <path d="M7 8H14M7 12H17M7 16H12" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  )
}
