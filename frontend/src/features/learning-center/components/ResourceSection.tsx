import { useEffect, useState } from 'react';
import { RotateCcw } from 'lucide-react';
import { colors } from '../lib/colors';
import { LearningResource, Podcast } from '../lib/types';
import { fetchLearningResources } from '../lib/resourcesApi';
import { fetchFinancePodcasts } from '../lib/podcastsApi';
import { Pagination } from './Pagination';

const PODCAST_TAB = 'Podcast';
const RESOURCES_PAGE_SIZE = 6;

function isSafeHttpUrl(url: string): boolean {
  try {
    const parsed = new URL(url);
    return parsed.protocol === 'http:' || parsed.protocol === 'https:';
  } catch {
    return false;
  }
}

export function ResourceSection() {
  const [activeTab, setActiveTab] = useState('All');
  const [resources, setResources] = useState<LearningResource[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [podcasts, setPodcasts] = useState<Podcast[]>([]);
  const [isLoadingPodcasts, setIsLoadingPodcasts] = useState(false);
  const [podcastError, setPodcastError] = useState<string | null>(null);
  const [podcastsLoaded, setPodcastsLoaded] = useState(false);

  const [resourcesPage, setResourcesPage] = useState(1);
  const [podcastsPage, setPodcastsPage] = useState(1);

  const handleTabChange = (tab: string) => {
    setActiveTab(tab);
    setResourcesPage(1);
    setPodcastsPage(1);
  };

  const loadResources = () => {
    setIsLoading(true);
    setError(null);
    fetchLearningResources()
      .then(setResources)
      .catch(() => setError('Failed to load learning resources.'))
      .finally(() => setIsLoading(false));
  };

  const loadPodcasts = () => {
    setIsLoadingPodcasts(true);
    setPodcastError(null);
    fetchFinancePodcasts()
      .then(result => {
        setPodcasts(result);
        setPodcastsLoaded(true);
      })
      .catch(() => setPodcastError('Failed to load podcasts.'))
      .finally(() => setIsLoadingPodcasts(false));
  };

  useEffect(() => {
    loadResources();
  }, []);

  useEffect(() => {
    if (activeTab === PODCAST_TAB && !podcastsLoaded && !isLoadingPodcasts) {
      loadPodcasts();
    }
  }, [activeTab, podcastsLoaded, isLoadingPodcasts]);

  const tabs = [
    'All',
    ...Array.from(new Set(resources.map(resource => resource.category))).sort(),
    PODCAST_TAB,
  ];

  const filteredResources =
    activeTab === 'All'
      ? resources
      : resources.filter(resource => resource.category === activeTab);

  const resourcesTotalPages = Math.max(1, Math.ceil(filteredResources.length / RESOURCES_PAGE_SIZE));
  const safeResourcesPage = Math.min(resourcesPage, resourcesTotalPages);
  const pagedResources = filteredResources.slice(
    (safeResourcesPage - 1) * RESOURCES_PAGE_SIZE,
    safeResourcesPage * RESOURCES_PAGE_SIZE,
  );

  const podcastsTotalPages = Math.max(1, Math.ceil(podcasts.length / RESOURCES_PAGE_SIZE));
  const safePodcastsPage = Math.min(podcastsPage, podcastsTotalPages);
  const pagedPodcasts = podcasts.slice(
    (safePodcastsPage - 1) * RESOURCES_PAGE_SIZE,
    safePodcastsPage * RESOURCES_PAGE_SIZE,
  );

  return (
    <div
      className="flex flex-col rounded-xl border"
      style={{
        backgroundColor: colors.surface,
        borderColor: colors.border,
      }}
    >
      {/* Header - Fixed */}
      <div className="flex-shrink-0 p-8 border-b" style={{ borderColor: colors.border }}>
        <div className="flex items-start justify-between mb-6">
          <div>
            <h2
              className="text-2xl font-serif font-bold mb-1"
              style={{ color: colors.emphasis }}
            >
              Learning Resources
            </h2>
            <p style={{ color: colors.accent }}>
              Finance and market news from trusted sources.
            </p>
          </div>
          <button
            onClick={activeTab === PODCAST_TAB ? loadPodcasts : loadResources}
            disabled={activeTab === PODCAST_TAB ? isLoadingPodcasts : isLoading}
            className="p-2 rounded-lg transition-all hover:shadow-md flex-shrink-0 disabled:opacity-60"
            style={{ color: colors.primary }}
            aria-label="Refresh resources"
          >
            <RotateCcw
              size={20}
              className={(activeTab === PODCAST_TAB ? isLoadingPodcasts : isLoading) ? 'animate-spin' : undefined}
            />
          </button>
        </div>

        {/* Tabs */}
        <div className="flex gap-4 border-b overflow-x-auto pb-0" style={{ borderColor: colors.border }}>
          {tabs.map(tab => (
            <button
              key={tab}
              onClick={() => handleTabChange(tab)}
              className="px-4 py-2 border-b-2 font-medium text-sm transition-colors whitespace-nowrap"
              style={{
                borderColor: tab === activeTab ? colors.primary : 'transparent',
                color: tab === activeTab ? colors.primary : colors.accent,
              }}
            >
              {tab[0] === tab[0].toUpperCase() ? tab : tab.charAt(0).toUpperCase() + tab.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="p-8 space-y-4">
        {activeTab === PODCAST_TAB ? (
          <>
            {podcastError && <p style={{ color: colors.accent }}>{podcastError}</p>}
            {!podcastError && isLoadingPodcasts && <p style={{ color: colors.accent }}>Loading podcasts...</p>}
            {!podcastError && !isLoadingPodcasts && podcasts.length === 0 && (
              <p style={{ color: colors.accent }}>No podcasts found.</p>
            )}
            {pagedPodcasts.map(podcast => {
              const safeUrl = isSafeHttpUrl(podcast.listennotesUrl) ? podcast.listennotesUrl : undefined;
              const safeImageUrl = podcast.image && isSafeHttpUrl(podcast.image) ? podcast.image : undefined;

              return (
                <a
                  key={podcast.id}
                  href={safeUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex gap-4 p-4 rounded-lg border transition-all hover:shadow-md cursor-pointer"
                  style={{
                    backgroundColor: colors.background,
                    borderColor: colors.border,
                  }}
                >
                  <div
                    className="w-24 h-24 rounded-lg flex-shrink-0"
                    style={{
                      backgroundColor: colors.secondary,
                      backgroundImage: safeImageUrl ? `url(${safeImageUrl})` : undefined,
                      backgroundSize: 'cover',
                      backgroundPosition: 'center',
                    }}
                  />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-2">
                      <span
                        className="text-xs font-semibold px-2 py-1 rounded"
                        style={{
                          background: 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
                          color: colors.emphasis,
                        }}
                      >
                        Podcast
                      </span>
                    </div>
                    <h4
                      className="font-medium text-sm mb-1 line-clamp-2"
                      style={{ color: colors.emphasis }}
                    >
                      {podcast.title}
                    </h4>
                    <p className="text-xs" style={{ color: colors.accent }}>
                      {podcast.publisher}
                    </p>
                  </div>
                </a>
              );
            })}
            <Pagination
              currentPage={safePodcastsPage}
              totalPages={podcastsTotalPages}
              onPageChange={setPodcastsPage}
            />
          </>
        ) : (
          <>
            {error && <p style={{ color: colors.accent }}>{error}</p>}
            {!error && isLoading && <p style={{ color: colors.accent }}>Loading resources...</p>}
            {!error && !isLoading && filteredResources.length === 0 && (
              <p style={{ color: colors.accent }}>No resources in this category yet.</p>
            )}
            {pagedResources.map(resource => {
              const safeUrl = isSafeHttpUrl(resource.url) ? resource.url : undefined;
              const safeImageUrl = resource.imageUrl && isSafeHttpUrl(resource.imageUrl) ? resource.imageUrl : undefined;

              return (
                <a
                  key={resource.id}
                  href={safeUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex gap-4 p-4 rounded-lg border transition-all hover:shadow-md cursor-pointer"
                  style={{
                    backgroundColor: colors.background,
                    borderColor: colors.border,
                  }}
                >
                  <div
                    className="w-24 h-24 rounded-lg flex-shrink-0"
                    style={{
                      backgroundColor: colors.secondary,
                      backgroundImage: safeImageUrl ? `url(${safeImageUrl})` : undefined,
                      backgroundSize: 'cover',
                      backgroundPosition: 'center',
                    }}
                  />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-2">
                      <span
                        className="text-xs font-semibold px-2 py-1 rounded"
                        style={{
                          background: 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
                          color: colors.emphasis,
                        }}
                      >
                        {resource.category}
                      </span>
                    </div>
                    <h4
                      className="font-medium text-sm mb-1 line-clamp-2"
                      style={{ color: colors.emphasis }}
                    >
                      {resource.title}
                    </h4>
                    <p className="text-xs" style={{ color: colors.accent }}>
                      {resource.source}
                    </p>
                  </div>
                </a>
              );
            })}
            <Pagination
              currentPage={safeResourcesPage}
              totalPages={resourcesTotalPages}
              onPageChange={setResourcesPage}
            />
          </>
        )}
      </div>
    </div>
  );
}

