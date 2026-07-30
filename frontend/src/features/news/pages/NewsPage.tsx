import { NewsFeed } from '../components/NewsFeed'

export function NewsPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">News</h1>
      <p className="mt-2 text-gray-600 dark:text-gray-400">
        Stay informed with the latest financial news across the market, earnings, and mergers.
      </p>

      <div className="mt-6">
        <NewsFeed />
      </div>
    </div>
  )
}
