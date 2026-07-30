import { api } from '../../lib/api'
import type { NewsArticle, NewsCategory } from './types'

export function getGeneralNews(category: NewsCategory, limit = 20): Promise<NewsArticle[]> {
  return api.get<NewsArticle[]>(`/api/v1/news/general?category=${category}&limit=${limit}`)
}
