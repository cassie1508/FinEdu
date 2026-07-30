import { useState } from 'react'
import { useGeneralNews } from '../hooks/useGeneralNews'
import { NEWS_CATEGORIES, type NewsArticle, type NewsCategory } from '../types'
import { formatRelativeTime } from '../lib/format'

// Full-page general market news feed — not tied to any particular ticker.
// Mirrors Finnhub's own category taxonomy (market/earnings/mergers; see
// generalNewsCategories in news.go), defaulting to the top/general feed.
export function NewsFeed() {
  const [category, setCategory] = useState<NewsCategory>('market')
  const { articles, isLoading, error } = useGeneralNews(category)
  const categoryLabel = NEWS_CATEGORIES.find((c) => c.value === category)?.label ?? ''

  return (
    <div>
      <div className="flex items-center gap-1 border-b border-gray-200 pb-3 dark:border-gray-800">
        {NEWS_CATEGORIES.map((c) => (
          <button
            key={c.value}
            type="button"
            onClick={() => setCategory(c.value)}
            className={`rounded-md px-3 py-1.5 text-sm font-medium transition ${
              c.value === category
                ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400'
                : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
            }`}
          >
            {c.label}
          </button>
        ))}
      </div>

      {!isLoading && error && (
        <div className="flex flex-col items-center gap-1 py-16 text-center">
          <p className="text-sm font-medium text-red-600 dark:text-red-400">Couldn't load news</p>
          <p className="text-xs text-gray-500 dark:text-gray-400">{error}</p>
        </div>
      )}

      {!isLoading && !error && articles.length === 0 && (
        <div className="py-16 text-center text-sm text-gray-500 dark:text-gray-400">
          No news available right now.
        </div>
      )}

      <div className="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {isLoading && Array.from({ length: 6 }).map((_, i) => <ArticleCardSkeleton key={i} />)}

        {!isLoading &&
          !error &&
          articles.map((article) => (
            <ArticleCard key={article.id} article={article} categoryLabel={categoryLabel} />
          ))}
      </div>
    </div>
  )
}

function ArticleCard({ article, categoryLabel }: { article: NewsArticle; categoryLabel: string }) {
  return (
    <a
      href={article.url}
      target="_blank"
      rel="noreferrer"
      className="flex flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm transition hover:shadow-md dark:border-gray-800 dark:bg-gray-900"
    >
      <div className="h-36 w-full shrink-0 overflow-hidden bg-gray-100 dark:bg-gray-800">
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
            <NewsIcon className="h-8 w-8" />
          </div>
        )}
      </div>
      <div className="flex flex-1 flex-col p-4">
        <div className="flex items-center gap-1.5 text-[11px]">
          <span className="font-medium text-indigo-600 dark:text-indigo-400">{categoryLabel}</span>
          <span className="text-gray-300 dark:text-gray-600">•</span>
          <span className="text-gray-400 dark:text-gray-500">{formatRelativeTime(article.publishedAt)}</span>
        </div>
        <p className="mt-1 line-clamp-2 text-sm font-medium text-gray-900 dark:text-gray-100">
          {article.headline}
        </p>
        {article.summary && (
          <p className="mt-1 line-clamp-2 text-xs text-gray-500 dark:text-gray-400">{article.summary}</p>
        )}
        <p className="mt-2 text-xs text-gray-400 dark:text-gray-500">{article.source}</p>
      </div>
    </a>
  )
}

function ArticleCardSkeleton() {
  return (
    <div className="overflow-hidden rounded-2xl border border-gray-200 dark:border-gray-800">
      <div className="h-36 w-full animate-pulse bg-gray-100 dark:bg-gray-800" />
      <div className="space-y-2 p-4">
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
