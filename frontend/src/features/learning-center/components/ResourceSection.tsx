import { useState } from 'react';
import { RotateCcw } from 'lucide-react';
import { colors } from '../lib/colors';

interface ResourceData {
  badge: string;
  title: string;
  publisher: string;
}

const resourcesData: ResourceData[] = [
  {
    badge: 'Article',
    title: 'South Korean chip giant SK Hynix raises $26.5B in US share sale',
    publisher: 'BBC Business',
  },
  {
    badge: 'Video',
    title: 'This Stock Will Go to $500!!',
    publisher: 'Financial Education',
  },
  {
    badge: 'Article',
    title: "Pressure builds on Europe's biggest port to be greener",
    publisher: 'Reuters',
  },
  {
    badge: 'Video',
    title: 'Microsoft is doing a hard reset of its Xbox gaming division',
    publisher: 'TechNews Daily',
  },
  {
    badge: 'Podcast',
    title: 'Understanding Cryptocurrency Markets in 2024',
    publisher: 'Crypto Insights Daily',
  },
  {
    badge: 'Article',
    title: 'Federal Reserve Signals Potential Rate Cuts Next Quarter',
    publisher: 'Financial Times',
  },
];

const tabs = ['All', 'Tutorials', 'Articles', 'Podcasts'];

export function ResourceSection() {
  const [activeTab, setActiveTab] = useState('All');

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
              Tutorials, articles, podcasts from trusted sources.
            </p>
          </div>
          <button
            className="p-2 rounded-lg transition-all hover:shadow-md flex-shrink-0"
            style={{ color: colors.primary }}
            aria-label="Refresh resources"
          >
            <RotateCcw size={20} />
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
              {tab}
            </button>
          ))}
        </div>
      </div>

      {/* Scrollable Content */}
      <div className="flex-1 overflow-y-auto p-8 space-y-4 scrollbar-custom">
        {resourcesData.map((resource, idx) => (
          <div
            key={idx}
            className="flex gap-4 p-4 rounded-lg border transition-all hover:shadow-md cursor-pointer"
            style={{
              backgroundColor: colors.background,
              borderColor: colors.border,
            }}
          >
            <div
              className="w-24 h-24 rounded-lg flex-shrink-0"
              style={{ backgroundColor: colors.secondary }}
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
                  {resource.badge}
                </span>
              </div>
              <h4
                className="font-medium text-sm mb-1 line-clamp-2"
                style={{ color: colors.emphasis }}
              >
                {resource.title}
              </h4>
              <p className="text-xs" style={{ color: colors.accent }}>
                {resource.publisher}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
