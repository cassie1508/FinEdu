import { useEffect, useState } from 'react';
import { RotateCcw } from 'lucide-react';
import { colors } from '../lib/colors';
import { LearningResource } from '../lib/types';
import { fetchLearningResources } from '../lib/resourcesApi';

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

  const loadResources = () => {
    setIsLoading(true);
    setError(null);
    fetchLearningResources()
      .then(setResources)
      .catch(() => setError('Failed to load learning resources.'))
      .finally(() => setIsLoading(false));
  };

  useEffect(() => {
    loadResources();
  }, []);

  const tabs = ['All', ...Array.from(new Set(resources.map(resource => resource.category))).sort()];

  const filteredResources =
    activeTab === 'All'
      ? resources
      : resources.filter(resource => resource.category === activeTab);

  return (
    <div
      className="flex flex-col h-full overflow-hidden rounded-xl border"
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
            onClick={loadResources}
            disabled={isLoading}
            className="p-2 rounded-lg transition-all hover:shadow-md flex-shrink-0 disabled:opacity-60"
            style={{ color: colors.primary }}
            aria-label="Refresh resources"
          >
            <RotateCcw size={20} className={isLoading ? 'animate-spin' : undefined} />
          </button>
        </div>

        {/* Tabs */}
        <div className="flex gap-4 border-b overflow-x-auto pb-0" style={{ borderColor: colors.border }}>
          {tabs.map(tab => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
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

      {/* Scrollable Content */}
      <div className="flex-1 overflow-y-auto p-8 space-y-4 scrollbar-custom">
        {error && <p style={{ color: colors.accent }}>{error}</p>}
        {!error && isLoading && <p style={{ color: colors.accent }}>Loading resources...</p>}
        {!error && !isLoading && filteredResources.length === 0 && (
          <p style={{ color: colors.accent }}>No resources in this category yet.</p>
        )}
        {filteredResources.map(resource => {
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
      </div>
    </div>
  );
}

