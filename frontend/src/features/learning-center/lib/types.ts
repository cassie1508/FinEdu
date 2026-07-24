export interface Flashcard {
  id: string;
  title: string;
  category: string;
  definition: string;
  example: string;
  whyItMatters: string;
  commonMisconceptions: string;
  reviewCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface LearningResource {
  id: number;
  title: string;
  category: string;
  summary: string;
  source: string;
  imageUrl: string;
  related: string[];
  publishedAt: string;
  url: string;
}

export type ResourceType = 'all' | 'video' | 'article' | 'podcast';

export interface ConversationMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}

export interface LearningProgress {
  flashcardsReviewed: number;
  flashcardsTotal: number;
  articlesRead: number;
  videosWatched: number;
  podcastsListened: number;
  bookmarksCount: number;
  favoritesCount: number;
}
