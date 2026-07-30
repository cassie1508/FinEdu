// Matches the backend's generalNewsCategories (news.go) — "market" is the
// catch-all/top-news feed, not a literal "Market" category.
export type NewsCategory = 'market' | 'earnings' | 'mergers'

export const NEWS_CATEGORIES: { value: NewsCategory; label: string }[] = [
  { value: 'market', label: 'Top News' },
  { value: 'earnings', label: 'Earnings' },
  { value: 'mergers', label: 'Mergers' },
]

export interface NewsArticle {
  id: string
  headline: string
  summary: string
  source: string
  url: string
  publishedAt: string
  imageUrl: string
}
