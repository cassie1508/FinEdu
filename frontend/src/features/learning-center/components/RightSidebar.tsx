import React, { useState, useEffect } from 'react';
import { BookOpen, ExternalLink } from 'lucide-react';
import { colors } from '../lib/colors';
import { api } from '../../../lib/api';

interface RightSidebarProps {
  progress: {
    flashcardsTotal: number;
  };
}

interface RecommendationResource {

  title: string;
  link: string;
}

interface RecommendationItem {
  flashcard_title: string;
  resources: RecommendationResource[];
}

interface RecommendationsResponse {
  success: boolean;
  user_id: string;
  titles_count: number;
  recommendations: RecommendationItem[];
}

export function RightSidebar({ progress }: RightSidebarProps) {
  const [recommendations, setRecommendations] = useState<RecommendationItem[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRecommendations = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await api.post<RecommendationsResponse>('/api/v1/recommendations', {
          max_results_per_title: 3,
        });
        if (response.success) {
          setRecommendations(response.recommendations || []);
        } else {
          setError('Failed to load recommendations');
        }
      } catch (err: any) {
        setError('Error loading recommendations');
      } finally {
        setIsLoading(false);
      }
    };

    fetchRecommendations();
  }, []);

  return (
    <div
      className="flex flex-col h-full p-6"
      style={{
        background: 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
      }}
    >

      {/* Recommended For You */}
      <div className="flex-1 flex flex-col min-h-0">
        <h3 className="text-sm font-semibold mb-3" style={{ color: colors.emphasis }}>
          Recommended For You
        </h3>

        {isLoading && (
          <div className="text-center py-4">
            <p className="text-xs" style={{ color: colors.accent }}>
              Loading recommendations...
            </p>
          </div>
        )}

        {error && (
          <div className="text-center py-4">
            <p className="text-xs" style={{ color: '#ef4444' }}>
              {error}
            </p>
          </div>
        )}

        {!isLoading && !error && recommendations.length === 0 && (
          <div className="text-center py-4">
            <p className="text-xs" style={{ color: colors.accent }}>
              No recommendations yet. Upload flashcards to get started!
            </p>
          </div>
        )}

        {!isLoading && !error && recommendations.length > 0 && (
          <div className="flex-1 space-y-3 overflow-y-auto">
            {recommendations.map((rec, idx) => (
              <div
                key={idx}
                className="rounded-lg border p-3"
                style={{
                  backgroundColor: colors.surface,
                  borderColor: colors.border,
                }}
              >
                <h4
                  className="text-xs font-semibold mb-2"
                  style={{ color: colors.emphasis }}
                >
                  {rec.flashcard_title}
                </h4>
                <div className="space-y-1">
                  {rec.resources && rec.resources.length > 0 ? (
                    rec.resources.map((resource, ridx) => (
                      <a
                        key={ridx}
                        href={resource.link}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-start gap-2 p-2 rounded transition-colors hover:opacity-80"
                        style={{
                          backgroundColor: colors.background,
                        }}
                      >
                        <ExternalLink
                          size={12}
                          style={{ color: colors.primary, marginTop: '2px', flexShrink: 0 }}
                        />
                        <span
                          className="text-xs line-clamp-2"
                          style={{ color: colors.text.primary }}
                        >
                          {resource.title}
                        </span>
                      </a>
                    ))
                  ) : (
                    <p className="text-xs" style={{ color: colors.accent }}>
                      No resources found
                    </p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
