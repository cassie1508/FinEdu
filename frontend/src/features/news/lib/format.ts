// News timestamps are ISO strings with real times, so this needs no
// timezone pinning — Date parses the offset that's already in the string.
export function formatRelativeTime(publishedAt: string): string {
  const diffMs = Date.now() - new Date(publishedAt).getTime()
  const minutes = Math.round(diffMs / 60_000)

  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`

  const hours = Math.round(minutes / 60)
  if (hours < 24) return `${hours}h ago`

  const days = Math.round(hours / 24)
  if (days < 7) return `${days}d ago`

  return new Date(publishedAt).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}
