import React from 'react';
import { BookOpen, Bookmark, Heart, TrendingUp } from 'lucide-react';
import { colors } from '../lib/colors';

interface RightSidebarProps {
  progress: {
    flashcardsReviewed: number;
    flashcardsTotal: number;
    articlesRead: number;
    videosWatched: number;
    podcastsListened: number;
    bookmarksCount: number;
    favoritesCount: number;
  };
}

interface StatCardProps {
  label: string;
  value: number;
  icon?: React.ReactNode;
}

function StatCard({ label, value, icon }: StatCardProps) {
  return (
    <div
      className="p-4 rounded-lg border"
      style={{
        backgroundColor: colors.surface,
        borderColor: colors.border,
      }}
    >
      <div className="flex items-start justify-between">
        <div>
          <p className="text-xs mb-1" style={{ color: colors.accent }}>
            {label}
          </p>
          <p
            className="text-2xl font-bold"
            style={{ color: colors.emphasis }}
          >
            {value}
          </p>
        </div>
        {icon && (
          <div style={{ color: colors.primary }}>
            {icon}
          </div>
        )}
      </div>
    </div>
  );
}

export function RightSidebar({ progress }: RightSidebarProps) {
  const progressPercentage = (progress.flashcardsReviewed / progress.flashcardsTotal) * 100;

  const stats = [
    { label: 'Cards', value: progress.flashcardsTotal, icon: <BookOpen size={20} /> },
    { label: 'Articles', value: progress.articlesRead, icon: null },
    { label: 'Videos', value: progress.videosWatched, icon: null },
    { label: 'Podcasts', value: progress.podcastsListened, icon: null },
    { label: 'Bookmarks', value: progress.bookmarksCount, icon: <Bookmark size={20} /> },
    { label: 'Favorites', value: progress.favoritesCount, icon: <Heart size={20} /> },
  ];

  const recommendations = [
    { title: 'Investing Foundations', lessons: 6, category: 'Investing' },
    { title: 'Budgeting & Saving', lessons: 5, category: 'Budgeting' },
    { title: 'Retirement Planning', lessons: 8, category: 'Retirement' },
  ];

  return (
    <div
      className="flex flex-col h-full"
      style={{
        background: 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
        borderColor: colors.border,
      }}
    >
      <div className="p-6 space-y-6 flex-1">
        {/* Learning Progress Card */}
        <div
          className="rounded-xl border p-6"
          style={{
            backgroundColor: colors.surface,
            borderColor: colors.border,
          }}
        >
          <div className="flex flex-col items-center mb-6">
            {/* Circular Progress Indicator */}
            <div className="relative w-24 h-24 mb-4">
              <svg className="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
                {/* Background circle */}
                <circle
                  cx="50"
                  cy="50"
                  r="45"
                  fill="none"
                  stroke={colors.border}
                  strokeWidth="4"
                />
                {/* Progress circle */}
                <circle
                  cx="50"
                  cy="50"
                  r="45"
                  fill="none"
                  stroke={colors.primary}
                  strokeWidth="4"
                  strokeDasharray={`${(progressPercentage / 100) * 283} 283`}
                  strokeLinecap="round"
                />
              </svg>
              <div className="absolute inset-0 flex items-center justify-center">
                <span
                  className="text-sm font-bold"
                  style={{ color: colors.emphasis }}
                >
                  {Math.round(progressPercentage)}%
                </span>
              </div>
            </div>
          </div>

          <h3 className="text-center font-semibold mb-1" style={{ color: colors.emphasis }}>
            Keep it up
          </h3>
          <p className="text-center text-xs mb-4" style={{ color: colors.accent }}>
            {progress.flashcardsReviewed} of {progress.flashcardsTotal} flashcards reviewed.
          </p>
          <p className="text-center text-xs" style={{ color: colors.text.secondary }}>
            Consistency beats intensity.
          </p>
        </div>

        {/* Statistics Cards Grid */}
        <div>
          <h3 className="text-sm font-semibold mb-3" style={{ color: colors.emphasis }}>
            Your Learning Stats
          </h3>
          <div className="grid grid-cols-2 gap-3">
            {stats.map((stat, idx) => (
              <StatCard
                key={idx}
                label={stat.label}
                value={stat.value}
                icon={stat.icon}
              />
            ))}
          </div>
        </div>

        {/* Concepts You Explore Most */}
        <div
          className="rounded-lg border p-4"
          style={{
            backgroundColor: colors.surface,
            borderColor: colors.border,
          }}
        >
          <h3 className="text-sm font-semibold mb-3" style={{ color: colors.emphasis }}>
            Concepts You Explore Most
          </h3>
          <div className="flex flex-col items-center justify-center py-8">
            <TrendingUp size={32} style={{ color: colors.accent, marginBottom: '8px' }} />
            <p className="text-xs text-center" style={{ color: colors.accent }}>
              Your most-studied concepts will appear here.
            </p>
          </div>
        </div>

        {/* Recommended For You */}
        <div>
          <h3 className="text-sm font-semibold mb-3" style={{ color: colors.emphasis }}>
            Recommended For You
          </h3>
          <div className="space-y-2">
            {recommendations.map((rec, idx) => (
              <div
                key={idx}
                className="p-3 rounded-lg border cursor-pointer transition-all hover:shadow-md"
                style={{
                  backgroundColor: colors.surface,
                  borderColor: colors.border,
                }}
              >
                <div className="flex items-start gap-2">
                  <BookOpen size={16} style={{ color: colors.primary, marginTop: '2px' }} />
                  <div className="flex-1 min-w-0">
                    <h4
                      className="text-xs font-medium"
                      style={{ color: colors.emphasis }}
                    >
                      {rec.title}
                    </h4>
                    <p className="text-xs mt-1" style={{ color: colors.accent }}>
                      {rec.lessons} lessons • {rec.category}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
