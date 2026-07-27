import { api } from '../../../lib/api';
import { Podcast } from './types';

interface ApiResponse<T> {
  data: T;
}

// Raw shapes come straight from the Listen Notes API, hence the snake_case.
interface RawPodcastResult {
  id: string;
  title_original?: string;
  title?: string;
  publisher_original?: string;
  publisher?: string;
  description_original?: string;
  description?: string;
  image?: string;
  thumbnail?: string;
  website?: string;
  listennotes_url?: string;
  total_episodes?: number;
  explicit_content?: boolean;
}

interface RawSearchResponse {
  results?: RawPodcastResult[];
}

const BASE_PATH = '/api/v1/learning/resources/podcast';

function mapPodcast(raw: RawPodcastResult): Podcast {
  return {
    id: raw.id,
    title: raw.title_original ?? raw.title ?? 'Untitled podcast',
    publisher: raw.publisher_original ?? raw.publisher ?? 'Unknown publisher',
    description: raw.description_original ?? raw.description ?? '',
    image: raw.image ?? raw.thumbnail ?? '',
    website: raw.website ?? '',
    listennotesUrl: raw.listennotes_url ?? '',
    totalEpisodes: raw.total_episodes,
    explicitContent: raw.explicit_content ?? false,
  };
}

async function searchPodcasts(query: string): Promise<Podcast[]> {
  const res = await api.get<ApiResponse<RawSearchResponse>>(
    `${BASE_PATH}?q=${encodeURIComponent(query)}`
  );
  const results = res.data.results ?? [];
  return results.map(mapPodcast);
}

// Fetches podcasts matching the "Finance" and "Business" genres and merges
// the two result sets, removing duplicates by podcast id.
export async function fetchFinancePodcasts(): Promise<Podcast[]> {
  const [financeResults, businessResults] = await Promise.all([
    searchPodcasts('Finance'),
    searchPodcasts('Business'),
  ]);

  const merged = new Map<string, Podcast>();
  for (const podcast of [...financeResults, ...businessResults]) {
    if (!merged.has(podcast.id)) {
      merged.set(podcast.id, podcast);
    }
  }
  return Array.from(merged.values());
}
